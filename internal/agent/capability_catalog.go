package agent

import (
	"fmt"
	"sort"
	"strings"
)

const requestToolAccessToolName = "request_tool_access"

type capabilityDefinition struct {
	ID          string
	Description string
	ToolNames   []string
}

var coreCapabilityDefinitions = []capabilityDefinition{
	{
		ID:          "filesystem",
		Description: "Read, write, and list local project files.",
		ToolNames:   []string{"read_file", "write_file", "list_dir"},
	},
	{
		ID:          "local_exec",
		Description: "Run local shell commands and tests.",
		ToolNames:   []string{"run_command"},
	},
	{
		ID:          "web",
		Description: "Search the web for current information.",
		ToolNames:   []string{"web_search"},
	},
	{
		ID:          "scheduler",
		Description: "Create and manage reminders and recurring jobs.",
		ToolNames:   []string{"create_schedule", "list_schedules", "pause_schedule", "resume_schedule", "delete_schedule"},
	},
	{
		ID:          "team",
		Description: "Delegate work and inspect team execution state.",
		ToolNames:   []string{"spawn_agent", "pause_team", "resume_team", "cancel_team", "team_status", "send_team_message", "read_team_inbox"},
	},
	{
		ID:          "skills",
		Description: "Install runtime skills.",
		ToolNames:   []string{"install_skill"},
	},
}

func BuildHiddenCapabilityCatalog(defs []Tool, allowed []string) []capabilityDefinition {
	if allowed == nil {
		return nil
	}

	defMap := make(map[string]Tool, len(defs))
	for _, def := range defs {
		defMap[def.Name] = def
	}

	allowedSet := make(map[string]bool, len(allowed))
	for _, name := range allowed {
		allowedSet[name] = true
	}

	var capabilities []capabilityDefinition
	for _, capability := range coreCapabilityDefinitions {
		hidden := collectHiddenToolNames(capability.ToolNames, defMap, allowedSet)
		if len(hidden) == 0 {
			continue
		}
		capabilities = append(capabilities, capabilityDefinition{
			ID:          capability.ID,
			Description: capability.Description,
			ToolNames:   hidden,
		})
	}

	for _, capability := range buildDynamicMCPCapabilities(defs, allowedSet) {
		capabilities = append(capabilities, capability)
	}

	sort.Slice(capabilities, func(i, j int) bool {
		return capabilities[i].ID < capabilities[j].ID
	})
	return capabilities
}

func buildDynamicMCPCapabilities(defs []Tool, allowedSet map[string]bool) []capabilityDefinition {
	grouped := make(map[string][]string)
	for _, def := range defs {
		serverName, _ := SplitMCPToolName(def.Name)
		if serverName == "" || allowedSet[def.Name] {
			continue
		}
		grouped[serverName] = append(grouped[serverName], def.Name)
	}

	var capabilities []capabilityDefinition
	for serverName, toolNames := range grouped {
		sort.Strings(toolNames)
		capabilities = append(capabilities, capabilityDefinition{
			ID:          "mcp:" + serverName,
			Description: fmt.Sprintf("Use tools from MCP server %s.", serverName),
			ToolNames:   toolNames,
		})
	}
	return capabilities
}

func collectHiddenToolNames(names []string, defs map[string]Tool, allowedSet map[string]bool) []string {
	out := make([]string, 0, len(names))
	for _, name := range names {
		if _, ok := defs[name]; !ok || allowedSet[name] {
			continue
		}
		out = append(out, name)
	}
	return out
}

func capabilityRequestTool(capabilities []capabilityDefinition) *Tool {
	if len(capabilities) == 0 {
		return nil
	}

	return &Tool{
		Name:        requestToolAccessToolName,
		Description: "Requests access to one hidden capability from the runtime catalog before using its tools.",
		JSONSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"capability": map[string]interface{}{
					"type":        "string",
					"description": "Exact capability id from the runtime catalog, such as filesystem, local_exec, web, scheduler, team, skills, or mcp:<server>.",
				},
				"reason": map[string]interface{}{
					"type":        "string",
					"description": "Short reason why the capability is needed for the current user request.",
				},
			},
			"required": []string{"capability"},
		},
	}
}

func resolveCapabilityToolNames(capabilities []capabilityDefinition, capabilityID string) []string {
	capabilityID = strings.ToLower(strings.TrimSpace(capabilityID))
	for _, capability := range capabilities {
		if capability.ID != capabilityID {
			continue
		}
		return append([]string(nil), capability.ToolNames...)
	}
	return nil
}

func buildCapabilityCatalogPrompt(capabilities []capabilityDefinition) string {
	if len(capabilities) == 0 {
		return ""
	}

	lines := []string{
		"# ON-DEMAND CAPABILITIES",
		"Hidden capabilities can be expanded on demand with `request_tool_access` when the current tool list is insufficient.",
		"Available capability ids:",
	}
	for _, capability := range capabilities {
		lines = append(lines, fmt.Sprintf("- %s: %s", capability.ID, capability.Description))
	}
	lines = append(lines, "If the current tool list is missing something you need, call `request_tool_access` with the exact capability id before claiming the runtime cannot do it.")
	return strings.Join(lines, "\n")
}
