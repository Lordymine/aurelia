package config

// legacyFileConfig supports reading the old flat-key JSON format.
type legacyFileConfig struct {
	LLMProvider            string  `json:"llm_provider"`
	LLMModel               string  `json:"llm_model"`
	STTProvider            string  `json:"stt_provider"`
	TelegramBotToken       string  `json:"telegram_bot_token"`
	TelegramAllowedUserIDs []int64 `json:"telegram_allowed_user_ids"`
	AnthropicAPIKey        string  `json:"anthropic_api_key"`
	GoogleAPIKey           string  `json:"google_api_key"`
	KiloAPIKey             string  `json:"kilo_api_key"`
	KimiAPIKey             string  `json:"kimi_api_key"`
	OpenRouterAPIKey       string  `json:"openrouter_api_key"`
	ZAIAPIKey              string  `json:"zai_api_key"`
	AlibabaAPIKey          string  `json:"alibaba_api_key"`
	OpenAIAPIKey           string  `json:"openai_api_key"`
	OpenAIAuthMode         string  `json:"openai_auth_mode"`
	GroqAPIKey             string  `json:"groq_api_key"`
	MaxIterations          int     `json:"max_iterations"`
	DBPath                 string  `json:"db_path"`
	MCPConfigPath          string  `json:"mcp_servers_config_path"`
}

// migrateLegacy converts a legacy flat-key config to the new schema.
func migrateLegacy(legacy legacyFileConfig) fileConfig {
	providers := make(map[string]ProviderConfig)

	maybeSet := func(name, key string) {
		if key != "" {
			providers[name] = ProviderConfig{APIKey: key}
		}
	}

	maybeSet("anthropic", legacy.AnthropicAPIKey)
	maybeSet("google", legacy.GoogleAPIKey)
	maybeSet("kilo", legacy.KiloAPIKey)
	maybeSet("kimi", legacy.KimiAPIKey)
	maybeSet("openrouter", legacy.OpenRouterAPIKey)
	maybeSet("zai", legacy.ZAIAPIKey)
	maybeSet("alibaba", legacy.AlibabaAPIKey)
	maybeSet("groq", legacy.GroqAPIKey)

	if legacy.OpenAIAPIKey != "" || legacy.OpenAIAuthMode != "" {
		providers["openai"] = ProviderConfig{
			APIKey:   legacy.OpenAIAPIKey,
			AuthMode: legacy.OpenAIAuthMode,
		}
	}

	return fileConfig{
		DefaultProvider:        legacy.LLMProvider,
		DefaultModel:           legacy.LLMModel,
		Providers:              providers,
		TelegramBotToken:       legacy.TelegramBotToken,
		TelegramAllowedUserIDs: legacy.TelegramAllowedUserIDs,
		STTProvider:            legacy.STTProvider,
		MaxIterations:          legacy.MaxIterations,
		DBPath:                 legacy.DBPath,
		MCPConfigPath:          legacy.MCPConfigPath,
	}
}
