package agent

import (
	"context"
	"fmt"
)

func (s *MasterTeamService) Rehydrate(ctx context.Context) error {
	manager, ok := s.manager.(*DefaultTeamManager)
	if !ok {
		return fmt.Errorf("rehydrate requires DefaultTeamManager")
	}

	teams, err := manager.store.listTeams(ctx)
	if err != nil {
		return err
	}

	for _, team := range teams {
		status, err := s.manager.GetTeamStatus(ctx, team.TeamID)
		if err != nil {
			return err
		}
		if status != TeamStatusActive {
			s.rememberRehydratedTeam(team.TeamKey, team.TeamID, team.UserID)
			continue
		}
		if err := manager.store.requeueRunningTasks(ctx, team.TeamID); err != nil {
			return err
		}
		s.rememberRehydratedTeam(team.TeamKey, team.TeamID, team.UserID)

		tasks, err := s.manager.ListTasks(ctx, team.TeamID)
		if err != nil {
			return err
		}
		for _, task := range tasks {
			if task.Status == TaskPending && task.AssignedAgent != nil && *task.AssignedAgent != "" {
				s.ensureWorkerLoop(team.TeamID, team.TeamKey, team.UserID, *task.AssignedAgent, "Rehydrated worker")
			}
		}
	}

	return nil
}

func (s *MasterTeamService) rememberRehydratedTeam(teamKey, teamID, userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.teamByKey[teamKey] = teamID
	s.userByKey[teamKey] = userID
	if _, ok := s.memberSeen[teamID]; !ok {
		s.memberSeen[teamID] = map[string]bool{MasterAgentName: true}
	}
}
