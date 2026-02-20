package cli

import (
	"github.com/spf13/cobra"
	"github.com/thisismeamir/hepsw/internal/configuration"
	"github.com/thisismeamir/hepsw/internal/index/client"
)

var syncCmd = &cobra.Command{
	Use:     "sync",
	Short:   "Synchronizes local database with remote database",
	Long:    `Synchronizes local database in ~/.hepsw/index.db with the remote database online, ensuring we have access to updates, and new software that can be installed.`,
	Version: "0.0.1",
	RunE: func(cmd *cobra.Command, args []string) error {
		return sync()
	},
}

func sync() error {
	config, configErr := configuration.GetConfiguration()
	if configErr != nil {
		return configErr
	}

	cl, err := client.New(&config.IndexConfig)

	if err != nil {
		return err
	}

	return cl.Sync(config)
}
