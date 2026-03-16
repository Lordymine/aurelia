package tools

import (
	"context"
	"fmt"

	"github.com/kocar/aurelia/internal/agent"
	"github.com/kocar/aurelia/internal/mcp"
)

// RegisterMCPTools registers all tools provided by the MCP Manager into the given ToolRegistry
func RegisterMCPTools(registry *agent.ToolRegistry, manager *mcp.Manager) {
	if manager == nil || registry == nil {
		return
	}

	for _, spec := range manager.ToolSpecs() {
		tool := agent.Tool{
			Name:        spec.RegistryName,
			Description: spec.Description,
			JSONSchema:  spec.Parameters,
		}

		// capture spec by value for closure
		s := spec 
		
		handler := func(ctx context.Context, args map[string]interface{}) (string, error) {
			result, err := manager.CallTool(ctx, s.ServerName, s.RemoteName, args)
			if err != nil {
				return "", fmt.Errorf("MCP tool %s/%s failed: %w", s.ServerName, s.RemoteName, err)
			}
			if result.IsError {
				return "", fmt.Errorf("MCP tool returned error content: %s", result.Content)
			}
			return result.Content, nil
		}

		registry.Register(tool, handler)
	}
}


