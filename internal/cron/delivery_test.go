package cron

import (
	"context"
	"fmt"
	"testing"
)

type mockSender struct {
	lastChatID int64
	lastText   string
	err        error
}

func (m *mockSender) Send(chatID int64, text string) error {
	m.lastChatID = chatID
	m.lastText = text
	return m.err
}

func TestTelegramDelivery_Success(t *testing.T) {
	sender := &mockSender{}
	d := NewTelegramDelivery(sender)

	job := CronJob{ID: "12345678-abcd", TargetChatID: 42}
	result := &ExecutionResult{Output: "hello world"}

	err := d.Deliver(context.Background(), job, result, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sender.lastChatID != 42 {
		t.Fatalf("expected chat 42, got %d", sender.lastChatID)
	}
	if sender.lastText == "" {
		t.Fatal("expected non-empty text")
	}
}

func TestTelegramDelivery_Error(t *testing.T) {
	sender := &mockSender{}
	d := NewTelegramDelivery(sender)

	job := CronJob{ID: "12345678-abcd", TargetChatID: 42}
	execErr := fmt.Errorf("bridge failed")

	err := d.Deliver(context.Background(), job, nil, execErr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sender.lastChatID != 42 {
		t.Fatalf("expected chat 42, got %d", sender.lastChatID)
	}
}

func TestTelegramDelivery_NoChatID(t *testing.T) {
	sender := &mockSender{}
	d := NewTelegramDelivery(sender)

	job := CronJob{ID: "12345678-abcd", TargetChatID: 0}
	result := &ExecutionResult{Output: "hello"}

	err := d.Deliver(context.Background(), job, result, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sender.lastText != "" {
		t.Fatal("expected no send for zero chat ID")
	}
}

func TestTelegramDelivery_EmptyOutput(t *testing.T) {
	sender := &mockSender{}
	d := NewTelegramDelivery(sender)

	job := CronJob{ID: "12345678-abcd", TargetChatID: 42}
	result := &ExecutionResult{Output: ""}

	err := d.Deliver(context.Background(), job, result, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sender.lastText != "" {
		t.Fatal("expected no send for empty output")
	}
}
