package agent

import (
	"context"
	"fmt"
	"slices"
	"strings"
)

func (s *MasterTeamService) BuildStatusSnapshot(ctx context.Context, teamKey string) (TeamStatusSnapshot, error) {
	return s.buildSnapshot(ctx, teamKey, "")
}

func (s *MasterTeamService) BuildExecutionStatusSnapshot(ctx context.Context, teamKey, runID string) (TeamStatusSnapshot, error) {
	return s.buildSnapshot(ctx, teamKey, runID)
}

func (s *MasterTeamService) buildSnapshot(ctx context.Context, teamKey, runID string) (TeamStatusSnapshot, error) {
	s.mu.Lock()
	teamID := s.teamByKey[teamKey]
	s.mu.Unlock()
	if teamID == "" {
		return TeamStatusSnapshot{}, fmt.Errorf("team not found for key %s", teamKey)
	}

	tasks, err := s.manager.ListTasks(ctx, teamID)
	if err != nil {
		return TeamStatusSnapshot{}, err
	}

	filtered := filterTasksByRunID(tasks, runID)
	teamStatus, err := s.manager.GetTeamStatus(ctx, teamID)
	if err != nil {
		return TeamStatusSnapshot{}, err
	}
	snapshot := TeamStatusSnapshot{TeamKey: teamKey, TeamID: teamID, TeamStatus: teamStatus}
	if len(filtered) == 0 {
		return snapshot, nil
	}

	for _, statuses := range groupTaskStatusesByRoot(filtered) {
		snapshot.TotalTasks++
		switch resolveLogicalStatus(statuses) {
		case TaskPending:
			snapshot.Pending++
		case TaskRunning:
			snapshot.Running++
		case TaskBlocked:
			snapshot.Blocked++
		case TaskCompleted:
			snapshot.Completed++
		case TaskFailed:
			snapshot.Failed++
		case TaskCancelled:
			snapshot.Cancelled++
		}
	}
	return snapshot, nil
}

func filterTasksByRunID(tasks []TeamTask, runID string) []TeamTask {
	if runID == "" {
		return tasks
	}
	filtered := make([]TeamTask, 0, len(tasks))
	for _, task := range tasks {
		if task.RunID == runID {
			filtered = append(filtered, task)
		}
	}
	return filtered
}

func groupTaskStatusesByRoot(tasks []TeamTask) map[string][]TaskStatus {
	taskByID := make(map[string]TeamTask, len(tasks))
	for _, task := range tasks {
		taskByID[task.ID] = task
	}

	groupStatuses := make(map[string][]TaskStatus)
	for _, task := range tasks {
		rootID := task.ID
		current := task
		seen := map[string]bool{current.ID: true}
		for current.ParentTaskID != nil && *current.ParentTaskID != "" {
			parent, ok := taskByID[*current.ParentTaskID]
			if !ok || seen[parent.ID] {
				break
			}
			rootID = parent.ID
			current = parent
			seen[current.ID] = true
		}
		groupStatuses[rootID] = append(groupStatuses[rootID], task.Status)
	}
	return groupStatuses
}

func resolveLogicalStatus(statuses []TaskStatus) TaskStatus {
	if len(statuses) == 0 {
		return TaskPending
	}

	switch {
	case slices.Contains(statuses, TaskRunning):
		return TaskRunning
	case slices.Contains(statuses, TaskBlocked):
		return TaskBlocked
	case slices.Contains(statuses, TaskPending):
		return TaskPending
	case slices.Contains(statuses, TaskCompleted):
		return TaskCompleted
	case slices.Contains(statuses, TaskFailed):
		return TaskFailed
	case slices.Contains(statuses, TaskCancelled):
		return TaskCancelled
	default:
		return statuses[0]
	}
}

func (s *MasterTeamService) formatMasterNotification(snapshot TeamStatusSnapshot, processedCount int, lines []string) string {
	statusLine := classicStatusLine(snapshot)
	humanStatusLine := fmt.Sprintf(
		"Equipe em `%s`: %d pendente(s), %d em andamento, %d bloqueada(s), %d concluida(s), %d com falha e %d cancelada(s).",
		fallbackTeamStatus(snapshot.TeamStatus),
		snapshot.Pending,
		snapshot.Running,
		snapshot.Blocked,
		snapshot.Completed,
		snapshot.Failed,
		snapshot.Cancelled,
	)

	body := strings.Join(lines, "\n")
	if snapshot.Pending == 0 && snapshot.Running == 0 && snapshot.Blocked == 0 && snapshot.TotalTasks > 0 {
		return fmt.Sprintf(
			"Fechei este ciclo do time com %s.\n%s\n%s\n\nConsolidei o que saiu deste run:\n%s\n\nEncerrando a operacao deste ciclo: parei os workers deste time e limpei o estado transitorio antes do proximo passo.",
			classifyFinalSnapshot(snapshot),
			statusLine,
			humanStatusLine,
			body,
		)
	}

	return fmt.Sprintf("Estou acompanhando o time.\n%s\n%s\n\nDesde a ultima atualizacao, %d task(s) andaram:\n%s", statusLine, humanStatusLine, processedCount, body)
}

func (s *MasterTeamService) finalizeTeamRunIfIdle(ctx context.Context, teamKey, runID string, snapshot *TeamStatusSnapshot) {
	if snapshot == nil || snapshot.TeamID == "" {
		return
	}
	if snapshot.Pending != 0 || snapshot.Running != 0 || snapshot.Blocked != 0 || snapshot.TotalTasks == 0 {
		return
	}

	finalStatus := TeamStatusCompleted
	if snapshot.Failed > 0 || snapshot.Cancelled > 0 {
		finalStatus = TeamStatusAttentionRequired
	}
	if err := s.manager.SetTeamStatus(ctx, snapshot.TeamID, finalStatus); err == nil {
		snapshot.TeamStatus = finalStatus
	}

	s.cancelWorkerLoops(snapshot.TeamID)
	s.clearRunRuntimeState(ctx, snapshot.TeamID, runID)
}

func (s *MasterTeamService) clearRunRuntimeState(ctx context.Context, teamID, runID string) {
	tasks, err := s.manager.ListTasks(ctx, teamID)
	if err != nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	for _, task := range tasks {
		if runID != "" && task.RunID != runID {
			continue
		}
		delete(s.recoveryCount, task.ID)
	}
}

func classifyFinalSnapshot(snapshot TeamStatusSnapshot) string {
	switch {
	case snapshot.Failed == 0 && snapshot.Cancelled == 0:
		return "sucesso total"
	case snapshot.Completed > 0:
		return "conclusao parcial"
	default:
		return "bloqueio terminal"
	}
}

func fallbackTeamStatus(status string) string {
	if strings.TrimSpace(status) == "" {
		return "active"
	}
	return status
}

func classicStatusLine(snapshot TeamStatusSnapshot) string {
	return fmt.Sprintf(
		"status: pending=%d running=%d blocked=%d completed=%d failed=%d cancelled=%d total=%d",
		snapshot.Pending,
		snapshot.Running,
		snapshot.Blocked,
		snapshot.Completed,
		snapshot.Failed,
		snapshot.Cancelled,
		snapshot.TotalTasks,
	)
}
