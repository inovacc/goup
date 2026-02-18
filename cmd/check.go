package cmd

import (
	"fmt"
	"os"

	"github.com/inovacc/goup/internal/goversion"
	"github.com/inovacc/goup/internal/platform"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for available Go updates",
	RunE: func(cmd *cobra.Command, args []string) error {
		p := platform.Detect()
		fmt.Printf("Platform: %s/%s\n", p.OS, p.Arch)

		installed, err := goversion.InstalledVersion()
		if err != nil {
			fmt.Println("Go is not currently installed.")
		} else {
			fmt.Printf("Installed: %s\n", installed)
		}

		latest, err := goversion.LatestStable()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching releases: %v\n", err)
			return err
		}

		fmt.Printf("Latest:    %s (stable)\n", latest.Version)

		if installed == "" || goversion.NeedsUpdate(installed, latest.Version) {
			fmt.Println("\nUpdate available! Run 'goup install' to update.")
		} else {
			fmt.Println("\nYou are up to date.")
		}

		file, err := goversion.FindFile(latest, p)
		if err != nil {
			return err
		}

		fmt.Printf("\nDownload:  %s (%d MB)\n", file.Filename, file.Size/(1024*1024))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
