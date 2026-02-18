package cmd

import (
	"fmt"

	"github.com/inovacc/goup/internal/goversion"
	"github.com/spf13/cobra"
)

var listAll bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available Go versions",
	RunE: func(cmd *cobra.Command, args []string) error {
		releases, err := goversion.FetchReleases()
		if err != nil {
			return err
		}

		for _, r := range releases {
			if !listAll && !r.Stable {
				continue
			}

			label := ""
			if !r.Stable {
				label = " (unstable)"
			}

			fmt.Printf("%s%s\n", r.Version, label)
		}

		return nil
	},
}

func init() {
	listCmd.Flags().BoolVarP(&listAll, "all", "a", false, "Include unstable releases")
	rootCmd.AddCommand(listCmd)
}
