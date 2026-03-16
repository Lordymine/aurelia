package config

import (
	"encoding/json"
	"os"
)

type MCPServerConfig struct {
	Name       string            `json:"name"`
	Enabled    bool              `json:"enabled"`
	Command    string            `json:"command"`
	Args       []string          `json:"args"`
	Env        map[string]string `json:"env"`
	WorkingDir string            `json:"workingDir"`
	Transport  string            `json:"transport"` // "stdio", "http"
	Endpoint   string            `json:"endpoint"`
	Headers    map[string]string `json:"headers"`
	TimeoutMS  int               `json:"timeoutMs"`
	AllowTools []string          `json:"allowTools"`
}

type MCPToolsConfig struct {
	Enabled          bool                       `json:"enabled"`
	ClientName       string                     `json:"clientName"`
	ClientVersion    string                     `json:"clientVersion"`
	ConnectTimeoutMS int                        `json:"connectTimeoutMs"`
	CallTimeoutMS    int                        `json:"callTimeoutMs"`
	Headers          map[string]string          `json:"headers"`
	Servers          map[string]MCPServerConfig `json:"mcpServers"`
}

// LoadMCPConfig reads an MCP servers JSON configuration file (Claude desktop style).
func LoadMCPConfig(path string) (*MCPToolsConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Provide an empty but valid disabled config if file missing
			return &MCPToolsConfig{Enabled: false}, nil
		}
		return nil, err
	}

	var cfg MCPToolsConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// By default, if the file exists and is parsed, we could assume it's enabled 
	// unless explicitly structured differently. For now, let's assume if there are servers, it's enabled.
	if len(cfg.Servers) > 0 {
		cfg.Enabled = true
	}

	// Assign names to servers from map keys if not set
	for name, srv := range cfg.Servers {
		if srv.Name == "" {
			srv.Name = name
		}
		// Default transport to stdio if command exists
		if srv.Transport == "" && srv.Command != "" {
			srv.Transport = "stdio"
		}
		// Let's assume servers listed are enabled by default
		srv.Enabled = true
		cfg.Servers[name] = srv
	}

	return &cfg, nil
}
