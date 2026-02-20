package reporters

import (
	"fmt"

	"github.com/thisismeamir/hepsw/internal/manifest"
	"gopkg.in/yaml.v3"
)

func generateJSONReport(m *manifest.Manifest) (string, error) {
	// Use existing YAML marshaling since manifest is already structured
	data, err := yaml.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("failed to marshal manifest: %w", err)
	}
	return string(data), nil
}
