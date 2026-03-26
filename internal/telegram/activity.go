package telegram

import (
	"log"
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
		if err := sender.Notify(recipient, action); err != nil {
			log.Printf("Failed to send typing indicator: %v", err)
		}

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := sender.Notify(recipient, action); err != nil {
				log.Printf("Failed to send typing indicator: %v", err)
			}
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
