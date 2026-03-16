package runtime

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ProjectRoot returns the per-project Aurelia runtime root under the target project.
func ProjectRoot(cwd string) string {
	return filepath.Join(cwd, defaultDir)
}

// ProjectSkills returns the per-project skills directory under the target project.
func ProjectSkills(cwd string) string {
	return filepath.Join(ProjectRoot(cwd), "skills")
}

// BootstrapProject ensures the target project contains the minimal Aurelia local runtime tree.
func BootstrapProject(cwd string) error {
	if strings.TrimSpace(cwd) == "" {
		return nil
	}

	dirs := []string{
		ProjectRoot(cwd),
		ProjectSkills(cwd),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return fmt.Errorf("runtime: bootstrap project directory %q: %w", dir, err)
		}
	}
	return nil
}


