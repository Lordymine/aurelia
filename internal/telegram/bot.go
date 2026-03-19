package telegram

import (
	"fmt"
	"log"
	"sync"
	"time"

	"gopkg.in/telebot.v3"

	"github.com/kocar/aurelia/internal/agent"
	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/memory"
	"github.com/kocar/aurelia/internal/persona"
	"github.com/kocar/aurelia/internal/skill"
	"github.com/kocar/aurelia/pkg/stt"
)

// BotController wires Telegram I/O to the application services.
type BotController struct {
	bot              *telebot.Bot
	config           *config.AppConfig
	memory           *memory.MemoryManager
	router           *skill.Router
	executor         *skill.Executor
	loader           *skill.Loader
	stt              stt.Transcriber
	canonical        *persona.CanonicalIdentityService
	bootstrapMu      sync.Mutex
	pendingBootstrap map[int64]bootstrapState
	albumMu          sync.Mutex
	pendingAlbums    map[string]*pendingAlbum
	mediaMu          sync.Mutex
	recentMedia      map[string]recentMedia
	personasDir      string
}

type pendingAlbum struct {
	ownerMessageID int
	caption        string
	photos         []albumPhoto
}

type albumPhoto struct {
	messageID int
	photo     telebot.Photo
}

type recentMedia struct {
	parts     []agent.ContentPart
	updatedAt time.Time
}

// NewBotController builds the Telegram controller.
func NewBotController(
	cfg *config.AppConfig,
	mem *memory.MemoryManager,
	r *skill.Router,
	e *skill.Executor,
	l *skill.Loader,
	s stt.Transcriber,
	canonical *persona.CanonicalIdentityService,
	personasDir string,
) (*BotController, error) {

	pref := telebot.Settings{
		Token:  cfg.TelegramBotToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	bc := &BotController{
		bot:              b,
		config:           cfg,
		memory:           mem,
		router:           r,
		executor:         e,
		loader:           l,
		stt:              s,
		canonical:        canonical,
		pendingBootstrap: make(map[int64]bootstrapState),
		pendingAlbums:    make(map[string]*pendingAlbum),
		recentMedia:      make(map[string]recentMedia),
		personasDir:      personasDir,
	}

	bc.setupRoutes()
	return bc, nil
}

// GetBot exposes the underlying Telebot instance.
func (bc *BotController) GetBot() *telebot.Bot {
	return bc.bot
}

// Start begins Telegram polling.
func (bc *BotController) Start() {
	log.Println("Starting Aurelia Telegram Bot...")
	bc.bot.Start()
}

// Stop ends Telegram polling.
func (bc *BotController) Stop() {
	bc.bot.Stop()
}

func (bc *BotController) isAllowedUser(userID int64) bool {
	if bc == nil || bc.config == nil {
		return false
	}
	for _, id := range bc.config.TelegramAllowedUserIDs {
		if id == userID {
			return true
		}
	}
	return false
}

func (bc *BotController) setupRoutes() {
	bc.bot.Use(bc.whitelistMiddleware())

	bc.setupBootstrapRoutes()
	bc.registerContentRoutes()
}
