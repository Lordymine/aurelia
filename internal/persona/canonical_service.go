package persona

import (
	"context"
	"time"

	"github.com/kocar/aurelia/internal/memory"
)

// CanonicalMemoryStore defines the storage used by canonical identity resolution.
type CanonicalMemoryStore interface {
	UpsertFact(ctx context.Context, fact memory.Fact) error
	GetFact(ctx context.Context, scope, entityID, key string) (memory.Fact, bool, error)
	ListFacts(ctx context.Context, scope, entityID string) ([]memory.Fact, error)
	AddNote(ctx context.Context, note memory.Note) error
	ListRecentNotes(ctx context.Context, conversationID string, limit int) ([]memory.Note, error)
	ListArchiveEntries(ctx context.Context, conversationID string, limit int) ([]memory.ArchiveEntry, error)
}

// CanonicalIdentityService centralizes canonical fact precedence, file sync and prompt building.
type CanonicalIdentityService struct {
	memory               CanonicalMemoryStore
	identityPath         string
	soulPath             string
	userPath             string
	ownerPlaybookPath    string
	lessonsLearnedPath   string
	projectPlaybookPath  string
	now                  func() time.Time
	location             *time.Location
}

type ScoredFact struct {
	Fact  memory.Fact
	Score int
}

type ScoredNote struct {
	Note  memory.Note
	Score int
}

type LongTermMemoryDebugReport struct {
	Query         string
	Tokens        []string
	SelectedFacts []ScoredFact
	SelectedNotes []ScoredNote
}

func NewCanonicalIdentityService(
	memory CanonicalMemoryStore,
	identityPath, soulPath, userPath string,
	ownerPlaybookPath, lessonsLearnedPath string,
	projectPlaybookPath string,
) *CanonicalIdentityService {
	return &CanonicalIdentityService{
		memory:              memory,
		identityPath:        identityPath,
		soulPath:            soulPath,
		userPath:            userPath,
		ownerPlaybookPath:   ownerPlaybookPath,
		lessonsLearnedPath:  lessonsLearnedPath,
		projectPlaybookPath: projectPlaybookPath,
		now:                 time.Now,
		location:            time.Local,
	}
}


