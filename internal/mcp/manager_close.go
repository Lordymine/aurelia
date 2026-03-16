package mcp

import (
	"fmt"
	"strings"
	"sync"
	"time"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type closeResult struct {
	name string
	err  error
}

func (m *Manager) Close() error {
	if m == nil {
		return nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	var errs []string
	for result := range closeSessions(m.servers) {
		if result.err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", result.name, result.err))
		}
	}

	m.servers = map[string]*serverSession{}
	m.tools = nil

	if len(errs) > 0 {
		return fmt.Errorf("failed to close MCP sessions: %s", strings.Join(errs, "; "))
	}
	return nil
}

func closeSessions(servers map[string]*serverSession) <-chan closeResult {
	results := make(chan closeResult, len(servers))
	var wg sync.WaitGroup

	for name, server := range servers {
		if server == nil || server.session == nil {
			continue
		}
		wg.Add(1)
		go func(n string, sess *mcpsdk.ClientSession) {
			defer wg.Done()
			done := make(chan error, 1)
			go func() { done <- sess.Close() }()
			select {
			case err := <-done:
				results <- closeResult{name: n, err: err}
			case <-time.After(5 * time.Second):
				results <- closeResult{name: n, err: fmt.Errorf("close timed out after 5s")}
			}
		}(name, server.session)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}
