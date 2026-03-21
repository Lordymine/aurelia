package telegram

import (
	"fmt"
	"log"
	"strings"

	"gopkg.in/telebot.v3"
)

func (bc *BotController) whitelistMiddleware() telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			sender := c.Sender()
			if sender == nil {
				return nil
			}
			if !bc.isAllowedUser(sender.ID) {
				log.Printf("blocked unauthorized user: %d\n", sender.ID)
				return nil
			}
			return next(c)
		}
	}
}

func (bc *BotController) registerContentRoutes() {
	bc.bot.Handle("/cwd", bc.handleCwdCommand)
	bc.bot.Handle("/reset", bc.handleResetCommand)
	bc.bot.Handle(telebot.OnText, bc.handleText)
	bc.bot.Handle(telebot.OnPhoto, bc.handlePhoto)
	bc.bot.Handle(telebot.OnDocument, bc.handleDocument)
	bc.bot.Handle(telebot.OnVoice, bc.handleVoice)
	bc.bot.Handle(telebot.OnAudio, bc.handleVoice)
}

func (bc *BotController) handleCwdCommand(c telebot.Context) error {
	args := strings.TrimSpace(c.Message().Payload)
	if args == "" {
		cwd := bc.sessions.GetCwd(c.Chat().ID)
		if cwd == "" {
			return SendText(bc.bot, c.Chat(), "Nenhum diretório configurado. Use: /cwd C:\\path\\to\\project")
		}
		return SendText(bc.bot, c.Chat(), fmt.Sprintf("Diretório atual: %s", cwd))
	}
	bc.sessions.SetCwd(c.Chat().ID, args)
	return SendText(bc.bot, c.Chat(), fmt.Sprintf("Diretório configurado: %s", args))
}

func (bc *BotController) handleResetCommand(c telebot.Context) error {
	bc.sessions.Clear(c.Chat().ID)
	return SendText(bc.bot, c.Chat(), "Sessão resetada. Próxima mensagem inicia conversa nova.")
}
