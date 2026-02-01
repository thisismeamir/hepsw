package manifest

import (
	"fmt"
	"strings"
)

// ManifestAccessor provides convenient access to manifest fields
type ManifestAccessor struct {
	manifest *Manifest
}

// NewManifestAccessor creates a new ManifestAccessor
func NewManifestAccessor(m *Manifest) *ManifestAccessor {
	return &ManifestAccessor{manifest: m}
}

// Basic Information

func (ma *ManifestAccessor) Name() string {
	return ma.manifest.Name
}

func (ma *ManifestAccessor) Version() string {
	return ma.manifest.Version
}

func (ma *ManifestAccessor) Description() string {
	return ma.manifest.Description
}

// Source Information

func (ma *ManifestAccessor) Source() SourceSpec {
	return ma.manifest.Source
}

func (ma *ManifestAccessor) SourceType() string {
	return ma.manifest.Source.Type
}

func (ma *ManifestAccessor) SourceURL() string {
	return ma.manifest.Source.Url
}

func (ma *ManifestAccessor) SourceTag() string {
	return ma.manifest.Source.Tag
}

func (ma *ManifestAccessor) SourceChecksum() string {
	return ma.manifest.Source.Checksum
}

// Metadata

func (ma *ManifestAccessor) Metadata() ManifestMetaData {
	return ma.manifest.Metadata
}

func (ma *ManifestAccessor) Authors() []string {
	return ma.manifest.Metadata.Authors
}

func (ma *ManifestAccessor) Homepage() string {
	return ma.manifest.Metadata.Homepage
}

func (ma *ManifestAccessor) License() string {
	return ma.manifest.Metadata.License
}

func (ma *ManifestAccessor) Documentation() string {
	return ma.manifest.Metadata.Documentation
}

// Dependencies

func (ma *ManifestAccessor) BuildDependencies() []Dependency {
	return ma.manifest.Specifications.Build.Dependencies
}

func (ma *ManifestAccessor) RuntimeDependencies() []Dependency {
	return ma.manifest.Specifications.Runtime.Dependencies
}

func (ma *ManifestAccessor) AllDependencies() []Dependency {
	deps := make([]Dependency, 0)
	deps = append(deps, ma.BuildDependencies()...)
	deps = append(deps, ma.RuntimeDependencies()...)
	return deps
}

// Get dependencies filtered by options
func (ma *ManifestAccessor) GetDependenciesForOptions(options []string) []Dependency {
	allDeps := ma.AllDependencies()
	filtered := make([]Dependency, 0)

	for _, dep := range allDeps {
		if len(dep.ForOptions) == 0 {
			// No option requirement, always include
			filtered = append(filtered, dep)
			continue
		}

		// Check if any of the required options are present
		for _, reqOpt := range dep.ForOptions {
			for _, opt := range options {
				if reqOpt == opt {
					filtered = append(filtered, dep)
					break
				}
			}
		}
	}

	return filtered
}

// Build Specifications

func (ma *ManifestAccessor) Toolchain() []Tool {
	return ma.manifest.Specifications.Build.Toolchain
}

func (ma *ManifestAccessor) Targets() []Targets {
	return ma.manifest.Specifications.Build.Targets
}

func (ma *ManifestAccessor) Options() []string {
	return ma.manifest.Specifications.Build.Options
}

func (ma *ManifestAccessor) BuildVariables() []map[string]string {
	return ma.manifest.Specifications.Build.Variables
}

// Get a specific build variable
func (ma *ManifestAccessor) GetBuildVariable(key string) (string, bool) {
	for _, varMap := range ma.manifest.Specifications.Build.Variables {
		if val, ok := varMap[key]; ok {
			return val, true
		}
	}
	return "", false
}

// Environment Variables

func (ma *ManifestAccessor) BuildEnvironment() map[string]string {
	return ma.manifest.Specifications.Environment.Build
}

func (ma *ManifestAccessor) RuntimeEnvironment() map[string]string {
	return ma.manifest.Specifications.Environment.Runtime
}

func (ma *ManifestAccessor) SelfEnvironment() map[string]string {
	return ma.manifest.Specifications.Environment.Self
}

func (ma *ManifestAccessor) GetEnvironmentVariable(scope, key string) (string, bool) {
	var envMap map[string]string
	switch scope {
	case "build":
		envMap = ma.BuildEnvironment()
	case "runtime":
		envMap = ma.RuntimeEnvironment()
	case "self":
		envMap = ma.SelfEnvironment()
	default:
		return "", false
	}

	val, ok := envMap[key]
	return val, ok
}

// Recipe Access

func (ma *ManifestAccessor) Recipe() Recipe {
	return ma.manifest.Recipe
}

func (ma *ManifestAccessor) ConfigurationSteps() []RecipeStep {
	return ma.manifest.Recipe.Configuration
}

func (ma *ManifestAccessor) BuildSteps() []RecipeStep {
	return ma.manifest.Recipe.Build
}

func (ma *ManifestAccessor) InstallSteps() []RecipeStep {
	return ma.manifest.Recipe.Install
}

func (ma *ManifestAccessor) UseSteps() []RecipeStep {
	return ma.manifest.Recipe.Use
}

// Get all recipe steps in order
func (ma *ManifestAccessor) AllRecipeSteps() []RecipeStep {
	steps := make([]RecipeStep, 0)
	steps = append(steps, ma.ConfigurationSteps()...)
	steps = append(steps, ma.BuildSteps()...)
	steps = append(steps, ma.InstallSteps()...)
	steps = append(steps, ma.UseSteps()...)
	return steps
}

// Get steps by phase
func (ma *ManifestAccessor) GetStepsByPhase(phase string) []RecipeStep {
	switch strings.ToLower(phase) {
	case "configuration", "configure":
		return ma.ConfigurationSteps()
	case "build":
		return ma.BuildSteps()
	case "install":
		return ma.InstallSteps()
	case "use":
		return ma.UseSteps()
	default:
		return []RecipeStep{}
	}
}

// Utility Methods

// HasDependency checks if a dependency exists
func (ma *ManifestAccessor) HasDependency(name string) bool {
	for _, dep := range ma.AllDependencies() {
		if dep.Name == name {
			return true
		}
	}
	return false
}

// GetDependency retrieves a specific dependency
func (ma *ManifestAccessor) GetDependency(name string) (Dependency, bool) {
	for _, dep := range ma.AllDependencies() {
		if dep.Name == name {
			return dep, true
		}
	}
	return Dependency{}, false
}

// HasOption checks if an option exists
func (ma *ManifestAccessor) HasOption(option string) bool {
	for _, opt := range ma.Options() {
		if opt == option {
			return true
		}
	}
	return false
}

// SupportsTarget checks if a target is supported
func (ma *ManifestAccessor) SupportsTarget(targetName string) bool {
	for _, target := range ma.Targets() {
		if target.Name == targetName {
			return true
		}
	}
	return false
}

// GetFullIdentifier returns name@version
func (ma *ManifestAccessor) GetFullIdentifier() string {
	return fmt.Sprintf("%s@%s", ma.Name(), ma.Version())
}
