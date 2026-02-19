package cli

import (
	"bufio"
	"os"

	"github.com/spf13/cobra"
	"github.com/thisismeamir/hepsw/internal/configuration"
)

var checkConfig = &cobra.Command{
	Use:   "check-config",
	Short: "Checks if the configuration is available, otherwise creates it.",
	Long:  `Checks if the configuration is available, otherwise creates it.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return ConfigInit()
	},
}

func ConfigInit() error {
	_, err := tryLoading()
	if err != nil {
		PrintError("Configuration file wasn't working. Do you want to create a new one? (y/n)")
		stdin := bufio.NewReader(os.Stdin)
		if stdin == nil {
			err := makeDefault()
			if err != nil {
				return err
			}
		} else if input, _ := stdin.ReadString('\n'); input == "y\n" {
			err := makeDefault()
			if err != nil {
				return err
			}
		} else {
			PrintError("Configuration file is not working, and you chose not to create a new one. Exiting.")
		}
		return nil
	}
	return nil
}

func makeDefault() error {
	return configuration.RestoreDefaultConfiguration()
}

func tryLoading() (*configuration.Configuration, error) {
	return configuration.GetConfiguration()
}
