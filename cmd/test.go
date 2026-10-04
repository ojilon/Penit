package cmd

import (
	"fmt"
	"strings"

	"github.com/ojilon/penit/internal/build"
	"github.com/ojilon/penit/internal/storage"
	"github.com/spf13/cobra"
)

var testDevice bool

var testCmd = &cobra.Command{
	Use:   "test <project>",
	Short: "Run unit tests (testDebugUnitTest)",
	Long: `Run Android unit tests for a configured project.

Streams Gradle output live and writes a full log under the data-root.
On failure, prints structured headlines.

Examples:
  penit test wayer
  penit test conductino_android`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		if testDevice {
			return fmt.Errorf("--device instrumented tests are not implemented yet")
		}

		cfg, _, err := loadConfig()
		if err != nil {
			return err
		}
		proj, err := cfg.ResolveProject(key)
		if err != nil {
			return err
		}
		if err := storage.EnsureProject(cfg.DataRoot, key); err != nil {
			return err
		}

		task := build.DefaultTestTask
		logPath := storage.NewBuildLogPath(cfg.DataRoot, key, "test")

		fmt.Println()
		fmt.Printf("Testing %s · %s\n", key, task)
		fmt.Printf("  project : %s\n", proj.Path)
		fmt.Printf("  log     : %s\n", logPath)
		fmt.Println(strings.Repeat("─", 60))

		res, err := build.Run(build.Options{
			ProjectRoot:    proj.Path,
			ProjectKey:     key,
			Task:           task,
			GradleUserHome: proj.GradleUserHome,
			LogPath:        logPath,
		})
		if err != nil {
			return err
		}

		fmt.Println(strings.Repeat("─", 60))
		if res.ExitCode != 0 {
			fmt.Println()
			fmt.Println(strings.Repeat("═", 40))
			fmt.Printf("TESTS FAILED  ·  %s\n", key)
			fmt.Println(strings.Repeat("═", 40))
			if len(res.Headlines) > 0 {
				for _, h := range res.Headlines {
					fmt.Println(h)
				}
			} else {
				fmt.Println("(no structured headlines extracted — see full log)")
			}
			fmt.Println()
			if res.LogPath != "" {
				fmt.Printf("Full log : %s\n", res.LogPath)
			}
			return fmt.Errorf("tests failed (exit %d)", res.ExitCode)
		}

		fmt.Printf("✓ %s finished · %s\n", res.Task, key)
		if res.LogPath != "" {
			fmt.Printf("  log: %s\n", res.LogPath)
		}
		return nil
	},
}

func init() {
	testCmd.Flags().BoolVar(&testDevice, "device", false, "run instrumented device tests (not yet implemented)")
	rootCmd.AddCommand(testCmd)
}
