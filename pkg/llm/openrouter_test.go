package llm

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListOpenRouterModels(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer secret" {
			t.Fatalf("Authorization = %q", got)
		}
		if got := r.Header.Get("HTTP-Referer"); got != openRouterReferer {
			t.Fatalf("HTTP-Referer = %q", got)
		}
		if got := r.Header.Get("X-Title"); got != openRouterTitle {
			t.Fatalf("X-Title = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"data": [
				{
					"id":"anthropic/claude-sonnet-4",
					"name":"Claude Sonnet 4",
					"pricing":{"prompt":"0.000003","completion":"0.000015"},
					"supported_parameters":["tools","max_tokens"],
					"architecture":{"input_modalities":["text","image"]}
				},
				{
					"id":"google/gemini-2.5-flash:free",
					"name":"Gemini 2.5 Flash",
					"pricing":{"prompt":"0","completion":"0"},
					"supported_parameters":["max_tokens"],
					"architecture":{"input_modalities":["text"]}
				}
			]
		}`))
	}))
	defer server.Close()

	models, err := listOpenRouterModels(context.Background(), "secret", server.URL, server.Client())
	if err != nil {
		t.Fatalf("listOpenRouterModels() error = %v", err)
	}
	if len(models) != 4 {
		t.Fatalf("expected 4 models, got %d", len(models))
	}
	if models[0].ID != "openrouter/auto" || models[1].ID != "openrouter/free" {
		t.Fatalf("unexpected router defaults: %+v", models[:2])
	}
	if !models[2].SupportsImageInput {
		t.Fatal("expected anthropic model to support image input")
	}
	if !models[2].SupportsTools {
		t.Fatal("expected anthropic model to advertise tools support")
	}
	if models[3].SupportsImageInput {
		t.Fatal("expected google test model to be text-only")
	}
	if !models[3].IsFree {
		t.Fatal("expected google test model to be marked as free")
	}
}

func TestNewOpenRouterProvider(t *testing.T) {
	t.Parallel()

	provider := NewOpenRouterProvider("secret", "openrouter/auto")
	if provider.baseURL != openRouterChatCompletionsURL {
		t.Fatalf("baseURL = %q", provider.baseURL)
	}
	if provider.headers["HTTP-Referer"] != openRouterReferer {
		t.Fatalf("HTTP-Referer = %q", provider.headers["HTTP-Referer"])
	}
}
