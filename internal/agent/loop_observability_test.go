package agent

import (
	"context"
	"strings"
	"testing"

	"github.com/kocar/aurelia/internal/observability"
)

type recordingProvider struct {
	response *ModelResponse
}

func (r *recordingProvider) GenerateContent(ctx context.Context, systemPrompt string, history []Message, tools []Tool) (*ModelResponse, error) {
	return r.response, nil
}

type recordingObserver struct {
	operations []observability.Operation
}

func (r *recordingObserver) RecordOperation(ctx context.Context, operation observability.Operation) error {
	r.operations = append(r.operations, operation)
	return nil
}

func (r *recordingObserver) ListRecentOperations(ctx context.Context, limit int) ([]observability.Operation, error) {
	return nil, nil
}

func (r *recordingObserver) ListFailedOperations(ctx context.Context, limit int) ([]observability.Operation, error) {
	return nil, nil
}

func TestLoop_Run_ObservesToolPayloadMetrics(t *testing.T) {
	t.Parallel()

	provider := &recordingProvider{response: &ModelResponse{Content: "ok"}}
	observer := &recordingObserver{}
	registry := NewToolRegistry()
	registry.Register(Tool{Name: "read_file", Description: "reads files"}, func(ctx context.Context, args map[string]interface{}) (string, error) {
		return "", nil
	})
	registry.Register(Tool{Name: "run_command", Description: "runs commands"}, func(ctx context.Context, args map[string]interface{}) (string, error) {
		return "", nil
	})

	loop := NewLoopWithObserver(provider, registry, 1, observer)
	if _, _, err := loop.Run(context.Background(), "base prompt", nil, []string{"read_file", "run_command"}); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	for _, operation := range observer.operations {
		if operation.Component == "agent.loop" && operation.Operation == "tool_context" {
			if !strings.Contains(operation.Summary, "tool_count=2") {
				t.Fatalf("expected tool_count in summary, got %q", operation.Summary)
			}
			if !strings.Contains(operation.Summary, "tool_payload_bytes=") {
				t.Fatalf("expected tool_payload_bytes in summary, got %q", operation.Summary)
			}
			return
		}
	}

	t.Fatalf("expected tool_context operation, got %+v", observer.operations)
}

func TestLoop_Run_ObservesOversizedToolOutput(t *testing.T) {
	t.Parallel()

	provider := &recordingProvider{response: &ModelResponse{
		ToolCalls: []ToolCall{{
			ID:   "call-1",
			Name: "read_file",
		}},
	}}
	observer := &recordingObserver{}
	registry := NewToolRegistry()
	registry.Register(Tool{Name: "read_file", Description: "reads files"}, func(ctx context.Context, args map[string]interface{}) (string, error) {
		return string(make([]rune, OversizedToolOutputThresholdChars+1)), nil
	})

	loop := NewLoopWithObserver(provider, registry, 1, observer)
	if _, _, err := loop.Run(context.Background(), "base prompt", nil, []string{"read_file"}); err == nil || !strings.Contains(err.Error(), "max iterations reached") {
		t.Fatalf("expected max iterations error after tool call, got %v", err)
	}

	for _, operation := range observer.operations {
		if operation.Component == "agent.tool" && operation.Operation == "read_file" {
			if !strings.Contains(operation.Summary, "oversized=true") {
				t.Fatalf("expected oversized summary, got %q", operation.Summary)
			}
			if !strings.Contains(operation.Summary, "raw_chars=") || !strings.Contains(operation.Summary, "compacted_chars=") {
				t.Fatalf("expected raw/compacted sizes in summary, got %q", operation.Summary)
			}
			return
		}
	}

	t.Fatalf("expected tool operation, got %+v", observer.operations)
}
