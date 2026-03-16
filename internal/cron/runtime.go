package cron

import (
	"context"
	"fmt"
	"strconv"

	"github.com/kocar/aurelia/internal/agent"
)

type AgentCronRuntime struct {
	executor      AgentExecutor
	basePrompt    string
	allowedTools  []string
	promptBuilder func(ctx context.Context, job CronJob) (string, []string, error)
}

func NewAgentCronRuntime(executor AgentExecutor, baseSystemPrompt string, allowedTools []string) *AgentCronRuntime {
	return &AgentCronRuntime{
		executor:     executor,
		basePrompt:   baseSystemPrompt,
		allowedTools: allowedTools,
	}
}

func NewAgentCronRuntimeWithPromptBuilder(executor AgentExecutor, baseSystemPrompt string, allowedTools []string, promptBuilder func(ctx context.Context, job CronJob) (string, []string, error)) *AgentCronRuntime {
	return &AgentCronRuntime{
		executor:      executor,
		basePrompt:    baseSystemPrompt,
		allowedTools:  allowedTools,
		promptBuilder: promptBuilder,
	}
}

func (r *AgentCronRuntime) ExecuteJob(ctx context.Context, job CronJob) (string, error) {
	ctx = agent.WithTeamContext(ctx, formatChatTeamKey(job.TargetChatID), job.OwnerUserID)

	systemPrompt := r.basePrompt
	allowedTools := r.allowedTools
	if r.promptBuilder != nil {
		prompt, tools, err := r.promptBuilder(ctx, job)
		if err == nil {
			systemPrompt = prompt
			allowedTools = tools
		}
	}

	_, finalAnswer, err := r.executor.Execute(ctx, systemPrompt, []agent.Message{{
		Role:    "user",
		Content: job.Prompt,
	}}, allowedTools)
	if err != nil {
		return "", err
	}
	return finalAnswer, nil
}

func formatChatTeamKey(chatID int64) string {
	return strconv.FormatInt(chatID, 10)
}

type DeliveryFunc func(ctx context.Context, job CronJob, output string, execErr error) error

type NotifyingRuntime struct {
	inner   Runtime
	deliver DeliveryFunc
}

func NewNotifyingRuntime(inner Runtime, deliver DeliveryFunc) *NotifyingRuntime {
	return &NotifyingRuntime{
		inner:   inner,
		deliver: deliver,
	}
}

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


