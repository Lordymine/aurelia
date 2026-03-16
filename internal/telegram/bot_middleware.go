package telegram

import (
	"log"

	"gopkg.in/telebot.v3"
)

func (bc *BotController) whitelistMiddleware() telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			senderID := c.Sender().ID
			if !bc.isAllowedUser(senderID) {
				log.Printf("blocked unauthorized user: %d\n", senderID)
				return nil
			}
			return next(c)
		}
	}
}

func (bc *BotController) registerContentRoutes() {
	bc.bot.Handle(telebot.OnText, bc.handleText)
	bc.bot.Handle(telebot.OnDocument, bc.handleDocument)
	bc.bot.Handle(telebot.OnVoice, bc.handleVoice)
	bc.bot.Handle(telebot.OnAudio, bc.handleVoice)
}
