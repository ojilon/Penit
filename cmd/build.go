package cmd

import (
	"fmt"
	"strings"

	"github.com/ojilon/penit/internal/build"
	"github.com/ojilon/penit/internal/storage"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build <project> [debug|release]",
	Short: "Run assembleDebug or assembleRelease via gradlew",
	Long: `Build an Android APK for a configured project.

Streams Gradle output live, writes a full log under the data-root,
and on failure prints structured headlines plus expected APK paths.

Examples:
  penit build wayer debug
  penit build wayer release
  penit build conductino_android debug`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		variant := build.VariantDebug
		if len(args) >= 2 {
			switch strings.ToLower(args[1]) {
			case "debug", "d":
				variant = build.VariantDebug
			case "release", "r":
				variant = build.VariantRelease
			default:
				return fmt.Errorf("unknown variant %q — use debug or release", args[1])
			}
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

		task := build.TaskForVariant(variant)
		logPath := storage.NewBuildLogPath(cfg.DataRoot, key, string(variant))

		fmt.Println()
		fmt.Printf("Building %s · %s · %s\n", key, variant, task)
		fmt.Printf("  project : %s\n", proj.Path)
		fmt.Printf("  log     : %s\n", logPath)
		fmt.Println(strings.Repeat("─", 60))

		res, err := build.Run(build.Options{
			ProjectRoot:    proj.Path,
			ProjectKey:     key,
			Task:           task,
			Variant:        variant,
			GradleUserHome: proj.GradleUserHome,
			LogPath:        logPath,
		})
		if err != nil {
			return err
		}

		fmt.Println(strings.Repeat("─", 60))
		if res.ExitCode != 0 {
			printBuildFailure(res, key, string(variant))
			return fmt.Errorf("build failed (exit %d)", res.ExitCode)
		}
		printBuildSuccess(res, key, string(variant))
		return nil
	},
}

func printBuildSuccess(res *build.Result, key, variant string) {
	fmt.Printf("✓ %s finished · %s · %s\n", res.Task, key, variant)
	if len(res.APKs) == 0 {
		fmt.Println("  (no APK found at expected paths — check Gradle output config)")
		for _, p := range res.Expected {
			fmt.Printf("  expected: %s\n", p)
		}
		return
	}
	for _, apk := range res.APKs {
		fmt.Printf("  %s  %s\n", filepathBase(apk.Path), build.FormatSize(apk.Size))
		fmt.Printf("  → %s\n", apk.Path)
	}
	if res.LogPath != "" {
		fmt.Printf("  log: %s\n", res.LogPath)
	}
}

func printBuildFailure(res *build.Result, key, variant string) {
	fmt.Println()
	fmt.Println(strings.Repeat("═", 40))
	fmt.Printf("BUILD FAILED  ·  %s  ·  %s\n", key, variant)
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
	if len(res.Expected) > 0 {
		fmt.Print("Expected :")
		for i, p := range res.Expected {
			if i == 0 {
				fmt.Printf(" %s\n", p)
			} else {
				fmt.Printf("           %s\n", p)
			}
		}
	}
}

func filepathBase(p string) string {
	for i := len(p) - 1; i >= 0; i-- {
		if p[i] == '/' || p[i] == '\\' {
			return p[i+1:]
		}
	}
	return p
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
