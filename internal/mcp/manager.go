package mcp

import (
	"context"
	"fmt"
	"regexp"
	"sync"
	"time"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

var invalidToolNameChars = regexp.MustCompile(`[^a-zA-Z0-9_]`)

type ToolSpec struct {
	RegistryName string
	ServerName   string
	RemoteName   string
	Description  string
	Parameters   map[string]interface{}
}

type CallResult struct {
	Content string
	IsError bool
}

type Caller interface {
	CallTool(ctx context.Context, serverName, remoteToolName string, args map[string]interface{}) (*CallResult, error)
}

type serverSession struct {
	name        string
	session     *mcpsdk.ClientSession
	callTimeout time.Duration
}

type Manager struct {
	mu      sync.RWMutex
	servers map[string]*serverSession
	tools   []ToolSpec
}

func (m *Manager) ToolSpecs() []ToolSpec {
	if m == nil {
		return nil
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]ToolSpec, len(m.tools))
	copy(out, m.tools)
	return out
}

func (m *Manager) CallTool(ctx context.Context, serverName, remoteToolName string, args map[string]interface{}) (*CallResult, error) {
	server, err := m.lookupServer(serverName)
	if err != nil {
		return nil, err
	}

	if args == nil {
		args = map[string]interface{}{}
	}

	callCtx, cancel := newCallContext(ctx, server.callTimeout)
	defer cancel()

	result, err := server.session.CallTool(callCtx, &mcpsdk.CallToolParams{
		Name:      remoteToolName,
		Arguments: args,
	})
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, fmt.Errorf("MCP server returned nil tool result")
	}

	return &CallResult{
		Content: formatCallToolResult(result),
		IsError: result.IsError,
	}, nil
}

func (m *Manager) lookupServer(serverName string) (*serverSession, error) {
	if m == nil {
		return nil, fmt.Errorf("MCP manager is not initialized")
	}

	m.mu.RLock()
	server, ok := m.servers[serverName]
	m.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("MCP server %q is not connected", serverName)
	}
	return server, nil
}

func newCallContext(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout <= 0 {
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, timeout)
}
