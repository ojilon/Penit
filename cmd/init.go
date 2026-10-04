package cmd

import (
	"fmt"

	"github.com/ojilon/penit/internal/config"
	"github.com/ojilon/penit/internal/storage"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create data-root and default config.json",
	Long: `Creates the Penit data-root directory and writes a default config.json
with placeholder entries for wayer, conductino_android, fdroid, and cooda.

Edit the paths to point at your local clones, then run: penit doctor`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Default()
		if dataRootFlag != "" {
			cfg.DataRoot = dataRootFlag
		}

		if err := storage.EnsureRoot(cfg.DataRoot); err != nil {
			return err
		}

		cfgPath := configFlag
		if cfgPath == "" {
			cfgPath = config.Path(cfg.DataRoot)
		}

		// Do not overwrite an existing config unless forced later.
		if _, err := config.Load(cfgPath); err == nil {
			// File may already exist with content; still ensure root.
			existing, loadErr := config.Load(cfgPath)
			if loadErr == nil && existing.DataRoot != "" {
				fmt.Printf("Config already exists: %s\n", cfgPath)
				fmt.Printf("Data root: %s\n", existing.DataRoot)
				fmt.Println("Edit project paths, then run: penit doctor")
				return nil
			}
		}

		if err := config.Save(cfgPath, cfg); err != nil {
			return err
		}

		fmt.Println("Penit initialised")
		fmt.Printf("  data-root : %s\n", cfg.DataRoot)
		fmt.Printf("  config    : %s\n", cfgPath)
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Println("  1. Edit config.json and set absolute paths for your projects")
		fmt.Println("  2. penit doctor")
		fmt.Println("  3. penit build wayer debug")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
