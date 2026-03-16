package persona

import (
	"context"
	"strings"

	"github.com/kocar/aurelia/internal/memory"
)

func (s *CanonicalIdentityService) ReprocessArchive(ctx context.Context, userID, conversationID string, limit int) error {
	if s.memory == nil {
		return nil
	}

	entries, err := s.memory.ListArchiveEntries(ctx, conversationID, limit)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if strings.TrimSpace(strings.ToLower(entry.Role)) != "user" {
			continue
		}

		if err := s.ApplyFacts(ctx, userID, extractArchiveFacts(entry.Content, userID)); err != nil {
			return err
		}
		if note, ok := extractArchiveNote(entry.Content, conversationID); ok {
			if err := s.memory.AddNote(ctx, note); err != nil {
				return err
			}
		}
	}

	return nil
}

func extractArchiveFacts(text, userID string) map[string]memory.Fact {
	return ExtractConversationFacts(text, userID)
}

func extractArchiveNote(text, conversationID string) (memory.Note, bool) {
	note, ok := ExtractConversationArchitectureNote(text, conversationID)
	if !ok {
		return memory.Note{}, false
	}
	note.Source = "archive-reprocess"
	return note, true
}


