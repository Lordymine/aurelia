package telegram

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/kocar/aurelia/internal/memory"
	"github.com/kocar/aurelia/internal/persona"
)

type fakeMemoryDebugService struct {
	calls []struct {
		userID         string
		conversationID string
		query          string
	}
	report *persona.LongTermMemoryDebugReport
	err    error
}

func (f *fakeMemoryDebugService) DebugLongTermMemory(ctx context.Context, userID, conversationID, query string) (*persona.LongTermMemoryDebugReport, error) {
	f.calls = append(f.calls, struct {
		userID         string
		conversationID string
		query          string
	}{
		userID:         userID,
		conversationID: conversationID,
		query:          query,
	})
	if f.err != nil {
		return nil, f.err
	}
	return f.report, nil
}

func TestMemoryCommandHandler_HandleText_Debug(t *testing.T) {
	t.Parallel()

	service := &fakeMemoryDebugService{
		report: &persona.LongTermMemoryDebugReport{
			Query:  "memoria minimalista",
			Tokens: []string{"memoria", "memory", "minimalista", "minimalism"},
			SelectedFacts: []persona.ScoredFact{
				{
					Fact: memory.Fact{Key: "project.memory.strategy", Value: "sqlite + facts + notes"},
					Score: 17,
				},
			},
			SelectedNotes: []persona.ScoredNote{
				{
					Note:  memory.Note{Topic: "architecture", Kind: "decision", Summary: "Decidido manter SQLite com facts e notes."},
					Score: 14,
				},
			},
		},
	}
	handler := NewMemoryCommandHandler(service)

	reply, err := handler.HandleText(context.Background(), "user-1", "conv-1", `/memory debug memoria minimalista`)
	if err != nil {
		t.Fatalf("HandleText() error = %v", err)
	}
	if len(service.calls) != 1 {
		t.Fatalf("expected one debug call, got %d", len(service.calls))
	}
	if service.calls[0].query != "memoria minimalista" {
		t.Fatalf("expected trimmed query, got %q", service.calls[0].query)
	}
	if !strings.Contains(reply, "project.memory.strategy") {
		t.Fatalf("expected fact in reply, got %q", reply)
	}
	if !strings.Contains(reply, "architecture/decision") {
		t.Fatalf("expected note in reply, got %q", reply)
	}
	if !strings.Contains(reply, "score=17") {
		t.Fatalf("expected fact score in reply, got %q", reply)
	}
}

func TestMemoryCommandHandler_HandleText_UsesEmptyQueryWhenNoDebugArgument(t *testing.T) {
	t.Parallel()

	service := &fakeMemoryDebugService{
		report: &persona.LongTermMemoryDebugReport{},
	}
	handler := NewMemoryCommandHandler(service)

	_, err := handler.HandleText(context.Background(), "user-1", "conv-1", `/memory`)
	if err != nil {
		t.Fatalf("HandleText() error = %v", err)
	}
	if len(service.calls) != 1 {
		t.Fatalf("expected one debug call, got %d", len(service.calls))
	}
	if service.calls[0].query != "" {
		t.Fatalf("expected empty query, got %q", service.calls[0].query)
	}
}

func TestMemoryCommandHandler_HandleText_NoMatches(t *testing.T) {
	t.Parallel()

	service := &fakeMemoryDebugService{
		report: &persona.LongTermMemoryDebugReport{
			Query:  "nada",
			Tokens: []string{"nada"},
		},
	}
	handler := NewMemoryCommandHandler(service)

	reply, err := handler.HandleText(context.Background(), "user-1", "conv-1", `/memory debug nada`)
	if err != nil {
		t.Fatalf("HandleText() error = %v", err)
	}
	if !strings.Contains(reply, "Nenhuma memoria relevante") {
		t.Fatalf("expected empty memory message, got %q", reply)
	}
}

func TestMemoryCommandHandler_HandleText_PropagatesServiceError(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("db unavailable")
	handler := NewMemoryCommandHandler(&fakeMemoryDebugService{err: expectedErr})

	_, err := handler.HandleText(context.Background(), "user-1", "conv-1", `/memory debug memoria`)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}



