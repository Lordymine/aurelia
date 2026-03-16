package skill

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Metadata represents the YAML frontmatter of a SKILL.md
type Metadata struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// Skill represents a loaded plugin/skill
type Skill struct {
	Metadata Metadata
	Content  string // The full markdown content including or excluding frontmatter
	DirPath  string // The folder path where it resides
}

// Loader is responsible for loading skills from the filesystem
type Loader struct {
	baseDirs []string
}

// NewLoader creates a new SkillLoader instance that scans one or more base directories.
// Each directory is scanned for skill subdirectories containing SKILL.md files.
// If the same skill name exists in multiple directories, the last directory wins.
func NewLoader(baseDirs ...string) *Loader {
	return &Loader{baseDirs: baseDirs}
}

// LoadAll reads all SKILL.md files from the subdirectories of each base directory.
// Absent directories are silently skipped. The LoadAll signature is unchanged so that
// all existing call sites (e.g. input_pipeline.go) continue to work without modification.
func (l *Loader) LoadAll() (map[string]Skill, error) {
	skills := make(map[string]Skill)

	for _, base := range l.baseDirs {
		if err := l.loadFrom(base, skills); err != nil {
			return nil, err
		}
	}

	return skills, nil
}

// loadFrom scans a single baseDir and populates the skills map.
// If the directory does not exist it returns nil (graceful skip).
// Other OS errors are propagated.
func (l *Loader) loadFrom(baseDir string, skills map[string]Skill) error {
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		if os.IsNotExist(err) {
			// It's perfectly fine if the folder doesn't exist yet
			return nil
		}
		return fmt.Errorf("failed to read skills directory %q: %w", baseDir, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue // Only look at subdirectories
		}

		skillDir := filepath.Join(baseDir, entry.Name())
		skillFile := filepath.Join(skillDir, "SKILL.md")

		data, err := os.ReadFile(skillFile)
		if err != nil {
			if os.IsNotExist(err) {
				continue // Skip ghost folders without SKILL.md
			}
			continue // Ignore other read errors as per PRD "pula iteração do loop e silencia"
		}

		metadata, contentStr, err := parseFrontmatter(data)
		if err != nil {
			continue // Skip structural failure of frontmatter
		}

		// Use the yaml defined "name" as key; duplicate names: last dir wins
		skills[metadata.Name] = Skill{
			Metadata: metadata,
			Content:  contentStr,
			DirPath:  skillDir,
		}
	}

	return nil
}

// parseFrontmatter extracts the YAML frontmatter and the rest of the markdown
func parseFrontmatter(data []byte) (Metadata, string, error) {
	const separator = "---"
	var meta Metadata

	content := string(data)
	// Find the first occurrence of ---
	if bytes.HasPrefix(data, []byte(separator)) {
		parts := bytes.SplitN(data, []byte(separator), 3)
		if len(parts) >= 3 {
			// parts[0] is empty, parts[1] is frontmatter, parts[2] is the rest
			err := yaml.Unmarshal(parts[1], &meta)
			if err != nil {
				return meta, content, fmt.Errorf("invalid yaml: %w", err)
			}
			
			// If meta.Name is not provided, it's structurally bad
			if meta.Name == "" {
				return meta, content, fmt.Errorf("missing name in frontmatter")
			}

			return meta, content, nil // Return full content or strings.TrimSpace(string(parts[2]))
		}
	}

	return meta, content, fmt.Errorf("no frontmatter found")
}
