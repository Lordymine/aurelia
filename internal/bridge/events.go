package bridge

// Event represents a single NDJSON event emitted by the Bridge on stdout.
// Not all fields are populated for every event type — only the fields relevant
// to the event's Type are set.
type Event struct {
	Type      string `json:"event"`
	RequestID string `json:"request_id,omitempty"`

	// system event
	SessionID string   `json:"session_id,omitempty"`
	Tools     []string `json:"tools,omitempty"`
	Model     string   `json:"model,omitempty"`

	// tool_use event
	Name  string `json:"name,omitempty"`
	Input any    `json:"input,omitempty"`

	// tool_result / assistant / result / error
	Content string `json:"content,omitempty"`
	Text    string `json:"text,omitempty"`
	Message string `json:"message,omitempty"`

	// result event
	CostUSD    float64 `json:"cost_usd,omitempty"`
	DurationMs int64   `json:"duration_ms,omitempty"`
	NumTurns   int     `json:"num_turns,omitempty"`
}

// IsTerminal returns true if the event signals the end of a request stream.
func (e Event) IsTerminal() bool {
	return e.Type == "result" || e.Type == "error" || e.Type == "pong"
}
