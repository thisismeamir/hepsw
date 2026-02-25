package manifestCmd

import "github.com/thisismeamir/hepsw/internal/manifest"

func getMinimalTemplate() *manifest.Manifest {
	return &manifest.Manifest{
		Name:        "example",
		Version:     "1.0.0",
		Description: "A simple example package",
		Source: manifest.SourceSpec{
			Type: "tarball",
			Url:  "https://example.com/example-1.0.0.tar.gz",
		},
		Recipe: manifest.Recipe{
			Configuration: []manifest.RecipeStep{
				{
					Name:    "Configure",
					Command: "./configure --prefix=${INSTALL_PREFIX}",
				},
			},
			Build: []manifest.RecipeStep{
				{
					Name:    "Build",
					Command: "make -j${NCORES}",
				},
			},
			Install: []manifest.RecipeStep{
				{
					Name:    "Install",
					Command: "make install",
				},
			},
		},
	}
}
