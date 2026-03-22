package telegram

import (
	"context"
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
	bc.bot.Handle("/start", bc.handleStartCommand)
	bc.bot.Handle("/help", bc.handleStartCommand)
	bc.bot.Handle("/cwd", bc.handleCwdCommand)
	bc.bot.Handle("/reset", bc.handleResetCommand)
	bc.bot.Handle("/cron", bc.handleCronCommand)
	bc.bot.Handle("/agents", bc.handleAgentsCommand)
	bc.bot.Handle(telebot.OnText, bc.handleText)
	bc.bot.Handle(telebot.OnPhoto, bc.handlePhoto)
	bc.bot.Handle(telebot.OnDocument, bc.handleDocument)
	bc.bot.Handle(telebot.OnVoice, bc.handleVoice)
	bc.bot.Handle(telebot.OnAudio, bc.handleVoice)
}

func (bc *BotController) registerSlashMenu() {
	commands := []telebot.Command{
		{Text: "cwd", Description: "Definir diretório de trabalho"},
		{Text: "reset", Description: "Resetar sessão (conversa nova)"},
		{Text: "cron", Description: "Gerenciar agendamentos"},
		{Text: "agents", Description: "Listar agentes disponíveis"},
		{Text: "help", Description: "Mostrar comandos disponíveis"},
	}
	if err := bc.bot.SetCommands(commands); err != nil {
		log.Printf("Failed to set bot commands: %v", err)
	}
}

func (bc *BotController) handleStartCommand(c telebot.Context) error {
	help := "Comandos disponíveis:\n\n" +
		"/cwd <path> — Definir diretório de trabalho\n" +
		"/reset — Resetar sessão (conversa nova)\n" +
		"/cron — Gerenciar agendamentos\n" +
		"/agents — Listar agentes disponíveis\n" +
		"/help — Mostrar esta mensagem\n\n" +
		"Ou simplesmente envie uma mensagem e eu respondo."
	return SendText(bc.bot, c.Chat(), help)
}

func (bc *BotController) handleAgentsCommand(c telebot.Context) error {
	if bc.agents == nil || len(bc.agents.Agents()) == 0 {
		return SendText(bc.bot, c.Chat(), "Nenhum agente configurado. Crie arquivos .md em ~/.aurelia/agents/")
	}
	var lines []string
	for _, a := range bc.agents.Agents() {
		line := fmt.Sprintf("• %s — %s", a.Name, a.Description)
		if a.Model != "" {
			line += fmt.Sprintf(" [%s]", a.Model)
		}
		if a.Schedule != "" {
			line += fmt.Sprintf(" (cron: %s)", a.Schedule)
		}
		lines = append(lines, line)
	}
	return SendText(bc.bot, c.Chat(), "Agentes disponíveis:\n\n"+strings.Join(lines, "\n"))
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

func (bc *BotController) handleCronCommand(c telebot.Context) error {
	if bc.cronHandler == nil {
		return SendText(bc.bot, c.Chat(), "Cron não está disponível.")
	}
	userID := fmt.Sprintf("%d", c.Sender().ID)
	chatID := c.Chat().ID
	text := c.Message().Text

	reply, err := bc.cronHandler.HandleText(context.Background(), userID, chatID, text)
	if err != nil {
		return SendError(bc.bot, c.Chat(), err.Error())
	}
	if reply != "" {
		return SendText(bc.bot, c.Chat(), reply)
	}
	return nil
}
