package memory

import (
	"context"
	"math"
	"testing"
)

func TestMockEmbedder_Deterministic(t *testing.T) {
	t.Parallel()

	emb := NewMockEmbedder(64)
	ctx := context.Background()

	v1, err := emb.Embed(ctx, "hello world")
	if err != nil {
		t.Fatalf("Embed() error = %v", err)
	}

	v2, err := emb.Embed(ctx, "hello world")
	if err != nil {
		t.Fatalf("Embed() error = %v", err)
	}

	if len(v1) != 64 {
		t.Fatalf("expected 64 dimensions, got %d", len(v1))
	}

	for i := range v1 {
		if v1[i] != v2[i] {
			t.Fatalf("vectors differ at index %d: %f != %f", i, v1[i], v2[i])
		}
	}
}

func TestMockEmbedder_DifferentTextsDifferentVectors(t *testing.T) {
	t.Parallel()

	emb := NewMockEmbedder(64)
	ctx := context.Background()

	v1, err := emb.Embed(ctx, "golang programming")
	if err != nil {
		t.Fatalf("Embed() error = %v", err)
	}

	v2, err := emb.Embed(ctx, "python scripting")
	if err != nil {
		t.Fatalf("Embed() error = %v", err)
	}

	same := true
	for i := range v1 {
		if v1[i] != v2[i] {
			same = false
			break
		}
	}
	if same {
		t.Fatal("expected different vectors for different texts")
	}
}

func TestMockEmbedder_SimilarTextsHigherSimilarity(t *testing.T) {
	t.Parallel()

	emb := NewMockEmbedder(128)
	ctx := context.Background()

	vCompiled, _ := emb.Embed(ctx, "Go is a compiled language")
	vProgramming, _ := emb.Embed(ctx, "compiled programming")
	vDarkMode, _ := emb.Embed(ctx, "user prefers dark mode")

	// "compiled programming" should be more similar to "Go is a compiled language"
	// than "user prefers dark mode" because they share the word "compiled"
	simRelated := cosineSimilarity(vCompiled, vProgramming)
	simUnrelated := cosineSimilarity(vCompiled, vDarkMode)

	if simRelated <= simUnrelated {
		t.Fatalf("expected related texts to have higher similarity: related=%f unrelated=%f", simRelated, simUnrelated)
	}
}

func TestMockEmbedder_Dimensions(t *testing.T) {
	t.Parallel()

	emb := NewMockEmbedder(256)
	if emb.Dimensions() != 256 {
		t.Fatalf("expected 256 dimensions, got %d", emb.Dimensions())
	}
}

func TestMockEmbedder_UnitVector(t *testing.T) {
	t.Parallel()

	emb := NewMockEmbedder(64)
	ctx := context.Background()

	v, err := emb.Embed(ctx, "some text")
	if err != nil {
		t.Fatalf("Embed() error = %v", err)
	}

	var norm float64
	for _, val := range v {
		norm += float64(val) * float64(val)
	}
	norm = math.Sqrt(norm)

	if math.Abs(norm-1.0) > 1e-5 {
		t.Fatalf("expected unit vector (norm=1.0), got norm=%f", norm)
	}
}

func TestCosineSimilarity_IdenticalVectors(t *testing.T) {
	t.Parallel()

	a := []float32{1, 2, 3}
	sim := cosineSimilarity(a, a)
	if math.Abs(sim-1.0) > 1e-6 {
		t.Fatalf("expected similarity 1.0 for identical vectors, got %f", sim)
	}
}

func TestCosineSimilarity_OrthogonalVectors(t *testing.T) {
	t.Parallel()

	a := []float32{1, 0}
	b := []float32{0, 1}
	sim := cosineSimilarity(a, b)
	if math.Abs(sim) > 1e-6 {
		t.Fatalf("expected similarity 0.0 for orthogonal vectors, got %f", sim)
	}
}

func TestCosineSimilarity_DifferentLengths(t *testing.T) {
	t.Parallel()

	a := []float32{1, 2, 3}
	b := []float32{1, 2}
	sim := cosineSimilarity(a, b)
	if sim != 0 {
		t.Fatalf("expected 0 for different length vectors, got %f", sim)
	}
}

func TestSplitWords(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		want  []string
	}{
		{"Hello World", []string{"hello", "world"}},
		{"go-lang", []string{"go", "lang"}},
		{"  spaces  ", []string{"spaces"}},
		{"UPPER lower", []string{"upper", "lower"}},
		{"", nil},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := splitWords(tt.input)
			if len(got) != len(tt.want) {
				t.Fatalf("splitWords(%q) = %v, want %v", tt.input, got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("splitWords(%q)[%d] = %q, want %q", tt.input, i, got[i], tt.want[i])
				}
			}
		})
	}
}
