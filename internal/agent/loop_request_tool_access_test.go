package agent

import (
	"context"
	"strings"
	"testing"
)

type queuedProvider struct {
	responses []*ModelResponse
	toolsSeen [][]string
}

func (q *queuedProvider) GenerateContent(ctx context.Context, systemPrompt string, history []Message, tools []Tool) (*ModelResponse, error) {
	names := make([]string, 0, len(tools))
	for _, tool := range tools {
		names = append(names, tool.Name)
	}
	q.toolsSeen = append(q.toolsSeen, names)

	if len(q.responses) == 0 {
		return &ModelResponse{Content: "ok"}, nil
	}
	response := q.responses[0]
	q.responses = q.responses[1:]
	return response, nil
}

func TestLoop_Run_ReentersAfterRequestToolAccess(t *testing.T) {
	t.Parallel()

	provider := &queuedProvider{
		responses: []*ModelResponse{
			{
				ToolCalls: []ToolCall{{
					ID:   "call-1",
					Name: requestToolAccessToolName,
					Arguments: map[string]interface{}{
						"capability": "mcp:playwright",
					},
				}},
			},
			{
				ToolCalls: []ToolCall{{
					ID:   "call-2",
					Name: "mcp_playwright_browser_navigate",
				}},
			},
			{Content: "feito"},
		},
	}

	registry := NewToolRegistry()
	registry.Register(Tool{Name: "read_file", Description: "reads files"}, func(ctx context.Context, args map[string]interface{}) (string, error) {
		return "", nil
	})
	registry.Register(Tool{Name: "mcp_playwright_browser_navigate", Description: "navigates browser"}, func(ctx context.Context, args map[string]interface{}) (string, error) {
		return "navegado", nil
	})

	loop := NewLoop(provider, registry, 4)
	ctx := WithDynamicToolAccess(context.Background(), true)
	history, finalAnswer, err := loop.Run(ctx, "base prompt", nil, []string{"read_file"})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if finalAnswer != "feito" {
		t.Fatalf("finalAnswer = %q", finalAnswer)
	}

	if len(provider.toolsSeen) < 2 {
		t.Fatalf("expected at least two provider calls, got %d", len(provider.toolsSeen))
	}
	if !containsTool(provider.toolsSeen[0], requestToolAccessToolName) {
		t.Fatalf("expected request tool access in initial tool set, got %v", provider.toolsSeen[0])
	}
	if containsTool(provider.toolsSeen[0], "mcp_playwright_browser_navigate") {
		t.Fatalf("expected playwright tool to stay hidden initially, got %v", provider.toolsSeen[0])
	}
	if !containsTool(provider.toolsSeen[1], "mcp_playwright_browser_navigate") {
		t.Fatalf("expected playwright tool after expansion, got %v", provider.toolsSeen[1])
	}

	foundExpansion := false
	for _, message := range history {
		if message.Role != "tool" {
			continue
		}
		if strings.Contains(message.Content, `"status":"expanded"`) {
			foundExpansion = true
			break
		}
	}
	if !foundExpansion {
		t.Fatalf("expected tool history to include expansion result, got %+v", history)
	}
}

func containsTool(names []string, target string) bool {
	for _, name := range names {
		if name == target {
			return true
		}
	}
	return false
}
