package telegram

import (
	"fmt"
	"log"
	"sync"
	"time"

	"gopkg.in/telebot.v3"

	"github.com/kocar/aurelia/internal/agents"
	"github.com/kocar/aurelia/internal/bridge"
	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/persona"
	"github.com/kocar/aurelia/internal/session"
	"github.com/kocar/aurelia/pkg/stt"
)

// BotController wires Telegram I/O to the application services.
type BotController struct {
	bot              *telebot.Bot
	config           *config.AppConfig
	bridge           *bridge.Bridge
	agents           *agents.Registry
	persona          *persona.CanonicalIdentityService
	stt              stt.Transcriber
	cronHandler      *CronCommandHandler
	sessions         *session.Store
	tracker          *session.Tracker
	personasDir      string
	exePath          string // path to aurelia binary for CLI instructions in system prompt
	bootstrapMu      sync.Mutex
	pendingBootstrap map[int64]bootstrapState
	albums           *albumBuffer
	bridgeFailures   bridgeFailureTracker
}

type albumBuffer struct {
	mu      sync.Mutex
	pending map[string]*pendingAlbum
}

func newAlbumBuffer() *albumBuffer {
	return &albumBuffer{
		pending: make(map[string]*pendingAlbum),
	}
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

// NewBotController builds the Telegram controller.
func NewBotController(
	cfg *config.AppConfig,
	br *bridge.Bridge,
	ag *agents.Registry,
	p *persona.CanonicalIdentityService,
	s stt.Transcriber,
	cronHandler *CronCommandHandler,
	personasDir string,
	exePath string,
	sessions *session.Store,
	tracker *session.Tracker,
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
		bridge:           br,
		agents:           ag,
		persona:          p,
		stt:              s,
		cronHandler:      cronHandler,
		sessions:         sessions,
		tracker:          tracker,
		personasDir:      personasDir,
		exePath:          exePath,
		pendingBootstrap: make(map[int64]bootstrapState),
		albums:           newAlbumBuffer(),
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
	bc.registerSlashMenu()
}
