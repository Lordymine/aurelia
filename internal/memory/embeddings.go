package memory

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"math"
	"net/http"
	"time"
)

// Embedder generates vector embeddings from text.
type Embedder interface {
	Embed(ctx context.Context, text string) ([]float32, error)
	Dimensions() int
}

// VoyageEmbedder calls Voyage AI's /v1/embeddings endpoint.
type VoyageEmbedder struct {
	apiKey string
	model  string
	dims   int
	client *http.Client
}

// NewVoyageEmbedder creates an embedder that calls Voyage AI.
func NewVoyageEmbedder(apiKey, model string) *VoyageEmbedder {
	return &VoyageEmbedder{
		apiKey: apiKey,
		model:  model,
		dims:   1024,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

type voyageRequest struct {
	Input []string `json:"input"`
	Model string   `json:"model"`
}

type voyageResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
}

func (e *VoyageEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
	body, err := json.Marshal(voyageRequest{
		Input: []string{text},
		Model: e.model,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.voyageai.com/v1/embeddings", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+e.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("voyage API returned status %d", resp.StatusCode)
	}

	var result voyageResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if len(result.Data) == 0 {
		return nil, fmt.Errorf("voyage API returned no embeddings")
	}

	return result.Data[0].Embedding, nil
}

func (e *VoyageEmbedder) Dimensions() int {
	return e.dims
}

// MockEmbedder produces deterministic vectors for testing.
// Same text always produces the same vector; different texts produce different vectors.
type MockEmbedder struct {
	dims int
}

// NewMockEmbedder creates a mock embedder with the given dimensions.
func NewMockEmbedder(dims int) *MockEmbedder {
	return &MockEmbedder{dims: dims}
}

func (e *MockEmbedder) Embed(_ context.Context, text string) ([]float32, error) {
	vec := make([]float32, e.dims)

	// Use FNV hash to seed deterministic values per word.
	// We hash each word and spread its influence across dimensions
	// so that texts with shared words produce similar vectors.
	words := splitWords(text)
	for _, word := range words {
		h := fnv.New64a()
		h.Write([]byte(word))
		seed := h.Sum64()
		for i := range vec {
			// Deterministic pseudo-random per (word, dimension)
			v := seed ^ uint64(i)*2654435761
			vec[i] += float32(int64(v%1000)-500) / 500.0
		}
	}

	// Normalize to unit vector
	var norm float64
	for _, v := range vec {
		norm += float64(v) * float64(v)
	}
	norm = math.Sqrt(norm)
	if norm > 0 {
		for i := range vec {
			vec[i] = float32(float64(vec[i]) / norm)
		}
	}

	return vec, nil
}

func (e *MockEmbedder) Dimensions() int {
	return e.dims
}

// splitWords splits text into lowercase words.
func splitWords(text string) []string {
	var words []string
	var current []byte
	for i := 0; i < len(text); i++ {
		c := text[i]
		if c >= 'A' && c <= 'Z' {
			current = append(current, c+32) // tolower
		} else if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') {
			current = append(current, c)
		} else {
			if len(current) > 0 {
				words = append(words, string(current))
				current = current[:0]
			}
		}
	}
	if len(current) > 0 {
		words = append(words, string(current))
	}
	return words
}
