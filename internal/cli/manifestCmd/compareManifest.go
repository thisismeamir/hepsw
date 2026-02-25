package manifestCmd

import (
	"fmt"

	"github.com/thisismeamir/hepsw/internal/manifest"
)

func compareManifests(m1, m2 *manifest.Manifest) []string {
	diffs := make([]string, 0)

	if m1.Name != m2.Name {
		diffs = append(diffs, fmt.Sprintf("Name: %s → %s", m1.Name, m2.Name))
	}

	if m1.Version != m2.Version {
		diffs = append(diffs, fmt.Sprintf("Version: %s → %s", m1.Version, m2.Version))
	}

	if m1.Description != m2.Description {
		diffs = append(diffs, "Description changed")
	}

	if m1.Source.Type != m2.Source.Type {
		diffs = append(diffs, fmt.Sprintf("Source type: %s → %s", m1.Source.Type, m2.Source.Type))
	}

	if m1.Source.Url != m2.Source.Url {
		diffs = append(diffs, "Source URL changed")
	}

	// Compare dependencies
	if len(m1.Specifications.Build.Dependencies) != len(m2.Specifications.Build.Dependencies) {
		diffs = append(diffs, fmt.Sprintf("Build dependencies count: %d → %d",
			len(m1.Specifications.Build.Dependencies),
			len(m2.Specifications.Build.Dependencies)))
	}

	// Compare recipe steps
	if len(m1.Recipe.Build) != len(m2.Recipe.Build) {
		diffs = append(diffs, fmt.Sprintf("Build steps count: %d → %d",
			len(m1.Recipe.Build), len(m2.Recipe.Build)))
	}

	return diffs
}
