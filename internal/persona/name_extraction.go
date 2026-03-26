package persona

import (
	"regexp"
	"strings"
)

var (
	nameMeChamoPattern = regexp.MustCompile(`(?i)\bme chamo\s+([A-Za-zÀ-ÿ][A-Za-zÀ-ÿ\s'-]{0,60})`)
	nameMeuNomePattern = regexp.MustCompile(`(?i)\bmeu nome e\s+([A-Za-zÀ-ÿ][A-Za-zÀ-ÿ\s'-]{0,60})`)
	nameSouPattern     = regexp.MustCompile(`(?i)\bsou\s+([A-Za-zÀ-ÿ][A-Za-zÀ-ÿ\s'-]{0,60})`)
)

// ExtractNameFromProfile extracts a user name from Portuguese profile text.
func ExtractNameFromProfile(profileText string) string {
	text := strings.TrimSpace(profileText)
	for _, pattern := range []*regexp.Regexp{nameMeChamoPattern, nameMeuNomePattern, nameSouPattern} {
		matches := pattern.FindStringSubmatch(text)
		if len(matches) < 2 {
			continue
		}

		name := strings.TrimSpace(matches[1])
		name = strings.TrimRight(name, ".,;:!?")
		for _, separator := range []string{" e quero ", " e prefiro ", " e gosto ", " e "} {
			if idx := strings.Index(strings.ToLower(name), separator); idx >= 0 {
				name = strings.TrimSpace(name[:idx])
				break
			}
		}
		name = strings.Join(strings.Fields(name), " ")
		if name != "" {
			return name
		}
	}

	return ""
}
