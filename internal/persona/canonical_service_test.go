package persona

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kocar/aurelia/internal/memory"
)

func TestCanonicalIdentityService_ConversationDoesNotOverrideBootstrap(t *testing.T) {
	mem := setupCanonicalMemory(t)
	service := newTestCanonicalService(t, mem)
	ctx := context.Background()

	err := service.ApplyFacts(ctx, "42", map[string]memory.Fact{
		"user.name": {Scope: "user", EntityID: "42", Key: "user.name", Value: "Rafael", Source: "bootstrap"},
	})
	if err != nil {
		t.Fatalf("ApplyFacts bootstrap error = %v", err)
	}
	err = service.ApplyFacts(ctx, "42", map[string]memory.Fact{
		"user.name": {Scope: "user", EntityID: "42", Key: "user.name", Value: "Rafa", Source: "conversation"},
	})
	if err != nil {
		t.Fatalf("ApplyFacts conversation error = %v", err)
	}

	fact, ok, err := mem.GetFact(ctx, "user", "42", "user.name")
	if err != nil || !ok {
		t.Fatalf("expected user.name fact, ok=%t err=%v", ok, err)
	}
	if fact.Value != "Rafael" {
		t.Fatalf("expected bootstrap fact to win, got %q", fact.Value)
	}
}

func TestCanonicalIdentityService_CorrectionOverridesBootstrap(t *testing.T) {
	mem := setupCanonicalMemory(t)
	service := newTestCanonicalService(t, mem)
	ctx := context.Background()

	_ = service.ApplyFacts(ctx, "42", map[string]memory.Fact{
		"user.name": {Scope: "user", EntityID: "42", Key: "user.name", Value: "Rafael", Source: "bootstrap"},
	})
	err := service.ApplyFacts(ctx, "42", map[string]memory.Fact{
		"user.name": {Scope: "user", EntityID: "42", Key: "user.name", Value: "Rafa", Source: "correction"},
	})
	if err != nil {
		t.Fatalf("ApplyFacts correction error = %v", err)
	}

	fact, ok, err := mem.GetFact(ctx, "user", "42", "user.name")
	if err != nil || !ok {
		t.Fatalf("expected user.name fact, ok=%t err=%v", ok, err)
	}
	if fact.Value != "Rafa" {
		t.Fatalf("expected correction fact to win, got %q", fact.Value)
	}
}

func TestCanonicalIdentityService_SyncFactsToFiles(t *testing.T) {
	mem := setupCanonicalMemory(t)
	service := newTestCanonicalService(t, mem)

	err := service.SyncFactsToFiles(map[string]memory.Fact{
		"user.name":                {Key: "user.name", Value: "Rafael"},
		"user.preferences.summary": {Key: "user.preferences.summary", Value: "Prefiro respostas diretas."},
		"agent.style.summary":      {Key: "agent.style.summary", Value: "Quero que voce seja mais direto."},
	})
	if err != nil {
		t.Fatalf("SyncFactsToFiles() error = %v", err)
	}

	userContent, _ := os.ReadFile(service.userPath)
	if !strings.Contains(string(userContent), "Nome: Rafael") {
		t.Fatalf("expected synced user file, got %q", string(userContent))
	}

	soulContent, _ := os.ReadFile(service.soulPath)
	if !strings.Contains(string(soulContent), "Preferencias persistentes do usuario sobre tom: Quero que voce seja mais direto.") {
		t.Fatalf("expected synced soul file, got %q", string(soulContent))
	}
}

func TestCanonicalIdentityService_BuildPromptUsesFacts(t *testing.T) {
	mem := setupCanonicalMemory(t)
	service := newTestCanonicalService(t, mem)
	ctx := context.Background()

	_ = mem.UpsertFact(ctx, memory.Fact{Scope: "agent", EntityID: "default", Key: "agent.role", Value: "Chief Architect", Source: "correction"})
	_ = mem.UpsertFact(ctx, memory.Fact{Scope: "user", EntityID: "42", Key: "user.name", Value: "Rafael", Source: "bootstrap"})
	_ = mem.AddNote(ctx, memory.Note{ConversationID: "42", Topic: "architecture", Kind: "decision", Summary: "Decidido manter o team lead fixo.", Importance: 8, Source: "conversation"})

	prompt, tools, err := service.BuildPrompt(ctx, "42", "42")
	if err != nil {
		t.Fatalf("BuildPrompt() error = %v", err)
	}
	if !strings.Contains(prompt, "Papel canonico do agente: Chief Architect") {
		t.Fatalf("expected resolved agent role, got %q", prompt)
	}
	if !strings.Contains(prompt, "Nome canonico do usuario: Rafael") {
		t.Fatalf("expected resolved user name, got %q", prompt)
	}
	if !strings.Contains(prompt, "# LONG-TERM MEMORY") {
		t.Fatalf("expected notes block, got %q", prompt)
	}
	if !strings.Contains(prompt, "# RUNTIME CONTEXT") {
		t.Fatalf("expected runtime context block, got %q", prompt)
	}
	if tools != nil {
		t.Fatalf("expected runtime tools to be resolved outside persona, got %v", tools)
	}
}

func TestCanonicalIdentityService_BuildPrompt_InjectsCurrentLocalDate(t *testing.T) {
	mem := setupCanonicalMemory(t)
	service := newTestCanonicalService(t, mem)
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

func TestCanonicalIdentityService_BuildPromptForQuery_RetrievesRelevantLongTermMemory(t *testing.T) {
	mem := setupCanonicalMemory(t)
	service := newTestCanonicalService(t, mem)
	ctx := context.Background()

	_ = mem.UpsertFact(ctx, memory.Fact{Scope: "project", EntityID: "default", Key: "project.memory.strategy", Value: "sqlite + facts + notes", Source: "conversation"})
	_ = mem.UpsertFact(ctx, memory.Fact{Scope: "user", EntityID: "42", Key: "user.preference.response_style", Value: "direto", Source: "conversation"})
	_ = mem.AddNote(ctx, memory.Note{ConversationID: "42", Topic: "memory", Kind: "decision", Summary: "Foi decidido evitar embeddings no core.", Importance: 8, Source: "conversation"})
	_ = mem.AddNote(ctx, memory.Note{ConversationID: "42", Topic: "frontend", Kind: "decision", Summary: "Usar layout bold.", Importance: 8, Source: "conversation"})

	prompt, _, err := service.BuildPromptForQuery(ctx, "42", "42", "o que decidimos sobre memoria de longo prazo?")
	if err != nil {
		t.Fatalf("BuildPromptForQuery() error = %v", err)
	}

	if !strings.Contains(prompt, "project.memory.strategy: sqlite + facts + notes") {
		t.Fatalf("expected relevant fact in prompt, got %q", prompt)
	}
	if !strings.Contains(prompt, "[memory/decision] Foi decidido evitar embeddings no core.") {
		t.Fatalf("expected relevant note in prompt, got %q", prompt)
	}
	if strings.Contains(prompt, "[frontend/decision] Usar layout bold.") {
		t.Fatalf("expected unrelated note to stay out of prompt, got %q", prompt)
	}
}

func TestCanonicalIdentityService_ReprocessArchivePromotesFactsAndNotes(t *testing.T) {
	mem := setupCanonicalMemory(t)
	service := newTestCanonicalService(t, mem)
	ctx := context.Background()

	_ = mem.AddArchiveEntry(ctx, memory.ArchiveEntry{
		ConversationID: "42",
		SessionID:      "42",
		Role:           "user",
		Content:        "Meu nome e Rafael e prefiro respostas diretas, sem floreios.",
		MessageType:    "chat",
	})
	_ = mem.AddArchiveEntry(ctx, memory.ArchiveEntry{
		ConversationID: "42",
		SessionID:      "42",
		Role:           "user",
		Content:        "Fica decidido que vamos manter SQLite com facts e notes.",
		MessageType:    "chat",
	})

	if err := service.ReprocessArchive(ctx, "42", "42", 20); err != nil {
		t.Fatalf("ReprocessArchive() error = %v", err)
	}

	nameFact, ok, err := mem.GetFact(ctx, "user", "42", "user.name")
	if err != nil || !ok {
		t.Fatalf("expected user.name fact, ok=%t err=%v", ok, err)
	}
	if nameFact.Value != "Rafael" {
		t.Fatalf("expected promoted user name, got %q", nameFact.Value)
	}

	notes, err := mem.ListRecentNotes(ctx, "42", 10)
	if err != nil {
		t.Fatalf("ListRecentNotes() error = %v", err)
	}
	if len(notes) == 0 {
		t.Fatal("expected promoted note from archive")
	}
}

func TestCanonicalIdentityService_ReprocessArchive_IsIdempotentForNotes(t *testing.T) {
	mem := setupCanonicalMemory(t)
	service := newTestCanonicalService(t, mem)
	ctx := context.Background()

	_ = mem.AddArchiveEntry(ctx, memory.ArchiveEntry{
		ConversationID: "42",
		SessionID:      "42",
		Role:           "user",
		Content:        "Fica decidido que vamos manter SQLite com facts e notes.",
		MessageType:    "chat",
	})

	if err := service.ReprocessArchive(ctx, "42", "42", 20); err != nil {
		t.Fatalf("first ReprocessArchive() error = %v", err)
	}
	if err := service.ReprocessArchive(ctx, "42", "42", 20); err != nil {
		t.Fatalf("second ReprocessArchive() error = %v", err)
	}

	notes, err := mem.ListRecentNotes(ctx, "42", 10)
	if err != nil {
		t.Fatalf("ListRecentNotes() error = %v", err)
	}
	if len(notes) != 1 {
		t.Fatalf("expected one deduplicated note after reprocess, got %d", len(notes))
	}
}

func TestCanonicalIdentityService_ReprocessArchivePromotesProjectFacts(t *testing.T) {
	mem := setupCanonicalMemory(t)
	service := newTestCanonicalService(t, mem)
	ctx := context.Background()

	_ = mem.AddArchiveEntry(ctx, memory.ArchiveEntry{
		ConversationID: "42",
		SessionID:      "42",
		Role:           "user",
		Content:        "Quero esse projeto minimalista, com monolito modular, usando SQLite com facts e notes, e nao quero depender de embeddings no core.",
		MessageType:    "chat",
	})

	if err := service.ReprocessArchive(ctx, "42", "42", 20); err != nil {
		t.Fatalf("ReprocessArchive() error = %v", err)
	}

	goalFact, ok, err := mem.GetFact(ctx, "project", "default", "project.goal")
	if err != nil || !ok {
		t.Fatalf("expected project.goal fact, ok=%t err=%v", ok, err)
	}
	if goalFact.Value != "manter o projeto minimalista" {
		t.Fatalf("expected project goal fact, got %q", goalFact.Value)
	}

	constraintFact, ok, err := mem.GetFact(ctx, "project", "default", "project.constraint.no_embeddings_core")
	if err != nil || !ok {
		t.Fatalf("expected constraint fact, ok=%t err=%v", ok, err)
	}
	if constraintFact.Value != "true" {
		t.Fatalf("expected no_embeddings_core=true, got %q", constraintFact.Value)
	}

	architectureFact, ok, err := mem.GetFact(ctx, "project", "default", "project.architecture.style")
	if err != nil || !ok {
		t.Fatalf("expected architecture fact, ok=%t err=%v", ok, err)
	}
	if architectureFact.Value != "monolito modular" {
		t.Fatalf("expected architecture style fact, got %q", architectureFact.Value)
	}

	memoryStrategy, ok, err := mem.GetFact(ctx, "project", "default", "project.memory.strategy")
	if err != nil || !ok {
		t.Fatalf("expected memory strategy fact, ok=%t err=%v", ok, err)
	}
	if memoryStrategy.Value != "sqlite + facts + notes" {
		t.Fatalf("expected memory strategy fact, got %q", memoryStrategy.Value)
	}
}

func TestCanonicalIdentityService_BuildPromptForQuery_PrioritizesThematicProjectFacts(t *testing.T) {
	mem := setupCanonicalMemory(t)
	service := newTestCanonicalService(t, mem)
	ctx := context.Background()

	_ = mem.UpsertFact(ctx, memory.Fact{Scope: "project", EntityID: "default", Key: "project.memory.strategy", Value: "sqlite + facts + notes", Source: "conversation"})
	_ = mem.UpsertFact(ctx, memory.Fact{Scope: "project", EntityID: "default", Key: "project.architecture.style", Value: "monolito modular", Source: "conversation"})
	_ = mem.UpsertFact(ctx, memory.Fact{Scope: "project", EntityID: "default", Key: "project.constraint.minimalism", Value: "true", Source: "conversation"})

	prompt, _, err := service.BuildPromptForQuery(ctx, "42", "42", "qual nossa estrategia de memoria e arquitetura minimalista?")
	if err != nil {
		t.Fatalf("BuildPromptForQuery() error = %v", err)
	}

	if !strings.Contains(prompt, "project.memory.strategy: sqlite + facts + notes") {
		t.Fatalf("expected memory strategy in prompt, got %q", prompt)
	}
	if !strings.Contains(prompt, "project.architecture.style: monolito modular") {
		t.Fatalf("expected architecture style in prompt, got %q", prompt)
	}
	if !strings.Contains(prompt, "project.constraint.minimalism: true") {
		t.Fatalf("expected minimalism constraint in prompt, got %q", prompt)
	}
}

func TestCanonicalIdentityService_DebugLongTermMemory(t *testing.T) {
	mem := setupCanonicalMemory(t)
	service := newTestCanonicalService(t, mem)
	ctx := context.Background()

	_ = mem.UpsertFact(ctx, memory.Fact{Scope: "project", EntityID: "default", Key: "project.memory.strategy", Value: "sqlite + facts + notes", Source: "conversation"})
	_ = mem.AddNote(ctx, memory.Note{ConversationID: "42", Topic: "memory", Kind: "decision", Summary: "Foi decidido evitar embeddings no core.", Importance: 8, Source: "conversation"})

	report, err := service.DebugLongTermMemory(ctx, "42", "42", "qual nossa estrategia de memoria?")
	if err != nil {
		t.Fatalf("DebugLongTermMemory() error = %v", err)
	}

	if len(report.Tokens) == 0 {
		t.Fatal("expected query tokens in debug report")
	}
	if len(report.SelectedFacts) == 0 {
		t.Fatal("expected selected facts in debug report")
	}
	if report.SelectedFacts[0].Fact.Key != "project.memory.strategy" {
		t.Fatalf("unexpected selected fact %+v", report.SelectedFacts[0])
	}
	if report.SelectedFacts[0].Score <= 0 {
		t.Fatalf("expected positive fact score, got %+v", report.SelectedFacts[0])
	}
	if len(report.SelectedNotes) == 0 {
		t.Fatal("expected selected notes in debug report")
	}
	if report.SelectedNotes[0].Score <= 0 {
		t.Fatalf("expected positive note score, got %+v", report.SelectedNotes[0])
	}
}

func setupCanonicalMemory(t *testing.T) *memory.MemoryManager {
	t.Helper()
	mem, err := memory.NewMemoryManager(filepath.Join(t.TempDir(), "canonical.db"), 5)
	if err != nil {
		t.Fatalf("NewMemoryManager() error = %v", err)
	}
	t.Cleanup(func() { _ = mem.Close() })
	return mem
}

func newTestCanonicalService(t *testing.T, mem *memory.MemoryManager) *CanonicalIdentityService {
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

	return NewCanonicalIdentityService(mem, identityPath, soulPath, userPath, "", "", "")
}

func TestCanonicalIdentityService_BuildPrompt_InjectsOwnerPlaybook(t *testing.T) {
	mem := setupCanonicalMemory(t)
	dir := t.TempDir()
	playbookPath := filepath.Join(dir, "OWNER_PLAYBOOK.md")
	playbookContent := "Always be direct and concise."
	if err := os.WriteFile(playbookPath, []byte(playbookContent), 0644); err != nil {
		t.Fatal(err)
	}

	service := newTestCanonicalServiceWithOwnerDocs(t, mem, playbookPath, "")
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

func TestCanonicalIdentityService_BuildPrompt_InjectsLessonsLearned(t *testing.T) {
	mem := setupCanonicalMemory(t)
	dir := t.TempDir()
	lessonsPath := filepath.Join(dir, "LESSONS_LEARNED.md")
	lessonsContent := "Never skip error handling."
	if err := os.WriteFile(lessonsPath, []byte(lessonsContent), 0644); err != nil {
		t.Fatal(err)
	}

	service := newTestCanonicalServiceWithOwnerDocs(t, mem, "", lessonsPath)
	prompt, _, err := service.BuildPromptForQuery(context.Background(), "42", "42", "test query")
	if err != nil {
		t.Fatalf("BuildPromptForQuery() error = %v", err)
	}
	if !strings.Contains(prompt, lessonsContent) {
		t.Fatalf("expected lessons content in prompt, got %q", prompt)
	}
}

func TestCanonicalIdentityService_BuildPrompt_ToleratesAbsentOwnerDocs(t *testing.T) {
	mem := setupCanonicalMemory(t)
	dir := t.TempDir()
	service := newTestCanonicalServiceWithOwnerDocs(t, mem,
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

func newTestCanonicalServiceWithOwnerDocs(t *testing.T, mem *memory.MemoryManager, ownerPlaybookPath, lessonsLearnedPath string) *CanonicalIdentityService {
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

	return NewCanonicalIdentityService(mem, identityPath, soulPath, userPath, ownerPlaybookPath, lessonsLearnedPath, "")
}

func newTestCanonicalServiceWithProjectPlaybook(t *testing.T, mem *memory.MemoryManager, ownerPlaybookPath, lessonsLearnedPath, projectPlaybookPath string) *CanonicalIdentityService {
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

	return NewCanonicalIdentityService(mem, identityPath, soulPath, userPath, ownerPlaybookPath, lessonsLearnedPath, projectPlaybookPath)
}

func TestCanonicalIdentityService_BuildPrompt_InjectsProjectPlaybook(t *testing.T) {
	mem := setupCanonicalMemory(t)
	dir := t.TempDir()
	playbookPath := filepath.Join(dir, "PROJECT_PLAYBOOK.md")
	playbookContent := "Use tabs not spaces."
	if err := os.WriteFile(playbookPath, []byte(playbookContent), 0644); err != nil {
		t.Fatal(err)
	}

	service := newTestCanonicalServiceWithProjectPlaybook(t, mem, "", "", playbookPath)
	prompt, _, err := service.BuildPromptForQuery(context.Background(), "42", "42", "test query")
	if err != nil {
		t.Fatalf("BuildPromptForQuery() error = %v", err)
	}
	if !strings.Contains(prompt, "# PROJECT CONTEXT") {
		t.Fatalf("expected PROJECT CONTEXT section in prompt, got %q", prompt)
	}
	if !strings.Contains(prompt, "## Project Playbook") {
		t.Fatalf("expected Project Playbook subsection in prompt, got %q", prompt)
	}
	if !strings.Contains(prompt, playbookContent) {
		t.Fatalf("expected project playbook content in prompt, got %q", prompt)
	}
}

func TestCanonicalIdentityService_BuildPrompt_ToleratesAbsentProjectPlaybook(t *testing.T) {
	mem := setupCanonicalMemory(t)
	dir := t.TempDir()
	service := newTestCanonicalServiceWithProjectPlaybook(t, mem, "", "", filepath.Join(dir, "nonexistent_PROJECT_PLAYBOOK.md"))
	prompt, _, err := service.BuildPromptForQuery(context.Background(), "42", "42", "test query")
	if err != nil {
		t.Fatalf("BuildPromptForQuery() error = %v", err)
	}
	if strings.Contains(prompt, "# PROJECT CONTEXT") {
		t.Fatalf("expected NO PROJECT CONTEXT section when file absent, got %q", prompt)
	}
}

func TestCanonicalIdentityService_BuildPrompt_ToleratesEmptyProjectPath(t *testing.T) {
	mem := setupCanonicalMemory(t)
	service := newTestCanonicalServiceWithProjectPlaybook(t, mem, "", "", "")
	prompt, _, err := service.BuildPromptForQuery(context.Background(), "42", "42", "test query")
	if err != nil {
		t.Fatalf("BuildPromptForQuery() error = %v", err)
	}
	if strings.Contains(prompt, "# PROJECT CONTEXT") {
		t.Fatalf("expected NO PROJECT CONTEXT section when path is empty, got %q", prompt)
	}
}

func TestCanonicalIdentityService_BuildPrompt_ProjectPlaybookNotPersistedToMemory(t *testing.T) {
	mem := setupCanonicalMemory(t)
	dir := t.TempDir()
	playbookPath := filepath.Join(dir, "PROJECT_PLAYBOOK.md")
	playbookContent := "Use tabs not spaces."
	if err := os.WriteFile(playbookPath, []byte(playbookContent), 0644); err != nil {
		t.Fatal(err)
	}

	service := newTestCanonicalServiceWithProjectPlaybook(t, mem, "", "", playbookPath)
	_, _, err := service.BuildPromptForQuery(context.Background(), "42", "42", "test query")
	if err != nil {
		t.Fatalf("BuildPromptForQuery() error = %v", err)
	}

	ctx := context.Background()
	facts, err := mem.ListFacts(ctx, "user", "42")
	if err != nil {
		t.Fatalf("ListFacts() error = %v", err)
	}
	for _, f := range facts {
		if strings.Contains(f.Value, playbookContent) {
			t.Fatalf("project playbook content was persisted to memory facts: %+v", f)
		}
	}

	notes, err := mem.ListRecentNotes(ctx, "42", 100)
	if err != nil {
		t.Fatalf("ListRecentNotes() error = %v", err)
	}
	for _, n := range notes {
		if strings.Contains(n.Summary, playbookContent) {
			t.Fatalf("project playbook content was persisted to memory notes: %+v", n)
		}
	}
}


