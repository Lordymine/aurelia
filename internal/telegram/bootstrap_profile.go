package telegram

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kocar/aurelia/internal/bridge"
	"github.com/kocar/aurelia/internal/persona"
	"gopkg.in/telebot.v3"
)

const bootstrapGenerateTimeout = 30 * time.Second

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

// completeBootstrapAssistant uses the LLM to generate IDENTITY.md and SOUL.md
// from the user's description of the desired assistant personality, then
// transitions to the profile step.
func (bc *BotController) completeBootstrapAssistant(c telebot.Context, state bootstrapState, text string) error {
	preset, err := bootstrapPresetForChoice(state.Choice)
	if err != nil {
		return SendContextText(c, bootstrapFailureMessage)
	}

	stopTyping := startChatActionLoop(bc.bot, c.Chat(), telebot.Typing, typingIndicatorInterval)
	defer stopTyping()

	prompt := buildAssistantGeneratePrompt(preset, text)
	generated, err := bc.bootstrapGenerate(prompt)
	if err != nil {
		log.Printf("Bootstrap assistant LLM error: %v", err)
		bc.setPendingBootstrap(c.Sender().ID, bootstrapState{Choice: state.Choice, Step: bootstrapStepProfile})
		return SendContextText(c, bootstrapProfileMessage)
	}

	log.Printf("Bootstrap assistant LLM output (%d chars): %.500s", len(generated), generated)

	identity, soul, err := parseGeneratedPersona(generated)
	if err != nil {
		log.Printf("Bootstrap assistant parse error: %v — raw output: %.300s", err, generated)
		bc.setPendingBootstrap(c.Sender().ID, bootstrapState{Choice: state.Choice, Step: bootstrapStepProfile})
		return SendContextText(c, bootstrapProfileMessage)
	}

	if err := writeGeneratedPersona(bc.personasDir, identity, soul); err != nil {
		log.Printf("Bootstrap assistant write error: %v", err)
		bc.setPendingBootstrap(c.Sender().ID, bootstrapState{Choice: state.Choice, Step: bootstrapStepProfile})
		return SendContextText(c, bootstrapProfileMessage)
	}

	bc.setPendingBootstrap(c.Sender().ID, bootstrapState{Choice: state.Choice, Step: bootstrapStepProfile})
	return SendContextText(c, bootstrapProfileMessage)
}

func (bc *BotController) completeBootstrapProfile(c telebot.Context, state bootstrapState, text string) error {
	fallbackName := bootstrapFallbackName(c.Sender())

	stopTyping := startChatActionLoop(bc.bot, c.Chat(), telebot.Typing, typingIndicatorInterval)
	defer stopTyping()

	prompt := buildUserGeneratePrompt(text, fallbackName)
	generated, err := bc.bootstrapGenerate(prompt)
	if err != nil {
		log.Printf("Bootstrap profile LLM error: %v — using template fallback", err)
		generated = buildUserTemplateFromProfile(text, fallbackName)
	}

	if err := os.WriteFile(filepath.Join(bc.personasDir, "USER.md"), []byte(generated), 0o644); err != nil {
		log.Printf("Bootstrap user profile write error: %v\n", err)
		return SendContextText(c, bootstrapFailureMessage)
	}

	return SendContextText(c, bootstrapSuccessMessage)
}

// bootstrapGenerate calls the LLM via bridge to generate persona content.
// It accumulates content from assistant events since the terminal result
// event may not carry the full text.
func (bc *BotController) bootstrapGenerate(prompt string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), bootstrapGenerateTimeout)
	defer cancel()

	ch, err := bc.bridge.Execute(ctx, bridge.Request{
		Command: "query",
		Prompt:  prompt,
		Options: bridge.RequestOptions{
			Model:          bc.config.DefaultModel,
			SystemPrompt:   "Voce e um gerador de arquivos de configuracao de persona. Responda apenas com o conteudo solicitado, sem explicacoes adicionais.",
			MaxTurns:       1,
			PermissionMode: "bypassPermissions",
		},
	})
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	for ev := range ch {
		switch ev.Type {
		case "assistant":
			content := ev.Text
			if content == "" {
				content = ev.Content
			}
			buf.WriteString(content)
		case "result":
			content := ev.Text
			if content == "" {
				content = ev.Content
			}
			if content != "" {
				buf.Reset()
				buf.WriteString(content)
			}
			return buf.String(), nil
		case "error":
			msg := ev.Message
			if msg == "" {
				msg = ev.Content
			}
			return "", fmt.Errorf("bridge error: %s", msg)
		}
	}

	if buf.Len() > 0 {
		return buf.String(), nil
	}
	return "", fmt.Errorf("bridge: no content received")
}

