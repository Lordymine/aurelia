package persona

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
)

// BuildPrompt builds a system prompt from persona files.
func (s *CanonicalIdentityService) BuildPrompt(ctx context.Context, userID, conversationID string) (string, []string, error) {
	return s.BuildPromptForQuery(ctx, userID, conversationID, "")
}

// BuildPromptForQuery builds a system prompt, optionally tailored to a query.
// Memory-based retrieval was removed — only persona files and runtime context are used.
func (s *CanonicalIdentityService) BuildPromptForQuery(ctx context.Context, userID, conversationID, query string) (string, []string, error) {
	_ = query
	_ = userID
	_ = conversationID

	p, identity, err := s.ResolveIdentity(ctx)
	if err != nil {
		return "", nil, err
	}

	return s.appendRuntimeContext(p.RenderSystemPrompt(identity)), nil, nil
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
