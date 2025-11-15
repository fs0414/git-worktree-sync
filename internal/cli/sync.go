package cli

import (
	"fmt"
	"os"

	"github.com/fs0414/git-worktree-sync/internal/config"
	"github.com/fs0414/git-worktree-sync/internal/git"
	"github.com/fs0414/git-worktree-sync/internal/sync"
	"github.com/spf13/cobra"
)

// SyncCmd creates the 'sync' command
func SyncCmd() *cobra.Command {
	var (
		copyMode bool
		force    bool
	)

	cmd := &cobra.Command{
		Use:   "sync [worktree-path]",
		Short: "Synchronize resources to an existing worktree",
		Long: `Synchronize resources from the main worktree to an existing worktree.
If no path is specified, the current directory is used.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var targetPath string
			if len(args) > 0 {
				targetPath = args[0]
			} else {
				var err error
				targetPath, err = os.Getwd()
				if err != nil {
					return fmt.Errorf("failed to get current directory: %w", err)
				}
			}
			return runSync(targetPath, copyMode, force)
		},
	}

	cmd.Flags().BoolVarP(&copyMode, "copy", "c", false, "Use copy mode instead of symlink")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite existing resources")

	return cmd
}

func runSync(targetPath string, copyMode, force bool) error {
	// Check if target is a git repository
	if !git.IsGitRepository(targetPath) {
		return fmt.Errorf("not a git repository: %s", targetPath)
	}

	// Get main worktree path
	mainPath, err := git.GetWorktreeMainPath(targetPath)
	if err != nil {
		return fmt.Errorf("failed to get main worktree path: %w", err)
	}

	// Load config from main worktree
	var cfg *config.Config
	if config.Exists(mainPath) {
		cfg, err = config.Load(mainPath)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
	} else {
		fmt.Println("⚠️  No .gwt.yml found in main worktree, using default configuration")
		cfg = config.GetDefaultConfig()
	}

	fmt.Printf("🔄 Syncing from main worktree: %s\n", mainPath)

	// Sync resources
	results, err := sync.SyncResources(cfg, mainPath, targetPath, copyMode)
	if err != nil {
		return fmt.Errorf("failed to sync resources: %w", err)
	}

	// Display sync results
	syncCount := 0
	for _, result := range results {
		if result.Mode == "skip" {
			fmt.Printf("⚠️  Skipped %s (not found in source)\n", result.Resource)
			continue
		}
		if result.Mode == "exists" {
			if force {
				// TODO: Implement force overwrite
				fmt.Printf("⚠️  %s already exists (force overwrite not yet implemented)\n", result.Resource)
			}
			continue
		}
		if !result.Success {
			fmt.Printf("✗ Failed to sync %s: %v\n", result.Resource, result.Error)
			continue
		}

		if result.Mode == "symlink" {
			fmt.Printf("✓ Linked %s\n", result.Resource)
		} else if result.Mode == "copy" {
			fmt.Printf("✓ Copied %s\n", result.Resource)
		}
		syncCount++
	}

	if syncCount > 0 {
		fmt.Println("\n✨ Sync complete!")
	} else {
		fmt.Println("\n✨ All resources already synced!")
	}

	return nil
}
