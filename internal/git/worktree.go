package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Worktree represents a git worktree
type Worktree struct {
	Path   string
	Branch string
	IsMain bool
}

// IsGitRepository checks if the current directory is a git repository
func IsGitRepository(dir string) bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = dir
	err := cmd.Run()
	return err == nil
}

// GetMainWorktreePath returns the path to the main worktree
func GetMainWorktreePath(currentDir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = currentDir
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get main worktree path: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// CreateWorktree creates a new git worktree
func CreateWorktree(branchName, path, baseBranch string) error {
	// Check if path already exists
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("path already exists: %s", path)
	}

	var args []string
	if baseBranch != "" {
		// Create new branch from base branch
		args = []string{"worktree", "add", "-b", branchName, path, baseBranch}
	} else {
		// Create new branch from current HEAD
		args = []string{"worktree", "add", "-b", branchName, path}
	}

	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create worktree: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// ListWorktrees returns a list of all worktrees
func ListWorktrees() ([]Worktree, error) {
	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list worktrees: %w", err)
	}

	return parseWorktreeList(string(output)), nil
}

func parseWorktreeList(output string) []Worktree {
	var worktrees []Worktree
	lines := strings.Split(strings.TrimSpace(output), "\n")

	var current *Worktree
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			if current != nil {
				worktrees = append(worktrees, *current)
				current = nil
			}
			continue
		}

		if strings.HasPrefix(line, "worktree ") {
			if current != nil {
				worktrees = append(worktrees, *current)
			}
			current = &Worktree{
				Path: strings.TrimPrefix(line, "worktree "),
			}
		} else if strings.HasPrefix(line, "branch ") {
			if current != nil {
				branchRef := strings.TrimPrefix(line, "branch ")
				// Extract branch name from refs/heads/branch-name
				parts := strings.Split(branchRef, "/")
				if len(parts) > 0 {
					current.Branch = parts[len(parts)-1]
				}
			}
		} else if strings.HasPrefix(line, "bare") {
			if current != nil {
				current.IsMain = true
			}
		}
	}

	// Add the last worktree if exists
	if current != nil {
		worktrees = append(worktrees, *current)
	}

	// Mark the first worktree as main if no bare repository
	if len(worktrees) > 0 && !worktrees[0].IsMain {
		worktrees[0].IsMain = true
	}

	return worktrees
}

// GetCurrentBranch returns the current branch name
func GetCurrentBranch(dir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// IsWorktree checks if the given directory is a git worktree
func IsWorktree(dir string) (bool, error) {
	// Check if it's a git repository first
	if !IsGitRepository(dir) {
		return false, nil
	}

	// Get the git directory
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	gitDir := strings.TrimSpace(string(output))

	// If git dir is .git/worktrees/*, it's a worktree
	return strings.Contains(gitDir, "worktrees"), nil
}

// GetWorktreeMainPath returns the main worktree path for a given worktree
func GetWorktreeMainPath(worktreeDir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--git-common-dir")
	cmd.Dir = worktreeDir
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git common dir: %w", err)
	}

	gitCommonDir := strings.TrimSpace(string(output))

	// The common dir points to the main .git directory
	// Get the parent of .git to get the main worktree path
	mainPath := filepath.Dir(gitCommonDir)

	return mainPath, nil
}

// BranchExists checks if a branch exists
func BranchExists(branchName string) bool {
	cmd := exec.Command("git", "rev-parse", "--verify", branchName)
	err := cmd.Run()
	return err == nil
}
