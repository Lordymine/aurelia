package bridge

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

// newMockBridge creates a Bridge that uses `node mock.js` and writes the given
// JavaScript body as mock.js in dir. The JS code should read from stdin and
// write NDJSON to stdout to simulate the bridge protocol.
func newMockBridge(t *testing.T, dir string, jsBody string) *Bridge {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, "mock.js"), []byte(jsBody), 0644); err != nil {
		t.Fatal(err)
	}
	b := New(dir)
	b.command = "node"
	b.args = []string{"mock.js"}
	return b
}

// mockReadAndRespond wraps JS code so it reads one line from stdin before
// executing the response logic. This ensures the process waits for Go to
// write the request before responding.
func mockReadAndRespond(responseJS string) string {
	return `
process.stdin.resume();
process.stdin.once('data', function() {
    ` + responseJS + `
    process.stdin.destroy();
});
`
}

func TestBridge_Execute_ParsesEvents(t *testing.T) {
	dir := t.TempDir()

	b := newMockBridge(t, dir, mockReadAndRespond(`
    process.stdout.write(JSON.stringify({event:"system",session_id:"test-123",tools:["Read"],model:"claude-3"}) + "\n");
    process.stdout.write(JSON.stringify({event:"assistant",text:"hello world"}) + "\n");
    process.stdout.write(JSON.stringify({event:"result",content:"done",cost_usd:0.01,session_id:"test-123",duration_ms:100,num_turns:1}) + "\n");
`))

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

	b := newMockBridge(t, dir, mockReadAndRespond(`
    process.stdout.write(JSON.stringify({event:"system",session_id:"s1"}) + "\n");
    process.stdout.write(JSON.stringify({event:"assistant",text:"thinking..."}) + "\n");
    process.stdout.write(JSON.stringify({event:"result",content:"final answer",cost_usd:0.05,duration_ms:500,num_turns:3}) + "\n");
`))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ev, err := b.ExecuteSync(ctx, Request{Command: "query", Prompt: "test"})
	if err != nil {
		t.Fatalf("ExecuteSync() error: %v", err)
	}
	if ev.Type != "result" {
		t.Errorf("Type = %q, want %q", ev.Type, "result")
	}
	if ev.Content != "final answer" {
		t.Errorf("Content = %q, want %q", ev.Content, "final answer")
	}
	if ev.CostUSD != 0.05 {
		t.Errorf("CostUSD = %f, want %f", ev.CostUSD, 0.05)
	}
}

func TestBridge_Execute_ErrorEvent(t *testing.T) {
	dir := t.TempDir()

	b := newMockBridge(t, dir, mockReadAndRespond(`
    process.stdout.write(JSON.stringify({event:"error",message:"something went wrong"}) + "\n");
`))

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
	if runtime.GOOS == "windows" {
		t.Skip("process kill semantics differ on Windows; skip for CI stability")
	}

	dir := t.TempDir()

	// Script that hangs forever.
	b := newMockBridge(t, dir, `
process.stdin.resume();
process.stdin.once('data', function() {
    setTimeout(function() {}, 60000);
});
`)

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

	b := newMockBridge(t, dir, mockReadAndRespond(`
    process.stdout.write(JSON.stringify({event:"pong"}) + "\n");
`))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := b.Ping(ctx); err != nil {
		t.Fatalf("Ping() error: %v", err)
	}
}

func TestBridge_Execute_ToolUseEvent(t *testing.T) {
	dir := t.TempDir()

	b := newMockBridge(t, dir, mockReadAndRespond(`
    process.stdout.write(JSON.stringify({event:"system",session_id:"s1"}) + "\n");
    process.stdout.write(JSON.stringify({event:"tool_use",name:"Read",input:{file_path:"/tmp/test.txt"}}) + "\n");
    process.stdout.write(JSON.stringify({event:"tool_result",content:"file contents here"}) + "\n");
    process.stdout.write(JSON.stringify({event:"result",content:"done"}) + "\n");
`))

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
		{"pong", false},
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

	b := newMockBridge(t, dir, mockReadAndRespond(`
    process.stdout.write(JSON.stringify({event:"error",message:"auth failed"}) + "\n");
`))

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
