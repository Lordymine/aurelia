package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/kocar/aurelia/internal/agent"
	"github.com/kocar/aurelia/internal/skill"
)

type SkillInstaller interface {
	Install(ctx context.Context, skillRef, scope string) (skill.InstallResult, error)
}

type InstallSkillTool struct {
	installer SkillInstaller
}

func NewInstallSkillTool(installer SkillInstaller) *InstallSkillTool {
	return &InstallSkillTool{installer: installer}
}

func (t *InstallSkillTool) Definition() agent.Tool {
	return agent.Tool{
		Name:        "install_skill",
		Description: "Baixa e instala uma skill do catalogo skills.sh no destino global (~/.aurelia/skills) ou no destino local do projeto atual (.aurelia/skills).",
		JSONSchema: objectSchema(
			map[string]any{
				"skill_ref": stringProperty("Referencia da skill no catalogo skills.sh, por exemplo owner/skill-name."),
				"scope":     stringProperty("Destino da instalacao: `global` ou `project`. O padrao e `global`."),
			},
			"skill_ref",
		),
	}
}

func (t *InstallSkillTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	if t.installer == nil {
		return "", fmt.Errorf("skill installer is not configured")
	}

	skillRef := optionalStringArg(args, "skill_ref")
	scope := optionalStringArg(args, "scope")
	if strings.TrimSpace(scope) == "" {
		scope = skill.InstallScopeGlobal
	}

	result, err := t.installer.Install(ctx, skillRef, scope)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"Instalei %d skill(s) do catalogo no escopo `%s`.\nDestino: `%s`.\nSkills: %s.",
		len(result.SkillNames),
		result.Scope,
		result.TargetDir,
		strings.Join(result.SkillNames, ", "),
	), nil
}


