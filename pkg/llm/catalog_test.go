package llm

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anthropics/anthropic-sdk-go/option"
)

func TestListModels_KimiReturnsFallbackCatalog(t *testing.T) {
	t.Parallel()

	models, err := ListModels(context.Background(), "kimi", ModelCatalogCredentials{})
	if err != nil {
		t.Fatalf("ListModels() error = %v", err)
	}
	if len(models) == 0 {
		t.Fatal("expected at least one Kimi model")
	}
	if models[0].ID == "" {
		t.Fatal("expected model id to be set")
	}
}

func TestFallbackModels_UnknownProviderReturnsEmpty(t *testing.T) {
	t.Parallel()

	models := FallbackModels("unknown")
	if len(models) != 0 {
		t.Fatalf("expected empty fallback catalog, got %v", models)
	}
}

func TestModelOptionLabel(t *testing.T) {
	t.Parallel()

	option := ModelOption{ID: "k2.5", Name: "Kimi K2.5"}
	if got := option.Label(); got != "Kimi K2.5 (k2.5)" {
		t.Fatalf("Label() = %q", got)
	}

	option = ModelOption{ID: "gpt-5.4", Name: "GPT-5.4", SupportsImageInput: true, SupportsTools: true, IsFree: true}
	if got := option.Label(); got != "GPT-5.4 (gpt-5.4) [vision, tools, free]" {
		t.Fatalf("Label() with badges = %q", got)
	}
}

func TestListAnthropicModels(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/models" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if got := r.Header.Get("x-api-key"); got != "secret" {
			t.Fatalf("x-api-key = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"data": [
				{"id":"claude-sonnet-4-6","display_name":"Claude Sonnet 4.6","created_at":"2026-01-01T00:00:00Z","type":"model"},
				{"id":"claude-haiku-4-5","display_name":"Claude Haiku 4.5","created_at":"2025-10-01T00:00:00Z","type":"model"}
			],
			"has_more": false,
			"first_id": "claude-sonnet-4-6",
			"last_id": "claude-haiku-4-5"
		}`))
	}))
	defer server.Close()

	models, err := listAnthropicModels(context.Background(), "secret", option.WithBaseURL(server.URL+"/"))
	if err != nil {
		t.Fatalf("listAnthropicModels() error = %v", err)
	}
	if len(models) != 2 {
		t.Fatalf("expected 2 models, got %d", len(models))
	}
	if models[0].ID != "claude-sonnet-4-6" {
		t.Fatalf("first model = %+v", models[0])
	}
}

func TestListGoogleModels(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if got := r.URL.Query().Get("key"); got != "secret" {
			t.Fatalf("key = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"models": [
				{
					"name": "models/gemini-2.5-flash",
					"displayName": "Gemini 2.5 Flash",
					"supportedGenerationMethods": ["generateContent"]
				},
				{
					"name": "models/text-embedding-004",
					"displayName": "Text Embedding 004",
					"supportedGenerationMethods": ["embedContent"]
				}
			]
		}`))
	}))
	defer server.Close()

	models, err := listGoogleModels(context.Background(), "secret", server.URL, server.Client())
	if err != nil {
		t.Fatalf("listGoogleModels() error = %v", err)
	}
	if len(models) != 1 {
		t.Fatalf("expected 1 model, got %d", len(models))
	}
	if models[0].ID != "gemini-2.5-flash" {
		t.Fatalf("first model = %+v", models[0])
	}
}

func TestListModels_KiloUsesRemoteCatalog(t *testing.T) {
	originalClient := http.DefaultClient
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/models" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer secret" {
			t.Fatalf("Authorization = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"data": [
				{"id":"openai/gpt-5.4","name":"GPT-5.4","owned_by":"openai","architecture":{"input_modalities":["text","image"]}},
				{"id":"z-ai/glm-5-turbo:free","name":"GLM 5 Turbo","owned_by":"z-ai","architecture":{"input_modalities":["text"]}}
			]
		}`))
	}))
	defer server.Close()
	http.DefaultClient = server.Client()
	defer func() { http.DefaultClient = originalClient }()

	originalURL := kiloModelsURLForTest(server.URL + "/models")
	defer originalURL()

	models, err := ListModels(context.Background(), "kilo", ModelCatalogCredentials{KiloAPIKey: "secret"})
	if err != nil {
		t.Fatalf("ListModels() error = %v", err)
	}
	if len(models) != 2 {
		t.Fatalf("expected 2 models, got %d", len(models))
	}
	if models[0].ID != "openai/gpt-5.4" {
		t.Fatalf("first model = %+v", models[0])
	}
	if models[0].Name != "GPT-5.4 · openai" {
		t.Fatalf("first model name = %q", models[0].Name)
	}
	if !models[0].SupportsImageInput {
		t.Fatal("expected first model to support image input")
	}
	if models[1].SupportsImageInput {
		t.Fatal("expected second model to be text-only")
	}
	if !models[1].IsFree {
		t.Fatal("expected second model to be marked as free")
	}
}

func TestListModels_OpenAICodexUsesFallbackCatalog(t *testing.T) {
	t.Parallel()

	models, err := ListModels(context.Background(), "openai", ModelCatalogCredentials{OpenAIAuthMode: "codex"})
	if err != nil {
		t.Fatalf("ListModels() error = %v", err)
	}
	if len(models) == 0 {
		t.Fatal("expected codex fallback catalog")
	}
	if models[0].ID != "gpt-5.4" {
		t.Fatalf("first model = %+v", models[0])
	}
}

func TestFallbackModels_ExposeVisionCapabilityLabels(t *testing.T) {
	t.Parallel()

	openAI := FallbackModels("openai")
	if len(openAI) == 0 || !openAI[0].SupportsImageInput {
		t.Fatalf("expected OpenAI fallback to advertise vision support: %+v", openAI)
	}

	zai := FallbackModels("zai")
	foundVision := false
	for _, model := range zai {
		if model.SupportsImageInput {
			foundVision = true
			break
		}
	}
	if !foundVision {
		t.Fatalf("expected at least one Z.ai fallback model to advertise vision support: %+v", zai)
	}
}
