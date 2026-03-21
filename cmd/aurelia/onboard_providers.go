package main

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/kocar/aurelia/internal/config"
)

// ModelOption is a selectable model entry for onboarding and config UIs.
type ModelOption struct {
	ID                 string
	Name               string
	SupportsImageInput bool
	SupportsTools      bool
	IsFree             bool
}

// Label returns a display label for the model option.
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

// ProviderSpec describes a supported LLM provider for onboarding and config.
type ProviderSpec struct {
	ID                  string
	Label               string
	DefaultModel        string
	APIKeyLabel         string
	APIKeyHelp          string
	SupportsModelSearch bool
	AuthModes           []string
}

var providerSpecs = []ProviderSpec{
	{
		ID:           "kimi",
		Label:        "Kimi",
		DefaultModel: "kimi-k2-thinking",
		APIKeyLabel:  "Kimi API key",
		APIKeyHelp:   "Used for the main LLM runtime.",
	},
	{
		ID:           "anthropic",
		Label:        "Anthropic",
		DefaultModel: "claude-sonnet-4-6",
		APIKeyLabel:  "Anthropic API key",
		APIKeyHelp:   "Used for the Anthropic LLM runtime.",
	},
	{
		ID:           "google",
		Label:        "Google",
		DefaultModel: "gemini-2.5-pro",
		APIKeyLabel:  "Google API key",
		APIKeyHelp:   "Used for the Google Gemini LLM runtime.",
	},
	{
		ID:                  "kilo",
		Label:               "Kilo Code",
		DefaultModel:        "openai/gpt-5.4",
		APIKeyLabel:         "Kilo API key",
		APIKeyHelp:          "Used for the Kilo Gateway LLM runtime.",
		SupportsModelSearch: true,
	},
	{
		ID:                  "openrouter",
		Label:               "OpenRouter",
		DefaultModel:        "openrouter/auto",
		APIKeyLabel:         "OpenRouter API key",
		APIKeyHelp:          "Used for the OpenRouter LLM runtime.",
		SupportsModelSearch: true,
	},
	{
		ID:           "zai",
		Label:        "Z.ai",
		DefaultModel: "glm-5",
		APIKeyLabel:  "Z.ai Coding Plan API key",
		APIKeyHelp:   "Used for the Z.ai GLM Coding Plan runtime.",
	},
	{
		ID:           "alibaba",
		Label:        "Alibaba",
		DefaultModel: "qwen3-coder-plus",
		APIKeyLabel:  "Alibaba Coding Plan API key",
		APIKeyHelp:   "Used for the Alibaba Coding Plan runtime.",
	},
	{
		ID:           "openai",
		Label:        "OpenAI",
		DefaultModel: "gpt-5.4",
		APIKeyLabel:  "OpenAI API key",
		APIKeyHelp:   "Used for the OpenAI LLM runtime.",
		AuthModes:    []string{"api_key", "codex"},
	},
}

// providers returns a copy of the available provider specs.
func providers() []ProviderSpec {
	specs := make([]ProviderSpec, len(providerSpecs))
	copy(specs, providerSpecs)
	return specs
}

// provider returns the spec for the given provider name.
func provider(name string) (ProviderSpec, bool) {
	normalized := config.NormalizeProvider(name)
	for _, spec := range providerSpecs {
		if spec.ID == normalized {
			return spec, true
		}
	}
	return ProviderSpec{}, false
}

// defaultModelForProvider returns the default model for the given provider.
func defaultModelForProvider(p string) string {
	spec, ok := provider(p)
	if !ok {
		return "kimi-k2-thinking"
	}
	return spec.DefaultModel
}

// providerChoices returns the list of provider IDs.
func providerChoices() []string {
	specs := providers()
	choices := make([]string, 0, len(specs))
	for _, spec := range specs {
		choices = append(choices, spec.ID)
	}
	return choices
}

// providerLabels returns the list of provider display labels.
func providerLabels() []string {
	specs := providers()
	labels := make([]string, 0, len(specs))
	for _, spec := range specs {
		labels = append(labels, spec.Label)
	}
	return labels
}

// listModels returns the available model options for a provider.
// Remote discovery was removed — returns curated fallback models only.
func listModels(_ context.Context, p string, creds ModelCatalogCredentials) ([]ModelOption, error) {
	_ = creds

	p = config.NormalizeProvider(p)
	if p == "openai" && creds.OpenAIAuthMode == "codex" {
		return fallbackModels("openai_codex"), nil
	}

	models := fallbackModels(p)
	if models == nil {
		return nil, fmt.Errorf("unsupported llm provider %q", p)
	}
	return models, nil
}

// fallbackModelList returns curated default models when discovery is unavailable.
func fallbackModelList(p string) []ModelOption {
	return fallbackModels(config.NormalizeProvider(p))
}

func fallbackModels(p string) []ModelOption {
	switch p {
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

var codexLookPath = exec.LookPath

// ensureCodexCLIAvailable checks that the codex CLI binary is reachable.
func ensureCodexCLIAvailable() error {
	_, err := codexLookPath("codex")
	return err
}
