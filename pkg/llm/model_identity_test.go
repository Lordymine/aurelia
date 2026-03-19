package llm

import "testing"

func TestNormalizeModelID(t *testing.T) {
	if got := NormalizeModelID("openai", "openai/gpt-5.4"); got != "gpt-5.4" {
		t.Fatalf("NormalizeModelID(openai, openai/gpt-5.4) = %q", got)
	}
	if got := NormalizeModelID("zai", "z-ai/glm-5"); got != "glm-5" {
		t.Fatalf("NormalizeModelID(zai, z-ai/glm-5) = %q", got)
	}
	if got := NormalizeModelID("kilo", "z-ai/glm-5-turbo"); got != "zai/glm-5-turbo" {
		t.Fatalf("NormalizeModelID(kilo, z-ai/glm-5-turbo) = %q", got)
	}
}

func TestProviderModelKey(t *testing.T) {
	if got := ProviderModelKey("openai", "openai/gpt-5.4"); got != "openai/gpt-5.4" {
		t.Fatalf("ProviderModelKey(openai, openai/gpt-5.4) = %q", got)
	}
	if got := ProviderModelKey("zai", "z-ai/glm-5"); got != "zai/glm-5" {
		t.Fatalf("ProviderModelKey(zai, z-ai/glm-5) = %q", got)
	}
}
