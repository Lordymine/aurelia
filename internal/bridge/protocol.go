package bridge

// TelegramPluginTools lists the Telegram MCP plugin tools that should be
// blocked when Aurelia is the Telegram bot (to prevent duplicate delivery).
var TelegramPluginTools = []string{
	"mcp__plugin_telegram_telegram__reply",
	"mcp__plugin_telegram_telegram__react",
	"mcp__plugin_telegram_telegram__edit_message",
}

// Request sent to Bridge process via stdin as JSON.
type Request struct {
	Command   string         `json:"command"`
	Prompt    string         `json:"prompt,omitempty"`
	RequestID string         `json:"request_id,omitempty"`
	Options   RequestOptions `json:"options,omitempty"`
}

// RequestOptions configures how the Bridge executes a query.
type RequestOptions struct {
	Model          string         `json:"model,omitempty"`
	Cwd            string         `json:"cwd,omitempty"`
	SystemPrompt   string         `json:"system_prompt,omitempty"`
	Resume         string         `json:"resume,omitempty"`
	MaxTurns       int            `json:"max_turns,omitempty"`
	PermissionMode string         `json:"permission_mode,omitempty"`
	MCPServers     map[string]any `json:"mcp_servers,omitempty"`
	AllowedTools   []string       `json:"allowed_tools,omitempty"`
	Continue       bool           `json:"continue,omitempty"`
	Agents         map[string]any `json:"agents,omitempty"`
	NoUserSettings bool           `json:"no_user_settings,omitempty"`
	DisabledTools  []string       `json:"disabled_tools,omitempty"`
}
