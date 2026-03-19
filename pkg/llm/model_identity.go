package llm

import "strings"

func NormalizeModelID(provider, model string) string {
	normalizedProvider := NormalizeProvider(provider)
	normalizedModel := strings.TrimSpace(strings.ToLower(model))
	normalizedModel = strings.ReplaceAll(normalizedModel, "z-ai/", "zai/")

	switch normalizedProvider {
	case "kimi", "anthropic", "google", "zai", "alibaba", "openai":
		prefix := normalizedProvider + "/"
		if strings.HasPrefix(normalizedModel, prefix) {
			return strings.TrimPrefix(normalizedModel, prefix)
		}
		return normalizedModel
	default:
		return normalizedModel
	}
}

func ProviderModelKey(provider, model string) string {
	return NormalizeProvider(provider) + "/" + NormalizeModelID(provider, model)
}
