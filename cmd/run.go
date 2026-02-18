package cmd

import (
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/inovacc/goup/internal/cache"
	"github.com/inovacc/goup/internal/goversion"
	"github.com/inovacc/goup/internal/installer"
	"github.com/inovacc/goup/internal/platform"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run <version> [-- command args...]",
	Short: "Run a command using a specific Go version",
	Long: `Download a Go version to a local cache if needed, then run a command using it.

Examples:
  goup run 1.22.0 -- go build ./...
  goup run go1.25.7 -- go test -v ./...
  goup run 1.22.0                        # just ensure version is cached`,
	Args:               cobra.MinimumNArgs(1),
	DisableFlagParsing: true,
	RunE:               runRun,
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func runRun(cmd *cobra.Command, args []string) error {
	version := normalizeVersion(args[0])

	// Split on "--" to find the command to run
	var cmdArgs []string

	for i, a := range args[1:] {
		if a == "--" {
			cmdArgs = args[i+2:]
			break
		}
	}

	if err := ensureCached(version); err != nil {
		return err
	}

	if len(cmdArgs) == 0 {
		cmd.Printf("%s is ready\n", version)
		return nil
	}

	return execWithVersion(version, cmdArgs)
}

func normalizeVersion(v string) string {
	if !strings.HasPrefix(v, "go") {
		return "go" + v
	}

	return v
}

func ensureCached(version string) error {
	releases, err := goversion.FetchReleases()
	if err != nil {
		return err
	}

	release, err := goversion.FindRelease(releases, version)
	if err != nil {
		return err
	}

	p := platform.Detect()

	file, err := goversion.FindArchiveFile(release, p)
	if err != nil {
		return err
	}

	// Skip download if cached and checksum matches
	if cache.Has(version) {
		stored := cache.Checksum(version)
		if stored == file.SHA256 {
			slog.Info("cached version is valid, skipping download", "version", version)
			return nil
		}

		slog.Info("checksum mismatch, re-downloading", "version", version)

		if err := cache.Remove(version); err != nil {
			return err
		}
	}

	slog.Info("downloading", "version", version)

	archivePath, err := installer.Download(file)
	if err != nil {
		return err
	}

	if err := cache.Install(version, archivePath); err != nil {
		return err
	}

	if err := cache.SaveChecksum(version, file.SHA256); err != nil {
		return err
	}

	slog.Info("cached", "version", version)

	return nil
}

func execWithVersion(version string, args []string) error {
	goBin, err := cache.GoBin(version)
	if err != nil {
		return err
	}

	vdir, err := cache.VersionDir(version)
	if err != nil {
		return err
	}

	goroot := filepath.Join(vdir, "go")

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Set GOROOT and prepend cached go/bin to PATH
	env := os.Environ()
	binDir := filepath.Dir(goBin)
	env = replaceOrAddEnv(env, "GOROOT", goroot)
	env = prependPath(env, binDir)
	cmd.Env = env

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}

		return err
	}

	return nil
}

func replaceOrAddEnv(env []string, key, value string) []string {
	prefix := key + "="
	for i, e := range env {
		if strings.HasPrefix(strings.ToUpper(e), strings.ToUpper(prefix)) {
			env[i] = prefix + value
			return env
		}
	}

	return append(env, prefix+value)
}

func prependPath(env []string, dir string) []string {
	for i, e := range env {
		upper := strings.ToUpper(e)
		if strings.HasPrefix(upper, "PATH=") {
			eqIdx := strings.Index(e, "=")
			env[i] = e[:eqIdx+1] + dir + string(os.PathListSeparator) + e[eqIdx+1:]

			return env
		}
	}

	return append(env, "PATH="+dir)
}
