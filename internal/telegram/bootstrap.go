package telegram

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/telebot.v3"
)

func (bc *BotController) setupBootstrapRoutes() {
	bc.bot.Handle("/start", bc.handleStart)
	bc.bot.Handle("\fbtn_coder", bc.handleBootstrapChoice("coder"))
	bc.bot.Handle("\fbtn_assist", bc.handleBootstrapChoice("assist"))
}

func (bc *BotController) handleStart(c telebot.Context) error {
	identityExists := bootstrapIdentityExists(bc.personasDir)

	message, menu := bootstrapStartResponse(identityExists)
	if menu == nil {
		return SendContextText(c, message)
	}
	return SendContextText(c, message, menu)
}

func (bc *BotController) handleBootstrapChoice(choice string) func(telebot.Context) error {
	return func(c telebot.Context) error {
		_ = bc.bot.Respond(c.Callback(), &telebot.CallbackResponse{})

		// Write base preset as fallback (will be overwritten by LLM generation)
		preset, err := bootstrapPresetForChoice(choice)
		if err != nil {
			return SendContextText(c, bootstrapFailureMessage)
		}
		if err := writeBootstrapPreset(bc.personasDir, preset); err != nil {
			log.Printf("Bootstrap error: %v\n", err)
			return SendContextText(c, bootstrapFailureMessage)
		}

		bc.setPendingBootstrap(c.Sender().ID, bootstrapState{Choice: choice, Step: bootstrapStepAssistant})
		return SendContextText(c, bootstrapAssistantMessage)
	}
}

func bootstrapIdentityExists(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, "IDENTITY.md"))
	return err == nil
}
