package agent

import "context"

type LeadRuntime struct {
	manager   TeamManager
	agentName string
}

func NewLeadRuntime(manager TeamManager, agentName string) *LeadRuntime {
	return &LeadRuntime{
		manager:   manager,
		agentName: agentName,
	}
}

func (l *LeadRuntime) CollectInbox(ctx context.Context, teamID string, limit int) ([]MailMessage, error) {
	return l.manager.PullMessages(ctx, teamID, l.agentName, limit)
}

func (l *LeadRuntime) CollectEvents(ctx context.Context, teamID string, limit int) ([]TaskEvent, error) {
	return l.manager.ListEvents(ctx, teamID, limit)
}
