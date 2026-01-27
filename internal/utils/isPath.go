package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func IsFilePath(input string) bool {
	ext := strings.ToLower(filepath.Ext(input))
	if ext != ".yaml" && ext != ".yml" {
		return false
	}

	info, err := os.Stat(input)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
