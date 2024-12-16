package utils

import (
	"fmt"
	"time"
)

// logMessage formats and prints log messages with a timestamp and color
func logMessage(level, message, color string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("%s%s - %s --> %s\033[0m\n", color, level, timestamp, message)
}

// logError logs an error message in red
func LogError(message string) {
	logMessage("ERROR", message, "\033[31m") // Red
}

// logInfo logs an info message in blue
func LogInfo(message string) {
	logMessage("INFO", message, "\033[34m") // Blue
}

// logSuccess logs a success message in green
func LogSuccess(message string) {
	logMessage("SUCCESS", message, "\033[32m") // Green
}
