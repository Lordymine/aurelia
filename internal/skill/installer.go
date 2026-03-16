package skill

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
)

const (
	InstallScopeGlobal  = "global"
	InstallScopeProject = "project"
)

type InstallResult struct {
	Scope      string
	TargetDir  string
	SkillNames []string
}

type commandRunner interface {
	Run(ctx context.Context, dir string, name string, args ...string) error
}

type execRunner struct{}

func (execRunner) Run(ctx context.Context, dir string, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(output)))
	}
	return nil
}

type Installer struct {
	globalSkillsDir  string
	projectSkillsDir string
	runner           commandRunner
}

func NewInstaller(globalSkillsDir, projectSkillsDir string) *Installer {
	return &Installer{
		globalSkillsDir:  globalSkillsDir,
		projectSkillsDir: projectSkillsDir,
		runner:           execRunner{},
	}
}

func (i *Installer) Install(ctx context.Context, skillRef, scope string) (InstallResult, error) {
	skillRef = strings.TrimSpace(skillRef)
	if skillRef == "" {
		return InstallResult{}, fmt.Errorf("skill_ref is required")
	}

	targetDir, err := i.resolveTargetDir(scope)
	if err != nil {
		return InstallResult{}, err
	}
	if err := os.MkdirAll(targetDir, 0700); err != nil {
		return InstallResult{}, fmt.Errorf("create target skills directory: %w", err)
	}

	tempDir, err := os.MkdirTemp("", "aurelia-skill-install-*")
	if err != nil {
		return InstallResult{}, fmt.Errorf("create temp skill install dir: %w", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	if err := i.runner.Run(ctx, tempDir, npxExecutable(), "-y", "skills", "add", skillRef); err != nil {
		return InstallResult{}, fmt.Errorf("install skill with skills.sh CLI: %w", err)
	}

	skillDirs, err := discoverInstalledSkills(tempDir)
	if err != nil {
		return InstallResult{}, err
	}
	if len(skillDirs) == 0 {
		return InstallResult{}, fmt.Errorf("skills CLI completed but no SKILL.md was found in the downloaded content")
	}

	result := InstallResult{Scope: scope, TargetDir: targetDir}
	for _, skillDir := range skillDirs {
		metadata, _, err := parseFrontmatterFromPath(filepath.Join(skillDir, "SKILL.md"))
		if err != nil {
			return InstallResult{}, err
		}
		destDir := filepath.Join(targetDir, filepath.Base(skillDir))
		if metadata.Name != "" {
			destDir = filepath.Join(targetDir, sanitizeSkillDirName(metadata.Name))
		}
		if err := copyDir(skillDir, destDir); err != nil {
			return InstallResult{}, err
		}
		result.SkillNames = append(result.SkillNames, metadata.Name)
	}
	slices.Sort(result.SkillNames)
	return result, nil
}

func (i *Installer) resolveTargetDir(scope string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(scope)) {
	case "", InstallScopeGlobal:
		if strings.TrimSpace(i.globalSkillsDir) == "" {
			return "", fmt.Errorf("global skills directory is not configured")
		}
		return i.globalSkillsDir, nil
	case InstallScopeProject:
		if strings.TrimSpace(i.projectSkillsDir) == "" {
			return "", fmt.Errorf("project-local skills directory is not available in the current working directory")
		}
		return i.projectSkillsDir, nil
	default:
		return "", fmt.Errorf("unsupported install scope %q", scope)
	}
}

func discoverInstalledSkills(root string) ([]string, error) {
	var skillDirs []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.EqualFold(d.Name(), "SKILL.md") {
			return nil
		}
		skillDirs = append(skillDirs, filepath.Dir(path))
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("scan installed skill content: %w", err)
	}
	return skillDirs, nil
}

func parseFrontmatterFromPath(path string) (Metadata, string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Metadata{}, "", fmt.Errorf("read installed skill metadata: %w", err)
	}
	meta, content, err := parseFrontmatter(data)
	if err != nil {
		return Metadata{}, "", fmt.Errorf("parse installed skill metadata: %w", err)
	}
	return meta, content, nil
}

func sanitizeSkillDirName(name string) string {
	name = strings.TrimSpace(strings.ToLower(name))
	replacer := strings.NewReplacer(" ", "-", "/", "-", "\\", "-", ":", "-", "*", "-", "?", "-", "\"", "-", "<", "-", ">", "-", "|", "-")
	name = replacer.Replace(name)
	name = strings.Trim(name, "-.")
	if name == "" {
		return "skill"
	}
	return name
}

func copyDir(srcDir, dstDir string) error {
	if err := os.RemoveAll(dstDir); err != nil {
		return fmt.Errorf("replace installed skill directory: %w", err)
	}
	if err := filepath.Walk(srcDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dstDir, rel)
		if info.IsDir() {
			return os.MkdirAll(target, 0700)
		}
		return copyFile(path, target, info.Mode())
	}); err != nil {
		return fmt.Errorf("copy installed skill: %w", err)
	}
	return nil
}

func copyFile(src, dst string, mode fs.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0700); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = in.Close() }()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}

func npxExecutable() string {
	if runtime.GOOS == "windows" {
		return "npx.cmd"
	}
	return "npx"
}


