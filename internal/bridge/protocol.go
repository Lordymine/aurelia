package bridge

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
}
