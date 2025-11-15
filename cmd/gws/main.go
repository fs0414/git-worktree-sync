package main

import (
	"fmt"
	"os"

	"github.com/fs0414/git-worktree-sync/internal/cli"
	"github.com/spf13/cobra"
)

var (
	version = "1.0.0"
	commit  = "dev"
	date    = "unknown"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "gws",
		Short: "Git worktree resource synchronization tool",
		Long: `gws (git-worktree-sync) is a tool for managing git worktrees with automatic resource synchronization.

It helps you create new worktrees and automatically sync resources like node_modules, .env files,
and other dependencies from your main worktree.`,
		Version: fmt.Sprintf("%s (commit: %s, built at: %s)", version, commit, date),
	}

	// Add subcommands
	rootCmd.AddCommand(cli.CreateCmd())
	rootCmd.AddCommand(cli.InitCmd())
	rootCmd.AddCommand(cli.SyncCmd())
	rootCmd.AddCommand(cli.ListCmd())

	// Execute
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
