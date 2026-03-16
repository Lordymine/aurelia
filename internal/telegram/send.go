package telegram

import "gopkg.in/telebot.v3"

type contextSender interface {
	Send(what interface{}, opts ...interface{}) error
}

func SendContextText(c contextSender, text string, opts ...interface{}) error {
	html := MarkdownToHTML(text)
	sendOpts := append([]interface{}{&telebot.SendOptions{ParseMode: telebot.ModeHTML}}, opts...)
	if err := c.Send(html, sendOpts...); err == nil {
		return nil
	}
	return c.Send(text, opts...)
}
