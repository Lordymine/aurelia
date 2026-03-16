package telegram

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/memory"
	"gopkg.in/telebot.v3"
)

func TestBuildUserTemplate_UsesTelegramName(t *testing.T) {
	user := &telebot.User{
		ID:        42,
		FirstName: "Rafael",
		LastName:  "Kocar",
		Username:  "rafa",
	}

	got := buildUserTemplate(user)

	if !strings.Contains(got, "Nome: Rafael Kocar") {
		t.Fatalf("expected full name in user template, got %q", got)
	}
	if strings.Contains(got, "Usuario 42") {
		t.Fatalf("should not write numeric placeholder, got %q", got)
	}
}

func TestBuildUserTemplate_FallsBackWithoutInventingIdentity(t *testing.T) {
	user := &telebot.User{ID: 42}

	got := buildUserTemplate(user)

	if !strings.Contains(got, "Nome: Nao definido") {
		t.Fatalf("expected unresolved placeholder, got %q", got)
	}
	if strings.Contains(got, "Usuario 42") {
		t.Fatalf("should not invent identity from telegram id, got %q", got)
	}
}

func TestBuildUserTemplateFromProfile_UsesConversationProfile(t *testing.T) {
	got := buildUserTemplateFromProfile("Me chamo Rafael e quero respostas diretas, sem floreios.", "rafa")

	if !strings.Contains(got, "Nome: Rafael") {
		t.Fatalf("expected extracted name, got %q", got)
	}
	if !strings.Contains(got, "Preferencias: Me chamo Rafael e quero respostas diretas, sem floreios.") {
		t.Fatalf("expected full profile text, got %q", got)
	}
}

func TestBuildUserTemplateFromProfile_FallsBackToTelegramName(t *testing.T) {
	got := buildUserTemplateFromProfile("Quero respostas diretas, sem floreios.", "rafa")

	if !strings.Contains(got, "Nome: rafa") {
		t.Fatalf("expected telegram fallback name, got %q", got)
	}
}

func TestExtractNameFromProfile(t *testing.T) {
	cases := map[string]string{
		"Me chamo Rafael e quero respostas diretas.": "Rafael",
		"Meu nome e Rafael Kocar.":                   "Rafael Kocar",
		"Sou Rafael.":                                "Rafael",
		"Quero respostas diretas.":                   "",
	}

	for input, want := range cases {
		if got := extractNameFromProfile(input); got != want {
			t.Fatalf("extractNameFromProfile(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestSeedBootstrapFacts_PersistsCanonicalFacts(t *testing.T) {
	mem := memoryForBootstrapTest(t)
	ctx := context.Background()

	if err := seedBootstrapFacts(ctx, mem, "42", "Lex", "Team Lead", "Rafael"); err != nil {
		t.Fatalf("seedBootstrapFacts() error = %v", err)
	}

	agentName, ok, err := mem.GetFact(ctx, "agent", "default", "agent.name")
	if err != nil || !ok {
		t.Fatalf("expected agent.name fact, ok=%t err=%v", ok, err)
	}
	if agentName.Value != "Lex" {
		t.Fatalf("expected agent name Lex, got %q", agentName.Value)
	}

	userName, ok, err := mem.GetFact(ctx, "user", "42", "user.name")
	if err != nil || !ok {
		t.Fatalf("expected user.name fact, ok=%t err=%v", ok, err)
	}
	if userName.Value != "Rafael" {
		t.Fatalf("expected user name Rafael, got %q", userName.Value)
	}
}

func TestExtractFactsFromConversation_PromotesUserName(t *testing.T) {
	facts := extractFactsFromConversation("Meu nome e Rafael Kocar.", "42")

	fact, ok := facts["user.name"]
	if !ok {
		t.Fatal("expected user.name fact")
	}
	if fact.Value != "Rafael Kocar" {
		t.Fatalf("expected user name fact, got %q", fact.Value)
	}
}

func TestExtractFactsFromConversation_PromotesUserPreference(t *testing.T) {
	facts := extractFactsFromConversation("Prefiro respostas diretas e sem floreios.", "42")

	fact, ok := facts["user.preferences.summary"]
	if !ok {
		t.Fatal("expected user.preferences.summary fact")
	}
	if fact.Value != "Prefiro respostas diretas e sem floreios." {
		t.Fatalf("unexpected preference summary %q", fact.Value)
	}
}

func TestExtractFactsFromConversation_PromotesAgentCorrections(t *testing.T) {
	facts := extractFactsFromConversation("Seu nome e Lex e seu papel e Arquiteto Chefe.", "42")

	nameFact, ok := facts["agent.name"]
	if !ok || nameFact.Value != "Lex" {
		t.Fatalf("expected agent.name Lex, got %+v", nameFact)
	}

	roleFact, ok := facts["agent.role"]
	if !ok || roleFact.Value != "Arquiteto Chefe" {
		t.Fatalf("expected agent.role Arquiteto Chefe, got %+v", roleFact)
	}
}

func TestPersistConversationFacts_UpsertsExtractedFacts(t *testing.T) {
	mem := memoryForBootstrapTest(t)
	ctx := context.Background()

	if err := persistConversationFacts(ctx, mem, "42", "Meu nome e Rafael e prefiro respostas diretas."); err != nil {
		t.Fatalf("persistConversationFacts() error = %v", err)
	}

	nameFact, ok, err := mem.GetFact(ctx, "user", "42", "user.name")
	if err != nil || !ok {
		t.Fatalf("expected user.name fact, ok=%t err=%v", ok, err)
	}
	if nameFact.Value != "Rafael" {
		t.Fatalf("expected Rafael, got %q", nameFact.Value)
	}

	prefFact, ok, err := mem.GetFact(ctx, "user", "42", "user.preferences.summary")
	if err != nil || !ok {
		t.Fatalf("expected user.preferences.summary fact, ok=%t err=%v", ok, err)
	}
	if prefFact.Value != "Meu nome e Rafael e prefiro respostas diretas." {
		t.Fatalf("unexpected preference fact %q", prefFact.Value)
	}
}

func TestExtractFactsFromConversation_PromotesSpecificPreferences(t *testing.T) {
	facts := extractFactsFromConversation("Prefiro respostas diretas, sem floreios e sem emojis.", "42")

	if fact, ok := facts["user.preference.response_style"]; !ok || fact.Value != "direto" {
		t.Fatalf("expected response style fact, got %+v", fact)
	}
	if fact, ok := facts["user.preference.no_fluff"]; !ok || fact.Value != "true" {
		t.Fatalf("expected no_fluff fact, got %+v", fact)
	}
	if fact, ok := facts["user.preference.no_emoji"]; !ok || fact.Value != "true" {
		t.Fatalf("expected no_emoji fact, got %+v", fact)
	}
}

func TestExtractFactsFromConversation_PromotesAgentStylePreference(t *testing.T) {
	facts := extractFactsFromConversation("Quero que voce seja mais direto e sarcastico.", "42")

	if fact, ok := facts["agent.style.summary"]; !ok || fact.Value != "Quero que voce seja mais direto e sarcastico." {
		t.Fatalf("expected agent.style.summary fact, got %+v", fact)
	}
}

func TestExtractFactsFromConversation_PromotesProjectFacts(t *testing.T) {
	facts := extractFactsFromConversation("Quero esse projeto minimalista, com monolito modular, usando SQLite com facts e notes, e nao quero depender de embeddings no core.", "42")

	if fact, ok := facts["project.goal"]; !ok || fact.Value != "manter o projeto minimalista" {
		t.Fatalf("expected project.goal fact, got %+v", fact)
	}
	if fact, ok := facts["project.constraint.no_embeddings_core"]; !ok || fact.Value != "true" {
		t.Fatalf("expected project.constraint.no_embeddings_core fact, got %+v", fact)
	}
	if fact, ok := facts["project.architecture.style"]; !ok || fact.Value != "monolito modular" {
		t.Fatalf("expected project.architecture.style fact, got %+v", fact)
	}
	if fact, ok := facts["project.memory.strategy"]; !ok || fact.Value != "sqlite + facts + notes" {
		t.Fatalf("expected project.memory.strategy fact, got %+v", fact)
	}
	if fact, ok := facts["project.constraint.minimalism"]; !ok || fact.Value != "true" {
		t.Fatalf("expected project.constraint.minimalism fact, got %+v", fact)
	}
}

func TestExtractArchitectureNoteFromConversation(t *testing.T) {
	note, ok := extractArchitectureNoteFromConversation("Fica decidido que vamos manter SQLite com facts e notes.", "42")
	if !ok {
		t.Fatal("expected architecture note")
	}
	if note.Topic != "architecture" || note.Kind != "decision" {
		t.Fatalf("unexpected note %+v", note)
	}
}

func TestSyncUserPersonaFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "USER.md")
	if err := os.WriteFile(path, []byte("# User\nNome: Nao definido\nFuso horario: Relativo a sua localidade.\n"), 0644); err != nil {
		t.Fatal(err)
	}

	facts := map[string]memory.Fact{
		"user.name":                      {Key: "user.name", Value: "Rafael"},
		"user.preferences.summary":      {Key: "user.preferences.summary", Value: "Prefiro respostas diretas."},
		"user.preference.response_style": {Key: "user.preference.response_style", Value: "direto"},
	}
	if err := syncUserPersonaFile(path, facts); err != nil {
		t.Fatalf("syncUserPersonaFile() error = %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	content := string(got)
	if !strings.Contains(content, "Nome: Rafael") {
		t.Fatalf("expected synced user name, got %q", content)
	}
	if !strings.Contains(content, "Preferencias: Prefiro respostas diretas.") {
		t.Fatalf("expected synced preference summary, got %q", content)
	}
	if !strings.Contains(content, "Estilo de resposta: direto") {
		t.Fatalf("expected synced response style, got %q", content)
	}
}

func TestSyncSoulPersonaFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "SOUL.md")
	if err := os.WriteFile(path, []byte("# Soul\nBase.\n"), 0644); err != nil {
		t.Fatal(err)
	}

	facts := map[string]memory.Fact{
		"agent.style.summary": {Key: "agent.style.summary", Value: "Quero que voce seja mais direto."},
	}
	if err := syncSoulPersonaFile(path, facts); err != nil {
		t.Fatalf("syncSoulPersonaFile() error = %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	content := string(got)
	if !strings.Contains(content, "Preferencias persistentes do usuario sobre tom: Quero que voce seja mais direto.") {
		t.Fatalf("expected synced soul preference, got %q", content)
	}
}

func memoryForBootstrapTest(t *testing.T) *memory.MemoryManager {
	t.Helper()
	mem, err := memory.NewMemoryManager(t.TempDir()+"\\bootstrap.db", 5)
	if err != nil {
		t.Fatalf("failed to create memory manager: %v", err)
	}
	t.Cleanup(func() {
		_ = mem.Close()
	})
	return mem
}

func TestBotController_IsAllowedUser(t *testing.T) {
	bc := &BotController{
		config: &config.AppConfig{
			TelegramAllowedUserIDs: []int64{42, 99},
		},
	}

	if !bc.isAllowedUser(42) {
		t.Fatal("expected user 42 to be allowed")
	}
	if bc.isAllowedUser(7) {
		t.Fatal("expected user 7 to be blocked")
	}
}

func TestBootstrapStartResponse_WhenAlreadyConfigured(t *testing.T) {
	message, menu := bootstrapStartResponse(true)

	if message != alreadyConfiguredMessage {
		t.Fatalf("unexpected configured message: %q", message)
	}
	if menu != nil {
		t.Fatalf("expected no menu when already configured, got %#v", menu)
	}
}

func TestBootstrapStartResponse_WhenBootstrapNeeded(t *testing.T) {
	message, menu := bootstrapStartResponse(false)

	if message != bootstrapWelcomeMessage {
		t.Fatalf("unexpected bootstrap welcome message: %q", message)
	}
	if menu == nil {
		t.Fatal("expected bootstrap menu")
	}
	if len(menu.InlineKeyboard) != 2 {
		t.Fatalf("expected two inline rows, got %d", len(menu.InlineKeyboard))
	}
	if len(menu.InlineKeyboard[0]) != 1 || len(menu.InlineKeyboard[1]) != 1 {
		t.Fatalf("expected one button per row, got %#v", menu.InlineKeyboard)
	}
	if menu.InlineKeyboard[0][0].Unique != "btn_coder" {
		t.Fatalf("expected coder callback button, got %#v", menu.InlineKeyboard[0][0])
	}
	if menu.InlineKeyboard[1][0].Unique != "btn_assist" {
		t.Fatalf("expected assist callback button, got %#v", menu.InlineKeyboard[1][0])
	}
}

func TestBootstrapIdentityExists_UsesGivenDir(t *testing.T) {
	dir := t.TempDir()

	// Initially no IDENTITY.md in dir — should return false.
	if bootstrapIdentityExists(dir) {
		t.Fatal("expected false for empty dir, got true")
	}

	// Create IDENTITY.md in dir — should return true.
	if err := os.WriteFile(filepath.Join(dir, "IDENTITY.md"), []byte("# Identity"), 0o644); err != nil {
		t.Fatalf("failed to create IDENTITY.md: %v", err)
	}
	if !bootstrapIdentityExists(dir) {
		t.Fatal("expected true after creating IDENTITY.md, got false")
	}

	// Different empty dir — should return false.
	if bootstrapIdentityExists(t.TempDir()) {
		t.Fatal("expected false for different empty dir, got true")
	}
}

func TestWriteBootstrapPreset_WritesToDir(t *testing.T) {
	dir := t.TempDir()
	preset, err := bootstrapPresetForChoice("coder")
	if err != nil {
		t.Fatalf("bootstrapPresetForChoice() error = %v", err)
	}

	if err := writeBootstrapPreset(dir, preset); err != nil {
		t.Fatalf("writeBootstrapPreset() error = %v", err)
	}

	// Files must exist in the given dir.
	if _, err := os.Stat(filepath.Join(dir, "IDENTITY.md")); err != nil {
		t.Fatalf("IDENTITY.md not found in dir: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "SOUL.md")); err != nil {
		t.Fatalf("SOUL.md not found in dir: %v", err)
	}

	// Files must NOT exist in CWD.
	if _, err := os.Stat("IDENTITY.md"); err == nil {
		t.Fatal("IDENTITY.md must not be written to CWD")
	}
	if _, err := os.Stat("SOUL.md"); err == nil {
		t.Fatal("SOUL.md must not be written to CWD")
	}
}


