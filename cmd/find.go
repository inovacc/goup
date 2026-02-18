package cmd

import (
	"fmt"

	"github.com/inovacc/goup/internal/finder"
	"github.com/spf13/cobra"
)

var findCmd = &cobra.Command{
	Use:   "find",
	Short: "Discover all Go installations on the system",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Println("Scanning for Go installations...")
		cmd.Println()

		installs, err := finder.Find()
		if err != nil {
			return fmt.Errorf("scan failed: %w", err)
		}

		if len(installs) == 0 {
			cmd.Println("No Go installations found.")
			return nil
		}

		cmd.Print(finder.FormatTable(installs))
		cmd.Println()
		cmd.Printf("Found %d Go installation(s).\n", len(installs))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(findCmd)
}
