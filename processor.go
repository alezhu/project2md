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

type Processor struct {
	projectPath    string
	defaultConfig  Config
	customConfig   Config
	outputFileName string
	verbose        bool
	showStats      bool
	noGit          bool
	gitIgnore      *GitIgnore
	stats          *Statistics
	files          []string
	writer         *bufio.Writer
}

func NewProcessor(
	projectPath string,
	defaultConfig Config,
	customConfig Config,
	outputFileName string,
	verbose bool,
	showStats bool,
	noGit bool,
) *Processor {
	return &Processor{
		projectPath:    projectPath,
		defaultConfig:  defaultConfig,
		customConfig:   customConfig,
		outputFileName: outputFileName,
		verbose:        verbose,
		showStats:      showStats,
		noGit:          noGit,
	}
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
	} else {
		result, exists = config.CodeExtensions[ext]
	}

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
func (p *Processor) Process() error {
	var outputFile string
	if filepath.IsAbs(p.outputFileName) {
		outputFile = p.outputFileName
	} else {
		outputFile = filepath.Clean(filepath.Join(p.projectPath, p.outputFileName))
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	p.writer = bufio.NewWriter(file)
	defer func() {
		_ = p.writer.Flush()
	}()

	// Load .gitignore patterns
	err = p.loadGitIgnore()
	if err != nil {
		return err
	}

	p.stats = &Statistics{
		StartTime: time.Now(),
	}

	// Write header
	err = p.writeHeader()
	if err != nil {
		return err
	}

	// Collect all files first for better organization
	err = p.findFiles()
	if err != nil {
		return fmt.Errorf("error walking directory: %w", err)
	}

	// Sort files for consistent output
	sort.Strings(p.files)

	// Process files
	err2 := p.processFiles()
	if err2 != nil {
		return err2
	}

	// Write statistics only if showStats is true
	if p.showStats {
		err = p.writeStats()
		if err != nil {
			return err
		}
	}

	fmt.Printf("\nArchive created: %s\n", outputFile)
	fmt.Printf("Statistics: %d files, %s", p.stats.ProcessedFiles, formatFileSize(p.stats.TotalSize))
	if p.showStats {
		duration := time.Since(p.stats.StartTime)
		fmt.Printf(", %v processing time", duration.Round(time.Millisecond))
	}
	fmt.Println()

	return nil
}

func (p *Processor) writeStats() error {
	duration := time.Since(p.stats.StartTime)
	if err := writeFileContent(p.writer, "---\n\n"); err != nil {
		return fmt.Errorf("failed to write statistics separator: %w", err)
	}

	if err := writeFileContent(p.writer, "## Statistics\n\n"); err != nil {
		return fmt.Errorf("failed to write statistics header: %w", err)
	}

	if err := writeFileContent(p.writer, "- **Files processed**: %d\n", p.stats.ProcessedFiles); err != nil {
		return fmt.Errorf("failed to write file count: %w", err)
	}

	if err := writeFileContent(p.writer, "- **Directories skipped**: %d\n", p.stats.SkippedDirs); err != nil {
		return fmt.Errorf("failed to write dir count: %w", err)
	}

	if err := writeFileContent(p.writer, "- **Total size**: %s\n", formatFileSize(p.stats.TotalSize)); err != nil {
		return fmt.Errorf("failed to write total size: %w", err)
	}

	if err := writeFileContent(p.writer, "- **Processing time**: %v\n", duration.Round(time.Millisecond)); err != nil {
		return fmt.Errorf("failed to write duration: %w", err)
	}
	return nil
}

func (p *Processor) processFiles() error {
	for _, path := range p.files {
		info, err := os.Stat(path)
		if err != nil {
			if p.verbose {
				log.Printf("Warning: cannot stat file %s: %v", path, err)
			}
			continue
		}

		// Get relative path
		relPath, err := filepath.Rel(p.projectPath, path)
		if err != nil {
			if p.verbose {
				log.Printf("Warning: cannot get relative path for %s: %v", path, err)
			}
			continue
		}

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			if p.verbose {
				log.Printf("Warning: cannot read file %s: %v", relPath, err)
			}
			continue
		}

		// Get language for syntax highlighting
		language := getLanguage(info.Name(), p.defaultConfig)

		// Write file section
		err2 := p.writeFileSection(relPath, info, language, content)
		if err2 != nil {
			return err2
		}

		p.stats.ProcessedFiles++
		p.stats.TotalSize += info.Size()

		if p.verbose {
			fmt.Printf("Processed: %s (%s)\n", relPath, formatFileSize(info.Size()))
		}
	}
	return nil
}

func (p *Processor) writeFileSection(relPath string, info os.FileInfo, language string, content []byte) error {
	if err := writeFileContent(p.writer, "=== %s ===\n\n", relPath); err != nil {
		return fmt.Errorf("failed to write file header: %w", err)
	}

	if p.showStats {
		if err := writeFileContent(
			p.writer,
			"*Size: %s, Modified: %s*\n\n",
			formatFileSize(info.Size()),
			info.ModTime().Format("2006-01-02 15:04:05"),
		); err != nil {
			return fmt.Errorf("failed to write file info: %w", err)
		}
	}

	if err := writeFileContent(p.writer, "```%s\n", language); err != nil {
		return fmt.Errorf("failed to write code block start: %w", err)
	}

	if err := writeFileContent(p.writer, "%s", string(content)); err != nil {
		return fmt.Errorf("failed to write file content: %w", err)
	}

	if !strings.HasSuffix(string(content), "\n") {
		if err := writeFileContent(p.writer, "\n"); err != nil {
			return fmt.Errorf("failed to write newline: %w", err)
		}
	}

	if err := writeFileContent(p.writer, "```\n\n"); err != nil {
		return fmt.Errorf("failed to write code block end: %w", err)
	}
	return nil
}

func (p *Processor) writeHeader() error {
	if err := writeFileContent(p.writer, "# Code Archive: %s\n\n", filepath.Base(p.projectPath)); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	if err := writeFileContent(p.writer, "Generated automatically from: `%s`\n", p.projectPath); err != nil {
		return fmt.Errorf("failed to write path: %w", err)
	}

	if err := writeFileContent(p.writer, "Generated at: %s\n\n", time.Now().Format("2006-01-02 15:04:05")); err != nil {
		return fmt.Errorf("failed to write timestamp: %w", err)
	}

	if err := writeFileContent(p.writer, "---\n\n"); err != nil {
		return fmt.Errorf("failed to write separator: %w", err)
	}
	return nil
}

func (p *Processor) loadGitIgnore() error {
	if p.noGit {
		p.gitIgnore = &GitIgnore{}
	} else {
		var err error
		p.gitIgnore, err = loadGitIgnore(p.projectPath)
		if err != nil {
			return fmt.Errorf("failed to write header: %w", err)
		}
	}
	return nil
}

func (p *Processor) findFiles() error {
	clear(p.files)
	err := filepath.Walk(p.projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if p.verbose {
				log.Printf("Warning: error accessing %s: %v", path, err)
			}
			return nil
		}

		// Skip directories in exclude list
		if info.IsDir() {
			return p.checkDirAllowed(path, info)
		}

		// Check if file should be processed
		return p.checkFileShouldBeProcessed(path, info)
	})
	return err
}

func (p *Processor) checkFileShouldBeProcessed(path string, info os.FileInfo) error {
	filePath := info.Name()
	relPath, _ := filepath.Rel(p.projectPath, path)
	shouldProcess := shouldProcessFile(relPath, p.customConfig)
	switch {
	case shouldProcess == -1:
		return nil
	case shouldProcess == 0:
		isCode, exists := isCodeFile(filePath, p.customConfig)
		if !isCode && exists {
			//If file explicit excluded in custom config, then skip file
			if p.verbose {
				fmt.Printf("Ignored by custom config: %s\n", filePath)
			}
			return nil
		}
		if !p.noGit && p.gitIgnore.isMatchPattern(filePath, false) {
			if p.verbose {
				fmt.Printf("Ignored by .gitignore: %s\n", filePath)
			}
			return nil
		}
		shouldProcess = shouldProcessFile(relPath, p.defaultConfig)
		if shouldProcess == -1 {
			if p.verbose {
				fmt.Printf("Ignored by default config: %s\n", filePath)
			}
			return nil
		} else if shouldProcess == 0 {
			if !isCode {
				isCode, _ = isCodeFile(filePath, p.defaultConfig)
			}
			if isCode {
				shouldProcess = 1
			}
		}
	}

	if shouldProcess == 1 {
		p.files = append(p.files, path)
	}
	return nil
}

func (p *Processor) checkDirAllowed(path string, info os.FileInfo) error {
	//Check customConfig first
	skip, exists := shouldSkipDir(info.Name(), p.customConfig)

	if skip {
		//If dir is explicitly skipped in custom config then skip this dir
		if p.verbose {
			fmt.Printf("Ignored by custom config: %s\n", path)
		}
		return filepath.SkipDir
	} else {
		if exists {
			//If dir is explicitly excluded from skip-dir-list in custom config then Process this dir
		} else {
			// Check .gitignore
			if !p.noGit && p.gitIgnore.isMatchPattern(path, true) {
				if p.verbose {
					fmt.Printf("Ignored by .gitignore: %s\n", path)
				}
				p.stats.SkippedDirs++
				return filepath.SkipDir
			}
			//Check defaultConfig
			if skip, _ = shouldSkipDir(info.Name(), p.defaultConfig); skip {
				if p.verbose {
					fmt.Printf("Skipping directory: %s\n", path)
				}
				p.stats.SkippedDirs++
				return filepath.SkipDir
			}
		}
	}
	return nil
}
