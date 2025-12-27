package manifest

// SourceSpec defines where and how to fetch the package source
type SourceSpec struct {
	Type     string      `yaml:"type"` // git, tarball, svn, etc.
	URL      string      `yaml:"url"`
	Ref      string      `yaml:"ref,omitempty"`      // branch, tag, or commit
	Checksum string      `yaml:"checksum,omitempty"` // for tarballs
	Patches  []PatchSpec `yaml:"patches,omitempty"`
	Mirror   []string    `yaml:"mirror,omitempty"` // fallback URLs
	Auth     *AuthSpec   `yaml:"auth,omitempty"`
}
