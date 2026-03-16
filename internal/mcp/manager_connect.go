package mcp

import (
	"fmt"
	"sync"
	"time"

	"github.com/kocar/aurelia/internal/config"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

func connectServers(
	cfg config.MCPToolsConfig,
	servers []config.MCPServerConfig,
	workspace string,
	connectTimeout time.Duration,
	defaultCallTimeout time.Duration,
) <-chan serverResult {
	results := make(chan serverResult, len(servers))
	var wg sync.WaitGroup

	for _, serverCfg := range servers {
		wg.Add(1)
		go func(sc config.MCPServerConfig) {
			defer wg.Done()
			results <- connectServerResult(cfg, sc, workspace, connectTimeout, defaultCallTimeout)
		}(serverCfg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

func connectServerResult(
	cfg config.MCPToolsConfig,
	serverCfg config.MCPServerConfig,
	workspace string,
	connectTimeout time.Duration,
	defaultCallTimeout time.Duration,
) serverResult {
	session, err := connectServer(connectTimeout, cfg, serverCfg, workspace)
	if err != nil {
		return serverResult{name: serverCfg.Name, err: fmt.Errorf("connect: %w", err)}
	}

	specs, err := discoverTools(session, serverCfg, map[string]int{})
	if err != nil {
		_ = session.Close()
		return serverResult{name: serverCfg.Name, err: fmt.Errorf("discover: %w", err)}
	}

	return serverResult{
		name:    serverCfg.Name,
		session: newServerSession(serverCfg, session, defaultCallTimeout),
		specs:   specs,
	}
}

func newServerSession(
	serverCfg config.MCPServerConfig,
	session *mcpsdk.ClientSession,
	defaultCallTimeout time.Duration,
) *serverSession {
	callTimeout := defaultCallTimeout
	if serverCfg.TimeoutMS > 0 {
		callTimeout = time.Duration(serverCfg.TimeoutMS) * time.Millisecond
	}

	return &serverSession{
		name:        serverCfg.Name,
		session:     session,
		callTimeout: callTimeout,
	}
}


