package llm

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anthropics/anthropic-sdk-go/option"

	"github.com/kocar/aurelia/internal/agent"
)

func TestAnthropicProviderGenerateContent_TextResponse(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/messages" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if got := r.Header.Get("x-api-key"); got != "secret" {
			t.Fatalf("x-api-key = %q", got)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("ReadAll() error = %v", err)
		}

		var payload map[string]any
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("Unmarshal() error = %v", err)
		}
		if payload["model"] != "claude-sonnet-4-6" {
			t.Fatalf("model = %v", payload["model"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "msg_123",
			"type": "message",
			"role": "assistant",
			"model": "claude-sonnet-4-6",
			"content": [
				{"type":"text","text":"Resposta Anthropic"}
			],
			"stop_reason": "end_turn",
			"usage": {"input_tokens": 10, "output_tokens": 5}
		}`))
	}))
	defer server.Close()

	provider := NewAnthropicProvider("secret", "claude-sonnet-4-6", option.WithBaseURL(server.URL+"/"))

	resp, err := provider.GenerateContent(context.Background(), "system", []agent.Message{
		{Role: "user", Content: "Oi"},
	}, nil)
	if err != nil {
		t.Fatalf("GenerateContent() error = %v", err)
	}
	if resp.Content != "Resposta Anthropic" {
		t.Fatalf("Content = %q", resp.Content)
	}
}

func TestAnthropicProviderGenerateContent_ToolUseResponse(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "msg_456",
			"type": "message",
			"role": "assistant",
			"model": "claude-sonnet-4-6",
			"content": [
				{"type":"tool_use","id":"toolu_1","name":"read_file","input":{"path":"README.md"}}
			],
			"stop_reason": "tool_use",
			"usage": {"input_tokens": 20, "output_tokens": 7}
		}`))
	}))
	defer server.Close()

	provider := NewAnthropicProvider("secret", "claude-sonnet-4-6", option.WithBaseURL(server.URL+"/"))

	resp, err := provider.GenerateContent(context.Background(), "", []agent.Message{
		{Role: "user", Content: "Leia o README"},
	}, []agent.Tool{{
		Name:        "read_file",
		Description: "Read a file",
		JSONSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"path": map[string]interface{}{"type": "string"},
			},
			"required": []interface{}{"path"},
		},
	}})
	if err != nil {
		t.Fatalf("GenerateContent() error = %v", err)
	}
	if len(resp.ToolCalls) != 1 {
		t.Fatalf("expected 1 tool call, got %d", len(resp.ToolCalls))
	}
	if resp.ToolCalls[0].Name != "read_file" {
		t.Fatalf("tool name = %q", resp.ToolCalls[0].Name)
	}
	if resp.ToolCalls[0].Arguments["path"] != "README.md" {
		t.Fatalf("tool args = %#v", resp.ToolCalls[0].Arguments)
	}
}

func TestBuildAnthropicUserBlocks_IncludesImage(t *testing.T) {
	t.Parallel()

	blocks := buildAnthropicUserBlocks(agent.Message{
		Role:    "user",
		Content: "describe",
		Parts: []agent.ContentPart{
			{Type: agent.ContentPartText, Text: "describe"},
			{Type: agent.ContentPartImage, MIMEType: "image/jpeg", Data: []byte("jpg-bytes")},
		},
	})

	if len(blocks) != 2 {
		t.Fatalf("expected 2 blocks, got %d", len(blocks))
	}
	if blocks[1].OfImage == nil {
		t.Fatal("expected image block")
	}
}
