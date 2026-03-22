package agents

import "testing"

func TestBuildSDKAgents_Empty(t *testing.T) {
	r := &Registry{agents: map[string]*Agent{}}
	got := BuildSDKAgents(r)
	if len(got) != 0 {
		t.Errorf("expected empty map, got %d entries", len(got))
	}
}

func TestBuildSDKAgents_SingleAgent(t *testing.T) {
	r := &Registry{agents: map[string]*Agent{
		"coder": {
			Name:         "coder",
			Description:  "writes code",
			Model:        "claude-sonnet-4-6",
			Prompt:       "You are a coder.",
			AllowedTools: []string{"Read", "Edit"},
		},
	}}
	got := BuildSDKAgents(r)
	if len(got) != 1 {
		t.Fatalf("expected 1 agent, got %d", len(got))
	}
	a, ok := got["coder"]
	if !ok {
		t.Fatal("missing 'coder' key")
	}
	m := a.(map[string]any)
	if m["description"] != "writes code" {
		t.Errorf("description = %q", m["description"])
	}
	if m["prompt"] != "You are a coder." {
		t.Errorf("prompt = %q", m["prompt"])
	}
	if m["model"] != "claude-sonnet-4-6" {
		t.Errorf("model = %q", m["model"])
	}
	tools, ok := m["tools"]
	if !ok {
		t.Fatal("missing 'tools' key")
	}
	toolSlice := tools.([]string)
	if len(toolSlice) != 2 || toolSlice[0] != "Read" {
		t.Errorf("tools = %v", toolSlice)
	}
}

func TestBuildSDKAgents_OmitsEmptyFields(t *testing.T) {
	r := &Registry{agents: map[string]*Agent{
		"simple": {
			Name:        "simple",
			Description: "a simple agent",
			Prompt:      "Be helpful.",
		},
	}}
	got := BuildSDKAgents(r)
	m := got["simple"].(map[string]any)
	if _, ok := m["model"]; ok {
		t.Error("model should be omitted when empty")
	}
	if _, ok := m["tools"]; ok {
		t.Error("tools should be omitted when empty")
	}
}
