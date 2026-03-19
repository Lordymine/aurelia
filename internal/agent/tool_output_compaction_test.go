package agent

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"github.com/kocar/aurelia/internal/observability"
)

type artifactRecorderStub struct {
	recordingObserver
	artifacts []observability.Artifact
}

func (a *artifactRecorderStub) RecordArtifact(ctx context.Context, artifact observability.Artifact) (int64, error) {
	a.artifacts = append(a.artifacts, artifact)
	return int64(len(a.artifacts)), nil
}

func TestCompactToolOutputForHistory_LeavesSmallOutputUntouched(t *testing.T) {
	t.Parallel()

	got := CompactToolOutputForHistory(context.Background(), nil, "read_file", "ok")
	if got.Content != "ok" {
		t.Fatalf("expected untouched output, got %q", got.Content)
	}
}

func TestCompactToolOutputForHistory_CompactsOversizedOutputAndRecordsArtifact(t *testing.T) {
	t.Parallel()

	recorder := &artifactRecorderStub{}
	ctx := WithRunContext(context.Background(), "run-1")
	ctx = WithTaskContext(ctx, "team-1", "task-1")
	ctx = WithAgentContext(ctx, "agent-1")
	raw := strings.Repeat("a", OversizedToolOutputThresholdChars+100)

	got := CompactToolOutputForHistory(ctx, recorder, "read_file", raw)
	if got.Content == raw {
		t.Fatalf("expected output to be compacted")
	}
	if !strings.Contains(got.Content, "tool output compacted") {
		t.Fatalf("expected compaction marker, got %q", got.Content)
	}
	if !strings.Contains(got.Content, "artifact_id=1") {
		t.Fatalf("expected artifact id in compacted output, got %q", got.Content)
	}
	if len(recorder.artifacts) != 1 {
		t.Fatalf("expected 1 recorded artifact, got %d", len(recorder.artifacts))
	}
	if recorder.artifacts[0].Content != raw {
		t.Fatalf("expected raw content to be preserved")
	}
	if got.RawChars <= got.CompactedChars {
		t.Fatalf("expected compaction to reduce size, got raw=%d compacted=%d", got.RawChars, got.CompactedChars)
	}
}

func TestCompactToolOutputForHistory_CompactedOutputStaysBounded(t *testing.T) {
	t.Parallel()

	got := CompactToolOutputForHistory(context.Background(), nil, "run_command", strings.Repeat("x", OversizedToolOutputThresholdChars+500))
	if len([]rune(got.Content)) > maxCompactedToolOutputRunes+3 {
		t.Fatalf("expected bounded compacted output, got %d runes", len([]rune(got.Content)))
	}
	if !strings.Contains(got.Content, strconv.Itoa(OversizedToolOutputThresholdChars+500)) {
		t.Fatalf("expected size marker in compacted output, got %q", got.Content)
	}
}
