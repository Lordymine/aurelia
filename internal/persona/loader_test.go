package persona

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadPersona(t *testing.T) {
	tempDir := t.TempDir()

	identityPath := filepath.Join(tempDir, "IDENTITY.md")
	soulPath := filepath.Join(tempDir, "SOUL.md")
	userPath := filepath.Join(tempDir, "USER.md")

	identityContent := `---
name: "TestAgent"
role: "Tester"
memory_window_size: 10
tools:
  - read_file
---
IDENTITY_BODY`

	soulContent := "SOUL_BODY"
	userContent := "USER_BODY"

	err := os.WriteFile(identityPath, []byte(identityContent), 0644)
	if err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(soulPath, []byte(soulContent), 0644)
	_ = os.WriteFile(userPath, []byte(userContent), 0644)

	persona, err := LoadPersona(identityPath, soulPath, userPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if persona.Config.Name != "TestAgent" {
		t.Errorf("expected Name 'TestAgent', got %q", persona.Config.Name)
	}

	if len(persona.Config.Tools) != 1 || persona.Config.Tools[0] != "read_file" {
		t.Errorf("expected tools ['read_file'], got %v", persona.Config.Tools)
	}

	if !strings.Contains(persona.SystemPrompt, "IDENTITY_BODY") {
		t.Errorf("prompt missing identity body: %s", persona.SystemPrompt)
	}
	if !strings.Contains(persona.SystemPrompt, "SOUL_BODY") {
		t.Errorf("prompt missing soul body: %s", persona.SystemPrompt)
	}
	if !strings.Contains(persona.SystemPrompt, "USER_BODY") {
		t.Errorf("prompt missing user body: %s", persona.SystemPrompt)
	}
}

func TestLoadPersona_MissingFiles(t *testing.T) {
	tempDir := t.TempDir()

	identityPath := filepath.Join(tempDir, "IDENTITY.md")
	_ = os.WriteFile(identityPath, []byte("test"), 0644)

	_, err := LoadPersona(identityPath, "bad_soul.md", "bad_user.md")
	if err == nil {
		t.Error("expected error for missing SOUL/USER files")
	}
}

func TestLoadPersona_IncludesCanonicalIdentityBlock(t *testing.T) {
	tempDir := t.TempDir()

	identityPath := filepath.Join(tempDir, "IDENTITY.md")
	soulPath := filepath.Join(tempDir, "SOUL.md")
	userPath := filepath.Join(tempDir, "USER.md")

	identityContent := `---
name: "Lex"
role: "Team Lead"
memory_window_size: 10
tools:
  - read_file
---
IDENTITY_BODY`

	if err := os.WriteFile(identityPath, []byte(identityContent), 0644); err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(soulPath, []byte("SOUL_BODY"), 0644)
	_ = os.WriteFile(userPath, []byte("# User\nNome: Rafael\nFuso horario: America/Sao_Paulo"), 0644)

	persona, err := LoadPersona(identityPath, soulPath, userPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(persona.SystemPrompt, "# CANONICAL IDENTITY") {
		t.Fatalf("expected canonical identity block, got %q", persona.SystemPrompt)
	}
	if !strings.Contains(persona.SystemPrompt, "Nome canonico do agente: Lex") {
		t.Fatalf("expected canonical agent name, got %q", persona.SystemPrompt)
	}
	if !strings.Contains(persona.SystemPrompt, "Nome canonico do usuario: Rafael") {
		t.Fatalf("expected canonical user name, got %q", persona.SystemPrompt)
	}
}

func TestLoadPersona_DoesNotPromotePlaceholderUserName(t *testing.T) {
	tempDir := t.TempDir()

	identityPath := filepath.Join(tempDir, "IDENTITY.md")
	soulPath := filepath.Join(tempDir, "SOUL.md")
	userPath := filepath.Join(tempDir, "USER.md")

	identityContent := `---
name: "Lex"
role: "Team Lead"
---
IDENTITY_BODY`

	if err := os.WriteFile(identityPath, []byte(identityContent), 0644); err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(soulPath, []byte("SOUL_BODY"), 0644)
	_ = os.WriteFile(userPath, []byte("# User\nNome: Usuario 12345"), 0644)

	persona, err := LoadPersona(identityPath, soulPath, userPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(persona.SystemPrompt, "Nome canonico do usuario: nao definido") {
		t.Fatalf("expected unresolved canonical user name, got %q", persona.SystemPrompt)
	}
}

// Tests for long-term memory, facts, notes, and retrieval were removed
// because they depend on internal/memory which was deleted.
// They will be rewritten when semantic memory is wired via the bridge.
