package telegram

import (
	"gopkg.in/telebot.v3"
)

func (bc *BotController) processInput(c telebot.Context, text string, parts [][]byte, requiresAudio bool) error {
	_ = parts

	if state, ok := bc.popPendingBootstrap(c.Sender().ID); ok {
		return bc.completeBootstrapProfile(c, state, text)
	}

	// TODO: wire bridge executor to process user input
	return SendText(bc.bot, c.Chat(), "Aurelia is not yet wired to an executor. Bridge integration pending.")
}
