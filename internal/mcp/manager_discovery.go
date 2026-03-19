package mcp

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/kocar/aurelia/internal/config"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

func discoverTools(session *mcpsdk.ClientSession, serverCfg config.MCPServerConfig, usedNames map[string]int) ([]ToolSpec, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	remoteTools, err := listAllTools(ctx, session)
	if err != nil {
		return nil, err
	}

	specs := make([]ToolSpec, 0, len(remoteTools))
	for _, remote := range remoteTools {
		if remote == nil || remote.Name == "" {
			continue
		}
		if !isToolAllowed(remote.Name, serverCfg.AllowTools) {
			continue
		}

		baseName := "mcp_" + sanitizeToolSegment(serverCfg.Name) + "_" + sanitizeToolSegment(remote.Name)
		registryName := uniqueToolName(baseName, usedNames)

		specs = append(specs, ToolSpec{
			RegistryName: registryName,
			ServerName:   serverCfg.Name,
			RemoteName:   remote.Name,
			Description:  buildToolDescription(serverCfg.Name, remote.Name, remote.Description),
			Parameters:   normalizeSchema(remote.InputSchema),
		})
	}

	return specs, nil
}

func listAllTools(ctx context.Context, session *mcpsdk.ClientSession) ([]*mcpsdk.Tool, error) {
	var all []*mcpsdk.Tool
	cursor := ""

	for {
		params := &mcpsdk.ListToolsParams{}
		if cursor != "" {
			params.Cursor = cursor
		}

		result, err := session.ListTools(ctx, params)
		if err != nil {
			return nil, err
		}
		if result == nil {
			break
		}

		all = append(all, result.Tools...)
		if result.NextCursor == "" || result.NextCursor == cursor {
			break
		}
		cursor = result.NextCursor
	}

	return all, nil
}

func buildToolDescription(serverName, remoteName, rawDescription string) string {
	description := strings.TrimSpace(rawDescription)
	if description == "" {
		return fmt.Sprintf("MCP tool %s.", remoteName)
	}
	description = strings.Join(strings.Fields(description), " ")
	if idx := strings.IndexAny(description, ".;:\n"); idx > 0 {
		description = strings.TrimSpace(description[:idx+1])
	}
	runes := []rune(description)
	if len(runes) > 140 {
		description = strings.TrimSpace(string(runes[:140])) + "..."
	}
	return description
}

func timeoutFromMS(ms int, fallback time.Duration) time.Duration {
	if ms <= 0 {
		return fallback
	}
	return time.Duration(ms) * time.Millisecond
}

func isToolAllowed(name string, allow []string) bool {
	if len(allow) == 0 {
		return true
	}
	for _, item := range allow {
		if strings.TrimSpace(item) == name {
			return true
		}
	}
	return false
}

func uniqueToolName(base string, used map[string]int) string {
	if used[base] == 0 {
		used[base] = 1
		return base
	}
	used[base]++
	return fmt.Sprintf("%s_%d", base, used[base])
}

func sanitizeToolSegment(raw string) string {
	normalized := strings.TrimSpace(raw)
	normalized = strings.ToLower(normalized)
	normalized = invalidToolNameChars.ReplaceAllString(normalized, "_")
	normalized = strings.Trim(normalized, "_")
	if normalized == "" {
		return "tool"
	}
	return normalized
}

func normalizeSchema(inputSchema any) map[string]interface{} {
	if inputSchema == nil {
		return map[string]interface{}{
			"type":                 "object",
			"additionalProperties": true,
		}
	}

	m, ok := inputSchema.(map[string]interface{})
	if !ok {
		return map[string]interface{}{
			"type":                 "object",
			"additionalProperties": true,
		}
	}

	if _, ok := m["type"]; !ok {
		m["type"] = "object"
	}
	if m["type"] == "object" {
		if _, hasProps := m["properties"]; !hasProps {
			if _, hasAdditional := m["additionalProperties"]; !hasAdditional {
				m["additionalProperties"] = true
			}
		}
	}

	return m
}
