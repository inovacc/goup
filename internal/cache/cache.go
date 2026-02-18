package cache

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Dir returns the goup cache root directory (~/.goup).
func Dir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home dir: %w", err)
	}

	return filepath.Join(home, ".goup"), nil
}

// VersionDir returns the directory for a specific Go version in the cache.
func VersionDir(version string) (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "versions", version), nil
}

// GoBin returns the path to the go binary for a cached version.
func GoBin(version string) (string, error) {
	vdir, err := VersionDir(version)
	if err != nil {
		return "", err
	}

	bin := "go"
	if runtime.GOOS == "windows" {
		bin = "go.exe"
	}

	return filepath.Join(vdir, "go", "bin", bin), nil
}

// Has reports whether the given version is already cached.
func Has(version string) bool {
	bin, err := GoBin(version)
	if err != nil {
		return false
	}

	_, err = os.Stat(bin)

	return err == nil
}

// Install extracts a Go archive into the cache directory for the given version.
// archivePath should be a .tar.gz (Linux/macOS) or .zip (Windows) file.
func Install(version string, archivePath string) error {
	if err := checkNotRoot(); err != nil {
		return err
	}

	vdir, err := VersionDir(version)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(vdir, 0o755); err != nil {
		return fmt.Errorf("create version dir: %w", err)
	}

	if strings.HasSuffix(archivePath, ".zip") {
		return extractZip(archivePath, vdir)
	}

	return extractTarGz(archivePath, vdir)
}

// Remove deletes a cached Go version.
func Remove(version string) error {
	if err := checkNotRoot(); err != nil {
		return err
	}

	vdir, err := VersionDir(version)
	if err != nil {
		return err
	}

	return os.RemoveAll(vdir)
}

// List returns all cached Go versions.
func List() ([]string, error) {
	dir, err := Dir()
	if err != nil {
		return nil, err
	}

	versionsDir := filepath.Join(dir, "versions")

	entries, err := os.ReadDir(versionsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, fmt.Errorf("read versions dir: %w", err)
	}

	var versions []string

	for _, e := range entries {
		if e.IsDir() {
			versions = append(versions, e.Name())
		}
	}

	return versions, nil
}

func extractTarGz(archivePath, destDir string) error {
	cmd := exec.Command("tar", "-C", destDir, "-xzf", archivePath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("extract tar.gz: %w: %s", err, output)
	}

	return nil
}

func extractZip(archivePath, destDir string) error {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return fmt.Errorf("open zip: %w", err)
	}

	defer func() { _ = r.Close() }()

	for _, f := range r.File {
		target := filepath.Join(destDir, f.Name)

		// Prevent zip slip
		if !strings.HasPrefix(filepath.Clean(target), filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid zip entry: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}

			continue
		}

		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}

		if err := extractZipFile(f, target); err != nil {
			return err
		}
	}

	return nil
}

func extractZipFile(f *zip.File, target string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}

	defer func() { _ = rc.Close() }()

	out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}

	defer func() { _ = out.Close() }()

	_, err = io.Copy(out, rc)

	return err
}
