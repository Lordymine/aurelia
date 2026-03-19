package llm

import "strings"

const unknownContextWindow = 0

var modelContextWindows = map[string]int{
	"anthropic/claude-sonnet-4-6":  200000,
	"anthropic/claude-opus-4-6":    200000,
	"anthropic/claude-haiku-4-5":   200000,
	"google/gemini-2.5-pro":        1048576,
	"google/gemini-2.5-flash":      1048576,
	"google/gemini-2.5-flash-lite": 1048576,
	"kimi/kimi-k2-thinking":        128000,
	"kimi/kimi-k2-thinking-turbo":  128000,
	"kimi/k2.5":                    128000,
	"kimi/moonshot-v1-vision":      128000,
	"kimi/moonshot-v1-8k":          8192,
	"kimi/moonshot-v1-32k":         32768,
	"kimi/moonshot-v1-128k":        131072,
	"openai/gpt-5.4":               400000,
	"openai/gpt-5-mini":            400000,
	"openai/gpt-5.2":               400000,
	"openai/gpt-5.2-codex":         400000,
	"openai/gpt-5.1":               400000,
	"openai/o4-mini":               200000,
	"openai/gpt-4.1":               1047576,
	"openrouter/openrouter/auto":   200000,
	"openrouter/openrouter/free":   128000,
	"zai/glm-5":                    128000,
	"zai/glm-4.7":                  128000,
	"zai/glm-4.6v":                 65536,
	"zai/glm-4.5-air":              65536,
	"alibaba/qwen3-coder-plus":     1000000,
	"alibaba/qwen3-coder-next":     262144,
	"alibaba/qwen-vl-max":          131072,
	"alibaba/qwen3.5-plus":         131072,
}

// ContextWindow returns the best-effort context window for a provider/model pair.
func ContextWindow(provider, model string) int {
	key := normalizeContextKey(provider, model)
	if limit, ok := modelContextWindows[key]; ok {
		return limit
	}

	switch strings.TrimSpace(strings.ToLower(provider)) {
	case "openrouter", "kilo":
		return inferGatewayModelContextWindow(model)
	default:
		return unknownContextWindow
	}
}

func normalizeContextKey(provider, model string) string {
	return strings.TrimSpace(strings.ToLower(provider)) + "/" + strings.TrimSpace(strings.ToLower(model))
}

func inferGatewayModelContextWindow(model string) int {
	lower := strings.TrimSpace(strings.ToLower(model))
	if strings.HasPrefix(lower, "openai/") || strings.Contains(lower, "/gpt-5") {
		return 400000
	}
	if strings.Contains(lower, "gpt-4.1") {
		return 1047576
	}
	if strings.Contains(lower, "claude") {
		return 200000
	}
	if strings.Contains(lower, "gemini-2.5") {
		return 1048576
	}
	if strings.Contains(lower, "qwen3-coder-plus") {
		return 1000000
	}
	if strings.Contains(lower, "glm-5") {
		return 128000
	}
	return unknownContextWindow
}
