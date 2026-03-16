package telegram

import (
	"errors"
	"testing"

	"gopkg.in/telebot.v3"
)

type stubContextSender struct {
	calls        []sendCall
	firstSendErr error
}

func (s *stubContextSender) Send(what interface{}, opts ...interface{}) error {
	s.calls = append(s.calls, sendCall{what: what, opts: opts})
	if len(s.calls) == 1 && s.firstSendErr != nil {
		return s.firstSendErr
	}
	return nil
}

func TestSendContextText_SendsTelegramHTML(t *testing.T) {
	sender := &stubContextSender{}

	if err := SendContextText(sender, "## Titulo\n\n- **item**"); err != nil {
		t.Fatalf("SendContextText returned error: %v", err)
	}

	if len(sender.calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(sender.calls))
	}

	payload, ok := sender.calls[0].what.(string)
	if !ok {
		t.Fatalf("expected string payload, got %T", sender.calls[0].what)
	}
	if !containsSubstring(payload, "<b>Titulo</b>") {
		t.Fatalf("expected html payload, got %s", payload)
	}

	options, ok := sender.calls[0].opts[0].(*telebot.SendOptions)
	if !ok {
		t.Fatalf("expected send options, got %T", sender.calls[0].opts[0])
	}
	if options.ParseMode != telebot.ModeHTML {
		t.Fatalf("expected html parse mode, got %q", options.ParseMode)
	}
}

func TestSendContextText_FallsBackToPlainText(t *testing.T) {
	sender := &stubContextSender{firstSendErr: errors.New("bad html")}

	if err := SendContextText(sender, "## Titulo"); err != nil {
		t.Fatalf("SendContextText returned error: %v", err)
	}

	if len(sender.calls) != 2 {
		t.Fatalf("expected 2 calls, got %d", len(sender.calls))
	}

	payload, ok := sender.calls[1].what.(string)
	if !ok {
		t.Fatalf("expected plain text payload, got %T", sender.calls[1].what)
	}
	if payload != "## Titulo" {
		t.Fatalf("expected plain text fallback, got %q", payload)
	}
}
