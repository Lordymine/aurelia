package telegram

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kocar/aurelia/internal/config"
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

	if bootstrapIdentityExists(dir) {
		t.Fatal("expected false for empty dir, got true")
	}

	if err := os.WriteFile(filepath.Join(dir, "IDENTITY.md"), []byte("# Identity"), 0o644); err != nil {
		t.Fatalf("failed to create IDENTITY.md: %v", err)
	}
	if !bootstrapIdentityExists(dir) {
		t.Fatal("expected true after creating IDENTITY.md, got false")
	}

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

	if _, err := os.Stat(filepath.Join(dir, "IDENTITY.md")); err != nil {
		t.Fatalf("IDENTITY.md not found in dir: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "SOUL.md")); err != nil {
		t.Fatalf("SOUL.md not found in dir: %v", err)
	}

	if _, err := os.Stat("IDENTITY.md"); err == nil {
		t.Fatal("IDENTITY.md must not be written to CWD")
	}
	if _, err := os.Stat("SOUL.md"); err == nil {
		t.Fatal("SOUL.md must not be written to CWD")
	}
}
