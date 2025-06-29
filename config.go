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
	defer func() {
		_ = file.Close()
	}()

	var customConfig Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&customConfig); err != nil {
		return config, fmt.Errorf("failed to decode config: %w", err)
	}

	// Apply overrides
	mergeMap(&config.CodeExtensions, customConfig.CodeExtensions)
	mergeMap(&config.SkipDirs, customConfig.SkipDirs)
	mergeMap(&config.Languages, customConfig.Languages)
	mergeMap(&config.Exclude, customConfig.Exclude)
	mergeMap(&config.Include, customConfig.Include)

	return config, nil
}
