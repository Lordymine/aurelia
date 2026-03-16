package mcp

import (
	"encoding/json"
	"fmt"
	"strings"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

func formatCallToolResult(result *mcpsdk.CallToolResult) string {
	parts := make([]string, 0, len(result.Content)+1)
	for _, block := range result.Content {
		switch content := block.(type) {
		case *mcpsdk.TextContent:
			if strings.TrimSpace(content.Text) != "" {
				parts = append(parts, content.Text)
			}
		case *mcpsdk.ImageContent:
			parts = append(parts, fmt.Sprintf("[image content mime=%s bytes=%d]", content.MIMEType, len(content.Data)))
		case *mcpsdk.EmbeddedResource:
			formatted := formatEmbeddedResource(content)
			if formatted != "" {
				parts = append(parts, formatted)
			}
		default:
			if raw, err := json.Marshal(block); err == nil && len(raw) > 0 {
				parts = append(parts, string(raw))
			}
		}
	}

	if len(parts) == 0 {
		if result.IsError {
			return "MCP tool failed with empty response"
		}
		return "MCP tool completed with empty response"
	}

	return strings.Join(parts, "\n")
}

func formatEmbeddedResource(content *mcpsdk.EmbeddedResource) string {
	if content == nil || content.Resource == nil {
		return ""
	}

	switch {
	case strings.TrimSpace(content.Resource.Text) != "":
		return content.Resource.Text
	case len(content.Resource.Blob) > 0:
		return fmt.Sprintf("[embedded resource %s blob bytes=%d]", content.Resource.URI, len(content.Resource.Blob))
	default:
		return fmt.Sprintf("[embedded resource %s]", content.Resource.URI)
	}
}
