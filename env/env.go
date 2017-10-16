// Package env provides simple convenient functions to work with
// environment variables.
package env

import (
	"fmt"
	"os"
	"strings"
)

// Must returns the value of the environment variable.
// It panics if the variable is not present.
// If the variable is present but the value is empty,
// it returns empty string.
func Must(variable string) string {
	variable = strings.TrimPrefix(variable, "$")
	value, ok := os.LookupEnv(variable)
	if !ok {
		panic(fmt.Sprintf("variable %s is not present in the environment", variable))
	}
	return value
}

// MustBool returns boolean value of the environment variable.
// It panics if variable is not present, or if value is neither true nor false.
func MustBool(variable string) bool {
	value := Must(variable)
	switch value {
	case "true":
		return true
	case "false":
		return false
	default:
		panic(fmt.Sprintf("environment variable %s must be either true or false, %s given", variable, value))
	}
}

// Get returns the value of the environment variable.
// If the variable is not present or is empty, returns defaultValue.
func Get(variable, defaultValue string) string {
	variable = strings.TrimPrefix(variable, "$")
	value := os.Getenv(variable)
	if value == "" {
		return defaultValue
	}
	return value
}

// Bool returns boolean value of the environment variable.
// If the variable is not present, is empty or is neither true nor false,
// returns defaultValue.
func Bool(variable string, defaultValue bool) bool {
	variable = strings.TrimPrefix(variable, "$")
	value := os.Getenv(variable)
	switch value {
	case "true":
		return true
	case "false":
		return false
	default:
		return defaultValue
	}
}
