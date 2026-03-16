package persona

import (
	"bytes"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func readPersonaFile(path, kind string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s file: %w", kind, err)
	}
	return content, nil
}

func parseIdentityFrontmatter(identityBytes []byte) (Config, string, error) {
	var config Config

	parts := bytes.SplitN(identityBytes, []byte("---"), 3)
	if len(parts) != 3 {
		return config, string(bytes.TrimSpace(identityBytes)), nil
	}

	if err := yaml.Unmarshal(parts[1], &config); err != nil {
		return Config{}, "", fmt.Errorf("failed to parse yaml frontmatter: %w", err)
	}
	return config, string(bytes.TrimSpace(parts[2])), nil
}
