# Code Archive: project2md

Generated automatically from: `e:\Projects\go\project2md`
Generated at: 2025-06-29 16:53:08

---

=== config.go ===

```go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config structures
type Config struct {
	CodeExtensions map[string]bool     `json:"code_extensions"`
	Exclude        map[string]struct{} `json:"-"`
	Include        map[string]struct{} `json:"-"`
	SkipDirs       map[string]bool     `json:"skip_dirs"`
	Languages      map[string]string   `json:"languages"`
}

func NewConfig() *Config {
	return &Config{
		CodeExtensions: map[string]bool{},
		Exclude:        map[string]struct{}{},
		Include:        map[string]struct{}{},
		SkipDirs:       map[string]bool{},
		Languages:      map[string]string{},
	}
}

func (c *Config) UnmarshalJSON(data []byte) error {
	type Alias Config
	aux := &struct {
		ExcludeArray []string `json:"exclude"`
		IncludeArray []string `json:"include"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	c.Exclude = make(map[string]struct{})
	for _, item := range aux.ExcludeArray {
		c.Exclude[filepath.Clean(item)] = struct{}{}
	}
	c.Include = make(map[string]struct{})
	for _, item := range aux.IncludeArray {
		c.Include[filepath.Clean(item)] = struct{}{}
	}

	return nil
}

// Default configuration
func getDefaultConfig() Config {
	return Config{
		CodeExtensions: map[string]bool{
			".bash":       true,
			".bat":        true,
			".bak":        false,
			".c":          true,
			".cc":         true,
			".cfg":        true,
			".clj":        true,
			".cmd":        true,
			".conf":       true,
			".cpp":        true,
			".cs":         true,
			".css":        true,
			".cxx":        true,
			".dockerfile": true,
			".env":        true,
			".fish":       true,
			".gitignore":  true,
			".go":         true,
			".h":          true,
			".hpp":        true,
			".hs":         true,
			".htm":        true,
			".html":       true,
			".ini":        true,
			".java":       true,
			".js":         true,
			".json":       true,
			".jsx":        true,
			".kt":         true,
			".less":       true,
			".log":        false,
			".m":          true,
			".md":         false,
			".ml":         true,
			".php":        true,
			".ps1":        true,
			".py":         true,
			".r":          true,
			".rb":         true,
			".rs":         true,
			".sass":       true,
			".scala":      true,
			".scss":       true,
			".sh":         true,
			".sql":        true,
			".swift":      true,
			".toml":       true,
			".ts":         true,
			".tsx":        true,
			".txt":        true,
			".xml":        true,
			".yaml":       true,
			".yml":        true,
			".zsh":        true,
			"makefile":    true,
		},
		SkipDirs: map[string]bool{
			".coverage":     true,
			".git":          true,
			".hg":           true,
			".idea":         true,
			".mypy_cache":   true,
			".nyc_output":   true,
			".pytest_cache": true,
			".svn":          true,
			".tox":          true,
			".vs":           true,
			".vscode":       true,
			"__pycache__":   true,
			"bin":           true,
			"build":         true,
			"coverage":      true,
			"dist":          true,
			"logs":          true,
			"node_modules":  true,
			"obj":           true,
			"out":           true,
			"target":        true,
			"temp":          true,
			"tmp":           true,
			"vendor":        true,
		},
		Languages: map[string]string{
			".bash":       "bash",
			".bat":        "batch",
			".c":          "c",
			".cc":         "cpp",
			".cfg":        "ini",
			".clj":        "clojure",
			".cmd":        "batch",
			".conf":       "ini",
			".cpp":        "cpp",
			".cs":         "csharp",
			".css":        "css",
			".cxx":        "cpp",
			".dockerfile": "dockerfile",
			".fish":       "bash",
			".go":         "go",
			".h":          "c",
			".hpp":        "c",
			".hs":         "haskell",
			".htm":        "html",
			".html":       "html",
			".ini":        "ini",
			".java":       "java",
			".js":         "javascript",
			".json":       "json",
			".jsx":        "jsx",
			".kt":         "kotlin",
			".less":       "less",
			".m":          "objectivec",
			".md":         "markdown",
			".ml":         "ocaml",
			".php":        "php",
			".ps1":        "powershell",
			".py":         "python",
			".r":          "r",
			".rb":         "ruby",
			".rs":         "rust",
			".sass":       "sass",
			".scala":      "scala",
			".scss":       "scss",
			".sh":         "bash",
			".sql":        "sql",
			".swift":      "swift",
			".toml":       "toml",
			".ts":         "typescript",
			".tsx":        "tsx",
			".xml":        "xml",
			".yaml":       "yaml",
			".yml":        "yaml",
			".zsh":        "bash",
			"makefile":    "makefile",
		},
		Exclude: map[string]struct{}{
			"project2md.config.json": {},
			"pnpm-lock.yaml":         {},
			"package-lock.json":      {},
			"yarn.lock":              {},
			".rsync-filter":          {},
			"LICENSE":                {},
		},
	}
}

// Load and apply user configuration overrides
func loadCustomConfig(config Config, configPath string) (Config, error) {
	if configPath == "" {
		return config, nil
	}

	file, err := os.Open(configPath)
	if err != nil {
		return config, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var customConfig Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&customConfig); err != nil {
		return config, fmt.Errorf("failed to decode config: %w", err)
	}

	// Apply overrides
	if len(customConfig.CodeExtensions) > 0 {
		if config.CodeExtensions == nil || len(config.CodeExtensions) == 0 {
			config.CodeExtensions = customConfig.CodeExtensions
		} else {
			for ext, enabled := range customConfig.CodeExtensions {
				config.CodeExtensions[ext] = enabled
			}
		}
	}

	if len(customConfig.SkipDirs) > 0 {
		if config.SkipDirs == nil || len(config.SkipDirs) == 0 {
			config.SkipDirs = customConfig.SkipDirs
		} else {
			for dir, skip := range customConfig.SkipDirs {
				config.SkipDirs[dir] = skip
			}
		}
	}

	if len(customConfig.Languages) > 0 {
		if config.Languages == nil || len(config.Languages) == 0 {
			config.Languages = customConfig.Languages
		} else {
			for ext, lang := range customConfig.Languages {
				config.Languages[ext] = lang
			}
		}
	}

	if len(customConfig.Exclude) > 0 {
		if config.Exclude == nil || len(config.Exclude) == 0 {
			config.Exclude = customConfig.Exclude
		} else {
			for patt := range customConfig.Exclude {
				config.Exclude[patt] = struct{}{}
			}
		}
	}

	if len(customConfig.Include) > 0 {
		if config.Include == nil || len(config.Include) == 0 {
			config.Include = customConfig.Include
		} else {
			for patt := range customConfig.Include {
				config.Include[patt] = struct{}{}
			}
		}
	}

	return config, nil
}
```

=== git.go ===

```go
package main

import (
	"bufio"
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
func loadGitIgnore(rootPath string) *GitIgnore {
	gi := &GitIgnore{}

	// Walk through all directories and load .gitignore files
	filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.Name() == ".gitignore" {
			file, err := os.Open(path)
			if err != nil {
				return nil
			}
			defer file.Close()

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

	return gi
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
```

=== main.go ===

```go
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Export default configuration to a file
func exportDefaultConfig(configPath string) error {
	config := getDefaultConfig()

	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	return nil
}

func main() {
	var (
		userConfigPath = flag.String("config", "", "Path to user configuration file")
		exportPath     = flag.String("export", "", "Export default configuration file to specified path")
		outputFileName = flag.String("output", "project.md", "Output file name")
		verbose        = flag.Bool("verbose", false, "Enable verbose output with file details")
		showStats      = flag.Bool("stat", false, "Include file size information and statistics in output")
		version        = flag.Bool("version", false, "Show version information")
		noGit          = flag.Bool("no-git", false, "Do not use .gitignore for exclude files")
	)
	flag.Parse()

	// Show version
	if *version {
		fmt.Println("Project2MD v1.0.0")
		fmt.Println("A tool to create markdown archives of source code projects")
		return
	}

	// Export configuration
	if *exportPath != "" {
		if err := exportDefaultConfig(*exportPath); err != nil {
			log.Fatalf("Error exporting configuration: %v", err)
		}
		fmt.Printf("Default configuration exported to: %s\n", *exportPath)
		fmt.Println("You can now edit this file to customize the archiving rules.")
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		exeFile := filepath.Base(os.Args[0])
		fmt.Printf("Usage: %s [options] <project_directory>\n\n", exeFile)
		fmt.Println("Options:")
		fmt.Println("  -config <path>      Path to user configuration file")
		fmt.Println("  -export <path>      Export default configuration to file")
		fmt.Println("  -output <filename>  Output file name (default: project.md)")
		fmt.Println("  -verbose            Enable verbose output with file details")
		fmt.Println("  -stat               Include file size information and statistics in output")
		fmt.Println("  -version            Show version information")
		fmt.Println("  -no-git             Do not use .gitignore for exclude files")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Printf("  %s ./my-project\n", exeFile)
		fmt.Printf("  %s -config config.json ./my-project\n", exeFile)
		fmt.Printf("  %s -export default-config.json\n", exeFile)
		fmt.Printf("  %s -verbose -stat -config -no-git config.json -output ./my-project/project.md ./my-project\n", exeFile)
		os.Exit(1)
	}

	projectPath := args[0]

	// Validate project directory
	if stat, err := os.Stat(projectPath); os.IsNotExist(err) {
		log.Fatalf("Error: directory '%s' does not exist", projectPath)
	} else if err != nil {
		log.Fatalf("Error: cannot access directory '%s': %v", projectPath, err)
	} else if !stat.IsDir() {
		projectPath = filepath.Dir(projectPath)
	}

	// Load configuration
	config := getDefaultConfig()

	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		log.Fatalf("failed to get absolute path: %v", err)
	}

	var customConfig = *NewConfig()
	projectConfigPath := filepath.Join(absPath, "project2md.config.json")
	exists, err := pathExists(projectConfigPath)
	if err != nil {
		log.Fatalf("Error using project configuration: %v", err)
	}
	if exists {
		fmt.Printf("Using project configuration: %s\n", projectConfigPath)
		customConfig, err = loadCustomConfig(customConfig, projectConfigPath)
	}

	if *userConfigPath != "" {
		customConfig, err = loadCustomConfig(customConfig, *userConfigPath)
		fmt.Printf("Using custom configuration: %s\n", *userConfigPath)
	}

	// Process project
	if err := processProject(absPath, config, customConfig, *outputFileName, *verbose, *showStats, *noGit); err != nil {
		log.Fatalf("Error processing project: %v", err)
	}
}
```

=== processor.go ===

```go
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
```

=== utils.go ===

```go
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
```

