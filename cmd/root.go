package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"

	dataRootFlag string
	configFlag   string
)

// SetVersionInfo is called from main with ldflags values.
func SetVersionInfo(v, bt, gc string) {
	version = v
	buildTime = bt
	gitCommit = gc
}

var rootCmd = &cobra.Command{
	Use:   "penit",
	Short: "Minimal Android APK build · sign · install · debug tool",
	Long: `Penit is a focused tool for building, signing, packaging, installing,
and debugging Android APKs from a small set of local Gradle projects
(wayer, conductino_android, fdroid, cooda, …).

It will be merged into ReleaseForge over time; until then it is usable standalone.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command.
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	return nil
}

func init() {
	rootCmd.PersistentFlags().StringVar(&dataRootFlag, "data-root", "", "override data-root directory")
	rootCmd.PersistentFlags().StringVar(&configFlag, "config", "", "override config file path")
}
