package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Statistics holds processing statistics
type Statistics struct {
	ProcessedFiles int
	SkippedDirs    int
	TotalSize      int64
	StartTime      time.Time
}

// Check if file should be processed as source code
func isCodeFile(filename string, config Config) (bool, bool) {
	var (
		result bool
		exists bool
	)
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		// Check files without extension
		base := strings.ToLower(filepath.Base(filename))
		result, exists = config.CodeExtensions[base]
	}
	result, exists = config.CodeExtensions[ext]
	return result, exists
}

// Check if directory should be skipped
func shouldSkipDir(dirname string, config Config) (bool, bool) {
	result, exists := config.SkipDirs[strings.ToLower(dirname)]
	return result, exists
}

func isFilenameInFileSet(filename string, set map[string]struct{}) bool {
	_, exists := set[filename]
	if exists {
		return exists
	}

	for patt := range set {
		patt = filepath.Clean(patt)
		if matchGlob(patt, filename) {
			return true
		}
	}
	return false
}

// check if file should be processed
// 1 - should be processed
// -1 - should not be processed
// 0 - unknown
func shouldProcessFile(filename string, config Config) int {
	// Convert to lowercase for case-insensitive matching on some systems
	if filepath.Separator == '\\' {
		filename = strings.ToLower(filename)
	}
	//White list has more priority
	if isFilenameInFileSet(filename, config.Include) {
		return 1
	}
	if isFilenameInFileSet(filename, config.Exclude) {
		return -1
	}
	return 0
}

// Get language identifier for syntax highlighting
func getLanguage(filename string, config Config) string {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		base := strings.ToLower(filepath.Base(filename))
		if lang, exists := config.Languages[base]; exists {
			return lang
		}
		return ""
	}

	if lang, exists := config.Languages[ext]; exists {
		return lang
	}
	return ""
}

// Process project directory and generate archive
func processProject(
	projectPath string,
	defaultConfig Config,
	customConfig Config,
	outputFileName string,
	verbose bool,
	showStats bool,
	noGit bool,
) error {

	var outputFile string
	if filepath.IsAbs(outputFileName) {
		outputFile = outputFileName
	} else {
		outputFile = filepath.Clean(filepath.Join(projectPath, outputFileName))
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// Load .gitignore patterns
	var gitIgnore *GitIgnore
	if noGit {
		gitIgnore = &GitIgnore{}
	} else {
		gitIgnore = loadGitIgnore(projectPath)
	}

	stats := &Statistics{
		StartTime: time.Now(),
	}

	// Write header
	if err := writeFileContent(writer, "# Code Archive: %s\n\n", filepath.Base(projectPath)); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	if err := writeFileContent(writer, "Generated automatically from: `%s`\n", projectPath); err != nil {
		return fmt.Errorf("failed to write path: %w", err)
	}

	if err := writeFileContent(writer, "Generated at: %s\n\n", time.Now().Format("2006-01-02 15:04:05")); err != nil {
		return fmt.Errorf("failed to write timestamp: %w", err)
	}

	if err := writeFileContent(writer, "---\n\n"); err != nil {
		return fmt.Errorf("failed to write separator: %w", err)
	}

	// Collect all files first for better organization
	var files []string
	err = filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if verbose {
				log.Printf("Warning: error accessing %s: %v", path, err)
			}
			return nil
		}

		// Skip directories in exclude list
		if info.IsDir() {
			//Check customConfig first
			skip, exists := shouldSkipDir(info.Name(), customConfig)

			if skip {
				//If dir is explicitly skiped in custom config then skip this dir
				if verbose {
					fmt.Printf("Ignored by custom config: %s\n", path)
				}
				return filepath.SkipDir
			} else {
				if exists {
					//If dir is explicitly excluded from skip-dir-list in custom config then process this dir
				} else {
					// Check .gitignore
					if !noGit && gitIgnore.isMatchPattern(path, true) {
						if verbose {
							fmt.Printf("Ignored by .gitignore: %s\n", path)
						}
						stats.SkippedDirs++
						return filepath.SkipDir
					}
					//Check defaultConfig
					if skip, _ = shouldSkipDir(info.Name(), defaultConfig); skip {
						if verbose {
							fmt.Printf("Skipping directory: %s\n", path)
						}
						stats.SkippedDirs++
						return filepath.SkipDir
					}
				}
			}
			return nil
		}

		// Check if file should be processed
		filePath := info.Name()
		relPath, _ := filepath.Rel(projectPath, path)
		shouldProcess := shouldProcessFile(relPath, customConfig)
		switch {
		case shouldProcess == -1:
			return nil
		case shouldProcess == 0:
			isCode, exists := isCodeFile(filePath, customConfig)
			if !isCode && exists {
				//If file explicit excluded in custom config, then skip file
				if verbose {
					fmt.Printf("Ignored by custom config: %s\n", filePath)
				}
				return nil
			}
			if !noGit && gitIgnore.isMatchPattern(filePath, false) {
				if verbose {
					fmt.Printf("Ignored by .gitignore: %s\n", filePath)
				}
				return nil
			}
			shouldProcess = shouldProcessFile(relPath, defaultConfig)
			if shouldProcess == -1 {
				if verbose {
					fmt.Printf("Ignored by default config: %s\n", filePath)
				}
				return nil
			} else if shouldProcess == 0 {
				if !isCode {
					isCode, exists = isCodeFile(filePath, defaultConfig)
				}
				if isCode {
					shouldProcess = 1
				}
			}
		}

		if shouldProcess == 1 {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking directory: %w", err)
	}

	// Sort files for consistent output
	sort.Strings(files)

	// Process files
	for _, path := range files {
		info, err := os.Stat(path)
		if err != nil {
			if verbose {
				log.Printf("Warning: cannot stat file %s: %v", path, err)
			}
			continue
		}

		// Get relative path
		relPath, err := filepath.Rel(projectPath, path)
		if err != nil {
			if verbose {
				log.Printf("Warning: cannot get relative path for %s: %v", path, err)
			}
			continue
		}

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			if verbose {
				log.Printf("Warning: cannot read file %s: %v", relPath, err)
			}
			continue
		}

		// Get language for syntax highlighting
		language := getLanguage(info.Name(), defaultConfig)

		// Write file section
		if err := writeFileContent(writer, "=== %s ===\n\n", relPath); err != nil {
			return fmt.Errorf("failed to write file header: %w", err)
		}

		if showStats {
			if err := writeFileContent(writer, "*Size: %s, Modified: %s*\n\n",
				formatFileSize(info.Size()),
				info.ModTime().Format("2006-01-02 15:04:05")); err != nil {
				return fmt.Errorf("failed to write file info: %w", err)
			}
		}

		if err := writeFileContent(writer, "```%s\n", language); err != nil {
			return fmt.Errorf("failed to write code block start: %w", err)
		}

		if err := writeFileContent(writer, "%s", string(content)); err != nil {
			return fmt.Errorf("failed to write file content: %w", err)
		}

		if !strings.HasSuffix(string(content), "\n") {
			if err := writeFileContent(writer, "\n"); err != nil {
				return fmt.Errorf("failed to write newline: %w", err)
			}
		}

		if err := writeFileContent(writer, "```\n\n"); err != nil {
			return fmt.Errorf("failed to write code block end: %w", err)
		}

		stats.ProcessedFiles++
		stats.TotalSize += info.Size()

		if verbose {
			fmt.Printf("Processed: %s (%s)\n", relPath, formatFileSize(info.Size()))
		}
	}

	// Write statistics only if showStats is true
	if showStats {
		duration := time.Since(stats.StartTime)
		if err := writeFileContent(writer, "---\n\n"); err != nil {
			return fmt.Errorf("failed to write statistics separator: %w", err)
		}

		if err := writeFileContent(writer, "## Statistics\n\n"); err != nil {
			return fmt.Errorf("failed to write statistics header: %w", err)
		}

		if err := writeFileContent(writer, "- **Files processed**: %d\n", stats.ProcessedFiles); err != nil {
			return fmt.Errorf("failed to write file count: %w", err)
		}

		if err := writeFileContent(writer, "- **Directories skipped**: %d\n", stats.SkippedDirs); err != nil {
			return fmt.Errorf("failed to write dir count: %w", err)
		}

		if err := writeFileContent(writer, "- **Total size**: %s\n", formatFileSize(stats.TotalSize)); err != nil {
			return fmt.Errorf("failed to write total size: %w", err)
		}

		if err := writeFileContent(writer, "- **Processing time**: %v\n", duration.Round(time.Millisecond)); err != nil {
			return fmt.Errorf("failed to write duration: %w", err)
		}
	}

	fmt.Printf("\nArchive created: %s\n", outputFile)
	fmt.Printf("Statistics: %d files, %s", stats.ProcessedFiles, formatFileSize(stats.TotalSize))
	if showStats {
		duration := time.Since(stats.StartTime)
		fmt.Printf(", %v processing time", duration.Round(time.Millisecond))
	}
	fmt.Println()

	return nil
}
