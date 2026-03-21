package llm

import "strings"

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

// Providers returns a copy of the available provider specs.
func Providers() []ProviderSpec {
	specs := make([]ProviderSpec, len(providerSpecs))
	copy(specs, providerSpecs)
	return specs
}

// Provider returns the spec for the given provider name.
func Provider(provider string) (ProviderSpec, bool) {
	normalized := NormalizeProvider(provider)
	for _, spec := range providerSpecs {
		if spec.ID == normalized {
			return spec, true
		}
	}
	return ProviderSpec{}, false
}

// NormalizeProvider returns a canonical lowercase provider name.
func NormalizeProvider(provider string) string {
	normalized := strings.TrimSpace(strings.ToLower(provider))
	if normalized == "" {
		return "kimi"
	}
	return normalized
}

// DefaultModelForProvider returns the default model for the given provider.
func DefaultModelForProvider(provider string) string {
	spec, ok := Provider(provider)
	if !ok {
		return "kimi-k2-thinking"
	}
	return spec.DefaultModel
}

// ProviderChoices returns the list of provider IDs.
func ProviderChoices() []string {
	specs := Providers()
	choices := make([]string, 0, len(specs))
	for _, spec := range specs {
		choices = append(choices, spec.ID)
	}
	return choices
}

// ProviderLabels returns the list of provider display labels.
func ProviderLabels() []string {
	specs := Providers()
	labels := make([]string, 0, len(specs))
	for _, spec := range specs {
		labels = append(labels, spec.Label)
	}
	return labels
}
