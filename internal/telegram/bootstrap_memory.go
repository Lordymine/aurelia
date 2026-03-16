package telegram

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kocar/aurelia/internal/agent"
	"github.com/kocar/aurelia/internal/memory"
	"github.com/kocar/aurelia/internal/persona"
	"gopkg.in/telebot.v3"
)

func (bc *BotController) seedBootstrapIdentity(c telebot.Context, preset bootstrapPreset) error {
	ctx := agent.WithTeamContext(context.Background(), fmt.Sprintf("%d", c.Sender().ID), fmt.Sprintf("%d", c.Sender().ID))
	if bc.canonical != nil {
		return bc.canonical.SeedIdentityFacts(ctx, fmt.Sprintf("%d", c.Sender().ID), persona.CanonicalIdentity{
			AgentName: preset.AgentName,
			AgentRole: preset.AgentRole,
		}, "bootstrap")
	}
	return seedBootstrapFacts(ctx, bc.memory, fmt.Sprintf("%d", c.Sender().ID), preset.AgentName, preset.AgentRole, "")
}

func seedBootstrapFacts(ctx context.Context, mem *memory.MemoryManager, senderID, agentName, agentRole, userName string) error {
	if mem == nil {
		return nil
	}

	if strings.TrimSpace(agentName) != "" {
		if err := mem.UpsertFact(ctx, memory.Fact{Scope: "agent", EntityID: "default", Key: "agent.name", Value: agentName, Source: "bootstrap"}); err != nil {
			return err
		}
	}
	if strings.TrimSpace(agentRole) != "" {
		if err := mem.UpsertFact(ctx, memory.Fact{Scope: "agent", EntityID: "default", Key: "agent.role", Value: agentRole, Source: "bootstrap"}); err != nil {
			return err
		}
	}
	if strings.TrimSpace(userName) != "" {
		if err := mem.UpsertFact(ctx, memory.Fact{Scope: "user", EntityID: senderID, Key: "user.name", Value: userName, Source: "bootstrap"}); err != nil {
			return err
		}
	}
	return nil
}

func extractFactsFromConversation(text, senderID string) map[string]memory.Fact {
	return persona.ExtractConversationFacts(text, senderID)
}

func persistConversationFacts(ctx context.Context, mem *memory.MemoryManager, senderID, text string) error {
	if mem == nil {
		return nil
	}
	for _, fact := range extractFactsFromConversation(text, senderID) {
		if err := mem.UpsertFact(ctx, fact); err != nil {
			return err
		}
	}
	return nil
}

func extractArchitectureNoteFromConversation(text, conversationID string) (memory.Note, bool) {
	return persona.ExtractConversationArchitectureNote(text, conversationID)
}

func persistConversationNote(ctx context.Context, mem *memory.MemoryManager, conversationID, text string) error {
	if mem == nil {
		return nil
	}
	note, ok := extractArchitectureNoteFromConversation(text, conversationID)
	if !ok {
		return nil
	}
	return mem.AddNote(ctx, note)
}

func syncUserPersonaFile(path string, facts map[string]memory.Fact) error {
	return persona.SyncUserFile(path, facts)
}

func syncSoulPersonaFile(path string, facts map[string]memory.Fact) error {
	return persona.SyncSoulFile(path, facts)
}

func (bc *BotController) completeBootstrapProfile(c telebot.Context, state bootstrapState, text string) error {
	_ = state

	userTemplate := buildUserTemplateFromProfile(text, bootstrapFallbackName(c.Sender()))
	if err := os.WriteFile(filepath.Join(bc.personasDir, "USER.md"), []byte(userTemplate), 0o644); err != nil {
		log.Printf("Bootstrap user profile write error: %v\n", err)
		return SendContextText(c, bootstrapFailureMessage)
	}

	p, err := persona.LoadPersona(filepath.Join(bc.personasDir, "IDENTITY.md"), filepath.Join(bc.personasDir, "SOUL.md"), filepath.Join(bc.personasDir, "USER.md"))
	if err == nil {
		ctx := agent.WithTeamContext(context.Background(), fmt.Sprintf("%d", c.Sender().ID), fmt.Sprintf("%d", c.Sender().ID))
		if bc.canonical != nil {
			if err := bc.canonical.SeedIdentityFacts(ctx, fmt.Sprintf("%d", c.Sender().ID), p.CanonicalIdentity, "bootstrap"); err != nil {
				log.Printf("Bootstrap user fact seed warning: %v\n", err)
			}
			bootstrapFacts := map[string]memory.Fact{
				"user.preferences.summary": {
					Scope:    "user",
					EntityID: fmt.Sprintf("%d", c.Sender().ID),
					Key:      "user.preferences.summary",
					Value:    strings.TrimSpace(text),
					Source:   "bootstrap",
				},
			}
			_ = bc.canonical.ApplyFactsAndSync(ctx, fmt.Sprintf("%d", c.Sender().ID), bootstrapFacts)
		} else if strings.TrimSpace(text) != "" {
			_ = bc.memory.UpsertFact(ctx, memory.Fact{
				Scope:    "user",
				EntityID: fmt.Sprintf("%d", c.Sender().ID),
				Key:      "user.preferences.summary",
				Value:    strings.TrimSpace(text),
				Source:   "bootstrap",
			})
		}
	}

	return SendContextText(c, bootstrapSuccessMessage)
}


