package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/kocar/aurelia/internal/runtime"
)

func TestLoad_CreatesDefaultAppConfigWhenMissing(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("AURELIA_HOME", tmpDir)

	r, err := runtime.New()
	if err != nil {
		t.Fatalf("runtime.New() unexpected error: %v", err)
	}

	cfg, err := Load(r)
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	if cfg.MaxIterations != defaultMaxIterations {
		t.Fatalf("MaxIterations = %d, want %d", cfg.MaxIterations, defaultMaxIterations)
	}
	if cfg.DefaultProvider != defaultLLMProvider {
		t.Fatalf("DefaultProvider = %q, want %q", cfg.DefaultProvider, defaultLLMProvider)
	}
	if cfg.DefaultModel != defaultModelForProvider(defaultLLMProvider) {
		t.Fatalf("DefaultModel = %q, want %q", cfg.DefaultModel, defaultModelForProvider(defaultLLMProvider))
	}
	if cfg.STTProvider != defaultSTTProvider {
		t.Fatalf("STTProvider = %q, want %q", cfg.STTProvider, defaultSTTProvider)
	}
	if cfg.MemoryWindowSize != defaultMemoryWindowSize {
		t.Fatalf("MemoryWindowSize = %d, want %d", cfg.MemoryWindowSize, defaultMemoryWindowSize)
	}
	if cfg.Providers == nil {
		t.Fatal("Providers should not be nil")
	}

	if _, err := os.Stat(r.AppConfig()); err != nil {
		t.Fatalf("expected app config file to be created: %v", err)
	}
}

func TestLoadConfig_NewSchema(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("AURELIA_HOME", tmpDir)

	r, err := runtime.New()
	if err != nil {
		t.Fatalf("runtime.New() unexpected error: %v", err)
	}

	cfgJSON := `{
		"default_provider": "anthropic",
		"default_model": "claude-sonnet-4-6",
		"providers": {
			"anthropic": {"api_key": "sk-ant-test"},
			"kimi": {"api_key": "sk-kimi-test", "base_url": "https://api.kimi.ai"}
		},
		"telegram_bot_token": "test-token",
		"telegram_allowed_user_ids": [123456],
		"embedding_provider": "voyage",
		"embedding_model": "voyage-3",
		"stt_provider": "groq",
		"max_iterations": 100,
		"memory_window_size": 15
	}`

	if err := os.MkdirAll(filepath.Dir(r.AppConfig()), 0o700); err != nil {
		t.Fatalf("MkdirAll() unexpected error: %v", err)
	}
	if err := os.WriteFile(r.AppConfig(), []byte(cfgJSON), 0o600); err != nil {
		t.Fatalf("WriteFile() unexpected error: %v", err)
	}

	cfg, err := Load(r)
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	if cfg.DefaultProvider != "anthropic" {
		t.Fatalf("DefaultProvider = %q, want %q", cfg.DefaultProvider, "anthropic")
	}
	if cfg.DefaultModel != "claude-sonnet-4-6" {
		t.Fatalf("DefaultModel = %q, want %q", cfg.DefaultModel, "claude-sonnet-4-6")
	}
	if cfg.ProviderAPIKey("anthropic") != "sk-ant-test" {
		t.Fatalf("anthropic api_key = %q", cfg.ProviderAPIKey("anthropic"))
	}
	if cfg.ProviderAPIKey("kimi") != "sk-kimi-test" {
		t.Fatalf("kimi api_key = %q", cfg.ProviderAPIKey("kimi"))
	}
	if cfg.ProviderBaseURL("kimi") != "https://api.kimi.ai" {
		t.Fatalf("kimi base_url = %q", cfg.ProviderBaseURL("kimi"))
	}
	if cfg.TelegramBotToken != "test-token" {
		t.Fatalf("TelegramBotToken = %q", cfg.TelegramBotToken)
	}
	if !reflect.DeepEqual(cfg.TelegramAllowedUserIDs, []int64{123456}) {
		t.Fatalf("TelegramAllowedUserIDs = %v", cfg.TelegramAllowedUserIDs)
	}
	if cfg.EmbeddingProvider != "voyage" {
		t.Fatalf("EmbeddingProvider = %q", cfg.EmbeddingProvider)
	}
	if cfg.EmbeddingModel != "voyage-3" {
		t.Fatalf("EmbeddingModel = %q", cfg.EmbeddingModel)
	}
	if cfg.MaxIterations != 100 {
		t.Fatalf("MaxIterations = %d", cfg.MaxIterations)
	}
	if cfg.MemoryWindowSize != 15 {
		t.Fatalf("MemoryWindowSize = %d", cfg.MemoryWindowSize)
	}
}

func TestSaveAndReload(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("AURELIA_HOME", tmpDir)

	r, err := runtime.New()
	if err != nil {
		t.Fatalf("runtime.New() unexpected error: %v", err)
	}

	original := EditableConfig{
		LLMProvider:            "anthropic",
		LLMModel:               "claude-sonnet-4-6",
		STTProvider:            "groq",
		TelegramBotToken:       "my-token",
		TelegramAllowedUserIDs: []int64{42, 99},
		AnthropicAPIKey:        "ant-key",
		KimiAPIKey:             "kimi-key",
		GroqAPIKey:             "groq-key",
		MaxIterations:          300,
		MemoryWindowSize:       25,
		EmbeddingProvider:      "voyage",
		EmbeddingModel:         "voyage-3",
		EmbeddingAPIKey:        "voyage-key",
	}

	if err := SaveEditable(r, original); err != nil {
		t.Fatalf("SaveEditable() error = %v", err)
	}

	cfg, err := Load(r)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.DefaultProvider != "anthropic" {
		t.Fatalf("DefaultProvider = %q", cfg.DefaultProvider)
	}
	if cfg.DefaultModel != "claude-sonnet-4-6" {
		t.Fatalf("DefaultModel = %q", cfg.DefaultModel)
	}
	if cfg.ProviderAPIKey("anthropic") != "ant-key" {
		t.Fatalf("anthropic key = %q", cfg.ProviderAPIKey("anthropic"))
	}
	if cfg.ProviderAPIKey("kimi") != "kimi-key" {
		t.Fatalf("kimi key = %q", cfg.ProviderAPIKey("kimi"))
	}
	if cfg.ProviderAPIKey("groq") != "groq-key" {
		t.Fatalf("groq key = %q", cfg.ProviderAPIKey("groq"))
	}
	if cfg.TelegramBotToken != "my-token" {
		t.Fatalf("TelegramBotToken = %q", cfg.TelegramBotToken)
	}
	if !reflect.DeepEqual(cfg.TelegramAllowedUserIDs, []int64{42, 99}) {
		t.Fatalf("TelegramAllowedUserIDs = %v", cfg.TelegramAllowedUserIDs)
	}
	if cfg.EmbeddingProvider != "voyage" {
		t.Fatalf("EmbeddingProvider = %q", cfg.EmbeddingProvider)
	}
	if cfg.EmbeddingModel != "voyage-3" {
		t.Fatalf("EmbeddingModel = %q", cfg.EmbeddingModel)
	}
	if cfg.EmbeddingAPIKey != "voyage-key" {
		t.Fatalf("EmbeddingAPIKey = %q", cfg.EmbeddingAPIKey)
	}
	if cfg.MaxIterations != 300 {
		t.Fatalf("MaxIterations = %d", cfg.MaxIterations)
	}
	if cfg.MemoryWindowSize != 25 {
		t.Fatalf("MemoryWindowSize = %d", cfg.MemoryWindowSize)
	}
	if cfg.DBPath != filepath.Join(tmpDir, "data", "aurelia.db") {
		t.Fatalf("DBPath = %q", cfg.DBPath)
	}
	if cfg.MCPConfigPath != filepath.Join(tmpDir, "config", "mcp_servers.json") {
		t.Fatalf("MCPConfigPath = %q", cfg.MCPConfigPath)
	}
}

func TestConfig_DefaultValues(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("AURELIA_HOME", tmpDir)

	r, err := runtime.New()
	if err != nil {
		t.Fatalf("runtime.New() unexpected error: %v", err)
	}

	// Write minimal config
	if err := os.MkdirAll(filepath.Dir(r.AppConfig()), 0o700); err != nil {
		t.Fatalf("MkdirAll() unexpected error: %v", err)
	}
	if err := os.WriteFile(r.AppConfig(), []byte(`{"telegram_bot_token":"abc"}`), 0o600); err != nil {
		t.Fatalf("WriteFile() unexpected error: %v", err)
	}

	cfg, err := Load(r)
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	if cfg.TelegramBotToken != "abc" {
		t.Fatalf("TelegramBotToken = %q, want %q", cfg.TelegramBotToken, "abc")
	}
	if cfg.DefaultProvider != defaultLLMProvider {
		t.Fatalf("DefaultProvider = %q, want %q", cfg.DefaultProvider, defaultLLMProvider)
	}
	if cfg.DefaultModel != defaultModelForProvider(defaultLLMProvider) {
		t.Fatalf("DefaultModel = %q, want %q", cfg.DefaultModel, defaultModelForProvider(defaultLLMProvider))
	}
	if cfg.STTProvider != defaultSTTProvider {
		t.Fatalf("STTProvider = %q, want %q", cfg.STTProvider, defaultSTTProvider)
	}
	if cfg.MaxIterations != defaultMaxIterations {
		t.Fatalf("MaxIterations = %d, want %d", cfg.MaxIterations, defaultMaxIterations)
	}
	if cfg.DBPath != filepath.Join(tmpDir, "data", "aurelia.db") {
		t.Fatalf("DBPath = %q, want instance default", cfg.DBPath)
	}
	if cfg.MemoryWindowSize != defaultMemoryWindowSize {
		t.Fatalf("MemoryWindowSize = %d, want %d", cfg.MemoryWindowSize, defaultMemoryWindowSize)
	}
	if cfg.MCPConfigPath != filepath.Join(tmpDir, "config", "mcp_servers.json") {
		t.Fatalf("MCPConfigPath = %q, want instance default", cfg.MCPConfigPath)
	}
}

func TestLoad_LegacyFormatMigration(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("AURELIA_HOME", tmpDir)

	r, err := runtime.New()
	if err != nil {
		t.Fatalf("runtime.New() unexpected error: %v", err)
	}

	legacyJSON := `{
		"llm_provider": "kimi",
		"llm_model": "moonshot-v1-8k",
		"stt_provider": "groq",
		"telegram_bot_token": "telegram-token",
		"telegram_allowed_user_ids": [1, 2, 3],
		"anthropic_api_key": "anthropic-key",
		"google_api_key": "google-key",
		"kimi_api_key": "kimi-key",
		"openai_api_key": "openai-key",
		"openai_auth_mode": "codex",
		"groq_api_key": "groq-key",
		"max_iterations": 321,
		"memory_window_size": 42
	}`

	if err := os.MkdirAll(filepath.Dir(r.AppConfig()), 0o700); err != nil {
		t.Fatalf("MkdirAll() unexpected error: %v", err)
	}
	if err := os.WriteFile(r.AppConfig(), []byte(legacyJSON), 0o600); err != nil {
		t.Fatalf("WriteFile() unexpected error: %v", err)
	}

	cfg, err := Load(r)
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	if cfg.DefaultProvider != "kimi" {
		t.Fatalf("DefaultProvider = %q, want %q", cfg.DefaultProvider, "kimi")
	}
	if cfg.DefaultModel != "moonshot-v1-8k" {
		t.Fatalf("DefaultModel = %q, want %q", cfg.DefaultModel, "moonshot-v1-8k")
	}
	if cfg.TelegramBotToken != "telegram-token" {
		t.Fatalf("TelegramBotToken = %q", cfg.TelegramBotToken)
	}
	if !reflect.DeepEqual(cfg.TelegramAllowedUserIDs, []int64{1, 2, 3}) {
		t.Fatalf("TelegramAllowedUserIDs = %v", cfg.TelegramAllowedUserIDs)
	}
	if cfg.ProviderAPIKey("anthropic") != "anthropic-key" {
		t.Fatalf("anthropic key = %q", cfg.ProviderAPIKey("anthropic"))
	}
	if cfg.ProviderAPIKey("google") != "google-key" {
		t.Fatalf("google key = %q", cfg.ProviderAPIKey("google"))
	}
	if cfg.ProviderAPIKey("kimi") != "kimi-key" {
		t.Fatalf("kimi key = %q", cfg.ProviderAPIKey("kimi"))
	}
	if cfg.ProviderAPIKey("openai") != "openai-key" {
		t.Fatalf("openai key = %q", cfg.ProviderAPIKey("openai"))
	}
	if cfg.ProviderAuthMode("openai") != "codex" {
		t.Fatalf("openai auth_mode = %q", cfg.ProviderAuthMode("openai"))
	}
	if cfg.ProviderAPIKey("groq") != "groq-key" {
		t.Fatalf("groq key = %q", cfg.ProviderAPIKey("groq"))
	}
	if cfg.MaxIterations != 321 {
		t.Fatalf("MaxIterations = %d", cfg.MaxIterations)
	}
	if cfg.MemoryWindowSize != 42 {
		t.Fatalf("MemoryWindowSize = %d", cfg.MemoryWindowSize)
	}

	// Verify the file was rewritten in new format
	data, err := os.ReadFile(r.AppConfig())
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	var rewritten fileConfig
	if err := json.Unmarshal(data, &rewritten); err != nil {
		t.Fatalf("Unmarshal rewritten config error = %v", err)
	}
	if rewritten.DefaultProvider != "kimi" {
		t.Fatalf("rewritten DefaultProvider = %q", rewritten.DefaultProvider)
	}
	if len(rewritten.Providers) == 0 {
		t.Fatal("rewritten config should have providers map")
	}
}

func TestSaveEditable_PreservesManagedPaths(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("AURELIA_HOME", tmpDir)

	r, err := runtime.New()
	if err != nil {
		t.Fatalf("runtime.New() unexpected error: %v", err)
	}

	if err := SaveEditable(r, EditableConfig{
		LLMProvider:            "kimi",
		LLMModel:               "kimi-k2-thinking",
		STTProvider:            "groq",
		TelegramBotToken:       "telegram-token",
		TelegramAllowedUserIDs: []int64{7, 8},
		AnthropicAPIKey:        "anthropic-key",
		GoogleAPIKey:           "google-key",
		KiloAPIKey:             "kilo-key",
		KimiAPIKey:             "kimi-key",
		OpenRouterAPIKey:       "openrouter-key",
		ZAIAPIKey:              "zai-key",
		AlibabaAPIKey:          "alibaba-key",
		GroqAPIKey:             "groq-key",
		MaxIterations:          900,
		MemoryWindowSize:       25,
	}); err != nil {
		t.Fatalf("SaveEditable() unexpected error: %v", err)
	}

	cfg, err := Load(r)
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	if cfg.DBPath != filepath.Join(tmpDir, "data", "aurelia.db") {
		t.Fatalf("DBPath = %q, want managed default", cfg.DBPath)
	}
	if cfg.DefaultProvider != "kimi" || cfg.DefaultModel != "kimi-k2-thinking" || cfg.STTProvider != "groq" {
		t.Fatalf("unexpected providers llm=%q model=%q stt=%q", cfg.DefaultProvider, cfg.DefaultModel, cfg.STTProvider)
	}
	if cfg.ProviderAPIKey("kilo") != "kilo-key" {
		t.Fatalf("kilo key = %q, want %q", cfg.ProviderAPIKey("kilo"), "kilo-key")
	}
	if cfg.MCPConfigPath != filepath.Join(tmpDir, "config", "mcp_servers.json") {
		t.Fatalf("MCPConfigPath = %q, want managed default", cfg.MCPConfigPath)
	}
	if !reflect.DeepEqual(cfg.TelegramAllowedUserIDs, []int64{7, 8}) {
		t.Fatalf("TelegramAllowedUserIDs = %v", cfg.TelegramAllowedUserIDs)
	}
}

func TestLoadEditable_RoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("AURELIA_HOME", tmpDir)

	r, err := runtime.New()
	if err != nil {
		t.Fatalf("runtime.New() unexpected error: %v", err)
	}

	original := EditableConfig{
		LLMProvider:            "anthropic",
		LLMModel:               "claude-sonnet-4-6",
		STTProvider:            "groq",
		TelegramBotToken:       "token",
		TelegramAllowedUserIDs: []int64{1},
		AnthropicAPIKey:        "ant-key",
		GroqAPIKey:             "groq-key",
		MaxIterations:          100,
		MemoryWindowSize:       10,
	}

	if err := SaveEditable(r, original); err != nil {
		t.Fatalf("SaveEditable() error = %v", err)
	}

	loaded, err := LoadEditable(r)
	if err != nil {
		t.Fatalf("LoadEditable() error = %v", err)
	}

	if loaded.LLMProvider != "anthropic" {
		t.Fatalf("LLMProvider = %q", loaded.LLMProvider)
	}
	if loaded.AnthropicAPIKey != "ant-key" {
		t.Fatalf("AnthropicAPIKey = %q", loaded.AnthropicAPIKey)
	}
	if loaded.GroqAPIKey != "groq-key" {
		t.Fatalf("GroqAPIKey = %q", loaded.GroqAPIKey)
	}
}

func TestProviderAPIKey_ReturnsEmptyForMissing(t *testing.T) {
	cfg := &AppConfig{
		Providers: map[string]ProviderConfig{
			"anthropic": {APIKey: "test-key"},
		},
	}

	if cfg.ProviderAPIKey("anthropic") != "test-key" {
		t.Fatalf("expected test-key, got %q", cfg.ProviderAPIKey("anthropic"))
	}
	if cfg.ProviderAPIKey("nonexistent") != "" {
		t.Fatalf("expected empty for nonexistent, got %q", cfg.ProviderAPIKey("nonexistent"))
	}
}
