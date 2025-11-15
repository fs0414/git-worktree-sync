package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "gws-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test config file
	configContent := `resources:
  symlink:
    - node_modules
    - dist
  copy:
    - .env
worktree_path: "../{branch}"
exclude:
  - "*.log"
`
	configPath := filepath.Join(tmpDir, ConfigFileName)
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	// Load config
	cfg, err := Load(tmpDir)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	// Verify config
	if len(cfg.Resources.Symlink) != 2 {
		t.Errorf("expected 2 symlink resources, got %d", len(cfg.Resources.Symlink))
	}
	if len(cfg.Resources.Copy) != 1 {
		t.Errorf("expected 1 copy resource, got %d", len(cfg.Resources.Copy))
	}
	if cfg.WorktreePath != "../{branch}" {
		t.Errorf("expected worktree_path '../{branch}', got %s", cfg.WorktreePath)
	}
}

func TestResolveWorktreePath(t *testing.T) {
	cfg := &Config{
		WorktreePath: "../{branch}",
	}

	result := cfg.ResolveWorktreePath("feature-branch")
	expected := "../feature-branch"

	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestGetDefaultConfig(t *testing.T) {
	cfg := GetDefaultConfig()

	if cfg == nil {
		t.Fatal("expected config, got nil")
	}

	if len(cfg.Resources.Symlink) == 0 {
		t.Error("expected default symlink resources")
	}

	if cfg.WorktreePath == "" {
		t.Error("expected default worktree path")
	}
}

func TestSave(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "gws-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := GetDefaultConfig()

	// Save config
	if err := cfg.Save(tmpDir); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// Verify file exists
	configPath := filepath.Join(tmpDir, ConfigFileName)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("config file was not created")
	}

	// Load and verify
	loadedCfg, err := Load(tmpDir)
	if err != nil {
		t.Fatalf("failed to load saved config: %v", err)
	}

	if len(loadedCfg.Resources.Symlink) != len(cfg.Resources.Symlink) {
		t.Error("loaded config does not match saved config")
	}
}
