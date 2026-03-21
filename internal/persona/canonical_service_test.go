package persona

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
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

func TestCanonicalIdentityService_BuildPrompt_ContainsPersonaSections(t *testing.T) {
	service := newTestCanonicalService(t)

	prompt, err := service.BuildPrompt()
	if err != nil {
		t.Fatalf("BuildPrompt() error = %v", err)
	}

	if !strings.Contains(prompt, "IDENTITY_BODY") {
		t.Fatalf("expected identity body in prompt, got %q", prompt)
	}
	if !strings.Contains(prompt, "# Soul") {
		t.Fatalf("expected soul content in prompt, got %q", prompt)
	}
	if !strings.Contains(prompt, "# User") {
		t.Fatalf("expected user content in prompt, got %q", prompt)
	}
	if !strings.Contains(prompt, "# CANONICAL IDENTITY") {
		t.Fatalf("expected canonical identity block in prompt, got %q", prompt)
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
	prompt, err := service.BuildPrompt()
	if err != nil {
		t.Fatalf("BuildPrompt() error = %v", err)
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
	prompt, err := service.BuildPrompt()
	if err != nil {
		t.Fatalf("BuildPrompt() error = %v", err)
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
	prompt, err := service.BuildPrompt()
	if err != nil {
		t.Fatalf("BuildPrompt() error = %v", err)
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
