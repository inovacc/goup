//go:build !windows

package finder

import (
	"os"
	"path/filepath"
	"runtime"
)

func commonDirs() []string {
	dirs := []string{
		"/usr/local/go",
		"/usr/lib/go",
		"/snap/go/current",
		"/home/linuxbrew/.linuxbrew/opt/go",
	}

	if runtime.GOOS == "darwin" {
		dirs = append(dirs, "/opt/homebrew/opt/go", "/usr/local/opt/go")
	}

	if home, err := os.UserHomeDir(); err == nil {
		dirs = append(dirs, filepath.Join(home, "go"))
	}

	return dirs
}

func commonGlobs() []string {
	var globs []string

	globs = append(globs, "/usr/lib/go-*")

	if home, err := os.UserHomeDir(); err == nil {
		globs = append(globs, filepath.Join(home, "sdk", "go*"))
	}

	return globs
}
