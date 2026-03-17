package skill

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeSkillMD creates a skill subdirectory with a SKILL.md file inside baseDir.
func writeSkillMD(t *testing.T, baseDir, skillDirName, name, description string) {
	t.Helper()
	skillDir := filepath.Join(baseDir, skillDirName)
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatalf("create skill dir: %v", err)
	}
	content := "---\nname: " + name + "\ndescription: " + description + "\n---\n\nSkill body for " + name + ".\n"
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0o644); err != nil {
		t.Fatalf("write SKILL.md: %v", err)
	}
}

// TestLoader_MultiDir: two temp dirs, each with one skill — both should appear in the map.
func TestLoader_MultiDir(t *testing.T) {
	dir1 := t.TempDir()
	dir2 := t.TempDir()

	writeSkillMD(t, dir1, "alpha", "Alpha Skill", "Does alpha things")
	writeSkillMD(t, dir2, "beta", "Beta Skill", "Does beta things")

	skills, err := NewLoader(dir1, dir2).LoadAll()
	if err != nil {
		t.Fatalf("LoadAll returned error: %v", err)
	}
	if _, ok := skills["Alpha Skill"]; !ok {
		t.Errorf("expected 'Alpha Skill' in map, got keys: %v", mapKeys(skills))
	}
	if _, ok := skills["Beta Skill"]; !ok {
		t.Errorf("expected 'Beta Skill' in map, got keys: %v", mapKeys(skills))
	}
	if len(skills) != 2 {
		t.Errorf("expected 2 skills, got %d", len(skills))
	}
}

// TestLoader_EmptyAllDirs: three empty temp dirs — must return empty map, nil error.
func TestLoader_EmptyAllDirs(t *testing.T) {
	dir1 := t.TempDir()
	dir2 := t.TempDir()
	dir3 := t.TempDir()

	skills, err := NewLoader(dir1, dir2, dir3).LoadAll()
	if err != nil {
		t.Fatalf("LoadAll returned error: %v", err)
	}
	if len(skills) != 0 {
		t.Errorf("expected empty map, got %d skills", len(skills))
	}
}

// TestLoader_SkipsAbsentDir: nonexistent path + valid dir — valid dir skills loaded, absent dir silently skipped, nil error.
func TestLoader_SkipsAbsentDir(t *testing.T) {
	validDir := t.TempDir()
	writeSkillMD(t, validDir, "gamma", "Gamma Skill", "Does gamma things")

	skills, err := NewLoader("/nonexistent/path/should/not/exist", validDir).LoadAll()
	if err != nil {
		t.Fatalf("LoadAll returned error: %v", err)
	}
	if _, ok := skills["Gamma Skill"]; !ok {
		t.Errorf("expected 'Gamma Skill' in map, got keys: %v", mapKeys(skills))
	}
	if len(skills) != 1 {
		t.Errorf("expected 1 skill, got %d", len(skills))
	}
}

// TestLoader_SingleDir: single dir, single skill — backward-compatible behavior.
func TestLoader_SingleDir(t *testing.T) {
	dir := t.TempDir()
	writeSkillMD(t, dir, "delta", "Delta Skill", "Does delta things")

	skills, err := NewLoader(dir).LoadAll()
	if err != nil {
		t.Fatalf("LoadAll returned error: %v", err)
	}
	if _, ok := skills["Delta Skill"]; !ok {
		t.Errorf("expected 'Delta Skill' in map, got keys: %v", mapKeys(skills))
	}
	if len(skills) != 1 {
		t.Errorf("expected 1 skill, got %d", len(skills))
	}
}

// TestLoader_DuplicateName: same skill name in two dirs — last dir (later in baseDirs) wins.
func TestLoader_DuplicateName(t *testing.T) {
	dir1 := t.TempDir()
	dir2 := t.TempDir()

	// Both dirs have a skill named "Shared Skill", but different descriptions.
	writeSkillMD(t, dir1, "shared", "Shared Skill", "From dir1")
	writeSkillMD(t, dir2, "shared", "Shared Skill", "From dir2")

	skills, err := NewLoader(dir1, dir2).LoadAll()
	if err != nil {
		t.Fatalf("LoadAll returned error: %v", err)
	}
	if len(skills) != 1 {
		t.Errorf("expected 1 skill (last wins), got %d", len(skills))
	}
	s, ok := skills["Shared Skill"]
	if !ok {
		t.Fatal("expected 'Shared Skill' in map")
	}
	// dir2 should win — DirPath should be inside dir2
	rel, err := filepath.Rel(dir2, s.DirPath)
	if err != nil {
		t.Fatalf("filepath.Rel() error = %v", err)
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		t.Errorf("expected DirPath under dir2 (%s), got %s", dir2, s.DirPath)
	}
}

// mapKeys is a helper to get sorted key list for error messages.
func mapKeys(m map[string]Skill) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
