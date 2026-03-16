package mcp

import (
	"testing"
	"time"

	"github.com/kocar/aurelia/internal/config"
)

func TestNewManager_Disabled(t *testing.T) {
	cfg := config.MCPToolsConfig{
		Enabled: false,
	}

	manager, err := NewManager(cfg, ".")
	if err != nil {
		t.Fatalf("expected nil error for disabled config, got %v", err)
	}
	if manager != nil {
		t.Fatalf("expected nil manager, got %v", manager)
	}
}

func TestNewManager_NoServers(t *testing.T) {
	cfg := config.MCPToolsConfig{
		Enabled: true,
		Servers: map[string]config.MCPServerConfig{},
	}

	manager, err := NewManager(cfg, ".")
	if err != nil {
		t.Fatalf("expected nil error for empty servers config, got %v", err)
	}
	if manager != nil {
		t.Fatalf("expected nil manager, got %v", manager)
	}
}

// removed unused test

func TestManagerClose(t *testing.T) {
	m := &Manager{
		servers: make(map[string]*serverSession),
	}
	
	err := m.Close()
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestTimeoutFromMS(t *testing.T) {
	val := timeoutFromMS(0, 5*time.Second)
	if val != 5*time.Second {
		t.Errorf("expected 5s fallback")
	}

	val = timeoutFromMS(2000, 5*time.Second)
	if val != 2*time.Second {
		t.Errorf("expected 2s")
	}
}


