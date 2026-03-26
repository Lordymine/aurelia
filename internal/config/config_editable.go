package config

import (
	"github.com/kocar/aurelia/internal/runtime"
)

// EditableConfig represents the user-editable portion of the runtime config.
// Keeps flat per-provider fields for backward compatibility with onboarding UI.
type EditableConfig struct {
	LLMProvider            string
	LLMModel               string
	STTProvider            string
	TelegramBotToken       string
	TelegramAllowedUserIDs []int64
	AnthropicAPIKey        string
	GoogleAPIKey           string
	KiloAPIKey             string
	KimiAPIKey             string
	OpenRouterAPIKey       string
	ZAIAPIKey              string
	AlibabaAPIKey          string
	AnthropicAuthMode      string
	GroqAPIKey             string
	MaxIterations          int
}

func (c EditableConfig) LLMAPIKey(provider string) string {
	switch NormalizeProvider(provider) {
	case "anthropic":
		return c.AnthropicAPIKey
	case "google":
		return c.GoogleAPIKey
	case "kilo":
		return c.KiloAPIKey
	case "openrouter":
		return c.OpenRouterAPIKey
	case "zai":
		return c.ZAIAPIKey
	case "alibaba":
		return c.AlibabaAPIKey
	default:
		return c.KimiAPIKey
	}
}

func (c *EditableConfig) SetLLMAPIKey(provider, value string) {
	switch NormalizeProvider(provider) {
	case "anthropic":
		c.AnthropicAPIKey = value
	case "google":
		c.GoogleAPIKey = value
	case "kilo":
		c.KiloAPIKey = value
	case "openrouter":
		c.OpenRouterAPIKey = value
	case "zai":
		c.ZAIAPIKey = value
	case "alibaba":
		c.AlibabaAPIKey = value
	default:
		c.KimiAPIKey = value
	}
}

// DefaultEditableConfig returns the default user-editable configuration values.
func DefaultEditableConfig() EditableConfig {
	return EditableConfig{
		LLMProvider:            defaultLLMProvider,
		LLMModel:               defaultModelForProvider(defaultLLMProvider),
		AnthropicAuthMode:      "api_key",
		STTProvider:            defaultSTTProvider,
		TelegramAllowedUserIDs: []int64{},
		MaxIterations:          defaultMaxIterations,
	}
}

// LoadEditable returns the editable config subset from the current app config.
func LoadEditable(r *runtime.PathResolver) (*EditableConfig, error) {
	cfg, err := Load(r)
	if err != nil {
		return nil, err
	}
	return appConfigToEditable(cfg), nil
}

// appConfigToEditable converts AppConfig to the flat EditableConfig.
func appConfigToEditable(cfg *AppConfig) *EditableConfig {
	anthropicAuthMode := cfg.ProviderAuthMode("anthropic")
	if anthropicAuthMode == "" {
		anthropicAuthMode = "api_key"
	}
	return &EditableConfig{
		LLMProvider:            cfg.DefaultProvider,
		LLMModel:               cfg.DefaultModel,
		STTProvider:            cfg.STTProvider,
		TelegramBotToken:       cfg.TelegramBotToken,
		TelegramAllowedUserIDs: append([]int64(nil), cfg.TelegramAllowedUserIDs...),
		AnthropicAuthMode:      anthropicAuthMode,
		AnthropicAPIKey:        cfg.ProviderAPIKey("anthropic"),
		GoogleAPIKey:           cfg.ProviderAPIKey("google"),
		KiloAPIKey:             cfg.ProviderAPIKey("kilo"),
		KimiAPIKey:             cfg.ProviderAPIKey("kimi"),
		OpenRouterAPIKey:       cfg.ProviderAPIKey("openrouter"),
		ZAIAPIKey:              cfg.ProviderAPIKey("zai"),
		AlibabaAPIKey:          cfg.ProviderAPIKey("alibaba"),
		GroqAPIKey:             cfg.ProviderAPIKey("groq"),
		MaxIterations:          cfg.MaxIterations,
	}
}

// SaveEditable updates the user-editable config subset while preserving managed paths.
func SaveEditable(r *runtime.PathResolver, editable EditableConfig) error {
	cfg := editableToFileConfig(editable)
	normalized := normalizeFileConfig(cfg, r)
	return writeConfigFile(r.AppConfig(), normalized)
}

// editableToFileConfig converts the flat EditableConfig to the new fileConfig.
func editableToFileConfig(editable EditableConfig) fileConfig {
	providers := make(map[string]ProviderConfig)

	maybeSet := func(name, key string) {
		if key != "" {
			providers[name] = ProviderConfig{APIKey: key}
		}
	}

	if editable.AnthropicAPIKey != "" || editable.AnthropicAuthMode != "" {
		providers["anthropic"] = ProviderConfig{
			APIKey:   editable.AnthropicAPIKey,
			AuthMode: editable.AnthropicAuthMode,
		}
	}
	maybeSet("google", editable.GoogleAPIKey)
	maybeSet("kilo", editable.KiloAPIKey)
	maybeSet("kimi", editable.KimiAPIKey)
	maybeSet("openrouter", editable.OpenRouterAPIKey)
	maybeSet("zai", editable.ZAIAPIKey)
	maybeSet("alibaba", editable.AlibabaAPIKey)
	maybeSet("groq", editable.GroqAPIKey)

	return fileConfig{
		DefaultProvider:        editable.LLMProvider,
		DefaultModel:           editable.LLMModel,
		Providers:              providers,
		STTProvider:            editable.STTProvider,
		TelegramBotToken:       editable.TelegramBotToken,
		TelegramAllowedUserIDs: append([]int64(nil), editable.TelegramAllowedUserIDs...),
		MaxIterations:          editable.MaxIterations,
	}
}

func sameFileConfig(a, b fileConfig) bool {
	if a.TelegramBotToken != b.TelegramBotToken ||
		a.DefaultProvider != b.DefaultProvider ||
		a.DefaultModel != b.DefaultModel ||
		a.STTProvider != b.STTProvider ||
		a.MaxIterations != b.MaxIterations ||
		a.DBPath != b.DBPath ||
		a.MCPConfigPath != b.MCPConfigPath {
		return false
	}
	if len(a.TelegramAllowedUserIDs) != len(b.TelegramAllowedUserIDs) {
		return false
	}
	for i := range a.TelegramAllowedUserIDs {
		if a.TelegramAllowedUserIDs[i] != b.TelegramAllowedUserIDs[i] {
			return false
		}
	}
	if len(a.Providers) != len(b.Providers) {
		return false
	}
	for k, v := range a.Providers {
		bv, ok := b.Providers[k]
		if !ok || v != bv {
			return false
		}
	}
	return true
}
