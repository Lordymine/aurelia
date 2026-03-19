package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/kocar/aurelia/internal/agent"
	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/cron"
	"github.com/kocar/aurelia/internal/mcp"
	"github.com/kocar/aurelia/internal/observability"
	"github.com/kocar/aurelia/internal/telegram"
	"github.com/kocar/aurelia/internal/tools"
	"gopkg.in/telebot.v3"
)

func buildToolRegistry() *agent.ToolRegistry {
	registry := agent.NewToolRegistry()
	tools.RegisterCoreTools(registry)
	return registry
}

func registerScheduleTools(registry *agent.ToolRegistry, cronStore *cron.SQLiteCronStore) *cron.Service {
	cronService := cron.NewService(cronStore, nil)
	createScheduleTool := tools.NewCreateScheduleTool(cronService)
	registry.Register(createScheduleTool.Definition(), createScheduleTool.Execute)
	listSchedulesTool := tools.NewListSchedulesTool(cronService)
	registry.Register(listSchedulesTool.Definition(), listSchedulesTool.Execute)
	pauseScheduleTool := tools.NewPauseScheduleTool(cronService)
	registry.Register(pauseScheduleTool.Definition(), pauseScheduleTool.Execute)
	resumeScheduleTool := tools.NewResumeScheduleTool(cronService)
	registry.Register(resumeScheduleTool.Definition(), resumeScheduleTool.Execute)
	deleteScheduleTool := tools.NewDeleteScheduleTool(cronService)
	registry.Register(deleteScheduleTool.Definition(), deleteScheduleTool.Execute)
	return cronService
}

func maybeRegisterMCPTools(cfg *config.AppConfig, registry *agent.ToolRegistry) (*mcp.Manager, error) {
	mcpPath := cfg.MCPConfigPath

	mcpCfg, err := config.LoadMCPConfig(mcpPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, loggableError("failed to load MCP config from %s: %v", mcpPath, err)
	}
	if mcpCfg == nil || !mcpCfg.Enabled {
		return nil, nil
	}

	cwd, _ := os.Getwd()
	mcpManager, err := mcp.NewManager(*mcpCfg, cwd)
	if err != nil {
		return nil, loggableError("failed to start MCP Manager: %v", err)
	}

	tools.RegisterMCPTools(registry, mcpManager)
	log.Printf("MCP Manager initialized with %d tools.", len(mcpManager.ToolSpecs()))
	return mcpManager, nil
}

func registerSpawnAgentTool(
	cfg *config.AppConfig,
	registry *agent.ToolRegistry,
	llmProvider agent.LLMProvider,
	bot *telegram.BotController,
	taskStore *agent.SQLiteTaskStore,
	observer observability.Recorder,
) error {
	teamManager, err := agent.NewTeamManagerWithObserver(taskStore, observer)
	if err != nil {
		return loggableError("initialize team manager: %v", err)
	}

	masterTeams, err := agent.NewMasterTeamService(
		teamManager,
		llmProvider,
		registry,
		cfg.MaxIterations,
		func(teamKey string, message string) {
			notifyMasterTeam(bot, teamKey, message)
		},
	)
	if err != nil {
		return loggableError("initialize master team service: %v", err)
	}

	if err := masterTeams.Rehydrate(context.Background()); err != nil {
		log.Printf("Warning: failed to rehydrate master teams: %v", err)
	}

	spawnAgent := tools.NewSpawnAgentTool(masterTeams)
	registry.Register(spawnAgent.Definition(), spawnAgent.Execute)
	pauseTeam := tools.NewPauseTeamTool(masterTeams)
	registry.Register(pauseTeam.Definition(), pauseTeam.Execute)
	resumeTeam := tools.NewResumeTeamTool(masterTeams)
	registry.Register(resumeTeam.Definition(), resumeTeam.Execute)
	cancelTeam := tools.NewCancelTeamTool(masterTeams)
	registry.Register(cancelTeam.Definition(), cancelTeam.Execute)
	teamStatus := tools.NewTeamStatusTool(masterTeams)
	registry.Register(teamStatus.Definition(), teamStatus.Execute)
	sendTeamMessage := tools.NewSendTeamMessageTool(teamManager)
	registry.Register(sendTeamMessage.Definition(), sendTeamMessage.Execute)
	readTeamInbox := tools.NewReadTeamInboxTool(teamManager)
	registry.Register(readTeamInbox.Definition(), readTeamInbox.Execute)
	return nil
}

func notifyMasterTeam(bot *telegram.BotController, teamKey, message string) {
	chatID, err := strconv.ParseInt(teamKey, 10, 64)
	if err != nil {
		log.Printf("Warning: invalid team key %q for master notification: %v", teamKey, err)
		return
	}

	if err := telegram.SendText(bot.GetBot(), &telebot.Chat{ID: chatID}, message); err != nil {
		log.Printf("Warning: failed to send master team update: %v", err)
	}
}

func loggableError(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}
