package skill

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

type stubCommandRunner struct {
	lastDir  string
	lastName string
	lastArgs []string
	runErr   error
	onRun    func(dir string) error
}

func (s *stubCommandRunner) Run(ctx context.Context, dir string, name string, args ...string) error {
	s.lastDir = dir
	s.lastName = name
	s.lastArgs = append([]string(nil), args...)
	if s.runErr != nil {
		return s.runErr
	}
	if s.onRun != nil {
		return s.onRun(dir)
	}
	return nil
}

func TestInstaller_Install_GlobalScope(t *testing.T) {
	globalDir := t.TempDir()
	projectDir := t.TempDir()
	runner := &stubCommandRunner{
		onRun: func(dir string) error {
			writeSkillMD(t, filepath.Join(dir, ".downloaded"), "alpha", "Alpha Skill", "Installed from stub")
			return nil
		},
	}
	installer := &Installer{
		globalSkillsDir:  globalDir,
		projectSkillsDir: projectDir,
		runner:           runner,
	}

	result, err := installer.Install(context.Background(), "demo/alpha", InstallScopeGlobal)
	if err != nil {
		t.Fatalf("Install() error = %v", err)
	}

	if result.TargetDir != globalDir {
		t.Fatalf("expected global target dir, got %q", result.TargetDir)
	}
	if len(result.SkillNames) != 1 || result.SkillNames[0] != "Alpha Skill" {
		t.Fatalf("unexpected installed skills: %#v", result.SkillNames)
	}
	if _, err := os.Stat(filepath.Join(globalDir, "alpha-skill", "SKILL.md")); err != nil {
		t.Fatalf("expected installed skill in global dir: %v", err)
	}
}

func TestInstaller_Install_ProjectScope(t *testing.T) {
	globalDir := t.TempDir()
	projectDir := t.TempDir()
	runner := &stubCommandRunner{
		onRun: func(dir string) error {
			writeSkillMD(t, filepath.Join(dir, ".downloaded"), "beta", "Beta Skill", "Installed from stub")
			return nil
		},
	}
	installer := &Installer{
		globalSkillsDir:  globalDir,
		projectSkillsDir: projectDir,
		runner:           runner,
	}

	result, err := installer.Install(context.Background(), "demo/beta", InstallScopeProject)
	if err != nil {
		t.Fatalf("Install() error = %v", err)
	}

	if result.TargetDir != projectDir {
		t.Fatalf("expected project target dir, got %q", result.TargetDir)
	}
	if _, err := os.Stat(filepath.Join(projectDir, "beta-skill", "SKILL.md")); err != nil {
		t.Fatalf("expected installed skill in project dir: %v", err)
	}
}

func TestInstaller_Install_FailsWhenProjectScopeUnavailable(t *testing.T) {
	installer := &Installer{
		globalSkillsDir: t.TempDir(),
		runner:          &stubCommandRunner{},
	}

	if _, err := installer.Install(context.Background(), "demo/beta", InstallScopeProject); err == nil {
		t.Fatal("expected project scope to fail without project skills dir")
	}
}
