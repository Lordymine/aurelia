package persona

import (
	"time"
)

// CanonicalIdentityService centralizes identity resolution and prompt building.
// Memory-backed features (facts, notes, retrieval) were removed — will be replaced
// by semantic memory via the bridge in a later task.
type CanonicalIdentityService struct {
	identityPath        string
	soulPath            string
	userPath            string
	ownerPlaybookPath   string
	lessonsLearnedPath  string
	projectPlaybookPath string
	now                 func() time.Time
	location            *time.Location
}

// NewCanonicalIdentityService creates a canonical identity service without memory backing.
func NewCanonicalIdentityService(
	identityPath, soulPath, userPath string,
	ownerPlaybookPath, lessonsLearnedPath string,
	projectPlaybookPath string,
) *CanonicalIdentityService {
	return &CanonicalIdentityService{
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
