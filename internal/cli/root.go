package cli

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thisismeamir/hepsw/utils"
)

var (
	// Global flags
	cfgFile   string
	workspace string
	verbose   bool
	quiet     bool

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
		printInfo("Run 'hepsw --help' to see available commands and options.")
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"config file (default is $HOME/.hepsw.yaml)")

	rootCmd.PersistentFlags().StringVarP(&workspace, "workspace", "w", "",
		"path to HepSW workspace (default is $HEPSW_WORKSPACE or ./hepsw-workspace)")

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false,
		"enable verbose output")

	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false,
		"suppress non-essential output")

	// Bind flags to viper (allows reading from config files and env vars)
	viper.BindPFlag("workspace", rootCmd.PersistentFlags().Lookup("workspace"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("quiet", rootCmd.PersistentFlags().Lookup("quiet"))

	// Add subcommands
	rootCmd.AddCommand(initCmd)
}

func hepswInit() {
	// Creating a hidden directory where we keep configurations, package index, etc.
	homeDir, err := os.UserHomeDir()
	if err != nil {
		printError("Unable to find user home directory: " + err.Error())
		os.Exit(1)
	}
	hepswConfigDir := homeDir + "/.hepsw"
	if _, err := os.Stat(hepswConfigDir); os.IsNotExist(err) {
		printInfo("Creating HepSW config directory: " + hepswConfigDir)
		err := os.MkdirAll(hepswConfigDir, 0755)
		if err != nil {
			printError("Unable to create HepSW config directory: " + err.Error())
			os.Exit(1)
		}
	}

	// Fetching package index from GitHub repo if not already present
	packageIndexPath := hepswConfigDir + "/package-index.yaml"
	if _, err := os.Stat(packageIndexPath); os.IsNotExist(err) {
		printInfo("Fetching package index...")
		err := utils.FetchPackageIndex(packageIndexPath)
		if err != nil {
			printError("Unable to fetch package index: " + err.Error())
			os.Exit(1)
		}
		printSuccess("Package index fetched successfully.")
	}
}

// Helper functions for colored output
func printSuccess(msg string) {
	fmt.Println(colorSuccess("[SUCC]"), msg)
}

func printError(msg string) {
	fmt.Fprintln(os.Stderr, colorError("[ERROR]"), msg)
}

func printWarning(msg string) {
	fmt.Println(colorWarning("[WARN]"), msg)
}

func printInfo(msg string) {
	if !quiet {
		fmt.Println(colorInfo("[INFO]"), msg)
	}
}

func printHeader(msg string) {
	if !quiet {
		fmt.Println(colorHeader("=====>"), msg)
	}
}
