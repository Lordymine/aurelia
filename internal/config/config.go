package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kocar/aurelia/internal/runtime"
)

// AppConfig holds all the environment variables needed for the application
type AppConfig struct {
	TelegramBotToken       string
	TelegramAllowedUserIDs []int64
	KimiAPIKey             string
	GroqAPIKey             string
	MaxIterations          int
	DBPath                 string
	MemoryWindowSize       int
	MCPConfigPath          string
}

// Load reads from .env and returns an AppConfig. The resolver is used to
// provide instance-directory defaults for DBPath and MCPConfigPath when the
// corresponding env vars are not set.
func Load(r *runtime.PathResolver) (*AppConfig, error) {
	// We ignore the error here because in production .env might not exist
	// and we might rely on system environment variables directly.
	_ = godotenv.Load()

	// Parse allowed users
	usersStr := os.Getenv("TELEGRAM_ALLOWED_USER_IDS")
	usersSplit := strings.Split(usersStr, ",")
	var allowedUsers []int64
	for _, u := range usersSplit {
		u = strings.TrimSpace(u)
		if u == "" {
			continue
		}
		id, err := strconv.ParseInt(u, 10, 64)
		if err != nil {
			log.Printf("Warning: failed to parse user ID '%s': %v\n", u, err)
			continue
		}
		allowedUsers = append(allowedUsers, id)
	}

	maxIters := 5
	if m := os.Getenv("MAX_ITERATIONS"); m != "" {
		if val, err := strconv.Atoi(m); err == nil {
			maxIters = val
		}
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = filepath.Join(r.Data(), "aurelia.db")
	}

	memWinSize := 20
	if mw := os.Getenv("MEMORY_WINDOW_SIZE"); mw != "" {
		if val, err := strconv.Atoi(mw); err == nil {
			memWinSize = val
		}
	}

	mcpConfigPath := os.Getenv("MCP_SERVERS_CONFIG_PATH")
	if mcpConfigPath == "" {
		mcpConfigPath = filepath.Join(r.Config(), "mcp_servers.json")
	}

	cfg := &AppConfig{
		TelegramBotToken:       os.Getenv("TELEGRAM_BOT_TOKEN"),
		TelegramAllowedUserIDs: allowedUsers,
		KimiAPIKey:             os.Getenv("KIMI_API_KEY"),
		GroqAPIKey:             os.Getenv("GROQ_API_KEY"),
		MaxIterations:          maxIters,
		DBPath:                 dbPath,
		MemoryWindowSize:       memWinSize,
		MCPConfigPath:          mcpConfigPath,
	}

	return cfg, nil
}


