package persona

import (
	"fmt"
	"strings"

	"github.com/kocar/aurelia/internal/memory"
)

func buildPromptBody(identityBody string, soulBytes, userBytes []byte) string {
	return fmt.Sprintf("%s\n\n%s\n\n%s",
		identityBody,
		string(bytesTrimSpace(soulBytes)),
		string(bytesTrimSpace(userBytes)),
	)
}

func buildSystemPrompt(identity CanonicalIdentity, promptBody string) string {
	return fmt.Sprintf("%s\n\n%s", buildCanonicalIdentityBlock(identity), promptBody)
}

func bytesTrimSpace(content []byte) []byte {
	return []byte(strings.TrimSpace(string(content)))
}

// RenderSystemPrompt assembles the final prompt with canonical identity and long-term memory.
func (p *Persona) RenderSystemPrompt(identity CanonicalIdentity, facts []memory.Fact, notes []memory.Note) string {
	if p == nil {
		return ""
	}

	sections := []string{buildCanonicalIdentityBlock(identity)}
	if memoryBlock := buildLongTermMemoryBlock(facts, notes); memoryBlock != "" {
		sections = append(sections, memoryBlock)
	}
	sections = append(sections, strings.TrimSpace(p.PromptBody))

	return strings.Join(sections, "\n\n")
}


