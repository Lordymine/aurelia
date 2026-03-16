package persona

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/kocar/aurelia/internal/memory"
)

func (s *CanonicalIdentityService) BuildPrompt(ctx context.Context, userID, conversationID string) (string, []string, error) {
	return s.BuildPromptForQuery(ctx, userID, conversationID, "")
}

func (s *CanonicalIdentityService) BuildPromptForQuery(ctx context.Context, userID, conversationID, query string) (string, []string, error) {
	p, identity, err := s.ResolveIdentity(ctx, userID)
	if err != nil {
		return "", nil, err
	}

	facts, notes, err := s.selectedLongTermMemory(ctx, userID, conversationID, query)
	if err != nil {
		return "", nil, err
	}

	// Runtime tools come from the live ToolRegistry, not from persona markdown.
	// The persona remains responsible for identity/tone/instructions, while the
	// execution layer injects the actual capabilities available in this process.
	return s.appendRuntimeContext(p.RenderSystemPrompt(identity, facts, notes)), nil, nil
}

func (s *CanonicalIdentityService) selectedLongTermMemory(ctx context.Context, userID, conversationID, query string) ([]memory.Fact, []memory.Note, error) {
	if s.memory == nil {
		return nil, nil, nil
	}

	report, err := s.DebugLongTermMemory(ctx, userID, conversationID, query)
	if err != nil {
		return nil, nil, err
	}

	var facts []memory.Fact
	var notes []memory.Note
	for _, scored := range report.SelectedFacts {
		facts = append(facts, scored.Fact)
	}
	for _, scored := range report.SelectedNotes {
		notes = append(notes, scored.Note)
	}
	return facts, notes, nil
}

func (s *CanonicalIdentityService) appendRuntimeContext(prompt string) string {
	nowFn := s.now
	if nowFn == nil {
		nowFn = time.Now
	}

	location := s.location
	if location == nil {
		location = time.Local
	}

	now := nowFn().In(location)
	runtimeBlock := strings.Join([]string{
		"# RUNTIME CONTEXT",
		fmt.Sprintf("Data local atual: %s", now.Format("2006-01-02")),
		fmt.Sprintf("Horario local atual: %s", now.Format(time.RFC3339)),
		fmt.Sprintf("Fuso horario atual: %s", location.String()),
		"Use essa referencia para interpretar corretamente pedidos relativos como hoje, ontem, amanha e horarios agendados.",
	}, "\n")

	ownerBlock := s.buildOwnerDocsBlock()
	if ownerBlock != "" {
		runtimeBlock = runtimeBlock + "\n\n" + ownerBlock
	}

	projectBlock := s.buildProjectBlock()
	if projectBlock != "" {
		runtimeBlock = runtimeBlock + "\n\n" + projectBlock
	}

	if strings.TrimSpace(prompt) == "" {
		return runtimeBlock
	}

	return runtimeBlock + "\n\n" + prompt
}

// buildProjectBlock reads the optional project playbook from the working directory
// and returns a PROJECT CONTEXT markdown block, or an empty string if the path is
// empty, the file does not exist, or its content is blank.
func (s *CanonicalIdentityService) buildProjectBlock() string {
	if s.projectPlaybookPath == "" {
		return ""
	}

	content, err := readOptionalFile(s.projectPlaybookPath)
	if err != nil {
		log.Printf("Warning: failed to read project playbook at %q: %v", s.projectPlaybookPath, err)
	}

	if strings.TrimSpace(content) == "" {
		return ""
	}

	return strings.Join([]string{"# PROJECT CONTEXT", "## Project Playbook", content}, "\n")
}

// buildOwnerDocsBlock reads optional owner documents and returns an OWNER CONTEXT
// markdown block, or an empty string if neither file exists or both are empty.
func (s *CanonicalIdentityService) buildOwnerDocsBlock() string {
	playbookContent, err := readOptionalFile(s.ownerPlaybookPath)
	if err != nil {
		log.Printf("Warning: failed to read owner playbook at %q: %v", s.ownerPlaybookPath, err)
	}

	lessonsContent, err := readOptionalFile(s.lessonsLearnedPath)
	if err != nil {
		log.Printf("Warning: failed to read lessons learned at %q: %v", s.lessonsLearnedPath, err)
	}

	if strings.TrimSpace(playbookContent) == "" && strings.TrimSpace(lessonsContent) == "" {
		return ""
	}

	var sections []string
	sections = append(sections, "# OWNER CONTEXT")

	if strings.TrimSpace(playbookContent) != "" {
		sections = append(sections, "## Owner Playbook")
		sections = append(sections, playbookContent)
	}

	if strings.TrimSpace(lessonsContent) != "" {
		sections = append(sections, "## Lessons Learned")
		sections = append(sections, lessonsContent)
	}

	return strings.Join(sections, "\n")
}


