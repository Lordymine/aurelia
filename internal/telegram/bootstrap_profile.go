package telegram

import (
	"fmt"
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

	return fmt.Sprintf("# User\nNome: %s\nFuso horario: Relativo a sua localidade.\n", name)
}

func buildUserTemplateFromProfile(profileText, fallbackName string) string {
	name := extractNameFromProfile(profileText)
	if name == "" {
		name = strings.TrimSpace(fallbackName)
	}
	if name == "" {
		name = "Nao definido"
	}

	return fmt.Sprintf("# User\nNome: %s\nFuso horario: Relativo a sua localidade.\nPreferencias: %s\n", name, strings.TrimSpace(profileText))
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


