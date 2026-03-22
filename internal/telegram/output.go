package telegram

import (
	"log"
	"strings"
	"time"

	"gopkg.in/telebot.v3"
)

const telegramMessageLimit = 3900
const interChunkDelay = 200 * time.Millisecond

type messageSender interface {
	Send(to telebot.Recipient, what interface{}, opts ...interface{}) (*telebot.Message, error)
}

func SendText(bot *telebot.Bot, chat *telebot.Chat, text string) error {
	return sendTextWithSender(bot, chat, text, telegramMessageLimit)
}

func sendTextWithSender(sender messageSender, chat *telebot.Chat, text string, limit int) error {
	chunks := splitTelegramMarkdown(text, limit)
	for _, chunk := range chunks {
		htmlChunk := MarkdownToHTML(chunk)
		_, err := sender.Send(chat, htmlChunk, &telebot.SendOptions{
			ParseMode: telebot.ModeHTML,
		})
		if err == nil {
			time.Sleep(interChunkDelay)
			continue
		}

		log.Printf("Send chunk with HTML failed (%v). Retrying as plain text...", err)
		_, err = sender.Send(chat, chunk)
		if err != nil {
			if floodErr, ok := err.(*telebot.FloodError); ok {
				log.Printf("Hit rate limit in chunk sending. Retrying in %v...", floodErr.RetryAfter)
				time.Sleep(time.Duration(floodErr.RetryAfter) * time.Second)
				if _, retryErr := sender.Send(chat, chunk); retryErr == nil {
					time.Sleep(interChunkDelay)
					continue
				}
			}
			return err
		}
		time.Sleep(interChunkDelay)
	}
	return nil
}

func splitTelegramMarkdown(text string, limit int) []string {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return []string{""}
	}

	var chunks []string
	remaining := trimmed
	for len([]rune(remaining)) > limit {
		splitAt := bestSplitIndex(remaining, limit)
		chunks = append(chunks, strings.TrimSpace(remaining[:splitAt]))
		remaining = strings.TrimSpace(remaining[splitAt:])
	}
	if remaining != "" {
		chunks = append(chunks, remaining)
	}
	return chunks
}

func bestSplitIndex(text string, limit int) int {
	runes := []rune(text)
	if len(runes) <= limit {
		return len(text)
	}

	candidates := []string{"\n\n", "\n", ". ", " "}
	window := string(runes[:limit])
	for _, candidate := range candidates {
		if idx := strings.LastIndex(window, candidate); idx > 0 {
			return idx
		}
	}
	return len(string(runes[:limit]))
}

func SendTextReply(bot *telebot.Bot, chat *telebot.Chat, text string, replyToID int) error {
	if replyToID == 0 {
		return SendText(bot, chat, text)
	}
	return sendTextReplyWithSender(bot, chat, text, telegramMessageLimit, replyToID)
}

func sendTextReplyWithSender(sender messageSender, chat *telebot.Chat, text string, limit int, replyToID int) error {
	chunks := splitTelegramMarkdown(text, limit)
	replyTo := &telebot.Message{ID: replyToID}

	for i, chunk := range chunks {
		htmlChunk := MarkdownToHTML(chunk)
		opts := &telebot.SendOptions{ParseMode: telebot.ModeHTML}
		// Only reply-to on the first chunk
		if i == 0 {
			opts.ReplyTo = replyTo
		}

		_, err := sender.Send(chat, htmlChunk, opts)
		if err == nil {
			time.Sleep(interChunkDelay)
			continue
		}

		log.Printf("Send chunk with HTML failed (%v). Retrying as plain text...", err)
		opts = &telebot.SendOptions{}
		if i == 0 {
			opts.ReplyTo = replyTo
		}
		_, err = sender.Send(chat, chunk, opts)
		if err != nil {
			return err
		}
		time.Sleep(interChunkDelay)
	}
	return nil
}

func ReactToMessage(bot *telebot.Bot, chat *telebot.Chat, messageID int, emoji string) {
	if messageID == 0 || chat == nil {
		return
	}
	msg := &telebot.Message{ID: messageID, Chat: chat}
	err := bot.React(chat, msg, telebot.ReactionOptions{
		Reactions: []telebot.Reaction{{Type: "emoji", Emoji: emoji}},
	})
	if err != nil {
		log.Printf("React error: %v", err)
	}
}

func SendError(bot *telebot.Bot, chat *telebot.Chat, errMsg string) error {
	return sendErrorWithSender(bot, chat, "Erro", errMsg)
}

func sendErrorWithSender(sender messageSender, chat *telebot.Chat, title, errMsg string) error {
	formatted := ErrorMessage(title, errMsg)
	_, err := sender.Send(chat, formatted, &telebot.SendOptions{
		ParseMode: telebot.ModeHTML,
	})
	if err == nil {
		return nil
	}

	log.Printf("Send error with HTML failed (%v). Retrying as plain text...", err)
	_, err = sender.Send(chat, title+"\n\n"+errMsg)
	return err
}
