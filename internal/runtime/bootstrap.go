package runtime

import (
	"fmt"
	"os"
)

// requiredDirs lists all directories Bootstrap must ensure exist.
// Order does not matter — MkdirAll creates parent paths as needed.
var requiredDirs = []func(*PathResolver) string{
	(*PathResolver).Config,
	(*PathResolver).Data,
	(*PathResolver).Memory,
	(*PathResolver).MemoryPersonas,
	(*PathResolver).Skills,
	(*PathResolver).Logs,
	(*PathResolver).Agents,
}

// Bootstrap creates the full instance directory tree with 0700 permissions.
// It is safe to call multiple times — existing directories and files are not modified.
// On Windows, the 0700 permission argument is accepted but has no effect (ACL-based permissions).
func Bootstrap(r *PathResolver) error {
	for _, dir := range requiredDirs {
		if err := os.MkdirAll(dir(r), 0700); err != nil {
			return fmt.Errorf("runtime: bootstrap failed to create %q: %w", dir(r), err)
		}
	}
	return nil
}
