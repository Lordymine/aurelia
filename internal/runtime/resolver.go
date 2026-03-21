package runtime

import (
	"fmt"
	"os"
	"path/filepath"
)

const envKey = "AURELIA_HOME"
const defaultDir = ".aurelia"

// PathResolver resolves and exposes all instance-directory paths.
type PathResolver struct {
	root string
}

// New returns a PathResolver whose root is:
//   - $AURELIA_HOME if set and non-empty
//   - $HOME/.aurelia otherwise
//
// Returns a descriptive error if $AURELIA_HOME is unset and os.UserHomeDir() fails.
func New() (*PathResolver, error) {
	if override := os.Getenv(envKey); override != "" {
		return &PathResolver{root: override}, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("runtime: cannot resolve instance root: os.UserHomeDir failed and %s is not set: %w", envKey, err)
	}
	return &PathResolver{root: filepath.Join(home, defaultDir)}, nil
}

// Root returns the instance root directory.
func (r *PathResolver) Root() string { return r.root }

// Config returns the path to the config/ subdirectory.
func (r *PathResolver) Config() string { return filepath.Join(r.root, "config") }

// AppConfig returns the path to the main app config JSON file.
func (r *PathResolver) AppConfig() string { return filepath.Join(r.Config(), "app.json") }

// Data returns the path to the data/ subdirectory.
func (r *PathResolver) Data() string { return filepath.Join(r.root, "data") }

// Memory returns the path to the memory/ subdirectory.
func (r *PathResolver) Memory() string { return filepath.Join(r.root, "memory") }

// MemoryPersonas returns the path to the memory/personas/ subdirectory.
func (r *PathResolver) MemoryPersonas() string { return filepath.Join(r.root, "memory", "personas") }

// Skills returns the path to the skills/ subdirectory.
func (r *PathResolver) Skills() string { return filepath.Join(r.root, "skills") }

// Logs returns the path to the logs/ subdirectory.
func (r *PathResolver) Logs() string { return filepath.Join(r.root, "logs") }

// Agents returns the path to the agents/ subdirectory.
func (r *PathResolver) Agents() string { return filepath.Join(r.root, "agents") }

// DBPath returns the path to a named database file inside the data/ subdirectory.
func (r *PathResolver) DBPath(name string) string { return filepath.Join(r.Data(), name) }
