package manifestCmd

import "github.com/thisismeamir/hepsw/internal/manifest"

func getGitTemplate() *manifest.Manifest {
	m := getMinimalTemplate()
	m.Source = manifest.SourceSpec{
		Type: "git",
		Url:  "https://github.com/example/example.git",
		Tag:  "main",
	}
	return m
}
