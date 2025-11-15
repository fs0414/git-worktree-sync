package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectProjectType(t *testing.T) {
	tests := []struct {
		name      string
		files     []string
		expected  ProjectType
	}{
		{
			name:     "Node.js project",
			files:    []string{"package.json"},
			expected: ProjectTypeNode,
		},
		{
			name:     "Go project",
			files:    []string{"go.mod"},
			expected: ProjectTypeGo,
		},
		{
			name:     "Rust project",
			files:    []string{"Cargo.toml"},
			expected: ProjectTypeRust,
		},
		{
			name:     "Rails project",
			files:    []string{"Gemfile", "config/application.rb"},
			expected: ProjectTypeRails,
		},
		{
			name:     "Default project",
			files:    []string{},
			expected: ProjectTypeDefault,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tmpDir, err := os.MkdirTemp("", "gws-test-*")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			// Create test files
			for _, file := range tt.files {
				filePath := filepath.Join(tmpDir, file)
				dir := filepath.Dir(filePath)
				if dir != tmpDir {
					if err := os.MkdirAll(dir, 0755); err != nil {
						t.Fatalf("failed to create directory: %v", err)
					}
				}
				if err := os.WriteFile(filePath, []byte{}, 0644); err != nil {
					t.Fatalf("failed to create file: %v", err)
				}
			}

			// Detect project type
			result := DetectProjectType(tmpDir)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGetTemplate(t *testing.T) {
	templates := []ProjectType{
		ProjectTypeNode,
		ProjectTypeRails,
		ProjectTypeGo,
		ProjectTypeRust,
		ProjectTypeDefault,
	}

	for _, tmpl := range templates {
		t.Run(string(tmpl), func(t *testing.T) {
			cfg := GetTemplate(tmpl)

			if cfg == nil {
				t.Fatal("expected config, got nil")
			}

			if cfg.WorktreePath == "" {
				t.Error("expected worktree path to be set")
			}

			// Each template should have at least some resources
			totalResources := len(cfg.Resources.Symlink) + len(cfg.Resources.Copy)
			if totalResources == 0 {
				t.Error("expected template to have resources")
			}
		})
	}
}
