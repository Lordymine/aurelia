package agent

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"
)

const MasterAgentName = "master"
const (
	TeamStatusActive            = "active"
	TeamStatusPaused            = "paused"
	TeamStatusCancelled         = "cancelled"
	TeamStatusCompleted         = "completed"
	TeamStatusAttentionRequired = "attention_required"
)

type LeadNotifier func(teamKey string, message string)

type MasterTeamService struct {
	manager             TeamManager
	llm                 LLMProvider
	registry            *ToolRegistry
	maxIterations       int
	notifyLead          LeadNotifier
	maxRecoveryAttempts int

	mu            sync.Mutex
	teamByKey     map[string]string
	userByKey     map[string]string
	memberSeen    map[string]map[string]bool
	workerLoops   map[string]*workerLoopHandle
	recoveryCount map[string]int
}

type workerLoopHandle struct {
	wake   chan struct{}
	cancel context.CancelFunc
}

const recoveryTaskPrefix = "recovery:"

func NewMasterTeamService(manager TeamManager, llm LLMProvider, registry *ToolRegistry, maxIterations int, notify LeadNotifier) (*MasterTeamService, error) {
	if manager == nil {
		return nil, fmt.Errorf("team manager is required")
	}
	if llm == nil {
		return nil, fmt.Errorf("llm provider is required")
	}
	if registry == nil {
		return nil, fmt.Errorf("tool registry is required")
	}

	return &MasterTeamService{
		manager:             manager,
		llm:                 llm,
		registry:            registry,
		maxIterations:       maxIterations,
		notifyLead:          notify,
		maxRecoveryAttempts: 1,
		teamByKey:           make(map[string]string),
		userByKey:           make(map[string]string),
		memberSeen:          make(map[string]map[string]bool),
		workerLoops:         make(map[string]*workerLoopHandle),
		recoveryCount:       make(map[string]int),
	}, nil
}

func (s *MasterTeamService) Spawn(ctx context.Context, teamKey, userID, agentName, roleDescription, taskPrompt string, allowedTools ...string) (string, error) {
	teamID, err := s.ensureTeam(ctx, teamKey, userID)
	if err != nil {
		return "", err
	}
	if err := s.reactivateTeamIfTerminal(ctx, teamID); err != nil {
		return "", err
	}
	if err := s.ensureTeammate(ctx, teamID, agentName, roleDescription); err != nil {
		return "", err
	}

	taskID := uuid.NewString()
	task := TeamTask{
		ID:           taskID,
		TeamID:       teamID,
		RunID:        s.resolveRunID(ctx, teamID),
		Title:        agentName,
		Prompt:       taskPrompt,
		Workdir:      s.resolveTaskWorkdir(ctx),
		AllowedTools: append([]string(nil), allowedTools...),
		Status:       TaskPending,
	}

	var dependsOn []string
	if parentTeamID, parentTaskID, ok := TaskContextFromContext(ctx); ok && parentTeamID == teamID {
		dependsOn = []string{parentTaskID}
	}
	if err := s.manager.CreateTask(ctx, task, dependsOn); err != nil {
		return "", err
	}

	s.ensureWorkerLoop(teamID, teamKey, userID, agentName, roleDescription)
	return taskID, nil
}

func (s *MasterTeamService) reactivateTeamIfTerminal(ctx context.Context, teamID string) error {
	status, err := s.manager.GetTeamStatus(ctx, teamID)
	if err != nil {
		return err
	}
	if !isTerminalTeamStatus(status) {
		return nil
	}
	return s.manager.SetTeamStatus(ctx, teamID, TeamStatusActive)
}

func (s *MasterTeamService) resolveRunID(ctx context.Context, teamID string) string {
	if runID, ok := RunContextFromContext(ctx); ok {
		return runID
	}
	if parentTeamID, parentTaskID, ok := TaskContextFromContext(ctx); ok && parentTeamID == teamID {
		if parentTask, err := s.manager.GetTask(ctx, teamID, parentTaskID); err == nil && parentTask != nil && parentTask.RunID != "" {
			return parentTask.RunID
		}
	}
	return uuid.NewString()
}

func (s *MasterTeamService) resolveTaskWorkdir(ctx context.Context) string {
	if workdir, ok := WorkdirFromContext(ctx); ok {
		return workdir
	}
	if parentTeamID, parentTaskID, ok := TaskContextFromContext(ctx); ok {
		if parentTask, err := s.manager.GetTask(ctx, parentTeamID, parentTaskID); err == nil && parentTask != nil {
			return parentTask.Workdir
		}
	}
	return ""
}

func (s *MasterTeamService) ensureTeam(ctx context.Context, teamKey, userID string) (string, error) {
	s.mu.Lock()
	teamID := s.teamByKey[teamKey]
	s.mu.Unlock()
	if teamID != "" {
		return teamID, nil
	}

	existingTeamID, err := s.manager.GetTeamIDByKey(ctx, teamKey)
	if err != nil {
		return "", err
	}
	if existingTeamID != "" {
		s.rememberTeam(teamKey, existingTeamID, userID)
		return existingTeamID, nil
	}

	createdTeamID, err := s.manager.CreateTeam(ctx, teamKey, userID, MasterAgentName)
	if err != nil {
		return "", err
	}
	s.rememberTeam(teamKey, createdTeamID, userID)
	return createdTeamID, nil
}

func (s *MasterTeamService) rememberTeam(teamKey, teamID, userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if existing := s.teamByKey[teamKey]; existing != "" {
		return
	}
	s.teamByKey[teamKey] = teamID
	s.userByKey[teamKey] = userID
	if _, ok := s.memberSeen[teamID]; !ok {
		s.memberSeen[teamID] = map[string]bool{MasterAgentName: true}
	}
}

func (s *MasterTeamService) ensureTeammate(ctx context.Context, teamID, agentName, roleDescription string) error {
	s.mu.Lock()
	teamMembers, ok := s.memberSeen[teamID]
	if !ok {
		teamMembers = make(map[string]bool)
		s.memberSeen[teamID] = teamMembers
	}
	if teamMembers[agentName] {
		s.mu.Unlock()
		return nil
	}
	teamMembers[agentName] = true
	s.mu.Unlock()

	return s.manager.RegisterTeammate(ctx, teamID, agentName, roleDescription)
}

func (s *MasterTeamService) notifyMaster(teamKey, message string) {
	if s.notifyLead != nil && strings.TrimSpace(message) != "" {
		s.notifyLead(teamKey, message)
	}
}

func (s *MasterTeamService) Pause(ctx context.Context, teamKey string) error {
	teamID, err := s.ensureExistingTeam(ctx, teamKey)
	if err != nil {
		return err
	}
	if err := s.manager.SetTeamStatus(ctx, teamID, TeamStatusPaused); err != nil {
		return err
	}
	s.notifyMaster(teamKey, "Parei de distribuir novas tasks para a equipe. O que ja estava rodando pode terminar, mas nada novo sera iniciado ate voce mandar retomar.")
	return nil
}

func (s *MasterTeamService) Resume(ctx context.Context, teamKey string) error {
	teamID, err := s.ensureExistingTeam(ctx, teamKey)
	if err != nil {
		return err
	}
	if err := s.manager.SetTeamStatus(ctx, teamID, TeamStatusActive); err != nil {
		return err
	}

	s.mu.Lock()
	userID := s.userByKey[teamKey]
	members := mapsClone(s.memberSeen[teamID])
	s.mu.Unlock()
	for agentName := range members {
		if agentName == MasterAgentName {
			continue
		}
		s.ensureWorkerLoop(teamID, teamKey, userID, agentName, "Resumed worker")
	}

	s.notifyMaster(teamKey, "Retomei a equipe. Vou voltar a distribuir as tasks pendentes e te aviso conforme o trabalho andar.")
	return nil
}

func (s *MasterTeamService) Cancel(ctx context.Context, teamKey, reason string) error {
	teamID, err := s.ensureExistingTeam(ctx, teamKey)
	if err != nil {
		return err
	}
	if strings.TrimSpace(reason) == "" {
		reason = "cancelado pelo usuario"
	}
	if err := s.manager.SetTeamStatus(ctx, teamID, TeamStatusCancelled); err != nil {
		return err
	}
	s.cancelWorkerLoops(teamID)
	if err := s.manager.CancelActiveTasks(ctx, teamID, reason); err != nil {
		return err
	}
	s.notifyMaster(teamKey, fmt.Sprintf("Cancelei a operacao do time. Interrompi os workers ativos e marquei as tasks abertas como canceladas. Motivo: %s.", reason))
	return nil
}

func (s *MasterTeamService) ensureExistingTeam(ctx context.Context, teamKey string) (string, error) {
	s.mu.Lock()
	teamID := s.teamByKey[teamKey]
	s.mu.Unlock()
	if teamID != "" {
		return teamID, nil
	}
	return s.manager.GetTeamIDByKey(ctx, teamKey)
}

func (s *MasterTeamService) cancelWorkerLoops(teamID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for key, handle := range s.workerLoops {
		if strings.HasPrefix(key, teamID+"::") && handle.cancel != nil {
			handle.cancel()
			delete(s.workerLoops, key)
		}
	}
}

func mapsClone(src map[string]bool) map[string]bool {
	dst := make(map[string]bool, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func isTerminalTeamStatus(status string) bool {
	switch status {
	case TeamStatusCompleted, TeamStatusCancelled, TeamStatusAttentionRequired:
		return true
	default:
		return false
	}
}
