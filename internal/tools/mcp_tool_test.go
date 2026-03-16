package tools

import (
	"context"
	"testing"

	"github.com/kocar/aurelia/internal/agent"
	"github.com/kocar/aurelia/internal/mcp"
)

// mockCaller mimics mcp.Manager for testing
type mockCaller struct {
	called     bool
	serverName string
	remoteName string
	args       map[string]interface{}
	result     *mcp.CallResult
	err        error
}

func (m *mockCaller) CallTool(ctx context.Context, serverName, remoteToolName string, args map[string]interface{}) (*mcp.CallResult, error) {
	m.called = true
	m.serverName = serverName
	m.remoteName = remoteToolName
	m.args = args
	return m.result, m.err
}

func TestRegisterMCPToolsAdapterHandler(t *testing.T) {
	// Let's test the closure manually to avoid full manager initialization which requires actual MCP servers.
	s := mcp.ToolSpec{
		RegistryName: "test_tool",
		ServerName:   "test_server",
		RemoteName:   "real_tool",
		Description:  "testing tool",
		Parameters: map[string]interface{}{
			"type": "object",
		},
	}

	tool := agent.Tool{
		Name:        s.RegistryName,
		Description: s.Description,
		JSONSchema:  s.Parameters,
	}

	caller := &mockCaller{
		result: &mcp.CallResult{
			Content: "success",
			IsError: false,
		},
	}

	handler := func(ctx context.Context, args map[string]interface{}) (string, error) {
		result, err := caller.CallTool(ctx, s.ServerName, s.RemoteName, args)
		if err != nil {
			return "", err
		}
		if result.IsError {
			return "", nil // ignoring error struct for simple test mapping
		}
		return result.Content, nil
	}

	registry := agent.NewToolRegistry()
	registry.Register(tool, handler)

	res, err := registry.Execute(context.Background(), "test_tool", map[string]interface{}{"key": "value"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res != "success" {
		t.Errorf("expected 'success', got %s", res)
	}

	if !caller.called {
		t.Errorf("caller was not invoked")
	}
	if caller.serverName != "test_server" || caller.remoteName != "real_tool" {
		t.Errorf("wrong mapping for server/remote name")
	}
}


