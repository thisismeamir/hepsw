package manifest

// AuthSpec defines authentication for private repositories
type AuthSpec struct {
	Type        string `yaml:"type"` // token, ssh-key, basic
	SecretRef   string `yaml:"secretRef,omitempty"`
	UsernameRef string `yaml:"usernameRef,omitempty"`
	PasswordRef string `yaml:"passwordRef,omitempty"`
}
