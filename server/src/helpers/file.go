package helpers

import (
	"os"
	"strings"
)

func FileNameSanitize(fileName string) string {
	// Replace special characters with underscores
	specialChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range specialChars {
		fileName = strings.ReplaceAll(fileName, char, "_")
	}

	return fileName
}

func BuildSafeFileName(text string, maxLength int, extension string) string {
	// Remove any leading or trailing whitespace from the text
	text = strings.TrimSpace(text)

	// Replace any special characters or spaces with underscores
	text = strings.ReplaceAll(text, " ", "_")
	text = FileNameSanitize(text)

	// Truncate the text if it exceeds the maximum length
	if maxLength > 0 && len(text) > maxLength {
		text = text[:maxLength]
	}

	// Add the file extension if provided
	if extension != "" {
		text += "." + strings.TrimPrefix(extension, ".")
	}

	return text
}

func FileExists(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}
