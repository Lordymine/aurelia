package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/kocar/aurelia/internal/agents"
	"github.com/kocar/aurelia/internal/bridge"
	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/cron"
	"github.com/kocar/aurelia/internal/memory"
	"github.com/kocar/aurelia/internal/persona"
	"github.com/kocar/aurelia/internal/runtime"
	"github.com/kocar/aurelia/internal/telegram"
	"github.com/kocar/aurelia/pkg/stt"
)

type app struct {
	resolver  *runtime.PathResolver
	bridge    *bridge.Bridge
	memory    *memory.Store
	agents    *agents.Registry
	cronStore *cron.SQLiteCronStore
	bot       *telegram.BotController
	scheduler *cron.Scheduler
	cronCtx   context.Context
	cronCancel context.CancelFunc
}

func bootstrapApp() (*app, error) {
	// 1. Resolve instance root and bootstrap directory tree
	resolver, err := runtime.New()
	if err != nil {
		return nil, fmt.Errorf("resolve instance root: %w", err)
	}
	if err := runtime.Bootstrap(resolver); err != nil {
		return nil, fmt.Errorf("bootstrap instance directory: %w", err)
	}

	// 2. Load config
	cfg, err := config.Load(resolver)
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	// 3. Set provider env vars for Bridge
	setProviderEnv(cfg)

	// 4. Create Bridge
	bridgeDir := findBridgeDir()
	br := bridge.New(bridgeDir)

	// 5. Create Embedder
	embedder := createEmbedder(cfg)

	// 6. Create Memory Store
	memStore, err := memory.NewStore(resolver.DBPath("memory.db"), embedder)
	if err != nil {
		return nil, fmt.Errorf("initialize memory store: %w", err)
	}

	// 7. Load Agent Registry
	agentReg, err := agents.Load(resolver.Agents())
	if err != nil {
		log.Printf("Warning: failed to load agents registry: %v (continuing without agents)", err)
		agentReg = nil
	}

	// 8. Build Persona service
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

	personaSvc := persona.NewCanonicalIdentityService(
		filepath.Join(personasDir, "IDENTITY.md"),
		filepath.Join(personasDir, "SOUL.md"),
		filepath.Join(personasDir, "USER.md"),
		ownerPlaybookPath,
		lessonsLearnedPath,
		projectPlaybookPath,
	)

	// 9. Create Cron Store
	cronStore, err := cron.NewSQLiteCronStore(resolver.DBPath("cron.db"))
	if err != nil {
		_ = memStore.Close()
		return nil, fmt.Errorf("initialize cron store: %w", err)
	}

	// 10. Create STT transcriber
	transcriber, err := buildTranscriber(cfg)
	if err != nil {
		_ = memStore.Close()
		_ = cronStore.Close()
		return nil, fmt.Errorf("initialize transcriber: %w", err)
	}

	// 11. Create Telegram BotController
	bot, err := telegram.NewBotController(
		cfg,
		br,
		agentReg,
		memStore,
		personaSvc,
		transcriber,
		personasDir,
	)
	if err != nil {
		_ = memStore.Close()
		_ = cronStore.Close()
		return nil, fmt.Errorf("initialize telegram bot: %w", err)
	}

	// 12. Create Cron Scheduler with BridgeCronRuntime
	var scheduler *cron.Scheduler
	if agentReg != nil {
		cronRuntime := cron.NewBridgeCronRuntime(
			&cron.BridgeAdapter{B: br},
			agentReg,
			personaSvc,
			memStore,
		)
		scheduler, err = cron.NewScheduler(cronStore, cronRuntime, nil, cron.SchedulerConfig{
			PollInterval: time.Minute,
		})
		if err != nil {
			_ = memStore.Close()
			_ = cronStore.Close()
			return nil, fmt.Errorf("initialize cron scheduler: %w", err)
		}

		// 13. Register scheduled agents from registry
		registerScheduledAgents(cronStore, agentReg)
	}

	cronCtx, cronCancel := context.WithCancel(context.Background())

	return &app{
		resolver:   resolver,
		bridge:     br,
		memory:     memStore,
		agents:     agentReg,
		cronStore:  cronStore,
		bot:        bot,
		scheduler:  scheduler,
		cronCtx:    cronCtx,
		cronCancel: cronCancel,
	}, nil
}

func (a *app) start() {
	if a.scheduler != nil {
		go func() {
			if err := a.scheduler.Start(a.cronCtx); err != nil && err != context.Canceled {
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
	if a.memory != nil {
		if err := a.memory.Close(); err != nil {
			log.Printf("Warning: failed to close memory store: %v", err)
		}
	}
	if a.cronStore != nil {
		if err := a.cronStore.Close(); err != nil {
			log.Printf("Warning: failed to close cron store: %v", err)
		}
	}
}

// setProviderEnv exports provider credentials as env vars for the Bridge process.
func setProviderEnv(cfg *config.AppConfig) {
	provider := cfg.DefaultProvider
	authMode := cfg.ProviderAuthMode(provider)

	// Subscription mode (Anthropic Max): SDK uses OAuth from ~/.claude/.credentials.json
	if provider == "anthropic" && authMode == "subscription" {
		os.Unsetenv("ANTHROPIC_API_KEY")
		os.Unsetenv("ANTHROPIC_BASE_URL")
		credPath := filepath.Join(os.Getenv("HOME"), ".claude", ".credentials.json")
		if _, err := os.Stat(credPath); os.IsNotExist(err) {
			log.Println("No Claude credentials found. Running 'claude login'...")
			cmd := exec.Command("claude", "login")
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				log.Fatalf("Claude login failed: %v. Run 'claude login' manually.", err)
			}
		}
		return
	}

	apiKey := cfg.ProviderAPIKey(provider)
	baseURL := cfg.ProviderBaseURL(provider)

	if apiKey != "" {
		os.Setenv("ANTHROPIC_API_KEY", apiKey)
	}
	if baseURL != "" {
		os.Setenv("ANTHROPIC_BASE_URL", baseURL)
		// Kimi and other non-Anthropic providers may not support tool_search
		os.Setenv("ENABLE_TOOL_SEARCH", "false")
	}
}

// findBridgeDir locates the bridge/ directory containing index.ts.
func findBridgeDir() string {
	candidates := []string{
		"bridge",
		filepath.Join(filepath.Dir(os.Args[0]), "bridge"),
	}
	for _, c := range candidates {
		if _, err := os.Stat(filepath.Join(c, "index.ts")); err == nil {
			return c
		}
	}
	return "bridge" // fallback
}

// createEmbedder builds the embedding provider from config.
func createEmbedder(cfg *config.AppConfig) memory.Embedder {
	apiKey := cfg.EmbeddingAPIKey
	if apiKey == "" {
		// Only use provider key if embedding provider is explicitly set
		if cfg.EmbeddingProvider != "" {
			apiKey = cfg.ProviderAPIKey(cfg.EmbeddingProvider)
		}
	}
	if apiKey == "" {
		log.Println("No embedding API key configured — using local word-hash embeddings")
		return memory.NewMockEmbedder(256)
	}
	model := cfg.EmbeddingModel
	if model == "" {
		model = "voyage-3"
	}
	return memory.NewVoyageEmbedder(apiKey, model)
}

func buildTranscriber(cfg *config.AppConfig) (stt.Transcriber, error) {
	switch cfg.STTProvider {
	case "", "groq":
		return stt.NewGroqTranscriber(cfg.ProviderAPIKey("groq")), nil
	default:
		return nil, fmt.Errorf("unsupported stt provider %q", cfg.STTProvider)
	}
}

// registerScheduledAgents syncs agent schedules into the cron store.
// Uses a deterministic job ID derived from agent name so that restarts
// skip agents that already have a job registered (idempotent).
func registerScheduledAgents(store *cron.SQLiteCronStore, reg *agents.Registry) {
	if reg == nil {
		return
	}
	svc := cron.NewService(store, nil)
	for _, a := range reg.Scheduled() {
		jobID := "scheduled-agent-" + a.Name

		// Skip if a job with this ID already exists.
		existing, err := store.GetJob(context.Background(), jobID)
		if err != nil {
			log.Printf("Warning: failed to check existing job for agent %q: %v", a.Name, err)
			continue
		}
		if existing != nil {
			log.Printf("Scheduled agent %q already registered (job %s), skipping", a.Name, jobID)
			continue
		}

		_, err = svc.CreateJob(context.Background(), cron.CronJob{
			ID:           jobID,
			AgentName:    a.Name,
			ScheduleType: "cron",
			CronExpr:     a.Schedule,
			Prompt:       a.Prompt,
		})
		if err != nil {
			log.Printf("Warning: failed to register scheduled agent %q: %v", a.Name, err)
		}
	}
}
