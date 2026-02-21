package manifestCmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thisismeamir/hepsw/internal/manifest"
	"github.com/thisismeamir/hepsw/internal/manifest/loader"
)

func runManifestEnv(cmd *cobra.Command, args []string) error {
	manifestSource := args[0]

	// Load manifest
	m, err := loader.LoadManifest(manifestSource)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	accessor := manifest.NewManifestAccessor(m)

	fmt.Printf("Environment Variables for %s@%s\n", accessor.Name(), accessor.Version())
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	showBuild := envScope == "all" || envScope == "build"
	showRuntime := envScope == "all" || envScope == "runtime"
	showSelf := envScope == "all" || envScope == "self"

	printed := false

	if showBuild {
		buildEnv := accessor.BuildEnvironment()
		if len(buildEnv) > 0 {
			fmt.Println("Build Environment:")
			for k, v := range buildEnv {
				fmt.Printf("  %s = %s\n", k, v)
			}
			fmt.Println()
			printed = true
		}
	}

	if showRuntime {
		runtimeEnv := accessor.RuntimeEnvironment()
		if len(runtimeEnv) > 0 {
			fmt.Println("Runtime Environment:")
			for k, v := range runtimeEnv {
				fmt.Printf("  %s = %s\n", k, v)
			}
			fmt.Println()
			printed = true
		}
	}

	if showSelf {
		selfEnv := accessor.SelfEnvironment()
		if len(selfEnv) > 0 {
			fmt.Println("Exported Environment (self):")
			for k, v := range selfEnv {
				fmt.Printf("  %s = %s\n", k, v)
			}
			fmt.Println()
			printed = true
		}
	}

	if !printed {
		fmt.Println("No environment variables defined.")
	}

	return nil
}
