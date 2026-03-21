package agents

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Registry holds all loaded agent definitions indexed by name.
type Registry struct {
	agents map[string]*Agent
}

// Load reads all .md files from dir and returns a Registry.
// Files without valid frontmatter (missing --- markers or missing name) are skipped.
func Load(dir string) (*Registry, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading agents dir: %w", err)
	}

	reg := &Registry{agents: make(map[string]*Agent)}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("reading agent file %s: %w", entry.Name(), err)
		}

		agent, err := parseAgentFile(data)
		if err != nil {
			return nil, fmt.Errorf("parsing agent file %s: %w", entry.Name(), err)
		}
		if agent == nil {
			// No valid frontmatter, skip.
			continue
		}

		reg.agents[agent.Name] = agent
	}

	return reg, nil
}

// Get returns the agent with the given name, or nil if not found.
func (r *Registry) Get(name string) *Agent {
	return r.agents[name]
}

// Agents returns all loaded agents sorted by name.
func (r *Registry) Agents() []*Agent {
	result := make([]*Agent, 0, len(r.agents))
	for _, a := range r.agents {
		result = append(result, a)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result
}

// Scheduled returns agents that have a schedule defined.
func (r *Registry) Scheduled() []*Agent {
	var result []*Agent
	for _, a := range r.agents {
		if a.Schedule != "" {
			result = append(result, a)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result
}

// Route checks if message starts with @agentname (case-insensitive)
// and returns the matching agent, or nil if no match.
func (r *Registry) Route(message string) *Agent {
	if !strings.HasPrefix(message, "@") {
		return nil
	}

	// Extract the name after @, before the first space.
	rest := message[1:]
	name := rest
	if idx := strings.IndexByte(rest, ' '); idx != -1 {
		name = rest[:idx]
	}
	name = strings.ToLower(name)

	return r.agents[name]
}

// ClassifyPrompt builds a prompt that asks an LLM to pick the best agent for a message.
// Returns empty string if no agents are loaded.
func (r *Registry) ClassifyPrompt(message string) string {
	if len(r.agents) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("Given these available agents:\n\n")
	for _, a := range r.Agents() {
		fmt.Fprintf(&sb, "- %s: %s\n", a.Name, a.Description)
	}
	fmt.Fprintf(&sb, "\nUser message: %q\n\n", message)
	sb.WriteString("Reply with ONLY the agent name that best matches, or 'none' if no agent is a good match. Reply with a single word.")
	return sb.String()
}

// parseAgentFile splits a markdown file on --- markers, parses YAML frontmatter,
// and extracts the prompt body. Returns nil if the file has no valid frontmatter
// or the name field is empty.
func parseAgentFile(data []byte) (*Agent, error) {
	parts := bytes.SplitN(data, []byte("---"), 3)
	if len(parts) != 3 {
		return nil, nil
	}

	var agent Agent
	if err := yaml.Unmarshal(parts[1], &agent); err != nil {
		return nil, fmt.Errorf("parsing yaml frontmatter: %w", err)
	}

	if agent.Name == "" {
		return nil, nil
	}

	agent.Prompt = string(bytes.TrimSpace(parts[2]))
	return &agent, nil
}
