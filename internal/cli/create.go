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

// CreateCmd creates the 'create' command
func CreateCmd() *cobra.Command {
	var (
		path       string
		copyMode   bool
		noSync     bool
		baseBranch string
	)

	cmd := &cobra.Command{
		Use:   "create <branch-name>",
		Short: "Create a new worktree with resource synchronization",
		Long: `Create a new git worktree and synchronize resources from the main worktree.
Resources to sync are defined in .gwt.yml configuration file.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			branchName := args[0]
			return runCreate(branchName, path, baseBranch, copyMode, noSync)
		},
	}

	cmd.Flags().StringVarP(&path, "path", "p", "", "Worktree path (default: ../<branch-name>)")
	cmd.Flags().BoolVarP(&copyMode, "copy", "c", false, "Use copy mode instead of symlink")
	cmd.Flags().BoolVar(&noSync, "no-sync", false, "Skip resource synchronization")
	cmd.Flags().StringVarP(&baseBranch, "base", "b", "", "Base branch for new branch")

	return cmd
}

func runCreate(branchName, path, baseBranch string, copyMode, noSync bool) error {
	// Check if we're in a git repository
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	if !git.IsGitRepository(currentDir) {
		return fmt.Errorf("not a git repository")
	}

	// Check if branch already exists
	if git.BranchExists(branchName) {
		return fmt.Errorf("branch '%s' already exists", branchName)
	}

	// Load config
	var cfg *config.Config
	if config.Exists(currentDir) {
		cfg, err = config.Load(currentDir)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
	} else {
		fmt.Println("⚠️  No .gwt.yml found, using default configuration")
		fmt.Println("   Run 'gws init' to create a configuration file")
		cfg = config.GetDefaultConfig()
	}

	// Resolve worktree path
	var worktreePath string
	if path != "" {
		worktreePath = path
	} else {
		worktreePath = cfg.ResolveWorktreePath(branchName)
	}

	// Make path absolute if it's relative
	if !filepath.IsAbs(worktreePath) {
		worktreePath = filepath.Join(filepath.Dir(currentDir), filepath.Base(worktreePath))
	}

	// Create worktree
	fmt.Printf("Creating worktree at %s...\n", worktreePath)
	if err := git.CreateWorktree(branchName, worktreePath, baseBranch); err != nil {
		return err
	}
	fmt.Println("✓ Created worktree")

	// Sync resources if not disabled
	if !noSync {
		fmt.Println("\nSynchronizing resources...")
		results, err := sync.SyncResources(cfg, currentDir, worktreePath, copyMode)
		if err != nil {
			return fmt.Errorf("failed to sync resources: %w", err)
		}

		// Display sync results
		for _, result := range results {
			if result.Mode == "skip" {
				fmt.Printf("⚠️  Skipped %s (not found in source)\n", result.Resource)
				continue
			}
			if result.Mode == "exists" {
				continue // Skip already existing resources
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
		}
	}

	fmt.Printf("\n✨ Done! Run: cd %s\n", worktreePath)
	return nil
}
