package cron

import (
	"context"
	"fmt"
	"log"
)

// ChatSender sends text messages to a chat by ID.
type ChatSender interface {
	Send(chatID int64, text string) error
}

// TelegramDelivery sends cron execution results to Telegram chats.
type TelegramDelivery struct {
	sender ChatSender
}

// NewTelegramDelivery creates a delivery that sends via the given sender.
func NewTelegramDelivery(sender ChatSender) *TelegramDelivery {
	return &TelegramDelivery{sender: sender}
}

// Deliver sends the cron job result or error to the target chat.
func (d *TelegramDelivery) Deliver(ctx context.Context, job CronJob, result *ExecutionResult, execErr error) error {
	if job.TargetChatID == 0 {
		log.Println("Cron delivery skipped: no chat ID")
		return nil
	}

	output := ""
	if result != nil {
		output = result.Output
	}
	log.Printf("Cron delivery: job=%s chat=%d output_len=%d err=%v", job.ID[:8], job.TargetChatID, len(output), execErr)

	if execErr != nil {
		return d.sender.Send(job.TargetChatID, fmt.Sprintf("❌ Cron job %s falhou: %v", job.ID[:8], execErr))
	}
	if output == "" {
		return nil
	}
	header := fmt.Sprintf("📋 Resultado agendamento (%s):\n\n", job.ID[:8])
	return d.sender.Send(job.TargetChatID, header+output)
}
