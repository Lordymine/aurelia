package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/knights-analytics/hugot"
	"github.com/knights-analytics/hugot/pipelines"
)

// HugotEmbedder generates embeddings locally using the all-MiniLM-L6-v2 model
// via Hugot's pure Go backend (no CGo required).
type HugotEmbedder struct {
	pipeline *pipelines.FeatureExtractionPipeline
	session  *hugot.Session
	dims     int
	mu       sync.Mutex
}

// NewHugotEmbedder creates a local embedder. modelDir is the path where the
// ONNX model will be downloaded/cached (e.g. ~/.aurelia/models/).
func NewHugotEmbedder(modelDir string) (*HugotEmbedder, error) {
	session, err := hugot.NewGoSession()
	if err != nil {
		return nil, fmt.Errorf("hugot session: %w", err)
	}

	// Download model if not already cached
	modelPath, err := hugot.DownloadModel(
		"sentence-transformers/all-MiniLM-L6-v2",
		modelDir,
		hugot.NewDownloadOptions(),
	)
	if err != nil {
		session.Destroy()
		return nil, fmt.Errorf("hugot download model: %w", err)
	}

	config := hugot.FeatureExtractionConfig{
		ModelPath: modelPath,
		Name:      "embeddings",
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
