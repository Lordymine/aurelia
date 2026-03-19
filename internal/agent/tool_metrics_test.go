package agent

import "testing"

func TestMeasureToolPayload(t *testing.T) {
	t.Parallel()

	metrics := MeasureToolPayload([]Tool{
		{
			Name:        "read_file",
			Description: "read a file",
			JSONSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]interface{}{"type": "string"},
				},
				"required": []string{"path"},
			},
		},
	})

	if metrics.Count != 1 {
		t.Fatalf("expected 1 tool, got %d", metrics.Count)
	}
	if metrics.SerializedBytes == 0 {
		t.Fatalf("expected serialized bytes to be > 0")
	}
}

func TestMeasureToolOutput(t *testing.T) {
	t.Parallel()

	small := MeasureToolOutput("ok")
	if small.Chars != 2 {
		t.Fatalf("expected 2 chars, got %d", small.Chars)
	}
	if small.Oversized {
		t.Fatalf("expected small output not to be oversized")
	}

	large := MeasureToolOutput(string(make([]rune, OversizedToolOutputThresholdChars+1)))
	if !large.Oversized {
		t.Fatalf("expected oversized output")
	}
}
