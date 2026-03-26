package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/knights-analytics/hugot"
	"github.com/knights-analytics/hugot/pipelines"
)

// HugotEmbedder generates embeddings locally using a pre-downloaded ONNX model
// via Hugot's pure Go backend (no CGo required).
type HugotEmbedder struct {
	pipeline *pipelines.FeatureExtractionPipeline
	session  *hugot.Session
	dims     int
	mu       sync.Mutex
}

// NewHugotEmbedder creates a local embedder from a pre-downloaded model directory.
// The directory must contain model.onnx, tokenizer.json, and config.json.
func NewHugotEmbedder(modelDir string) (*HugotEmbedder, error) {
	session, err := hugot.NewGoSession()
	if err != nil {
		return nil, fmt.Errorf("hugot session: %w", err)
	}

	config := hugot.FeatureExtractionConfig{
		ModelPath:    modelDir,
		Name:         "embeddings",
		OnnxFilename: "model.onnx",
	}

	pipeline, err := hugot.NewPipeline(session, config)
	if err != nil {
		session.Destroy()
		return nil, fmt.Errorf("hugot pipeline: %w", err)
	}

	return &HugotEmbedder{
		pipeline: pipeline,
		session:  session,
		dims:     384, // all-MiniLM-L6-v2 output dimension
	}, nil
}

func (e *HugotEmbedder) Embed(_ context.Context, text string) ([]float32, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	result, err := e.pipeline.RunPipeline([]string{text})
	if err != nil {
		return nil, fmt.Errorf("hugot embed: %w", err)
	}
	if len(result.Embeddings) == 0 {
		return nil, fmt.Errorf("hugot returned no embeddings")
	}
	return result.Embeddings[0], nil
}

func (e *HugotEmbedder) Dimensions() int {
	return e.dims
}

// Close releases the Hugot session resources.
func (e *HugotEmbedder) Close() error {
	return e.session.Destroy()
}
