package config

import (
	"os"
	"path/filepath"
)

// ProjectType represents a project type
type ProjectType string

const (
	ProjectTypeNode    ProjectType = "node"
	ProjectTypeRails   ProjectType = "rails"
	ProjectTypeGo      ProjectType = "go"
	ProjectTypeRust    ProjectType = "rust"
	ProjectTypeDefault ProjectType = "default"
)

// GetTemplate returns a configuration template for the specified project type
func GetTemplate(projectType ProjectType) *Config {
	switch projectType {
	case ProjectTypeNode:
		return getNodeTemplate()
	case ProjectTypeRails:
		return getRailsTemplate()
	case ProjectTypeGo:
		return getGoTemplate()
	case ProjectTypeRust:
		return getRustTemplate()
	default:
		return GetDefaultConfig()
	}
}

func getNodeTemplate() *Config {
	return &Config{
		Resources: Resources{
			Symlink: []string{
				"node_modules",
				".pnpm-store",
				"dist",
			},
			Copy: []string{
				".env",
				".env.local",
			},
		},
		WorktreePath: "../{branch}",
		Exclude:      []string{"*.log", "tmp/*"},
	}
}

func getRailsTemplate() *Config {
	return &Config{
		Resources: Resources{
			Symlink: []string{
				"vendor",
				"node_modules",
				"tmp",
			},
			Copy: []string{
				".env",
				"config/master.key",
			},
		},
		WorktreePath: "../{branch}",
		Exclude:      []string{"*.log"},
	}
}

func getGoTemplate() *Config {
	return &Config{
		Resources: Resources{
			Symlink: []string{
				"vendor",
			},
			Copy: []string{
				".env",
			},
		},
		WorktreePath: "../{branch}",
		Exclude:      []string{"*.log", "tmp/*"},
	}
}

func getRustTemplate() *Config {
	return &Config{
		Resources: Resources{
			Symlink: []string{
				"target",
			},
			Copy: []string{
				".env",
			},
		},
		WorktreePath: "../{branch}",
		Exclude:      []string{"*.log", "tmp/*"},
	}
}

// DetectProjectType attempts to detect the project type based on files in the directory
func DetectProjectType(dir string) ProjectType {
	// Check for Node.js
	if fileExists(filepath.Join(dir, "package.json")) {
		return ProjectTypeNode
	}

	// Check for Rails
	if fileExists(filepath.Join(dir, "Gemfile")) && fileExists(filepath.Join(dir, "config", "application.rb")) {
		return ProjectTypeRails
	}

	// Check for Go
	if fileExists(filepath.Join(dir, "go.mod")) {
		return ProjectTypeGo
	}

	// Check for Rust
	if fileExists(filepath.Join(dir, "Cargo.toml")) {
		return ProjectTypeRust
	}

	return ProjectTypeDefault
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
