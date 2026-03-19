package agent

import (
	"sort"
	"strings"
)

const maxToolDescriptionRunes = 140

func CompactToolsForPrompt(tools []Tool) []Tool {
	if len(tools) == 0 {
		return tools
	}

	compacted := make([]Tool, 0, len(tools))
	for _, tool := range tools {
		compacted = append(compacted, Tool{
			Name:        strings.TrimSpace(tool.Name),
			Description: compactToolDescription(tool.Description),
			JSONSchema:  pruneToolSchema(tool.JSONSchema),
		})
	}
	sort.Slice(compacted, func(i, j int) bool {
		return compacted[i].Name < compacted[j].Name
	})
	return compacted
}

func compactToolDescription(description string) string {
	description = strings.TrimSpace(strings.Join(strings.Fields(description), " "))
	if description == "" {
		return ""
	}
	if idx := strings.IndexAny(description, ".;:\n"); idx > 0 {
		description = strings.TrimSpace(description[:idx+1])
	}
	runes := []rune(description)
	if len(runes) <= maxToolDescriptionRunes {
		return description
	}
	return strings.TrimSpace(string(runes[:maxToolDescriptionRunes])) + "..."
}

func pruneToolSchema(schema map[string]interface{}) map[string]interface{} {
	if len(schema) == 0 {
		return map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
		}
	}

	out := make(map[string]interface{})
	if schemaType, ok := schema["type"]; ok {
		out["type"] = schemaType
	}
	if out["type"] == nil {
		out["type"] = "object"
	}

	if properties, ok := schema["properties"].(map[string]interface{}); ok {
		prunedProps := make(map[string]interface{}, len(properties))
		for key, raw := range properties {
			if propertySchema, ok := raw.(map[string]interface{}); ok {
				prunedProps[key] = pruneToolSchema(propertySchema)
			}
		}
		out["properties"] = prunedProps
	}

	if items, ok := schema["items"].(map[string]interface{}); ok {
		out["items"] = pruneToolSchema(items)
	}
	if required, ok := schema["required"]; ok {
		out["required"] = required
	}
	if enumValues, ok := schema["enum"]; ok {
		out["enum"] = enumValues
	}
	if additional, ok := schema["additionalProperties"]; ok {
		switch typed := additional.(type) {
		case map[string]interface{}:
			out["additionalProperties"] = pruneToolSchema(typed)
		default:
			out["additionalProperties"] = typed
		}
	}

	return out
}
