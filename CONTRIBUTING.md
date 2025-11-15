# Contributing to git-worktree-sync

Thank you for your interest in contributing to git-worktree-sync! This document provides guidelines and instructions for contributing.

## Development Setup

### Prerequisites

- Go 1.24 or later
- Git 2.5 or later
- Make (optional, but recommended)

### Getting Started

1. **Fork and clone the repository**

```bash
git clone https://github.com/YOUR_USERNAME/git-worktree-sync.git
cd git-worktree-sync
```

2. **Install dependencies**

```bash
go mod download
# or
make dev-init
```

3. **Build the project**

```bash
make build
# or
go build -o gws ./cmd/gws
```

4. **Run tests**

```bash
make test
# or
go test ./...
```

## Project Structure

```
git-worktree-sync/
├── cmd/
│   └── gws/
│       └── main.go           # Entry point
├── internal/
│   ├── config/               # Configuration management
│   │   ├── config.go         # Config file loading/saving
│   │   └── templates.go      # Project templates
│   ├── git/                  # Git worktree operations
│   │   └── worktree.go
│   ├── sync/                 # Resource synchronization
│   │   └── sync.go
│   └── cli/                  # CLI commands
│       ├── create.go
│       ├── init.go
│       ├── sync.go
│       └── list.go
├── .gwt.yml                  # Example config
├── README.md
├── LICENSE
└── Makefile
```

## Making Changes

### Code Style

- Follow standard Go conventions and style guidelines
- Run `go fmt` before committing
- Use meaningful variable and function names
- Add comments for exported functions and complex logic

```bash
# Format code
make fmt
# or
go fmt ./...
```

### Testing

- Write tests for new features
- Ensure existing tests pass
- Aim for at least 70% code coverage

```bash
# Run tests
make test

# View coverage
make coverage
```

### Commits

- Write clear, descriptive commit messages
- Use present tense ("Add feature" not "Added feature")
- Reference issues when applicable

Good commit message examples:
```
Add support for global configuration
Fix symlink creation on Windows
Update documentation for init command
```

## Submitting Changes

1. **Create a new branch**

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/bug-description
```

2. **Make your changes**

- Write code
- Add tests
- Update documentation if needed

3. **Test your changes**

```bash
make test
make build
./build/gws --help  # Verify it works
```

4. **Commit your changes**

```bash
git add .
git commit -m "Brief description of changes"
```

5. **Push to your fork**

```bash
git push origin feature/your-feature-name
```

6. **Create a Pull Request**

- Go to the original repository on GitHub
- Click "New Pull Request"
- Select your branch
- Fill in the PR template with:
  - Description of changes
  - Related issues
  - Testing done

## Development Guidelines

### Adding New Commands

1. Create a new file in `internal/cli/`
2. Implement the command using the Cobra framework
3. Register the command in `cmd/gws/main.go`
4. Add tests
5. Update README.md

Example:
```go
// internal/cli/mycommand.go
package cli

import "github.com/spf13/cobra"

func MyCommandCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "mycommand",
        Short: "Description",
        RunE: func(cmd *cobra.Command, args []string) error {
            return runMyCommand()
        },
    }
    return cmd
}
```

### Adding New Templates

1. Add the project type constant in `internal/config/templates.go`
2. Implement the template function
3. Add to `GetTemplate()` switch statement
4. Add tests
5. Update documentation

### Adding New Config Options

1. Update the `Config` struct in `internal/config/config.go`
2. Handle the new option in relevant commands
3. Update `.gwt.yml` example
4. Update documentation

## Testing Locally

### Manual Testing

Create a test repository:

```bash
# Create test repo
mkdir /tmp/test-repo
cd /tmp/test-repo
git init
echo "test" > README.md
git add README.md
git commit -m "Initial commit"

# Copy your built binary
cp /path/to/git-worktree-sync/build/gws .

# Test commands
./gws init
./gws create test-branch
./gws list
```

### Integration Testing

Test with real projects:

- Node.js project with `node_modules`
- Go project with `vendor`
- Project with `.env` files

## Reporting Issues

When reporting bugs, please include:

- OS and version
- Go version
- Git version
- Steps to reproduce
- Expected behavior
- Actual behavior
- Error messages (if any)

## Feature Requests

We welcome feature requests! Please:

- Check if the feature already exists or is planned
- Describe the use case
- Explain why it would be useful
- Provide examples if possible

## Questions?

- Open an issue for questions about development
- Check existing issues for similar questions
- Be respectful and patient

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

## Thank You!

Your contributions make this project better for everyone. Thank you for taking the time to contribute!
