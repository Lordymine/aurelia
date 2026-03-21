package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kocar/aurelia/internal/runtime"
)

// normalizeProvider returns a canonical lowercase provider name.
// Duplicated from cmd/aurelia to keep internal/config free of pkg/llm.
func normalizeProvider(provider string) string {
	normalized := strings.TrimSpace(strings.ToLower(provider))
	if normalized == "" {
		return "kimi"
	}
	return normalized
}

// defaultModelForProvider returns the default model for the given provider.
func defaultModelForProvider(provider string) string {
	switch normalizeProvider(provider) {
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

// AppConfig holds all runtime configuration needed for the application.
type AppConfig struct {
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
	DBPath                 string
	MemoryWindowSize       int
	MCPConfigPath          string
}

type fileConfig struct {
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
}

func (c EditableConfig) LLMAPIKey(provider string) string {
	switch normalizeProvider(provider) {
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
	switch normalizeProvider(provider) {
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
	defaults := DefaultFileConfig(r)

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
		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("decode app config: %w", err)
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

func defaultFileConfig(r *runtime.PathResolver) fileConfig {
	return fileConfig{
		LLMProvider:            defaultLLMProvider,
		LLMModel:               defaultModelForProvider(defaultLLMProvider),
		OpenAIAuthMode:         "api_key",
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

// DefaultFileConfig returns the full default config including instance paths.
func DefaultFileConfig(r *runtime.PathResolver) fileConfig {
	return defaultFileConfig(r)
}

// LoadEditable returns the editable config subset from the current app config.
func LoadEditable(r *runtime.PathResolver) (*EditableConfig, error) {
	cfg, err := Load(r)
	if err != nil {
		return nil, err
	}
	return &EditableConfig{
		LLMProvider:            cfg.LLMProvider,
		LLMModel:               cfg.LLMModel,
		STTProvider:            cfg.STTProvider,
		TelegramBotToken:       cfg.TelegramBotToken,
		TelegramAllowedUserIDs: append([]int64(nil), cfg.TelegramAllowedUserIDs...),
		AnthropicAPIKey:        cfg.AnthropicAPIKey,
		GoogleAPIKey:           cfg.GoogleAPIKey,
		KiloAPIKey:             cfg.KiloAPIKey,
		KimiAPIKey:             cfg.KimiAPIKey,
		OpenRouterAPIKey:       cfg.OpenRouterAPIKey,
		ZAIAPIKey:              cfg.ZAIAPIKey,
		AlibabaAPIKey:          cfg.AlibabaAPIKey,
		OpenAIAPIKey:           cfg.OpenAIAPIKey,
		OpenAIAuthMode:         cfg.OpenAIAuthMode,
		GroqAPIKey:             cfg.GroqAPIKey,
		MaxIterations:          cfg.MaxIterations,
		MemoryWindowSize:       cfg.MemoryWindowSize,
	}, nil
}

// SaveEditable updates the user-editable config subset while preserving managed paths.
func SaveEditable(r *runtime.PathResolver, editable EditableConfig) error {
	cfg := normalizeFileConfig(fileConfig{
		LLMProvider:            editable.LLMProvider,
		LLMModel:               editable.LLMModel,
		STTProvider:            editable.STTProvider,
		TelegramBotToken:       editable.TelegramBotToken,
		TelegramAllowedUserIDs: append([]int64(nil), editable.TelegramAllowedUserIDs...),
		AnthropicAPIKey:        editable.AnthropicAPIKey,
		GoogleAPIKey:           editable.GoogleAPIKey,
		KiloAPIKey:             editable.KiloAPIKey,
		KimiAPIKey:             editable.KimiAPIKey,
		OpenRouterAPIKey:       editable.OpenRouterAPIKey,
		ZAIAPIKey:              editable.ZAIAPIKey,
		AlibabaAPIKey:          editable.AlibabaAPIKey,
		OpenAIAPIKey:           editable.OpenAIAPIKey,
		OpenAIAuthMode:         editable.OpenAIAuthMode,
		GroqAPIKey:             editable.GroqAPIKey,
		MaxIterations:          editable.MaxIterations,
		MemoryWindowSize:       editable.MemoryWindowSize,
	}, r)
	return writeConfigFile(r.AppConfig(), cfg)
}

func normalizeFileConfig(cfg fileConfig, r *runtime.PathResolver) fileConfig {
	defaults := defaultFileConfig(r)
	if cfg.TelegramAllowedUserIDs == nil {
		cfg.TelegramAllowedUserIDs = defaults.TelegramAllowedUserIDs
	}
	if cfg.LLMProvider == "" {
		cfg.LLMProvider = defaults.LLMProvider
	}
	if cfg.LLMModel == "" {
		cfg.LLMModel = defaultModelForProvider(cfg.LLMProvider)
	}
	if cfg.OpenAIAuthMode == "" {
		cfg.OpenAIAuthMode = defaults.OpenAIAuthMode
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
		LLMProvider:            cfg.LLMProvider,
		LLMModel:               cfg.LLMModel,
		STTProvider:            cfg.STTProvider,
		TelegramBotToken:       cfg.TelegramBotToken,
		TelegramAllowedUserIDs: cfg.TelegramAllowedUserIDs,
		AnthropicAPIKey:        cfg.AnthropicAPIKey,
		GoogleAPIKey:           cfg.GoogleAPIKey,
		KiloAPIKey:             cfg.KiloAPIKey,
		KimiAPIKey:             cfg.KimiAPIKey,
		OpenRouterAPIKey:       cfg.OpenRouterAPIKey,
		ZAIAPIKey:              cfg.ZAIAPIKey,
		AlibabaAPIKey:          cfg.AlibabaAPIKey,
		OpenAIAPIKey:           cfg.OpenAIAPIKey,
		OpenAIAuthMode:         cfg.OpenAIAuthMode,
		GroqAPIKey:             cfg.GroqAPIKey,
		MaxIterations:          cfg.MaxIterations,
		DBPath:                 cfg.DBPath,
		MemoryWindowSize:       cfg.MemoryWindowSize,
		MCPConfigPath:          cfg.MCPConfigPath,
	}
}

func defaultLLMModelForProvider(provider string) string {
	return defaultModelForProvider(provider)
}

func sameFileConfig(a, b fileConfig) bool {
	if a.TelegramBotToken != b.TelegramBotToken ||
		a.LLMProvider != b.LLMProvider ||
		a.LLMModel != b.LLMModel ||
		a.STTProvider != b.STTProvider ||
		a.AnthropicAPIKey != b.AnthropicAPIKey ||
		a.GoogleAPIKey != b.GoogleAPIKey ||
		a.KiloAPIKey != b.KiloAPIKey ||
		a.KimiAPIKey != b.KimiAPIKey ||
		a.OpenRouterAPIKey != b.OpenRouterAPIKey ||
		a.ZAIAPIKey != b.ZAIAPIKey ||
		a.AlibabaAPIKey != b.AlibabaAPIKey ||
		a.OpenAIAPIKey != b.OpenAIAPIKey ||
		a.OpenAIAuthMode != b.OpenAIAuthMode ||
		a.GroqAPIKey != b.GroqAPIKey ||
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
	return true
}

