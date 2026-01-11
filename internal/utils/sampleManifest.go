package utils

import (
	"github.com/thisismeamir/hepsw/internal/manifest"
)

func CreateSampleManifest() manifest.Manifest {
	sample := manifest.Manifest{
		Name:        "sample",
		Version:     "1.1.1",
		Description: "A Sample Manifest to Start from",
		Source: manifest.SourceSpec{
			Type: "git",
			Url:  "https://github.com/example/sample.git",
		},
		Metadata: manifest.ManifestMetaData{
			Authors:       []string{"Kid A", "Amir H. Ebrahimnezhad"},
			Homepage:      "https://github.com/thisismeamir/hepsw",
			License:       "Apache-2.0",
			Documentation: "https://hepsw.readthedocs.io",
		},
		Specifications: manifest.Specifications{
			Build:       manifest.BuildSpecification{},
			Runtime:     manifest.RuntimeSpecification{},
			Environment: manifest.EnvironmentSpecification{},
		},
		Recipe: manifest.Recipe{
			Configuration: []manifest.RecipeStep{
				{
					Name: "Set X",
					Set: map[string]string{
						"X":    "Y",
						"Here": "There",
					},
				},
				{
					Name:    "Run X",
					Command: "x",
				},
			},
			Build: []manifest.RecipeStep{
				{
					Name: "Set X",
					Set: map[string]string{
						"X":    "Y",
						"Here": "There",
					},
				},
				{
					Name:    "Run X",
					Command: "x",
				},
			},
			Install: []manifest.RecipeStep{
				{
					Name: "Install",
					Set: map[string]string{
						"X":    "Install",
						"Here": "There",
					},
				},
				{
					Name:    "Run X",
					Command: "x",
				},
			},
			Use: []manifest.RecipeStep{
				{
					Name: "Use Case",
					Set: map[string]string{
						"X":    "Y",
						"Here": "There",
					},
				},
				{
					Name:    "Case Use",
					Command: "x",
				},
			},
		},
	}
	return sample
}
