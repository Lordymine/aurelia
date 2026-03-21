package bridge

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
)

// Bridge spawns the TypeScript bridge process and communicates via stdin/stdout.
// Each Execute call spawns a fresh process — the bridge is not long-lived from
// Go's perspective in v1.
type Bridge struct {
	bridgeDir string // directory containing bridge/index.ts

	// command and args override the default "npx tsx index.ts" for testing.
	command string
	args    []string
}

// New creates a Bridge that will spawn processes in bridgeDir.
func New(bridgeDir string) *Bridge {
	return &Bridge{
		bridgeDir: bridgeDir,
		command:   "npx",
		args:      []string{"tsx", "index.ts"},
	}
}

// Execute sends a request to a freshly spawned Bridge process and streams
// parsed NDJSON events back through the returned channel. The channel is
// closed when the process exits or the context is cancelled.
func (b *Bridge) Execute(ctx context.Context, req Request) (<-chan Event, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("bridge: marshal request: %w", err)
	}

	cmd := exec.CommandContext(ctx, b.command, b.args...)
	cmd.Dir = b.bridgeDir

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("bridge: stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("bridge: stdout pipe: %w", err)
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("bridge: start process: %w", err)
	}

	// Write request and close stdin so the bridge knows no more input is coming.
	if _, err := stdin.Write(append(payload, '\n')); err != nil {
		_ = cmd.Process.Kill()
		return nil, fmt.Errorf("bridge: write request: %w", err)
	}
	stdin.Close()

	ch := make(chan Event, 16)

	go func() {
		defer close(ch)
		defer cmd.Wait() //nolint:errcheck // exit code checked via events

		scanner := bufio.NewScanner(stdout)
		// Allow large lines (some tool results can be big).
		scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

		for scanner.Scan() {
			line := scanner.Bytes()
			if len(line) == 0 {
				continue
			}

			var ev Event
			if err := json.Unmarshal(line, &ev); err != nil {
				// Emit a synthetic error event for unparseable lines.
				select {
				case ch <- Event{Type: "error", Message: fmt.Sprintf("bridge: unmarshal event: %v (line: %s)", err, line)}:
				case <-ctx.Done():
					return
				}
				continue
			}

			select {
			case ch <- ev:
			case <-ctx.Done():
				return
			}
		}
	}()

	return ch, nil
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
