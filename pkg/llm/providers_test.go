package llm

import "testing"

func TestProviderChoicesAndLabelsStayAligned(t *testing.T) {
	choices := ProviderChoices()
	labels := ProviderLabels()
	if len(choices) == 0 {
		t.Fatal("expected provider choices")
	}
	if len(choices) != len(labels) {
		t.Fatalf("provider choices/labels length mismatch: %d vs %d", len(choices), len(labels))
	}
}

func TestDefaultModelForProvider(t *testing.T) {
	if got := DefaultModelForProvider("openai"); got != "gpt-5.4" {
		t.Fatalf("DefaultModelForProvider(openai) = %q", got)
	}
	if got := DefaultModelForProvider(""); got != "kimi-k2-thinking" {
		t.Fatalf("DefaultModelForProvider(\"\") = %q", got)
	}
}

func TestProviderLookup(t *testing.T) {
	spec, ok := Provider("openrouter")
	if !ok {
		t.Fatal("expected openrouter provider spec")
	}
	if !spec.SupportsModelSearch {
		t.Fatal("expected openrouter to support model search")
	}
	if spec.APIKeyLabel == "" || spec.Label == "" {
		t.Fatalf("unexpected provider spec %+v", spec)
	}
}
