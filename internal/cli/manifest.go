package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thisismeamir/hepsw/internal/manifest"
	"github.com/thisismeamir/hepsw/internal/manifest/loader"
	"github.com/thisismeamir/hepsw/internal/manifest/reporters"
	"gopkg.in/yaml.v3"
)

// manifestCmd represents the manifest command
var manifestCmd = &cobra.Command{
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
	manifestCmd.AddCommand(manifestFetchCmd)
	manifestCmd.AddCommand(manifestStripCmd)
	manifestCmd.AddCommand(manifestExplainCmd)
	manifestCmd.AddCommand(manifestGraphCmd)

	manifestFetchCmd.Flags().StringVarP(&fetchDestination, "dest", "d", ".",
		"Destination directory for downloaded manifest")

	manifestStripCmd.Flags().BoolVar(&stripMinimal, "minimal", false,
		"Create absolute minimal manifest (only required fields)")

	manifestCmd.AddCommand(manifestNewCmd)
	manifestCmd.AddCommand(manifestInitCmd)

	manifestNewCmd.Flags().StringVarP(&manifestNewTemplate, "template", "t", "minimal",
		"Template to use (minimal, cmake, autotools, git, tarball)")
	manifestNewCmd.Flags().StringVarP(&manifestNewOutput, "output", "o", "",
		"Output file path (default: <name>.yaml)")

	manifestCmd.AddCommand(manifestDepsCmd)
	manifestCmd.AddCommand(manifestOptionsCmd)
	manifestCmd.AddCommand(manifestEnvCmd)
	manifestCmd.AddCommand(manifestSourceCmd)
	manifestCmd.AddCommand(manifestRecipeCmd)

	manifestDepsCmd.Flags().StringVarP(&depsFormat, "format", "f", "tree",
		"Output format (tree, list, json)")
	manifestDepsCmd.Flags().BoolVar(&depsShowOptions, "show-options", false,
		"Show which options enable each dependency")

	manifestEnvCmd.Flags().StringVarP(&envScope, "scope", "s", "all",
		"Environment scope (build, runtime, self, all)")

	manifestRecipeCmd.Flags().StringSliceVarP(&walkOptions, "options", "o", []string{},
		"Build options to enable")

	manifestCmd.AddCommand(manifestFormatCmd)
	manifestCmd.AddCommand(manifestDiffCmd)
	manifestCmd.AddCommand(manifestExportCmd)
	manifestCmd.AddCommand(manifestFlattenCmd)
	manifestCmd.AddCommand(manifestLintCmd)
	manifestCmd.AddCommand(manifestCheckCmd)

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

	manifestCmd.AddCommand(manifestValidateCmd)
	manifestCmd.AddCommand(manifestShowCmd)
	manifestCmd.AddCommand(manifestInspectCmd)
	manifestCmd.AddCommand(manifestWalkCmd)

	manifestValidateCmd.Flags().BoolVarP(&validateVerbose, "verbose", "v", false,
		"Show detailed validation information")

	manifestShowCmd.Flags().StringVarP(&showFormat, "format", "f", "text",
		"Output format (text, yaml, json)")

	manifestWalkCmd.Flags().StringSliceVarP(&walkOptions, "options", "o", []string{},
		"Build options to enable (e.g., with-ssl,with-gui)")
	manifestWalkCmd.Flags().StringToStringVarP(&walkVariables, "var", "V", map[string]string{},
		"Set variables (e.g., NCORES=8)")
}

func runManifestFetch(cmd *cobra.Command, args []string) error {
	reference := args[0]

	fmt.Printf("Fetching manifest: %s\n", reference)

	// Load from index
	m, err := loader.LoadManifestFromIndex(reference)
	if err != nil {
		return fmt.Errorf("failed to fetch manifest: %w", err)
	}

	// Determine output filename
	filename := fmt.Sprintf("%s-%s.yaml", m.Name, m.Version)
	outputPath := filepath.Join(fetchDestination, filename)

	// Save manifest
	if err := loader.SaveManifest(m, outputPath); err != nil {
		return fmt.Errorf("failed to save manifest: %w", err)
	}

	fmt.Printf("Downloaded to: %s\n", outputPath)
	fmt.Printf("Package: %s@%s\n", m.Name, m.Version)

	return nil
}

func runManifestStrip(cmd *cobra.Command, args []string) error {
	manifestSource := args[0]

	// Load manifest
	m, err := loader.LoadManifest(manifestSource)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	// Create stripped version
	stripped := stripManifest(m, stripMinimal)

	// Generate output filename
	ext := filepath.Ext(manifestSource)
	base := strings.TrimSuffix(manifestSource, ext)
	outputPath := base + ".minimal" + ext

	// Save stripped manifest
	if err := loader.SaveManifest(stripped, outputPath); err != nil {
		return fmt.Errorf("failed to save stripped manifest: %w", err)
	}

	fmt.Printf("Stripped manifest saved to: %s\n", outputPath)

	return nil
}

func runManifestExplain(cmd *cobra.Command, args []string) error {
	manifestSource := args[0]
	target := args[1]

	// Load manifest
	m, err := loader.LoadManifest(manifestSource)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	accessor := manifest.NewManifestAccessor(m)

	fmt.Printf("Explaining: %s in %s@%s\n", target, m.Name, m.Version)
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	// Check if it's a dependency
	if dep, found := accessor.GetDependency(target); found {
		fmt.Printf("Dependency: %s\n", dep.Name)
		fmt.Printf("Version: %s\n", dep.Version)

		if len(dep.ForOptions) > 0 {
			fmt.Printf("\nThis dependency is CONDITIONAL.\n")
			fmt.Printf("It will be included when these options are enabled:\n")
			for _, opt := range dep.ForOptions {
				fmt.Printf("  - %s\n", opt)
			}
		} else {
			fmt.Printf("\nThis dependency is ALWAYS included.\n")
		}

		if len(dep.WithOptions) > 0 {
			fmt.Printf("\nThis dependency REQUIRES these options:\n")
			for _, opt := range dep.WithOptions {
				fmt.Printf("  - %s\n", opt)
			}
		}

		if dep.IsOptional {
			fmt.Printf("\nThis dependency is OPTIONAL.\n")
		}

		return nil
	}

	// Check if it's an option
	if accessor.HasOption(target) {
		fmt.Printf("Option: %s\n\n", target)
		fmt.Println("This option affects the following dependencies:")

		affected := false
		for _, dep := range accessor.AllDependencies() {
			for _, opt := range dep.ForOptions {
				if opt == target {
					fmt.Printf("  - %s %s (enables)\n", dep.Name, dep.Version)
					affected = true
				}
			}
		}

		if !affected {
			fmt.Println("  (no dependencies are affected)")
		}

		return nil
	}

	return fmt.Errorf("target not found: %s", target)
}

func runManifestGraph(cmd *cobra.Command, args []string) error {
	manifestSource := args[0]

	// Load manifest
	m, err := loader.LoadManifest(manifestSource)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	// Generate graph
	graph := generateDependencyGraph(m)

	fmt.Println(graph)

	return nil
}

// Helper functions

func stripManifest(m *manifest.Manifest, minimal bool) *manifest.Manifest {
	stripped := &manifest.Manifest{
		Name:        m.Name,
		Version:     m.Version,
		Description: m.Description,
		Source:      m.Source,
		Recipe:      m.Recipe,
	}

	if !minimal {
		// Keep specifications if not minimal
		stripped.Specifications = m.Specifications
		stripped.Metadata = m.Metadata
	} else {
		// Minimal: only keep essential build specs
		stripped.Specifications.Build.Toolchain = m.Specifications.Build.Toolchain
	}

	return stripped
}

func generateDependencyGraph(m *manifest.Manifest) string {
	var sb strings.Builder

	sb.WriteString("# Dependency Graph\n\n")
	sb.WriteString(fmt.Sprintf("Package: %s@%s\n", m.Name, m.Version))
	sb.WriteString("\n")

	// Build dependencies
	sb.WriteString("Build Dependencies:\n")
	if len(m.Specifications.Build.Dependencies) > 0 {
		for _, dep := range m.Specifications.Build.Dependencies {
			sb.WriteString(fmt.Sprintf("  %s --> %s %s", m.Name, dep.Name, dep.Version))
			if len(dep.ForOptions) > 0 {
				sb.WriteString(fmt.Sprintf(" [when: %s]", strings.Join(dep.ForOptions, ", ")))
			}
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString("  (none)\n")
	}

	sb.WriteString("\n")

	// Runtime dependencies
	sb.WriteString("Runtime Dependencies:\n")
	if len(m.Specifications.Runtime.Dependencies) > 0 {
		for _, dep := range m.Specifications.Runtime.Dependencies {
			sb.WriteString(fmt.Sprintf("  %s ==> %s %s", m.Name, dep.Name, dep.Version))
			if len(dep.ForOptions) > 0 {
				sb.WriteString(fmt.Sprintf(" [when: %s]", strings.Join(dep.ForOptions, ", ")))
			}
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString("  (none)\n")
	}

	sb.WriteString("\n")

	// Recipe flow
	sb.WriteString("Recipe Flow:\n")
	sb.WriteString("  Configuration --> Build --> Install --> Use\n")

	return sb.String()
}

func runManifestNew(cmd *cobra.Command, args []string) error {
	name := args[0]

	// Get template
	tmpl, err := getTemplate(manifestNewTemplate)
	if err != nil {
		return err
	}

	// Customize template with name
	tmpl.Name = name

	// Determine output path
	outputPath := manifestNewOutput
	if outputPath == "" {
		outputPath = fmt.Sprintf("%s.yaml", name)
	}

	// Save manifest
	if err := loader.SaveManifest(tmpl, outputPath); err != nil {
		return fmt.Errorf("failed to save manifest: %w", err)
	}

	fmt.Printf("Created manifest: %s\n", outputPath)
	fmt.Printf("Template: %s\n", manifestNewTemplate)
	return nil
}

func runManifestInit(cmd *cobra.Command, args []string) error {
	fmt.Println("Initializing new HepSW manifest...")
	fmt.Println()

	// Interactive prompts
	name := prompt("Package name", getCurrentDirName())
	version := prompt("Version", "0.1.0")
	description := prompt("Description", "")

	sourceType := promptChoice("Source type", []string{"git", "tarball", "svn", "local"}, "git")
	sourceURL := prompt("Source URL", "")

	var sourceTag string
	if sourceType == "git" {
		sourceTag = prompt("Git tag/branch", "main")
	}

	buildSystem := promptChoice("Build system", []string{"cmake", "autotools", "make", "custom"}, "cmake")

	// Create manifest
	m := &manifest.Manifest{
		Name:        name,
		Version:     version,
		Description: description,
		Source: manifest.SourceSpec{
			Type: sourceType,
			Url:  sourceURL,
			Tag:  sourceTag,
		},
		Specifications: manifest.Specifications{
			Build: manifest.BuildSpecification{
				Toolchain: []manifest.Tool{},
			},
		},
		Recipe: createRecipeForBuildSystem(buildSystem),
	}

	// Save manifest
	outputPath := fmt.Sprintf("%s.yaml", name)
	if err := loader.SaveManifest(m, outputPath); err != nil {
		return fmt.Errorf("failed to save manifest: %w", err)
	}

	fmt.Printf("\nCreated manifest: %s\n", outputPath)
	return nil
}

func getTemplate(templateName string) (*manifest.Manifest, error) {
	switch templateName {
	case "minimal":
		return getMinimalTemplate(), nil
	case "cmake":
		return getCMakeTemplate(), nil
	case "autotools":
		return getAutotoolsTemplate(), nil
	case "git":
		return getGitTemplate(), nil
	case "tarball":
		return getTarballTemplate(), nil
	default:
		return nil, fmt.Errorf("unknown template: %s", templateName)
	}
}

func getMinimalTemplate() *manifest.Manifest {
	return &manifest.Manifest{
		Name:        "example",
		Version:     "1.0.0",
		Description: "A simple example package",
		Source: manifest.SourceSpec{
			Type: "tarball",
			Url:  "https://example.com/example-1.0.0.tar.gz",
		},
		Recipe: manifest.Recipe{
			Configuration: []manifest.RecipeStep{
				{
					Name:    "Configure",
					Command: "./configure --prefix=${INSTALL_PREFIX}",
				},
			},
			Build: []manifest.RecipeStep{
				{
					Name:    "Build",
					Command: "make -j${NCORES}",
				},
			},
			Install: []manifest.RecipeStep{
				{
					Name:    "Install",
					Command: "make install",
				},
			},
		},
	}
}

func getCMakeTemplate() *manifest.Manifest {
	return &manifest.Manifest{
		Name:        "example",
		Version:     "1.0.0",
		Description: "A CMake-based example package",
		Source: manifest.SourceSpec{
			Type: "git",
			Url:  "https://github.com/example/example.git",
			Tag:  "v1.0.0",
		},
		Specifications: manifest.Specifications{
			Build: manifest.BuildSpecification{
				Toolchain: []manifest.Tool{
					{Name: "cmake", Version: ">=3.15"},
					{Name: "gcc", Version: ">=9.0"},
				},
			},
		},
		Recipe: manifest.Recipe{
			Configuration: []manifest.RecipeStep{
				{
					Name:    "Create build directory",
					Command: "mkdir -p build && cd build",
				},
				{
					Name:    "Configure with CMake",
					Command: "cmake .. -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX=${INSTALL_PREFIX}",
				},
			},
			Build: []manifest.RecipeStep{
				{
					Name:    "Build",
					Command: "cmake --build . -j${NCORES}",
				},
			},
			Install: []manifest.RecipeStep{
				{
					Name:    "Install",
					Command: "cmake --install .",
				},
			},
		},
	}
}

func getAutotoolsTemplate() *manifest.Manifest {
	return &manifest.Manifest{
		Name:        "example",
		Version:     "1.0.0",
		Description: "An Autotools-based example package",
		Source: manifest.SourceSpec{
			Type: "tarball",
			Url:  "https://example.com/example-1.0.0.tar.gz",
		},
		Specifications: manifest.Specifications{
			Build: manifest.BuildSpecification{
				Toolchain: []manifest.Tool{
					{Name: "gcc", Version: ">=7.0"},
					{Name: "autoconf", Version: ">=2.69"},
					{Name: "automake", Version: ">=1.16"},
				},
			},
		},
		Recipe: manifest.Recipe{
			Configuration: []manifest.RecipeStep{
				{
					Name:    "Configure",
					Command: "./configure --prefix=${INSTALL_PREFIX}",
				},
			},
			Build: []manifest.RecipeStep{
				{
					Name:    "Build",
					Command: "make -j${NCORES}",
				},
			},
			Install: []manifest.RecipeStep{
				{
					Name:    "Install",
					Command: "make install",
				},
			},
		},
	}
}

func getGitTemplate() *manifest.Manifest {
	m := getMinimalTemplate()
	m.Source = manifest.SourceSpec{
		Type: "git",
		Url:  "https://github.com/example/example.git",
		Tag:  "main",
	}
	return m
}

func getTarballTemplate() *manifest.Manifest {
	return getMinimalTemplate()
}

func createRecipeForBuildSystem(buildSystem string) manifest.Recipe {
	switch buildSystem {
	case "cmake":
		return getCMakeTemplate().Recipe
	case "autotools":
		return getAutotoolsTemplate().Recipe
	default:
		return getMinimalTemplate().Recipe
	}
}

// Helper functions

func prompt(question, defaultValue string) string {
	if defaultValue != "" {
		fmt.Printf("%s [%s]: ", question, defaultValue)
	} else {
		fmt.Printf("%s: ", question)
	}

	var answer string
	fmt.Scanln(&answer)

	if answer == "" && defaultValue != "" {
		return defaultValue
	}
	return answer
}

func promptChoice(question string, choices []string, defaultValue string) string {
	fmt.Printf("%s (%s) [%s]: ", question, strings.Join(choices, "/"), defaultValue)

	var answer string
	fmt.Scanln(&answer)

	if answer == "" {
		return defaultValue
	}

	// Validate choice
	for _, choice := range choices {
		if answer == choice {
			return answer
		}
	}

	return defaultValue
}

func getCurrentDirName() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "example"
	}
	return filepath.Base(cwd)
}

func runManifestDeps(cmd *cobra.Command, args []string) error {
	manifestSource := args[0]

	// Load manifest
	m, err := loader.LoadManifest(manifestSource)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	accessor := manifest.NewManifestAccessor(m)

	fmt.Printf("Dependencies for %s@%s\n", accessor.Name(), accessor.Version())
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	// Build dependencies
	buildDeps := accessor.BuildDependencies()
	if len(buildDeps) > 0 {
		fmt.Println("Build Dependencies:")
		for _, dep := range buildDeps {
			fmt.Printf("  ├─ %s %s", dep.Name, dep.Version)
			if depsShowOptions && len(dep.ForOptions) > 0 {
				fmt.Printf(" [for: %s]", strings.Join(dep.ForOptions, ", "))
			}
			if dep.IsOptional {
				fmt.Print(" (optional)")
			}
			fmt.Println()
		}
		fmt.Println()
	}

	// Runtime dependencies
	runtimeDeps := accessor.RuntimeDependencies()
	if len(runtimeDeps) > 0 {
		fmt.Println("Runtime Dependencies:")
		for _, dep := range runtimeDeps {
			fmt.Printf("  ├─ %s %s", dep.Name, dep.Version)
			if depsShowOptions && len(dep.ForOptions) > 0 {
				fmt.Printf(" [for: %s]", strings.Join(dep.ForOptions, ", "))
			}
			fmt.Println()
		}
		fmt.Println()
	}

	if len(buildDeps) == 0 && len(runtimeDeps) == 0 {
		fmt.Println("No dependencies found.")
	}

	return nil
}

func runManifestOptions(cmd *cobra.Command, args []string) error {
	manifestSource := args[0]

	// Load manifest
	m, err := loader.LoadManifest(manifestSource)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	accessor := manifest.NewManifestAccessor(m)

	fmt.Printf("Options for %s@%s\n", accessor.Name(), accessor.Version())
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	options := accessor.Options()
	if len(options) == 0 {
		fmt.Println("No build options defined.")
		return nil
	}

	// Show each option and its effects
	for _, opt := range options {
		fmt.Printf("Option: %s\n", opt)

		// Find dependencies enabled by this option
		enabledDeps := make([]string, 0)
		for _, dep := range accessor.AllDependencies() {
			for _, reqOpt := range dep.ForOptions {
				if reqOpt == opt {
					enabledDeps = append(enabledDeps, fmt.Sprintf("%s %s", dep.Name, dep.Version))
					break
				}
			}
		}

		if len(enabledDeps) > 0 {
			fmt.Println("  Enables dependencies:")
			for _, dep := range enabledDeps {
				fmt.Printf("    - %s\n", dep)
			}
		} else {
			fmt.Println("  No dependencies affected")
		}
		fmt.Println()
	}

	return nil
}

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

func runManifestSource(cmd *cobra.Command, args []string) error {
	manifestSource := args[0]

	// Load manifest
	m, err := loader.LoadManifest(manifestSource)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	accessor := manifest.NewManifestAccessor(m)

	fmt.Printf("Source Information for %s@%s\n", accessor.Name(), accessor.Version())
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	source := accessor.Source()

	fmt.Printf("Type:     %s\n", source.Type)
	fmt.Printf("URL:      %s\n", source.Url)

	if source.Tag != "" {
		fmt.Printf("Tag:      %s\n", source.Tag)
	}

	if source.Checksum != "" {
		fmt.Printf("Checksum: %s\n", source.Checksum)

		// Parse checksum
		parts := strings.Split(source.Checksum, ":")
		if len(parts) == 2 {
			fmt.Printf("  Algorithm: %s\n", parts[0])
			fmt.Printf("  Hash:      %s\n", parts[1])
		}
	} else {
		fmt.Println("Checksum: (not provided)")
		fmt.Println("  Warning: Checksums are recommended for reproducibility")
	}

	return nil
}

func runManifestRecipe(cmd *cobra.Command, args []string) error {
	manifestSource := args[0]

	// Load manifest
	m, err := loader.LoadManifest(manifestSource)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	accessor := manifest.NewManifestAccessor(m)

	fmt.Printf("Recipe for %s@%s\n", accessor.Name(), accessor.Version())
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	// Print each phase
	phases := []struct {
		name  string
		steps []manifest.RecipeStep
	}{
		{"Configuration", accessor.ConfigurationSteps()},
		{"Build", accessor.BuildSteps()},
		{"Install", accessor.InstallSteps()},
		{"Use", accessor.UseSteps()},
	}

	for _, phase := range phases {
		if len(phase.steps) == 0 {
			continue
		}

		fmt.Printf("=== %s Phase ===\n\n", phase.name)

		for i, step := range phase.steps {
			fmt.Printf("%d. %s\n", i+1, step.Name)

			if step.Command != "" {
				fmt.Printf("   Command: %s\n", step.Command)
			}

			if step.Script != "" {
				fmt.Printf("   Script: %s\n", step.Script)
				if len(step.Args) > 0 {
					fmt.Printf("   Args: %v\n", step.Args)
				}
			}

			if step.WorkingDir != "" {
				fmt.Printf("   Working Directory: %s\n", step.WorkingDir)
			}

			if step.If != "" {
				fmt.Printf("   Condition: %s\n", step.If)
			}

			if step.Set != nil {
				fmt.Println("   Sets variables:")
				for k, v := range step.Set {
					fmt.Printf("     %s = %s\n", k, v)
				}
			}

			fmt.Println()
		}
	}

	return nil
}

func runManifestFormat(cmd *cobra.Command, args []string) error {
	manifestPath := args[0]

	// Load manifest
	m, err := loader.LoadManifest(manifestPath)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	// Format and save
	data, err := yaml.Marshal(m)
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}

	// Write back to file
	if err := os.WriteFile(manifestPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write formatted manifest: %w", err)
	}

	fmt.Printf("Formatted manifest: %s\n", manifestPath)
	return nil
}

func runManifestDiff(cmd *cobra.Command, args []string) error {
	manifest1Path := args[0]
	manifest2Path := args[1]

	// Load both manifests
	m1, err := loader.LoadManifest(manifest1Path)
	if err != nil {
		return fmt.Errorf("failed to load first manifest: %w", err)
	}

	m2, err := loader.LoadManifest(manifest2Path)
	if err != nil {
		return fmt.Errorf("failed to load second manifest: %w", err)
	}

	// Compare manifests
	fmt.Printf("Comparing:\n  %s (%s@%s)\n  %s (%s@%s)\n\n",
		manifest1Path, m1.Name, m1.Version,
		manifest2Path, m2.Name, m2.Version)

	differences := compareManifests(m1, m2)

	if len(differences) == 0 {
		fmt.Println("No differences found.")
		return nil
	}

	fmt.Println("Differences:")
	for _, diff := range differences {
		fmt.Printf("  %s\n", diff)
	}

	return nil
}

func runManifestExport(cmd *cobra.Command, args []string) error {
	manifestSource := args[0]

	// Load manifest
	m, err := loader.LoadManifest(manifestSource)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	var output string

	switch exportFormat {
	case "json":
		data, err := json.MarshalIndent(m, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		output = string(data)

	case "yaml":
		data, err := yaml.Marshal(m)
		if err != nil {
			return fmt.Errorf("failed to marshal YAML: %w", err)
		}
		output = string(data)

	case "lockfile":
		output = generateLockfile(m)

	case "report":
		output, err = reporters.GenerateReport(m, reporters.FormatText)
		if err != nil {
			return fmt.Errorf("failed to generate report: %w", err)
		}

	default:
		return fmt.Errorf("unsupported export format: %s", exportFormat)
	}

	// Write output
	if exportOutput != "" {
		if err := os.WriteFile(exportOutput, []byte(output), 0644); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
		fmt.Printf("Exported to: %s\n", exportOutput)
	} else {
		fmt.Print(output)
	}

	return nil
}

func runManifestFlatten(cmd *cobra.Command, args []string) error {
	manifestSource := args[0]

	// Load manifest
	m, err := loader.LoadManifest(manifestSource)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	// Create flattened version
	flattened := flattenManifest(m, flattenOptions)

	// Output
	data, err := yaml.Marshal(flattened)
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}

	if flattenOutput != "" {
		if err := os.WriteFile(flattenOutput, data, 0644); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
		fmt.Printf("Flattened manifest written to: %s\n", flattenOutput)
	} else {
		fmt.Print(string(data))
	}

	return nil
}

func runManifestLint(cmd *cobra.Command, args []string) error {
	manifestSource := args[0]

	// Load manifest
	m, err := loader.LoadManifest(manifestSource)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	fmt.Printf("Linting manifest: %s@%s\n", m.Name, m.Version)
	fmt.Println()

	// Run validation (includes many lint checks)
	result := manifest.ValidateManifest(m)

	// Additional lint-specific checks
	lintIssues := performLintChecks(m)

	// Combine results
	allIssues := append(result.Warnings, lintIssues...)

	if len(allIssues) > 0 {
		fmt.Println("Lint Issues:")
		for _, issue := range allIssues {
			fmt.Printf("  ⚠ %s\n", issue)
		}
		fmt.Printf("\nFound %d issue(s)\n", len(allIssues))
	} else {
		fmt.Println("✓ No lint issues found")
	}

	return nil
}

func runManifestCheck(cmd *cobra.Command, args []string) error {
	manifestSource := args[0]

	// Load manifest
	m, err := loader.LoadManifest(manifestSource)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	fmt.Printf("Checking manifest: %s@%s\n", m.Name, m.Version)
	fmt.Println()

	// Run comprehensive validation
	result := manifest.ValidateManifest(m)

	// Additional deep checks
	deepIssues := performDeepChecks(m)

	// Print all issues
	hasIssues := false

	if len(result.Errors) > 0 {
		fmt.Println("ERRORS:")
		for _, err := range result.Errors {
			fmt.Printf("  ✗ %s\n", err)
		}
		fmt.Println()
		hasIssues = true
	}

	if len(result.Warnings) > 0 {
		fmt.Println("WARNINGS:")
		for _, warn := range result.Warnings {
			fmt.Printf("  ⚠ %s\n", warn)
		}
		fmt.Println()
		hasIssues = true
	}

	if len(deepIssues) > 0 {
		fmt.Println("CONSISTENCY ISSUES:")
		for _, issue := range deepIssues {
			fmt.Printf("  ⚠ %s\n", issue)
		}
		fmt.Println()
		hasIssues = true
	}

	if !hasIssues {
		fmt.Println("✓ All checks passed")
	}

	return nil
}

// Helper functions

func compareManifests(m1, m2 *manifest.Manifest) []string {
	diffs := make([]string, 0)

	if m1.Name != m2.Name {
		diffs = append(diffs, fmt.Sprintf("Name: %s → %s", m1.Name, m2.Name))
	}

	if m1.Version != m2.Version {
		diffs = append(diffs, fmt.Sprintf("Version: %s → %s", m1.Version, m2.Version))
	}

	if m1.Description != m2.Description {
		diffs = append(diffs, "Description changed")
	}

	if m1.Source.Type != m2.Source.Type {
		diffs = append(diffs, fmt.Sprintf("Source type: %s → %s", m1.Source.Type, m2.Source.Type))
	}

	if m1.Source.Url != m2.Source.Url {
		diffs = append(diffs, "Source URL changed")
	}

	// Compare dependencies
	if len(m1.Specifications.Build.Dependencies) != len(m2.Specifications.Build.Dependencies) {
		diffs = append(diffs, fmt.Sprintf("Build dependencies count: %d → %d",
			len(m1.Specifications.Build.Dependencies),
			len(m2.Specifications.Build.Dependencies)))
	}

	// Compare recipe steps
	if len(m1.Recipe.Build) != len(m2.Recipe.Build) {
		diffs = append(diffs, fmt.Sprintf("Build steps count: %d → %d",
			len(m1.Recipe.Build), len(m2.Recipe.Build)))
	}

	return diffs
}

func generateLockfile(m *manifest.Manifest) string {
	var sb strings.Builder

	sb.WriteString("# HepSW Lockfile\n")
	sb.WriteString(fmt.Sprintf("# Generated for: %s@%s\n\n", m.Name, m.Version))

	sb.WriteString("[[package]]\n")
	sb.WriteString(fmt.Sprintf("name = \"%s\"\n", m.Name))
	sb.WriteString(fmt.Sprintf("version = \"%s\"\n", m.Version))
	sb.WriteString(fmt.Sprintf("source = \"%s\"\n", m.Source.Url))
	if m.Source.Checksum != "" {
		sb.WriteString(fmt.Sprintf("checksum = \"%s\"\n", m.Source.Checksum))
	}
	sb.WriteString("\n")

	// Add dependencies
	for _, dep := range m.Specifications.Build.Dependencies {
		sb.WriteString("[[dependency]]\n")
		sb.WriteString(fmt.Sprintf("name = \"%s\"\n", dep.Name))
		sb.WriteString(fmt.Sprintf("version = \"%s\"\n", dep.Version))
		sb.WriteString("type = \"build\"\n\n")
	}

	for _, dep := range m.Specifications.Runtime.Dependencies {
		sb.WriteString("[[dependency]]\n")
		sb.WriteString(fmt.Sprintf("name = \"%s\"\n", dep.Name))
		sb.WriteString(fmt.Sprintf("version = \"%s\"\n", dep.Version))
		sb.WriteString("type = \"runtime\"\n\n")
	}

	return sb.String()
}

func flattenManifest(m *manifest.Manifest, options []string) *manifest.Manifest {
	// Create a copy
	flattened := *m

	// Filter dependencies based on options
	accessor := manifest.NewManifestAccessor(m)
	flattened.Specifications.Build.Dependencies = accessor.GetDependenciesForOptions(options)

	// Remove conditional steps from recipe
	// (This is a simplified version - full implementation would evaluate all conditionals)

	return &flattened
}

func performLintChecks(m *manifest.Manifest) []manifest.ValidationError {
	issues := make([]manifest.ValidationError, 0)

	// Check for overly long description
	if len(m.Description) > 200 {
		issues = append(issues, manifest.ValidationError{
			Field:    "description",
			Message:  "Description is very long (>200 chars), consider shortening",
			Severity: "warning",
		})
	}

	// Check for missing homepage in metadata
	if m.Metadata.Homepage == "" {
		issues = append(issues, manifest.ValidationError{
			Field:    "metadata.homepage",
			Message:  "Homepage URL is recommended for discoverability",
			Severity: "info",
		})
	}

	// Check for empty recipe phases
	if len(m.Recipe.Build) == 0 {
		issues = append(issues, manifest.ValidationError{
			Field:    "recipe.build",
			Message:  "Build phase has no steps",
			Severity: "warning",
		})
	}

	return issues
}

func performDeepChecks(m *manifest.Manifest) []string {
	issues := make([]string, 0)

	// Check for incompatible options
	// (This would require more complex logic to detect actual incompatibilities)

	// Check for circular dependencies in options
	// (Simplified check)

	// Check for missing toolchain for recipe steps
	hasCMake := false
	for _, tool := range m.Specifications.Build.Toolchain {
		if tool.Name == "cmake" {
			hasCMake = true
			break
		}
	}

	// Check if recipe uses cmake commands
	for _, step := range m.Recipe.Configuration {
		if strings.Contains(strings.ToLower(step.Command), "cmake") && !hasCMake {
			issues = append(issues, "Recipe uses cmake but cmake is not in toolchain")
			break
		}
	}

	return issues
}

func runManifestValidate(cmd *cobra.Command, args []string) error {
	manifestSource := args[0]

	// Load manifest
	m, err := loader.LoadManifest(manifestSource)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	fmt.Printf("Validating manifest: %s@%s\n", m.Name, m.Version)
	fmt.Println()

	// Validate
	result := manifest.ValidateManifest(m)

	// Print results
	if len(result.Errors) > 0 {
		fmt.Println("ERRORS:")
		for _, err := range result.Errors {
			fmt.Printf("  ✗ %s\n", err)
		}
		fmt.Println()
	}

	if len(result.Warnings) > 0 {
		fmt.Println("WARNINGS:")
		for _, warn := range result.Warnings {
			fmt.Printf("  ⚠ %s\n", warn)
		}
		fmt.Println()
	}

	if validateVerbose && len(result.Info) > 0 {
		fmt.Println("INFO:")
		for _, info := range result.Info {
			fmt.Printf("  ℹ %s\n", info)
		}
		fmt.Println()
	}

	// Summary
	if result.Valid {
		fmt.Println("✓ Manifest is valid")
		return nil
	} else {
		fmt.Printf("✗ Manifest validation failed with %d error(s)\n", len(result.Errors))
		return fmt.Errorf("validation failed")
	}
}

func runManifestShow(cmd *cobra.Command, args []string) error {
	manifestSource := args[0]

	// Load manifest
	m, err := loader.LoadManifest(manifestSource)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	// Generate report based on format
	var output string
	switch showFormat {
	case "text":
		output, err = reporters.GenerateReport(m, reporters.FormatText)
	case "yaml":
		output, err = reporters.GenerateReport(m, reporters.FormatJSON) // Will output YAML
	case "json":
		output, err = reporters.GenerateReport(m, reporters.FormatJSON)
	case "markdown", "md":
		output, err = reporters.GenerateReport(m, reporters.FormatMarkdown)
	default:
		return fmt.Errorf("unsupported format: %s", showFormat)
	}

	if err != nil {
		return fmt.Errorf("failed to generate report: %w", err)
	}

	fmt.Print(output)
	return nil
}

func runManifestInspect(cmd *cobra.Command, args []string) error {
	manifestSource := args[0]
	field := ""
	if len(args) > 1 {
		field = args[1]
	}

	// Load manifest
	m, err := loader.LoadManifest(manifestSource)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	accessor := manifest.NewManifestAccessor(m)

	// If no field specified, show summary
	if field == "" {
		fmt.Printf("Manifest: %s@%s\n", accessor.Name(), accessor.Version())
		fmt.Printf("Description: %s\n\n", accessor.Description())

		fmt.Println("Available fields:")
		fmt.Println("  name, version, description")
		fmt.Println("  source, source.type, source.url, source.tag, source.checksum")
		fmt.Println("  metadata, metadata.authors, metadata.license, metadata.homepage")
		fmt.Println("  specifications.build.dependencies")
		fmt.Println("  specifications.runtime.dependencies")
		fmt.Println("  specifications.build.toolchain")
		fmt.Println("  specifications.build.targets")
		fmt.Println("  specifications.build.options")
		fmt.Println("  recipe, recipe.configuration, recipe.build, recipe.install")
		return nil
	}

	// Inspect specific field
	value, err := inspectField(accessor, field)
	if err != nil {
		return err
	}

	fmt.Println(value)
	return nil
}

func inspectField(accessor *manifest.ManifestAccessor, field string) (string, error) {
	switch field {
	case "name":
		return accessor.Name(), nil
	case "version":
		return accessor.Version(), nil
	case "description":
		return accessor.Description(), nil
	case "source":
		return fmt.Sprintf("Type: %s\nURL: %s\nTag: %s",
			accessor.SourceType(), accessor.SourceURL(), accessor.SourceTag()), nil
	case "source.type":
		return accessor.SourceType(), nil
	case "source.url":
		return accessor.SourceURL(), nil
	case "source.tag":
		return accessor.SourceTag(), nil
	case "source.checksum":
		return accessor.SourceChecksum(), nil
	case "metadata.authors":
		return formatList(accessor.Authors()), nil
	case "metadata.license":
		return accessor.License(), nil
	case "metadata.homepage":
		return accessor.Homepage(), nil
	case "specifications.build.dependencies":
		return formatDependencies(accessor.BuildDependencies()), nil
	case "specifications.runtime.dependencies":
		return formatDependencies(accessor.RuntimeDependencies()), nil
	case "specifications.build.toolchain":
		return formatToolchain(accessor.Toolchain()), nil
	case "specifications.build.options":
		return formatList(accessor.Options()), nil
	default:
		return "", fmt.Errorf("unknown field: %s", field)
	}
}

func runManifestWalk(cmd *cobra.Command, args []string) error {
	manifestSource := args[0]

	// Load manifest
	m, err := loader.LoadManifest(manifestSource)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	// Walk manifest
	result, err := manifest.WalkManifest(m, walkOptions, walkVariables)
	if err != nil {
		return fmt.Errorf("failed to walk manifest: %w", err)
	}

	// Print result
	output := manifest.PrintWalkResult(result)
	fmt.Print(output)

	return nil
}

// Helper functions

func formatList(items []string) string {
	if len(items) == 0 {
		return "(none)"
	}
	result := ""
	for _, item := range items {
		result += fmt.Sprintf("  - %s\n", item)
	}
	return result
}

func formatDependencies(deps []manifest.Dependency) string {
	if len(deps) == 0 {
		return "(none)"
	}
	result := ""
	for _, dep := range deps {
		result += fmt.Sprintf("  - %s %s", dep.Name, dep.Version)
		if len(dep.ForOptions) > 0 {
			result += fmt.Sprintf(" (for: %v)", dep.ForOptions)
		}
		result += "\n"
	}
	return result
}

func formatToolchain(tools []manifest.Tool) string {
	if len(tools) == 0 {
		return "(none)"
	}
	result := ""
	for _, tool := range tools {
		result += fmt.Sprintf("  - %s %s\n", tool.Name, tool.Version)
	}
	return result
}
