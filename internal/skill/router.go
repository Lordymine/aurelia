package skill

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/kocar/aurelia/internal/agent"
)

// Router handles intent classification to pick a specific skill
type Router struct {
	llm agent.LLMProvider
}

// NewRouter constructs a router
func NewRouter(llm agent.LLMProvider) *Router {
	return &Router{llm: llm}
}

// Route identifies if any given skill is suitable for the user prompt
func (r *Router) Route(ctx context.Context, prompt string, availableSkills map[string]Skill) (string, error) {
	if len(availableSkills) == 0 {
		return "", nil // no skills to pick
	}

	var descriptions []string
	for name, skill := range availableSkills {
		descriptions = append(descriptions, fmt.Sprintf("- Name: %s\n  Desc: %s", name, skill.Metadata.Description))
	}

	systemPrompt := fmt.Sprintf(`You are a precise classifier. You are evaluating a user query against a list of available skills.
If the query exactly matches the intent of a skill, return ONLY valid JSON containing the skillName.
If no skill matches meaningfully, return {"skillName": null}.

Available skills:
%s

Example Output:
{"skillName": "git-manager"}
`, strings.Join(descriptions, "\n"))

	history := []agent.Message{
		{Role: "user", Content: prompt},
	}

	// Request from LLM without tools just as a fast classifier (Passo Zero)
	resp, err := r.llm.GenerateContent(ctx, systemPrompt, history, nil)
	if err != nil {
		log.Printf("Router error: %v", err)
		return "", nil // Fallback gracefully to null intent
	}

	raw := resp.Content

	// Sometimes LLMs wrap json in markdown
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimSuffix(raw, "```")
	raw = strings.TrimSpace(raw)

	var output struct {
		SkillName *string `json:"skillName"`
	}

	if err := json.Unmarshal([]byte(raw), &output); err != nil {
		log.Printf("Router parse error: %v (Raw: %s)", err, raw)
		return "", nil
	}

	if output.SkillName == nil {
		return "", nil
	}

	requested := *output.SkillName
	if _, ok := availableSkills[requested]; ok {
		return requested, nil
	}

	return "", nil
}


