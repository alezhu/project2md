# Project2MD

A command-line tool that automatically generates markdown archives of source code projects. Perfect for documentation, code reviews, sharing project structures, or creating backups of your codebase in a readable format.

## Features

- 🚀 **Fast and Efficient**: Written in Go for optimal performance
- 📁 **Smart File Detection**: Automatically identifies code files by extension
- 🎯 **Configurable Filtering**: Extensive configuration options for include/exclude patterns
- 🔍 **Git Integration**: Respects `.gitignore` files automatically
- 📊 **Statistics**: Optional file size and processing statistics
- 🎨 **Syntax Highlighting**: Proper language detection for Markdown code blocks
- ⚙️ **Flexible Configuration**: JSON-based configuration with project-level overrides
- 📝 **Clean Output**: Well-formatted markdown with file organization

## Installation

### From Source

```bash
git clone https://github.com/alezhu/project2md.git
cd project2md
go build -o project2md
```

### Binary Release

Download the latest binary from the [releases page](https://github.com/alezhu/project2md/releases).

## Usage

```
project2md [options] <project_directory>

Options:
  -config <path>      Path to user configuration file
  -export <path>      Only export default configuration to file (<project_directory> ignored) 
  -output <filename>  Output file name (default: project.md)
  -verbose            Enable verbose output with file details
  -stat               Include file size information and statistics in output
  -version            Show version information
  -no-git             Do not use .gitignore for exclude files
```

### Examples

```bash
# Basic project archiving
./project2md ./my-project

# With custom configuration
./project2md -config custom-config.json ./my-project

# Verbose mode with statistics
./project2md -verbose -stat ./my-project

# Custom output location
./project2md -output ./archive/project-backup.md ./my-project

# Ignore .gitignore files
./project2md -no-git ./my-project
```

## Configuration

Project2MD uses a flexible JSON configuration system with three levels of precedence:

1. **Default Configuration**: Built-in sensible defaults
2. **Project Configuration**: `project2md.config.json` in the project root
3. **User Configuration**: Custom config file specified with `-config`

### Supported File Types

The tool automatically recognizes these file types:

**Programming Languages**: `.go`, `.py`, `.js`, `.ts`, `.java`, `.c`, `.cpp`, `.cs`, `.rs`, `.rb`, `.php`, `.swift`, `.kt`, `.scala`

**Web Technologies**: `.html`, `.css`, `.scss`, `.sass`, `.less`, `.jsx`, `.tsx`

**Configuration**: `.json`, `.yaml`, `.yml`, `.toml`, `.ini`, `.cfg`, `.conf`, `.env`

**Scripts**: `.sh`, `.bash`, `.zsh`, `.fish`, `.bat`, `.cmd`, `.ps1`

**Documentation**: `.md`, `.txt`, `.xml`

**Build Files**: `Makefile`, `Dockerfile`

### Configuration File Format

Export the default configuration to see all available options:

```bash
./project2md -export config.json
```

Example configuration:

```json
{
  "code_extensions": {
    ".go": true,
    ".py": true,
    ".js": true,
    ".md": false,
    ".log": false
  },
  "skip_dirs": {
    "node_modules": true,
    ".git": true,
    "dist": true,
    "build": true
  },
  "languages": {
    ".go": "go",
    ".py": "python",
    ".js": "javascript"
  },
  "exclude": [
    "*.log",
    "package-lock.json",
    "yarn.lock"
  ],
  "include": [
    "special-file.txt"
  ]
}
```

### Configuration Properties

- **`code_extensions`**: Map of file extensions to include/exclude (`true`/`false`)
- **`skip_dirs`**: Directories to skip during processing
- **`languages`**: File extension to syntax highlighting language mapping
- **`exclude`**: File patterns to exclude (supports glob patterns)
- **`include`**: File patterns to force include (takes precedence over exclude)

## Output Format

The generated markdown file includes:

- **Header**: Project name, source path, and generation timestamp
- **File Sections**: Each file with syntax-highlighted code blocks
- **Statistics** (optional): Processing statistics and file information

Example output structure:

```markdown
# Code Archive: my-project

Generated automatically from: `/path/to/my-project`
Generated at: 2025-06-29 16:53:08

---

=== src/main.go ===

*Size: 2.5 KB, Modified: 2025-06-29 15:30:22*

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Changelog

### v1.0.0
- Initial release
- Core functionality for project archiving
- JSON configuration system
- Git integration
- Statistics and verbose output
- Cross-platform support
