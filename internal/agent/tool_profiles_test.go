package agent

import (
	"reflect"
	"testing"
)

func TestResolveAllowedToolsForQuery_ChatOnlyByDefault(t *testing.T) {
	t.Parallel()

	allowed := ResolveAllowedToolsForQuery("oi, tudo bem?", nil)
	if len(allowed) != 0 {
		t.Fatalf("expected no tools for plain chat, got %v", allowed)
	}
}

func TestResolveAllowedToolsForQuery_LocalExecution(t *testing.T) {
	t.Parallel()

	allowed := ResolveAllowedToolsForQuery("rode os testes e veja os arquivos do repo", nil)
	expected := []string{"list_dir", "read_file", "run_command", "write_file"}
	if !reflect.DeepEqual(allowed, expected) {
		t.Fatalf("expected %v, got %v", expected, allowed)
	}
}

func TestResolveAllowedToolsForQuery_MergesExplicitTools(t *testing.T) {
	t.Parallel()

	allowed := ResolveAllowedToolsForQuery("pesquise na internet", []string{"custom_tool"})
	expected := []string{"custom_tool", "web_search"}
	if !reflect.DeepEqual(allowed, expected) {
		t.Fatalf("expected %v, got %v", expected, allowed)
	}
}

func TestResolveAllowedToolsForQuery_SchedulerAndTeamOps(t *testing.T) {
	t.Parallel()

	allowed := ResolveAllowedToolsForQuery("agende uma rotina e depois mostre o status do time", nil)
	expected := []string{
		"cancel_team",
		"create_schedule",
		"delete_schedule",
		"list_schedules",
		"pause_schedule",
		"pause_team",
		"read_team_inbox",
		"resume_schedule",
		"resume_team",
		"send_team_message",
		"spawn_agent",
		"team_status",
	}
	if !reflect.DeepEqual(allowed, expected) {
		t.Fatalf("expected %v, got %v", expected, allowed)
	}
}

func TestResolveAllowedToolsForQuery_DoesNotExposeMCPByDefault(t *testing.T) {
	t.Parallel()

	allowed := ResolveAllowedToolsForQuery("oi, quero entender esse projeto", nil)
	for _, tool := range allowed {
		if len(tool) >= 4 && tool[:4] == "mcp_" {
			t.Fatalf("expected no MCP tools by default, got %v", allowed)
		}
	}
}

func TestResolveAllowedToolsForQuery_PreservesExplicitMCPTool(t *testing.T) {
	t.Parallel()

	allowed := ResolveAllowedToolsForQuery("oi", []string{"mcp_docs_search"})
	expected := []string{"mcp_docs_search"}
	if !reflect.DeepEqual(allowed, expected) {
		t.Fatalf("expected %v, got %v", expected, allowed)
	}
}

func TestResolveAllowedToolsForWorker_DefaultsToMailboxAndExecution(t *testing.T) {
	t.Parallel()

	allowed := ResolveAllowedToolsForWorker("worker", "generalist", "faca a tarefa", nil)
	expected := []string{"read_file", "read_team_inbox", "run_command", "send_team_message", "write_file"}
	if !reflect.DeepEqual(allowed, expected) {
		t.Fatalf("expected %v, got %v", expected, allowed)
	}
}

func TestResolveAllowedToolsForWorker_ResearcherProfile(t *testing.T) {
	t.Parallel()

	allowed := ResolveAllowedToolsForWorker("researcher", "pesquisa externa", "buscar documentacao atual", nil)
	expected := []string{"read_file", "read_team_inbox", "send_team_message", "web_search"}
	if !reflect.DeepEqual(allowed, expected) {
		t.Fatalf("expected %v, got %v", expected, allowed)
	}
}
