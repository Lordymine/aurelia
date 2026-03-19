package llm

import (
	"context"
	"fmt"

	"github.com/kocar/aurelia/internal/agent"
)

type RuntimeConfig struct {
	Provider         string
	Model            string
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

type ProviderRuntime interface {
	agent.LLMProvider
	Close()
}

func BuildProvider(cfg RuntimeConfig) (ProviderRuntime, error) {
	switch NormalizeProvider(cfg.Provider) {
	case "anthropic":
		return NewAnthropicProvider(cfg.AnthropicAPIKey, cfg.Model), nil
	case "google":
		return NewGeminiProvider(context.Background(), cfg.GoogleAPIKey, cfg.Model)
	case "kilo":
		return NewKiloProvider(cfg.KiloAPIKey, cfg.Model), nil
	case "openrouter":
		return NewOpenRouterProvider(cfg.OpenRouterAPIKey, cfg.Model), nil
	case "zai":
		return NewZAIProvider(cfg.ZAIAPIKey, cfg.Model), nil
	case "alibaba":
		return NewAlibabaProvider(cfg.AlibabaAPIKey, cfg.Model), nil
	case "openai":
		if cfg.OpenAIAuthMode == "codex" {
			if err := EnsureCodexCLIAvailable(); err != nil {
				return nil, err
			}
			return NewCodexCLIProvider(cfg.Model)
		}
		return NewOpenAIProvider(cfg.OpenAIAPIKey, cfg.Model), nil
	case "kimi":
		return NewKimiProvider(cfg.KimiAPIKey, cfg.Model), nil
	default:
		return nil, fmt.Errorf("unsupported llm provider %q", cfg.Provider)
	}
}
