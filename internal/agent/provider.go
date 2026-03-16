package agent

import (
	"context"
	"fmt"
)

// ToolCall represents a function call requested by the LLM
type ToolCall struct {
	ID        string
	Name      string
	Arguments map[string]interface{}
}

// ModelResponse is the standardized output from any LLM Provider
type ModelResponse struct {
	Content          string
	ReasoningContent string
	ToolCalls        []ToolCall
}

// LLMProvider is the interface for different AI providers (Gemini, DeepSeek, etc)
type LLMProvider interface {
	// GenerateContent sends a prompt with history and available tools, returning a response
	GenerateContent(ctx context.Context, systemPrompt string, history []Message, tools []Tool) (*ModelResponse, error)
}

// Message is the standard internal representation of a chat message in the loop
type Message struct {
	Role             string // "user", "assistant", "system", "tool"
	Content          string
	ReasoningContent string
	// ToolCallID is used when Role == "tool" to map the observation to the correct call
	ToolCallID string
	// ToolCalls is used when Role == "assistant" to remember the calls it requested
	ToolCalls []ToolCall
}

// Tool describes a function available to the LLM
type Tool struct {
	Name        string
	Description string
	// JSONSchema is typically passed as a map or struct to the provider
	JSONSchema map[string]interface{}
}

// ToolRegistry manages available tools
type ToolRegistry struct {
	tools map[string]func(context.Context, map[string]interface{}) (string, error)
	defs  []Tool
}

// NewToolRegistry creates a new registry
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]func(context.Context, map[string]interface{}) (string, error)),
		defs:  make([]Tool, 0),
	}
}

// Register adds a tool to the registry
func (r *ToolRegistry) Register(tool Tool, handler func(context.Context, map[string]interface{}) (string, error)) {
	r.tools[tool.Name] = handler
	r.defs = append(r.defs, tool)
}

// GetDefinitions returns the slice of tool definitions for the LLM
func (r *ToolRegistry) GetDefinitions() []Tool {
	return r.defs
}

// FilterDefinitions returns the tool definitions allowed for this execution.
// An empty allowed list means "all tools registered in the runtime".
func (r *ToolRegistry) FilterDefinitions(allowed []string) []Tool {
	if len(allowed) == 0 {
		return append([]Tool(nil), r.defs...)
	}

	allowedMap := make(map[string]bool)
	for _, a := range allowed {
		allowedMap[a] = true
	}

	var filtered []Tool
	for _, t := range r.defs {
		if allowedMap[t.Name] {
			filtered = append(filtered, t)
		}
	}
	return filtered
}

// Execute runs a specific tool
func (r *ToolRegistry) Execute(ctx context.Context, name string, args map[string]interface{}) (string, error) {
	handler, exists := r.tools[name]
	if !exists {
		return "", fmt.Errorf("tool %s not found in registry", name)
	}
	return handler(ctx, args)
}
