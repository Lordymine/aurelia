package telegram

import (
	"testing"
	"time"

	"github.com/kocar/aurelia/internal/agent"
	"gopkg.in/telebot.v3"
)

func TestIsSupportedImageDocument(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		filename string
		mimeType string
		want     bool
	}{
		{name: "mime image", filename: "scan.bin", mimeType: "image/png", want: true},
		{name: "extension fallback", filename: "photo.webp", mimeType: "", want: true},
		{name: "pdf is not image", filename: "report.pdf", mimeType: "application/pdf", want: false},
		{name: "markdown is not image", filename: "notes.md", mimeType: "text/markdown", want: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isSupportedImageDocument(tc.filename, tc.mimeType); got != tc.want {
				t.Fatalf("isSupportedImageDocument(%q, %q) = %t, want %t", tc.filename, tc.mimeType, got, tc.want)
			}
		})
	}
}

func TestDetectImageMIMEType(t *testing.T) {
	t.Parallel()

	if got := detectImageMIMEType("photo.png", ""); got != "image/png" {
		t.Fatalf("expected extension-based mime image/png, got %q", got)
	}
	if got := detectImageMIMEType("photo.bin", "image/webp"); got != "image/webp" {
		t.Fatalf("expected explicit mime to win, got %q", got)
	}
	if got := detectImageMIMEType("photo.unknown", ""); got != "image/jpeg" {
		t.Fatalf("expected jpeg fallback, got %q", got)
	}
}

func TestStoreAndFlushAlbumPhotos(t *testing.T) {
	t.Parallel()

	bc := &BotController{pendingAlbums: make(map[string]*pendingAlbum)}

	firstOwner := bc.storeAlbumPhoto("album-1", 12, "", telebot.Photo{File: telebot.File{FileID: "b"}})
	secondOwner := bc.storeAlbumPhoto("album-1", 10, "Legenda do album", telebot.Photo{File: telebot.File{FileID: "a"}})

	if !firstOwner {
		t.Fatal("expected first photo in album to become owner")
	}
	if secondOwner {
		t.Fatal("expected subsequent photo not to become owner")
	}

	caption, photos, ok := bc.flushAlbumPhotos("album-1")
	if !ok {
		t.Fatal("expected album flush to succeed")
	}
	if caption != "Legenda do album" {
		t.Fatalf("expected album caption to be preserved, got %q", caption)
	}
	if len(photos) != 2 {
		t.Fatalf("expected 2 photos, got %d", len(photos))
	}
	if photos[0].messageID != 10 || photos[1].messageID != 12 {
		t.Fatalf("expected photos sorted by message id, got %+v", photos)
	}
	if _, _, ok := bc.flushAlbumPhotos("album-1"); ok {
		t.Fatal("expected album to be removed after flush")
	}
}

func TestInputSessionPersistedContent_MultipleImages(t *testing.T) {
	t.Parallel()

	session := inputSession{
		text: "Compare estas referencias.",
		message: agent.Message{
			Role:    "user",
			Content: "Compare estas referencias.",
			Parts: []agent.ContentPart{
				{Type: agent.ContentPartText, Text: "Compare estas referencias."},
				{Type: agent.ContentPartImage, MIMEType: "image/jpeg", Data: []byte("a")},
				{Type: agent.ContentPartImage, MIMEType: "image/png", Data: []byte("b")},
			},
		},
	}

	if got := session.persistedContent(); got != "Compare estas referencias. [imagem anexada: 2]" {
		t.Fatalf("unexpected persisted content %q", got)
	}
}

func TestShouldReuseRecentMedia(t *testing.T) {
	t.Parallel()

	cases := map[string]bool{
		"verifique a imagem que mandei agora": true,
		"analise o pdf em anexo":              true,
		"me responda em uma frase":            false,
	}

	for input, want := range cases {
		if got := shouldReuseRecentMedia(input); got != want {
			t.Fatalf("shouldReuseRecentMedia(%q) = %t, want %t", input, got, want)
		}
	}
}

func TestStoreAndLoadRecentMedia(t *testing.T) {
	t.Parallel()

	bc := &BotController{recentMedia: make(map[string]recentMedia)}
	session := inputSession{
		convID: "42",
		message: agent.Message{
			Parts: []agent.ContentPart{
				{Type: agent.ContentPartText, Text: "Analise a imagem"},
				{Type: agent.ContentPartImage, MIMEType: "image/jpeg", Data: []byte("jpg")},
			},
		},
	}

	bc.storeRecentMedia(session)

	media, ok := bc.loadRecentMedia("42")
	if !ok {
		t.Fatal("expected recent media to be available")
	}
	if len(media.parts) != 1 || media.parts[0].Type != agent.ContentPartImage {
		t.Fatalf("unexpected media parts %+v", media.parts)
	}
}

func TestLoadRecentMedia_Expires(t *testing.T) {
	t.Parallel()

	bc := &BotController{
		recentMedia: map[string]recentMedia{
			"42": {
				parts:     []agent.ContentPart{{Type: agent.ContentPartImage, MIMEType: "image/jpeg", Data: []byte("jpg")}},
				updatedAt: time.Now().Add(-4 * time.Minute),
			},
		},
	}

	if _, ok := bc.loadRecentMedia("42"); ok {
		t.Fatal("expected stale media to expire")
	}
}

func TestAttachRecentMediaIfRelevant(t *testing.T) {
	t.Parallel()

	bc := &BotController{
		recentMedia: map[string]recentMedia{
			"42": {
				parts:     []agent.ContentPart{{Type: agent.ContentPartImage, MIMEType: "image/jpeg", Data: []byte("jpg")}},
				updatedAt: time.Now(),
			},
		},
	}
	ctx := fakeTelebotContext{senderID: 42}

	parts := bc.attachRecentMediaIfRelevant(ctx, "verifique a imagem correta", nil)
	if len(parts) != 2 {
		t.Fatalf("expected text + cached image, got %+v", parts)
	}
	if parts[0].Type != agent.ContentPartText || parts[1].Type != agent.ContentPartImage {
		t.Fatalf("unexpected attached parts %+v", parts)
	}
}

type fakeTelebotContext struct {
	telebot.Context
	senderID int64
}

func (f fakeTelebotContext) Sender() *telebot.User {
	return &telebot.User{ID: f.senderID}
}

func (f fakeTelebotContext) Chat() *telebot.Chat {
	return &telebot.Chat{ID: f.senderID}
}

func (f fakeTelebotContext) Message() *telebot.Message {
	return &telebot.Message{Sender: f.Sender(), Chat: f.Chat()}
}

func (f fakeTelebotContext) Recipient() telebot.Recipient {
	return fakeRecipient{id: f.senderID}
}
