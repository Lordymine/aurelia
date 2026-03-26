package telegram

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBootstrapState_HasStepField(t *testing.T) {
	state := bootstrapState{Choice: "coder", Step: "assistant"}
	if state.Step != "assistant" {
		t.Fatalf("expected step 'assistant', got %q", state.Step)
	}

	state2 := bootstrapState{Choice: "assist", Step: "profile"}
	if state2.Step != "profile" {
		t.Fatalf("expected step 'profile', got %q", state2.Step)
	}
}

func TestParseGeneratedPersona_ValidOutput(t *testing.T) {
	output := "===IDENTITY===\n# Identity content here\n===SOUL===\n# Soul content here\n"

	identity, soul, err := parseGeneratedPersona(output)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(identity, "# Identity content here") {
		t.Fatalf("expected identity content, got %q", identity)
	}
	if !strings.Contains(soul, "# Soul content here") {
		t.Fatalf("expected soul content, got %q", soul)
	}
}

func TestParseGeneratedPersona_MissingSoulMarker(t *testing.T) {
	output := "===IDENTITY===\n# Identity content here\n"

	_, _, err := parseGeneratedPersona(output)
	if err == nil {
		t.Fatal("expected error for missing SOUL marker")
	}
}

func TestParseGeneratedPersona_MissingIdentityMarker(t *testing.T) {
	output := "some random text without markers"

	_, _, err := parseGeneratedPersona(output)
	if err == nil {
		t.Fatal("expected error for missing IDENTITY marker")
	}
}

func TestParseGeneratedPersona_TrimsWhitespace(t *testing.T) {
	output := "===IDENTITY===\n\n  # Identity  \n\n===SOUL===\n\n  # Soul  \n\n"

	identity, soul, err := parseGeneratedPersona(output)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if identity != "# Identity" {
		t.Fatalf("expected trimmed identity, got %q", identity)
	}
	if soul != "# Soul" {
		t.Fatalf("expected trimmed soul, got %q", soul)
	}
}

func TestParseGeneratedPersona_ToleratesSpacesInMarkers(t *testing.T) {
	output := "=== IDENTITY ===\n# Identity\n=== SOUL ===\n# Soul\n"

	identity, soul, err := parseGeneratedPersona(output)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if identity != "# Identity" {
		t.Fatalf("expected identity, got %q", identity)
	}
	if soul != "# Soul" {
		t.Fatalf("expected soul, got %q", soul)
	}
}

func TestParseGeneratedPersona_ToleratesCodeFences(t *testing.T) {
	output := "```markdown\n===IDENTITY===\n# Identity\n===SOUL===\n# Soul\n```"

	identity, soul, err := parseGeneratedPersona(output)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if identity != "# Identity" {
		t.Fatalf("expected identity, got %q", identity)
	}
	if soul != "# Soul" {
		t.Fatalf("expected soul, got %q", soul)
	}
}

func TestParseGeneratedPersona_ToleratesHashHeaders(t *testing.T) {
	output := "## IDENTITY\n# My Identity\n## SOUL\n# My Soul\n"

	identity, soul, err := parseGeneratedPersona(output)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if identity != "# My Identity" {
		t.Fatalf("expected identity, got %q", identity)
	}
	if soul != "# My Soul" {
		t.Fatalf("expected soul, got %q", soul)
	}
}

func TestStripCodeFences(t *testing.T) {
	cases := map[string]string{
		"```\ncontent\n```":         "content",
		"```markdown\ncontent\n```": "content",
		"no fences":                 "no fences",
	}
	for input, want := range cases {
		got := stripCodeFences(input)
		if got != want {
			t.Fatalf("stripCodeFences(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestBuildAssistantGeneratePrompt_IncludesPresetAndDescription(t *testing.T) {
	preset, _ := bootstrapPresetForChoice("coder")
	prompt := buildAssistantGeneratePrompt(preset, "Quero um assistente sarcastico que fale como pirata")

	if !strings.Contains(prompt, "Aurelia Coder") || !strings.Contains(prompt, "Programacao") {
		t.Fatalf("expected preset context in prompt, got %q", prompt)
	}
	if !strings.Contains(prompt, "sarcastico") || !strings.Contains(prompt, "pirata") {
		t.Fatalf("expected user description in prompt, got %q", prompt)
	}
	if !strings.Contains(prompt, "===IDENTITY===") || !strings.Contains(prompt, "===SOUL===") {
		t.Fatalf("expected output format markers in prompt, got %q", prompt)
	}
}

func TestBuildUserGeneratePrompt_IncludesDescriptionAndFallback(t *testing.T) {
	prompt := buildUserGeneratePrompt("Me chamo Rafael, sou dev e quero respostas diretas", "RafaKocar")

	if !strings.Contains(prompt, "Rafael") {
		t.Fatalf("expected user description in prompt, got %q", prompt)
	}
	if !strings.Contains(prompt, "RafaKocar") {
		t.Fatalf("expected fallback name in prompt, got %q", prompt)
	}
}

func TestWriteGeneratedPersona_WritesFiles(t *testing.T) {
	dir := t.TempDir()
	identity := "---\nname: \"Test\"\nrole: \"Tester\"\n---\n\n# Identity"
	soul := "# Soul\nBe nice."

	err := writeGeneratedPersona(dir, identity, soul)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	gotIdentity, err := os.ReadFile(filepath.Join(dir, "IDENTITY.md"))
	if err != nil {
		t.Fatalf("failed to read IDENTITY.md: %v", err)
	}
	if string(gotIdentity) != identity {
		t.Fatalf("IDENTITY.md content mismatch: got %q, want %q", string(gotIdentity), identity)
	}

	gotSoul, err := os.ReadFile(filepath.Join(dir, "SOUL.md"))
	if err != nil {
		t.Fatalf("failed to read SOUL.md: %v", err)
	}
	if string(gotSoul) != soul {
		t.Fatalf("SOUL.md content mismatch: got %q, want %q", string(gotSoul), soul)
	}
}

func TestBootstrapStepAssistant_MessageConstant(t *testing.T) {
	if bootstrapAssistantMessage == "" {
		t.Fatal("bootstrapAssistantMessage must not be empty")
	}
}

func TestHandleBootstrapChoice_SetsStepAssistant(t *testing.T) {
	bc := &BotController{
		pendingBootstrap: make(map[int64]bootstrapState),
	}

	bc.setPendingBootstrap(42, bootstrapState{Choice: "coder", Step: bootstrapStepAssistant})

	state, ok := bc.popPendingBootstrap(42)
	if !ok {
		t.Fatal("expected pending bootstrap state")
	}
	if state.Step != bootstrapStepAssistant {
		t.Fatalf("expected step %q, got %q", bootstrapStepAssistant, state.Step)
	}
	if state.Choice != "coder" {
		t.Fatalf("expected choice 'coder', got %q", state.Choice)
	}
}
