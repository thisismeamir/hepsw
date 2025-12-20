package cli

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thisismeamir/hepsw/internal/remote"
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"config file (default is $HOME/.hepsw.yaml)")

	rootCmd.PersistentFlags().StringVarP(&workspace, "workspace", "w", "",
		"path to HepSW workspace (default is $HEPSW_WORKSPACE or ./hepsw-workspace)")

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false,
		"enable verbose output")

	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false,
		"suppress non-essential output")

	// Bind flags to viper (allows reading from config files and env vars)
	_ = viper.BindPFlag("workspace", rootCmd.PersistentFlags().Lookup("workspace"))
	_ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	_ = viper.BindPFlag("quiet", rootCmd.PersistentFlags().Lookup("quiet"))

	// Add subcommands
	rootCmd.AddCommand(initCmd)
}

func hepswInit() {
	utils.PrintHeader()
	// Creating a hidden directory where we keep configurations, package index, etc.
	homeDir, err := os.UserHomeDir()
	if err != nil {
		PrintError("Unable to find user home directory: " + err.Error())
		os.Exit(1)
	}
	hepswConfigDir := homeDir + "/.hepsw"
	if _, err := os.Stat(hepswConfigDir); os.IsNotExist(err) {
		PrintInfo("Creating HepSW config directory: " + hepswConfigDir)
		err := os.MkdirAll(hepswConfigDir, 0755)
		if err != nil {
			PrintError("Unable to create HepSW config directory: " + err.Error())
			os.Exit(1)
		}
	}

	// HepSW package index directory
	hepswPackageIndexDir := hepswConfigDir + "/package-index"
	// Check if Git is available
	if !utils.CheckGit() {
		PrintError("Git is not installed or not available in PATH. Please install Git to proceed.")
		os.Exit(1)
	}
	if _, err := os.Stat(hepswPackageIndexDir); os.IsNotExist(err) {
		PrintWarning("No package index repository at " + hepswPackageIndexDir)
		PrintInfo("Cloning official HepSW package index into directory: " + hepswPackageIndexDir)
		err := remote.CloneRepo(hepswPackageIndexDir, "master")
		if err != nil {
			PrintError("Unable to clone HepSW package index: " + err.Error())
			os.Exit(1)
		}
		PrintSuccess("HepSW package index is cloned at " + hepswPackageIndexDir)
	} else {
		PrintInfo("HepSW package index found at " + hepswPackageIndexDir)
		// Checking Health:
		repo, err := remote.OpenRepo(hepswPackageIndexDir)
		if err != nil {
			PrintError("Unable to open HepSW package index repository: " + err.Error())
			os.Exit(1)
		}
		hasChanges, err := remote.HasLocalChanges(repo)
		if err != nil {
			PrintError("Unable to check local changes in HepSW package index: " + err.Error())
			os.Exit(1)
		}
		if hasChanges {
			PrintWarning("Local changes detected in HepSW package index. Resetting to HEAD.")
			err := remote.ResetLocalChanges(repo)
			if err != nil {
				PrintError("Unable to reset local changes in HepSW package index: " + err.Error())
				os.Exit(1)
			}
			PrintSuccess("Local changes reset in HepSW package index.")
		}
		// Check updates:
		PrintInfo("Checking for available updates")
		hasUpdate, err := remote.HasRemoteUpdates(repo, "master")
		if hasUpdate {
			PrintInfo("Updates are available for HepSW package index.")
			err := remote.FetchRemote(repo, "master")
			if err != nil {
				PrintError("Unable to fetch updates from HepSW package index: " + err.Error())
			}
			err = remote.PullChanges(repo, "master")
			if err != nil {
				PrintError("Unable to pull updates from HepSW package index: " + err.Error())
			}
		} else {
			PrintSuccess("HepSW package index is already up to date.")
		}
	}
}

// Helper functions for colored output
func PrintSuccess(msg string) {
	fmt.Println(colorSuccess("[SUCC]"), msg)
}

func PrintError(msg string) {
	fmt.Fprintln(os.Stderr, colorError("[ERR!]"), msg)
}

func PrintWarning(msg string) {
	fmt.Println(colorWarning("[WARN]"), msg)
}

func PrintInfo(msg string) {
	if !quiet {
		fmt.Println(colorInfo("[INFO]"), msg)
	}
}

func PrintSection(msg string) {
	if !quiet {
		fmt.Println(colorHeader("=====>"), msg)
	}
}
