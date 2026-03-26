package persona

import (
	"fmt"
	"strings"
)

func buildPromptBody(identityBody string, soulBytes, userBytes []byte) string {
	return fmt.Sprintf("%s\n\n%s\n\n%s",
		identityBody,
		strings.TrimSpace(string(soulBytes)),
		strings.TrimSpace(string(userBytes)),
	)
}

func buildSystemPrompt(identity CanonicalIdentity, promptBody string) string {
	return fmt.Sprintf("%s\n\n%s", buildCanonicalIdentityBlock(identity), promptBody)
}
