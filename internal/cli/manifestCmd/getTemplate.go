package manifestCmd

import (
	"fmt"

	"github.com/thisismeamir/hepsw/internal/manifest"
)

func getTemplate(templateName string) (*manifest.Manifest, error) {
	switch templateName {
	case "minimal":
		return getMinimalTemplate(), nil
	case "cmake":
		return getCMakeTemplate(), nil
	case "autotools":
		return getAutotoolsTemplate(), nil
	case "git":
		return getGitTemplate(), nil
	case "tarball":
		return getTarballTemplate(), nil
	default:
		return nil, fmt.Errorf("unknown template: %s", templateName)
	}
}
