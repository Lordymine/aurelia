package memory

import (
	"context"
	"database/sql"
	"encoding/binary"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

// Memory represents a stored memory entry.
type Memory struct {
	ID        string
	Content   string
	Category  string // "fact", "conversation", "decision", "preference"
	Agent     string
	CreatedAt time.Time
}

// Store provides semantic memory with embedding-based search.
type Store struct {
	db       *sql.DB
	embedder Embedder
}

// NewStore creates a new memory store backed by SQLite.
func NewStore(dbPath string, embedder Embedder) (*Store, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if _, err := db.Exec(schemaSQL); err != nil {
		db.Close()
		return nil, fmt.Errorf("create schema: %w", err)
	}

	return &Store{db: db, embedder: embedder}, nil
}

// Save stores a new memory with its embedding.
func (s *Store) Save(ctx context.Context, content, category, agent string) error {
	id := uuid.New().String()
	now := time.Now().UTC()

	embedding, err := s.embedder.Embed(ctx, content)
	if err != nil {
		return fmt.Errorf("generate embedding: %w", err)
	}

	blob := serializeEmbedding(embedding)

	_, err = s.db.ExecContext(ctx,
		`INSERT INTO memories (id, content, category, agent, created_at, embedding) VALUES (?, ?, ?, ?, ?, ?)`,
		id, content, category, agent, now.Format(time.RFC3339), blob,
	)
	if err != nil {
		return fmt.Errorf("insert memory: %w", err)
	}

	return nil
}

// Search finds the most similar memories to the query using cosine similarity.
func (s *Store) Search(ctx context.Context, query string, limit int) ([]Memory, error) {
	queryEmbedding, err := s.embedder.Embed(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("generate query embedding: %w", err)
	}

	rows, err := s.db.QueryContext(ctx,
		`SELECT id, content, category, agent, created_at, embedding FROM memories WHERE embedding IS NOT NULL`,
	)
	if err != nil {
		return nil, fmt.Errorf("query memories: %w", err)
	}
	defer rows.Close()

	type scored struct {
		memory Memory
		score  float64
	}

	var results []scored
	for rows.Next() {
		var m Memory
		var createdAt string
		var blob []byte

		if err := rows.Scan(&m.ID, &m.Content, &m.Category, &m.Agent, &createdAt, &blob); err != nil {
			return nil, fmt.Errorf("scan memory: %w", err)
		}

		m.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
		if err != nil {
			return nil, fmt.Errorf("parse created_at: %w", err)
		}

		embedding := deserializeEmbedding(blob)
		score := cosineSimilarity(queryEmbedding, embedding)
		results = append(results, scored{memory: m, score: score})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	if limit > len(results) {
		limit = len(results)
	}

	memories := make([]Memory, limit)
	for i := 0; i < limit; i++ {
		memories[i] = results[i].memory
	}

	return memories, nil
}

// Inject searches for relevant memories and formats them as a markdown block.
func (s *Store) Inject(ctx context.Context, query string, limit int) (string, error) {
	memories, err := s.Search(ctx, query, limit)
	if err != nil {
		return "", err
	}

	if len(memories) == 0 {
		return "", nil
	}

	var b strings.Builder
	b.WriteString("## Relevant Memories\n\n")
	for _, m := range memories {
		b.WriteString(fmt.Sprintf("- [%s] %s (%s)\n", m.Category, m.Content, m.CreatedAt.Format("2006-01-02")))
	}

	return b.String(), nil
}

// Close closes the underlying database connection.
func (s *Store) Close() error {
	return s.db.Close()
}

// serializeEmbedding converts []float32 to binary BLOB.
func serializeEmbedding(vec []float32) []byte {
	buf := make([]byte, len(vec)*4)
	for i, v := range vec {
		binary.LittleEndian.PutUint32(buf[i*4:], math.Float32bits(v))
	}
	return buf
}

// deserializeEmbedding converts binary BLOB back to []float32.
func deserializeEmbedding(blob []byte) []float32 {
	vec := make([]float32, len(blob)/4)
	for i := range vec {
		vec[i] = math.Float32frombits(binary.LittleEndian.Uint32(blob[i*4:]))
	}
	return vec
}

// cosineSimilarity computes dot(a, b) / (norm(a) * norm(b)).
func cosineSimilarity(a, b []float32) float64 {
	if len(a) != len(b) {
		return 0
	}

	var dot, normA, normB float64
	for i := range a {
		dot += float64(a[i]) * float64(b[i])
		normA += float64(a[i]) * float64(a[i])
		normB += float64(b[i]) * float64(b[i])
	}

	denom := math.Sqrt(normA) * math.Sqrt(normB)
	if denom == 0 {
		return 0
	}

	return dot / denom
}
