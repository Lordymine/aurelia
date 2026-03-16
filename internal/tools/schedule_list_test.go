package tools

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/kocar/aurelia/internal/agent"
	cronpkg "github.com/kocar/aurelia/internal/cron"
)

type fakeScheduleLister struct {
	chatIDs []int64
	jobs    []cronpkg.CronJob
	err     error
}

func (f *fakeScheduleLister) ListJobs(ctx context.Context, chatID int64) ([]cronpkg.CronJob, error) {
	f.chatIDs = append(f.chatIDs, chatID)
	if f.err != nil {
		return nil, f.err
	}
	return f.jobs, nil
}

func TestListSchedulesTool_Execute_ListsJobsForCurrentChat(t *testing.T) {
	t.Parallel()

	runAt := time.Date(2026, 3, 12, 9, 0, 0, 0, time.UTC)
	lister := &fakeScheduleLister{
		jobs: []cronpkg.CronJob{
			{ID: "job-a", ScheduleType: "cron", CronExpr: "0 8 * * *", Prompt: "Resumo diario", Active: true, LastStatus: "idle"},
			{ID: "job-b", ScheduleType: "once", RunAt: &runAt, Prompt: "Lembrete unico", Active: false, LastStatus: "success"},
		},
	}
	tool := NewListSchedulesTool(lister)

	ctx := agent.WithTeamContext(context.Background(), "12345", "user-1")
	result, err := tool.Execute(ctx, nil)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if len(lister.chatIDs) != 1 || lister.chatIDs[0] != 12345 {
		t.Fatalf("expected one list call for chat 12345, got %#v", lister.chatIDs)
	}
	if !strings.Contains(result, "job-a") || !strings.Contains(result, "job-b") {
		t.Fatalf("expected job ids in result, got %q", result)
	}
}

func TestListSchedulesTool_Execute_ReturnsFriendlyEmptyState(t *testing.T) {
	t.Parallel()

	tool := NewListSchedulesTool(&fakeScheduleLister{})
	ctx := agent.WithTeamContext(context.Background(), "55", "user-2")

	result, err := tool.Execute(ctx, nil)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(strings.ToLower(result), "nenhum") {
		t.Fatalf("expected empty state in result, got %q", result)
	}
}


