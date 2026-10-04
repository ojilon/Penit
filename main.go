package main

import (
	"os"

	"github.com/ojilon/penit/cmd"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	cmd.SetVersionInfo(Version, BuildTime, GitCommit)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
