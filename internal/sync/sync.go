package sync

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/fs0414/git-worktree-sync/internal/config"
)

// SyncMode represents how resources should be synced
type SyncMode int

const (
	SyncModeSymlink SyncMode = iota
	SyncModeCopy
)

// SyncResult represents the result of a sync operation
type SyncResult struct {
	Resource string
	Mode     string
	Success  bool
	Error    error
}

// SyncResources synchronizes resources from source to destination based on config
func SyncResources(cfg *config.Config, sourceDir, destDir string, forceCopy bool) ([]SyncResult, error) {
	var results []SyncResult

	// Sync symlink resources
	if !forceCopy {
		for _, resource := range cfg.Resources.Symlink {
			result := syncResource(resource, sourceDir, destDir, SyncModeSymlink)
			results = append(results, result)
		}
	} else {
		// If force copy, treat symlink resources as copy
		for _, resource := range cfg.Resources.Symlink {
			result := syncResource(resource, sourceDir, destDir, SyncModeCopy)
			results = append(results, result)
		}
	}

	// Sync copy resources
	for _, resource := range cfg.Resources.Copy {
		result := syncResource(resource, sourceDir, destDir, SyncModeCopy)
		results = append(results, result)
	}

	return results, nil
}

func syncResource(resource, sourceDir, destDir string, mode SyncMode) SyncResult {
	sourcePath := filepath.Join(sourceDir, resource)
	destPath := filepath.Join(destDir, resource)

	result := SyncResult{
		Resource: resource,
		Success:  false,
	}

	// Check if source exists
	sourceInfo, err := os.Lstat(sourcePath)
	if err != nil {
		if os.IsNotExist(err) {
			result.Error = fmt.Errorf("source does not exist: %s", resource)
			result.Mode = "skip"
			return result
		}
		result.Error = fmt.Errorf("failed to stat source: %w", err)
		result.Mode = "error"
		return result
	}

	// Check if destination already exists
	if _, err := os.Lstat(destPath); err == nil {
		// Destination exists, skip
		result.Success = true
		result.Mode = "exists"
		return result
	}

	// Ensure parent directory exists
	parentDir := filepath.Dir(destPath)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		result.Error = fmt.Errorf("failed to create parent directory: %w", err)
		result.Mode = "error"
		return result
	}

	// Perform sync based on mode
	if mode == SyncModeSymlink {
		result.Mode = "symlink"
		if err := createSymlink(sourcePath, destPath); err != nil {
			result.Error = err
			return result
		}
	} else {
		result.Mode = "copy"
		if sourceInfo.IsDir() {
			if err := copyDir(sourcePath, destPath); err != nil {
				result.Error = err
				return result
			}
		} else {
			if err := copyFile(sourcePath, destPath); err != nil {
				result.Error = err
				return result
			}
		}
	}

	result.Success = true
	return result
}

func createSymlink(source, dest string) error {
	// Use absolute path for symlink
	absSource, err := filepath.Abs(source)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	if err := os.Symlink(absSource, dest); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	return nil
}

func copyFile(source, dest string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	// Copy file permissions
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return fmt.Errorf("failed to stat source file: %w", err)
	}

	if err := os.Chmod(dest, sourceInfo.Mode()); err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	return nil
}

func copyDir(source, dest string) error {
	// Get source directory info
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return fmt.Errorf("failed to stat source directory: %w", err)
	}

	// Create destination directory
	if err := os.MkdirAll(dest, sourceInfo.Mode()); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Read directory entries
	entries, err := os.ReadDir(source)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	// Copy each entry
	for _, entry := range entries {
		sourcePath := filepath.Join(source, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			if err := copyDir(sourcePath, destPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(sourcePath, destPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// CheckSyncStatus checks if resources are synced in the destination
func CheckSyncStatus(cfg *config.Config, sourceDir, destDir string) (bool, error) {
	// Check all resources
	allResources := append(cfg.Resources.Symlink, cfg.Resources.Copy...)

	for _, resource := range allResources {
		destPath := filepath.Join(destDir, resource)
		if _, err := os.Lstat(destPath); os.IsNotExist(err) {
			return false, nil
		}
	}

	return true, nil
}
