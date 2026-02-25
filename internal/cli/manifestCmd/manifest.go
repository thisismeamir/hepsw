package manifestCmd

import (
	"github.com/spf13/cobra"
)

// ManifestCmd represents the manifest command
var ManifestCmd = &cobra.Command{
	Use:   "manifest",
	Short: "Manifest management commands",
	Long: `Commands for creating, validating, inspecting, and managing HepSW manifests.
	
Manifests can be loaded from local files or from the HepSW Package Index Repository.
Use path/to/manifest.yaml for local files or package@version for index references.`,
}

var (
	fetchDestination    string
	stripMinimal        bool
	manifestNewTemplate string
	manifestNewOutput   string
	depsFormat          string
	depsShowOptions     bool
	envScope            string
	formatIndent        int
	exportFormat        string
	exportOutput        string
	flattenOptions      []string
	flattenOutput       string
	diffFormat          string
	validateVerbose     bool
	walkOptions         []string
	walkVariables       map[string]string
	showFormat          string
)

// manifestFetchCmd fetches manifest from registry
var manifestFetchCmd = &cobra.Command{
	Use:   "fetch [package@version]",
	Short: "Download a manifest from a registry or repository",
	Long: `Fetch a manifest from the HepSW Package Index Repository.
	
Example:
  hepsw manifest fetch root@6.30.02
  hepsw manifest fetch python@3.11`,
	Args: cobra.ExactArgs(1),
	RunE: runManifestFetch,
}

// manifestStripCmd strips non-essential fields
var manifestStripCmd = &cobra.Command{
	Use:   "strip [manifest]",
	Short: "Remove non-essential fields to produce a minimal manifest",
	Long: `Strip a manifest to include only essential fields.
	
This creates a minimal version suitable for distribution or as a template.`,
	Args: cobra.ExactArgs(1),
	RunE: runManifestStrip,
}

// manifestExplainCmd explains dependency inclusion
var manifestExplainCmd = &cobra.Command{
	Use:   "explain [manifest] [dependency]",
	Short: "Explain why a dependency, option, or variable is included or excluded",
	Long: `Provide detailed explanation for why a specific dependency or option
is included in the build.`,
	Args: cobra.ExactArgs(2),
	RunE: runManifestExplain,
}

// manifestGraphCmd generates dependency graph
var manifestGraphCmd = &cobra.Command{
	Use:   "graph [manifest]",
	Short: "Output a graph representation of the manifest",
	Long: `Generate a graph showing dependencies, recipe steps, and variable flow.
	
Output formats:
  - dot:  GraphViz DOT format
  - mermaid: Mermaid diagram
  - text: Simple text representation`,
	Args: cobra.ExactArgs(1),
	RunE: runManifestGraph,
}

// manifestNewCmd creates a new manifest from template
var manifestNewCmd = &cobra.Command{
	Use:   "new [name]",
	Short: "Create a new manifest file from a template",
	Long: `Create a new manifest file from a predefined template.

Available templates:
  - minimal:    Basic manifest with essential fields only
  - cmake:      Manifest for CMake-based projects
  - autotools:  Manifest for Autotools-based projects  
  - git:        Manifest for Git-based source projects
  - tarball:    Manifest for tarball-based projects`,
	Args: cobra.ExactArgs(1),
	RunE: runManifestNew,
}

// manifestInitCmd initializes manifest interactively
var manifestInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a manifest in the current directory using interactive prompts",
	Long:  `Creates a new manifest file by asking questions interactively.`,
	RunE:  runManifestInit,
}

// manifestDepsCmd shows dependencies
var manifestDepsCmd = &cobra.Command{
	Use:   "deps [manifest]",
	Short: "Resolve and display dependency graph from a manifest",
	Long: `Display all dependencies (build and runtime) from a manifest.
	
Can show dependencies filtered by build options.`,
	Args: cobra.ExactArgs(1),
	RunE: runManifestDeps,
}

// manifestOptionsCmd lists options
var manifestOptionsCmd = &cobra.Command{
	Use:   "options [manifest]",
	Short: "List all extensions/options and their conditional effects",
	Long:  `Display all available build options and show which dependencies they enable.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runManifestOptions,
}

// manifestEnvCmd shows environment variables
var manifestEnvCmd = &cobra.Command{
	Use:   "env [manifest]",
	Short: "Show all environment variables produced or consumed by the manifest",
	Long: `Display environment variables used by the manifest.
	
Scopes:
  - build:   Variables needed during build
  - runtime: Variables needed at runtime
  - self:    Variables exported by this package
  - all:     All environment variables (default)`,
	Args: cobra.ExactArgs(1),
	RunE: runManifestEnv,
}

// manifestSourceCmd inspects source metadata
var manifestSourceCmd = &cobra.Command{
	Use:   "source [manifest]",
	Short: "Inspect and verify source metadata (URLs, checksums, SCM refs)",
	Long:  `Display detailed information about the package source.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runManifestSource,
}

// manifestRecipeCmd shows recipe steps
var manifestRecipeCmd = &cobra.Command{
	Use:   "recipe [manifest]",
	Short: "Display the expanded recipe steps after variable interpolation",
	Long: `Show all recipe steps with variables expanded.
	
Use --options to see how different build options affect the recipe.`,
	Args: cobra.ExactArgs(1),
	RunE: runManifestRecipe,
}

// manifestFormatCmd formats a manifest
var manifestFormatCmd = &cobra.Command{
	Use:   "format [manifest]",
	Short: "Auto-format and normalize a manifest file",
	Long: `Automatically format a manifest file with proper indentation and ordering.
	
This command normalizes the structure and ensures consistent formatting.`,
	Args: cobra.ExactArgs(1),
	RunE: runManifestFormat,
}

// manifestDiffCmd compares manifests
var manifestDiffCmd = &cobra.Command{
	Use:   "diff [manifest1] [manifest2]",
	Short: "Compare two manifests and show structural and semantic differences",
	Long:  `Compare two manifest files and display their differences.`,
	Args:  cobra.ExactArgs(2),
	RunE:  runManifestDiff,
}

// manifestExportCmd exports manifest to different formats
var manifestExportCmd = &cobra.Command{
	Use:   "export [manifest]",
	Short: "Convert a manifest into another format (JSON, lockfile, build spec)",
	Long: `Export a manifest to different formats:
  - json:     Export as JSON
  - yaml:     Export as YAML (default)
  - lockfile: Generate a lockfile
  - report:   Generate a detailed report`,
	Args: cobra.ExactArgs(1),
	RunE: runManifestExport,
}

// manifestFlattenCmd flattens a manifest
var manifestFlattenCmd = &cobra.Command{
	Use:   "flatten [manifest]",
	Short: "Produce a fully-expanded manifest with all conditionals resolved",
	Long: `Flatten a manifest by resolving all conditional dependencies and variables.
	
This produces a manifest with all options applied and conditionals evaluated.`,
	Args: cobra.ExactArgs(1),
	RunE: runManifestFlatten,
}

// manifestLintCmd lints a manifest
var manifestLintCmd = &cobra.Command{
	Use:   "lint [manifest]",
	Short: "Enforce style and best-practice rules on a manifest",
	Long: `Check a manifest for style issues and best practices.
	
This includes:
  - Naming conventions
  - Field ordering
  - Recommended fields
  - Common pitfalls`,
	Args: cobra.ExactArgs(1),
	RunE: runManifestLint,
}

// manifestCheckCmd performs deep consistency checks
var manifestCheckCmd = &cobra.Command{
	Use:   "check [manifest]",
	Short: "Run deep consistency checks",
	Long: `Perform comprehensive consistency checks including:
  - Dependency cycles
  - Incompatible options
  - Target mismatches
  - Version conflicts`,
	Args: cobra.ExactArgs(1),
	RunE: runManifestCheck,
}

// manifestValidateCmd validates a manifest
var manifestValidateCmd = &cobra.Command{
	Use:   "validate [manifest]",
	Short: "Validate a manifest against the schema and report errors",
	Long: `Validates a manifest file for structural and semantic correctness.
	
Checks include:
  - Required fields presence
  - Correct data types and formats
  - Dependency resolution
  - Version constraints
  - Recipe step validity
  - Environment variable definitions
  - Overall schema adherence`,
	Args: cobra.ExactArgs(1),
	RunE: runManifestValidate,
}

// manifestShowCmd displays a manifest
var manifestShowCmd = &cobra.Command{
	Use:   "show [manifest]",
	Short: "Print a normalized, fully-resolved view of a manifest",
	Long:  `Displays the manifest in a human-readable format with all fields resolved.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runManifestShow,
}

// manifestInspectCmd inspects specific fields
var manifestInspectCmd = &cobra.Command{
	Use:   "inspect [manifest] [field]",
	Short: "Query specific fields or paths inside a manifest",
	Long: `Inspect and display specific fields from a manifest.
	
Examples:
  hepsw manifest inspect myapp.yaml name
  hepsw manifest inspect myapp.yaml source.url
  hepsw manifest inspect myapp.yaml specifications.build.dependencies`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runManifestInspect,
}

// manifestWalkCmd walks through recipe steps
var manifestWalkCmd = &cobra.Command{
	Use:   "walk [manifest]",
	Short: "Simulate the build process without executing steps",
	Long: `Walk through the recipe steps and display what would be executed.
	
This command simulates the build process, showing which steps would run
based on the provided options and variables.`,
	Args: cobra.ExactArgs(1),
	RunE: runManifestWalk,
}

func init() {
	ManifestCmd.AddCommand(manifestFetchCmd)
	ManifestCmd.AddCommand(manifestStripCmd)
	ManifestCmd.AddCommand(manifestExplainCmd)
	ManifestCmd.AddCommand(manifestGraphCmd)

	manifestFetchCmd.Flags().StringVarP(&fetchDestination, "dest", "d", ".",
		"Destination directory for downloaded manifest")

	manifestStripCmd.Flags().BoolVar(&stripMinimal, "minimal", false,
		"Create absolute minimal manifest (only required fields)")

	ManifestCmd.AddCommand(manifestNewCmd)
	ManifestCmd.AddCommand(manifestInitCmd)

	manifestNewCmd.Flags().StringVarP(&manifestNewTemplate, "template", "t", "minimal",
		"Template to use (minimal, cmake, autotools, git, tarball)")
	manifestNewCmd.Flags().StringVarP(&manifestNewOutput, "output", "o", "",
		"Output file path (default: <name>.yaml)")

	ManifestCmd.AddCommand(manifestDepsCmd)
	ManifestCmd.AddCommand(manifestOptionsCmd)
	ManifestCmd.AddCommand(manifestEnvCmd)
	ManifestCmd.AddCommand(manifestSourceCmd)
	ManifestCmd.AddCommand(manifestRecipeCmd)

	manifestDepsCmd.Flags().StringVarP(&depsFormat, "format", "f", "tree",
		"Output format (tree, list, json)")
	manifestDepsCmd.Flags().BoolVar(&depsShowOptions, "show-options", false,
		"Show which options enable each dependency")

	manifestEnvCmd.Flags().StringVarP(&envScope, "scope", "s", "all",
		"Environment scope (build, runtime, self, all)")

	manifestRecipeCmd.Flags().StringSliceVarP(&walkOptions, "options", "o", []string{},
		"Build options to enable")

	ManifestCmd.AddCommand(manifestFormatCmd)
	ManifestCmd.AddCommand(manifestDiffCmd)
	ManifestCmd.AddCommand(manifestExportCmd)
	ManifestCmd.AddCommand(manifestFlattenCmd)
	ManifestCmd.AddCommand(manifestLintCmd)
	ManifestCmd.AddCommand(manifestCheckCmd)

	manifestFormatCmd.Flags().IntVar(&formatIndent, "indent", 2,
		"Number of spaces for indentation")

	manifestDiffCmd.Flags().StringVarP(&diffFormat, "format", "f", "text",
		"Output format (text, json)")

	manifestExportCmd.Flags().StringVarP(&exportFormat, "format", "f", "yaml",
		"Export format (json, yaml, lockfile, report)")
	manifestExportCmd.Flags().StringVarP(&exportOutput, "output", "o", "",
		"Output file (default: stdout)")

	manifestFlattenCmd.Flags().StringSliceVarP(&flattenOptions, "options", "o", []string{},
		"Options to apply when flattening")
	manifestFlattenCmd.Flags().StringVar(&flattenOutput, "output", "",
		"Output file (default: stdout)")

	ManifestCmd.AddCommand(manifestValidateCmd)
	ManifestCmd.AddCommand(manifestShowCmd)
	ManifestCmd.AddCommand(manifestInspectCmd)
	ManifestCmd.AddCommand(manifestWalkCmd)

	manifestValidateCmd.Flags().BoolVarP(&validateVerbose, "verbose", "v", false,
		"Show detailed validation information")

	manifestShowCmd.Flags().StringVarP(&showFormat, "format", "f", "text",
		"Output format (text, yaml, json)")

	manifestWalkCmd.Flags().StringSliceVarP(&walkOptions, "options", "o", []string{},
		"Build options to enable (e.g., with-ssl,with-gui)")
	manifestWalkCmd.Flags().StringToStringVarP(&walkVariables, "var", "V", map[string]string{},
		"Set variables (e.g., NCORES=8)")
}
