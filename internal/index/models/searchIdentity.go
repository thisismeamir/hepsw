package models

import (
	"fmt"
	"strings"
)

type searchPackageIdentity struct {
	Name    string
	Version string
}

func GetSearchPackageIdentity(input string) (*searchPackageIdentity, error) {
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

	return &searchPackageIdentity{
		Name:    name,
		Version: version,
	}, nil
}
