package telegram

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// parseGeneratedPersona extracts IDENTITY and SOUL sections from LLM output
// delimited by ===IDENTITY=== and ===SOUL=== markers. The parser is tolerant
// of common LLM quirks: extra whitespace, backtick wrapping, and markdown
// code fences around the markers.
func parseGeneratedPersona(output string) (identity, soul string, err error) {
	// Strip common LLM wrapping (```markdown ... ``` blocks)
	cleaned := stripCodeFences(output)

	identityIdx := findMarker(cleaned, "IDENTITY")
	if identityIdx == -1 {
		return "", "", fmt.Errorf("missing IDENTITY marker in LLM output")
	}

	soulIdx := findMarker(cleaned, "SOUL")
	if soulIdx == -1 {
		return "", "", fmt.Errorf("missing SOUL marker in LLM output")
	}

	// Extract content after each marker line
	identityStart := strings.IndexByte(cleaned[identityIdx:], '\n')
	if identityStart == -1 {
		return "", "", fmt.Errorf("no content after IDENTITY marker")
	}
	identityStart += identityIdx + 1

	soulStart := strings.IndexByte(cleaned[soulIdx:], '\n')
	if soulStart == -1 {
		return "", "", fmt.Errorf("no content after SOUL marker")
	}
	soulStart += soulIdx + 1

	identity = strings.TrimSpace(cleaned[identityStart:soulIdx])
	soul = strings.TrimSpace(cleaned[soulStart:])

	return identity, soul, nil
}

// findMarker locates a section marker like ===IDENTITY=== in the output,
// tolerating variations: "===IDENTITY===", "=== IDENTITY ===", "IDENTITY:",
// "## IDENTITY", etc.
func findMarker(s, name string) int {
	upper := strings.ToUpper(name)

	// Try exact markers first
	for _, pattern := range []string{
		"===" + upper + "===",
		"=== " + upper + " ===",
		"---" + upper + "---",
		"--- " + upper + " ---",
	} {
		if idx := strings.Index(strings.ToUpper(s), pattern); idx != -1 {
			return idx
		}
	}

	// Try line-based patterns: "## IDENTITY" or "IDENTITY:"
	for i, line := range strings.Split(s, "\n") {
		trimmed := strings.TrimSpace(strings.ToUpper(line))
		trimmed = strings.TrimLeft(trimmed, "#= -")
		trimmed = strings.TrimRight(trimmed, "= -:")
		trimmed = strings.TrimSpace(trimmed)
		if trimmed == upper {
			// Calculate byte offset
			offset := 0
			for _, prev := range strings.Split(s, "\n")[:i] {
				offset += len(prev) + 1
			}
			return offset
		}
	}

	return -1
}

// stripCodeFences removes wrapping ``` or ```markdown fences from LLM output.
func stripCodeFences(s string) string {
	lines := strings.Split(s, "\n")
	if len(lines) < 2 {
		return s
	}

	first := strings.TrimSpace(lines[0])
	if strings.HasPrefix(first, "```") {
		lines = lines[1:]
	}
	if len(lines) > 0 {
		last := strings.TrimSpace(lines[len(lines)-1])
		if last == "```" {
			lines = lines[:len(lines)-1]
		}
	}

	return strings.Join(lines, "\n")
}

// buildAssistantGeneratePrompt builds the user prompt for the LLM to generate
// IDENTITY.md and SOUL.md content based on the chosen preset and user description.
func buildAssistantGeneratePrompt(preset bootstrapPreset, description string) string {
	return fmt.Sprintf(`Voce e um gerador de arquivos de persona para um agente de IA.

O usuario escolheu o modo: %s (role: %s).
O template base do IDENTITY e:
%s

A descricao do usuario sobre como o assistente deve ser:
%s

Gere o conteudo final dos dois arquivos de persona. O IDENTITY.md DEVE manter o frontmatter YAML (---name/role---) e a secao de Agendamentos OBRIGATORIO do template base. Adapte o resto da personalidade e comportamento conforme a descricao do usuario.

O SOUL.md deve descrever a personalidade e tom do assistente baseado na descricao.

Responda EXATAMENTE neste formato (sem blocos de codigo, sem explicacao):

===IDENTITY===
(conteudo completo do IDENTITY.md)
===SOUL===
(conteudo completo do SOUL.md)`, preset.AgentName, preset.AgentRole, preset.IdentityTemplate, description)
}

// buildUserGeneratePrompt builds the user prompt for the LLM to generate
// USER.md content based on the user's self-description.
func buildUserGeneratePrompt(description, fallbackName string) string {
	return fmt.Sprintf(`Voce e um gerador de arquivo de perfil de usuario para um agente de IA.

A descricao que o usuario deu sobre si mesmo:
%s

Nome de fallback do Telegram: %s

Gere o conteudo do arquivo USER.md extraindo:
- Nome do usuario (use o fallback se nao encontrar na descricao)
- Fuso horario (se mencionado, senao "Relativo a sua localidade")
- Preferencias de comunicacao e trabalho
- Qualquer outro detalhe relevante sobre o usuario

Responda APENAS com o conteudo do arquivo USER.md em Markdown, sem blocos de codigo e sem explicacao. Comece com "# User".`, description, fallbackName)
}

// writeGeneratedPersona writes the generated IDENTITY.md and SOUL.md to the given directory.
func writeGeneratedPersona(dir, identity, soul string) error {
	if err := os.WriteFile(filepath.Join(dir, "IDENTITY.md"), []byte(identity), 0o644); err != nil {
		return fmt.Errorf("write IDENTITY.md: %w", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "SOUL.md"), []byte(soul), 0o644); err != nil {
		return fmt.Errorf("write SOUL.md: %w", err)
	}
	return nil
}
