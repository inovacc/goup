package cmd

import (
	"fmt"
	"os"

	"github.com/inovacc/goupdater/internal/goversion"
	"github.com/inovacc/goupdater/internal/installer"
	"github.com/inovacc/goupdater/internal/platform"
	"github.com/spf13/cobra"
)

var forceInstall bool

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Download and install the latest Go version",
	RunE: func(cmd *cobra.Command, args []string) error {
		p := platform.Detect()
		fmt.Printf("Platform: %s/%s\n", p.OS, p.Arch)

		installed, _ := goversion.InstalledVersion()
		if installed != "" {
			fmt.Printf("Installed: %s\n", installed)
		}

		latest, err := goversion.LatestStable()
		if err != nil {
			return fmt.Errorf("fetch latest: %w", err)
		}

		fmt.Printf("Latest:    %s\n", latest.Version)

		if !forceInstall && installed != "" && !goversion.NeedsUpdate(installed, latest.Version) {
			fmt.Println("Already up to date. Use --force to reinstall.")
			return nil
		}

		file, err := goversion.FindFile(latest, p)
		if err != nil {
			return err
		}

		filePath, err := installer.Download(file)
		if err != nil {
			return err
		}
		defer installer.Cleanup(filePath)

		if err := installer.Install(filePath); err != nil {
			fmt.Fprintf(os.Stderr, "Installation failed: %v\n", err)
			return err
		}

		fmt.Printf("Successfully installed %s\n", latest.Version)

		return nil
	},
}

func init() {
	installCmd.Flags().BoolVarP(&forceInstall, "force", "f", false, "Force reinstall even if up to date")
	rootCmd.AddCommand(installCmd)
}
