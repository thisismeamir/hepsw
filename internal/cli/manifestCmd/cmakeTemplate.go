package manifestCmd

import "github.com/thisismeamir/hepsw/internal/manifest"

func getCMakeTemplate() *manifest.Manifest {
	return &manifest.Manifest{
		Name:        "example",
		Version:     "1.0.0",
		Description: "A CMake-based example package",
		Source: manifest.SourceSpec{
			Type: "git",
			Url:  "https://github.com/example/example.git",
			Tag:  "v1.0.0",
		},
		Specifications: manifest.Specifications{
			Build: manifest.BuildSpecification{
				Toolchain: []manifest.Tool{
					{Name: "cmake", Version: ">=3.15"},
					{Name: "gcc", Version: ">=9.0"},
				},
			},
		},
		Recipe: manifest.Recipe{
			Configuration: []manifest.RecipeStep{
				{
					Name:    "Create build directory",
					Command: "mkdir -p build && cd build",
				},
				{
					Name:    "Configure with CMake",
					Command: "cmake .. -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX=${INSTALL_PREFIX}",
				},
			},
			Build: []manifest.RecipeStep{
				{
					Name:    "Build",
					Command: "cmake --build . -j${NCORES}",
				},
			},
			Install: []manifest.RecipeStep{
				{
					Name:    "Install",
					Command: "cmake --install .",
				},
			},
		},
	}
}
