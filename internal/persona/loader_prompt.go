package persona

import (
	"fmt"
	"strings"
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

// RenderSystemPrompt assembles the final prompt with canonical identity.
func (p *Persona) RenderSystemPrompt(identity CanonicalIdentity) string {
	if p == nil {
		return ""
	}

	sections := []string{buildCanonicalIdentityBlock(identity)}
	sections = append(sections, strings.TrimSpace(p.PromptBody))

	return strings.Join(sections, "\n\n")
}
