package cron

import (
	"context"
	"errors"
	"testing"

	"github.com/kocar/aurelia/internal/agent"
)

type fakeCronAgentExecutor struct {
	lastCtx          context.Context
	lastSystemPrompt string
	lastHistory      []agent.Message
	lastAllowedTools []string
	finalAnswer      string
	err              error
}

func (f *fakeCronAgentExecutor) Execute(ctx context.Context, systemPrompt string, history []agent.Message, allowedTools []string) ([]agent.Message, string, error) {
	f.lastCtx = ctx
	f.lastSystemPrompt = systemPrompt
	f.lastHistory = history
	f.lastAllowedTools = allowedTools
	return nil, f.finalAnswer, f.err
}

func TestCronRuntime_ExecuteJob_RunsAgentWithPromptAndContext(t *testing.T) {
	t.Parallel()

	executor := &fakeCronAgentExecutor{
		finalAnswer: "daily summary ready",
	}

	runtime := NewAgentCronRuntime(executor, "base system prompt", []string{"web_search", "run_command"})
	job := CronJob{
		ID:           "job-1",
		OwnerUserID:  "user-1",
		TargetChatID: 12345,
		ScheduleType: "cron",
		CronExpr:     "0 8 * * *",
		Prompt:       "Me entregue o resumo diario",
		Active:       true,
	}

	answer, err := runtime.ExecuteJob(context.Background(), job)
	if err != nil {
		t.Fatalf("ExecuteJob() error = %v", err)
	}
	if answer != "daily summary ready" {
		t.Fatalf("unexpected final answer: %q", answer)
	}
	if executor.lastCtx == nil {
		t.Fatalf("expected runtime to forward context")
	}
	if executor.lastSystemPrompt != "base system prompt" {
		t.Fatalf("unexpected system prompt: %q", executor.lastSystemPrompt)
	}
	if len(executor.lastAllowedTools) != 2 {
		t.Fatalf("unexpected allowed tools: %#v", executor.lastAllowedTools)
	}
	if len(executor.lastHistory) != 1 {
		t.Fatalf("expected one synthetic user message, got %d", len(executor.lastHistory))
	}
	if executor.lastHistory[0].Role != "user" || executor.lastHistory[0].Content != "Me entregue o resumo diario" {
		t.Fatalf("unexpected synthetic history: %#v", executor.lastHistory)
	}

	teamKey, userID, ok := agent.TeamContextFromContext(executor.lastCtx)
	if !ok {
		t.Fatalf("expected runtime to propagate team context")
	}
	if teamKey != "12345" || userID != "user-1" {
		t.Fatalf("unexpected team context: teamKey=%q userID=%q", teamKey, userID)
	}
}

func TestCronRuntime_ExecuteJob_ResolvesAllowedToolsFromPrompt(t *testing.T) {
	t.Parallel()

	executor := &fakeCronAgentExecutor{finalAnswer: "ok"}
	runtime := NewAgentCronRuntime(executor, "base system prompt", nil)
	job := CronJob{
		ID:           "job-heuristic",
		OwnerUserID:  "user-2",
		TargetChatID: 42,
		ScheduleType: "once",
		Prompt:       "rode os testes do projeto",
		Active:       true,
	}

	if _, err := runtime.ExecuteJob(context.Background(), job); err != nil {
		t.Fatalf("ExecuteJob() error = %v", err)
	}

	want := []string{"list_dir", "read_file", "run_command", "write_file"}
	if len(executor.lastAllowedTools) != len(want) {
		t.Fatalf("expected %v, got %#v", want, executor.lastAllowedTools)
	}
	for i := range want {
		if executor.lastAllowedTools[i] != want[i] {
			t.Fatalf("expected %v, got %#v", want, executor.lastAllowedTools)
		}
	}
}

func TestCronRuntime_ExecuteJob_PropagatesExecutorError(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("loop failed")
	executor := &fakeCronAgentExecutor{
		err: expectedErr,
	}

	runtime := NewAgentCronRuntime(executor, "base system prompt", []string{"web_search"})
	job := CronJob{
		ID:           "job-2",
		OwnerUserID:  "user-1",
		TargetChatID: 999,
		ScheduleType: "once",
		Prompt:       "Falhe",
		Active:       true,
	}

	_, err := runtime.ExecuteJob(context.Background(), job)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected executor error %v, got %v", expectedErr, err)
	}
}

func TestCronRuntime_ExecuteJob_RebuildsPromptPerExecution(t *testing.T) {
	t.Parallel()

	executor := &fakeCronAgentExecutor{
		finalAnswer: "ok",
	}
	buildCount := 0
	runtime := NewAgentCronRuntimeWithPromptBuilder(
		executor,
		"fallback prompt",
		[]string{"fallback"},
		func(ctx context.Context, job CronJob) (string, []string, error) {
			buildCount++
			return "prompt build #" + string(rune('0'+buildCount)), []string{"web_search", "run_command"}, nil
		},
	)
	job := CronJob{
		ID:           "job-3",
		OwnerUserID:  "user-99",
		TargetChatID: 777,
		ScheduleType: "cron",
		CronExpr:     "0 8 * * *",
		Prompt:       "Pesquisar noticias de hoje",
		Active:       true,
	}

	if _, err := runtime.ExecuteJob(context.Background(), job); err != nil {
		t.Fatalf("first ExecuteJob() error = %v", err)
	}
	if executor.lastSystemPrompt != "prompt build #1" {
		t.Fatalf("expected first rebuilt prompt, got %q", executor.lastSystemPrompt)
	}

	if _, err := runtime.ExecuteJob(context.Background(), job); err != nil {
		t.Fatalf("second ExecuteJob() error = %v", err)
	}
	if executor.lastSystemPrompt != "prompt build #2" {
		t.Fatalf("expected second rebuilt prompt, got %q", executor.lastSystemPrompt)
	}
	if buildCount != 2 {
		t.Fatalf("expected prompt builder to run twice, got %d", buildCount)
	}
}
