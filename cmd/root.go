package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goup",
	Short: "Detect, download, and install the latest Go version",
	Long:  `goup checks the official Go download API, detects your OS and architecture, and installs the latest stable Go release.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
