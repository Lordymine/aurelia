package llm

import "testing"

func TestContextWindow(t *testing.T) {
	t.Parallel()

	if got := ContextWindow("openai", "gpt-5.4"); got != 400000 {
		t.Fatalf("ContextWindow(openai, gpt-5.4) = %d", got)
	}
	if got := ContextWindow("google", "gemini-2.5-pro"); got != 1048576 {
		t.Fatalf("ContextWindow(google, gemini-2.5-pro) = %d", got)
	}
	if got := ContextWindow("openrouter", "anthropic/claude-sonnet-4.6"); got != 200000 {
		t.Fatalf("ContextWindow(openrouter, anthropic/claude-sonnet-4.6) = %d", got)
	}
	if got := ContextWindow("kilo", "z-ai/glm-5-turbo"); got != 128000 {
		t.Fatalf("ContextWindow(kilo, z-ai/glm-5-turbo) = %d", got)
	}
	if got := ContextWindow("unknown", "mystery"); got != 0 {
		t.Fatalf("ContextWindow(unknown, mystery) = %d", got)
	}
}
