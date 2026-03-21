package persona

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func newTestCanonicalService(t *testing.T) *CanonicalIdentityService {
	t.Helper()
	dir := t.TempDir()
	identityPath := filepath.Join(dir, "IDENTITY.md")
	soulPath := filepath.Join(dir, "SOUL.md")
	userPath := filepath.Join(dir, "USER.md")

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
	if err := os.WriteFile(soulPath, []byte("# Soul\nBase.\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(userPath, []byte("# User\nNome: Nao definido\nFuso horario: Relativo a sua localidade.\n"), 0644); err != nil {
		t.Fatal(err)
	}

	return NewCanonicalIdentityService(identityPath, soulPath, userPath, "", "", "")
}

func TestCanonicalIdentityService_BuildPrompt_InjectsCurrentLocalDate(t *testing.T) {
	service := newTestCanonicalService(t)
	loc := time.FixedZone("America/Sao_Paulo", -3*60*60)
	service.location = loc
	service.now = func() time.Time {
		return time.Date(2026, time.March, 13, 9, 45, 0, 0, time.UTC)
	}

	prompt, _, err := service.BuildPrompt(context.Background(), "42", "42")
	if err != nil {
		t.Fatalf("BuildPrompt() error = %v", err)
	}

	if !strings.Contains(prompt, "Data local atual: 2026-03-13") {
		t.Fatalf("expected current local date in prompt, got %q", prompt)
	}
	if !strings.Contains(prompt, "Horario local atual: 2026-03-13T06:45:00-03:00") {
		t.Fatalf("expected localized timestamp in prompt, got %q", prompt)
	}
	if !strings.Contains(prompt, "Fuso horario atual: America/Sao_Paulo") {
		t.Fatalf("expected timezone in prompt, got %q", prompt)
	}
}

func TestCanonicalIdentityService_BuildPrompt_InjectsOwnerPlaybook(t *testing.T) {
	dir := t.TempDir()
	playbookPath := filepath.Join(dir, "OWNER_PLAYBOOK.md")
	playbookContent := "Always be direct and concise."
	if err := os.WriteFile(playbookPath, []byte(playbookContent), 0644); err != nil {
		t.Fatal(err)
	}

	service := newTestCanonicalServiceWithOwnerDocs(t, playbookPath, "")
	prompt, _, err := service.BuildPromptForQuery(context.Background(), "42", "42", "test query")
	if err != nil {
		t.Fatalf("BuildPromptForQuery() error = %v", err)
	}
	if !strings.Contains(prompt, "# OWNER CONTEXT") {
		t.Fatalf("expected OWNER CONTEXT section in prompt, got %q", prompt)
	}
	if !strings.Contains(prompt, playbookContent) {
		t.Fatalf("expected playbook content in prompt, got %q", prompt)
	}
}

func TestCanonicalIdentityService_BuildPrompt_InjectsProjectPlaybook(t *testing.T) {
	dir := t.TempDir()
	playbookPath := filepath.Join(dir, "PROJECT_PLAYBOOK.md")
	playbookContent := "Use tabs not spaces."
	if err := os.WriteFile(playbookPath, []byte(playbookContent), 0644); err != nil {
		t.Fatal(err)
	}

	service := newTestCanonicalServiceWithProjectPlaybook(t, "", "", playbookPath)
	prompt, _, err := service.BuildPromptForQuery(context.Background(), "42", "42", "test query")
	if err != nil {
		t.Fatalf("BuildPromptForQuery() error = %v", err)
	}
	if !strings.Contains(prompt, "# PROJECT CONTEXT") {
		t.Fatalf("expected PROJECT CONTEXT section in prompt, got %q", prompt)
	}
	if !strings.Contains(prompt, playbookContent) {
		t.Fatalf("expected project playbook content in prompt, got %q", prompt)
	}
}

func TestCanonicalIdentityService_BuildPrompt_ToleratesAbsentOwnerDocs(t *testing.T) {
	dir := t.TempDir()
	service := newTestCanonicalServiceWithOwnerDocs(t,
		filepath.Join(dir, "nonexistent_OWNER_PLAYBOOK.md"),
		filepath.Join(dir, "nonexistent_LESSONS_LEARNED.md"),
	)
	prompt, _, err := service.BuildPromptForQuery(context.Background(), "42", "42", "test query")
	if err != nil {
		t.Fatalf("BuildPromptForQuery() error = %v", err)
	}
	if strings.Contains(prompt, "# OWNER CONTEXT") {
		t.Fatalf("expected NO OWNER CONTEXT section when files absent, got %q", prompt)
	}
}

func newTestCanonicalServiceWithOwnerDocs(t *testing.T, ownerPlaybookPath, lessonsLearnedPath string) *CanonicalIdentityService {
	t.Helper()
	dir := t.TempDir()
	identityPath := filepath.Join(dir, "IDENTITY.md")
	soulPath := filepath.Join(dir, "SOUL.md")
	userPath := filepath.Join(dir, "USER.md")

	identityContent := `---
name: "Lex"
role: "Team Lead"
---
IDENTITY_BODY`
	if err := os.WriteFile(identityPath, []byte(identityContent), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(soulPath, []byte("# Soul\nBase.\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(userPath, []byte("# User\nNome: Nao definido\n"), 0644); err != nil {
		t.Fatal(err)
	}

	return NewCanonicalIdentityService(identityPath, soulPath, userPath, ownerPlaybookPath, lessonsLearnedPath, "")
}

func newTestCanonicalServiceWithProjectPlaybook(t *testing.T, ownerPlaybookPath, lessonsLearnedPath, projectPlaybookPath string) *CanonicalIdentityService {
	t.Helper()
	dir := t.TempDir()
	identityPath := filepath.Join(dir, "IDENTITY.md")
	soulPath := filepath.Join(dir, "SOUL.md")
	userPath := filepath.Join(dir, "USER.md")

	identityContent := `---
name: "Lex"
role: "Team Lead"
---
IDENTITY_BODY`
	if err := os.WriteFile(identityPath, []byte(identityContent), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(soulPath, []byte("# Soul\nBase.\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(userPath, []byte("# User\nNome: Nao definido\n"), 0644); err != nil {
		t.Fatal(err)
	}

	return NewCanonicalIdentityService(identityPath, soulPath, userPath, ownerPlaybookPath, lessonsLearnedPath, projectPlaybookPath)
}

// Tests for memory-backed features (facts, notes, retrieval, archive reprocess) were removed
// because they depend on internal/memory which was deleted.
