package persona

import (
	"sort"
	"strings"

	"github.com/kocar/aurelia/internal/memory"
)

func rankNotes(notes []memory.Note, tokens []string, limit int) []ScoredNote {
	if len(notes) == 0 {
		return nil
	}
	var scored []ScoredNote
	for _, note := range notes {
		score := note.Importance
		if len(tokens) == 0 {
			score += 1
		}
		topic := strings.ToLower(note.Topic)
		summary := strings.ToLower(note.Summary)
		for _, token := range tokens {
			if strings.Contains(topic, token) {
				score += 5
			}
			if strings.Contains(summary, token) {
				score += 3
			}
		}
		score += thematicNoteBonus(topic, summary, tokens)
		if score > note.Importance || len(tokens) == 0 {
			scored = append(scored, ScoredNote{Note: note, Score: score})
		}
	}
	sort.SliceStable(scored, func(i, j int) bool { return scored[i].Score > scored[j].Score })
	if limit <= 0 || limit > len(scored) {
		limit = len(scored)
	}
	return append([]ScoredNote(nil), scored[:limit]...)
}

func thematicFactBonus(key string, tokens []string) int {
	themes := map[string][]string{
		"memory":       {"project.memory.", "project.constraint.no_embeddings_core"},
		"architecture": {"project.architecture.", "project.goal"},
		"minimalism":   {"project.constraint.minimalism", "project.goal"},
		"strategy":     {"project.memory.strategy"},
	}

	score := 0
	for _, token := range tokens {
		for _, prefix := range themes[token] {
			if strings.Contains(key, prefix) {
				score += 4
			}
		}
	}
	return score
}

func thematicNoteBonus(topic, summary string, tokens []string) int {
	score := 0
	for _, token := range tokens {
		switch token {
		case "memory", "strategy":
			if strings.Contains(topic, "memory") || strings.Contains(summary, "sqlite") || strings.Contains(summary, "embedding") {
				score += 4
			}
		case "architecture":
			if strings.Contains(topic, "architecture") || strings.Contains(summary, "monolito") || strings.Contains(summary, "arquitet") {
				score += 4
			}
		case "minimalism":
			if strings.Contains(summary, "minimal") {
				score += 4
			}
		}
	}
	return score
}


