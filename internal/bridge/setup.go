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
	bundlePath := filepath.Join(targetDir, "bundle.ts")
	nodeModules := filepath.Join(targetDir, "node_modules")

	// Always update bundle.ts to match embedded version
	needsNpmInstall := false
	if _, err := os.Stat(nodeModules); os.IsNotExist(err) {
		needsNpmInstall = true
	}

	existing, readErr := os.ReadFile(bundlePath)
	bundleUpToDate := readErr == nil && string(existing) == string(bundleJS)

	if bundleUpToDate && !needsNpmInstall {
		return targetDir, nil
	}

	if needsNpmInstall {
		log.Println("Setting up Bridge for first time...")
	} else if !bundleUpToDate {
		log.Println("Updating Bridge bundle...")
	}

	if err := os.MkdirAll(targetDir, 0700); err != nil {
		return "", fmt.Errorf("create bridge dir: %w", err)
	}

	// Always write latest bundle.ts
	if !bundleUpToDate {
		if err := os.WriteFile(bundlePath, bundleJS, 0600); err != nil {
			return "", fmt.Errorf("write bundle.ts: %w", err)
		}
	}

	// npm install only if node_modules missing
	if needsNpmInstall {
		pkgPath := filepath.Join(targetDir, "package.json")
		if err := os.WriteFile(pkgPath, []byte(bridgePackageJSON), 0600); err != nil {
			return "", fmt.Errorf("write package.json: %w", err)
		}

		log.Println("Installing Claude Agent SDK (npm install)...")
		cmd := exec.Command("npm", "install", "--production", "--no-optional")
		cmd.Dir = targetDir
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("npm install failed: %w", err)
		}
	}

	log.Println("Bridge setup complete.")
	return targetDir, nil
}
