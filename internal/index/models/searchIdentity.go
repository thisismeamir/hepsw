package models

import (
	"fmt"
	"strings"
)

type SearchPackageIdentity struct {
	Name    string
	Version string
}

func GetSearchPackageIdentity(input string) (*SearchPackageIdentity, error) {
	// Split the input into name and version
	parts := strings.SplitN(input, ":", 2)
	if len(parts) == 0 {
		return nil, fmt.Errorf("Invalid search Query: %s", input)
	}

	name := parts[0]
	version := "latest" // Default to latest if not specified
	if len(parts) > 1 {
		version = parts[1]
	}

	return &SearchPackageIdentity{
		Name:    name,
		Version: version,
	}, nil
}
