package cmd

import (
	"strings"

	"github.com/inovacc/goup/internal/cache"
	"github.com/inovacc/goup/internal/goversion"
	"github.com/spf13/cobra"
)

var listAll bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available Go versions",
	RunE: func(cmd *cobra.Command, args []string) error {
		var releases []goversion.Release
		var err error

		if listAll {
			releases, err = goversion.FetchAllReleases()
		} else {
			releases, err = goversion.FetchReleases()
		}

		if err != nil {
			return err
		}

		installed, _ := goversion.InstalledVersion()
		cached, _ := cache.List()
		cachedSet := make(map[string]bool, len(cached))
		for _, v := range cached {
			cachedSet[v] = true
		}

		// Group releases by minor version (e.g., "1.22", "1.21")
		type group struct {
			minor    string
			versions []goversion.Release
		}

		var groups []group
		groupIdx := make(map[string]int)

		for _, r := range releases {
			minor := minorVersion(r.Version)
			if idx, ok := groupIdx[minor]; ok {
				groups[idx].versions = append(groups[idx].versions, r)
			} else {
				groupIdx[minor] = len(groups)
				groups = append(groups, group{minor: minor, versions: []goversion.Release{r}})
			}
		}

		for _, g := range groups {
			cmd.Printf("Go %s:\n", g.minor)
			for _, r := range g.versions {
				marker := "  "
				if r.Version == installed {
					marker = "* "
				}

				labels := ""
				if !r.Stable {
					labels += " (unstable)"
				}
				if cachedSet[r.Version] {
					labels += " [cached]"
				}
				if r.Version == installed {
					labels += " [installed]"
				}

				cmd.Printf("  %s%s%s\n", marker, r.Version, labels)
			}
		}

		return nil
	},
}

func init() {
	listCmd.Flags().BoolVarP(&listAll, "all", "a", false, "Include all releases (old and unstable)")
	rootCmd.AddCommand(listCmd)
}

// minorVersion extracts the minor version from a Go version string.
// e.g., "go1.22.3" -> "1.22", "go1.22rc1" -> "1.22", "go1.22" -> "1.22"
func minorVersion(version string) string {
	v := strings.TrimPrefix(version, "go")

	// Find first dot
	dotIdx := strings.Index(v, ".")
	if dotIdx == -1 {
		return v
	}

	rest := v[dotIdx+1:]

	// Find second dot or first non-digit (for rc/beta suffixes)
	for i, c := range rest {
		if c == '.' || (c < '0' || c > '9') {
			return v[:dotIdx+1+i]
		}
	}

	return v
}
