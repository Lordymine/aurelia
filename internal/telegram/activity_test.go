package telegram

import (
	"sync"
	"testing"
	"time"

	"gopkg.in/telebot.v3"
)

type fakeActionSender struct {
	mu      sync.Mutex
	calls   int
	actions []telebot.ChatAction
}

func (f *fakeActionSender) Notify(to telebot.Recipient, action telebot.ChatAction, until ...int) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls++
	f.actions = append(f.actions, action)
	return nil
}

type fakeRecipient struct {
	id int64
}

func (f fakeRecipient) Recipient() string {
	return ""
}

func TestChatActionLoop_SendsRepeatedNotificationsUntilStopped(t *testing.T) {
	t.Parallel()

	sender := &fakeActionSender{}
	recipient := fakeRecipient{id: 1}

	stop := startChatActionLoop(sender, recipient, telebot.Typing, 15*time.Millisecond)
	time.Sleep(55 * time.Millisecond)
	stop()

	sender.mu.Lock()
	callsAfterStop := sender.calls
	sender.mu.Unlock()

	if callsAfterStop < 2 {
		t.Fatalf("expected repeated notifications, got %d", callsAfterStop)
	}

	time.Sleep(40 * time.Millisecond)

	sender.mu.Lock()
	defer sender.mu.Unlock()
	if sender.calls != callsAfterStop {
		t.Fatalf("expected notifications to stop after stop(), before=%d after=%d", callsAfterStop, sender.calls)
	}
}
