package manifestCmd

import "github.com/thisismeamir/hepsw/internal/manifest"

func createRecipeForBuildSystem(buildSystem string) manifest.Recipe {
	switch buildSystem {
	case "cmake":
		return getCMakeTemplate().Recipe
	case "autotools":
		return getAutotoolsTemplate().Recipe
	default:
		return getMinimalTemplate().Recipe
	}
}
