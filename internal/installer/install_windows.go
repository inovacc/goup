package installer

import (
	"fmt"
	"os"
	"os/exec"
)

// Install runs the platform-specific installer for the downloaded file.
func Install(filePath string) error {
	fmt.Println("Installing via MSI...")

	cmd := exec.Command("msiexec", "/i", filePath, "/quiet", "ADDLOCAL=ALL")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("msi install: %w", err)
	}

	fmt.Println("Installation complete.")

	return nil
}

// Cleanup removes the downloaded installer file.
func Cleanup(filePath string) {
	_ = os.Remove(filePath)
}
