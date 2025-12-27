package manifest

// PatchSpec defines a patch to apply to the source
type PatchSpec struct {
	URL      string `yaml:"url,omitempty"`
	Path     string `yaml:"path,omitempty"` // local path in index repo
	Checksum string `yaml:"checksum,omitempty"`
	Level    int    `yaml:"level,omitempty"` // patch -p level
}
