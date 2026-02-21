package utils

import (
	"fmt"

	"github.com/thisismeamir/hepsw/internal/manifest"
)

func FormatToolChain(tools []manifest.Tool) string {
	if len(tools) == 0 {
		return "(none)"
	}
	result := ""
	for _, tool := range tools {
		result += fmt.Sprintf("  - %s %s\n", tool.Name, tool.Version)
	}
	return result
}
