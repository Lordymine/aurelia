package stt

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGroqTranscriberIsAvailable(t *testing.T) {
	if NewGroqTranscriber("").IsAvailable() {
		t.Fatal("expected transcriber without api key to be unavailable")
	}
	if !NewGroqTranscriber("secret").IsAvailable() {
		t.Fatal("expected transcriber with api key to be available")
	}
}

func TestGroqTranscriberTranscribeSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer secret" {
			t.Fatalf("unexpected authorization header: %q", got)
		}
		if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			t.Fatalf("unexpected content type: %q", r.Header.Get("Content-Type"))
		}
		_, _ = w.Write([]byte(`{"text":"hello world","language":"en","duration":1.5}`))
	}))
	defer server.Close()

	audioPath := filepath.Join(t.TempDir(), "sample.wav")
	if err := os.WriteFile(audioPath, []byte("audio"), 0o600); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	transcriber := NewGroqTranscriber("secret")
	transcriber.apiBase = server.URL

	result, err := transcriber.Transcribe(context.Background(), audioPath)
	if err != nil {
		t.Fatalf("Transcribe returned error: %v", err)
	}
	if result != "hello world" {
		t.Fatalf("unexpected transcription result: %s", result)
	}
}

func TestGroqTranscriberTranscribeAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad request", http.StatusBadRequest)
	}))
	defer server.Close()

	audioPath := filepath.Join(t.TempDir(), "sample.wav")
	if err := os.WriteFile(audioPath, []byte("audio"), 0o600); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	transcriber := NewGroqTranscriber("secret")
	transcriber.apiBase = server.URL

	_, err := transcriber.Transcribe(context.Background(), audioPath)
	if err == nil {
		t.Fatal("expected api error")
	}
}
