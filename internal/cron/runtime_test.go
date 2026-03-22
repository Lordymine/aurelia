package cron

import (
	"context"
	"errors"
	"testing"

	"github.com/kocar/aurelia/internal/agents"
	"github.com/kocar/aurelia/internal/bridge"
)

// --- fakes ---

type fakeBridgeExecutor struct {
	lastReq bridge.Request
	result  *bridge.Event
	err     error
}

func (f *fakeBridgeExecutor) Execute(_ context.Context, req bridge.Request) (*bridge.Event, error) {
	f.lastReq = req
	return f.result, f.err
}

type fakeRegistry struct {
	agents map[string]*agents.Agent
}

func (f *fakeRegistry) Get(name string) *agents.Agent {
	return f.agents[name]
}

type fakePersona struct {
	prompt string
	err    error
}

func (f *fakePersona) BuildPrompt() (string, error) {
	return f.prompt, f.err
}

type fakeMemory struct {
	injectResult string
	injectErr    error
	savedContent string
	savedAgent   string
}

func (f *fakeMemory) Inject(_ context.Context, _ string, _ int) (string, error) {
	return f.injectResult, f.injectErr
}

func (f *fakeMemory) Save(_ context.Context, content, _ string, agent string) error {
	f.savedContent = content
	f.savedAgent = agent
	return nil
}

// --- tests ---

func TestBridgeCronRuntime_ExecuteJob(t *testing.T) {
	t.Parallel()

	executor := &fakeBridgeExecutor{
		result: &bridge.Event{Type: "result", Content: "daily summary ready"},
	}
	registry := &fakeRegistry{agents: map[string]*agents.Agent{
		"news": {
			Name:         "news",
			Model:        "claude-sonnet-4-20250514",
			Prompt:       "You are a news agent.",
			AllowedTools: []string{"web_search"},
		},
	}}
	persona := &fakePersona{prompt: "I am Aurelia."}
	mem := &fakeMemory{injectResult: "## Relevant Memories\n\n- [fact] something"}

	runtime := NewBridgeCronRuntime(executor, registry, persona, mem)

	job := CronJob{
		ID:           "job-1",
		AgentName:    "news",
		ScheduleType: "cron",
		CronExpr:     "0 8 * * *",
		Prompt:       "Resumo diario de noticias",
		Active:       true,
	}

	result, err := runtime.ExecuteJob(context.Background(), job)
	if err != nil {
		t.Fatalf("ExecuteJob() error = %v", err)
	}
	if result.Output != "daily summary ready" {
		t.Fatalf("unexpected output: %q", result.Output)
	}

	// Verify bridge request
	if executor.lastReq.Command != "query" {
		t.Fatalf("expected command %q, got %q", "query", executor.lastReq.Command)
	}
	if executor.lastReq.Prompt != "Resumo diario de noticias" {
		t.Fatalf("unexpected prompt: %q", executor.lastReq.Prompt)
	}
	if executor.lastReq.Options.Model != "claude-sonnet-4-20250514" {
		t.Fatalf("unexpected model: %q", executor.lastReq.Options.Model)
	}
	if executor.lastReq.Options.PermissionMode != "bypassPermissions" {
		t.Fatalf("unexpected permission mode: %q", executor.lastReq.Options.PermissionMode)
	}

	// System prompt should contain persona + agent prompt + memory
	sp := executor.lastReq.Options.SystemPrompt
	if sp == "" {
		t.Fatal("system prompt is empty")
	}
	if !contains(sp, "I am Aurelia.") {
		t.Fatalf("system prompt missing persona: %q", sp)
	}
	if !contains(sp, "You are a news agent.") {
		t.Fatalf("system prompt missing agent prompt: %q", sp)
	}
	if !contains(sp, "Relevant Memories") {
		t.Fatalf("system prompt missing memory context: %q", sp)
	}

	// Verify memory was saved
	if mem.savedContent != "daily summary ready" {
		t.Fatalf("expected memory save with content %q, got %q", "daily summary ready", mem.savedContent)
	}
	if mem.savedAgent != "news" {
		t.Fatalf("expected memory save with agent %q, got %q", "news", mem.savedAgent)
	}
}

func TestBridgeCronRuntime_NoAgent(t *testing.T) {
	t.Parallel()

	executor := &fakeBridgeExecutor{
		result: &bridge.Event{Type: "result", Content: "done without agent"},
	}
	registry := &fakeRegistry{agents: map[string]*agents.Agent{}}
	persona := &fakePersona{prompt: "base"}
	mem := &fakeMemory{}

	runtime := NewBridgeCronRuntime(executor, registry, persona, mem)

	job := CronJob{
		ID:     "job-2",
		Prompt: "test",
	}

	result, err := runtime.ExecuteJob(context.Background(), job)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Output != "done without agent" {
		t.Fatalf("output = %q", result.Output)
	}
}

func TestBridgeCronRuntime_BridgeError(t *testing.T) {
	t.Parallel()

	executor := &fakeBridgeExecutor{
		result: &bridge.Event{Type: "error", Message: "timeout"},
	}
	registry := &fakeRegistry{agents: map[string]*agents.Agent{
		"test": {Name: "test", Prompt: "test agent"},
	}}
	persona := &fakePersona{prompt: "base"}
	mem := &fakeMemory{}

	runtime := NewBridgeCronRuntime(executor, registry, persona, mem)

	job := CronJob{ID: "job-3", AgentName: "test", Prompt: "test"}

	_, err := runtime.ExecuteJob(context.Background(), job)
	if err == nil {
		t.Fatal("expected error for bridge error event")
	}
	if !contains(err.Error(), "bridge error") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBridgeCronRuntime_BridgeExecuteFailure(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("connection refused")
	executor := &fakeBridgeExecutor{err: expectedErr}
	registry := &fakeRegistry{agents: map[string]*agents.Agent{
		"test": {Name: "test", Prompt: "test agent"},
	}}
	persona := &fakePersona{prompt: "base"}
	mem := &fakeMemory{}

	runtime := NewBridgeCronRuntime(executor, registry, persona, mem)

	job := CronJob{ID: "job-4", AgentName: "test", Prompt: "test"}

	_, err := runtime.ExecuteJob(context.Background(), job)
	if err == nil {
		t.Fatal("expected error")
	}
	if !contains(err.Error(), "bridge execute") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBridgeCronRuntime_PersonaError(t *testing.T) {
	t.Parallel()

	executor := &fakeBridgeExecutor{}
	registry := &fakeRegistry{agents: map[string]*agents.Agent{
		"test": {Name: "test", Prompt: "test agent"},
	}}
	persona := &fakePersona{err: errors.New("file not found")}
	mem := &fakeMemory{}

	runtime := NewBridgeCronRuntime(executor, registry, persona, mem)

	job := CronJob{ID: "job-5", AgentName: "test", Prompt: "test"}

	_, err := runtime.ExecuteJob(context.Background(), job)
	if err == nil {
		t.Fatal("expected error for persona failure")
	}
	if !contains(err.Error(), "build persona prompt") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBridgeCronRuntime_MemoryInjectFailureNonFatal(t *testing.T) {
	t.Parallel()

	executor := &fakeBridgeExecutor{
		result: &bridge.Event{Type: "result", Content: "ok"},
	}
	registry := &fakeRegistry{agents: map[string]*agents.Agent{
		"test": {Name: "test", Prompt: "test agent"},
	}}
	persona := &fakePersona{prompt: "base"}
	mem := &fakeMemory{injectErr: errors.New("embed failed")}

	runtime := NewBridgeCronRuntime(executor, registry, persona, mem)

	job := CronJob{ID: "job-6", AgentName: "test", Prompt: "test"}

	// Should succeed despite memory inject failure
	result, err := runtime.ExecuteJob(context.Background(), job)
	if err != nil {
		t.Fatalf("ExecuteJob() error = %v", err)
	}
	if result.Output != "ok" {
		t.Fatalf("unexpected output: %q", result.Output)
	}
}

func TestNotifyingRuntime_Delivers(t *testing.T) {
	t.Parallel()

	inner := &stubRuntime{result: &ExecutionResult{Output: "hello"}, err: nil}
	var delivered bool
	nr := NewNotifyingRuntime(inner, func(_ context.Context, _ CronJob, result *ExecutionResult, _ error) error {
		delivered = true
		if result.Output != "hello" {
			t.Fatalf("unexpected output in delivery: %q", result.Output)
		}
		return nil
	})

	job := CronJob{ID: "job-n1", AgentName: "test", Prompt: "test"}
	result, err := nr.ExecuteJob(context.Background(), job)
	if err != nil {
		t.Fatalf("ExecuteJob() error = %v", err)
	}
	if result.Output != "hello" {
		t.Fatalf("unexpected output: %q", result.Output)
	}
	if !delivered {
		t.Fatal("delivery func was not called")
	}
}

func TestNotifyingRuntime_NilInner(t *testing.T) {
	t.Parallel()

	nr := NewNotifyingRuntime(nil, nil)
	_, err := nr.ExecuteJob(context.Background(), CronJob{})
	if err == nil {
		t.Fatal("expected error for nil inner runtime")
	}
}

// --- helpers ---

type stubRuntime struct {
	result *ExecutionResult
	err    error
}

func (s *stubRuntime) ExecuteJob(_ context.Context, _ CronJob) (*ExecutionResult, error) {
	return s.result, s.err
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
