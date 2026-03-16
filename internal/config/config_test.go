package config

import (
	"path/filepath"
	"testing"

	"github.com/kocar/aurelia/internal/runtime"
)

func TestLoad_DBPath_DefaultsToInstance(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("AURELIA_HOME", tmpDir)
	t.Setenv("DB_PATH", "")

	r, err := runtime.New()
	if err != nil {
		t.Fatalf("runtime.New() unexpected error: %v", err)
	}

	cfg, err := Load(r)
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	want := filepath.Join(tmpDir, "data", "aurelia.db")
	if cfg.DBPath != want {
		t.Errorf("DBPath = %q, want %q", cfg.DBPath, want)
	}
}

func TestLoad_DBPath_EnvVarOverrides(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("AURELIA_HOME", tmpDir)
	t.Setenv("DB_PATH", "/custom/path.db")

	r, err := runtime.New()
	if err != nil {
		t.Fatalf("runtime.New() unexpected error: %v", err)
	}

	cfg, err := Load(r)
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	want := "/custom/path.db"
	if cfg.DBPath != want {
		t.Errorf("DBPath = %q, want %q", cfg.DBPath, want)
	}
}

func TestLoad_MCPPath_EnvVarOverrides(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("AURELIA_HOME", tmpDir)
	t.Setenv("MCP_SERVERS_CONFIG_PATH", "/custom/mcp.json")

	r, err := runtime.New()
	if err != nil {
		t.Fatalf("runtime.New() unexpected error: %v", err)
	}

	cfg, err := Load(r)
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	want := "/custom/mcp.json"
	if cfg.MCPConfigPath != want {
		t.Errorf("MCPConfigPath = %q, want %q", cfg.MCPConfigPath, want)
	}
}

func TestLoad_MCPPath_DefaultsToInstance(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("AURELIA_HOME", tmpDir)
	t.Setenv("MCP_SERVERS_CONFIG_PATH", "")

	r, err := runtime.New()
	if err != nil {
		t.Fatalf("runtime.New() unexpected error: %v", err)
	}

	cfg, err := Load(r)
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	want := filepath.Join(tmpDir, "config", "mcp_servers.json")
	if cfg.MCPConfigPath != want {
		t.Errorf("MCPConfigPath = %q, want %q", cfg.MCPConfigPath, want)
	}
}


