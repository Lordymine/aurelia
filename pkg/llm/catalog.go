package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// ModelOption is a selectable model entry for onboarding and config UIs.
type ModelOption struct {
	ID                 string
	Name               string
	SupportsImageInput bool
	SupportsTools      bool
	IsFree             bool
}

func (m ModelOption) Label() string {
	badges := make([]string, 0, 3)
	if m.SupportsImageInput {
		badges = append(badges, "vision")
	}
	if m.SupportsTools {
		badges = append(badges, "tools")
	}
	if m.IsFree {
		badges = append(badges, "free")
	}
	suffix := ""
	if len(badges) != 0 {
		suffix = " [" + strings.Join(badges, ", ") + "]"
	}
	if m.Name == "" || m.Name == m.ID {
		return m.ID + suffix
	}
	return fmt.Sprintf("%s (%s)%s", m.Name, m.ID, suffix)
}

// ModelCatalogCredentials carries provider-specific credentials used by model catalogs.
type ModelCatalogCredentials struct {
	AnthropicAPIKey  string
	GoogleAPIKey     string
	KiloAPIKey       string
	KimiAPIKey       string
	OpenRouterAPIKey string
	ZAIAPIKey        string
	AlibabaAPIKey    string
	OpenAIAPIKey     string
	OpenAIAuthMode   string
}

const googleModelsURL = "https://generativelanguage.googleapis.com/v1beta/models"

// ListModels returns the available model options for a provider.
// Providers may use remote discovery or curated fallbacks internally.
func ListModels(ctx context.Context, provider string, creds ModelCatalogCredentials) ([]ModelOption, error) {
	_ = ctx
	_ = creds

	provider = NormalizeProvider(provider)

	switch provider {
	case "anthropic":
		if creds.AnthropicAPIKey != "" {
			return listAnthropicModels(ctx, creds.AnthropicAPIKey)
		}
		return fallbackModels("anthropic"), nil
	case "google":
		if creds.GoogleAPIKey != "" {
			return listGoogleModels(ctx, creds.GoogleAPIKey, googleModelsURL, http.DefaultClient)
		}
		return fallbackModels("google"), nil
	case "kilo":
		models, err := listKiloModels(ctx, creds.KiloAPIKey, kiloModelsCatalogURL, http.DefaultClient)
		if err == nil && len(models) != 0 {
			return models, nil
		}
		return fallbackModels("kilo"), nil
	case "openrouter":
		return listOpenRouterModels(ctx, creds.OpenRouterAPIKey, openRouterModelsURL, http.DefaultClient)
	case "zai":
		return fallbackModels("zai"), nil
	case "alibaba":
		return fallbackModels("alibaba"), nil
	case "openai":
		if creds.OpenAIAuthMode == "codex" {
			return fallbackModels("openai_codex"), nil
		}
		if creds.OpenAIAPIKey != "" {
			return listOpenAIModels(ctx, creds.OpenAIAPIKey, openAIModelsURL, http.DefaultClient)
		}
		return fallbackModels("openai"), nil
	case "kimi":
		return fallbackModels("kimi"), nil
	default:
		return nil, fmt.Errorf("unsupported llm provider %q", provider)
	}
}

// FallbackModels returns curated default models when discovery is unavailable.
func FallbackModels(provider string) []ModelOption {
	return fallbackModels(NormalizeProvider(provider))
}

func fallbackModels(provider string) []ModelOption {
	switch provider {
	case "anthropic":
		return []ModelOption{
			{ID: "claude-sonnet-4-6", Name: "Claude Sonnet 4.6", SupportsImageInput: true},
			{ID: "claude-opus-4-6", Name: "Claude Opus 4.6", SupportsImageInput: true},
			{ID: "claude-haiku-4-5", Name: "Claude Haiku 4.5", SupportsImageInput: true},
		}
	case "google":
		return []ModelOption{
			{ID: "gemini-2.5-pro", Name: "Gemini 2.5 Pro", SupportsImageInput: true},
			{ID: "gemini-2.5-flash", Name: "Gemini 2.5 Flash", SupportsImageInput: true},
			{ID: "gemini-2.5-flash-lite", Name: "Gemini 2.5 Flash-Lite", SupportsImageInput: true},
		}
	case "kilo":
		return []ModelOption{
			{ID: "openai/gpt-5.4", Name: "OpenAI: GPT-5.4", SupportsImageInput: true},
			{ID: "anthropic/claude-sonnet-4.6", Name: "Anthropic: Claude Sonnet 4.6", SupportsImageInput: true},
			{ID: "google/gemini-3.1-pro-preview", Name: "Google: Gemini 3.1 Pro Preview", SupportsImageInput: true},
			{ID: "zai/glm-4.6v", Name: "Z.ai: GLM 4.6V", SupportsImageInput: true},
			{ID: "zai/glm-5-turbo", Name: "Z.ai: GLM 5 Turbo"},
		}
	case "openrouter":
		return []ModelOption{
			{ID: "openrouter/auto", Name: "OpenRouter Auto"},
			{ID: "openrouter/free", Name: "OpenRouter Free Router", IsFree: true},
		}
	case "zai":
		return []ModelOption{
			{ID: "glm-5", Name: "GLM-5"},
			{ID: "glm-4.7", Name: "GLM-4.7"},
			{ID: "glm-4.6v", Name: "GLM-4.6V", SupportsImageInput: true},
			{ID: "glm-4.5-air", Name: "GLM-4.5 Air"},
		}
	case "alibaba":
		return []ModelOption{
			{ID: "qwen3-coder-plus", Name: "Qwen3 Coder Plus"},
			{ID: "qwen3-coder-next", Name: "Qwen3 Coder Next"},
			{ID: "qwen-vl-max", Name: "Qwen VL Max", SupportsImageInput: true},
			{ID: "qwen3.5-plus", Name: "Qwen3.5 Plus"},
		}
	case "openai":
		return []ModelOption{
			{ID: "gpt-5.4", Name: "GPT-5.4", SupportsImageInput: true, SupportsTools: true},
			{ID: "gpt-5-mini", Name: "GPT-5 mini", SupportsImageInput: true, SupportsTools: true},
			{ID: "o4-mini", Name: "o4-mini", SupportsImageInput: true, SupportsTools: true},
		}
	case "openai_codex":
		return []ModelOption{
			{ID: "gpt-5.4", Name: "GPT-5.4", SupportsImageInput: true, SupportsTools: true},
			{ID: "gpt-5-mini", Name: "GPT-5 mini", SupportsImageInput: true, SupportsTools: true},
			{ID: "gpt-5.2-codex", Name: "GPT-5.2-Codex"},
			{ID: "o4-mini", Name: "o4-mini", SupportsImageInput: true, SupportsTools: true},
		}
	case "kimi":
		return []ModelOption{
			{ID: "kimi-k2-thinking", Name: "Kimi K2 Thinking"},
			{ID: "kimi-k2-thinking-turbo", Name: "Kimi K2 Thinking Turbo"},
			{ID: "k2.5", Name: "Kimi K2.5"},
			{ID: "moonshot-v1-vision", Name: "Moonshot Vision", SupportsImageInput: true},
			{ID: "moonshot-v1-8k", Name: "Moonshot v1 8K"},
			{ID: "moonshot-v1-32k", Name: "Moonshot v1 32K"},
			{ID: "moonshot-v1-128k", Name: "Moonshot v1 128K"},
		}
	default:
		return nil
	}
}

func listAnthropicModels(ctx context.Context, apiKey string, opts ...option.RequestOption) ([]ModelOption, error) {
	requestOptions := []option.RequestOption{option.WithAPIKey(apiKey)}
	requestOptions = append(requestOptions, opts...)

	client := anthropic.NewClient(requestOptions...)
	pager := client.Models.ListAutoPaging(ctx, anthropic.ModelListParams{})

	var models []ModelOption
	for pager.Next() {
		model := pager.Current()
		models = append(models, ModelOption{
			ID:                 model.ID,
			Name:               model.DisplayName,
			SupportsImageInput: true,
		})
	}
	if err := pager.Err(); err != nil {
		return nil, err
	}
	return models, nil
}

func listGoogleModels(ctx context.Context, apiKey string, baseURL string, client *http.Client) ([]ModelOption, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"?key="+apiKey, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google models list returned %s", resp.Status)
	}

	var payload struct {
		Models []struct {
			Name                      string   `json:"name"`
			DisplayName               string   `json:"displayName"`
			SupportedGenerationMethod []string `json:"supportedGenerationMethods"`
		} `json:"models"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	var models []ModelOption
	for _, model := range payload.Models {
		id := strings.TrimPrefix(model.Name, "models/")
		if !supportsGenerateContent(model.SupportedGenerationMethod) || !looksLikeGoogleChatModel(id) {
			continue
		}
		models = append(models, ModelOption{
			ID:                 id,
			Name:               model.DisplayName,
			SupportsImageInput: true,
		})
	}
	sort.SliceStable(models, func(i, j int) bool {
		return rankGoogleModel(models[i].ID) < rankGoogleModel(models[j].ID)
	})
	return models, nil
}

func supportsGenerateContent(methods []string) bool {
	for _, method := range methods {
		if method == "generateContent" {
			return true
		}
	}
	return false
}

func looksLikeGoogleChatModel(modelID string) bool {
	if !strings.HasPrefix(modelID, "gemini-") {
		return false
	}
	blocked := []string{"preview", "exp", "tts", "live", "image", "embedding", "aqa"}
	for _, token := range blocked {
		if strings.Contains(modelID, token) {
			return false
		}
	}
	return true
}

func rankGoogleModel(modelID string) int {
	switch modelID {
	case "gemini-2.5-pro":
		return 0
	case "gemini-2.5-flash":
		return 1
	case "gemini-2.5-flash-lite":
		return 2
	default:
		return 10
	}
}
