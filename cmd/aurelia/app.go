package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/kocar/aurelia/internal/agent"
	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/cron"
	"github.com/kocar/aurelia/internal/mcp"
	"github.com/kocar/aurelia/internal/memory"
	"github.com/kocar/aurelia/internal/observability"
	"github.com/kocar/aurelia/internal/persona"
	"github.com/kocar/aurelia/internal/runtime"
	"github.com/kocar/aurelia/internal/skill"
	"github.com/kocar/aurelia/internal/telegram"
	"github.com/kocar/aurelia/internal/tools"
	"github.com/kocar/aurelia/pkg/llm"
	"github.com/kocar/aurelia/pkg/stt"
	"gopkg.in/telebot.v3"
)

type app struct {
	resolver      *runtime.PathResolver
	mem           *memory.MemoryManager
	cronStore     *cron.SQLiteCronStore
	opsStore      *observability.SQLiteStore
	mcpManager    *mcp.Manager
	taskStore     *agent.SQLiteTaskStore
	llmProvider   closableLLMProvider
	bot           *telegram.BotController
	cronScheduler *cron.Scheduler
	cronCtx       context.Context
	cronCancel    context.CancelFunc
}

type closableLLMProvider interface {
	agent.LLMProvider
	Close()
}

func bootstrapApp() (*app, error) {
	resolver, err := runtime.New()
	if err != nil {
		return nil, fmt.Errorf("resolve instance root: %w", err)
	}
	if err := runtime.Bootstrap(resolver); err != nil {
		return nil, fmt.Errorf("bootstrap instance directory: %w", err)
	}

	cfg, err := config.Load(resolver)
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	mem, err := memory.NewMemoryManager(cfg.DBPath, cfg.MemoryWindowSize)
	if err != nil {
		return nil, fmt.Errorf("initialize memory manager: %w", err)
	}
	opsStore, err := observability.NewSQLiteStore(cfg.DBPath)
	if err != nil {
		_ = mem.Close()
		return nil, fmt.Errorf("initialize observability store: %w", err)
	}

	personasDir := resolver.MemoryPersonas()
	memoryDir := resolver.Memory()
	ownerPlaybookPath := filepath.Join(memoryDir, "OWNER_PLAYBOOK.md")
	lessonsLearnedPath := filepath.Join(memoryDir, "LESSONS_LEARNED.md")
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Warning: failed to resolve working directory for project playbook: %v", err)
		cwd = ""
	}
	if err := runtime.BootstrapProject(cwd); err != nil {
		log.Printf("Warning: failed to bootstrap project-local Aurelia directory: %v", err)
	}
	var projectPlaybookPath string
	var projectSkillsDir string
	if cwd != "" {
		projectPlaybookPath = filepath.Join(cwd, "docs", "PROJECT_PLAYBOOK.md")
		projectSkillsDir = runtime.ProjectSkills(cwd)
	}
	llmProvider, err := buildLLMProvider(cfg, resolver)
	if err != nil {
		_ = opsStore.Close()
		_ = mem.Close()
		return nil, fmt.Errorf("initialize llm provider: %w", err)
	}
	canonicalService := persona.NewCanonicalIdentityService(
		mem,
		filepath.Join(personasDir, "IDENTITY.md"),
		filepath.Join(personasDir, "SOUL.md"),
		filepath.Join(personasDir, "USER.md"),
		ownerPlaybookPath,
		lessonsLearnedPath,
		projectPlaybookPath,
	)
	registry := buildToolRegistry()

	cronStore, err := cron.NewSQLiteCronStore(cfg.DBPath)
	if err != nil {
		_ = opsStore.Close()
		_ = mem.Close()
		return nil, fmt.Errorf("initialize cron store: %w", err)
	}

	registerScheduleTools(registry, cronStore)
	mcpManager, err := maybeRegisterMCPTools(cfg, registry)
	if err != nil {
		log.Printf("Warning: %v", err)
	}

	loop := agent.NewLoopWithObserver(llmProvider, registry, cfg.MaxIterations, opsStore)
	skillLoader := skill.NewLoader(resolver.Skills(), projectSkillsDir)
	skillRouter := skill.NewRouter(llmProvider)
	skillExecutor := skill.NewExecutor(loop)
	skillInstaller := skill.NewInstaller(resolver.Skills(), projectSkillsDir)
	transcriber, err := buildTranscriber(cfg)
	if err != nil {
		_ = cronStore.Close()
		_ = opsStore.Close()
		_ = mem.Close()
		return nil, fmt.Errorf("initialize transcriber: %w", err)
	}

	installSkillTool := tools.NewInstallSkillTool(skillInstaller)
	registry.Register(installSkillTool.Definition(), installSkillTool.Execute)

	bot, err := telegram.NewBotController(
		cfg,
		mem,
		skillRouter,
		skillExecutor,
		skillLoader,
		transcriber,
		canonicalService,
		personasDir,
		opsStore,
	)
	if err != nil {
		_ = cronStore.Close()
		_ = opsStore.Close()
		_ = mem.Close()
		return nil, fmt.Errorf("initialize telegram block: %w", err)
	}

	taskStore, err := agent.NewSQLiteTaskStore(cfg.DBPath + ".teams")
	if err != nil {
		_ = cronStore.Close()
		_ = opsStore.Close()
		_ = mem.Close()
		return nil, fmt.Errorf("initialize team task store: %w", err)
	}

	if err := registerSpawnAgentTool(cfg, registry, llmProvider, bot, taskStore, opsStore); err != nil {
		_ = taskStore.Close()
		_ = cronStore.Close()
		_ = opsStore.Close()
		_ = mem.Close()
		return nil, err
	}

	cronScheduler, cronCtx, cronCancel, err := buildCronScheduler(cronStore, loop, canonicalService, bot)
	if err != nil {
		_ = taskStore.Close()
		_ = cronStore.Close()
		_ = opsStore.Close()
		_ = mem.Close()
		return nil, fmt.Errorf("initialize cron scheduler: %w", err)
	}

	return &app{
		resolver:      resolver,
		mem:           mem,
		cronStore:     cronStore,
		opsStore:      opsStore,
		mcpManager:    mcpManager,
		taskStore:     taskStore,
		llmProvider:   llmProvider,
		bot:           bot,
		cronScheduler: cronScheduler,
		cronCtx:       cronCtx,
		cronCancel:    cronCancel,
	}, nil
}

func buildLLMProvider(cfg *config.AppConfig, resolver *runtime.PathResolver) (closableLLMProvider, error) {
	_ = resolver
	switch cfg.LLMProvider {
	case "anthropic":
		return llm.NewAnthropicProvider(cfg.AnthropicAPIKey, cfg.LLMModel), nil
	case "google":
		return llm.NewGeminiProvider(context.Background(), cfg.GoogleAPIKey, cfg.LLMModel)
	case "kilo":
		return llm.NewKiloProvider(cfg.KiloAPIKey, cfg.LLMModel), nil
	case "openrouter":
		return llm.NewOpenRouterProvider(cfg.OpenRouterAPIKey, cfg.LLMModel), nil
	case "zai":
		return llm.NewZAIProvider(cfg.ZAIAPIKey, cfg.LLMModel), nil
	case "alibaba":
		return llm.NewAlibabaProvider(cfg.AlibabaAPIKey, cfg.LLMModel), nil
	case "openai":
		if cfg.OpenAIAuthMode == "codex" {
			if err := llm.EnsureCodexCLIAvailable(); err != nil {
				return nil, err
			}
			return llm.NewCodexCLIProvider(cfg.LLMModel)
		}
		return llm.NewOpenAIProvider(cfg.OpenAIAPIKey, cfg.LLMModel), nil
	case "", "kimi":
		return llm.NewKimiProvider(cfg.KimiAPIKey, cfg.LLMModel), nil
	default:
		return nil, fmt.Errorf("unsupported llm provider %q", cfg.LLMProvider)
	}
}

func buildTranscriber(cfg *config.AppConfig) (stt.Transcriber, error) {
	switch cfg.STTProvider {
	case "", "groq":
		return stt.NewGroqTranscriber(cfg.GroqAPIKey), nil
	default:
		return nil, fmt.Errorf("unsupported stt provider %q", cfg.STTProvider)
	}
}

func (a *app) start() {
	go func() {
		if err := a.cronScheduler.Start(a.cronCtx); err != nil && err != context.Canceled {
			log.Printf("Warning: cron scheduler stopped with error: %v", err)
		}
	}()
	go a.bot.Start()
}

func (a *app) shutdown(ctx context.Context) {
	if a.cronCancel != nil {
		a.cronCancel()
	}
	if a.bot != nil {
		a.bot.Stop()
	}
	_ = ctx
}

func (a *app) close() {
	if a.taskStore != nil {
		if err := a.taskStore.Close(); err != nil {
			log.Printf("Warning: failed to close team task store: %v", err)
		}
	}
	if a.mcpManager != nil {
		if err := a.mcpManager.Close(); err != nil {
			log.Printf("Warning: failed to close MCP manager: %v", err)
		}
	}
	if a.cronStore != nil {
		if err := a.cronStore.Close(); err != nil {
			log.Printf("Warning: failed to close cron store: %v", err)
		}
	}
	if a.opsStore != nil {
		if err := a.opsStore.Close(); err != nil {
			log.Printf("Warning: failed to close observability store: %v", err)
		}
	}
	if a.mem != nil {
		if err := a.mem.Close(); err != nil {
			log.Printf("Warning: failed to close memory manager: %v", err)
		}
	}
	if a.llmProvider != nil {
		a.llmProvider.Close()
	}
}

func buildCronScheduler(
	cronStore *cron.SQLiteCronStore,
	loop *agent.Loop,
	canonicalService *persona.CanonicalIdentityService,
	bot *telegram.BotController,
) (*cron.Scheduler, context.Context, context.CancelFunc, error) {
	cronPrompt, cronAllowedTools := loadCronPromptConfig(canonicalService)
	cronRuntime := cron.NewAgentCronRuntimeWithPromptBuilder(
		&loopExecutorAdapter{loop: loop},
		cronPrompt,
		cronAllowedTools,
		func(ctx context.Context, job cron.CronJob) (string, []string, error) {
			if canonicalService == nil {
				return cronPrompt, cronAllowedTools, nil
			}
			return canonicalService.BuildPromptForQuery(ctx, job.OwnerUserID, strconv.FormatInt(job.TargetChatID, 10), job.Prompt)
		},
	)

	notifyingRuntime := cron.NewNotifyingRuntime(cronRuntime, func(ctx context.Context, job cron.CronJob, output string, execErr error) error {
		chat := &telebot.Chat{ID: job.TargetChatID}
		if execErr != nil {
			return telegram.SendText(bot.GetBot(), chat, "Falha na rotina agendada:\n"+execErr.Error())
		}
		return telegram.SendText(bot.GetBot(), chat, output)
	})

	scheduler, err := cron.NewScheduler(cronStore, notifyingRuntime, nil, cron.SchedulerConfig{
		PollInterval: time.Minute,
		Observer:     bot.GetOpsStore(),
	})
	if err != nil {
		return nil, nil, nil, err
	}

	cronCtx, cancel := context.WithCancel(context.Background())
	return scheduler, cronCtx, cancel, nil
}
