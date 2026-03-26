package cron

import (
	"context"
	"fmt"

	"github.com/kocar/aurelia/internal/agents"
	"github.com/kocar/aurelia/internal/bridge"
)

// BridgeCronRuntime executes cron jobs via the Claude Code bridge,
// resolving agent config from the registry and injecting persona prompt.
type BridgeCronRuntime struct {
	bridge  BridgeExecutor
	agents  AgentRegistry
	persona PersonaBuilder
}

// AgentRegistry resolves agent definitions by name.
type AgentRegistry interface {
	Get(name string) *agents.Agent
}

// PersonaBuilder builds the base system prompt from persona files.
type PersonaBuilder interface {
	BuildPrompt() (string, error)
}

// NewBridgeCronRuntime creates a runtime that executes jobs via Bridge
// with agent config and persona prompt.
func NewBridgeCronRuntime(
	b BridgeExecutor,
	ag AgentRegistry,
	p PersonaBuilder,
) *BridgeCronRuntime {
	return &BridgeCronRuntime{
		bridge:  b,
		agents:  ag,
		persona: p,
	}
}

// ExecuteJob builds the system prompt with persona and memory context,
// optionally resolves an agent, executes via Bridge, and saves the result.
func (r *BridgeCronRuntime) ExecuteJob(ctx context.Context, job CronJob) (*ExecutionResult, error) {
	// 1. Build system prompt from persona
	basePrompt, err := r.persona.BuildPrompt()
	if err != nil {
		return nil, fmt.Errorf("build persona prompt: %w", err)
	}
	systemPrompt := basePrompt

	// 2. Build request options — block Telegram plugin tools to prevent
	// delivery via wrong bot. All other user MCPs/plugins remain available.
	opts := bridge.RequestOptions{
		SystemPrompt:   systemPrompt,
		PermissionMode: "bypassPermissions",
		DisabledTools: bridge.TelegramPluginTools,
	}

	// 3. Apply agent config if available
	agent := r.agents.Get(job.AgentName)
	if agent != nil {
		systemPrompt += "\n\n" + agent.Prompt
		opts.SystemPrompt = systemPrompt
		opts.Model = agent.Model
		opts.Cwd = agent.Cwd
		opts.AllowedTools = agent.AllowedTools
		opts.MCPServers = agent.MCPServers
	}

	// 4. Execute via Bridge
	ev, err := r.bridge.Execute(ctx, bridge.Request{
		Command: "query",
		Prompt:  job.Prompt,
		Options: opts,
	})
	if err != nil {
		return nil, fmt.Errorf("bridge execute: %w", err)
	}
	if ev.Type == "error" {
		return nil, fmt.Errorf("bridge error: %s", ev.Message)
	}

	return &ExecutionResult{
		Output:    ev.Content,
		SessionID: ev.SessionID,
		CostUSD:   ev.CostUSD,
		NumTurns:  ev.NumTurns,
	}, nil
}

// BridgeAdapter wraps *bridge.Bridge to satisfy BridgeExecutor.
type BridgeAdapter struct {
	B *bridge.Bridge
}

// Execute calls bridge.ExecuteSync and returns the terminal event.
func (a *BridgeAdapter) Execute(ctx context.Context, req bridge.Request) (*bridge.Event, error) {
	return a.B.ExecuteSync(ctx, req)
}

// DeliveryFunc is called after a job completes to deliver its output.
type DeliveryFunc func(ctx context.Context, job CronJob, result *ExecutionResult, execErr error) error

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
func (r *NotifyingRuntime) ExecuteJob(ctx context.Context, job CronJob) (*ExecutionResult, error) {
	if r.inner == nil {
		return nil, fmt.Errorf("inner runtime is required")
	}

	result, err := r.inner.ExecuteJob(ctx, job)
	if r.deliver != nil {
		if deliverErr := r.deliver(ctx, job, result, err); deliverErr != nil {
			return result, deliverErr
		}
	}
	return result, err
}
