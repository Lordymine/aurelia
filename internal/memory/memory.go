package memory

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

// Message represents a single chat message
type Message struct {
	ID             int64
	ConversationID string
	Role           string
	Content        string
	CreatedAt      time.Time
}

// Fact represents a durable canonical fact stored between sessions.
type Fact struct {
	ID        int64
	Scope     string
	EntityID  string
	Key       string
	Value     string
	Source    string
	UpdatedAt time.Time
}

// Note represents a compact long-term memory note.
type Note struct {
	ID             int64
	ConversationID string
	Topic          string
	Summary        string
	Kind           string
	Importance     int
	Source         string
	CreatedAt      time.Time
}

// ArchiveEntry stores raw long-term conversation history outside the working window.
type ArchiveEntry struct {
	ID             int64
	ConversationID string
	SessionID      string
	Role           string
	Content        string
	MessageType    string
	CreatedAt      time.Time
}

// MemoryManager handles interactions with the SQLite database
type MemoryManager struct {
	db               *sql.DB
	memoryWindowSize int
}

// NewMemoryManager creates a new MemoryManager instance and initializes tables
func NewMemoryManager(dbPath string, memoryWindowSize int) (*MemoryManager, error) {
	db, err := sql.Open("sqlite", dbPath+"?_journal_mode=WAL")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err := initializeDB(db); err != nil {
		return nil, err
	}

	return &MemoryManager{
		db:               db,
		memoryWindowSize: memoryWindowSize,
	}, nil
}

// Close closes the database connection
func (m *MemoryManager) Close() error {
	return m.db.Close()
}
