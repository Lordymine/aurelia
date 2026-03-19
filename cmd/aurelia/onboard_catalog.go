package main

import (
	"context"
	"strings"

	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/pkg/llm"
)

func loadModelOptions(cfg config.EditableConfig) []llm.ModelOption {
	options, _ := llmModelCatalog(context.Background(), cfg.LLMProvider, modelCatalogCredentials(cfg))
	if len(options) != 0 {
		return options
	}
	return llm.FallbackModels(cfg.LLMProvider)
}

func resolveModelOptions(cfg config.EditableConfig) ([]llm.ModelOption, string) {
	options, err := llmModelCatalog(context.Background(), cfg.LLMProvider, modelCatalogCredentials(cfg))
	if err == nil && len(options) != 0 {
		return options, "provider catalog"
	}

	fallback := llm.FallbackModels(cfg.LLMProvider)
	if len(fallback) != 0 {
		return fallback, "curated fallback"
	}
	return nil, "no catalog available"
}

func filterModelOptions(cfg config.EditableConfig, options []llm.ModelOption, filter string, capability modelCapabilityFilter) []llm.ModelOption {
	filter = strings.ToLower(strings.TrimSpace(filter))

	filtered := make([]llm.ModelOption, 0, len(options))
	for _, option := range options {
		if !matchesCapabilityFilter(option, capability) {
			continue
		}
		if filter != "" && usesProviderModelSearch(cfg) && !matchesModelFilter(option, filter) {
			continue
		}
		filtered = append(filtered, option)
	}
	return filtered
}

func matchesModelFilter(option llm.ModelOption, filter string) bool {
	candidates := []string{
		strings.ToLower(option.ID),
		strings.ToLower(option.Name),
		strings.ToLower(option.Label()),
		strings.ToLower(openRouterProviderName(option.ID)),
	}
	for _, candidate := range candidates {
		if strings.Contains(candidate, filter) {
			return true
		}
	}
	return false
}

func openRouterProviderName(modelID string) string {
	prefix, _, ok := strings.Cut(modelID, "/")
	if !ok {
		return ""
	}
	return prefix
}

func selectedModelIndex(options []llm.ModelOption, current string) int {
	for i, option := range options {
		if option.ID == current {
			return i
		}
	}
	return 0
}

func modelCatalogCredentials(cfg config.EditableConfig) llm.ModelCatalogCredentials {
	return llm.ModelCatalogCredentials{
		AnthropicAPIKey:  cfg.AnthropicAPIKey,
		GoogleAPIKey:     cfg.GoogleAPIKey,
		KiloAPIKey:       cfg.KiloAPIKey,
		KimiAPIKey:       cfg.KimiAPIKey,
		OpenRouterAPIKey: cfg.OpenRouterAPIKey,
		ZAIAPIKey:        cfg.ZAIAPIKey,
		AlibabaAPIKey:    cfg.AlibabaAPIKey,
		OpenAIAPIKey:     cfg.OpenAIAPIKey,
		OpenAIAuthMode:   cfg.OpenAIAuthMode,
	}
}

func matchesCapabilityFilter(option llm.ModelOption, capability modelCapabilityFilter) bool {
	switch capability {
	case modelCapabilityVision:
		return option.SupportsImageInput
	case modelCapabilityTools:
		return option.SupportsTools
	case modelCapabilityFree:
		return option.IsFree
	default:
		return true
	}
}
