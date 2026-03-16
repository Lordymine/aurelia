package persona

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"unicode"

	"github.com/kocar/aurelia/internal/memory"
)

func (s *CanonicalIdentityService) DebugLongTermMemory(ctx context.Context, userID, conversationID, query string) (*LongTermMemoryDebugReport, error) {
	report := &LongTermMemoryDebugReport{Query: query}
	if s.memory == nil {
		return report, nil
	}

	facts, notes, tokens, err := s.retrieveLongTermMemory(ctx, userID, conversationID, query)
	if err != nil {
		return nil, err
	}
	report.Tokens = tokens
	report.SelectedFacts = facts
	report.SelectedNotes = notes
	return report, nil
}

func (s *CanonicalIdentityService) retrieveLongTermMemory(ctx context.Context, userID, conversationID, query string) ([]ScoredFact, []ScoredNote, []string, error) {
	candidatesFacts, err := s.loadLongTermFactCandidates(ctx, userID)
	if err != nil {
		return nil, nil, nil, err
	}

	notes, err := s.loadLongTermNotes(ctx, conversationID)
	if err != nil {
		return nil, nil, nil, err
	}

	tokens := tokenizeQuery(query)
	return rankFacts(candidatesFacts, tokens, 5), rankNotes(notes, tokens, 5), tokens, nil
}

func (s *CanonicalIdentityService) loadLongTermFactCandidates(ctx context.Context, userID string) ([]memory.Fact, error) {
	var candidatesFacts []memory.Fact
	for _, pair := range []struct {
		scope    string
		entityID string
	}{
		{scope: "project", entityID: "default"},
		{scope: "agent", entityID: "default"},
		{scope: "user", entityID: userID},
	} {
		if strings.TrimSpace(pair.entityID) == "" {
			continue
		}
		facts, err := s.memory.ListFacts(ctx, pair.scope, pair.entityID)
		if err != nil {
			return nil, fmt.Errorf("failed to list facts: %w", err)
		}
		candidatesFacts = append(candidatesFacts, facts...)
	}
	return candidatesFacts, nil
}

func (s *CanonicalIdentityService) loadLongTermNotes(ctx context.Context, conversationID string) ([]memory.Note, error) {
	notes, err := s.memory.ListRecentNotes(ctx, conversationID, 20)
	if err != nil && strings.TrimSpace(conversationID) != "" {
		return nil, fmt.Errorf("failed to load recent notes: %w", err)
	}
	return notes, nil
}

func tokenizeQuery(query string) []string {
	normalized := strings.ToLower(strings.TrimSpace(query))
	if normalized == "" {
		return nil
	}

	var b strings.Builder
	for _, r := range normalized {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) {
			b.WriteRune(r)
			continue
		}
		b.WriteRune(' ')
	}

	stopwords := map[string]struct{}{
		"o": {}, "a": {}, "os": {}, "as": {}, "de": {}, "do": {}, "da": {}, "e": {}, "que": {}, "sobre": {}, "um": {}, "uma": {}, "ja": {}, "qual": {},
	}
	parts := strings.Fields(b.String())
	seen := make(map[string]struct{})
	var tokens []string
	for _, part := range parts {
		if _, skip := stopwords[part]; skip || len(part) < 3 {
			continue
		}
		if _, ok := seen[part]; ok {
			continue
		}
		seen[part] = struct{}{}
		tokens = append(tokens, part)
	}
	return expandQueryAliases(tokens)
}

func expandQueryAliases(tokens []string) []string {
	aliases := map[string][]string{
		"memoria":     {"memory", "strategy"},
		"longo":       {"long", "longterm"},
		"prazo":       {"term"},
		"decidimos":   {"decision"},
		"decisao":     {"decision"},
		"arquitetura": {"architecture"},
		"minimalista": {"minimalism"},
		"minimalismo": {"minimalism"},
		"estrategia":  {"strategy"},
	}

	seen := make(map[string]struct{})
	var expanded []string
	for _, token := range tokens {
		if _, ok := seen[token]; !ok {
			seen[token] = struct{}{}
			expanded = append(expanded, token)
		}
		for _, alias := range aliases[token] {
			if _, ok := seen[alias]; ok {
				continue
			}
			seen[alias] = struct{}{}
			expanded = append(expanded, alias)
		}
	}
	return expanded
}

func rankFacts(facts []memory.Fact, tokens []string, limit int) []ScoredFact {
	if len(facts) == 0 {
		return nil
	}
	var scored []ScoredFact
	for _, fact := range facts {
		score := 0
		if len(tokens) == 0 {
			score = 1
		}
		key := strings.ToLower(fact.Key)
		value := strings.ToLower(fact.Value)
		for _, token := range tokens {
			if strings.Contains(key, token) {
				score += 5
			}
			if strings.Contains(value, token) {
				score += 3
			}
		}
		score += thematicFactBonus(key, tokens)
		if score > 0 {
			scored = append(scored, ScoredFact{Fact: fact, Score: score})
		}
	}
	sort.SliceStable(scored, func(i, j int) bool { return scored[i].Score > scored[j].Score })
	if limit <= 0 || limit > len(scored) {
		limit = len(scored)
	}
	return append([]ScoredFact(nil), scored[:limit]...)
}


