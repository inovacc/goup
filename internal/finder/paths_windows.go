package finder

import (
	"os"
	"path/filepath"
)

func commonDirs() []string {
	dirs := []string{
		`C:\Go`,
		`C:\Program Files\Go`,
		`C:\Program Files (x86)\Go`,
	}

	if local := os.Getenv("LOCALAPPDATA"); local != "" {
		dirs = append(dirs, filepath.Join(local, "Programs", "Go"))
	}

	if home, err := os.UserHomeDir(); err == nil {
		dirs = append(dirs, filepath.Join(home, "go"))
	}

	// Chocolatey
	dirs = append(dirs, `C:\ProgramData\chocolatey\lib\golang\tools`)

	return dirs
}

func commonGlobs() []string {
	var globs []string

	if home, err := os.UserHomeDir(); err == nil {
		globs = append(globs,
			filepath.Join(home, "sdk", "go*"),
			filepath.Join(home, "scoop", "apps", "go", "*"),
		)
	}

	return globs
}
