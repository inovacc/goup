package cmd

import (
	"os"

	"github.com/inovacc/goup/internal/goversion"
	"github.com/inovacc/goup/internal/installer"
	"github.com/inovacc/goup/internal/platform"
	"github.com/spf13/cobra"
)

var forceInstall bool

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Download and install the latest Go version",
	RunE: func(cmd *cobra.Command, args []string) error {
		p := platform.Detect()
		cmd.Printf("Platform: %s/%s\n", p.OS, p.Arch)

		installed, _ := goversion.InstalledVersion()
		if installed != "" {
			cmd.Printf("Installed: %s\n", installed)
		}

		latest, err := goversion.LatestStable()
		if err != nil {
			return err
		}

		cmd.Printf("Latest:    %s\n", latest.Version)

		if !forceInstall && installed != "" && !goversion.NeedsUpdate(installed, latest.Version) {
			cmd.Println("Already up to date. Use --force to reinstall.")
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

		if err := installer.Install(filePath); err != nil {
			cmd.PrintErrln("Installation failed:", err)
			os.Exit(1)
		}

		cmd.Printf("Successfully installed %s\n", latest.Version)

		return nil
	},
}

func init() {
	installCmd.Flags().BoolVarP(&forceInstall, "force", "f", false, "Force reinstall even if up to date")
	rootCmd.AddCommand(installCmd)
}
