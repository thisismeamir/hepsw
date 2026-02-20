package manifest

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// WalkResult contains information about a manifest walk
type WalkResult struct {
	Manifest  *Manifest
	Phases    []PhaseWalk
	Variables map[string]string
}

// PhaseWalk represents the walk through a recipe phase
type PhaseWalk struct {
	Name  string
	Steps []StepWalk
}

// StepWalk represents a single recipe step walk
type StepWalk struct {
	Index       int
	Name        string
	Command     string
	Script      string
	Args        []string
	WorkingDir  string
	Conditional string
	WillExecute bool
	Reason      string
}

// WalkManifest simulates walking through a manifest's recipe
func WalkManifest(m *Manifest, options []string, variables map[string]string) (*WalkResult, error) {
	result := &WalkResult{
		Manifest:  m,
		Phases:    make([]PhaseWalk, 0),
		Variables: make(map[string]string),
	}

	// Initialize default variables
	initializeDefaultVariables(result.Variables, m)

	// Merge provided variables
	for k, v := range variables {
		result.Variables[k] = v
	}

	// Walk through each phase
	phases := []struct {
		name  string
		steps []RecipeStep
	}{
		{"configuration", m.Recipe.Configuration},
		{"build", m.Recipe.Build},
		{"install", m.Recipe.Install},
		{"use", m.Recipe.Use},
	}

	for _, phase := range phases {
		if len(phase.steps) > 0 {
			phaseWalk := walkPhase(phase.name, phase.steps, options, result.Variables)
			result.Phases = append(result.Phases, phaseWalk)
		}
	}

	return result, nil
}

func walkPhase(name string, steps []RecipeStep, options []string, variables map[string]string) PhaseWalk {
	phaseWalk := PhaseWalk{
		Name:  name,
		Steps: make([]StepWalk, 0),
	}

	for i, step := range steps {
		stepWalk := StepWalk{
			Index:       i,
			Name:        step.Name,
			Command:     step.Command,
			Script:      step.Script,
			Args:        step.Args,
			WorkingDir:  step.WorkingDir,
			Conditional: step.If,
			WillExecute: true,
			Reason:      "",
		}

		// Evaluate conditional
		if step.If != "" {
			willExecute, reason := evaluateConditional(step.If, options, variables)
			stepWalk.WillExecute = willExecute
			stepWalk.Reason = reason
		}

		// Process set variables
		if step.Set != nil {
			for k, v := range step.Set {
				expanded := expandVariables(v, variables)
				variables[k] = expanded
				stepWalk.Reason = fmt.Sprintf("Sets variables: %v", step.Set)
			}
		}

		phaseWalk.Steps = append(phaseWalk.Steps, stepWalk)
	}

	return phaseWalk
}

func evaluateConditional(condition string, options []string, variables map[string]string) (bool, string) {
	// Expand variables in condition
	expanded := expandVariables(condition, variables)

	// Check for negation
	negated := false
	if strings.HasPrefix(expanded, "!") {
		negated = true
		expanded = strings.TrimPrefix(expanded, "!")
	}

	// Remove ${} wrappers
	expanded = strings.TrimPrefix(expanded, "${")
	expanded = strings.TrimSuffix(expanded, "}")

	// Check if it's an option
	if strings.HasPrefix(expanded, "OPTIONS_") {
		optName := strings.TrimPrefix(expanded, "OPTIONS_")
		optName = strings.ReplaceAll(optName, "_", "-")
		optName = strings.ToLower(optName)

		hasOption := false
		for _, opt := range options {
			if strings.ToLower(opt) == optName {
				hasOption = true
				break
			}
		}

		if negated {
			hasOption = !hasOption
		}

		reason := fmt.Sprintf("Option '%s' is %s", optName, map[bool]string{true: "enabled", false: "disabled"}[hasOption])
		return hasOption, reason
	}

	// Check if it's a variable
	if val, ok := variables[expanded]; ok {
		result := val != "" && val != "0" && val != "false"
		if negated {
			result = !result
		}
		reason := fmt.Sprintf("Variable '%s' is '%s'", expanded, val)
		return result, reason
	}

	// Default to true if we can't evaluate
	return true, fmt.Sprintf("Cannot evaluate condition: %s", condition)
}

func initializeDefaultVariables(variables map[string]string, m *Manifest) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	hepswPath := filepath.Join(homeDir, ".hepsw")

	variables["PACKAGE_NAME"] = m.Name
	variables["PACKAGE_VERSION"] = m.Version
	variables["SOURCE_TYPE"] = m.Source.Type
	variables["SOURCE_URL"] = m.Source.Url
	variables["INSTALL_PREFIX"] = filepath.Join(hepswPath, "installs", m.Name, m.Version)
	variables["SOURCE_DIR"] = filepath.Join(hepswPath, "sources", m.Name, m.Version)
	variables["BUILD_DIR"] = filepath.Join(hepswPath, "builds", m.Name, m.Version)

	// Add build variables from manifest
	for _, varMap := range m.Specifications.Build.Variables {
		for k, v := range varMap {
			variables[k] = v
		}
	}
	return nil
}

func expandVariables(input string, variables map[string]string) string {
	result := input

	// Replace ${var} patterns
	for k, v := range variables {
		patterns := []string{
			fmt.Sprintf("${%s}", k),
			fmt.Sprintf("$%s", k),
		}
		for _, pattern := range patterns {
			result = strings.ReplaceAll(result, pattern, v)
		}
	}

	return result
}

// PrintWalkResult prints a human-readable walk result
func PrintWalkResult(result *WalkResult) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Walking manifest: %s@%s\n", result.Manifest.Name, result.Manifest.Version))
	sb.WriteString(fmt.Sprintf("Description: %s\n\n", result.Manifest.Description))

	sb.WriteString("Variables:\n")
	for k, v := range result.Variables {
		sb.WriteString(fmt.Sprintf("  %s = %s\n", k, v))
	}
	sb.WriteString("\n")

	for _, phase := range result.Phases {
		sb.WriteString(fmt.Sprintf("=== Phase: %s ===\n", strings.ToUpper(phase.Name)))

		for _, step := range phase.Steps {
			status := "✓"
			if !step.WillExecute {
				status = "✗"
			}

			sb.WriteString(fmt.Sprintf("\n[%s] Step %d: %s\n", status, step.Index+1, step.Name))

			if step.Conditional != "" {
				sb.WriteString(fmt.Sprintf("    Condition: %s\n", step.Conditional))
				if step.Reason != "" {
					sb.WriteString(fmt.Sprintf("    Reason: %s\n", step.Reason))
				}
			}

			if step.WillExecute {
				if step.Command != "" {
					expanded := expandVariables(step.Command, result.Variables)
					sb.WriteString(fmt.Sprintf("    Command: %s\n", expanded))
				}
				if step.Script != "" {
					sb.WriteString(fmt.Sprintf("    Script: %s\n", step.Script))
					if len(step.Args) > 0 {
						sb.WriteString(fmt.Sprintf("    Args: %v\n", step.Args))
					}
				}
				if step.WorkingDir != "" {
					expanded := expandVariables(step.WorkingDir, result.Variables)
					sb.WriteString(fmt.Sprintf("    Working Dir: %s\n", expanded))
				}
			} else {
				sb.WriteString(fmt.Sprintf("    Skipped: %s\n", step.Reason))
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
