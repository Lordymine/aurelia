package cron

import (
	"context"
	"errors"
	"testing"
)

type fakeBridgeExecutor struct {
	lastCtx          context.Context
	lastSystemPrompt string
	lastUserPrompt   string
	result           string
	err              error
}

func (f *fakeBridgeExecutor) Execute(ctx context.Context, systemPrompt string, userPrompt string) (string, error) {
	f.lastCtx = ctx
	f.lastSystemPrompt = systemPrompt
	f.lastUserPrompt = userPrompt
	return f.result, f.err
}

func TestBridgeCronRuntime_ExecuteJob_RunsWithPromptAndContext(t *testing.T) {
	t.Parallel()

	executor := &fakeBridgeExecutor{result: "daily summary ready"}
	runtime := NewBridgeCronRuntime(executor, "base system prompt")
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
	if executor.lastSystemPrompt != "base system prompt" {
		t.Fatalf("unexpected system prompt: %q", executor.lastSystemPrompt)
	}
	if executor.lastUserPrompt != "Me entregue o resumo diario" {
		t.Fatalf("unexpected user prompt: %q", executor.lastUserPrompt)
	}
}

func TestBridgeCronRuntime_ExecuteJob_PropagatesExecutorError(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("bridge failed")
	executor := &fakeBridgeExecutor{err: expectedErr}
	runtime := NewBridgeCronRuntime(executor, "base system prompt")
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

func TestBridgeCronRuntime_ExecuteJob_RebuildsPromptPerExecution(t *testing.T) {
	t.Parallel()

	executor := &fakeBridgeExecutor{result: "ok"}
	buildCount := 0
	runtime := NewBridgeCronRuntimeWithPromptBuilder(
		executor,
		"fallback prompt",
		func(ctx context.Context, job CronJob) (string, error) {
			buildCount++
			return "prompt build #" + string(rune('0'+buildCount)), nil
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
