package persona

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kocar/aurelia/internal/memory"
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

	// Validate System Prompt assembly
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
	
	// Create only one file to force error
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

	soulContent := "SOUL_BODY"
	userContent := `# User
Nome: Rafael
Fuso horario: America/Sao_Paulo`

	if err := os.WriteFile(identityPath, []byte(identityContent), 0644); err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(soulPath, []byte(soulContent), 0644)
	_ = os.WriteFile(userPath, []byte(userContent), 0644)

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
	if !strings.Contains(persona.SystemPrompt, "Papel canonico do agente: Team Lead") {
		t.Fatalf("expected canonical agent role, got %q", persona.SystemPrompt)
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
memory_window_size: 10
tools:
  - read_file
---
IDENTITY_BODY`

	soulContent := "SOUL_BODY"
	userContent := `# User
Nome: Usuario 12345
Fuso horario: America/Sao_Paulo`

	if err := os.WriteFile(identityPath, []byte(identityContent), 0644); err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(soulPath, []byte(soulContent), 0644)
	_ = os.WriteFile(userPath, []byte(userContent), 0644)

	persona, err := LoadPersona(identityPath, soulPath, userPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(persona.SystemPrompt, "Nome canonico do usuario: Usuario 12345") {
		t.Fatalf("placeholder user name must not become canonical: %q", persona.SystemPrompt)
	}
	if !strings.Contains(persona.SystemPrompt, "Nome canonico do usuario: nao definido") {
		t.Fatalf("expected unresolved canonical user name, got %q", persona.SystemPrompt)
	}
}

func TestPersonaRenderSystemPrompt_UsesResolvedIdentity(t *testing.T) {
	persona := &Persona{
		Config: Config{Name: "Lex", Role: "Team Lead"},
		SystemPrompt: "# CANONICAL IDENTITY\nNome canonico do agente: Lex\nPapel canonico do agente: Team Lead\nNome canonico do usuario: nao definido\n\nIDENTITY_BODY\n\nSOUL_BODY\n\n# User\nNome: Nao definido",
		PromptBody: "IDENTITY_BODY\n\nSOUL_BODY\n\n# User\nNome: Nao definido",
		CanonicalIdentity: CanonicalIdentity{
			AgentName: "Lex",
			AgentRole: "Team Lead",
			UserName:  "nao definido",
		},
	}

	got := persona.RenderSystemPrompt(CanonicalIdentity{
		AgentName: "Lex",
		AgentRole: "Chief Architect",
		UserName:  "Rafael",
	}, nil, nil)

	if !strings.Contains(got, "Papel canonico do agente: Chief Architect") {
		t.Fatalf("expected overridden agent role, got %q", got)
	}
	if !strings.Contains(got, "Nome canonico do usuario: Rafael") {
		t.Fatalf("expected overridden user name, got %q", got)
	}
	if strings.Contains(got, "Nome canonico do usuario: nao definido") {
		t.Fatalf("expected canonical user name to be replaced, got %q", got)
	}
}

func TestPersonaRenderSystemPrompt_IncludesLongTermNotes(t *testing.T) {
	persona := &Persona{
		Config:     Config{Name: "Lex", Role: "Team Lead"},
		PromptBody: "IDENTITY_BODY\n\nSOUL_BODY\n\n# User\nNome: Rafael",
	}

	got := persona.RenderSystemPrompt(CanonicalIdentity{
		AgentName: "Lex",
		AgentRole: "Team Lead",
		UserName:  "Rafael",
	}, []memory.Fact{
		{
			Key:   "project.memory.strategy",
			Value: "sqlite + facts + notes",
		},
	}, []memory.Note{
		{
			Topic:      "architecture",
			Kind:       "decision",
			Summary:    "Decidido manter SQLite com facts e notes.",
			Importance: 8,
		},
		{
			Topic:      "memory",
			Kind:       "decision",
			Summary:    "Embeddings ficam opcionais, nao no core.",
			Importance: 7,
		},
	})

	if !strings.Contains(got, "# LONG-TERM MEMORY") {
		t.Fatalf("expected long-term memory block, got %q", got)
	}
	if !strings.Contains(got, "project.memory.strategy: sqlite + facts + notes") {
		t.Fatalf("expected fact in memory block, got %q", got)
	}
	if !strings.Contains(got, "[architecture/decision] Decidido manter SQLite com facts e notes.") {
		t.Fatalf("expected architecture note, got %q", got)
	}
	if !strings.Contains(got, "[memory/decision] Embeddings ficam opcionais, nao no core.") {
		t.Fatalf("expected memory note, got %q", got)
	}
}


