package agent

import "context"

type executionContextKey string

const (
	teamKeyContextKey  executionContextKey = "team_key"
	userKeyContextKey  executionContextKey = "user_key"
	teamIDContextKey   executionContextKey = "team_id"
	taskIDContextKey   executionContextKey = "task_id"
	runIDContextKey    executionContextKey = "run_id"
	agentKeyContextKey executionContextKey = "agent_name"
	workdirContextKey  executionContextKey = "workdir"
)

func WithTeamContext(ctx context.Context, teamKey, userID string) context.Context {
	ctx = context.WithValue(ctx, teamKeyContextKey, teamKey)
	return context.WithValue(ctx, userKeyContextKey, userID)
}

func TeamContextFromContext(ctx context.Context) (teamKey string, userID string, ok bool) {
	if ctx == nil {
		return "", "", false
	}

	teamKey, teamOK := ctx.Value(teamKeyContextKey).(string)
	userID, userOK := ctx.Value(userKeyContextKey).(string)
	if !teamOK || !userOK || teamKey == "" || userID == "" {
		return "", "", false
	}

	return teamKey, userID, true
}

func WithTaskContext(ctx context.Context, teamID, taskID string) context.Context {
	ctx = context.WithValue(ctx, teamIDContextKey, teamID)
	return context.WithValue(ctx, taskIDContextKey, taskID)
}

func TaskContextFromContext(ctx context.Context) (teamID string, taskID string, ok bool) {
	if ctx == nil {
		return "", "", false
	}

	teamID, teamOK := ctx.Value(teamIDContextKey).(string)
	taskID, taskOK := ctx.Value(taskIDContextKey).(string)
	if !teamOK || !taskOK || teamID == "" || taskID == "" {
		return "", "", false
	}

	return teamID, taskID, true
}

func WithRunContext(ctx context.Context, runID string) context.Context {
	return context.WithValue(ctx, runIDContextKey, runID)
}

func WithAgentContext(ctx context.Context, agentName string) context.Context {
	return context.WithValue(ctx, agentKeyContextKey, agentName)
}

func WithWorkdirContext(ctx context.Context, workdir string) context.Context {
	return context.WithValue(ctx, workdirContextKey, workdir)
}

func AgentContextFromContext(ctx context.Context) (agentName string, ok bool) {
	if ctx == nil {
		return "", false
	}

	agentName, ok = ctx.Value(agentKeyContextKey).(string)
	if !ok || agentName == "" {
		return "", false
	}

	return agentName, true
}

func RunContextFromContext(ctx context.Context) (runID string, ok bool) {
	if ctx == nil {
		return "", false
	}

	runID, ok = ctx.Value(runIDContextKey).(string)
	if !ok || runID == "" {
		return "", false
	}

	return runID, true
}

func WorkdirFromContext(ctx context.Context) (workdir string, ok bool) {
	if ctx == nil {
		return "", false
	}

	workdir, ok = ctx.Value(workdirContextKey).(string)
	if !ok || workdir == "" {
		return "", false
	}

	return workdir, true
}
