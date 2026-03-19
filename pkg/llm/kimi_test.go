package llm

import (
	"strings"
	"testing"

	"github.com/kocar/aurelia/internal/agent"
)

func TestExtractToolCallsFromContent_ParsesSinglePseudoToolCall(t *testing.T) {
	t.Parallel()

	content := `Calling tools: run_command:4座{"command":"Invoke-RestMethod -Uri \"http://localhost:3000/cep/01310100\" -Method GET","workdir":"D:\\projetos\\cep-api"}`

	toolCalls, cleanedContent := extractToolCallsFromContent(content)
	if len(toolCalls) != 1 {
		t.Fatalf("expected 1 parsed tool call, got %d", len(toolCalls))
	}
	if toolCalls[0].Name != "run_command" {
		t.Fatalf("expected run_command, got %q", toolCalls[0].Name)
	}
	if toolCalls[0].Arguments["workdir"] != `D:\projetos\cep-api` {
		t.Fatalf("unexpected workdir: %#v", toolCalls[0].Arguments["workdir"])
	}
	if cleanedContent != "" {
		t.Fatalf("expected cleaned content to be empty, got %q", cleanedContent)
	}
}

func TestExtractToolCallsFromContent_IgnoresRegularContent(t *testing.T) {
	t.Parallel()

	toolCalls, cleanedContent := extractToolCallsFromContent("Resposta final normal")
	if len(toolCalls) != 0 {
		t.Fatalf("expected no parsed tool calls, got %d", len(toolCalls))
	}
	if cleanedContent != "Resposta final normal" {
		t.Fatalf("unexpected cleaned content: %q", cleanedContent)
	}
}

func TestExtractToolCallsFromContent_ParsesPseudoToolCallWithEmojiAndSuffix(t *testing.T) {
	t.Parallel()

	content := "Calling tools: run_command:1️⃣{\"command\": \"try { Invoke-RestMethod -Uri \\\"http://localhost:3000/health\\\" -Method GET -TimeoutSec 3 } catch { Write-Host \\\"down\\\" }\", \"workdir\": \"D:\\\\projetos\\\\cep-api\"} ️🔄"

	toolCalls, cleanedContent := extractToolCallsFromContent(content)
	if len(toolCalls) != 1 {
		t.Fatalf("expected 1 parsed tool call, got %d", len(toolCalls))
	}
	if toolCalls[0].Name != "run_command" {
		t.Fatalf("expected run_command, got %q", toolCalls[0].Name)
	}
	if toolCalls[0].ID == "" {
		t.Fatalf("expected non-empty tool call id")
	}
	if toolCalls[0].Arguments["workdir"] != `D:\projetos\cep-api` {
		t.Fatalf("unexpected workdir: %#v", toolCalls[0].Arguments["workdir"])
	}
	if cleanedContent != "" {
		t.Fatalf("expected cleaned content to be empty, got %q", cleanedContent)
	}
}

func TestKimiProviderApplyFallbackToolCalls_RejectsMalformedCallingToolsContent(t *testing.T) {
	t.Parallel()

	provider := &KimiProvider{}
	result := &agent.ModelResponse{
		Content: "Calling tools: list_dir, read_file, read_file, read_file",
	}

	err := provider.applyFallbackToolCalls(result)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "malformed tool-call content") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestKimiProviderApplyFallbackToolCalls_KeepsStructuredToolCallsWithoutSyntheticContent(t *testing.T) {
	t.Parallel()

	provider := &KimiProvider{}
	result := &agent.ModelResponse{
		ToolCalls: []agent.ToolCall{{
			ID:   "call-1",
			Name: "read_file",
		}},
	}

	if err := provider.applyFallbackToolCalls(result); err != nil {
		t.Fatalf("applyFallbackToolCalls() error = %v", err)
	}
	if result.Content != "" {
		t.Fatalf("expected content to remain empty, got %q", result.Content)
	}
}
