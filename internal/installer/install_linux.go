package installer

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
)

// Install extracts the tar.gz archive to /usr/local, replacing any existing Go installation.
// Requires root privileges on Linux.
func Install(filePath string) error {
	if os.Geteuid() != 0 {
		return fmt.Errorf("install requires root privileges, run with sudo")
	}

	slog.Info("removing old Go installation")

	if err := os.RemoveAll("/usr/local/go"); err != nil {
		return fmt.Errorf("remove old go: %w", err)
	}

	slog.Info("extracting archive to /usr/local")

	cmd := exec.Command("tar", "-C", "/usr/local", "-xzf", filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("extract archive: %w", err)
	}

	slog.Info("installation complete, ensure /usr/local/go/bin is in your PATH")

	return nil
}

// Cleanup removes the downloaded archive file.
func Cleanup(filePath string) {
	_ = os.Remove(filePath)
}
