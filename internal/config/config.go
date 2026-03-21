package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kocar/aurelia/internal/runtime"
)

// NormalizeProvider returns a canonical lowercase provider name.
// Duplicated from cmd/aurelia to keep internal/config free of pkg/llm.
func NormalizeProvider(provider string) string {
	normalized := strings.TrimSpace(strings.ToLower(provider))
	if normalized == "" {
		return "kimi"
	}
	return normalized
}

// defaultModelForProvider returns the default model for the given provider.
func defaultModelForProvider(provider string) string {
	switch NormalizeProvider(provider) {
	case "anthropic":
		return "claude-sonnet-4-6"
	case "google":
		return "gemini-2.5-pro"
	case "kilo":
		return "openai/gpt-5.4"
	case "openrouter":
		return "openrouter/auto"
	case "zai":
		return "glm-5"
	case "alibaba":
		return "qwen3-coder-plus"
	case "openai":
		return "gpt-5.4"
	default:
		return "kimi-k2-thinking"
	}
}

const (
	defaultMaxIterations    = 500
	defaultMemoryWindowSize = 20
	defaultLLMProvider      = "kimi"
	defaultLLMModel         = "kimi-k2-thinking"
	defaultSTTProvider      = "groq"
)

// ProviderConfig holds credentials and endpoint for a single LLM provider.
type ProviderConfig struct {
	APIKey   string `json:"api_key"`
	BaseURL  string `json:"base_url,omitempty"`
	AuthMode string `json:"auth_mode,omitempty"`
}

// AppConfig holds all runtime configuration needed for the application.
type AppConfig struct {
	DefaultProvider string                    `json:"default_provider"`
	DefaultModel    string                    `json:"default_model"`
	Providers       map[string]ProviderConfig `json:"providers"`

	TelegramBotToken       string  `json:"telegram_bot_token"`
	TelegramAllowedUserIDs []int64 `json:"telegram_allowed_user_ids"`

	EmbeddingProvider string `json:"embedding_provider"`
	EmbeddingModel    string `json:"embedding_model"`
	EmbeddingAPIKey   string `json:"embedding_api_key,omitempty"`

	STTProvider string `json:"stt_provider"`

	MaxIterations    int    `json:"max_iterations"`
	DBPath           string `json:"db_path"`
	MemoryWindowSize int    `json:"memory_window_size"`
	MCPConfigPath    string `json:"mcp_servers_config_path"`
}

// ProviderAPIKey returns the API key for the given provider, or empty string.
func (c *AppConfig) ProviderAPIKey(provider string) string {
	p, ok := c.Providers[NormalizeProvider(provider)]
	if !ok {
		return ""
	}
	return p.APIKey
}

// ProviderBaseURL returns the base URL for the given provider, or empty string.
func (c *AppConfig) ProviderBaseURL(provider string) string {
	p, ok := c.Providers[NormalizeProvider(provider)]
	if !ok {
		return ""
	}
	return p.BaseURL
}

// ProviderAuthMode returns the auth mode for the given provider, or empty string.
func (c *AppConfig) ProviderAuthMode(provider string) string {
	p, ok := c.Providers[NormalizeProvider(provider)]
	if !ok {
		return ""
	}
	return p.AuthMode
}

// fileConfig is the JSON structure written to disk (new schema).
type fileConfig struct {
	DefaultProvider string                    `json:"default_provider"`
	DefaultModel    string                    `json:"default_model"`
	Providers       map[string]ProviderConfig `json:"providers"`

	TelegramBotToken       string  `json:"telegram_bot_token"`
	TelegramAllowedUserIDs []int64 `json:"telegram_allowed_user_ids"`

	EmbeddingProvider string `json:"embedding_provider,omitempty"`
	EmbeddingModel    string `json:"embedding_model,omitempty"`
	EmbeddingAPIKey   string `json:"embedding_api_key,omitempty"`

	STTProvider string `json:"stt_provider"`

	MaxIterations    int    `json:"max_iterations"`
	DBPath           string `json:"db_path"`
	MemoryWindowSize int    `json:"memory_window_size"`
	MCPConfigPath    string `json:"mcp_servers_config_path"`
}

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
	MemoryWindowSize       int     `json:"memory_window_size"`
	MCPConfigPath          string  `json:"mcp_servers_config_path"`
}

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
	OpenAIAPIKey           string
	OpenAIAuthMode         string
	GroqAPIKey             string
	MaxIterations          int
	MemoryWindowSize       int
	EmbeddingProvider      string
	EmbeddingModel         string
	EmbeddingAPIKey        string
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
	case "openai":
		return c.OpenAIAPIKey
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
	case "openai":
		c.OpenAIAPIKey = value
	default:
		c.KimiAPIKey = value
	}
}

// Load reads the instance-local JSON config, creates it with defaults when
// missing, and returns the normalized runtime config.
func Load(r *runtime.PathResolver) (*AppConfig, error) {
	path := r.AppConfig()
	defaults := defaultFileConfig(r)

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			if err := writeConfigFile(path, defaults); err != nil {
				return nil, err
			}
			return toAppConfig(defaults), nil
		}
		return nil, fmt.Errorf("stat app config: %w", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read app config: %w", err)
	}

	cfg := defaults
	if len(data) != 0 {
		// Try new schema first
		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("decode app config: %w", err)
		}

		// Detect legacy format: if providers map is empty but legacy fields present
		if len(cfg.Providers) == 0 {
			var legacy legacyFileConfig
			if err := json.Unmarshal(data, &legacy); err == nil {
				cfg = migrateLegacy(legacy, r)
			}
		}
	}

	normalized := normalizeFileConfig(cfg, r)
	if !sameFileConfig(normalized, cfg) {
		if err := writeConfigFile(path, normalized); err != nil {
			return nil, err
		}
	}

	return toAppConfig(normalized), nil
}

// migrateLegacy converts a legacy flat-key config to the new schema.
func migrateLegacy(legacy legacyFileConfig, r *runtime.PathResolver) fileConfig {
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
		MemoryWindowSize:       legacy.MemoryWindowSize,
		MCPConfigPath:          legacy.MCPConfigPath,
	}
}

func defaultFileConfig(r *runtime.PathResolver) fileConfig {
	return fileConfig{
		DefaultProvider:        defaultLLMProvider,
		DefaultModel:           defaultModelForProvider(defaultLLMProvider),
		Providers:              map[string]ProviderConfig{},
		STTProvider:            defaultSTTProvider,
		TelegramAllowedUserIDs: []int64{},
		MaxIterations:          defaultMaxIterations,
		DBPath:                 filepath.Join(r.Data(), "aurelia.db"),
		MemoryWindowSize:       defaultMemoryWindowSize,
		MCPConfigPath:          filepath.Join(r.Config(), "mcp_servers.json"),
	}
}

// DefaultEditableConfig returns the default user-editable configuration values.
func DefaultEditableConfig() EditableConfig {
	return EditableConfig{
		LLMProvider:            defaultLLMProvider,
		LLMModel:               defaultModelForProvider(defaultLLMProvider),
		OpenAIAuthMode:         "api_key",
		STTProvider:            defaultSTTProvider,
		TelegramAllowedUserIDs: []int64{},
		MaxIterations:          defaultMaxIterations,
		MemoryWindowSize:       defaultMemoryWindowSize,
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
	openAIAuthMode := cfg.ProviderAuthMode("openai")
	if openAIAuthMode == "" {
		openAIAuthMode = "api_key"
	}
	return &EditableConfig{
		LLMProvider:            cfg.DefaultProvider,
		LLMModel:               cfg.DefaultModel,
		STTProvider:            cfg.STTProvider,
		TelegramBotToken:       cfg.TelegramBotToken,
		TelegramAllowedUserIDs: append([]int64(nil), cfg.TelegramAllowedUserIDs...),
		AnthropicAPIKey:        cfg.ProviderAPIKey("anthropic"),
		GoogleAPIKey:           cfg.ProviderAPIKey("google"),
		KiloAPIKey:             cfg.ProviderAPIKey("kilo"),
		KimiAPIKey:             cfg.ProviderAPIKey("kimi"),
		OpenRouterAPIKey:       cfg.ProviderAPIKey("openrouter"),
		ZAIAPIKey:              cfg.ProviderAPIKey("zai"),
		AlibabaAPIKey:          cfg.ProviderAPIKey("alibaba"),
		OpenAIAPIKey:           cfg.ProviderAPIKey("openai"),
		OpenAIAuthMode:         openAIAuthMode,
		GroqAPIKey:             cfg.ProviderAPIKey("groq"),
		MaxIterations:          cfg.MaxIterations,
		MemoryWindowSize:       cfg.MemoryWindowSize,
		EmbeddingProvider:      cfg.EmbeddingProvider,
		EmbeddingModel:         cfg.EmbeddingModel,
		EmbeddingAPIKey:        cfg.EmbeddingAPIKey,
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

	maybeSet("anthropic", editable.AnthropicAPIKey)
	maybeSet("google", editable.GoogleAPIKey)
	maybeSet("kilo", editable.KiloAPIKey)
	maybeSet("kimi", editable.KimiAPIKey)
	maybeSet("openrouter", editable.OpenRouterAPIKey)
	maybeSet("zai", editable.ZAIAPIKey)
	maybeSet("alibaba", editable.AlibabaAPIKey)
	maybeSet("groq", editable.GroqAPIKey)

	if editable.OpenAIAPIKey != "" || editable.OpenAIAuthMode != "" {
		providers["openai"] = ProviderConfig{
			APIKey:   editable.OpenAIAPIKey,
			AuthMode: editable.OpenAIAuthMode,
		}
	}

	return fileConfig{
		DefaultProvider:        editable.LLMProvider,
		DefaultModel:           editable.LLMModel,
		Providers:              providers,
		STTProvider:            editable.STTProvider,
		TelegramBotToken:       editable.TelegramBotToken,
		TelegramAllowedUserIDs: append([]int64(nil), editable.TelegramAllowedUserIDs...),
		EmbeddingProvider:      editable.EmbeddingProvider,
		EmbeddingModel:         editable.EmbeddingModel,
		EmbeddingAPIKey:        editable.EmbeddingAPIKey,
		MaxIterations:          editable.MaxIterations,
		MemoryWindowSize:       editable.MemoryWindowSize,
	}
}

func normalizeFileConfig(cfg fileConfig, r *runtime.PathResolver) fileConfig {
	defaults := defaultFileConfig(r)
	if cfg.TelegramAllowedUserIDs == nil {
		cfg.TelegramAllowedUserIDs = defaults.TelegramAllowedUserIDs
	}
	if cfg.DefaultProvider == "" {
		cfg.DefaultProvider = defaults.DefaultProvider
	}
	if cfg.DefaultModel == "" {
		cfg.DefaultModel = defaultModelForProvider(cfg.DefaultProvider)
	}
	if cfg.STTProvider == "" {
		cfg.STTProvider = defaults.STTProvider
	}
	if cfg.MaxIterations <= 0 {
		cfg.MaxIterations = defaults.MaxIterations
	}
	if cfg.DBPath == "" {
		cfg.DBPath = defaults.DBPath
	}
	if cfg.MemoryWindowSize <= 0 {
		cfg.MemoryWindowSize = defaults.MemoryWindowSize
	}
	if cfg.MCPConfigPath == "" {
		cfg.MCPConfigPath = defaults.MCPConfigPath
	}
	if cfg.Providers == nil {
		cfg.Providers = map[string]ProviderConfig{}
	}
	return cfg
}

func writeConfigFile(path string, cfg fileConfig) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create app config dir: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("encode app config: %w", err)
	}
	data = append(data, '\n')

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write app config: %w", err)
	}
	return nil
}

func toAppConfig(cfg fileConfig) *AppConfig {
	return &AppConfig{
		DefaultProvider:        cfg.DefaultProvider,
		DefaultModel:           cfg.DefaultModel,
		Providers:              cfg.Providers,
		TelegramBotToken:       cfg.TelegramBotToken,
		TelegramAllowedUserIDs: cfg.TelegramAllowedUserIDs,
		EmbeddingProvider:      cfg.EmbeddingProvider,
		EmbeddingModel:         cfg.EmbeddingModel,
		EmbeddingAPIKey:        cfg.EmbeddingAPIKey,
		STTProvider:            cfg.STTProvider,
		MaxIterations:          cfg.MaxIterations,
		DBPath:                 cfg.DBPath,
		MemoryWindowSize:       cfg.MemoryWindowSize,
		MCPConfigPath:          cfg.MCPConfigPath,
	}
}

func sameFileConfig(a, b fileConfig) bool {
	if a.TelegramBotToken != b.TelegramBotToken ||
		a.DefaultProvider != b.DefaultProvider ||
		a.DefaultModel != b.DefaultModel ||
		a.STTProvider != b.STTProvider ||
		a.EmbeddingProvider != b.EmbeddingProvider ||
		a.EmbeddingModel != b.EmbeddingModel ||
		a.EmbeddingAPIKey != b.EmbeddingAPIKey ||
		a.MaxIterations != b.MaxIterations ||
		a.DBPath != b.DBPath ||
		a.MemoryWindowSize != b.MemoryWindowSize ||
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
