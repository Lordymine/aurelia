package agent

import (
	"context"
	"fmt"
	"strings"
)

func (s *MasterTeamService) ensureWorkerLoop(teamID, teamKey, userID, agentName, roleDescription string) {
	workerKey := teamID + "::" + agentName

	s.mu.Lock()
	handle, exists := s.workerLoops[workerKey]
	if !exists {
		loopCtx, cancel := context.WithCancel(context.Background())
		handle = &workerLoopHandle{wake: make(chan struct{}, 1), cancel: cancel}
		s.workerLoops[workerKey] = handle
		go s.runWorkerLoop(loopCtx, handle, teamID, teamKey, userID, agentName, roleDescription)
	}
	s.mu.Unlock()

	select {
	case handle.wake <- struct{}{}:
	default:
	}
}

func (s *MasterTeamService) runWorkerLoop(loopCtx context.Context, handle *workerLoopHandle, teamID, teamKey, userID, agentName, roleDescription string) {
	workerCtx := WithTeamContext(loopCtx, teamKey, userID)
	workerCtx = WithAgentContext(workerCtx, agentName)
	executor := &loopTaskExecutor{
		llm:           s.llm,
		registry:      s.registry,
		maxIterations: s.maxIterations,
		agentName:     agentName,
		roleDesc:      roleDescription,
	}
	worker := NewWorkerRuntime(agentName, s.manager, executor)

	for {
		select {
		case <-loopCtx.Done():
			return
		case _, ok := <-handle.wake:
			if !ok {
				return
			}
		}
		processedCount, err := worker.RunUntilIdle(workerCtx, teamID)
		if err != nil {
			if loopCtx.Err() != nil {
				return
			}
			s.notifyMaster(teamKey, fmt.Sprintf("Precisei interromper o worker `%s` por causa de uma falha operacional: %v", agentName, err))
			continue
		}
		if processedCount == 0 {
			continue
		}
		s.flushLeadInbox(workerCtx, teamID, teamKey, processedCount)
	}
}

func (s *MasterTeamService) flushLeadInbox(ctx context.Context, teamID, teamKey string, processedCount int) {
	lead := NewLeadRuntime(s.manager, MasterAgentName)
	updates, err := lead.CollectInbox(ctx, teamID, 50)
	if err != nil {
		s.notifyMaster(teamKey, fmt.Sprintf("Perdi uma atualizacao interna do time enquanto consolidava o andamento: %v", err))
		return
	}
	if len(updates) == 0 {
		return
	}

	lines, runID := s.summarizeLeadUpdates(ctx, teamID, teamKey, updates)
	if len(lines) == 0 {
		return
	}

	snapshot, err := s.BuildExecutionStatusSnapshot(ctx, teamKey, runID)
	if err != nil {
		s.notifyMaster(teamKey, fmt.Sprintf("O time andou, mas eu nao consegui montar o panorama consolidado agora.\n\nMovimentos recentes:\n%s", strings.Join(lines, "\n")))
		return
	}

	s.finalizeTeamRunIfIdle(ctx, teamKey, runID, &snapshot)
	s.notifyMaster(teamKey, s.formatMasterNotification(snapshot, processedCount, lines))
}

func (s *MasterTeamService) summarizeLeadUpdates(ctx context.Context, teamID, teamKey string, updates []MailMessage) ([]string, string) {
	var lines []string
	for _, update := range updates {
		switch update.Kind {
		case "result", "status_update":
			lines = append(lines, fmt.Sprintf("- %s: %s", update.FromAgent, update.Body))
		case "blocker":
			lines = append(lines, fmt.Sprintf("- %s encontrou um bloqueio: %s", update.FromAgent, update.Body))
			s.scheduleRecoveryTask(ctx, teamID, teamKey, update)
		}
	}
	return lines, s.resolveRunIDFromUpdates(ctx, teamID, updates)
}

func (s *MasterTeamService) resolveRunIDFromUpdates(ctx context.Context, teamID string, updates []MailMessage) string {
	if len(updates) == 0 {
		return ""
	}
	last := updates[len(updates)-1]
	if last.TaskID == nil {
		return ""
	}
	task, err := s.manager.GetTask(ctx, teamID, *last.TaskID)
	if err != nil || task == nil {
		return ""
	}
	return task.RunID
}
