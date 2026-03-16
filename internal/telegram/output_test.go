package telegram

import (
	"errors"
	"testing"

	"gopkg.in/telebot.v3"
)

type sendCall struct {
	to   telebot.Recipient
	what interface{}
	opts []interface{}
}

type stubSender struct {
	calls        []sendCall
	firstSendErr error
}

func (s *stubSender) Send(to telebot.Recipient, what interface{}, opts ...interface{}) (*telebot.Message, error) {
	s.calls = append(s.calls, sendCall{to: to, what: what, opts: opts})
	if len(s.calls) == 1 && s.firstSendErr != nil {
		return nil, s.firstSendErr
	}
	return &telebot.Message{}, nil
}

func TestSendText_SendsTelegramHTML(t *testing.T) {
	sender := &stubSender{}
	chat := &telebot.Chat{ID: 123}

	if err := sendTextWithSender(sender, chat, "## Title\n\n- **item**", 200); err != nil {
		t.Fatalf("sendTextWithSender returned error: %v", err)
	}

	if len(sender.calls) != 1 {
		t.Fatalf("expected 1 send call, got %d", len(sender.calls))
	}

	text, ok := sender.calls[0].what.(string)
	if !ok {
		t.Fatalf("expected sent payload to be string, got %T", sender.calls[0].what)
	}
	if !containsSubstring(text, "<b>Title</b>") {
		t.Fatalf("expected html formatted text, got: %s", text)
	}

	options, ok := sender.calls[0].opts[0].(*telebot.SendOptions)
	if !ok {
		t.Fatalf("expected first option to be *telebot.SendOptions, got %T", sender.calls[0].opts[0])
	}
	if options.ParseMode != telebot.ModeHTML {
		t.Fatalf("expected parse mode %q, got %q", telebot.ModeHTML, options.ParseMode)
	}
}

func TestSendText_FallsBackToPlainTextWhenHTMLSendFails(t *testing.T) {
	sender := &stubSender{firstSendErr: errors.New("bad html")}
	chat := &telebot.Chat{ID: 123}

	if err := sendTextWithSender(sender, chat, "## Title", 200); err != nil {
		t.Fatalf("sendTextWithSender returned error: %v", err)
	}

	if len(sender.calls) != 2 {
		t.Fatalf("expected 2 send calls, got %d", len(sender.calls))
	}

	if _, ok := sender.calls[1].what.(string); !ok {
		t.Fatalf("expected plain text fallback payload, got %T", sender.calls[1].what)
	}
	if len(sender.calls[1].opts) != 0 {
		t.Fatalf("expected fallback send without options, got %d opts", len(sender.calls[1].opts))
	}
}

func TestSplitTelegramMarkdown_PrefersParagraphBoundaries(t *testing.T) {
	text := "primeiro bloco\n\nsegundo bloco muito maior para obrigar split"

	chunks := splitTelegramMarkdown(text, 35)
	if len(chunks) < 2 {
		t.Fatalf("expected multiple chunks, got %d", len(chunks))
	}
	if chunks[0] != "primeiro bloco" {
		t.Fatalf("expected first chunk to stop at paragraph boundary, got %q", chunks[0])
	}
	for _, chunk := range chunks {
		if len([]rune(chunk)) > 35 {
			t.Fatalf("chunk exceeded limit: %q", chunk)
		}
	}
}

func TestSendError_SendsFormattedHTML(t *testing.T) {
	sender := &stubSender{}
	chat := &telebot.Chat{ID: 123}

	if err := sendErrorWithSender(sender, chat, "Erro", "max iterations reached"); err != nil {
		t.Fatalf("sendErrorWithSender returned error: %v", err)
	}

	if len(sender.calls) != 1 {
		t.Fatalf("expected 1 send call, got %d", len(sender.calls))
	}

	text, ok := sender.calls[0].what.(string)
	if !ok {
		t.Fatalf("expected sent payload to be string, got %T", sender.calls[0].what)
	}
	if !containsSubstring(text, "<b>Erro</b>") {
		t.Fatalf("expected bold html title, got: %s", text)
	}
	if !containsSubstring(text, "max iterations reached") {
		t.Fatalf("expected error body in payload, got: %s", text)
	}

	options, ok := sender.calls[0].opts[0].(*telebot.SendOptions)
	if !ok {
		t.Fatalf("expected first option to be *telebot.SendOptions, got %T", sender.calls[0].opts[0])
	}
	if options.ParseMode != telebot.ModeHTML {
		t.Fatalf("expected parse mode %q, got %q", telebot.ModeHTML, options.ParseMode)
	}
}

func TestSendError_FallsBackToPlainTextWhenHTMLSendFails(t *testing.T) {
	sender := &stubSender{firstSendErr: errors.New("bad html")}
	chat := &telebot.Chat{ID: 123}

	if err := sendErrorWithSender(sender, chat, "Erro", "max iterations reached"); err != nil {
		t.Fatalf("sendErrorWithSender returned error: %v", err)
	}

	if len(sender.calls) != 2 {
		t.Fatalf("expected 2 send calls, got %d", len(sender.calls))
	}

	payload, ok := sender.calls[1].what.(string)
	if !ok {
		t.Fatalf("expected plain text fallback payload, got %T", sender.calls[1].what)
	}
	if payload != "Erro\n\nmax iterations reached" {
		t.Fatalf("unexpected fallback payload: %q", payload)
	}
}
