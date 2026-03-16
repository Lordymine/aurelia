package agent

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/google/uuid"
)

func (s *MasterTeamService) scheduleRecoveryTask(ctx context.Context, teamID, teamKey string, update MailMessage) {
	if update.TaskID == nil {
		return
	}

	baseTask, attempts, ok := s.prepareRecoveryBaseTask(ctx, teamID, teamKey, *update.TaskID)
	if !ok {
		return
	}

	recoveryPrompt := s.buildRecoveryPrompt(ctx, teamID, baseTask, update.Body, attempts)
	recoveryAgent := s.selectRecoveryAgent(teamID, *baseTask.AssignedAgent)
	recoveryTask := TeamTask{
		ID:            uuid.NewString(),
		TeamID:        teamID,
		RunID:         baseTask.RunID,
		ParentTaskID:  &baseTask.ID,
		Title:         recoveryTaskPrefix + baseTask.Title,
		Prompt:        recoveryPrompt,
		Workdir:       baseTask.Workdir,
		AssignedAgent: &recoveryAgent,
		Status:        TaskPending,
	}
	if err := s.manager.CreateTask(ctx, recoveryTask, nil); err != nil {
		return
	}

	s.notifyRecoveryPlanning(ctx, teamKey, update.FromAgent, baseTask, attempts)

	s.mu.Lock()
	userID := s.userByKey[teamKey]
	s.mu.Unlock()
	s.ensureWorkerLoop(teamID, teamKey, userID, recoveryAgent, "Recovery specialist")
}

func (s *MasterTeamService) prepareRecoveryBaseTask(ctx context.Context, teamID, teamKey, failedTaskID string) (*TeamTask, int, bool) {
	originalTask, err := s.manager.GetTask(ctx, teamID, failedTaskID)
	if err != nil || originalTask == nil {
		return nil, 0, false
	}

	baseTask := originalTask
	if strings.HasPrefix(originalTask.Title, recoveryTaskPrefix) && originalTask.ParentTaskID != nil {
		if parentTask, err := s.manager.GetTask(ctx, teamID, *originalTask.ParentTaskID); err == nil && parentTask != nil {
			baseTask = parentTask
		}
	}
	if baseTask.AssignedAgent == nil || *baseTask.AssignedAgent == "" {
		return nil, 0, false
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	attempts := s.recoveryCount[baseTask.ID]
	if attempts >= s.maxRecoveryAttempts {
		s.notifyMaster(teamKey, fmt.Sprintf("Task `%s` escalated after %d failed recovery attempt(s). Tentei replanejar isso algumas vezes, mas o time nao conseguiu se recuperar sozinho. Vou parar por aqui e esperar sua orientacao.", baseTask.Title, attempts))
		return nil, 0, false
	}
	s.recoveryCount[baseTask.ID] = attempts + 1
	return baseTask, attempts + 1, true
}

func (s *MasterTeamService) notifyRecoveryPlanning(ctx context.Context, teamKey, failedAgent string, baseTask *TeamTask, attempt int) {
	snapshot, err := s.BuildExecutionStatusSnapshot(ctx, teamKey, baseTask.RunID)
	if err != nil {
		s.notifyMaster(teamKey, fmt.Sprintf("O especialista `%s` falhou enquanto trabalhava em `%s`. Ja iniciei o replanejamento e estou tentando uma nova abordagem (tentativa %d).", failedAgent, baseTask.Title, attempt))
		return
	}

	s.notifyMaster(teamKey, fmt.Sprintf(
		"O especialista `%s` falhou em `%s`, entao eu iniciei o replanejamento da equipe e abri uma tentativa %d.\n%s\nEquipe em `%s`: %d pendente(s), %d em andamento, %d bloqueada(s), %d concluida(s), %d com falha e %d cancelada(s).",
		failedAgent,
		baseTask.Title,
		attempt,
		classicStatusLine(snapshot),
		fallbackTeamStatus(snapshot.TeamStatus),
		snapshot.Pending,
		snapshot.Running,
		snapshot.Blocked,
		snapshot.Completed,
		snapshot.Failed,
		snapshot.Cancelled,
	))
}

func (s *MasterTeamService) buildRecoveryPrompt(ctx context.Context, teamID string, baseTask *TeamTask, failureReason string, attempt int) string {
	events, err := s.manager.ListEvents(ctx, teamID, 100)
	if err != nil {
		return fmt.Sprintf("Diagnostique a falha da task %q. Erro observado: %s. Proponha uma correcao objetiva ou um proximo passo seguro.", baseTask.Title, failureReason)
	}

	var history []string
	for _, event := range events {
		if event.TaskID == nil || *event.TaskID != baseTask.ID {
			continue
		}
		payload := strings.TrimSpace(event.Payload)
		if payload == "" {
			history = append(history, fmt.Sprintf("- %s por %s", event.EventType, event.AgentName))
			continue
		}
		history = append(history, fmt.Sprintf("- %s por %s: %s", event.EventType, event.AgentName, payload))
	}
	if len(history) == 0 {
		return fmt.Sprintf("Diagnostique a falha da task %q. Erro observado: %s. Proponha uma correcao objetiva ou um proximo passo seguro.", baseTask.Title, failureReason)
	}

	return fmt.Sprintf(
		"Diagnostique a falha da task %q. Esta e a tentativa de replanejamento %d.\nErro observado: %s.\nHistorico recente:\n%s\nProponha uma correcao objetiva ou um proximo passo seguro.",
		baseTask.Title,
		attempt,
		failureReason,
		strings.Join(history, "\n"),
	)
}

func (s *MasterTeamService) selectRecoveryAgent(teamID, failedAgent string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	members := s.memberSeen[teamID]
	if len(members) == 0 {
		return failedAgent
	}

	var candidates []string
	for member := range members {
		if member == "" || member == MasterAgentName || member == failedAgent {
			continue
		}
		candidates = append(candidates, member)
	}
	slices.Sort(candidates)
	if len(candidates) > 0 {
		return candidates[0]
	}
	return failedAgent
}
