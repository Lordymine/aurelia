package agent

import (
	"context"
	"fmt"
	"strings"

	"github.com/kocar/aurelia/internal/observability"
)

const (
	maxCompactedToolOutputRunes = 1200
	toolOutputPreviewRunes      = 400
)

type CompactedToolOutput struct {
	Content        string
	ArtifactID     int64
	RawChars       int
	CompactedChars int
	Oversized      bool
}

func CompactToolOutputForHistory(ctx context.Context, recorder observability.Recorder, toolName, content string) CompactedToolOutput {
	metrics := MeasureToolOutput(content)
	if !metrics.Oversized {
		return CompactedToolOutput{
			Content:        content,
			RawChars:       metrics.Chars,
			CompactedChars: metrics.Chars,
			Oversized:      false,
		}
	}

	artifactID := int64(0)
	if artifactRecorder, ok := recorder.(observability.ArtifactRecorder); ok {
		id, err := artifactRecorder.RecordArtifact(ctx, observability.Artifact{
			RunID:     ContextFields(ctx)["run_id"],
			TeamID:    ContextFields(ctx)["team_id"],
			TaskID:    ContextFields(ctx)["task_id"],
			AgentName: ContextFields(ctx)["agent"],
			Component: "agent.tool",
			Operation: toolName,
			Content:   content,
		})
		if err == nil {
			artifactID = id
		}
	}

	head := compactRunes(strings.TrimSpace(content), toolOutputPreviewRunes)
	tail := compactTailRunes(strings.TrimSpace(content), toolOutputPreviewRunes)
	lines := []string{
		fmt.Sprintf("[tool output compacted: tool=%s total_chars=%d]", toolName, metrics.Chars),
	}
	if artifactID > 0 {
		lines[0] += fmt.Sprintf(" [artifact_id=%d]", artifactID)
	}
	if head != "" {
		lines = append(lines, "head:")
		lines = append(lines, head)
	}
	if tail != "" && tail != head {
		lines = append(lines, "tail:")
		lines = append(lines, tail)
	}

	compacted := strings.Join(lines, "\n")
	compacted = compactRunes(compacted, maxCompactedToolOutputRunes)
	return CompactedToolOutput{
		Content:        compacted,
		ArtifactID:     artifactID,
		RawChars:       metrics.Chars,
		CompactedChars: len([]rune(compacted)),
		Oversized:      true,
	}
}

func compactRunes(text string, max int) string {
	runes := []rune(text)
	if len(runes) <= max {
		return text
	}
	return strings.TrimSpace(string(runes[:max])) + "..."
}

func compactTailRunes(text string, max int) string {
	runes := []rune(text)
	if len(runes) <= max {
		return text
	}
	return "..." + strings.TrimSpace(string(runes[len(runes)-max:]))
}
