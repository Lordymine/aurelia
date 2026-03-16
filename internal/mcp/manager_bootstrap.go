package mcp

import (
	"fmt"
	"log"
	"time"

	"github.com/kocar/aurelia/internal/config"
)

type serverResult struct {
	name    string
	session *serverSession
	specs   []ToolSpec
	err     error
}

func NewManager(cfg config.MCPToolsConfig, workspace string) (*Manager, error) {
	if !cfg.Enabled || len(cfg.Servers) == 0 {
		return nil, nil
	}

	connectTimeout := timeoutFromMS(cfg.ConnectTimeoutMS, 10*time.Second)
	defaultCallTimeout := timeoutFromMS(cfg.CallTimeoutMS, 60*time.Second)
	enabledCfgs := enabledMCPServers(cfg.Servers)
	if len(enabledCfgs) == 0 {
		return nil, nil
	}

	manager := &Manager{servers: make(map[string]*serverSession)}
	manager.connectEnabledServers(cfg, enabledCfgs, workspace, connectTimeout, defaultCallTimeout)
	if len(manager.servers) == 0 {
		_ = manager.Close()
		return nil, fmt.Errorf("no MCP servers connected successfully")
	}

	return manager, nil
}

func (m *Manager) connectEnabledServers(
	cfg config.MCPToolsConfig,
	servers []config.MCPServerConfig,
	workspace string,
	connectTimeout time.Duration,
	defaultCallTimeout time.Duration,
) {
	for result := range connectServers(cfg, servers, workspace, connectTimeout, defaultCallTimeout) {
		if result.err != nil {
			log.Printf("[MCP] Failed to connect server %q: %v\n", result.name, result.err)
			continue
		}

		m.servers[result.name] = result.session
		m.tools = append(m.tools, result.specs...)
		log.Printf("[MCP] Connected server %q with %d tools\n", result.name, len(result.specs))
	}
}

func enabledMCPServers(servers map[string]config.MCPServerConfig) []config.MCPServerConfig {
	enabled := make([]config.MCPServerConfig, 0, len(servers))
	for _, serverCfg := range servers {
		if serverCfg.Enabled {
			enabled = append(enabled, serverCfg)
		}
	}
	return enabled
}


