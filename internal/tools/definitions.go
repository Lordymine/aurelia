package tools

import "github.com/kocar/aurelia/internal/agent"

func ReadFileDefinition() agent.Tool {
	return agent.Tool{
		Name:        "read_file",
		Description: "Le o conteudo de um arquivo local. Se `workdir` for informado, caminhos relativos serao resolvidos a partir dele.",
		JSONSchema: objectSchema(
			map[string]any{
				"path":    stringProperty(""),
				"workdir": stringProperty("Diretorio base opcional para resolver caminhos relativos."),
			},
			"path",
		),
	}
}

func WriteFileDefinition() agent.Tool {
	return agent.Tool{
		Name:        "write_file",
		Description: "Escreve conteudo integral em um arquivo. Se `workdir` for informado, caminhos relativos serao resolvidos a partir dele.",
		JSONSchema: objectSchema(
			map[string]any{
				"path":    stringProperty(""),
				"content": stringProperty(""),
				"workdir": stringProperty("Diretorio base opcional para resolver caminhos relativos."),
			},
			"path",
			"content",
		),
	}
}

func ListDirDefinition() agent.Tool {
	return agent.Tool{
		Name:        "list_dir",
		Description: "Lista os arquivos dentro de um diretorio. Se `workdir` for informado, caminhos relativos serao resolvidos a partir dele.",
		JSONSchema: objectSchema(
			map[string]any{
				"path":    stringProperty(""),
				"workdir": stringProperty("Diretorio base opcional para resolver caminhos relativos."),
			},
			"path",
		),
	}
}

func WebSearchDefinition() agent.Tool {
	return agent.Tool{
		Name:        "web_search",
		Description: "Pesquisa na internet usando a engine do DuckDuckGo e extrai resultados textuais.",
		JSONSchema: objectSchema(
			map[string]any{
				"query": stringProperty(""),
				"count": numberProperty("Maximo de resultados (ate 10)"),
			},
			"query",
		),
	}
}

func RunCommandDefinition() agent.Tool {
	return agent.Tool{
		Name:        "run_command",
		Description: "Executa um comando local de forma controlada e retorna stdout, stderr, exit code e timeout em JSON. Em Windows, use sintaxe de PowerShell no comando e prefira informar `workdir` ao operar em outro projeto.",
		JSONSchema: objectSchema(
			map[string]any{
				"command":         stringProperty(""),
				"workdir":         stringProperty(""),
				"timeout_seconds": numberProperty(""),
			},
			"command",
		),
	}
}

func RegisterCoreTools(registry *agent.ToolRegistry) {
	if registry == nil {
		return
	}

	registry.Register(ReadFileDefinition(), ReadFileHandler)
	registry.Register(WriteFileDefinition(), WriteFileHandler)
	registry.Register(ListDirDefinition(), ListDirHandler)
	registry.Register(WebSearchDefinition(), WebSearchHandler)
	registry.Register(RunCommandDefinition(), RunCommandHandler)
}

func objectSchema(properties map[string]any, required ...string) map[string]any {
	schema := map[string]any{
		"type":       "object",
		"properties": properties,
	}
	if len(required) > 0 {
		schema["required"] = required
	}
	return schema
}

func stringProperty(description string) map[string]any {
	property := map[string]any{"type": "string"}
	if description != "" {
		property["description"] = description
	}
	return property
}

func numberProperty(description string) map[string]any {
	property := map[string]any{"type": "number"}
	if description != "" {
		property["description"] = description
	}
	return property
}


