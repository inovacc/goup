package goversion

import (
	"fmt"
	"os/exec"
	"strings"
)

// InstalledVersion returns the currently installed Go version string (e.g., "go1.23.0"),
// or an empty string if Go is not installed.
func InstalledVersion() (string, error) {
	cmd := exec.Command("go", "version")

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("go not found: %w", err)
	}

	// Output: "go version go1.23.0 windows/amd64"
	parts := strings.Fields(string(out))
	if len(parts) < 3 {
		return "", fmt.Errorf("unexpected go version output: %s", string(out))
	}

	return parts[2], nil
}

// NeedsUpdate returns true if the installed version differs from the target.
func NeedsUpdate(installed, target string) bool {
	return installed != target
}
