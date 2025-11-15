package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	ConfigFileName       = ".gwt.yml"
	GlobalConfigDir      = ".config/gws"
	GlobalConfigFileName = "config.yml"
)

// Config represents the .gwt.yml configuration file
type Config struct {
	Resources     Resources         `yaml:"resources"`
	WorktreePath  string            `yaml:"worktree_path"`
	Exclude       []string          `yaml:"exclude"`
	Hooks         map[string]string `yaml:"hooks,omitempty"`
}

// Resources defines which resources to sync
type Resources struct {
	Symlink []string `yaml:"symlink"`
	Copy    []string `yaml:"copy"`
}

// Load loads the configuration from .gwt.yml in the given directory
func Load(dir string) (*Config, error) {
	configPath := filepath.Join(dir, ConfigFileName)
	return LoadFromPath(configPath)
}

// LoadFromPath loads configuration from a specific path
func LoadFromPath(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set default worktree path if not specified
	if cfg.WorktreePath == "" {
		cfg.WorktreePath = "../{branch}"
	}

	return &cfg, nil
}

// Save saves the configuration to .gwt.yml in the given directory
func (c *Config) Save(dir string) error {
	configPath := filepath.Join(dir, ConfigFileName)
	return c.SaveToPath(configPath)
}

// SaveToPath saves configuration to a specific path
func (c *Config) SaveToPath(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Exists checks if a config file exists in the given directory
func Exists(dir string) bool {
	configPath := filepath.Join(dir, ConfigFileName)
	_, err := os.Stat(configPath)
	return err == nil
}

// ResolveWorktreePath resolves the worktree path with branch name substitution
func (c *Config) ResolveWorktreePath(branchName string) string {
	return strings.ReplaceAll(c.WorktreePath, "{branch}", branchName)
}

// LoadGlobalConfig loads the global configuration from ~/.config/gws/config.yml
func LoadGlobalConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, GlobalConfigDir, GlobalConfigFileName)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Return default config if global config doesn't exist
		return GetDefaultConfig(), nil
	}

	return LoadFromPath(configPath)
}

// GetDefaultConfig returns a default configuration
func GetDefaultConfig() *Config {
	return &Config{
		Resources: Resources{
			Symlink: []string{"node_modules"},
			Copy:    []string{".env"},
		},
		WorktreePath: "../{branch}",
		Exclude:      []string{"*.log", "tmp/*"},
	}
}
