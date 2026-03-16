package stt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Transcriber is the interface for STT engines
type Transcriber interface {
	Transcribe(ctx context.Context, audioFilePath string) (string, error)
	IsAvailable() bool
}

// GroqTranscriber implements Transcriber using the Groq Whisper API
type GroqTranscriber struct {
	apiKey     string
	apiBase    string
	httpClient *http.Client
}

// NewGroqTranscriber creates a new Groq STT client
func NewGroqTranscriber(apiKey string) *GroqTranscriber {
	return &GroqTranscriber{
		apiKey:  apiKey,
		apiBase: "https://api.groq.com/openai/v1",
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// Transcribe converts local audio file to text
func (t *GroqTranscriber) Transcribe(ctx context.Context, audioFilePath string) (string, error) {
	log.Printf("Starting Groq transcription for: %s\n", audioFilePath)

	audioFile, err := os.Open(audioFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to open audio file: %w", err)
	}
	defer func() { _ = audioFile.Close() }()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("file", filepath.Base(audioFilePath))
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}

	_, err = io.Copy(part, audioFile)
	if err != nil {
		return "", fmt.Errorf("failed to copy file content: %w", err)
	}

	if err := writer.WriteField("model", "whisper-large-v3"); err != nil {
		return "", fmt.Errorf("failed to write model field: %w", err)
	}

	if err := writer.WriteField("response_format", "json"); err != nil {
		return "", fmt.Errorf("failed to write response_format field: %w", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close multipart writer: %w", err)
	}

	url := t.apiBase + "/audio/transcriptions"
	req, err := http.NewRequestWithContext(ctx, "POST", url, &requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+t.apiKey)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	log.Printf("Groq transcription completed successfully\n")

	return result.Text, nil
}

// IsAvailable returns true if the transcriber has an API key configured
func (t *GroqTranscriber) IsAvailable() bool {
	return t.apiKey != ""
}
