package tools

import (
	"context"
	"strings"
	"testing"

	"github.com/kocar/aurelia/internal/skill"
)

type stubSkillInstaller struct {
	lastRef   string
	lastScope string
	result    skill.InstallResult
	err       error
}

func (s *stubSkillInstaller) Install(ctx context.Context, skillRef, scope string) (skill.InstallResult, error) {
	s.lastRef = skillRef
	s.lastScope = scope
	return s.result, s.err
}

func TestInstallSkillTool_Execute_DefaultsToGlobalScope(t *testing.T) {
	installer := &stubSkillInstaller{
		result: skill.InstallResult{
			Scope:      skill.InstallScopeGlobal,
			TargetDir:  "/home/test/.aurelia/skills",
			SkillNames: []string{"Alpha Skill"},
		},
	}
	tool := NewInstallSkillTool(installer)

	out, err := tool.Execute(context.Background(), map[string]interface{}{
		"skill_ref": "demo/alpha",
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if installer.lastScope != skill.InstallScopeGlobal {
		t.Fatalf("expected default global scope, got %q", installer.lastScope)
	}
	if !strings.Contains(out, "Alpha Skill") || !strings.Contains(out, "`global`") {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestInstallSkillTool_Execute_ForwardsProjectScope(t *testing.T) {
	installer := &stubSkillInstaller{
		result: skill.InstallResult{
			Scope:      skill.InstallScopeProject,
			TargetDir:  "/repo/.aurelia/skills",
			SkillNames: []string{"Beta Skill"},
		},
	}
	tool := NewInstallSkillTool(installer)

	if _, err := tool.Execute(context.Background(), map[string]interface{}{
		"skill_ref": "demo/beta",
		"scope":     skill.InstallScopeProject,
	}); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if installer.lastScope != skill.InstallScopeProject {
		t.Fatalf("expected project scope, got %q", installer.lastScope)
	}
}


