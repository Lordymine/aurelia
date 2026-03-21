package persona

import (
	"fmt"
	"regexp"
	"strings"
)

// Config holds the persona frontmatter configuration.
type Config struct {
	Name             string   `yaml:"name"`
	Role             string   `yaml:"role"`
	MemoryWindowSize int      `yaml:"memory_window_size"`
	Tools            []string `yaml:"tools"`
}

// Persona is the parsed identity package used by the prompt builder.
type Persona struct {
	Config            Config
	SystemPrompt      string
	PromptBody        string
	CanonicalIdentity CanonicalIdentity
}

var placeholderUserPattern = regexp.MustCompile(`(?i)^usuario\s+\d+$`)

// CanonicalIdentity holds resolved identity values before prompt assembly.
type CanonicalIdentity struct {
	AgentName string
	AgentRole string
	UserName  string
}

// LoadPersona reads IDENTITY.md, SOUL.md and USER.md into a single persona.
func LoadPersona(identityPath, soulPath, userPath string) (*Persona, error) {
	identityBytes, err := readPersonaFile(identityPath, "identity")
	if err != nil {
		return nil, err
	}
	soulBytes, err := readPersonaFile(soulPath, "soul")
	if err != nil {
		return nil, err
	}
	userBytes, err := readPersonaFile(userPath, "user")
	if err != nil {
		return nil, err
	}

	config, identityBody, err := parseIdentityFrontmatter(identityBytes)
	if err != nil {
		return nil, err
	}

	promptBody := buildPromptBody(identityBody, soulBytes, userBytes)

	canonicalIdentity := CanonicalIdentity{
		AgentName: canonicalValue(config.Name),
		AgentRole: canonicalValue(config.Role),
		UserName:  canonicalUserName(userBytes),
	}

	sysPrompt := buildSystemPrompt(canonicalIdentity, promptBody)

	return &Persona{
		Config:            config,
		SystemPrompt:      sysPrompt,
		PromptBody:        promptBody,
		CanonicalIdentity: canonicalIdentity,
	}, nil
}

func buildCanonicalIdentityBlock(identity CanonicalIdentity) string {
	lines := []string{
		"# CANONICAL IDENTITY",
		fmt.Sprintf("Nome canonico do agente: %s", canonicalValue(identity.AgentName)),
		fmt.Sprintf("Papel canonico do agente: %s", canonicalValue(identity.AgentRole)),
		fmt.Sprintf("Nome canonico do usuario: %s", canonicalValue(identity.UserName)),
		"Esses fatos canonicos tem prioridade sobre historico conversacional e placeholders.",
	}

	return strings.Join(lines, "\n")
}

func canonicalUserName(userBytes []byte) string {
	name := extractCanonicalUserName(string(userBytes))
	if name == "" {
		return "nao definido"
	}
	return name
}

func extractCanonicalUserName(userContent string) string {
	for _, line := range strings.Split(userContent, "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(strings.ToLower(trimmed), "nome:") {
			continue
		}

		name := strings.TrimSpace(strings.TrimPrefix(trimmed, "Nome:"))
		name = strings.TrimSpace(strings.Trim(name, "*"))
		if name == "" {
			return ""
		}
		if placeholderUserPattern.MatchString(name) {
			return ""
		}

		return name
	}

	return ""
}

func canonicalValue(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "nao definido"
	}
	return value
}
