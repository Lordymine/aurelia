package memory

import (
	"context"
	"hash/fnv"
	"math"
)

// Embedder generates vector embeddings from text.
type Embedder interface {
	Embed(ctx context.Context, text string) ([]float32, error)
	Dimensions() int
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

	words := splitWords(text)
	for _, word := range words {
		h := fnv.New64a()
		h.Write([]byte(word))
		seed := h.Sum64()
		for i := range vec {
			v := seed ^ uint64(i)*2654435761
			vec[i] += float32(int64(v%1000)-500) / 500.0
		}
	}

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
			current = append(current, c+32)
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
