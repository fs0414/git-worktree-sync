package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fs0414/git-worktree-sync/internal/config"
	"github.com/fs0414/git-worktree-sync/internal/git"
	"github.com/fs0414/git-worktree-sync/internal/sync"
	"github.com/spf13/cobra"
)

// ListCmd creates the 'list' command
func ListCmd() *cobra.Command {
	var verbose bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all worktrees and their sync status",
		Long:  `Display all git worktrees with their paths and synchronization status.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(verbose)
		},
	}

	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show detailed information")

	return cmd
}

func runList(verbose bool) error {
	// Check if we're in a git repository
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	if !git.IsGitRepository(currentDir) {
		return fmt.Errorf("not a git repository")
	}

	// Get repository name
	repoPath, err := git.GetMainWorktreePath(currentDir)
	if err != nil {
		return fmt.Errorf("failed to get repository path: %w", err)
	}
	repoName := filepath.Base(repoPath)

	// List all worktrees
	worktrees, err := git.ListWorktrees()
	if err != nil {
		return fmt.Errorf("failed to list worktrees: %w", err)
	}

	if len(worktrees) == 0 {
		fmt.Println("No worktrees found")
		return nil
	}

	// Load config
	var cfg *config.Config
	if config.Exists(repoPath) {
		cfg, err = config.Load(repoPath)
		if err != nil {
			fmt.Printf("⚠️  Failed to load config: %v\n", err)
			cfg = config.GetDefaultConfig()
		}
	} else {
		cfg = config.GetDefaultConfig()
	}

	fmt.Printf("📂 Worktrees for repository: %s\n\n", repoName)

	var mainPath string
	syncedCount := 0
	notSyncedCount := 0

	for _, wt := range worktrees {
		// Display worktree info
		branch := wt.Branch
		if branch == "" {
			branch = "(detached)"
		}

		var status string
		var syncStatus bool

		if wt.IsMain {
			status = "(main worktree)"
			mainPath = wt.Path
		} else {
			// Check sync status
			syncStatus, err = sync.CheckSyncStatus(cfg, mainPath, wt.Path)
			if err != nil {
				status = "(error checking status)"
			} else if syncStatus {
				status = "(synced)"
				syncedCount++
			} else {
				status = "(not synced)"
				notSyncedCount++
			}
		}

		// Format output
		var icon string
		if wt.IsMain {
			icon = "  "
		} else if syncStatus {
			icon = "✓ "
		} else {
			icon = "✗ "
		}

		// Adjust spacing for alignment
		branchDisplay := fmt.Sprintf("%-20s", branch)
		pathDisplay := fmt.Sprintf("%-40s", wt.Path)

		fmt.Printf("%s%s %s %s\n", icon, branchDisplay, pathDisplay, status)

		// Show verbose info if requested
		if verbose && !wt.IsMain {
			fmt.Printf("   Resources: ")
			if syncStatus {
				fmt.Printf("All synced\n")
			} else {
				fmt.Printf("Missing resources\n")
			}
		}
	}

	// Summary
	fmt.Printf("\nTotal: %d worktrees", len(worktrees))
	if syncedCount > 0 || notSyncedCount > 0 {
		fmt.Printf(" (%d synced, %d not synced)", syncedCount, notSyncedCount)
	}
	fmt.Println()

	if notSyncedCount > 0 {
		fmt.Println("\nRun 'gws sync <path>' to sync unsynced worktrees")
	}

	return nil
}
