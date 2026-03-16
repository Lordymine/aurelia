package agent

import (
	"context"
	"database/sql"
	"fmt"
)

func (s *SQLiteTaskStore) unblockDependentsTx(ctx context.Context, tx *sql.Tx, teamID, completedTaskID string) error {
	dependentIDs, err := s.listBlockedDependentsTx(ctx, tx, teamID, completedTaskID)
	if err != nil {
		return err
	}

	for _, dependentID := range dependentIDs {
		if err := s.unblockDependentIfReadyTx(ctx, tx, teamID, completedTaskID, dependentID); err != nil {
			return err
		}
	}

	return nil
}

func (s *SQLiteTaskStore) listBlockedDependentsTx(ctx context.Context, tx *sql.Tx, teamID, dependencyTaskID string) ([]string, error) {
	return s.listDependentsByStatusTx(ctx, tx, teamID, dependencyTaskID, TaskBlocked)
}

func (s *SQLiteTaskStore) unblockDependentIfReadyTx(ctx context.Context, tx *sql.Tx, teamID, completedTaskID, dependentID string) error {
	ready, err := s.areAllDependenciesCompletedTx(ctx, tx, dependentID)
	if err != nil {
		return err
	}
	if !ready {
		return nil
	}

	if err := s.updateTaskStatusTx(ctx, tx, teamID, dependentID, TaskBlocked, TaskPending, ""); err != nil {
		return err
	}

	return s.insertTaskEventTx(ctx, tx, TaskEvent{
		TeamID:    teamID,
		TaskID:    &dependentID,
		AgentName: MasterAgentName,
		EventType: "task_unblocked",
		Payload:   fmt.Sprintf("dependency %s completed", completedTaskID),
	})
}
