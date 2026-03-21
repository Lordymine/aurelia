package memory

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// GeminiEmbedder calls Google's Gemini embedding API.
type GeminiEmbedder struct {
	apiKey string
	model  string
	dims   int
	client *http.Client
}

// NewGeminiEmbedder creates an embedder that uses Google Gemini's embedding API.
// Model should be "gemini-embedding-001" or "text-embedding-004".
func NewGeminiEmbedder(apiKey, model string) *GeminiEmbedder {
	if model == "" {
		model = "gemini-embedding-001"
	}
	return &GeminiEmbedder{
		apiKey: apiKey,
		model:  model,
		dims:   768, // gemini-embedding-001 default
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

type geminiEmbedRequest struct {
	Model   string            `json:"model"`
	Content geminiContentPart `json:"content"`
}

type geminiContentPart struct {
	Parts []geminiTextPart `json:"parts"`
}

type geminiTextPart struct {
	Text string `json:"text"`
}

type geminiEmbedResponse struct {
	Embedding struct {
		Values []float32 `json:"values"`
	} `json:"embedding"`
}

func (e *GeminiEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:embedContent?key=%s", e.model, e.apiKey)

	body, err := json.Marshal(geminiEmbedRequest{
		Model: "models/" + e.model,
		Content: geminiContentPart{
			Parts: []geminiTextPart{{Text: text}},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gemini API returned status %d", resp.StatusCode)
	}

	var result geminiEmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if len(result.Embedding.Values) == 0 {
		return nil, fmt.Errorf("gemini API returned no embeddings")
	}

	e.dims = len(result.Embedding.Values)
	return result.Embedding.Values, nil
}

func (e *GeminiEmbedder) Dimensions() int {
	return e.dims
}
