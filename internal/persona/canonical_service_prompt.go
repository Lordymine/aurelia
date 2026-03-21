package persona

import (
	"log"
	"strings"
)

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
