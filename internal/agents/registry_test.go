package agents

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func TestRegistry_LoadAgents(t *testing.T) {
	dir := t.TempDir()

	content := `---
name: prospector
description: Busca leads
model: kimi-k2-thinking
schedule: "0 9 * * 1"
allowed_tools: ["WebSearch", "WebFetch"]
---

Voce eh um agente de prospeccao.
`
	if err := os.WriteFile(filepath.Join(dir, "prospector.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	reg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	agents := reg.Agents()
	if len(agents) != 1 {
		t.Fatalf("expected 1 agent, got %d", len(agents))
	}

	a := reg.Get("prospector")
	if a == nil {
		t.Fatal("Get('prospector') returned nil")
	}

	if a.Name != "prospector" {
		t.Errorf("expected name 'prospector', got %q", a.Name)
	}
	if a.Description != "Busca leads" {
		t.Errorf("expected description 'Busca leads', got %q", a.Description)
	}
	if a.Model != "kimi-k2-thinking" {
		t.Errorf("expected model 'kimi-k2-thinking', got %q", a.Model)
	}
	if a.Schedule != "0 9 * * 1" {
		t.Errorf("expected schedule '0 9 * * 1', got %q", a.Schedule)
	}
	if len(a.AllowedTools) != 2 || a.AllowedTools[0] != "WebSearch" || a.AllowedTools[1] != "WebFetch" {
		t.Errorf("expected allowed_tools [WebSearch, WebFetch], got %v", a.AllowedTools)
	}
	if a.Prompt != "Voce eh um agente de prospeccao." {
		t.Errorf("expected prompt body trimmed, got %q", a.Prompt)
	}
}

func TestRegistry_LoadMultipleAgents(t *testing.T) {
	dir := t.TempDir()

	writeAgent(t, dir, "alpha.md", `---
name: alpha
description: Agent Alpha
---
Alpha prompt.
`)
	writeAgent(t, dir, "beta.md", `---
name: beta
description: Agent Beta
model: gpt-4o
---
Beta prompt.
`)

	reg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	agents := reg.Agents()
	if len(agents) != 2 {
		t.Fatalf("expected 2 agents, got %d", len(agents))
	}

	if reg.Get("alpha") == nil {
		t.Error("Get('alpha') returned nil")
	}
	if reg.Get("beta") == nil {
		t.Error("Get('beta') returned nil")
	}
	if reg.Get("nonexistent") != nil {
		t.Error("Get('nonexistent') should return nil")
	}
}

func TestRegistry_Route(t *testing.T) {
	dir := t.TempDir()

	writeAgent(t, dir, "prospector.md", `---
name: prospector
description: Busca leads
---
Prompt.
`)
	writeAgent(t, dir, "writer.md", `---
name: writer
description: Writes content
---
Prompt.
`)

	reg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	tests := []struct {
		message string
		want    string // empty means nil
	}{
		{"@prospector find leads for ACME", "prospector"},
		{"@Prospector find leads", "prospector"},        // case-insensitive
		{"@PROSPECTOR find leads", "prospector"},        // all caps
		{"@writer create a blog post", "writer"},
		{"hello world", ""},                              // no @ prefix
		{"@unknown do something", ""},                    // unknown agent
		{"email me @prospector", ""},                     // @ not at start
	}

	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			got := reg.Route(tt.message)
			if tt.want == "" {
				if got != nil {
					t.Errorf("expected nil, got agent %q", got.Name)
				}
			} else {
				if got == nil {
					t.Fatalf("expected agent %q, got nil", tt.want)
				}
				if got.Name != tt.want {
					t.Errorf("expected agent %q, got %q", tt.want, got.Name)
				}
			}
		})
	}
}

func TestRegistry_Scheduled(t *testing.T) {
	dir := t.TempDir()

	writeAgent(t, dir, "scheduled.md", `---
name: scheduled
description: Has schedule
schedule: "0 9 * * 1"
---
Prompt.
`)
	writeAgent(t, dir, "ondemand.md", `---
name: ondemand
description: No schedule
---
Prompt.
`)

	reg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	scheduled := reg.Scheduled()
	if len(scheduled) != 1 {
		t.Fatalf("expected 1 scheduled agent, got %d", len(scheduled))
	}
	if scheduled[0].Name != "scheduled" {
		t.Errorf("expected 'scheduled', got %q", scheduled[0].Name)
	}
}

func TestRegistry_EmptyDir(t *testing.T) {
	dir := t.TempDir()

	reg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if len(reg.Agents()) != 0 {
		t.Errorf("expected 0 agents, got %d", len(reg.Agents()))
	}
}

func TestRegistry_MalformedMarkdown(t *testing.T) {
	dir := t.TempDir()

	// File without frontmatter — should be skipped gracefully
	writeAgent(t, dir, "nofrontmatter.md", `Just some text without frontmatter.`)

	// File with valid frontmatter
	writeAgent(t, dir, "valid.md", `---
name: valid
description: Valid agent
---
Valid prompt.
`)

	reg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	agents := reg.Agents()
	if len(agents) != 1 {
		t.Fatalf("expected 1 agent (skipping malformed), got %d", len(agents))
	}
	if agents[0].Name != "valid" {
		t.Errorf("expected 'valid', got %q", agents[0].Name)
	}
}

func TestRegistry_IgnoresNonMdFiles(t *testing.T) {
	dir := t.TempDir()

	writeAgent(t, dir, "agent.md", `---
name: agent
description: An agent
---
Prompt.
`)
	// Write a non-md file that should be ignored
	if err := os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("not an agent"), 0644); err != nil {
		t.Fatal(err)
	}

	reg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if len(reg.Agents()) != 1 {
		t.Errorf("expected 1 agent, got %d", len(reg.Agents()))
	}
}

func TestRegistry_MCPServers(t *testing.T) {
	dir := t.TempDir()

	writeAgent(t, dir, "devops.md", `---
name: devops
description: DevOps agent
mcp_servers:
  docker:
    command: docker
    args: ["mcp", "gateway"]
---
DevOps prompt.
`)

	reg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	a := reg.Get("devops")
	if a == nil {
		t.Fatal("Get('devops') returned nil")
	}
	if a.MCPServers == nil {
		t.Fatal("expected mcp_servers to be populated")
	}
	if _, ok := a.MCPServers["docker"]; !ok {
		t.Error("expected 'docker' key in mcp_servers")
	}
}

func TestRegistry_AgentsReturnsSorted(t *testing.T) {
	dir := t.TempDir()

	writeAgent(t, dir, "zeta.md", `---
name: zeta
description: Zeta
---
Prompt.
`)
	writeAgent(t, dir, "alpha.md", `---
name: alpha
description: Alpha
---
Prompt.
`)

	reg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	agents := reg.Agents()
	names := make([]string, len(agents))
	for i, a := range agents {
		names[i] = a.Name
	}

	if !sort.StringsAreSorted(names) {
		t.Errorf("expected agents sorted by name, got %v", names)
	}
}

func TestRegistry_ClassifyPrompt(t *testing.T) {
	dir := t.TempDir()

	writeAgent(t, dir, "prospector.md", `---
name: prospector
description: Busca leads e prospecta clientes
---
Prompt.
`)

	reg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	prompt := reg.ClassifyPrompt("find new leads for ACME Corp")
	if prompt == "" {
		t.Fatal("expected non-empty classify prompt")
	}
	if !strings.Contains(prompt, "prospector") {
		t.Errorf("expected prompt to contain agent name 'prospector', got %q", prompt)
	}
	if !strings.Contains(prompt, "Busca leads") {
		t.Errorf("expected prompt to contain agent description, got %q", prompt)
	}
	if !strings.Contains(prompt, "find new leads for ACME Corp") {
		t.Errorf("expected prompt to contain user message, got %q", prompt)
	}
}

func TestRegistry_ClassifyPrompt_Empty(t *testing.T) {
	dir := t.TempDir()

	reg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	prompt := reg.ClassifyPrompt("hello")
	if prompt != "" {
		t.Errorf("expected empty prompt for empty registry, got %q", prompt)
	}
}

func TestLoadNormalizesKeys(t *testing.T) {
	dir := t.TempDir()

	// Agent file with mixed-case name in YAML
	writeAgent(t, dir, "myagent.md", `---
name: MyAgent
description: test
---
Prompt.
`)

	reg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	// Get with exact case should work
	if reg.Get("MyAgent") == nil {
		t.Fatal("expected Get('MyAgent') to find agent")
	}
	// Get with lowercase should also work
	if reg.Get("myagent") == nil {
		t.Fatal("expected Get('myagent') to find agent")
	}
	// Get with uppercase should also work
	if reg.Get("MYAGENT") == nil {
		t.Fatal("expected Get('MYAGENT') to find agent")
	}
}

func TestRouteCaseInsensitive(t *testing.T) {
	dir := t.TempDir()

	writeAgent(t, dir, "myagent.md", `---
name: MyAgent
description: test
---
Prompt.
`)

	reg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	agent := reg.Route("@myagent hello")
	if agent == nil {
		t.Fatal("expected agent to be found with case-insensitive lookup")
	}
	if agent.Name != "MyAgent" {
		t.Errorf("expected Name to preserve original case 'MyAgent', got %q", agent.Name)
	}

	agent = reg.Route("@MYAGENT hello")
	if agent == nil {
		t.Fatal("expected agent to be found with uppercase lookup")
	}
}

func TestClassify(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "coder.md"), []byte("---\nname: coder\ndescription: writes code\n---\nYou write code."), 0644)
	reg, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}

	mockClassify := func(ctx context.Context, system, prompt string) (string, error) {
		return "coder", nil
	}
	agent := reg.Classify(context.Background(), "write me a function", mockClassify)
	if agent == nil || agent.Name != "coder" {
		t.Fatalf("expected coder agent, got %v", agent)
	}
}

func TestClassifyNone(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "coder.md"), []byte("---\nname: coder\ndescription: writes code\n---\nYou write code."), 0644)
	reg, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}

	mockClassify := func(ctx context.Context, system, prompt string) (string, error) {
		return "none", nil
	}
	agent := reg.Classify(context.Background(), "hello", mockClassify)
	if agent != nil {
		t.Fatal("expected nil for 'none' classification")
	}
}

func writeAgent(t *testing.T, dir, filename, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, filename), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}
