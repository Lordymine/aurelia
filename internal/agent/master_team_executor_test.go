package agent

import (
	"context"
	"reflect"
	"sort"
	"testing"
)

func TestLoopTaskExecutor_UsesWorkerFallbackProfileWhenTaskHasNoAllowedTools(t *testing.T) {
	t.Parallel()

	provider := &capturingToolProvider{response: &ModelResponse{Content: "done"}}
	registry := NewToolRegistry()
	registry.Register(Tool{Name: "web_search"}, func(ctx context.Context, args map[string]interface{}) (string, error) { return "", nil })
	registry.Register(Tool{Name: "read_file"}, func(ctx context.Context, args map[string]interface{}) (string, error) { return "", nil })
	registry.Register(Tool{Name: "write_file"}, func(ctx context.Context, args map[string]interface{}) (string, error) { return "", nil })
	registry.Register(Tool{Name: "run_command"}, func(ctx context.Context, args map[string]interface{}) (string, error) { return "", nil })
	registry.Register(Tool{Name: "send_team_message"}, func(ctx context.Context, args map[string]interface{}) (string, error) { return "", nil })
	registry.Register(Tool{Name: "read_team_inbox"}, func(ctx context.Context, args map[string]interface{}) (string, error) { return "", nil })

	executor := &loopTaskExecutor{
		llm:       provider,
		registry:  registry,
		agentName: "researcher",
		roleDesc:  "pesquisa externa",
	}

	_, err := executor.ExecuteTask(context.Background(), TeamTask{
		ID:     "task-1",
		Prompt: "buscar documentacao atual",
	})
	if err != nil {
		t.Fatalf("ExecuteTask() error = %v", err)
	}

	got := provider.CapturedTools()
	want := []string{"read_file", "read_team_inbox", "send_team_message", "web_search"}
	sort.Strings(got)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}
