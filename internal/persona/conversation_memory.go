package persona

import (
	"regexp"
	"strings"

	"github.com/kocar/aurelia/internal/memory"
)

var (
	conversationMeChamoPattern = regexp.MustCompile(`(?i)\bme chamo\s+([A-Za-zÀ-ÿ][A-Za-zÀ-ÿ\s'-]{0,60})`)
	conversationMeuNomePattern = regexp.MustCompile(`(?i)\bmeu nome e\s+([A-Za-zÀ-ÿ][A-Za-zÀ-ÿ\s'-]{0,60})`)
	conversationSouPattern     = regexp.MustCompile(`(?i)\bsou\s+([A-Za-zÀ-ÿ][A-Za-zÀ-ÿ\s'-]{0,60})`)
	conversationSeuNomePattern = regexp.MustCompile(`(?i)\bseu nome e\s+([A-Za-zÀ-ÿ][A-Za-zÀ-ÿ\s'-]{0,60})`)
	conversationSeuPapelPattern = regexp.MustCompile(`(?i)\bseu papel e\s+([A-Za-zÀ-ÿ][A-Za-zÀ-ÿ\s'-]{0,60})`)
)

func ExtractNameFromProfile(profileText string) string {
	text := strings.TrimSpace(profileText)
	for _, pattern := range []*regexp.Regexp{conversationMeChamoPattern, conversationMeuNomePattern, conversationSouPattern} {
		matches := pattern.FindStringSubmatch(text)
		if len(matches) < 2 {
			continue
		}

		name := strings.TrimSpace(matches[1])
		name = strings.TrimRight(name, ".,;:!?")
		for _, separator := range []string{" e quero ", " e prefiro ", " e gosto ", " e "} {
			if idx := strings.Index(strings.ToLower(name), separator); idx >= 0 {
				name = strings.TrimSpace(name[:idx])
				break
			}
		}
		name = strings.Join(strings.Fields(name), " ")
		if name != "" {
			return name
		}
	}

	return ""
}

func ExtractConversationFacts(text, senderID string) map[string]memory.Fact {
	facts := make(map[string]memory.Fact)
	trimmedText := strings.TrimSpace(text)
	if trimmedText == "" {
		return facts
	}

	if userName := ExtractNameFromProfile(trimmedText); userName != "" {
		facts["user.name"] = memory.Fact{Scope: "user", EntityID: senderID, Key: "user.name", Value: userName, Source: "conversation"}
	}

	lowered := strings.ToLower(trimmedText)
	if strings.Contains(lowered, "prefiro ") || strings.Contains(lowered, "quero respostas") || strings.Contains(lowered, "sem floreios") || strings.Contains(lowered, "sem emojis") {
		facts["user.preferences.summary"] = memory.Fact{Scope: "user", EntityID: senderID, Key: "user.preferences.summary", Value: trimmedText, Source: "conversation"}
	}
	if strings.Contains(lowered, "diret") {
		facts["user.preference.response_style"] = memory.Fact{Scope: "user", EntityID: senderID, Key: "user.preference.response_style", Value: "direto", Source: "conversation"}
	}
	if strings.Contains(lowered, "sem floreios") {
		facts["user.preference.no_fluff"] = memory.Fact{Scope: "user", EntityID: senderID, Key: "user.preference.no_fluff", Value: "true", Source: "conversation"}
	}
	if strings.Contains(lowered, "sem emojis") {
		facts["user.preference.no_emoji"] = memory.Fact{Scope: "user", EntityID: senderID, Key: "user.preference.no_emoji", Value: "true", Source: "conversation"}
	}
	if agentName := extractConversationPatternValue(trimmedText, conversationSeuNomePattern); agentName != "" {
		facts["agent.name"] = memory.Fact{Scope: "agent", EntityID: "default", Key: "agent.name", Value: agentName, Source: "conversation"}
	}
	if agentRole := extractConversationPatternValue(trimmedText, conversationSeuPapelPattern); agentRole != "" {
		facts["agent.role"] = memory.Fact{Scope: "agent", EntityID: "default", Key: "agent.role", Value: agentRole, Source: "conversation"}
	}
	if strings.Contains(lowered, "quero que voce seja") || strings.Contains(lowered, "seja mais direto") || strings.Contains(lowered, "seja mais sarcast") {
		facts["agent.style.summary"] = memory.Fact{Scope: "agent", EntityID: "default", Key: "agent.style.summary", Value: trimmedText, Source: "conversation"}
	}
	if strings.Contains(lowered, "projeto minimalista") || strings.Contains(lowered, "manter isso minimalista") || strings.Contains(lowered, "quero esse projeto minimalista") {
		facts["project.goal"] = memory.Fact{Scope: "project", EntityID: "default", Key: "project.goal", Value: "manter o projeto minimalista", Source: "conversation"}
	}
	if strings.Contains(lowered, "nao quero depender de embeddings no core") || strings.Contains(lowered, "evitar embeddings no core") {
		facts["project.constraint.no_embeddings_core"] = memory.Fact{Scope: "project", EntityID: "default", Key: "project.constraint.no_embeddings_core", Value: "true", Source: "conversation"}
	}
	if strings.Contains(lowered, "monolito modular") {
		facts["project.architecture.style"] = memory.Fact{Scope: "project", EntityID: "default", Key: "project.architecture.style", Value: "monolito modular", Source: "conversation"}
	}
	if strings.Contains(lowered, "sqlite com facts e notes") || strings.Contains(lowered, "sqlite + facts + notes") || strings.Contains(lowered, "usando sqlite com facts e notes") {
		facts["project.memory.strategy"] = memory.Fact{Scope: "project", EntityID: "default", Key: "project.memory.strategy", Value: "sqlite + facts + notes", Source: "conversation"}
	}
	if strings.Contains(lowered, "minimalista") {
		facts["project.constraint.minimalism"] = memory.Fact{Scope: "project", EntityID: "default", Key: "project.constraint.minimalism", Value: "true", Source: "conversation"}
	}

	return facts
}

func ExtractConversationArchitectureNote(text, conversationID string) (memory.Note, bool) {
	trimmedText := strings.TrimSpace(text)
	lowered := strings.ToLower(trimmedText)
	if trimmedText == "" {
		return memory.Note{}, false
	}
	if !strings.Contains(lowered, "decid") && !strings.Contains(lowered, "vamos manter") && !strings.Contains(lowered, "arquitet") && !strings.Contains(lowered, "memoria") {
		return memory.Note{}, false
	}

	return memory.Note{
		ConversationID: conversationID,
		Topic:          "architecture",
		Summary:        trimmedText,
		Kind:           "decision",
		Importance:     8,
		Source:         "conversation",
	}, true
}

func extractConversationPatternValue(text string, pattern *regexp.Regexp) string {
	matches := pattern.FindStringSubmatch(strings.TrimSpace(text))
	if len(matches) < 2 {
		return ""
	}

	value := strings.TrimSpace(matches[1])
	value = strings.TrimRight(value, ".,;:!?")
	for _, separator := range []string{" e seu ", " e voce ", " e "} {
		if idx := strings.Index(strings.ToLower(value), separator); idx >= 0 {
			value = strings.TrimSpace(value[:idx])
			break
		}
	}
	return strings.Join(strings.Fields(value), " ")
}


