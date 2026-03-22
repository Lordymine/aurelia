package bridge

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const bridgePackageJSON = `{
  "name": "aurelia-bridge",
  "version": "1.0.0",
  "private": true,
  "type": "module",
  "dependencies": {
    "@anthropic-ai/claude-agent-sdk": "latest"
  }
}
`

// EnsureBridge checks if the bridge is set up at targetDir. If not,
// creates it with package.json and runs npm install. Returns the
// directory path. bundleJS should be the compiled bridge source.
func EnsureBridge(targetDir string, bundleJS []byte) (string, error) {
	bundlePath := filepath.Join(targetDir, "bundle.js")
	nodeModules := filepath.Join(targetDir, "node_modules")

	// Already set up?
	if _, err := os.Stat(bundlePath); err == nil {
		if _, err := os.Stat(nodeModules); err == nil {
			return targetDir, nil
		}
	}

	log.Println("Setting up Bridge for first time...")

	if err := os.MkdirAll(targetDir, 0700); err != nil {
		return "", fmt.Errorf("create bridge dir: %w", err)
	}

	// Write package.json
	pkgPath := filepath.Join(targetDir, "package.json")
	if err := os.WriteFile(pkgPath, []byte(bridgePackageJSON), 0600); err != nil {
		return "", fmt.Errorf("write package.json: %w", err)
	}

	// Write bundle.js
	if err := os.WriteFile(bundlePath, bundleJS, 0600); err != nil {
		return "", fmt.Errorf("write bundle.js: %w", err)
	}

	// npm install (production only — just the SDK)
	log.Println("Installing Claude Agent SDK (npm install)...")
	cmd := exec.Command("npm", "install", "--production", "--no-optional")
	cmd.Dir = targetDir
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("npm install failed: %w", err)
	}

	log.Println("Bridge setup complete.")
	return targetDir, nil
}
