package manifestCmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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

func getDefaultNameOrNameOfDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "example"
	}
	return filepath.Base(cwd)
}
