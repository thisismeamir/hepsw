package manifest

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

type ValidationResult struct {
	Valid    bool
	Errors   []ValidationError
	Warnings []ValidationError
	Info     []ValidationError
}
type ValidationError struct {
	Field    string
	Message  string
	Severity string // "error", "warning", "info"
}

func (ve ValidationError) String() string {
	return fmt.Sprintf("[%s] %s: %s", strings.ToUpper(ve.Severity), ve.Field, ve.Message)
}

func (vr *ValidationResult) AddError(field, message string) {
	vr.Valid = false
	vr.Errors = append(vr.Errors, ValidationError{
		Field:    field,
		Message:  message,
		Severity: "error",
	})
}

// AddWarning adds a warning to the validation result
func (vr *ValidationResult) AddWarning(field, message string) {
	vr.Warnings = append(vr.Warnings, ValidationError{
		Field:    field,
		Message:  message,
		Severity: "warning",
	})
}

// AddInfo adds an info message to the validation result
func (vr *ValidationResult) AddInfo(field, message string) {
	vr.Info = append(vr.Info, ValidationError{
		Field:    field,
		Message:  message,
		Severity: "info",
	})
}

// HasIssues returns true if there are any errors or warnings
func (vr *ValidationResult) HasIssues() bool {
	return len(vr.Errors) > 0 || len(vr.Warnings) > 0
}

// ValidateManifest performs comprehensive validation of a manifest
func ValidateManifest(m *Manifest) *ValidationResult {
	result := &ValidationResult{Valid: true}

	// Validate basic fields
	validateBasicFields(m, result)

	// Validate source
	validateSource(m, result)

	// Validate metadata
	validateMetadata(m, result)

	// Validate specifications
	validateSpecifications(m, result)

	// Validate recipe
	validateRecipe(m, result)

	// Check for circular dependencies
	validateDependencies(m, result)

	return result
}

func validateBasicFields(m *Manifest, result *ValidationResult) {
	// Name validation
	if m.Name == "" {
		result.AddError("name", "name is required")
	} else if !isValidPackageName(m.Name) {
		result.AddError("name", "name must be lowercase alphanumeric with hyphens only")
	}

	// Version validation
	if m.Version == "" {
		result.AddError("version", "version is required")
	} else if !isValidVersion(m.Version) {
		result.AddWarning("version", "version should follow semantic versioning (x.y.z)")
	}

	// Description validation
	if m.Description == "" {
		result.AddWarning("description", "description is recommended")
	} else if len(m.Description) < 10 {
		result.AddWarning("description", "description should be more descriptive")
	}
}

func validateSource(m *Manifest, result *ValidationResult) {
	if m.Source.Type == "" {
		result.AddError("source.type", "source type is required")
		return
	}

	validTypes := []string{"tarball", "git", "svn", "local"}
	if !contains(validTypes, m.Source.Type) {
		result.AddError("source.type", fmt.Sprintf("invalid source type: %s (must be one of: %s)",
			m.Source.Type, strings.Join(validTypes, ", ")))
	}

	if m.Source.Url == "" {
		result.AddError("source.url", "source URL is required")
	} else {
		// Validate URL format
		if _, err := url.Parse(m.Source.Url); err != nil {
			result.AddError("source.url", "invalid URL format")
		}
	}

	// Git-specific validation
	if m.Source.Type == "git" && m.Source.Tag == "" {
		result.AddWarning("source.tag", "git sources should specify a tag or branch")
	}

	// Checksum validation
	if m.Source.Checksum == "" {
		result.AddWarning("source.checksum", "checksum is recommended for reproducibility")
	} else if !isValidChecksum(m.Source.Checksum) {
		result.AddError("source.checksum", "invalid checksum format (should be algorithm:hash)")
	}
}

func validateMetadata(m *Manifest, result *ValidationResult) {
	if len(m.Metadata.Authors) == 0 {
		result.AddInfo("metadata.authors", "author information is recommended")
	}

	if m.Metadata.License == "" {
		result.AddWarning("metadata.license", "license information is recommended")
	}

	if m.Metadata.Homepage == "" {
		result.AddInfo("metadata.homepage", "homepage URL is recommended")
	}
}

func validateSpecifications(m *Manifest, result *ValidationResult) {
	// Validate toolchain
	if len(m.Specifications.Build.Toolchain) == 0 {
		result.AddWarning("specifications.build.toolchain", "toolchain requirements should be specified")
	}

	// Validate targets
	if len(m.Specifications.Build.Targets) == 0 {
		result.AddInfo("specifications.build.targets", "no build targets specified")
	}

	// Validate dependencies
	for i, dep := range m.Specifications.Build.Dependencies {
		validateDependency(dep, fmt.Sprintf("specifications.build.dependencies[%d]", i), result)
	}

	for i, dep := range m.Specifications.Runtime.Dependencies {
		validateDependency(dep, fmt.Sprintf("specifications.runtime.dependencies[%d]", i), result)
	}

	// Check for duplicate dependencies
	checkDuplicateDependencies(m, result)
}

func validateDependency(dep Dependency, prefix string, result *ValidationResult) {
	if dep.Name == "" {
		result.AddError(prefix+".name", "dependency name is required")
	}

	if dep.Version == "" {
		result.AddWarning(prefix+".version", "dependency version constraint is recommended")
	}

	// Validate version constraint format
	if dep.Version != "" && !isValidVersionConstraint(dep.Version) {
		result.AddWarning(prefix+".version", "version constraint format may be invalid")
	}
}

func validateRecipe(m *Manifest, result *ValidationResult) {
	// Check if recipe has any steps
	totalSteps := len(m.Recipe.Configuration) + len(m.Recipe.Build) +
		len(m.Recipe.Install) + len(m.Recipe.Use)

	if totalSteps == 0 {
		result.AddError("recipe", "recipe must contain at least one step")
		return
	}

	// Validate each phase
	validateRecipePhase(m.Recipe.Configuration, "recipe.configuration", result)
	validateRecipePhase(m.Recipe.Build, "recipe.build", result)
	validateRecipePhase(m.Recipe.Install, "recipe.install", result)
	validateRecipePhase(m.Recipe.Use, "recipe.use", result)

	// Check for common issues
	if len(m.Recipe.Build) == 0 {
		result.AddWarning("recipe.build", "no build steps defined")
	}

	if len(m.Recipe.Install) == 0 {
		result.AddWarning("recipe.install", "no install steps defined")
	}
}

func validateRecipePhase(steps []RecipeStep, prefix string, result *ValidationResult) {
	for i, step := range steps {
		stepPrefix := fmt.Sprintf("%s[%d]", prefix, i)

		if step.Name == "" {
			result.AddWarning(stepPrefix+".name", "step name is recommended")
		}

		// Check if step has either command or script
		if step.Command == "" && step.Script == "" && step.Set == nil {
			result.AddError(stepPrefix, "step must have command, script, or set field")
		}

		// Both command and script should not be set
		if step.Command != "" && step.Script != "" {
			result.AddError(stepPrefix, "step cannot have both command and script")
		}

		// Validate conditional syntax if present
		if step.If != "" && !isValidConditional(step.If) {
			result.AddWarning(stepPrefix+".if", "conditional syntax may be invalid")
		}
	}
}

func validateDependencies(m *Manifest, result *ValidationResult) {
	// Build dependency graph
	deps := make(map[string][]string)

	for _, dep := range m.Specifications.Build.Dependencies {
		if _, ok := deps[m.Name]; !ok {
			deps[m.Name] = []string{}
		}
		deps[m.Name] = append(deps[m.Name], dep.Name)
	}

	for _, dep := range m.Specifications.Runtime.Dependencies {
		if _, ok := deps[m.Name]; !ok {
			deps[m.Name] = []string{}
		}
		if !contains(deps[m.Name], dep.Name) {
			deps[m.Name] = append(deps[m.Name], dep.Name)
		}
	}

	// Note: Full circular dependency detection would require loading all dependent manifests
	// For now, we just check for self-dependencies
	for _, depName := range deps[m.Name] {
		if depName == m.Name {
			result.AddError("dependencies", "package cannot depend on itself")
		}
	}
}

func checkDuplicateDependencies(m *Manifest, result *ValidationResult) {
	seen := make(map[string]bool)

	for _, dep := range m.Specifications.Build.Dependencies {
		key := dep.Name
		if seen[key] {
			result.AddWarning("specifications.build.dependencies",
				fmt.Sprintf("duplicate dependency: %s", dep.Name))
		}
		seen[key] = true
	}

	seen = make(map[string]bool)
	for _, dep := range m.Specifications.Runtime.Dependencies {
		key := dep.Name
		if seen[key] {
			result.AddWarning("specifications.runtime.dependencies",
				fmt.Sprintf("duplicate dependency: %s", dep.Name))
		}
		seen[key] = true
	}
}

// Validation helper functions

func isValidPackageName(name string) bool {
	// Package names should be lowercase, alphanumeric with hyphens
	matched, _ := regexp.MatchString(`^[a-z0-9][a-z0-9-]*[a-z0-9]$|^[a-z0-9]$`, name)
	return matched
}

func isValidVersion(version string) bool {
	// Basic semantic versioning check
	matched, _ := regexp.MatchString(`^\d+\.\d+\.\d+`, version)
	return matched
}

func isValidChecksum(checksum string) bool {
	// Format: algorithm:hash
	parts := strings.Split(checksum, ":")
	if len(parts) != 2 {
		return false
	}

	validAlgos := []string{"md5", "sha1", "sha256", "sha512"}
	return contains(validAlgos, parts[0])
}

func isValidVersionConstraint(constraint string) bool {
	// Basic version constraint validation
	// Supports: >=1.0.0, <2.0.0, ==1.5.0, >=1.0.0,<2.0.0
	matched, _ := regexp.MatchString(`^[><=!]+\d+\.\d+`, constraint)
	return matched
}

func isValidConditional(cond string) bool {
	// Basic conditional syntax validation
	// Should contain ${...} variables
	return strings.Contains(cond, "${")
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
