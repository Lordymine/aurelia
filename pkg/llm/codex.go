package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kocar/aurelia/internal/agent"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

var codexLookPath = exec.LookPath
var codexCallerFactory = connectCodexMCP

func CodexLookPathForTest(fn func(string) (string, error)) func() {
	previous := codexLookPath
	codexLookPath = fn
	return func() {
		codexLookPath = previous
	}
}

func CodexCallerFactoryForTest(fn func() (codexToolCaller, error)) func() {
	previous := codexCallerFactory
	codexCallerFactory = fn
	return func() {
		codexCallerFactory = previous
	}
}

func UseNoopCodexCallerForTest() func() {
	return CodexCallerFactoryForTest(func() (codexToolCaller, error) {
		return &noopCodexCaller{}, nil
	})
}

type codexToolResponse struct {
	Content  string
	ThreadID string
	IsError  bool
}

type codexToolCaller interface {
	CallTool(ctx context.Context, name string, args map[string]any) (*codexToolResponse, error)
	Close() error
}

type noopCodexCaller struct{}

func (noopCodexCaller) CallTool(_ context.Context, _ string, _ map[string]any) (*codexToolResponse, error) {
	return &codexToolResponse{Content: `{"content":"","reasoning_content":"","tool_calls":[]}`}, nil
}

func (noopCodexCaller) Close() error { return nil }

type codexMCPCaller struct {
	session *mcpsdk.ClientSession
}

func connectCodexMCP() (codexToolCaller, error) {
	cmd := exec.Command("codex", "mcp-server")
	transport := &mcpsdk.CommandTransport{Command: cmd}
	client := mcpsdk.NewClient(&mcpsdk.Implementation{
		Name:    "aurelia-codex",
		Version: "1.0.0",
	}, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		return nil, fmt.Errorf("connect codex mcp-server: %w", err)
	}
	return &codexMCPCaller{session: session}, nil
}

func (c *codexMCPCaller) Close() error {
	if c == nil || c.session == nil {
		return nil
	}
	return c.session.Close()
}

func (c *codexMCPCaller) CallTool(ctx context.Context, name string, args map[string]any) (*codexToolResponse, error) {
	result, err := c.session.CallTool(ctx, &mcpsdk.CallToolParams{
		Name:      name,
		Arguments: args,
	})
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, fmt.Errorf("codex mcp returned nil tool result")
	}

	resp := &codexToolResponse{
		Content:  extractCodexContent(result),
		ThreadID: extractCodexThreadID(result),
		IsError:  result.IsError,
	}
	return resp, nil
}

// CodexCLIProvider uses the local Codex MCP server as an experimental provider backend.
type CodexCLIProvider struct {
	model   string
	caller  codexToolCaller
	mu      sync.Mutex
	threads map[string]string
}

func NewCodexCLIProvider(model string) (*CodexCLIProvider, error) {
	caller, err := codexCallerFactory()
	if err != nil {
		return nil, err
	}
	return &CodexCLIProvider{
		model:   model,
		caller:  caller,
		threads: make(map[string]string),
	}, nil
}

func (p *CodexCLIProvider) Close() {
	if p == nil || p.caller == nil {
		return
	}
	_ = p.caller.Close()
}

func EnsureCodexCLIAvailable() error {
	if _, err := codexLookPath("codex"); err != nil {
		return fmt.Errorf("codex CLI not found in PATH")
	}
	return nil
}

func (p *CodexCLIProvider) GenerateContent(
	ctx context.Context,
	systemPrompt string,
	history []agent.Message,
	tools []agent.Tool,
) (*agent.ModelResponse, error) {
	if historyHasImages(history) {
		return nil, VisionUnsupportedError{provider: "openai_codex", model: p.model}
	}

	prompt, err := buildCodexPrompt(systemPrompt, history, tools)
	if err != nil {
		return nil, err
	}

	runID, _ := agent.RunContextFromContext(ctx)
	threadID := p.threadID(runID)

	toolName := "codex"
	args := map[string]any{
		"prompt":                 prompt,
		"approval-policy":        "never",
		"sandbox":                "read-only",
		"model":                  p.model,
		"developer-instructions": "Return only the requested response payload; do not run shell commands or modify files yourself.",
	}
	if workdir, ok := agent.WorkdirFromContext(ctx); ok {
		args["cwd"] = workdir
	}
	if threadID != "" {
		toolName = "codex-reply"
		args = map[string]any{
			"threadId": threadID,
			"prompt":   prompt,
		}
	}

	result, err := p.caller.CallTool(ctx, toolName, args)
	if err != nil {
		return nil, fmt.Errorf("codex mcp call failed: %w", err)
	}
	if result.IsError {
		return nil, fmt.Errorf("codex mcp tool returned an error: %s", result.Content)
	}
	if runID != "" && result.ThreadID != "" {
		p.setThreadID(runID, result.ThreadID)
	}

	return parseCodexResponse([]byte(result.Content))
}

func (p *CodexCLIProvider) threadID(runID string) string {
	if runID == "" {
		return ""
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.threads[runID]
}

func (p *CodexCLIProvider) setThreadID(runID, threadID string) {
	if runID == "" || threadID == "" {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.threads[runID] = threadID
}

func buildCodexPrompt(systemPrompt string, history []agent.Message, tools []agent.Tool) (string, error) {
	historyJSON, err := json.Marshal(history)
	if err != nil {
		return "", fmt.Errorf("marshal history: %w", err)
	}

	toolSpecJSON, err := json.Marshal(tools)
	if err != nil {
		return "", fmt.Errorf("marshal tools: %w", err)
	}

	return strings.TrimSpace(fmt.Sprintf(`
You are acting as an OpenAI Codex backend for Aurelia.
Return ONLY valid JSON. Do not wrap it in markdown fences.

JSON schema:
{
  "content": "string",
  "reasoning_content": "string",
  "tool_calls": [
    {
      "id": "string",
      "name": "tool-name",
      "arguments": {}
    }
  ]
}

Rules:
- If no tool is needed, return "tool_calls": [].
- If a tool is needed, keep "content" brief and fill "tool_calls".
- Arguments must be a JSON object.
- Use only the provided tools.

SYSTEM PROMPT:
%s

HISTORY JSON:
%s

TOOLS JSON:
%s
`, systemPrompt, string(historyJSON), string(toolSpecJSON))), nil
}

func extractCodexContent(result *mcpsdk.CallToolResult) string {
	if result == nil {
		return ""
	}
	if content, ok := extractStructuredString(result.StructuredContent, "content"); ok && content != "" {
		return content
	}
	var parts []string
	for _, item := range result.Content {
		if text, ok := item.(*mcpsdk.TextContent); ok && strings.TrimSpace(text.Text) != "" {
			parts = append(parts, text.Text)
		}
	}
	return strings.Join(parts, "\n")
}

func extractCodexThreadID(result *mcpsdk.CallToolResult) string {
	if result == nil {
		return ""
	}
	threadID, _ := extractStructuredString(result.StructuredContent, "threadId")
	return threadID
}

func extractStructuredString(value any, key string) (string, bool) {
	m, ok := value.(map[string]any)
	if !ok {
		return "", false
	}
	raw, ok := m[key]
	if !ok {
		return "", false
	}
	text, ok := raw.(string)
	if !ok {
		return "", false
	}
	return text, true
}

func parseCodexResponse(raw []byte) (*agent.ModelResponse, error) {
	clean := strings.TrimSpace(string(raw))
	clean = strings.TrimPrefix(clean, "```json")
	clean = strings.TrimPrefix(clean, "```")
	clean = strings.TrimSuffix(clean, "```")
	clean = strings.TrimSpace(clean)

	var payload struct {
		Content          string `json:"content"`
		ReasoningContent string `json:"reasoning_content"`
		ToolCalls        []struct {
			ID        string         `json:"id"`
			Name      string         `json:"name"`
			Arguments map[string]any `json:"arguments"`
		} `json:"tool_calls"`
	}
	if err := json.Unmarshal([]byte(clean), &payload); err != nil {
		return nil, fmt.Errorf("decode codex response: %w", err)
	}

	response := &agent.ModelResponse{
		Content:          payload.Content,
		ReasoningContent: payload.ReasoningContent,
	}
	for _, call := range payload.ToolCalls {
		id := strings.TrimSpace(call.ID)
		if id == "" {
			id = uuid.NewString()
		}
		response.ToolCalls = append(response.ToolCalls, agent.ToolCall{
			ID:        id,
			Name:      call.Name,
			Arguments: call.Arguments,
		})
	}
	return response, nil
}
