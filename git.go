package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GitIgnore represents a .gitignore parser
type GitIgnore struct {
	patterns []GitIgnorePattern
}

// GitIgnorePattern represents a single .gitignore pattern
type GitIgnorePattern struct {
	pattern   string
	isNegated bool
	isDir     bool
	basePath  string
}

// Parse .gitignore patterns
func parseGitIgnorePattern(line, basePath string) *GitIgnorePattern {
	line = strings.TrimSpace(line)

	// Skip empty lines and comments
	if line == "" || strings.HasPrefix(line, "#") {
		return nil
	}

	pattern := &GitIgnorePattern{
		basePath: basePath,
	}

	// Check for negation
	if strings.HasPrefix(line, "!") {
		pattern.isNegated = true
		line = line[1:]
	}

	// Check if pattern is for directories only
	if strings.HasSuffix(line, "/") {
		pattern.isDir = true
		line = strings.TrimSuffix(line, "/")
	}

	pattern.pattern = line
	return pattern
}

// Load .gitignore files recursively
func loadGitIgnore(rootPath string) (*GitIgnore, error) {
	gi := &GitIgnore{}

	// Walk through all directories and load .gitignore files
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.Name() == ".gitignore" {
			file, err := os.Open(path)
			if err != nil {
				return nil
			}
			defer func() {
				_ = file.Close()
			}()

			basePath := filepath.Dir(path)
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				if pattern := parseGitIgnorePattern(scanner.Text(), basePath); pattern != nil {
					gi.patterns = append(gi.patterns, *pattern)
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error walking directory: %w", err)
	}

	return gi, nil
}

// Check if path matches a gitignore pattern
func (gi *GitIgnore) matchesPattern(path string, isDir bool, pattern GitIgnorePattern) bool {
	// Get relative path from pattern's base directory
	relPath, err := filepath.Rel(pattern.basePath, path)
	if err != nil {
		return false
	}

	// If pattern is for directories only, check if path is directory
	if pattern.isDir && !isDir {
		return false
	}

	patternStr := pattern.pattern

	// Handle absolute patterns (starting with /)
	if strings.HasPrefix(patternStr, "/") {
		patternStr = patternStr[1:]
		// For absolute patterns, match from the pattern's base directory
		return matchGlob(patternStr, relPath)
	}

	// For relative patterns, check if any part of the path matches
	pathParts := strings.Split(relPath, string(filepath.Separator))

	// Check if the pattern matches the full relative path
	if matchGlob(patternStr, relPath) {
		return true
	}

	// Check if pattern matches any directory component
	for i := 0; i < len(pathParts); i++ {
		subPath := strings.Join(pathParts[i:], string(filepath.Separator))
		if matchGlob(patternStr, subPath) {
			return true
		}

		// Also check individual components
		if matchGlob(patternStr, pathParts[i]) {
			return true
		}
	}

	return false
}

// Check if path should be ignored by .gitignore
func (gi *GitIgnore) isMatchPattern(path string, isDir bool) bool {
	result := false

	// Process patterns in order
	for _, pattern := range gi.patterns {
		if gi.matchesPattern(path, isDir, pattern) {
			if pattern.isNegated {
				result = false // Negated pattern includes the file
			} else {
				result = true // Normal pattern excludes the file
			}
		}
	}

	return result
}
