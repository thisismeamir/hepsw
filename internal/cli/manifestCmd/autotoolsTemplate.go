package manifestCmd

import "github.com/thisismeamir/hepsw/internal/manifest"

func getAutotoolsTemplate() *manifest.Manifest {
	return &manifest.Manifest{
		Name:        "example",
		Version:     "1.0.0",
		Description: "An Autotools-based example package",
		Source: manifest.SourceSpec{
			Type: "tarball",
			Url:  "https://example.com/example-1.0.0.tar.gz",
		},
		Specifications: manifest.Specifications{
			Build: manifest.BuildSpecification{
				Toolchain: []manifest.Tool{
					{Name: "gcc", Version: ">=7.0"},
					{Name: "autoconf", Version: ">=2.69"},
					{Name: "automake", Version: ">=1.16"},
				},
			},
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
