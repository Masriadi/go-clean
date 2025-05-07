package utils

import (
	"strings"
	"unicode"
)

// StringToFileName converts a string to a directory name
func StringToFileName(s string) string {
	// Handle empty string
	if len(s) == 0 {
		return s
	}

	var result strings.Builder
	for i, char := range s {
		if i > 0 && unicode.IsUpper(char) {
			result.WriteRune('_')
			result.WriteRune(unicode.ToLower(char))
		} else {
			result.WriteRune(unicode.ToLower(char))
		}
	}
	return result.String()
}

// StringToEntityName converts a string to an entity name
func StringToEntityName(s string) string {
	// Handle empty string
	if len(s) == 0 {
		return s
	}

	// Convert first character to upper case and concatenate with rest of string
	return strings.ToUpper(s[:1]) + s[1:]
}

// StringToInstanceName converts a string to an instance name
func StringToInstanceName(s string) string {
	// Handle empty string
	if len(s) == 0 {
		return s
	}

	// Convert first character to lower case and concatenate with rest of string
	return strings.ToLower(s[:1]) + s[1:]
}
