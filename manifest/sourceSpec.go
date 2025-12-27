package manifest

import manifest2 "github.com/thisismeamir/hepsw/internal/manifest"

// SourceSpec defines where and how to fetch the package source
type SourceSpec struct {
	Type     string                `yaml:"type"` // git, tarball, svn, etc.
	URL      string                `yaml:"url"`
	Ref      string                `yaml:"ref,omitempty"`      // branch, tag, or commit
	Checksum string                `yaml:"checksum,omitempty"` // for tarballs
	Patches  []manifest2.PatchSpec `yaml:"patches,omitempty"`
	Mirror   []string              `yaml:"mirror,omitempty"` // fallback URLs
	Auth     *manifest2.AuthSpec   `yaml:"auth,omitempty"`
}
