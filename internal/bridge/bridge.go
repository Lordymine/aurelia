package bridge

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"sync/atomic"
)

// safeClose closes a channel, recovering from panic if already closed.
func safeClose(ch chan Event) {
	defer func() { recover() }()
	close(ch)
}

// Bridge manages a long-lived TypeScript bridge process and communicates via
// stdin/stdout using NDJSON. Multiple requests are multiplexed over a single
// process using request_id correlation.
type Bridge struct {
	bridgeDir string // directory containing bridge/index.ts

	// command and args override the default "npx tsx index.ts" for testing.
	command string
	args    []string

	mu      sync.Mutex // guards stdin writes and process lifecycle
	cmd     *exec.Cmd
	stdin   io.WriteCloser
	scanner *bufio.Scanner

	// pending maps request_id → channel for routing events.
	pending   map[string]chan Event
	pendingMu sync.Mutex

	started  bool
	stopping bool

	// reqCounter generates unique request IDs.
	reqCounter atomic.Uint64

	// done is closed when the reader goroutine exits.
	done chan struct{}
}

// New creates a Bridge that runs the given bundlePath with node.
// If bundlePath is empty, falls back to npx tsx index.ts in bridgeDir.
func New(bridgeDir string, bundlePath string) *Bridge {
	cmd := "node"
	args := []string{bundlePath}
	if bundlePath == "" {
		cmd = "npx"
		args = []string{"tsx", "index.ts"}
	}
	return &Bridge{
		bridgeDir: bridgeDir,
		command:   cmd,
		args:      args,
		pending:   make(map[string]chan Event),
	}
}

// Start launches the bridge process. Safe to call multiple times — no-op if
// already running.
func (b *Bridge) Start() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.startLocked()
}

func (b *Bridge) startLocked() error {
	if b.started {
		return nil
	}

	cmd := exec.Command(b.command, b.args...)
	cmd.Dir = b.bridgeDir

	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("bridge: stdin pipe: %w", err)
	}

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("bridge: stdout pipe: %w", err)
	}

	// Stderr goes to parent stderr for debugging.
	cmd.Stderr = nil // inherits parent stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("bridge: start process: %w", err)
	}

	b.cmd = cmd
	b.stdin = stdinPipe
	b.scanner = bufio.NewScanner(stdoutPipe)
	b.scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	b.started = true
	b.stopping = false
	b.done = make(chan struct{})

	go b.readLoop()

	return nil
}

// readLoop runs in a goroutine, reading stdout and routing events to pending
// request channels. When the process exits, all pending channels are closed.
func (b *Bridge) readLoop() {
	defer close(b.done)

	for b.scanner.Scan() {
		line := b.scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var ev Event
		if err := json.Unmarshal(line, &ev); err != nil {
			// Can't route without request_id, skip.
			continue
		}

		rid := ev.RequestID

		b.pendingMu.Lock()
		ch, ok := b.pending[rid]
		b.pendingMu.Unlock()

		if !ok {
			// No listener for this request_id — drop event.
			continue
		}

		// Non-blocking send — channel has buffer.
		select {
		case ch <- ev:
		default:
			// Buffer full — drop event rather than block reader.
		}

		if ev.IsTerminal() {
			b.pendingMu.Lock()
			delete(b.pending, rid)
			b.pendingMu.Unlock()
			safeClose(ch)
		}
	}

	// Process exited or stdout closed — close all pending channels.
	b.pendingMu.Lock()
	for rid, ch := range b.pending {
		safeClose(ch)
		delete(b.pending, rid)
	}
	b.pendingMu.Unlock()

	b.mu.Lock()
	b.started = false
	b.cmd = nil
	b.mu.Unlock()
}

// Stop kills the bridge process. Safe to call multiple times.
func (b *Bridge) Stop() {
	b.mu.Lock()
	if !b.started || b.stopping {
		b.mu.Unlock()
		return
	}
	b.stopping = true
	stdin := b.stdin
	cmd := b.cmd
	done := b.done
	b.mu.Unlock()

	// Close stdin — the TS bridge exits on stdin close.
	if stdin != nil {
		_ = stdin.Close()
	}

	// Wait for reader goroutine to finish (it will close all pending channels).
	if done != nil {
		<-done
	}

	// Ensure process is reaped.
	if cmd != nil {
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
		_ = cmd.Wait()
	}

	// Reset state so the bridge can be restarted.
	b.mu.Lock()
	b.started = false
	b.stopping = false
	b.cmd = nil
	b.stdin = nil
	b.scanner = nil
	b.pendingMu.Lock()
	b.pending = make(map[string]chan Event)
	b.pendingMu.Unlock()
	b.mu.Unlock()
}

// Execute sends a request to the long-lived Bridge process and returns a
// channel of events for that request. The process stays alive after the
// request completes.
func (b *Bridge) Execute(ctx context.Context, req Request) (<-chan Event, error) {
	b.mu.Lock()
	if !b.started {
		if err := b.startLocked(); err != nil {
			b.mu.Unlock()
			return nil, err
		}
	}
	b.mu.Unlock()

	// Assign request_id if not set.
	if req.RequestID == "" {
		req.RequestID = fmt.Sprintf("req-%d", b.reqCounter.Add(1))
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("bridge: marshal request: %w", err)
	}

	ch := make(chan Event, 16)

	b.pendingMu.Lock()
	b.pending[req.RequestID] = ch
	b.pendingMu.Unlock()

	// Write request to stdin (don't close stdin!).
	b.mu.Lock()
	if !b.started {
		b.mu.Unlock()
		b.pendingMu.Lock()
		delete(b.pending, req.RequestID)
		b.pendingMu.Unlock()
		safeClose(ch)
		return nil, fmt.Errorf("bridge: process died before write")
	}
	_, err = b.stdin.Write(append(payload, '\n'))
	b.mu.Unlock()

	if err != nil {
		b.pendingMu.Lock()
		delete(b.pending, req.RequestID)
		b.pendingMu.Unlock()
		safeClose(ch)
		return nil, fmt.Errorf("bridge: write request: %w", err)
	}

	// Wrap channel with context cancellation.
	out := make(chan Event, 16)
	go func() {
		defer close(out)
		for {
			select {
			case ev, ok := <-ch:
				if !ok {
					return
				}
				select {
				case out <- ev:
				case <-ctx.Done():
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return out, nil
}

// ExecuteSync sends a request and blocks until a terminal event (result or error)
// is received. It returns that event. Intermediate events are discarded.
func (b *Bridge) ExecuteSync(ctx context.Context, req Request) (*Event, error) {
	ch, err := b.Execute(ctx, req)
	if err != nil {
		return nil, err
	}

	var last *Event
	for ev := range ch {
		ev := ev
		last = &ev
		if ev.IsTerminal() {
			// Drain remaining events (shouldn't be any, but be safe).
			go func() {
				for range ch { //nolint:revive
				}
			}()
			return last, nil
		}
	}

	if last != nil {
		return last, nil
	}
	return nil, fmt.Errorf("bridge: process exited without producing any events")
}

// Ping verifies the bridge process can start and respond to a ping command.
func (b *Bridge) Ping(ctx context.Context) error {
	ev, err := b.ExecuteSync(ctx, Request{Command: "ping"})
	if err != nil {
		return fmt.Errorf("bridge: ping failed: %w", err)
	}
	if ev.Type == "error" {
		return fmt.Errorf("bridge: ping returned error: %s", ev.Message)
	}
	if ev.Type != "pong" {
		return fmt.Errorf("bridge: ping expected pong, got %q", ev.Type)
	}
	return nil
}
