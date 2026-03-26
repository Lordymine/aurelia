package persona

// CanonicalIdentityService loads persona files and builds a system prompt.
type CanonicalIdentityService struct {
	identityPath        string
	soulPath            string
	userPath            string
	ownerPlaybookPath   string
	lessonsLearnedPath  string
	projectPlaybookPath string
}

// NewCanonicalIdentityService creates a persona loader.
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
	}
}

// BuildPrompt loads persona files and returns the assembled system prompt.
func (s *CanonicalIdentityService) BuildPrompt() (string, error) {
	p, err := LoadPersona(s.identityPath, s.soulPath, s.userPath)
	if err != nil {
		return "", err
	}

	prompt := p.SystemPrompt

	ownerBlock := s.buildOwnerDocsBlock()
	if ownerBlock != "" {
		prompt = prompt + "\n\n" + ownerBlock
	}

	projectBlock := s.buildProjectBlock()
	if projectBlock != "" {
		prompt = prompt + "\n\n" + projectBlock
	}

	return prompt, nil
}
