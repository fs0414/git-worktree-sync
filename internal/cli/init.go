package cli

import (
	"fmt"
	"os"

	"github.com/fs0414/git-worktree-sync/internal/config"
	"github.com/spf13/cobra"
)

// InitCmd creates the 'init' command
func InitCmd() *cobra.Command {
	var (
		template string
		force    bool
	)

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize .gwt.yml configuration file",
		Long: `Create a .gwt.yml configuration file in the current directory.
The file will be created based on auto-detected project type or specified template.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInit(template, force)
		},
	}

	cmd.Flags().StringVarP(&template, "template", "t", "", "Template to use (node, rails, go, rust, default)")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite existing .gwt.yml")

	return cmd
}

func runInit(templateName string, force bool) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Check if config already exists
	if config.Exists(currentDir) && !force {
		return fmt.Errorf(".gwt.yml already exists. Use --force to overwrite")
	}

	var cfg *config.Config
	var projectType config.ProjectType

	if templateName != "" {
		// Use specified template
		projectType = config.ProjectType(templateName)
		cfg = config.GetTemplate(projectType)
		fmt.Printf("Using template: %s\n", templateName)
	} else {
		// Auto-detect project type
		projectType = config.DetectProjectType(currentDir)
		cfg = config.GetTemplate(projectType)

		if projectType != config.ProjectTypeDefault {
			fmt.Printf("🔍 Detected project type: %s\n", projectType)
		} else {
			fmt.Println("🔍 Using default template")
		}
	}

	// Save config
	if err := cfg.Save(currentDir); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println("✓ Created .gwt.yml")
	fmt.Println("📝 Edit .gwt.yml to customize resource sync settings")

	return nil
}
