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
	defer func() {
		_ = file.Close()
	}()

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
		if err != nil {
			log.Fatalf("Error using project configuration: %v", err)
		}
	}

	if *userConfigPath != "" {
		fmt.Printf("Using custom configuration: %s\n", *userConfigPath)
		customConfig, err = loadCustomConfig(customConfig, *userConfigPath)
		if err != nil {
			log.Fatalf("Error using custom configuration: %v", err)
		}
	}

	// Process project
	processor := NewProcessor(absPath, config, customConfig, *outputFileName, *verbose, *showStats, *noGit)
	if err := processor.Process(); err != nil {
		log.Fatalf("Error processing project: %v", err)
	}
}
