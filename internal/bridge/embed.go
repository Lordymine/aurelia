package bridge

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed bundle.js
var bundleJS []byte

// ExtractBundle writes the embedded bundle.js to the target directory
// if it doesn't already exist or has changed. Returns the path to the
// extracted bundle.js file.
func ExtractBundle(targetDir string) (string, error) {
	if err := os.MkdirAll(targetDir, 0700); err != nil {
		return "", fmt.Errorf("create bridge dir: %w", err)
	}

	bundlePath := filepath.Join(targetDir, "bundle.js")

	// Check if existing file matches
	existing, err := os.ReadFile(bundlePath)
	if err == nil && len(existing) == len(bundleJS) {
		return bundlePath, nil // already up to date
	}

	if err := os.WriteFile(bundlePath, bundleJS, 0600); err != nil {
		return "", fmt.Errorf("write bundle.js: %w", err)
	}

	return bundlePath, nil
}
