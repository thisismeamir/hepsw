package utils

import "fmt"

func FormatList(items []string) string {
	if len(items) == 0 {
		return "(none)"
	}
	result := ""
	for _, item := range items {
		result += fmt.Sprintf("  - %s\n", item)
	}
	return result
}
