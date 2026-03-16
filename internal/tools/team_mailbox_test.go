package tools

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/kocar/aurelia/internal/agent"
)

type mailboxManagerStub struct {
	posted []agent.MailMessage
	pulled []agent.MailMessage
}

func (m *mailboxManagerStub) PostMessage(ctx context.Context, msg agent.MailMessage) error {
	m.posted = append(m.posted, msg)
	return nil
}

func (m *mailboxManagerStub) PullMessages(ctx context.Context, teamID, agentName string, limit int) ([]agent.MailMessage, error) {
	return append([]agent.MailMessage(nil), m.pulled...), nil
}

func TestSendTeamMessageTool_Execute_PostsMailboxMessage(t *testing.T) {
	manager := &mailboxManagerStub{}
	tool := NewSendTeamMessageTool(manager)

	ctx := agent.WithAgentContext(
		agent.WithTaskContext(context.Background(), "team-1", "task-1"),
		"worker-a",
	)

	_, err := tool.Execute(ctx, map[string]interface{}{
		"to_agent": "worker-b",
		"body":     "preciso que voce valide a API",
		"kind":     "question",
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if len(manager.posted) != 1 {
		t.Fatalf("expected one mailbox message, got %d", len(manager.posted))
	}
	got := manager.posted[0]
	if got.TeamID != "team-1" || got.FromAgent != "worker-a" || got.ToAgent != "worker-b" {
		t.Fatalf("unexpected mailbox routing: %#v", got)
	}
	if got.TaskID == nil || *got.TaskID != "task-1" {
		t.Fatalf("expected task linkage, got %#v", got.TaskID)
	}
	if got.Kind != "question" {
		t.Fatalf("expected kind question, got %q", got.Kind)
	}
}

func TestReadTeamInboxTool_Execute_ReturnsJSONMessages(t *testing.T) {
	manager := &mailboxManagerStub{
		pulled: []agent.MailMessage{
			{
				ID:        "mail-1",
				TeamID:    "team-1",
				FromAgent: "worker-b",
				ToAgent:   "worker-a",
				Kind:      "status_update",
				Body:      "ja terminei a busca",
			},
		},
	}
	tool := NewReadTeamInboxTool(manager)

	ctx := agent.WithAgentContext(
		agent.WithTaskContext(context.Background(), "team-1", "task-1"),
		"worker-a",
	)

	result, err := tool.Execute(ctx, map[string]interface{}{"limit": float64(5)})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	var messages []agent.MailMessage
	if err := json.Unmarshal([]byte(result), &messages); err != nil {
		t.Fatalf("expected JSON payload, got %q error=%v", result, err)
	}
	if len(messages) != 1 || messages[0].FromAgent != "worker-b" {
		t.Fatalf("unexpected inbox payload: %#v", messages)
	}
}


