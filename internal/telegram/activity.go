package telegram

import (
	"sync"
	"time"

	"gopkg.in/telebot.v3"
)

type actionSender interface {
	Notify(to telebot.Recipient, action telebot.ChatAction, until ...int) error
}

func startChatActionLoop(sender actionSender, recipient telebot.Recipient, action telebot.ChatAction, interval time.Duration) func() {
	if sender == nil || recipient == nil {
		return func() {}
	}
	if interval <= 0 {
		interval = typingIndicatorInterval
	}

	done := make(chan struct{})
	var once sync.Once

	go func() {
		_ = sender.Notify(recipient, action)

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				_ = sender.Notify(recipient, action)
			case <-done:
				return
			}
		}
	}()

	return func() {
		once.Do(func() {
			close(done)
		})
	}
}
