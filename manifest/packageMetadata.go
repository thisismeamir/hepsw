package manifest

import "time"

// PackageMetadata contains identifying and descriptive information
type PackageMetadata struct {
	Name        string            `yaml:"name"`
	Version     string            `yaml:"version"`
	Description string            `yaml:"description,omitempty"`
	Homepage    string            `yaml:"homepage,omitempty"`
	License     string            `yaml:"license,omitempty"`
	Maintainers []string          `yaml:"maintainers,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty"`
	CreatedAt   time.Time         `yaml:"createdAt,omitempty"`
	UpdatedAt   time.Time         `yaml:"updatedAt,omitempty"`
}
