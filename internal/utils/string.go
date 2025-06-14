package utils

import (
	"strings"
	"unicode"
)

// StringToFileName converts a CamelCase or PascalCase string into a snake_case string.
// It adds underscores before uppercase letters (except the first character) and lowercases everything.
// Handles abbreviations like "HTTPServer" -> "http_server"
//
// Example:
//   StringToFileName("UserProfile")     => "user_profile"
//   StringToFileName("MyHTTPServer")    => "my_http_server"
//   StringToFileName("simple")          => "simple"
//   StringToFileName("")                => ""
func StringToFileName(s string) string {
	if len(s) == 0 {
		return s
	}

	var result strings.Builder
	for i, char := range s {
		if unicode.IsUpper(char) {
			if i > 0 && (unicode.IsLower(rune(s[i-1])) || (i+1 < len(s) && unicode.IsLower(rune(s[i+1])))) {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(char))
		} else {
			result.WriteRune(char)
		}
	}
	return result.String()
}

// StringToDirName converts a CamelCase or PascalCase string into a kebab-case string.
// It inserts hyphens before uppercase letters (except the first character) and lowercases everything.
//
// Example:
//   StringToDirName("UserProfile")     => "user-profile"
//   StringToDirName("MyHTTPServer")    => "my-http-server"
//   StringToDirName("simple")          => "simple"
//   StringToDirName("")                => ""
func StringToDirName(s string) string {
	if len(s) == 0 {
		return s
	}

	var result strings.Builder
	for i, char := range s {
		if unicode.IsUpper(char) {
			if i > 0 && (unicode.IsLower(rune(s[i-1])) || (i+1 < len(s) && unicode.IsLower(rune(s[i+1])))) {
				result.WriteRune('-')
			}
			result.WriteRune(unicode.ToLower(char))
		} else {
			result.WriteRune(char)
		}
	}
	return result.String()
}


// StringToEntityName converts the first character of a string to uppercase.
// Safe for Unicode characters.
//
// Example:
//   StringToEntityName("user")   => "User"
//   StringToEntityName("User")   => "User"
//   StringToEntityName("üser")   => "Üser"
//   StringToEntityName("")       => ""
func StringToEntityName(s string) string {
	if len(s) == 0 {
		return s
	}

	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// StringToInstanceName converts the first character of a string to lowercase.
// Safe for Unicode characters.
//
// Example:
//   StringToInstanceName("User")   => "user"
//   StringToInstanceName("user")   => "user"
//   StringToInstanceName("Üser")   => "üser"
//   StringToInstanceName("")       => ""
func StringToInstanceName(s string) string {
	if len(s) == 0 {
		return s
	}

	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}
