# Changelog Generator

An interactive command-line tool for generating structured changelog entries in Markdown format with Git integration.

## Features

- Interactive prompts with colorful UI
- Intuitive navigation (arrow keys, vim-style keys)
- Git integration (branch, commits, user info)
- Multiple output formats:
  - Generate markdown file
  - Copy Bitbucket PR text to clipboard
  - Display Bitbucket PR text
- Customizable change types with validation
- Inline yes/no selection
- Checklist for PR readiness
- Works from any directory

## Installation

```bash
go install github.com/abirhasanmubin/changelog-go@latest
```

Or clone and build:

```bash
git clone https://github.com/abirhasanmubin/changelog-go.git
cd changelog-go
go build -o changelog-go
```

## Usage

### As a CLI tool

```bash
changelog-go
```

### Navigation Controls

**Multi-select options:**
- ↑/↓ or j/k: Navigate up/down
- SPACE: Toggle selection
- a: Toggle all options
- ENTER: Confirm selection

**Single-select options:**
- ↑/↓ or j/k: Navigate up/down
- ENTER: Confirm selection

**Yes/No questions:**
- ←/→ or h/l: Navigate left/right
- ENTER: Confirm selection

**Text input:**
- Type normally for single-line input
- For multi-line input, type "EOF" on a new line to finish

### As a Go package

```go
import "github.com/abirhasanmubin/changelog-go/prompt"

func main() {
    prompt.Generate()
}
```

## Output Formats

### 1. Generate File
Creates a structured changelog in `.logs/.changelog/` with sections for:
- Title and description
- Type of changes (bug fix, feature, etc.)
- Motivation and implementation details
- To-do items before merge
- Model changes
- Testing instructions
- PR checklist
- Git commit history

### 2. Bitbucket PR Format
Generates optimized content for Bitbucket pull requests:
- **Copy to clipboard**: Automatically copies PR description to system clipboard
- **Display text**: Shows formatted PR content for manual copying
- Compact format with checkmarks (✅) for selected change types
- Bold section headers for better readability
- Excludes commit history (handled by Bitbucket)

## UI Features

- **Colorful interface** with syntax highlighting
- **Multi-select options** with ✓ checkmarks
- **Inline yes/no selection** using ←/→ arrows or h/l keys
- **Question prefixes** with ? symbols for clarity
- **Validation** requiring at least one change type
- **Visual feedback** with error messages and success indicators

## Example Output

### File Format
```markdown
## Title

Fix user authentication bug

## Type of change (Check all that apply)

- [x] Bug fix
- [ ] New feature
- [ ] Code refactor

## Checklist

- [x] I have performed a self-review of my code
- [x] I have added tests that prove my fix is effective
```

### Bitbucket PR Format
```markdown
Fix user authentication bug in login flow

**Motivation:**
Users were experiencing login failures due to token validation issues

**Type of change:**
- ✅ Bug fix

**Checklist:**
- [x] I have performed a self-review of my code
- [x] I have added tests that prove my fix is effective
```

## Development

### Running Tests

```bash
# Unit tests
go test ./...

# Integration tests
go test -tags=integration ./...

# With coverage
go test ./... -cover
```

### Project Structure

```
├── .logs/         # Generated changelog output
├── changelog/     # Core changelog logic
├── command/       # Git command execution
├── input/         # User input handling with validation
├── prompt/        # Interactive prompts with colors
├── ui/            # User interface components
│   ├── base.go        # Shared UI functionality
│   ├── colors.go      # Color constants
│   ├── boolean.go     # Inline yes/no selection
│   ├── multiselect.go # Multi-option selection with validation
│   ├── singleselect.go # Single option selection
│   └── terminal.go    # Terminal control
├── utils/         # Utility functions
│   └── clipboard.go   # Clipboard operations
├── main.go        # CLI entry point
├── go.mod         # Go module definition
└── *_test.go      # Test files
```

### Recent Optimizations

- **Modular Architecture**: Extracted common UI functionality into base components
- **Reduced Code Duplication**: Consolidated similar functions across packages
- **Improved Error Handling**: Consistent error messages with color coding
- **Streamlined Dependencies**: Updated go.mod with proper versioning
- **Enhanced Maintainability**: Separated concerns and improved code organization

## License

MIT License - see [LICENSE](LICENSE) file for details.
