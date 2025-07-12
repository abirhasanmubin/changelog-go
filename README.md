# Changelog Generator

An interactive command-line tool for generating structured changelog entries in Markdown format with Git integration.

## Features

- Interactive prompts with colorful UI
- Intuitive navigation (arrow keys, vim-style keys)
- Git integration (branch, commits, user info)
- Markdown output with structured sections
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

## Generated Output

The tool creates a structured changelog in `.logs/.changelog/` with sections for:

- Title and description
- Type of changes (bug fix, feature, etc.)
- Motivation and implementation details
- To-do items before merge
- Model changes
- Testing instructions
- PR checklist
- Git commit history

## UI Features

- **Colorful interface** with syntax highlighting
- **Multi-select options** with ✓ checkmarks
- **Inline yes/no selection** using ←/→ arrows or h/l keys
- **Question prefixes** with ? symbols for clarity
- **Validation** requiring at least one change type
- **Visual feedback** with error messages and success indicators

## Example Output

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
│   ├── colors.go      # Color constants
│   ├── boolean.go     # Inline yes/no selection
│   ├── multiselect.go # Multi-option selection with validation
│   └── terminal.go    # Terminal control
├── main.go        # CLI entry point
└── *_test.go      # Test files
```

## License

MIT License - see [LICENSE](LICENSE) file for details.
