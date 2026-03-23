package telegram

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kocar/aurelia/internal/persona"
	"gopkg.in/telebot.v3"
)

func buildUserTemplate(user *telebot.User) string {
	name := "Nao definido"
	if user != nil {
		fullName := strings.TrimSpace(strings.Join([]string{strings.TrimSpace(user.FirstName), strings.TrimSpace(user.LastName)}, " "))
		switch {
		case fullName != "":
			name = fullName
		case strings.TrimSpace(user.Username) != "":
			name = strings.TrimSpace(user.Username)
		}
	}

	return "# User\nNome: " + name + "\nFuso horario: Relativo a sua localidade.\n"
}

func buildUserTemplateFromProfile(profileText, fallbackName string) string {
	name := extractNameFromProfile(profileText)
	if name == "" {
		name = strings.TrimSpace(fallbackName)
	}
	if name == "" {
		name = "Nao definido"
	}

	return "# User\nNome: " + name + "\nFuso horario: Relativo a sua localidade.\nPreferencias: " + strings.TrimSpace(profileText) + "\n"
}

func extractNameFromProfile(profileText string) string {
	return persona.ExtractNameFromProfile(profileText)
}

func bootstrapFallbackName(user *telebot.User) string {
	if user == nil {
		return "Nao definido"
	}
	fallbackName := strings.TrimSpace(strings.Join([]string{strings.TrimSpace(user.FirstName), strings.TrimSpace(user.LastName)}, " "))
	if fallbackName == "" {
		fallbackName = strings.TrimSpace(user.Username)
	}
	if fallbackName == "" {
		return "Nao definido"
	}
	return fallbackName
}

func (bc *BotController) completeBootstrapProfile(c telebot.Context, state bootstrapState, text string) error {
	userTemplate := buildUserTemplateFromProfile(text, bootstrapFallbackName(c.Sender()))
	if err := os.WriteFile(filepath.Join(bc.personasDir, "USER.md"), []byte(userTemplate), 0o644); err != nil {
		log.Printf("Bootstrap user profile write error: %v\n", err)
		return SendContextText(c, bootstrapFailureMessage)
	}

	// TODO: seed identity facts via bridge when memory is wired
	return SendContextText(c, bootstrapSuccessMessage)
}

