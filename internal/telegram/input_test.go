package telegram

import (
	"testing"

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

// Tests for inputSession, recentMedia, and attachRecentMediaIfRelevant were removed
// because they depend on agent.Message and agent.ContentPart which no longer exist.
// They will be rewritten when the bridge executor is wired.
