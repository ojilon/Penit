package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print Penit version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("penit %s\n", version)
		fmt.Printf("  build : %s\n", buildTime)
		fmt.Printf("  commit: %s\n", gitCommit)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
