package agent

import (
	"context"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"
)

func newTestTeamManager(t *testing.T) TeamManager {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "team-manager.db")

	store, err := NewSQLiteTaskStore(dbPath)
	if err != nil {
		t.Fatalf("failed to create sqlite task store: %v", err)
	}
	t.Cleanup(func() {
		_ = store.Close()
	})

	manager, err := NewTeamManager(store)
	if err != nil {
		t.Fatalf("failed to create team manager: %v", err)
	}

	return manager
}

func TestTeamManager_ClaimNextTask_IsAtomic(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := newTestTeamManager(t)

	teamID, err := manager.CreateTeam(ctx, "team-user-1", "user-1", "master")
	if err != nil {
		t.Fatalf("CreateTeam() error = %v", err)
	}

	if err := manager.RegisterTeammate(ctx, teamID, "worker-a", "Implements code"); err != nil {
		t.Fatalf("RegisterTeammate(worker-a) error = %v", err)
	}
	if err := manager.RegisterTeammate(ctx, teamID, "worker-b", "Reviews code"); err != nil {
		t.Fatalf("RegisterTeammate(worker-b) error = %v", err)
	}

	task := TeamTask{
		ID:     "task-root",
		TeamID: teamID,
		Title:  "Root Task",
		Prompt: "Do the first runnable task",
		Status: TaskPending,
	}

	if err := manager.CreateTask(ctx, task, nil); err != nil {
		t.Fatalf("CreateTask() error = %v", err)
	}

	type claimResult struct {
		task *TeamTask
		err  error
	}

	results := make(chan claimResult, 2)
	var wg sync.WaitGroup
	for _, agentName := range []string{"worker-a", "worker-b"} {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			task, err := manager.ClaimNextTask(ctx, teamID, name)
			results <- claimResult{task: task, err: err}
		}(agentName)
	}

	wg.Wait()
	close(results)

	var claimed []*TeamTask
	for result := range results {
		if result.err != nil {
			t.Fatalf("ClaimNextTask() error = %v", result.err)
		}
		if result.task != nil {
			claimed = append(claimed, result.task)
		}
	}

	if len(claimed) != 1 {
		t.Fatalf("expected exactly one successful claimant, got %d", len(claimed))
	}
	if claimed[0].ID != "task-root" {
		t.Fatalf("expected task-root to be claimed, got %q", claimed[0].ID)
	}
	if claimed[0].AssignedAgent == nil || *claimed[0].AssignedAgent == "" {
		t.Fatalf("expected claimed task to record assigned agent")
	}
}

func TestTeamManager_Dependencies_BlockUntilPrerequisiteCompletes(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := newTestTeamManager(t)

	teamID, err := manager.CreateTeam(ctx, "team-user-2", "user-2", "master")
	if err != nil {
		t.Fatalf("CreateTeam() error = %v", err)
	}

	if err := manager.RegisterTeammate(ctx, teamID, "worker-a", "Planner"); err != nil {
		t.Fatalf("RegisterTeammate() error = %v", err)
	}

	taskA := TeamTask{
		ID:     "task-a",
		TeamID: teamID,
		Title:  "Task A",
		Prompt: "Finish first",
		Status: TaskPending,
	}
	taskB := TeamTask{
		ID:     "task-b",
		TeamID: teamID,
		Title:  "Task B",
		Prompt: "Depends on task A",
		Status: TaskPending,
	}

	if err := manager.CreateTask(ctx, taskA, nil); err != nil {
		t.Fatalf("CreateTask(taskA) error = %v", err)
	}
	if err := manager.CreateTask(ctx, taskB, []string{"task-a"}); err != nil {
		t.Fatalf("CreateTask(taskB) error = %v", err)
	}

	firstClaim, err := manager.ClaimNextTask(ctx, teamID, "worker-a")
	if err != nil {
		t.Fatalf("ClaimNextTask(first) error = %v", err)
	}
	if firstClaim == nil || firstClaim.ID != "task-a" {
		t.Fatalf("expected first runnable task to be task-a, got %#v", firstClaim)
	}

	noClaimYet, err := manager.ClaimNextTask(ctx, teamID, "worker-a")
	if err != nil {
		t.Fatalf("ClaimNextTask(second before completion) error = %v", err)
	}
	if noClaimYet != nil {
		t.Fatalf("expected task-b to stay blocked until task-a completes, got %#v", noClaimYet)
	}

	if err := manager.CompleteTask(ctx, teamID, "task-a", "worker-a", "done"); err != nil {
		t.Fatalf("CompleteTask(task-a) error = %v", err)
	}

	secondClaim, err := manager.ClaimNextTask(ctx, teamID, "worker-a")
	if err != nil {
		t.Fatalf("ClaimNextTask(after completion) error = %v", err)
	}
	if secondClaim == nil || secondClaim.ID != "task-b" {
		t.Fatalf("expected task-b to unlock after task-a completion, got %#v", secondClaim)
	}
}

func TestTeamManager_Mailbox_DeliversAndConsumesMessages(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := newTestTeamManager(t)

	teamID, err := manager.CreateTeam(ctx, "team-user-3", "user-3", "master")
	if err != nil {
		t.Fatalf("CreateTeam() error = %v", err)
	}

	msg := MailMessage{
		ID:        "mail-1",
		TeamID:    teamID,
		FromAgent: "worker-a",
		ToAgent:   "master",
		Kind:      "result",
		Body:      "task completed with findings",
	}

	if err := manager.PostMessage(ctx, msg); err != nil {
		t.Fatalf("PostMessage() error = %v", err)
	}

	firstPull, err := manager.PullMessages(ctx, teamID, "master", 10)
	if err != nil {
		t.Fatalf("PullMessages(first) error = %v", err)
	}
	if len(firstPull) != 1 {
		t.Fatalf("expected one message in first pull, got %d", len(firstPull))
	}
	if firstPull[0].Body != "task completed with findings" {
		t.Fatalf("unexpected mailbox body: %q", firstPull[0].Body)
	}
	if firstPull[0].ConsumedAt == nil {
		t.Fatalf("expected pulled message to be marked consumed")
	}

	secondPull, err := manager.PullMessages(ctx, teamID, "master", 10)
	if err != nil {
		t.Fatalf("PullMessages(second) error = %v", err)
	}
	if len(secondPull) != 0 {
		t.Fatalf("expected consumed messages to not be redelivered, got %d", len(secondPull))
	}
}

func TestTeamManager_CompleteTask_AllowsWorkerToPickNextRunnableTask(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := newTestTeamManager(t)

	teamID, err := manager.CreateTeam(ctx, "team-user-4", "user-4", "master")
	if err != nil {
		t.Fatalf("CreateTeam() error = %v", err)
	}

	if err := manager.RegisterTeammate(ctx, teamID, "worker-a", "Executes tasks sequentially"); err != nil {
		t.Fatalf("RegisterTeammate() error = %v", err)
	}

	if err := manager.CreateTask(ctx, TeamTask{
		ID:     "task-1",
		TeamID: teamID,
		Title:  "Prepare data",
		Prompt: "First task",
		Status: TaskPending,
	}, nil); err != nil {
		t.Fatalf("CreateTask(task-1) error = %v", err)
	}

	if err := manager.CreateTask(ctx, TeamTask{
		ID:     "task-2",
		TeamID: teamID,
		Title:  "Analyze data",
		Prompt: "Second task",
		Status: TaskPending,
	}, []string{"task-1"}); err != nil {
		t.Fatalf("CreateTask(task-2) error = %v", err)
	}

	claimedTask, err := manager.ClaimNextTask(ctx, teamID, "worker-a")
	if err != nil {
		t.Fatalf("ClaimNextTask(task-1) error = %v", err)
	}
	if claimedTask == nil || claimedTask.ID != "task-1" {
		t.Fatalf("expected worker to claim task-1 first, got %#v", claimedTask)
	}

	if err := manager.CompleteTask(ctx, teamID, "task-1", "worker-a", "data prepared"); err != nil {
		t.Fatalf("CompleteTask(task-1) error = %v", err)
	}

	nextTask, err := manager.ClaimNextTask(ctx, teamID, "worker-a")
	if err != nil {
		t.Fatalf("ClaimNextTask(task-2) error = %v", err)
	}
	if nextTask == nil || nextTask.ID != "task-2" {
		t.Fatalf("expected worker to automatically become eligible for task-2, got %#v", nextTask)
	}
}

func TestTeamManager_TaskEvents_RecordLifecycleInOrder(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := newTestTeamManager(t)

	teamID, err := manager.CreateTeam(ctx, "team-user-5", "user-5", "master")
	if err != nil {
		t.Fatalf("CreateTeam() error = %v", err)
	}

	if err := manager.RegisterTeammate(ctx, teamID, "worker-a", "Executes lifecycle"); err != nil {
		t.Fatalf("RegisterTeammate() error = %v", err)
	}

	if err := manager.CreateTask(ctx, TeamTask{
		ID:     "task-events-1",
		TeamID: teamID,
		Title:  "Lifecycle task",
		Prompt: "execute lifecycle",
		Status: TaskPending,
	}, nil); err != nil {
		t.Fatalf("CreateTask() error = %v", err)
	}

	claimed, err := manager.ClaimNextTask(ctx, teamID, "worker-a")
	if err != nil {
		t.Fatalf("ClaimNextTask() error = %v", err)
	}
	if claimed == nil {
		t.Fatalf("expected task to be claimed")
	}

	if err := manager.CompleteTask(ctx, teamID, claimed.ID, "worker-a", "done"); err != nil {
		t.Fatalf("CompleteTask() error = %v", err)
	}

	events, err := manager.ListEvents(ctx, teamID, 10)
	if err != nil {
		t.Fatalf("ListEvents() error = %v", err)
	}
	if len(events) != 3 {
		t.Fatalf("expected 3 task lifecycle events, got %d", len(events))
	}

	gotTypes := []string{events[0].EventType, events[1].EventType, events[2].EventType}
	wantTypes := []string{"task_created", "task_claimed", "task_completed"}
	if !slices.Equal(gotTypes, wantTypes) {
		t.Fatalf("unexpected event order: got %v want %v", gotTypes, wantTypes)
	}

	if events[1].AgentName != "worker-a" {
		t.Fatalf("expected claim event agent to be worker-a, got %q", events[1].AgentName)
	}
	if events[2].TaskID == nil || *events[2].TaskID != "task-events-1" {
		t.Fatalf("expected complete event to reference task-events-1, got %#v", events[2].TaskID)
	}
}

func TestTeamManager_CreateTask_WithDependenciesStartsBlocked(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := newTestTeamManager(t)

	teamID, err := manager.CreateTeam(ctx, "team-user-6", "user-6", "master")
	if err != nil {
		t.Fatalf("CreateTeam() error = %v", err)
	}

	if err := manager.CreateTask(ctx, TeamTask{
		ID:     "task-root",
		TeamID: teamID,
		Title:  "Root",
		Prompt: "Do root work",
		Status: TaskPending,
	}, nil); err != nil {
		t.Fatalf("CreateTask(root) error = %v", err)
	}

	if err := manager.CreateTask(ctx, TeamTask{
		ID:     "task-child",
		TeamID: teamID,
		Title:  "Child",
		Prompt: "Wait for root",
		Status: TaskPending,
	}, []string{"task-root"}); err != nil {
		t.Fatalf("CreateTask(child) error = %v", err)
	}

	task, err := manager.GetTask(ctx, teamID, "task-child")
	if err != nil {
		t.Fatalf("GetTask(child) error = %v", err)
	}
	if task == nil {
		t.Fatal("expected child task to exist")
	}
	if task.Status != TaskBlocked {
		t.Fatalf("expected dependent task to start blocked, got %q", task.Status)
	}
	if !strings.Contains(strings.ToLower(task.ErrorMessage), "dependency") {
		t.Fatalf("expected blocked reason to mention dependency, got %q", task.ErrorMessage)
	}

	events, err := manager.ListEvents(ctx, teamID, 10)
	if err != nil {
		t.Fatalf("ListEvents() error = %v", err)
	}
	if len(events) < 3 {
		t.Fatalf("expected blocked lifecycle events, got %d", len(events))
	}
	if events[2].EventType != "task_blocked" {
		t.Fatalf("expected third event to be task_blocked, got %q", events[2].EventType)
	}
}

func TestTeamManager_CompleteTask_UnblocksDependents(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := newTestTeamManager(t)

	teamID, err := manager.CreateTeam(ctx, "team-user-7", "user-7", "master")
	if err != nil {
		t.Fatalf("CreateTeam() error = %v", err)
	}

	if err := manager.RegisterTeammate(ctx, teamID, "worker-a", "Executes"); err != nil {
		t.Fatalf("RegisterTeammate() error = %v", err)
	}

	if err := manager.CreateTask(ctx, TeamTask{
		ID:     "task-a",
		TeamID: teamID,
		Title:  "Task A",
		Prompt: "Finish A",
		Status: TaskPending,
	}, nil); err != nil {
		t.Fatalf("CreateTask(task-a) error = %v", err)
	}

	if err := manager.CreateTask(ctx, TeamTask{
		ID:     "task-b",
		TeamID: teamID,
		Title:  "Task B",
		Prompt: "Depends on A",
		Status: TaskPending,
	}, []string{"task-a"}); err != nil {
		t.Fatalf("CreateTask(task-b) error = %v", err)
	}

	claimed, err := manager.ClaimNextTask(ctx, teamID, "worker-a")
	if err != nil {
		t.Fatalf("ClaimNextTask(task-a) error = %v", err)
	}
	if claimed == nil || claimed.ID != "task-a" {
		t.Fatalf("expected task-a claim, got %#v", claimed)
	}

	if err := manager.CompleteTask(ctx, teamID, "task-a", "worker-a", "done"); err != nil {
		t.Fatalf("CompleteTask(task-a) error = %v", err)
	}

	task, err := manager.GetTask(ctx, teamID, "task-b")
	if err != nil {
		t.Fatalf("GetTask(task-b) error = %v", err)
	}
	if task == nil {
		t.Fatal("expected task-b to exist")
	}
	if task.Status != TaskPending {
		t.Fatalf("expected task-b to be unblocked into pending, got %q", task.Status)
	}
	if task.ErrorMessage != "" {
		t.Fatalf("expected unblock to clear reason, got %q", task.ErrorMessage)
	}

	events, err := manager.ListEvents(ctx, teamID, 10)
	if err != nil {
		t.Fatalf("ListEvents() error = %v", err)
	}
	var sawUnblocked bool
	for _, event := range events {
		if event.TaskID != nil && *event.TaskID == "task-b" && event.EventType == "task_unblocked" {
			sawUnblocked = true
			break
		}
	}
	if !sawUnblocked {
		t.Fatal("expected dependent task to emit task_unblocked event")
	}
}

func TestTeamManager_FailTask_CancelsBlockedDependents(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := newTestTeamManager(t)

	teamID, err := manager.CreateTeam(ctx, "team-user-8", "user-8", "master")
	if err != nil {
		t.Fatalf("CreateTeam() error = %v", err)
	}

	if err := manager.RegisterTeammate(ctx, teamID, "worker-a", "Executes"); err != nil {
		t.Fatalf("RegisterTeammate() error = %v", err)
	}

	if err := manager.CreateTask(ctx, TeamTask{
		ID:     "task-a",
		TeamID: teamID,
		Title:  "Task A",
		Prompt: "Finish A",
		Status: TaskPending,
	}, nil); err != nil {
		t.Fatalf("CreateTask(task-a) error = %v", err)
	}

	if err := manager.CreateTask(ctx, TeamTask{
		ID:     "task-b",
		TeamID: teamID,
		Title:  "Task B",
		Prompt: "Depends on A",
		Status: TaskPending,
	}, []string{"task-a"}); err != nil {
		t.Fatalf("CreateTask(task-b) error = %v", err)
	}

	claimed, err := manager.ClaimNextTask(ctx, teamID, "worker-a")
	if err != nil {
		t.Fatalf("ClaimNextTask(task-a) error = %v", err)
	}
	if claimed == nil || claimed.ID != "task-a" {
		t.Fatalf("expected task-a claim, got %#v", claimed)
	}

	if err := manager.FailTask(ctx, teamID, "task-a", "worker-a", "boom"); err != nil {
		t.Fatalf("FailTask(task-a) error = %v", err)
	}

	task, err := manager.GetTask(ctx, teamID, "task-b")
	if err != nil {
		t.Fatalf("GetTask(task-b) error = %v", err)
	}
	if task == nil {
		t.Fatal("expected task-b to exist")
	}
	if task.Status != TaskCancelled {
		t.Fatalf("expected dependent task to be cancelled, got %q", task.Status)
	}
	if !strings.Contains(strings.ToLower(task.ErrorMessage), "task-a") {
		t.Fatalf("expected cancellation reason to mention upstream task, got %q", task.ErrorMessage)
	}

	events, err := manager.ListEvents(ctx, teamID, 10)
	if err != nil {
		t.Fatalf("ListEvents() error = %v", err)
	}
	var sawCancelled bool
	for _, event := range events {
		if event.TaskID != nil && *event.TaskID == "task-b" && event.EventType == "task_cancelled" {
			sawCancelled = true
			break
		}
	}
	if !sawCancelled {
		t.Fatal("expected dependent task cancellation event")
	}
}

func TestTeamManager_CreateRecoveryTask_ReopensCancelledDependentsAsBlocked(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := newTestTeamManager(t)

	teamID, err := manager.CreateTeam(ctx, "team-user-9", "user-9", "master")
	if err != nil {
		t.Fatalf("CreateTeam() error = %v", err)
	}

	if err := manager.RegisterTeammate(ctx, teamID, "worker-a", "Executes"); err != nil {
		t.Fatalf("RegisterTeammate() error = %v", err)
	}

	if err := manager.CreateTask(ctx, TeamTask{
		ID:     "task-a",
		TeamID: teamID,
		Title:  "Task A",
		Prompt: "Finish A",
		Status: TaskPending,
	}, nil); err != nil {
		t.Fatalf("CreateTask(task-a) error = %v", err)
	}
	if err := manager.CreateTask(ctx, TeamTask{
		ID:     "task-b",
		TeamID: teamID,
		Title:  "Task B",
		Prompt: "Depends on A",
		Status: TaskPending,
	}, []string{"task-a"}); err != nil {
		t.Fatalf("CreateTask(task-b) error = %v", err)
	}

	claimed, err := manager.ClaimNextTask(ctx, teamID, "worker-a")
	if err != nil {
		t.Fatalf("ClaimNextTask(task-a) error = %v", err)
	}
	if claimed == nil {
		t.Fatal("expected task-a claim")
	}

	if err := manager.FailTask(ctx, teamID, "task-a", "worker-a", "boom"); err != nil {
		t.Fatalf("FailTask(task-a) error = %v", err)
	}

	recoveryParent := "task-a"
	assignedAgent := "worker-a"
	if err := manager.CreateTask(ctx, TeamTask{
		ID:            "recovery-a",
		TeamID:        teamID,
		ParentTaskID:  &recoveryParent,
		Title:         "recovery:task-a",
		Prompt:        "Recover A",
		AssignedAgent: &assignedAgent,
		Status:        TaskPending,
	}, nil); err != nil {
		t.Fatalf("CreateTask(recovery-a) error = %v", err)
	}

	task, err := manager.GetTask(ctx, teamID, "task-b")
	if err != nil {
		t.Fatalf("GetTask(task-b) error = %v", err)
	}
	if task == nil {
		t.Fatal("expected task-b to exist")
	}
	if task.Status != TaskBlocked {
		t.Fatalf("expected dependent to reopen as blocked during recovery, got %q", task.Status)
	}
	if !strings.Contains(strings.ToLower(task.ErrorMessage), "recovery") {
		t.Fatalf("expected blocked reason to mention recovery, got %q", task.ErrorMessage)
	}

	events, err := manager.ListEvents(ctx, teamID, 20)
	if err != nil {
		t.Fatalf("ListEvents() error = %v", err)
	}
	var sawReopened bool
	for _, event := range events {
		if event.TaskID != nil && *event.TaskID == "task-b" && event.EventType == "task_reopened" {
			sawReopened = true
			break
		}
	}
	if !sawReopened {
		t.Fatal("expected dependent task reopen event")
	}
}

func TestTeamManager_CompleteRecoveryTask_ReopensDependentsToPending(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := newTestTeamManager(t)

	teamID, err := manager.CreateTeam(ctx, "team-user-10", "user-10", "master")
	if err != nil {
		t.Fatalf("CreateTeam() error = %v", err)
	}

	if err := manager.RegisterTeammate(ctx, teamID, "worker-a", "Executes"); err != nil {
		t.Fatalf("RegisterTeammate() error = %v", err)
	}

	if err := manager.CreateTask(ctx, TeamTask{
		ID:     "task-a",
		TeamID: teamID,
		Title:  "Task A",
		Prompt: "Finish A",
		Status: TaskPending,
	}, nil); err != nil {
		t.Fatalf("CreateTask(task-a) error = %v", err)
	}
	if err := manager.CreateTask(ctx, TeamTask{
		ID:     "task-b",
		TeamID: teamID,
		Title:  "Task B",
		Prompt: "Depends on A",
		Status: TaskPending,
	}, []string{"task-a"}); err != nil {
		t.Fatalf("CreateTask(task-b) error = %v", err)
	}

	claimed, err := manager.ClaimNextTask(ctx, teamID, "worker-a")
	if err != nil {
		t.Fatalf("ClaimNextTask(task-a) error = %v", err)
	}
	if claimed == nil {
		t.Fatal("expected task-a claim")
	}

	if err := manager.FailTask(ctx, teamID, "task-a", "worker-a", "boom"); err != nil {
		t.Fatalf("FailTask(task-a) error = %v", err)
	}

	recoveryParent := "task-a"
	assignedAgent := "worker-a"
	if err := manager.CreateTask(ctx, TeamTask{
		ID:            "recovery-a",
		TeamID:        teamID,
		ParentTaskID:  &recoveryParent,
		Title:         "recovery:task-a",
		Prompt:        "Recover A",
		AssignedAgent: &assignedAgent,
		Status:        TaskPending,
	}, nil); err != nil {
		t.Fatalf("CreateTask(recovery-a) error = %v", err)
	}

	recoveryClaim, err := manager.ClaimNextTask(ctx, teamID, "worker-a")
	if err != nil {
		t.Fatalf("ClaimNextTask(recovery-a) error = %v", err)
	}
	if recoveryClaim == nil || recoveryClaim.ID != "recovery-a" {
		t.Fatalf("expected recovery task claim, got %#v", recoveryClaim)
	}

	if err := manager.CompleteTask(ctx, teamID, "recovery-a", "worker-a", "recovered"); err != nil {
		t.Fatalf("CompleteTask(recovery-a) error = %v", err)
	}

	task, err := manager.GetTask(ctx, teamID, "task-b")
	if err != nil {
		t.Fatalf("GetTask(task-b) error = %v", err)
	}
	if task == nil {
		t.Fatal("expected task-b to exist")
	}
	if task.Status != TaskPending {
		t.Fatalf("expected dependent to reopen as pending after recovery, got %q", task.Status)
	}

	events, err := manager.ListEvents(ctx, teamID, 20)
	if err != nil {
		t.Fatalf("ListEvents() error = %v", err)
	}
	var sawUnblocked bool
	for _, event := range events {
		if event.TaskID != nil && *event.TaskID == "task-b" && event.EventType == "task_unblocked" {
			sawUnblocked = true
			break
		}
	}
	if !sawUnblocked {
		t.Fatal("expected dependent task unblocked after recovery completion")
	}
}

func TestTeamManager_CreateTask_RejectsMissingDependency(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := newTestTeamManager(t)

	teamID, err := manager.CreateTeam(ctx, "team-user-11", "user-11", "master")
	if err != nil {
		t.Fatalf("CreateTeam() error = %v", err)
	}

	err = manager.CreateTask(ctx, TeamTask{
		ID:     "task-a",
		TeamID: teamID,
		Title:  "Task A",
		Prompt: "Depends on missing task",
		Status: TaskPending,
	}, []string{"missing-task"})
	if err == nil {
		t.Fatal("expected missing dependency error")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "dependency") {
		t.Fatalf("expected dependency error, got %v", err)
	}
}

func TestTeamManager_CreateTask_RejectsSelfDependency(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := newTestTeamManager(t)

	teamID, err := manager.CreateTeam(ctx, "team-user-12", "user-12", "master")
	if err != nil {
		t.Fatalf("CreateTeam() error = %v", err)
	}

	err = manager.CreateTask(ctx, TeamTask{
		ID:     "task-a",
		TeamID: teamID,
		Title:  "Task A",
		Prompt: "Depends on itself",
		Status: TaskPending,
	}, []string{"task-a"})
	if err == nil {
		t.Fatal("expected self dependency error")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "cycle") && !strings.Contains(strings.ToLower(err.Error()), "itself") {
		t.Fatalf("expected cycle/self dependency error, got %v", err)
	}
}

func TestTeamManager_CreateTask_RejectsIndirectCycle(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := newTestTeamManager(t)

	teamID, err := manager.CreateTeam(ctx, "team-user-13", "user-13", "master")
	if err != nil {
		t.Fatalf("CreateTeam() error = %v", err)
	}

	if err := manager.CreateTask(ctx, TeamTask{
		ID:     "task-a",
		TeamID: teamID,
		Title:  "Task A",
		Prompt: "A",
		Status: TaskPending,
	}, nil); err != nil {
		t.Fatalf("CreateTask(task-a) error = %v", err)
	}

	if err := manager.CreateTask(ctx, TeamTask{
		ID:     "task-b",
		TeamID: teamID,
		Title:  "Task B",
		Prompt: "B",
		Status: TaskPending,
	}, []string{"task-a"}); err != nil {
		t.Fatalf("CreateTask(task-b) error = %v", err)
	}

	parentTaskID := "task-a"
	err = manager.CreateTask(ctx, TeamTask{
		ID:           "task-a-retry",
		TeamID:       teamID,
		ParentTaskID: &parentTaskID,
		Title:        "Task A Retry",
		Prompt:       "Would close cycle",
		Status:       TaskPending,
	}, []string{"task-b"})
	if err == nil {
		t.Fatal("expected indirect cycle error")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "cycle") {
		t.Fatalf("expected cycle error, got %v", err)
	}
}

func TestTeamManager_MultipleDependencies_FanInRequiresAllParents(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := newTestTeamManager(t)

	teamID, err := manager.CreateTeam(ctx, "team-user-14", "user-14", "master")
	if err != nil {
		t.Fatalf("CreateTeam() error = %v", err)
	}

	if err := manager.RegisterTeammate(ctx, teamID, "worker-a", "Executes"); err != nil {
		t.Fatalf("RegisterTeammate() error = %v", err)
	}

	for _, id := range []string{"task-a", "task-b"} {
		if err := manager.CreateTask(ctx, TeamTask{
			ID:     id,
			TeamID: teamID,
			Title:  id,
			Prompt: id,
			Status: TaskPending,
		}, nil); err != nil {
			t.Fatalf("CreateTask(%s) error = %v", id, err)
		}
	}

	if err := manager.CreateTask(ctx, TeamTask{
		ID:     "task-c",
		TeamID: teamID,
		Title:  "Task C",
		Prompt: "Needs A and B",
		Status: TaskPending,
	}, []string{"task-a", "task-b"}); err != nil {
		t.Fatalf("CreateTask(task-c) error = %v", err)
	}

	taskC, err := manager.GetTask(ctx, teamID, "task-c")
	if err != nil {
		t.Fatalf("GetTask(task-c) error = %v", err)
	}
	if taskC == nil || taskC.Status != TaskBlocked {
		t.Fatalf("expected task-c to start blocked, got %#v", taskC)
	}

	first, err := manager.ClaimNextTask(ctx, teamID, "worker-a")
	if err != nil {
		t.Fatalf("ClaimNextTask(first) error = %v", err)
	}
	if first == nil {
		t.Fatal("expected first claim")
	}
	if err := manager.CompleteTask(ctx, teamID, first.ID, "worker-a", "done"); err != nil {
		t.Fatalf("CompleteTask(first) error = %v", err)
	}

	taskC, err = manager.GetTask(ctx, teamID, "task-c")
	if err != nil {
		t.Fatalf("GetTask(task-c after first parent) error = %v", err)
	}
	if taskC == nil || taskC.Status != TaskBlocked {
		t.Fatalf("expected task-c to remain blocked until all parents complete, got %#v", taskC)
	}

	second, err := manager.ClaimNextTask(ctx, teamID, "worker-a")
	if err != nil {
		t.Fatalf("ClaimNextTask(second) error = %v", err)
	}
	if second == nil {
		t.Fatal("expected second claim")
	}
	if err := manager.CompleteTask(ctx, teamID, second.ID, "worker-a", "done"); err != nil {
		t.Fatalf("CompleteTask(second) error = %v", err)
	}

	taskC, err = manager.GetTask(ctx, teamID, "task-c")
	if err != nil {
		t.Fatalf("GetTask(task-c after both parents) error = %v", err)
	}
	if taskC == nil || taskC.Status != TaskPending {
		t.Fatalf("expected task-c to unblock after both parents complete, got %#v", taskC)
	}
}

func TestTeamManager_ClaimNextTask_RequeuesExpiredWorkerLease(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	rawManager := newTestTeamManager(t)
	manager := rawManager.(*DefaultTeamManager)

	teamID, err := manager.CreateTeam(ctx, "team-user-15", "user-15", "master")
	if err != nil {
		t.Fatalf("CreateTeam() error = %v", err)
	}

	if err := manager.RegisterTeammate(ctx, teamID, "worker-a", "Executes"); err != nil {
		t.Fatalf("RegisterTeammate(worker-a) error = %v", err)
	}
	if err := manager.RegisterTeammate(ctx, teamID, "worker-b", "Recovers expired work"); err != nil {
		t.Fatalf("RegisterTeammate(worker-b) error = %v", err)
	}

	if err := manager.CreateTask(ctx, TeamTask{
		ID:     "task-expired",
		TeamID: teamID,
		Title:  "Lease-sensitive task",
		Prompt: "Can be reclaimed if lease expires",
		Status: TaskPending,
	}, nil); err != nil {
		t.Fatalf("CreateTask() error = %v", err)
	}

	claimed, err := manager.ClaimNextTask(ctx, teamID, "worker-a")
	if err != nil {
		t.Fatalf("ClaimNextTask(worker-a) error = %v", err)
	}
	if claimed == nil || claimed.ID != "task-expired" {
		t.Fatalf("expected worker-a to claim task-expired, got %#v", claimed)
	}

	expiredAt := time.Now().UTC().Add(-time.Minute)
	if _, err := manager.store.db.ExecContext(ctx,
		`UPDATE team_members SET status = ?, lease_expires_at = ?, last_heartbeat_at = ? WHERE team_id = ? AND agent_name = ?`,
		"stale", expiredAt, expiredAt, teamID, "worker-a",
	); err != nil {
		t.Fatalf("force expire lease error = %v", err)
	}

	reclaimed, err := manager.ClaimNextTask(ctx, teamID, "worker-b")
	if err != nil {
		t.Fatalf("ClaimNextTask(worker-b) error = %v", err)
	}
	if reclaimed == nil || reclaimed.ID != "task-expired" {
		t.Fatalf("expected worker-b to reclaim expired task, got %#v", reclaimed)
	}
	if reclaimed.AssignedAgent == nil || *reclaimed.AssignedAgent != "worker-b" {
		t.Fatalf("expected reassigned lease owner worker-b, got %#v", reclaimed.AssignedAgent)
	}

	task, err := manager.GetTask(ctx, teamID, "task-expired")
	if err != nil {
		t.Fatalf("GetTask() error = %v", err)
	}
	if task == nil || task.Status != TaskRunning {
		t.Fatalf("expected task to be running under new owner, got %#v", task)
	}

	events, err := manager.ListEvents(ctx, teamID, 20)
	if err != nil {
		t.Fatalf("ListEvents() error = %v", err)
	}
	var sawLeaseRequeue bool
	for _, event := range events {
		if event.TaskID != nil && *event.TaskID == "task-expired" && event.EventType == "task_requeued" && strings.Contains(strings.ToLower(event.Payload), "lease") {
			sawLeaseRequeue = true
			break
		}
	}
	if !sawLeaseRequeue {
		t.Fatal("expected task_requeued event mentioning expired lease")
	}
}
