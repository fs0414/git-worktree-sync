# git-worktree-sync (gws)

A Git worktree resource synchronization tool that automatically syncs resources like `node_modules`, `.env` files, and other dependencies when creating new worktrees.

## Features

- 🚀 **Automatic Resource Sync**: Automatically sync resources from your main worktree when creating new ones
- 🔗 **Flexible Sync Modes**: Choose between symlink (default) or copy mode
- ⚙️ **Project-Aware**: Auto-detects project type (Node.js, Rails, Go, Rust) and suggests appropriate templates
- 📝 **Simple Configuration**: Use `.gwt.yml` to define which resources to sync
- 📊 **Status Tracking**: View all worktrees and their sync status

## Installation

### Using go install

```bash
go install github.com/fs0414/git-worktree-sync/cmd/gws@latest
```

### From Source

```bash
git clone https://github.com/fs0414/git-worktree-sync.git
cd git-worktree-sync
go build -o gws ./cmd/gws
sudo mv gws /usr/local/bin/
```

## Quick Start

1. **Initialize configuration** in your project:

```bash
gws init
```

This creates a `.gwt.yml` file with default settings based on your project type.

2. **Create a new worktree** with resource sync:

```bash
gws create feature-branch
```

This creates a new worktree and automatically syncs configured resources.

3. **List all worktrees** and their sync status:

```bash
gws list
```

4. **Sync resources** to an existing worktree:

```bash
gws sync /path/to/worktree
```

## Commands

### `gws init`

Initialize a `.gwt.yml` configuration file in your project.

```bash
gws init                    # Auto-detect project type
gws init -t node            # Use Node.js template
gws init -t rails           # Use Rails template
gws init -t go              # Use Go template
gws init -t rust            # Use Rust template
gws init --force            # Overwrite existing config
```

### `gws create <branch-name>`

Create a new git worktree with resource synchronization.

```bash
gws create feature-branch                    # Create in default location
gws create feature-branch -p /path/to/dir    # Specify custom path
gws create feature-branch --copy             # Use copy instead of symlink
gws create feature-branch -b main            # Create from main branch
gws create feature-branch --no-sync          # Skip resource sync
```

### `gws list`

List all worktrees and their sync status.

```bash
gws list              # Show all worktrees
gws list -v           # Verbose mode with details
```

### `gws sync [path]`

Synchronize resources to an existing worktree.

```bash
gws sync                      # Sync current directory
gws sync /path/to/worktree    # Sync specific worktree
gws sync --copy               # Use copy mode
gws sync --force              # Overwrite existing resources
```

## Configuration

Create a `.gwt.yml` file in your project root:

```yaml
# Resources to sync via symlink
resources:
  symlink:
    - node_modules
    - .pnpm-store
    - dist

  # Resources to copy
  copy:
    - .env
    - .env.local
    - .env.development

# Worktree path template ({branch} is replaced with branch name)
worktree_path: "../{branch}"

# Exclude patterns (glob supported)
exclude:
  - "*.log"
  - "tmp/*"
```

## Templates

`gws` includes built-in templates for common project types:

### Node.js
```yaml
resources:
  symlink:
    - node_modules
    - .pnpm-store
    - dist
  copy:
    - .env
    - .env.local
```

### Rails
```yaml
resources:
  symlink:
    - vendor
    - node_modules
    - tmp
  copy:
    - .env
    - config/master.key
```

### Go
```yaml
resources:
  symlink:
    - vendor
  copy:
    - .env
```

### Rust
```yaml
resources:
  symlink:
    - target
  copy:
    - .env
```

## Examples

### Create a feature branch worktree

```bash
# Initialize config if not already done
gws init

# Create new worktree
gws create feature-auth
# ✓ Created worktree at ../feature-auth
# ✓ Linked node_modules
# ✓ Copied .env
# ✨ Done! Run: cd ../feature-auth

cd ../feature-auth
```

### Work on multiple branches simultaneously

```bash
gws create feature-a
gws create feature-b
gws create bugfix-123

gws list
# 📂 Worktrees for repository: my-project
#
#   main                  /Users/user/dev/my-project          (main worktree)
# ✓ feature-a            /Users/user/dev/feature-a           (synced)
# ✓ feature-b            /Users/user/dev/feature-b           (synced)
# ✓ bugfix-123           /Users/user/dev/bugfix-123          (synced)
```

### Sync resources to existing worktree

```bash
# Add new resource to .gwt.yml
echo "    - build" >> .gwt.yml

# Sync to existing worktree
gws sync ../feature-branch
# 🔄 Syncing from main worktree: /path/to/main
# ✓ Linked build
# ✨ Sync complete!
```

## Why gws?

When working with git worktrees, you often need to recreate your development environment (install dependencies, copy config files, etc.) for each worktree. `gws` automates this process by:

1. **Saving time**: No need to manually copy or symlink resources
2. **Reducing disk usage**: Symlink large directories like `node_modules` instead of duplicating them
3. **Maintaining consistency**: Ensure all worktrees have the same resources
4. **Simplifying workflow**: One command to create a fully-configured worktree

## Comparison with Manual Workflow

### Without gws:
```bash
git worktree add ../feature-branch
cd ../feature-branch
npm install  # Time-consuming, uses disk space
cp ../ main/.env .
cp ../main/.env.local .
# ... repeat for each resource
```

### With gws:
```bash
gws create feature-branch
# Done! Everything is synced automatically
```

## Requirements

- Git 2.5+ (for worktree support)
- Go 1.24+ (for building from source)

## Troubleshooting

### Config file not found

If you see "No .gwt.yml found", run `gws init` to create a configuration file.

### Symlink creation failed

On Windows, you may need administrator privileges to create symlinks. Alternatively, use `--copy` flag:

```bash
gws create feature-branch --copy
```

### Resource not found

If a resource specified in `.gwt.yml` doesn't exist in the main worktree, it will be skipped with a warning.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details

## Author

Created by [fs0414](https://github.com/fs0414)

## Links

- [GitHub Repository](https://github.com/fs0414/git-worktree-sync)
- [Issue Tracker](https://github.com/fs0414/git-worktree-sync/issues)
