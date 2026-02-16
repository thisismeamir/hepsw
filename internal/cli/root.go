package cli

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/thisismeamir/hepsw/internal/configuration"
	"github.com/thisismeamir/hepsw/internal/utils"
)

var (
	// Global flags
	verbose bool
	quiet   bool

	// Color functions (define once, use everywhere)
	colorError   = color.New(color.FgRed, color.Bold).SprintFunc()
	colorSuccess = color.New(color.FgGreen, color.Bold).SprintFunc()
	colorWarning = color.New(color.FgYellow).SprintFunc()
	colorInfo    = color.New(color.FgCyan).SprintFunc()
	colorHeader  = color.New(color.FgMagenta, color.Bold).SprintFunc()
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hepsw",
	Short: "HEP Software Stack Manager",
	Long: `HepSW is a source-first, reproducible software framework for building,
packaging, and composing High Energy Physics (HEP) software stacks on Linux systems.`,
	Version: "0.0.1",
	// Uncomment if you want a default action
	Run: func(cmd *cobra.Command, args []string) {
		utils.PrintHeader()
		fmt.Println("")
		PrintInfo("Run 'hepsw --help' to see available commands and options.")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Initialize runs before command execution
	cobra.OnInitialize(hepswInit)

	// Persistent flags (available to all subcommands)

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false,
		"enable verbose output")

	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false,
		"suppress non-essential output")

	// Add subcommands
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(checkConfig)
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(sampleGenCmd)
}

func hepswInit() {
	err := configuration.ConfigHealth()
	if err != nil {
		PrintWarning(err.Error())
	}
}

// Helper functions for colored output

func PrintSuccess(msg ...string) {
	fmt.Println(colorSuccess("[SUCC] "), msg)
}

func PrintError(msg ...string) {
	// TODO: What the hell is to be handled for a print function anyway?
	fprintln, err := fmt.Fprintln(os.Stderr, colorError("[ERR!] "), msg)
	if err != nil {
		return
	}
	_ = fprintln
}

func PrintWarning(msg ...string) {
	fmt.Println(colorWarning("[WARN] "), msg)
}

func PrintInfo(msg ...string) {
	if !quiet {
		fmt.Println(colorInfo("[INFO] "), msg)
	}
}

func PrintSection(msg ...string) {
	if !quiet {
		fmt.Println(colorHeader("=> "), msg)
	}
}

func PrintBullet(msg ...string) {
	if !quiet {
		fmt.Println(colorHeader("  ●"), msg)
	}
}
