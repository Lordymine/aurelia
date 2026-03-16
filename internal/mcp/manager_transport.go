package mcp

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/kocar/aurelia/internal/config"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

func connectServer(connectTimeout time.Duration, cfg config.MCPToolsConfig, serverCfg config.MCPServerConfig, workspace string) (*mcpsdk.ClientSession, error) {
	transport, err := buildTransport(cfg, serverCfg, workspace)
	if err != nil {
		return nil, err
	}

	client := mcpsdk.NewClient(&mcpsdk.Implementation{
		Name:    mcpClientName(cfg),
		Version: mcpClientVersion(cfg),
	}, nil)

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	return client.Connect(ctx, transport, nil)
}

func buildTransport(cfg config.MCPToolsConfig, serverCfg config.MCPServerConfig, workspace string) (mcpsdk.Transport, error) {
	switch strings.ToLower(strings.TrimSpace(serverCfg.Transport)) {
	case "stdio":
		return buildStdioTransport(serverCfg, workspace)
	case "streamable_http", "streamable-http", "http":
		return buildHTTPTransport(cfg, serverCfg)
	default:
		return nil, fmt.Errorf("unsupported MCP transport %q", serverCfg.Transport)
	}
}

func buildStdioTransport(serverCfg config.MCPServerConfig, workspace string) (mcpsdk.Transport, error) {
	if strings.TrimSpace(serverCfg.Command) == "" {
		return nil, fmt.Errorf("stdio transport requires command")
	}

	cmd := exec.Command(serverCfg.Command, serverCfg.Args...)
	if serverCfg.WorkingDir != "" {
		dir := serverCfg.WorkingDir
		if !filepath.IsAbs(dir) {
			dir = filepath.Join(workspace, dir)
		}
		cmd.Dir = dir
	}

	cmd.Env = os.Environ()
	for key, value := range serverCfg.Env {
		if strings.Contains(key, "=") {
			return nil, fmt.Errorf("invalid MCP env key %q: must not contain '='", key)
		}
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	return &mcpsdk.CommandTransport{Command: cmd}, nil
}

func buildHTTPTransport(cfg config.MCPToolsConfig, serverCfg config.MCPServerConfig) (mcpsdk.Transport, error) {
	if strings.TrimSpace(serverCfg.Endpoint) == "" {
		return nil, fmt.Errorf("streamable_http transport requires endpoint")
	}

	client := &http.Client{}
	mergedHeaders := mergeHeaders(cfg.Headers, serverCfg.Headers)
	if len(mergedHeaders) > 0 {
		client.Transport = &headerRoundTripper{
			base:    http.DefaultTransport,
			headers: mergedHeaders,
		}
	}

	return &mcpsdk.StreamableClientTransport{
		Endpoint:             serverCfg.Endpoint,
		HTTPClient:           client,
		DisableStandaloneSSE: true,
	}, nil
}

func mcpClientName(cfg config.MCPToolsConfig) string {
	clientName := strings.TrimSpace(cfg.ClientName)
	if clientName == "" {
		return "aurelia"
	}
	return clientName
}

func mcpClientVersion(cfg config.MCPToolsConfig) string {
	clientVersion := strings.TrimSpace(cfg.ClientVersion)
	if clientVersion == "" {
		return "1.0.0"
	}
	return clientVersion
}

func mergeHeaders(globalHeaders, localHeaders map[string]string) map[string]string {
	if len(globalHeaders) == 0 && len(localHeaders) == 0 {
		return nil
	}

	out := make(map[string]string, len(globalHeaders)+len(localHeaders))
	for key, value := range globalHeaders {
		if strings.TrimSpace(value) != "" {
			out[key] = value
		}
	}
	for key, value := range localHeaders {
		if strings.TrimSpace(value) != "" {
			out[key] = value
		}
	}
	return out
}

type headerRoundTripper struct {
	base    http.RoundTripper
	headers map[string]string
}

func (rt *headerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	base := rt.base
	if base == nil {
		base = http.DefaultTransport
	}

	clone := req.Clone(req.Context())
	for key, value := range rt.headers {
		clone.Header.Set(key, value)
	}
	return base.RoundTrip(clone)
}


