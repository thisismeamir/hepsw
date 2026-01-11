package manifest

import (
	"path/filepath"

	"github.com/spf13/viper"
)

func (manifest *Manifest) SaveManifest(name string, path string) error {
	save := viper.New()
	save.SetConfigName(name)
	save.SetConfigType("yaml")

	save.Set("name", manifest.Name)
	save.Set("version", manifest.Version)
	save.Set("description", manifest.Description)
	save.Set("source", manifest.Source)
	save.Set("metadata", manifest.Metadata)
	save.Set("specifications", manifest.Specifications)
	save.Set("recipe", manifest.Recipe)

	err := save.WriteConfigAs(filepath.Join(filepath.Dir(path), name))
	if err != nil {
		return err
	}
	return nil
}
