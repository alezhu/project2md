package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Simple glob matching (supports * and ?)
func matchGlob(pattern, name string) bool {
	// Convert to lowercase for case-insensitive matching on some systems
	if filepath.Separator == '\\' {
		pattern = strings.ToLower(pattern)
		name = strings.ToLower(name)
	}

	return matchGlobRecursive(pattern, name)
}

func matchGlobRecursive(pattern, name string) bool {
	for len(pattern) > 0 {
		switch pattern[0] {
		case '*':
			// Skip consecutive asterisks
			for len(pattern) > 0 && pattern[0] == '*' {
				pattern = pattern[1:]
			}

			if len(pattern) == 0 {
				return true
			}

			// Try matching the rest of the pattern at each position in name
			for i := 0; i <= len(name); i++ {
				if matchGlobRecursive(pattern, name[i:]) {
					return true
				}
			}
			return false

		case '?':
			if len(name) == 0 {
				return false
			}
			pattern = pattern[1:]
			name = name[1:]

		default:
			if len(name) == 0 || pattern[0] != name[0] {
				return false
			}
			pattern = pattern[1:]
			name = name[1:]
		}
	}

	return len(name) == 0
}

// Format file size in human-readable format
func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// Write file content with proper error handling
func writeFileContent(writer io.Writer, format string, args ...interface{}) error {
	_, err := fmt.Fprintf(writer, format, args...)
	return err
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil // путь существует
	}
	if os.IsNotExist(err) {
		return false, nil // путь не существует
	}
	return false, err // другая ошибка
}
