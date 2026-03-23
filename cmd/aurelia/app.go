package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/telebot.v3"

	"github.com/kocar/aurelia/internal/agents"
	"github.com/kocar/aurelia/internal/bridge"
	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/cron"
	"github.com/kocar/aurelia/internal/persona"
	"github.com/kocar/aurelia/internal/runtime"
	"github.com/kocar/aurelia/internal/session"
	"github.com/kocar/aurelia/internal/telegram"
	"github.com/kocar/aurelia/pkg/stt"
)

type app struct {
	resolver   *runtime.PathResolver
	bridge     *bridge.Bridge
	agents     *agents.Registry
	cronStore  *cron.SQLiteCronStore
	bot        *telegram.BotController
	scheduler  *cron.Scheduler
	cronCtx    context.Context
	cronCancel context.CancelFunc
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
	setProviderEnv(cfg)

	br := setupBridge()
	personaSvc := setupPersona(resolver)

	agentReg, err := agents.Load(resolver.Agents())
	if err != nil {
		log.Printf("Warning: failed to load agents registry: %v (continuing without agents)", err)
		agentReg = nil
	}

	cronStore, err := cron.NewSQLiteCronStore(resolver.DBPath("cron.db"))
	if err != nil {
		return nil, fmt.Errorf("initialize cron store: %w", err)
	}

	transcriber, err := buildTranscriber(cfg)
	if err != nil {
		if closeErr := cronStore.Close(); closeErr != nil {
			log.Printf("Warning: failed to close cron store: %v", closeErr)
		}
		return nil, fmt.Errorf("initialize transcriber: %w", err)
	}

	cronSvc := cron.NewService(cronStore, nil)
	cronHandler := telegram.NewCronCommandHandler(cronSvc)
	exePath, err := os.Executable()
	if err != nil {
		log.Printf("Warning: failed to resolve executable path: %v", err)
	}
	sessions := session.NewStore()
	tracker := session.NewTracker()

	bot, err := telegram.NewBotController(
		cfg, br, agentReg, personaSvc, transcriber,
		cronHandler, resolver.MemoryPersonas(), exePath, sessions, tracker,
	)
	if err != nil {
		if closeErr := cronStore.Close(); closeErr != nil {
			log.Printf("Warning: failed to close cron store: %v", closeErr)
		}
		return nil, fmt.Errorf("initialize telegram bot: %w", err)
	}

	scheduler, err := setupCronScheduler(cronStore, br, agentReg, personaSvc, bot)
	if err != nil {
		if closeErr := cronStore.Close(); closeErr != nil {
			log.Printf("Warning: failed to close cron store: %v", closeErr)
		}
		return nil, fmt.Errorf("initialize cron scheduler: %w", err)
	}

	cronCtx, cronCancel := context.WithCancel(context.Background())

	return &app{
		resolver:   resolver,
		bridge:     br,
		agents:     agentReg,
		cronStore:  cronStore,
		bot:        bot,
		scheduler:  scheduler,
		cronCtx:    cronCtx,
		cronCancel: cronCancel,
	}, nil
}

// setupBridge creates the Bridge, ensuring ~/.aurelia/bridge/ is bootstrapped.
func setupBridge() *bridge.Bridge {
	home, _ := os.UserHomeDir()
	aureliBridgeDir := filepath.Join(home, ".aurelia", "bridge")
	if _, setupErr := bridge.EnsureBridge(aureliBridgeDir, bridge.EmbeddedBundleJS); setupErr != nil {
		log.Printf("Warning: bridge auto-setup failed: %v", setupErr)
	}
	bridgeDir := findBridgeDir()
	if bridgeDir == "" {
		bridgeDir = aureliBridgeDir
	}
	bundlePath := filepath.Join(bridgeDir, "bundle.js")
	if _, err := os.Stat(bundlePath); os.IsNotExist(err) {
		bundlePath = ""
	}
	return bridge.New(bridgeDir, bundlePath)
}

// setupPersona builds the canonical identity service from persona and playbook files.
func setupPersona(resolver *runtime.PathResolver) *persona.CanonicalIdentityService {
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

	return persona.NewCanonicalIdentityService(
		filepath.Join(personasDir, "IDENTITY.md"),
		filepath.Join(personasDir, "SOUL.md"),
		filepath.Join(personasDir, "USER.md"),
		ownerPlaybookPath,
		lessonsLearnedPath,
		projectPlaybookPath,
	)
}

// telegramChatSender adapts a telebot.Bot to the cron.ChatSender interface.
type telegramChatSender struct {
	bot *telebot.Bot
}

func (s *telegramChatSender) Send(chatID int64, text string) error {
	chat := &telebot.Chat{ID: chatID}
	return telegram.SendText(s.bot, chat, text)
}

// setupCronScheduler creates the cron scheduler with Telegram delivery.
// Returns nil scheduler if agentReg is nil.
func setupCronScheduler(
	cronStore *cron.SQLiteCronStore,
	br *bridge.Bridge,
	agentReg *agents.Registry,
	personaSvc *persona.CanonicalIdentityService,
	bot *telegram.BotController,
) (*cron.Scheduler, error) {
	if agentReg == nil {
		return nil, nil
	}

	cronRuntime := cron.NewBridgeCronRuntime(
		&cron.BridgeAdapter{B: br},
		agentReg,
		personaSvc,
	)

	delivery := cron.NewTelegramDelivery(&telegramChatSender{bot: bot.GetBot()})
	deliverFn := func(ctx context.Context, job cron.CronJob, result *cron.ExecutionResult, execErr error) error {
		return delivery.Deliver(ctx, job, result, execErr)
	}

	notifyingRuntime := cron.NewNotifyingRuntime(cronRuntime, deliverFn)
	scheduler, err := cron.NewScheduler(cronStore, notifyingRuntime, nil, cron.SchedulerConfig{
		PollInterval: 15 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	registerScheduledAgents(cronStore, agentReg)
	return scheduler, nil
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
		done := make(chan struct{})
		go func() {
			a.bot.Stop()
			close(done)
		}()
		select {
		case <-done:
		case <-ctx.Done():
			log.Println("Warning: bot shutdown timed out")
		}
	}
}

func (a *app) close() {
	if a.bridge != nil {
		a.bridge.Stop()
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
		home, _ := os.UserHomeDir()
		credPath := filepath.Join(home, ".claude", ".credentials.json")
		if _, err := os.Stat(credPath); os.IsNotExist(err) {
			log.Fatalf("Anthropic subscription requires Claude login. Run 'claude login' first.")
		}
		return
	}

	apiKey := cfg.ProviderAPIKey(provider)
	baseURL := cfg.ProviderBaseURL(provider)

	// Auto-set base URL for known providers if not explicitly configured
	if baseURL == "" {
		switch config.NormalizeProvider(provider) {
		case "kimi":
			baseURL = "https://api.kimi.com/coding/"
		case "openrouter":
			baseURL = "https://openrouter.ai/api/v1"
		case "zai":
			baseURL = "https://api.z.ai/api/anthropic"
		case "alibaba":
			baseURL = "https://dashscope-intl.aliyuncs.com/apps/anthropic"
		}
	}

	if apiKey != "" {
		os.Setenv("ANTHROPIC_API_KEY", apiKey)
	}
	if baseURL != "" {
		os.Setenv("ANTHROPIC_BASE_URL", baseURL)
		os.Setenv("ENABLE_TOOL_SEARCH", "false")
	}
}

// findBridgeDir locates the bridge/ directory containing bundle.js or index.ts.
func findBridgeDir() string {
	home, _ := os.UserHomeDir()
	candidates := []string{
		"bridge",                                          // relative to cwd (development)
		filepath.Join(filepath.Dir(os.Args[0]), "bridge"), // next to executable
		filepath.Join(home, ".aurelia", "bridge"),          // user data dir
	}
	for _, c := range candidates {
		if _, err := os.Stat(filepath.Join(c, "bundle.js")); err == nil {
			return c
		}
		if _, err := os.Stat(filepath.Join(c, "index.ts")); err == nil {
			return c
		}
	}
	return "" // not found — triggers auto-setup
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
