package agent

import "testing"

func TestCompactToolsForPrompt_PrunesSchemaAndDescription(t *testing.T) {
	t.Parallel()

	tools := CompactToolsForPrompt([]Tool{{
		Name:        "read_file",
		Description: "Read a local file from disk with optional working directory context. This extra sentence should be dropped.",
		JSONSchema: map[string]interface{}{
			"type":        "object",
			"title":       "Read File Input",
			"description": "Long schema description",
			"properties": map[string]interface{}{
				"path": map[string]interface{}{
					"type":        "string",
					"description": "Path to the file",
					"default":     "README.md",
				},
			},
			"required": []string{"path"},
			"examples": []interface{}{"README.md"},
		},
	}})

	if len(tools) != 1 {
		t.Fatalf("expected 1 tool, got %d", len(tools))
	}
	if tools[0].Description != "Read a local file from disk with optional working directory context." {
		t.Fatalf("unexpected compacted description %q", tools[0].Description)
	}
	if _, ok := tools[0].JSONSchema["title"]; ok {
		t.Fatalf("expected title to be pruned")
	}
	props := tools[0].JSONSchema["properties"].(map[string]interface{})
	pathSchema := props["path"].(map[string]interface{})
	if _, ok := pathSchema["description"]; ok {
		t.Fatalf("expected property description to be pruned")
	}
	if pathSchema["type"] != "string" {
		t.Fatalf("expected property type to remain, got %#v", pathSchema)
	}
}

func TestCompactToolsForPrompt_EmptySchemaGetsObjectFallback(t *testing.T) {
	t.Parallel()

	tools := CompactToolsForPrompt([]Tool{{Name: "noop"}})
	if tools[0].JSONSchema["type"] != "object" {
		t.Fatalf("expected object fallback, got %#v", tools[0].JSONSchema)
	}
}

func TestCompactToolsForPrompt_SortsToolsByName(t *testing.T) {
	t.Parallel()

	tools := CompactToolsForPrompt([]Tool{
		{Name: "write_file"},
		{Name: "read_file"},
		{Name: "list_dir"},
	})

	if tools[0].Name != "list_dir" || tools[1].Name != "read_file" || tools[2].Name != "write_file" {
		t.Fatalf("expected stable sorted order, got %#v", tools)
	}
}
