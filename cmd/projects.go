package cmd

import (
	"fmt"
	"os"

	"github.com/ojilon/penit/internal/config"
	"github.com/spf13/cobra"
)

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "List or manage known Android projects",
}

var projectsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _, err := loadConfig()
		if err != nil {
			return err
		}
		if len(cfg.Projects) == 0 {
			fmt.Println("No projects configured.")
			return nil
		}
		fmt.Printf("%-22s %-28s %s\n", "KEY", "APP NAME", "PATH")
		for _, key := range cfg.ProjectKeys() {
			p := cfg.Projects[key]
			path := p.Path
			if path == "" {
				path = "(unset)"
			}
			fmt.Printf("%-22s %-28s %s\n", key, p.AppName, path)
		}
		return nil
	},
}

var (
	addPath    string
	addAppName string
)

var projectsAddCmd = &cobra.Command{
	Use:   "add <key>",
	Short: "Add or update a project entry",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		cfg, cfgPath, err := loadConfig()
		if err != nil {
			return err
		}
		p, exists := cfg.Projects[key]
		if !exists {
			p = config.Project{}
		}
		if addPath != "" {
			p.Path = addPath
		}
		if addAppName != "" {
			p.AppName = addAppName
		}
		if p.AppName == "" {
			p.AppName = key
		}
		cfg.Projects[key] = p
		if err := config.Save(cfgPath, cfg); err != nil {
			return err
		}
		fmt.Printf("Saved project %q → %s\n", key, cfgPath)
		if p.Path == "" {
			fmt.Println("Warning: path still empty — set with --path")
		}
		return nil
	},
}

func init() {
	projectsAddCmd.Flags().StringVar(&addPath, "path", "", "absolute path to the project root")
	projectsAddCmd.Flags().StringVar(&addAppName, "app-name", "", "display / artifact name")
	projectsCmd.AddCommand(projectsListCmd, projectsAddCmd)
	rootCmd.AddCommand(projectsCmd)
}

// ensure unused import does not break if os is needed later
var _ = os.DevNull
