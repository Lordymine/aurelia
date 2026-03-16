package persona

import (
	"fmt"
	"os"
	"strings"

	"github.com/kocar/aurelia/internal/memory"
)

func (s *CanonicalIdentityService) SyncFactsToFiles(facts map[string]memory.Fact) error {
	if len(facts) == 0 {
		return nil
	}
	if err := SyncUserFile(s.userPath, facts); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := SyncSoulFile(s.soulPath, facts); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func SyncUserFile(path string, facts map[string]memory.Fact) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	name := resolveUserFileName(string(content), facts)
	preferences := resolveUserFilePreference(string(content), facts)
	responseStyle := resolveUserFileResponseStyle(string(content), facts)

	lines := []string{"# User", fmt.Sprintf("Nome: %s", name), "Fuso horario: Relativo a sua localidade."}
	if preferences != "" {
		lines = append(lines, fmt.Sprintf("Preferencias: %s", preferences))
	}
	if responseStyle != "" {
		lines = append(lines, fmt.Sprintf("Estilo de resposta: %s", responseStyle))
	}

	return os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0o644)
}

func resolveUserFileName(content string, facts map[string]memory.Fact) string {
	if fact, ok := facts["user.name"]; ok && strings.TrimSpace(fact.Value) != "" {
		return fact.Value
	}
	name := extractCanonicalUserName(content)
	if name == "" {
		return "Nao definido"
	}
	return name
}

func resolveUserFilePreference(content string, facts map[string]memory.Fact) string {
	if fact, ok := facts["user.preferences.summary"]; ok {
		return fact.Value
	}
	return extractLineValue(content, "Preferencias:")
}

func resolveUserFileResponseStyle(content string, facts map[string]memory.Fact) string {
	if fact, ok := facts["user.preference.response_style"]; ok {
		return fact.Value
	}
	return extractLineValue(content, "Estilo de resposta:")
}

func SyncSoulFile(path string, facts map[string]memory.Fact) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	base := strings.TrimSpace(string(content))
	styleSummary := resolveSoulStyleSummary(base, facts)
	if styleSummary == "" {
		return nil
	}

	lines := []string{stripSoulPreference(base), fmt.Sprintf("Preferencias persistentes do usuario sobre tom: %s", styleSummary)}
	return os.WriteFile(path, []byte(strings.TrimSpace(strings.Join(lines, "\n"))+"\n"), 0o644)
}

func resolveSoulStyleSummary(content string, facts map[string]memory.Fact) string {
	if fact, ok := facts["agent.style.summary"]; ok {
		return fact.Value
	}
	return extractSoulPreference(content)
}

func extractLineValue(content, prefix string) string {
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(trimmed, prefix))
		}
	}
	return ""
}

func extractSoulPreference(content string) string {
	return extractLineValue(content, "Preferencias persistentes do usuario sobre tom:")
}

func stripSoulPreference(content string) string {
	var lines []string
	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "Preferencias persistentes do usuario sobre tom:") {
			continue
		}
		lines = append(lines, line)
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}


