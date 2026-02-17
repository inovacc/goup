package installer

import (
	"fmt"
	"os"
	"os/exec"
)

// Install extracts the tar.gz archive to /usr/local, replacing any existing Go installation.
func Install(filePath string) error {
	fmt.Println("Removing old Go installation...")

	if err := os.RemoveAll("/usr/local/go"); err != nil {
		return fmt.Errorf("remove old go: %w", err)
	}

	fmt.Println("Extracting archive to /usr/local...")

	cmd := exec.Command("tar", "-C", "/usr/local", "-xzf", filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("extract archive: %w", err)
	}

	fmt.Println("Installation complete. Ensure /usr/local/go/bin is in your PATH.")

	return nil
}

// Cleanup removes the downloaded archive file.
func Cleanup(filePath string) {
	_ = os.Remove(filePath)
}
