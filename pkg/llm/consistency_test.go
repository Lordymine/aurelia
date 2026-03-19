package llm

import "testing"

func TestProviderDefaultsStayConsistentWithFallbacks(t *testing.T) {
	for _, spec := range Providers() {
		fallbacks := FallbackModels(spec.ID)
		if len(fallbacks) == 0 {
			t.Fatalf("provider %q has no fallback models", spec.ID)
		}

		found := false
		for _, model := range fallbacks {
			if NormalizeModelID(spec.ID, model.ID) == NormalizeModelID(spec.ID, spec.DefaultModel) {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("provider %q default model %q not present in fallbacks %+v", spec.ID, spec.DefaultModel, fallbacks)
		}
	}
}

func TestFallbackModelsHaveContextWindowCoverage(t *testing.T) {
	for _, spec := range Providers() {
		for _, model := range FallbackModels(spec.ID) {
			if got := ContextWindow(spec.ID, model.ID); got <= 0 {
				t.Fatalf("ContextWindow(%q, %q) = %d", spec.ID, model.ID, got)
			}
		}
	}
}
