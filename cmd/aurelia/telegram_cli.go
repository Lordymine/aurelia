package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/runtime"
	"gopkg.in/telebot.v3"
)

func runTelegramCLI(args []string) error {
	if len(args) == 0 {
		return printTelegramUsage()
	}

	resolver, err := runtime.New()
	if err != nil {
		return fmt.Errorf("resolve instance root: %w", err)
	}

	cfg, err := config.Load(resolver)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	bot, err := telebot.NewBot(telebot.Settings{
		Token:  cfg.TelegramBotToken,
		Poller: &telebot.LongPoller{Timeout: 1 * time.Second},
	})
	if err != nil {
		return fmt.Errorf("create bot: %w", err)
	}

	switch args[0] {
	case "react":
		// aurelia telegram react <chat-id> <message-id> <emoji>
		if len(args) < 4 {
			return fmt.Errorf("usage: aurelia telegram react <chat-id> <message-id> <emoji>")
		}
		chatID, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid chat-id: %w", err)
		}
		msgID, err := strconv.Atoi(args[2])
		if err != nil {
			return fmt.Errorf("invalid message-id: %w", err)
		}
		emoji := args[3]

		chat := &telebot.Chat{ID: chatID}
		msg := &telebot.Message{ID: msgID, Chat: chat}
		return bot.React(chat, msg, telebot.ReactionOptions{
			Reactions: []telebot.Reaction{{Type: "emoji", Emoji: emoji}},
		})

	case "send":
		// aurelia telegram send <chat-id> <text>
		if len(args) < 3 {
			return fmt.Errorf("usage: aurelia telegram send <chat-id> <text>")
		}
		chatID, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid chat-id: %w", err)
		}
		text := args[2]
		chat := &telebot.Chat{ID: chatID}
		_, err = bot.Send(chat, text)
		return err

	case "reply":
		// aurelia telegram reply <chat-id> <message-id> <text>
		if len(args) < 4 {
			return fmt.Errorf("usage: aurelia telegram reply <chat-id> <message-id> <text>")
		}
		chatID, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid chat-id: %w", err)
		}
		msgID, err := strconv.Atoi(args[2])
		if err != nil {
			return fmt.Errorf("invalid message-id: %w", err)
		}
		text := args[3]
		chat := &telebot.Chat{ID: chatID}
		_, err = bot.Send(chat, text, &telebot.SendOptions{
			ReplyTo: &telebot.Message{ID: msgID},
		})
		return err

	default:
		return printTelegramUsage()
	}
}

func printTelegramUsage() error {
	fmt.Println("Usage:")
	fmt.Println("  aurelia telegram react <chat-id> <message-id> <emoji>")
	fmt.Println("  aurelia telegram send <chat-id> <text>")
	fmt.Println("  aurelia telegram reply <chat-id> <message-id> <text>")
	return nil
}
