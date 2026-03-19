package agent

import (
	"reflect"
	"testing"
)

func TestBuildHiddenCapabilityCatalog_IncludesDynamicMCPServers(t *testing.T) {
	t.Parallel()

	defs := []Tool{
		{Name: "read_file"},
		{Name: "run_command"},
		{Name: "mcp_context7_get_library_docs"},
		{Name: "mcp_playwright_browser_navigate"},
	}

	capabilities := BuildHiddenCapabilityCatalog(defs, []string{"read_file"})
	ids := make([]string, 0, len(capabilities))
	for _, capability := range capabilities {
		ids = append(ids, capability.ID)
	}

	expected := []string{"local_exec", "mcp:context7", "mcp:playwright"}
	if !reflect.DeepEqual(ids, expected) {
		t.Fatalf("expected %v, got %v", expected, ids)
	}
}

func TestResolveCapabilityToolNames_ReturnsDynamicMCPTools(t *testing.T) {
	t.Parallel()

	capabilities := []capabilityDefinition{
		{ID: "mcp:context7", ToolNames: []string{"mcp_context7_get_library_docs", "mcp_context7_resolve_library_id"}},
	}

	got := resolveCapabilityToolNames(capabilities, "mcp:context7")
	expected := []string{"mcp_context7_get_library_docs", "mcp_context7_resolve_library_id"}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}
