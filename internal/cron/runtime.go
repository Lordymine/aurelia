package cron

import (
	"context"
	"fmt"
)

// BridgeCronRuntime executes cron jobs via a BridgeExecutor.
type BridgeCronRuntime struct {
	executor      BridgeExecutor
	basePrompt    string
	promptBuilder func(ctx context.Context, job CronJob) (string, error)
}

// NewBridgeCronRuntime creates a runtime that uses the bridge executor.
func NewBridgeCronRuntime(executor BridgeExecutor, baseSystemPrompt string) *BridgeCronRuntime {
	return &BridgeCronRuntime{
		executor:   executor,
		basePrompt: baseSystemPrompt,
	}
}

// NewBridgeCronRuntimeWithPromptBuilder creates a runtime with dynamic prompt building.
func NewBridgeCronRuntimeWithPromptBuilder(executor BridgeExecutor, baseSystemPrompt string, promptBuilder func(ctx context.Context, job CronJob) (string, error)) *BridgeCronRuntime {
	return &BridgeCronRuntime{
		executor:      executor,
		basePrompt:    baseSystemPrompt,
		promptBuilder: promptBuilder,
	}
}

// ExecuteJob runs the job via the bridge executor.
func (r *BridgeCronRuntime) ExecuteJob(ctx context.Context, job CronJob) (string, error) {
	systemPrompt := r.basePrompt
	if r.promptBuilder != nil {
		prompt, err := r.promptBuilder(ctx, job)
		if err == nil {
			systemPrompt = prompt
		}
	}

	return r.executor.Execute(ctx, systemPrompt, job.Prompt)
}

// DeliveryFunc is called after a job completes to deliver its output.
type DeliveryFunc func(ctx context.Context, job CronJob, output string, execErr error) error

// NotifyingRuntime wraps a Runtime and delivers results after execution.
type NotifyingRuntime struct {
	inner   Runtime
	deliver DeliveryFunc
}

// NewNotifyingRuntime wraps an inner runtime with delivery notification.
func NewNotifyingRuntime(inner Runtime, deliver DeliveryFunc) *NotifyingRuntime {
	return &NotifyingRuntime{
		inner:   inner,
		deliver: deliver,
	}
}

// ExecuteJob runs the inner runtime and delivers the result.
func (r *NotifyingRuntime) ExecuteJob(ctx context.Context, job CronJob) (string, error) {
	if r.inner == nil {
		return "", fmt.Errorf("inner runtime is required")
	}

	output, err := r.inner.ExecuteJob(ctx, job)
	if r.deliver != nil {
		if deliverErr := r.deliver(ctx, job, output, err); deliverErr != nil {
			return output, deliverErr
		}
	}
	return output, err
}
