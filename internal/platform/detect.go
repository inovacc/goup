// Package platform detects the current operating system and architecture.
package platform

import "runtime"

// Info holds the detected platform information.
type Info struct {
	OS   string
	Arch string
}

// Detect returns the current platform OS and architecture.
func Detect() Info {
	return Info{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}
}

// InstallerKind returns the preferred installer kind for the platform.
func (i Info) InstallerKind() string {
	switch i.OS {
	case "windows":
		return "installer"
	case "darwin":
		return "installer"
	default:
		return "archive"
	}
}

// FileExtension returns the expected file extension for the platform installer.
func (i Info) FileExtension() string {
	switch i.OS {
	case "windows":
		return ".msi"
	case "darwin":
		return ".pkg"
	default:
		return ".tar.gz"
	}
}
