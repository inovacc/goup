package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goupdater",
	Short: "Detect, download, and install the latest Go version",
	Long:  `goupdater checks the official Go download API, detects your OS and architecture, and installs the latest stable Go release.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
