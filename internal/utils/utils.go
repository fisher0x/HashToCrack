package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// CleanLine removes CRLF and trims whitespace
func CleanLine(line string) string {
	// Remove carriage returns (Windows line endings)
	line = strings.ReplaceAll(line, "\r", "")
	return strings.TrimSpace(line)
}

// BoolToIncluded converts a boolean to "Included" or "Excluded" string
func BoolToIncluded(b bool) string {
	if b {
		return "Included"
	}
	return "Excluded"
}

// EnsureDir creates the directory for a file path if it doesn't exist
func EnsureDir(filePath string) error {
	dir := filepath.Dir(filePath)
	if dir != "" && dir != "." {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}
