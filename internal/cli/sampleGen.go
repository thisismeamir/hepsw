package cli

import (
	"github.com/spf13/cobra"
	"github.com/thisismeamir/hepsw/internal/utils"
)

var sampleGenCmd = &cobra.Command{
	Use:   "sample-gen",
	Short: "Generate sample manifest file for testing.",
	Long: `Generate sample manifest file for testing.
hepsw sample-gen <name> <path>
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Getting the first argument as the name and second as the path:
		if len(args) != 2 {
			PrintError("You need to specify a name and a path")
			return nil
		}

		manifest := utils.CreateSampleManifest()

		err := manifest.SaveManifest(args[0], args[1])

		return err
	},
}
