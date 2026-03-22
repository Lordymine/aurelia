package agents

// BuildSDKAgents converts all agents in the registry to the format expected
// by the Claude Agent SDK's "agents" query option. Each agent becomes a map
// with keys: description, prompt, model (optional), tools (optional).
func BuildSDKAgents(r *Registry) map[string]any {
	all := r.Agents()
	if len(all) == 0 {
		return nil
	}
	result := make(map[string]any, len(all))
	for _, a := range all {
		def := map[string]any{
			"description": a.Description,
			"prompt":      a.Prompt,
		}
		if a.Model != "" {
			def["model"] = a.Model
		}
		if len(a.AllowedTools) > 0 {
			def["tools"] = a.AllowedTools
		}
		result[a.Name] = def
	}
	return result
}
