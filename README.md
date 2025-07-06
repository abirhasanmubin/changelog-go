# Changelog Generator

An interactive command-line tool for generating structured changelog entries in Markdown format with Git integration.

## Features

- Interactive prompts for changelog creation
- Git integration (branch, commits, user info)
- Markdown output with structured sections
- Customizable change types
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
├── changelog/     # Core changelog logic
├── command/       # Git command execution
├── input/         # User input handling
├── prompt/        # Interactive prompts
└── main.go        # CLI entry point
```

## License

MIT License - see [LICENSE](LICENSE) file for details.
