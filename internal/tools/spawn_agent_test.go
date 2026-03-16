package tools

import (
	"context"
	"strings"
	"testing"

	"github.com/kocar/aurelia/internal/agent"
)

type MockSpawner struct {
	CapturedTeamKey string
	CapturedUserID  string
	CapturedName    string
	CapturedTask    string
	CapturedRole    string
	CapturedWorkdir string
	CapturedTools   []string
	taskID          string
	err             error
}

func (m *MockSpawner) Spawn(ctx context.Context, teamKey, userID, agentName, roleDescription, taskPrompt string, allowedTools ...string) (string, error) {
	m.CapturedTeamKey = teamKey
	m.CapturedUserID = userID
	m.CapturedName = agentName
	m.CapturedRole = roleDescription
	m.CapturedTask = taskPrompt
	m.CapturedWorkdir, _ = agent.WorkdirFromContext(ctx)
	m.CapturedTools = append([]string(nil), allowedTools...)
	if m.taskID == "" {
		m.taskID = "task-123"
	}
	return m.taskID, m.err
}

func TestSpawnAgentTool_Execute(t *testing.T) {
	mockSpawner := &MockSpawner{}
	tool := NewSpawnAgentTool(mockSpawner)

	args := map[string]interface{}{
		"agent_name":       "TestReviewer",
		"role_description": "Review PRs",
		"task_prompt":      "Check main.go",
	}

	ctx := agent.WithTeamContext(context.Background(), "team-key-1", "user-1")
	result, err := tool.Execute(ctx, args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "Acionei o especialista") || !strings.Contains(result, "TestReviewer") {
		t.Errorf("unexpected execute output: %s", result)
	}

	if mockSpawner.CapturedName != "TestReviewer" {
		t.Errorf("Spawner captured wrong agent name: %s", mockSpawner.CapturedName)
	}
	if mockSpawner.CapturedTeamKey != "team-key-1" || mockSpawner.CapturedUserID != "user-1" {
		t.Errorf("Spawner captured wrong context: team=%s user=%s", mockSpawner.CapturedTeamKey, mockSpawner.CapturedUserID)
	}
	if mockSpawner.CapturedTask != "Check main.go" || mockSpawner.CapturedRole != "Review PRs" {
		t.Errorf("Spawner captured wrong task payload: task=%s role=%s", mockSpawner.CapturedTask, mockSpawner.CapturedRole)
	}
}

func TestSpawnAgentTool_Execute_ForwardsAllowedTools(t *testing.T) {
	mockSpawner := &MockSpawner{}
	tool := NewSpawnAgentTool(mockSpawner)

	args := map[string]interface{}{
		"agent_name":       "WebResearcher",
		"role_description": "Busca fatos externos",
		"task_prompt":      "Pesquisar docs atuais",
		"allowed_tools":    []any{"web_search", "read_file"},
	}

	ctx := agent.WithTeamContext(context.Background(), "team-key-2", "user-2")
	if _, err := tool.Execute(ctx, args); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(mockSpawner.CapturedTools) != 2 {
		t.Fatalf("expected 2 allowed tools, got %#v", mockSpawner.CapturedTools)
	}
	if mockSpawner.CapturedTools[0] != "web_search" || mockSpawner.CapturedTools[1] != "read_file" {
		t.Fatalf("unexpected allowed tools: %#v", mockSpawner.CapturedTools)
	}
}

func TestSpawnAgentTool_Execute_ForwardsWorkdirThroughContext(t *testing.T) {
	mockSpawner := &MockSpawner{}
	tool := NewSpawnAgentTool(mockSpawner)

	args := map[string]interface{}{
		"agent_name":       "Builder",
		"role_description": "Implementa feature",
		"task_prompt":      "Corrigir bug",
		"workdir":          `C:\projetos\api-alvo`,
	}

	ctx := agent.WithTeamContext(context.Background(), "team-key-6", "user-6")
	if _, err := tool.Execute(ctx, args); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if mockSpawner.CapturedWorkdir != `C:\projetos\api-alvo` {
		t.Fatalf("expected forwarded workdir, got %q", mockSpawner.CapturedWorkdir)
	}
}

func TestSpawnAgentTool_Execute_AppliesResearcherDefaultToolProfile(t *testing.T) {
	mockSpawner := &MockSpawner{}
	tool := NewSpawnAgentTool(mockSpawner)

	args := map[string]interface{}{
		"agent_name":       "WebResearcher",
		"role_description": "Especialista em pesquisa na internet",
		"task_prompt":      "Busque documentacao atual e referencias externas",
	}

	ctx := agent.WithTeamContext(context.Background(), "team-key-3", "user-3")
	if _, err := tool.Execute(ctx, args); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"web_search", "read_file", "send_team_message", "read_team_inbox"}
	if len(mockSpawner.CapturedTools) != len(want) {
		t.Fatalf("unexpected allowed tools count: got %#v want %#v", mockSpawner.CapturedTools, want)
	}
	for i := range want {
		if mockSpawner.CapturedTools[i] != want[i] {
			t.Fatalf("unexpected researcher tool profile: got %#v want %#v", mockSpawner.CapturedTools, want)
		}
	}
}

func TestSpawnAgentTool_Execute_AppliesImplementerDefaultToolProfile(t *testing.T) {
	mockSpawner := &MockSpawner{}
	tool := NewSpawnAgentTool(mockSpawner)

	args := map[string]interface{}{
		"agent_name":       "Builder",
		"role_description": "Implementa codigo da feature",
		"task_prompt":      "Corrigir bug e validar no ambiente",
	}

	ctx := agent.WithTeamContext(context.Background(), "team-key-4", "user-4")
	if _, err := tool.Execute(ctx, args); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"read_file", "write_file", "run_command", "send_team_message", "read_team_inbox"}
	if len(mockSpawner.CapturedTools) != len(want) {
		t.Fatalf("unexpected allowed tools count: got %#v want %#v", mockSpawner.CapturedTools, want)
	}
	for i := range want {
		if mockSpawner.CapturedTools[i] != want[i] {
			t.Fatalf("unexpected implementer tool profile: got %#v want %#v", mockSpawner.CapturedTools, want)
		}
	}
}

func TestSpawnAgentTool_Execute_AppliesReviewerDefaultToolProfile(t *testing.T) {
	mockSpawner := &MockSpawner{}
	tool := NewSpawnAgentTool(mockSpawner)

	args := map[string]interface{}{
		"agent_name":       "Reviewer",
		"role_description": "Revisor tecnico e QA",
		"task_prompt":      "Validar os testes e auditar a entrega",
	}

	ctx := agent.WithTeamContext(context.Background(), "team-key-5", "user-5")
	if _, err := tool.Execute(ctx, args); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"read_file", "run_command", "send_team_message", "read_team_inbox"}
	if len(mockSpawner.CapturedTools) != len(want) {
		t.Fatalf("unexpected allowed tools count: got %#v want %#v", mockSpawner.CapturedTools, want)
	}
	for i := range want {
		if mockSpawner.CapturedTools[i] != want[i] {
			t.Fatalf("unexpected reviewer tool profile: got %#v want %#v", mockSpawner.CapturedTools, want)
		}
	}
}


