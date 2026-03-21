package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/cron"
	"github.com/kocar/aurelia/internal/persona"
	"github.com/kocar/aurelia/internal/runtime"
	"github.com/kocar/aurelia/internal/telegram"
	"github.com/kocar/aurelia/pkg/stt"
)

type app struct {
	resolver      *runtime.PathResolver
	cronStore     *cron.SQLiteCronStore
	bot           *telegram.BotController
	cronScheduler *cron.Scheduler
	cronCtx       context.Context
	cronCancel    context.CancelFunc
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
	if cwd != "" {
		projectPlaybookPath = filepath.Join(cwd, "docs", "PROJECT_PLAYBOOK.md")
	}

	canonicalService := persona.NewCanonicalIdentityService(
		filepath.Join(personasDir, "IDENTITY.md"),
		filepath.Join(personasDir, "SOUL.md"),
		filepath.Join(personasDir, "USER.md"),
		ownerPlaybookPath,
		lessonsLearnedPath,
		projectPlaybookPath,
	)

	cronStore, err := cron.NewSQLiteCronStore(cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("initialize cron store: %w", err)
	}

	transcriber, err := buildTranscriber(cfg)
	if err != nil {
		_ = cronStore.Close()
		return nil, fmt.Errorf("initialize transcriber: %w", err)
	}

	// TODO(task-10): wire bridge, agents registry, and memory store
	bot, err := telegram.NewBotController(
		cfg,
		nil, // bridge — wired in task 10
		nil, // agents registry — wired in task 10
		nil, // memory store — wired in task 10
		canonicalService,
		transcriber,
		personasDir,
	)
	if err != nil {
		_ = cronStore.Close()
		return nil, fmt.Errorf("initialize telegram block: %w", err)
	}

	// TODO: wire cron scheduler with bridge executor
	_ = cronStore
	_ = bot

	return &app{
		resolver:  resolver,
		cronStore: cronStore,
		bot:       bot,
	}, nil
}

func buildTranscriber(cfg *config.AppConfig) (stt.Transcriber, error) {
	switch cfg.STTProvider {
	case "", "groq":
		return stt.NewGroqTranscriber(cfg.ProviderAPIKey("groq")), nil
	default:
		return nil, fmt.Errorf("unsupported stt provider %q", cfg.STTProvider)
	}
}

func (a *app) start() {
	if a.cronScheduler != nil {
		go func() {
			if err := a.cronScheduler.Start(a.cronCtx); err != nil && err != context.Canceled {
				log.Printf("Warning: cron scheduler stopped with error: %v", err)
			}
		}()
	}
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
	if a.cronStore != nil {
		if err := a.cronStore.Close(); err != nil {
			log.Printf("Warning: failed to close cron store: %v", err)
		}
	}
}
