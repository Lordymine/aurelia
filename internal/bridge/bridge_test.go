package bridge

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// longLivedMockJS returns JavaScript that acts as a long-lived bridge process:
// reads multiple lines from stdin, responds to each with events including request_id.
const longLivedMockJS = `
const readline = require('readline');
const rl = readline.createInterface({ input: process.stdin, terminal: false });

rl.on('line', (line) => {
    let req;
    try {
        req = JSON.parse(line);
    } catch(e) {
        process.stdout.write(JSON.stringify({event:"error",message:"invalid json"}) + "\n");
        return;
    }

    const rid = req.request_id || "";

    if (req.command === "ping") {
        process.stdout.write(JSON.stringify({event:"pong",request_id:rid}) + "\n");
    } else if (req.command === "query") {
        process.stdout.write(JSON.stringify({event:"system",request_id:rid,session_id:"test-123",tools:["Read"],model:"claude-3"}) + "\n");
        process.stdout.write(JSON.stringify({event:"assistant",request_id:rid,text:"hello world"}) + "\n");
        process.stdout.write(JSON.stringify({event:"result",request_id:rid,content:"done",cost_usd:0.01,session_id:"test-123",duration_ms:100,num_turns:1}) + "\n");
    } else {
        process.stdout.write(JSON.stringify({event:"error",request_id:rid,message:"unknown command: " + req.command}) + "\n");
    }
});

rl.on('close', () => {
    process.exit(0);
});
`

// newMockBridge creates a long-lived Bridge that uses `node mock.js`.
func newMockBridge(t *testing.T, dir string, jsBody string) *Bridge {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, "mock.js"), []byte(jsBody), 0644); err != nil {
		t.Fatal(err)
	}
	b := New(dir)
	b.command = "node"
	b.args = []string{"mock.js"}
	t.Cleanup(func() { b.Stop() })
	return b
}

func TestBridge_Execute_ParsesEvents(t *testing.T) {
	dir := t.TempDir()
	b := newMockBridge(t, dir, longLivedMockJS)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ch, err := b.Execute(ctx, Request{Command: "query", Prompt: "test"})
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	var events []Event
	for ev := range ch {
		events = append(events, ev)
	}

	if len(events) != 3 {
		t.Fatalf("expected 3 events, got %d: %+v", len(events), events)
	}

	// system event
	if events[0].Type != "system" {
		t.Errorf("event[0].Type = %q, want %q", events[0].Type, "system")
	}
	if events[0].SessionID != "test-123" {
		t.Errorf("event[0].SessionID = %q, want %q", events[0].SessionID, "test-123")
	}
	if len(events[0].Tools) != 1 || events[0].Tools[0] != "Read" {
		t.Errorf("event[0].Tools = %v, want [Read]", events[0].Tools)
	}
	if events[0].Model != "claude-3" {
		t.Errorf("event[0].Model = %q, want %q", events[0].Model, "claude-3")
	}

	// assistant event
	if events[1].Type != "assistant" {
		t.Errorf("event[1].Type = %q, want %q", events[1].Type, "assistant")
	}
	if events[1].Text != "hello world" {
		t.Errorf("event[1].Text = %q, want %q", events[1].Text, "hello world")
	}

	// result event
	if events[2].Type != "result" {
		t.Errorf("event[2].Type = %q, want %q", events[2].Type, "result")
	}
	if events[2].Content != "done" {
		t.Errorf("event[2].Content = %q, want %q", events[2].Content, "done")
	}
	if events[2].CostUSD != 0.01 {
		t.Errorf("event[2].CostUSD = %f, want %f", events[2].CostUSD, 0.01)
	}
	if events[2].NumTurns != 1 {
		t.Errorf("event[2].NumTurns = %d, want %d", events[2].NumTurns, 1)
	}
}

func TestBridge_ExecuteSync_ReturnsResult(t *testing.T) {
	dir := t.TempDir()
	b := newMockBridge(t, dir, longLivedMockJS)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ev, err := b.ExecuteSync(ctx, Request{Command: "query", Prompt: "test"})
	if err != nil {
		t.Fatalf("ExecuteSync() error: %v", err)
	}
	if ev.Type != "result" {
		t.Errorf("Type = %q, want %q", ev.Type, "result")
	}
	if ev.Content != "done" {
		t.Errorf("Content = %q, want %q", ev.Content, "done")
	}
	if ev.CostUSD != 0.01 {
		t.Errorf("CostUSD = %f, want %f", ev.CostUSD, 0.01)
	}
}

func TestBridge_Execute_ErrorEvent(t *testing.T) {
	dir := t.TempDir()

	errorMockJS := `
const readline = require('readline');
const rl = readline.createInterface({ input: process.stdin, terminal: false });
rl.on('line', (line) => {
    const req = JSON.parse(line);
    const rid = req.request_id || "";
    process.stdout.write(JSON.stringify({event:"error",request_id:rid,message:"something went wrong"}) + "\n");
});
rl.on('close', () => process.exit(0));
`
	b := newMockBridge(t, dir, errorMockJS)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ch, err := b.Execute(ctx, Request{Command: "query", Prompt: "test"})
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	var events []Event
	for ev := range ch {
		events = append(events, ev)
	}

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Type != "error" {
		t.Errorf("Type = %q, want %q", events[0].Type, "error")
	}
	if events[0].Message != "something went wrong" {
		t.Errorf("Message = %q, want %q", events[0].Message, "something went wrong")
	}
}

func TestBridge_Execute_ContextCancel(t *testing.T) {
	dir := t.TempDir()

	// Script that reads requests but never responds — simulates a hanging query.
	hangMockJS := `
const readline = require('readline');
const rl = readline.createInterface({ input: process.stdin, terminal: false });
rl.on('line', () => {
    // Intentionally don't respond — simulate hang.
});
rl.on('close', () => process.exit(0));
`
	b := newMockBridge(t, dir, hangMockJS)

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	ch, err := b.Execute(ctx, Request{Command: "query", Prompt: "test"})
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	select {
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for channel to close after context cancel")
	case _, ok := <-ch:
		if ok {
			for range ch {
			}
		}
	}
}

func TestBridge_Ping(t *testing.T) {
	dir := t.TempDir()
	b := newMockBridge(t, dir, longLivedMockJS)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := b.Ping(ctx); err != nil {
		t.Fatalf("Ping() error: %v", err)
	}
}

func TestBridge_Execute_ToolUseEvent(t *testing.T) {
	dir := t.TempDir()

	toolMockJS := `
const readline = require('readline');
const rl = readline.createInterface({ input: process.stdin, terminal: false });
rl.on('line', (line) => {
    const req = JSON.parse(line);
    const rid = req.request_id || "";
    process.stdout.write(JSON.stringify({event:"system",request_id:rid,session_id:"s1"}) + "\n");
    process.stdout.write(JSON.stringify({event:"tool_use",request_id:rid,name:"Read",input:{file_path:"/tmp/test.txt"}}) + "\n");
    process.stdout.write(JSON.stringify({event:"tool_result",request_id:rid,content:"file contents here"}) + "\n");
    process.stdout.write(JSON.stringify({event:"result",request_id:rid,content:"done"}) + "\n");
});
rl.on('close', () => process.exit(0));
`
	b := newMockBridge(t, dir, toolMockJS)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ch, err := b.Execute(ctx, Request{Command: "query", Prompt: "test"})
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	var events []Event
	for ev := range ch {
		events = append(events, ev)
	}

	if len(events) != 4 {
		t.Fatalf("expected 4 events, got %d", len(events))
	}
	if events[1].Type != "tool_use" {
		t.Errorf("event[1].Type = %q, want %q", events[1].Type, "tool_use")
	}
	if events[1].Name != "Read" {
		t.Errorf("event[1].Name = %q, want %q", events[1].Name, "Read")
	}
	if events[2].Type != "tool_result" {
		t.Errorf("event[2].Type = %q, want %q", events[2].Type, "tool_result")
	}
	if events[2].Content != "file contents here" {
		t.Errorf("event[2].Content = %q, want %q", events[2].Content, "file contents here")
	}
}

func TestEvent_IsTerminal(t *testing.T) {
	tests := []struct {
		eventType string
		want      bool
	}{
		{"result", true},
		{"error", true},
		{"system", false},
		{"assistant", false},
		{"tool_use", false},
		{"tool_result", false},
		{"pong", true},
	}
	for _, tt := range tests {
		ev := Event{Type: tt.eventType}
		if got := ev.IsTerminal(); got != tt.want {
			t.Errorf("Event{Type:%q}.IsTerminal() = %v, want %v", tt.eventType, got, tt.want)
		}
	}
}

func TestBridge_ExecuteSync_ErrorEvent(t *testing.T) {
	dir := t.TempDir()

	errorMockJS := `
const readline = require('readline');
const rl = readline.createInterface({ input: process.stdin, terminal: false });
rl.on('line', (line) => {
    const req = JSON.parse(line);
    const rid = req.request_id || "";
    process.stdout.write(JSON.stringify({event:"error",request_id:rid,message:"auth failed"}) + "\n");
});
rl.on('close', () => process.exit(0));
`
	b := newMockBridge(t, dir, errorMockJS)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ev, err := b.ExecuteSync(ctx, Request{Command: "query", Prompt: "test"})
	if err != nil {
		t.Fatalf("ExecuteSync() error: %v", err)
	}
	if ev.Type != "error" {
		t.Errorf("Type = %q, want %q", ev.Type, "error")
	}
	if ev.Message != "auth failed" {
		t.Errorf("Message = %q, want %q", ev.Message, "auth failed")
	}
}

func TestBridge_LongLived_MultipleRequests(t *testing.T) {
	dir := t.TempDir()
	b := newMockBridge(t, dir, longLivedMockJS)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Request 1: ping
	if err := b.Ping(ctx); err != nil {
		t.Fatalf("Ping 1 error: %v", err)
	}

	// Request 2: query
	ev, err := b.ExecuteSync(ctx, Request{Command: "query", Prompt: "first"})
	if err != nil {
		t.Fatalf("Query 1 error: %v", err)
	}
	if ev.Type != "result" {
		t.Errorf("Query 1: Type = %q, want %q", ev.Type, "result")
	}

	// Request 3: another query on the SAME process
	ev, err = b.ExecuteSync(ctx, Request{Command: "query", Prompt: "second"})
	if err != nil {
		t.Fatalf("Query 2 error: %v", err)
	}
	if ev.Type != "result" {
		t.Errorf("Query 2: Type = %q, want %q", ev.Type, "result")
	}

	// Request 4: ping again
	if err := b.Ping(ctx); err != nil {
		t.Fatalf("Ping 2 error: %v", err)
	}
}

func TestBridge_Stop_And_Restart(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "mock.js"), []byte(longLivedMockJS), 0644); err != nil {
		t.Fatal(err)
	}

	b := New(dir)
	b.command = "node"
	b.args = []string{"mock.js"}
	t.Cleanup(func() { b.Stop() })

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Start and use
	if err := b.Ping(ctx); err != nil {
		t.Fatalf("Ping 1 error: %v", err)
	}

	// Stop
	b.Stop()

	// Use again — should auto-restart
	if err := b.Ping(ctx); err != nil {
		t.Fatalf("Ping after restart error: %v", err)
	}
}
