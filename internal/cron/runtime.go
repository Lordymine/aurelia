package cron

import (
	"context"
	"fmt"

	"github.com/kocar/aurelia/internal/agents"
	"github.com/kocar/aurelia/internal/bridge"
)

// BridgeCronRuntime executes cron jobs via the Claude Code bridge,
// resolving agent config from the registry and injecting persona + memory.
type BridgeCronRuntime struct {
	bridge  BridgeExecutor
	agents  AgentRegistry
	persona PersonaBuilder
	memory  MemoryStore
}

// AgentRegistry resolves agent definitions by name.
type AgentRegistry interface {
	Get(name string) *agents.Agent
}

// PersonaBuilder builds the base system prompt from persona files.
type PersonaBuilder interface {
	BuildPrompt() (string, error)
}

// MemoryStore injects relevant memories and saves new ones.
type MemoryStore interface {
	Inject(ctx context.Context, query string, limit int) (string, error)
	Save(ctx context.Context, content, category, agent string) error
}

// NewBridgeCronRuntime creates a runtime that executes jobs via Bridge
// with agent config, persona prompt, and memory context.
func NewBridgeCronRuntime(
	b BridgeExecutor,
	ag AgentRegistry,
	p PersonaBuilder,
	m MemoryStore,
) *BridgeCronRuntime {
	return &BridgeCronRuntime{
		bridge:  b,
		agents:  ag,
		persona: p,
		memory:  m,
	}
}

// ExecuteJob builds the system prompt with persona and memory context,
// optionally resolves an agent, executes via Bridge, and saves the result.
func (r *BridgeCronRuntime) ExecuteJob(ctx context.Context, job CronJob) (string, error) {
	// 1. Build system prompt from persona
	basePrompt, err := r.persona.BuildPrompt()
	if err != nil {
		return "", fmt.Errorf("build persona prompt: %w", err)
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

	// 4. Inject memory context (best effort)
	if memCtx, err := r.memory.Inject(ctx, job.Prompt, 10); err == nil && memCtx != "" {
		opts.SystemPrompt += "\n\n" + memCtx
	}

	// 5. Execute via Bridge
	result, err := r.bridge.Execute(ctx, bridge.Request{
		Command: "query",
		Prompt:  job.Prompt,
		Options: opts,
	})
	if err != nil {
		return "", fmt.Errorf("bridge execute: %w", err)
	}
	if result.Type == "error" {
		return "", fmt.Errorf("bridge error: %s", result.Message)
	}

	// 6. Save result to memory (best effort)
	agentName := job.AgentName
	if agentName == "" {
		agentName = "cron"
	}
	_ = r.memory.Save(ctx, result.Content, "conversation", agentName)

	return result.Content, nil
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
