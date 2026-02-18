package installer

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	"golang.org/x/sys/windows"
)

// Install runs the platform-specific installer for the downloaded file.
// If not running as admin, re-launches msiexec elevated via UAC prompt.
func Install(filePath string) error {
	slog.Info("installing via MSI")

	if !isAdmin() {
		slog.Info("requesting admin privileges")
		return installElevated(filePath)
	}

	cmd := exec.Command("msiexec", "/i", filePath, "/quiet", "ADDLOCAL=ALL")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("msi install: %w", err)
	}

	slog.Info("installation complete")

	return nil
}

func isAdmin() bool {
	token := windows.GetCurrentProcessToken()
	return token.IsElevated()
}

func installElevated(filePath string) error {
	args := fmt.Sprintf("/i %q /quiet ADDLOCAL=ALL", filePath)

	cmd := exec.Command("powershell", "-Command",
		fmt.Sprintf(`Start-Process msiexec -ArgumentList '%s' -Verb RunAs -Wait`, args))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("elevated msi install: %w", err)
	}

	slog.Info("installation complete")

	return nil
}

// Cleanup removes the downloaded installer file.
func Cleanup(filePath string) {
	_ = os.Remove(filePath)
}
