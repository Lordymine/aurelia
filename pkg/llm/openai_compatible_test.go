package llm

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kocar/aurelia/internal/agent"
)

func TestOpenAICompatibleProvider_GenerateContent(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer secret" {
			t.Fatalf("Authorization = %q", got)
		}
		if got := r.Header.Get("User-Agent"); got != "Aurelia-Test" {
			t.Fatalf("User-Agent = %q", got)
		}
		if got := r.Header.Get("HTTP-Referer"); got != "https://example.test" {
			t.Fatalf("HTTP-Referer = %q", got)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("ReadAll() error = %v", err)
		}

		var payload map[string]any
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("Unmarshal() error = %v", err)
		}
		if payload["model"] != "model-x" {
			t.Fatalf("model = %v", payload["model"])
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"choices": [
				{
					"message": {
						"role": "assistant",
						"content": "ok",
						"tool_calls": [
							{
								"id": "call_1",
								"type": "function",
								"function": {
									"name": "read_file",
									"arguments": "{\"path\":\"README.md\"}"
								}
							}
						]
					}
				}
			]
		}`))
	}))
	defer server.Close()

	provider := NewOpenAICompatibleProvider(OpenAICompatibleConfig{
		Provider:  "openai",
		APIKey:    "secret",
		BaseURL:   server.URL,
		Model:     "model-x",
		UserAgent: "Aurelia-Test",
		Headers: map[string]string{
			"HTTP-Referer": "https://example.test",
		},
	})

	resp, err := provider.GenerateContent(context.Background(), "system", []agent.Message{
		{Role: "user", Content: "hello"},
	}, []agent.Tool{{
		Name:        "read_file",
		Description: "Read a file",
		JSONSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"path": map[string]interface{}{"type": "string"},
			},
		},
	}})
	if err != nil {
		t.Fatalf("GenerateContent() error = %v", err)
	}
	if resp.Content != "ok" {
		t.Fatalf("Content = %q", resp.Content)
	}
	if len(resp.ToolCalls) != 1 {
		t.Fatalf("expected 1 tool call, got %d", len(resp.ToolCalls))
	}
	if resp.ToolCalls[0].Name != "read_file" {
		t.Fatalf("tool name = %q", resp.ToolCalls[0].Name)
	}
}

func TestBuildOpenAICompatibleRequest_IncludesTools(t *testing.T) {
	t.Parallel()

	reqBody, err := buildOpenAICompatibleRequest("model-x", "system", []agent.Message{
		{Role: "user", Content: "hello"},
	}, []agent.Tool{{
		Name:        "run_command",
		Description: "Run a command",
		JSONSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"command": map[string]interface{}{"type": "string"},
			},
		},
	}})
	if err != nil {
		t.Fatalf("buildOpenAICompatibleRequest() error = %v", err)
	}

	if reqBody["model"] != "model-x" {
		t.Fatalf("model = %v", reqBody["model"])
	}
	if _, ok := reqBody["tools"]; !ok {
		t.Fatal("expected tools in request body")
	}
}

func TestBuildOpenAICompatibleRequest_IncludesImageParts(t *testing.T) {
	t.Parallel()

	reqBody, err := buildOpenAICompatibleRequest("model-x", "system", []agent.Message{
		{
			Role:    "user",
			Content: "describe",
			Parts: []agent.ContentPart{
				{Type: agent.ContentPartText, Text: "describe"},
				{Type: agent.ContentPartImage, MIMEType: "image/jpeg", Data: []byte("jpg-bytes")},
			},
		},
	}, nil)
	if err != nil {
		t.Fatalf("buildOpenAICompatibleRequest() error = %v", err)
	}

	messages, ok := reqBody["messages"].([]map[string]any)
	if !ok {
		t.Fatalf("messages type = %T", reqBody["messages"])
	}
	content, ok := messages[1]["content"].([]map[string]any)
	if !ok {
		t.Fatalf("content type = %T", messages[1]["content"])
	}
	if content[1]["type"] != "image_url" {
		t.Fatalf("image content type = %v", content[1]["type"])
	}
}

func TestOpenAICompatibleProvider_GenerateContent_MapsVision404ToNaturalFallback(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"message":"No endpoints found that support image input","code":404}}`))
	}))
	defer server.Close()

	provider := NewOpenAICompatibleProvider(OpenAICompatibleConfig{
		Provider: "kilo",
		APIKey:   "secret",
		BaseURL:  server.URL,
		Model:    "glm-5-turbo",
	})

	_, err := provider.GenerateContent(context.Background(), "system", []agent.Message{{
		Role:    "user",
		Content: "analise",
		Parts: []agent.ContentPart{
			{Type: agent.ContentPartText, Text: "analise"},
			{Type: agent.ContentPartImage, MIMEType: "image/jpeg", Data: []byte("jpg-bytes")},
		},
	}}, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if _, ok := err.(VisionUnsupportedError); !ok {
		t.Fatalf("error type = %T", err)
	}
}
