package utils

import (
	"fmt"

	"github.com/thisismeamir/hepsw/internal/manifest"
)

func FormatDependencies(deps []manifest.Dependency) string {
	if len(deps) == 0 {
		return "(none)"
	}
	result := ""
	for _, dep := range deps {
		result += fmt.Sprintf("  - %s %s", dep.Name, dep.Version)
		if len(dep.ForOptions) > 0 {
			result += fmt.Sprintf(" (for: %v)", dep.ForOptions)
		}
		result += "\n"
	}
	return result
}
