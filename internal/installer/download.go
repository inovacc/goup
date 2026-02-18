// Package installer downloads and installs Go releases.
package installer

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/inovacc/goup/internal/goversion"
)

const downloadBaseURL = "https://go.dev/dl/"

// Download fetches a Go release file to a temporary directory and verifies its checksum.
// Returns the path to the downloaded file.
func Download(file *goversion.File) (string, error) {
	url := downloadBaseURL + file.Filename
	destDir := os.TempDir()
	destPath := filepath.Join(destDir, file.Filename)

	fmt.Printf("Downloading %s (%d MB)...\n", file.Filename, file.Size/(1024*1024))

	client := &http.Client{Timeout: 10 * time.Minute}

	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("download: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed: %s", resp.Status)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("create file: %w", err)
	}
	defer func() { _ = out.Close() }()

	hasher := sha256.New()
	writer := io.MultiWriter(out, hasher)

	if _, err := io.Copy(writer, resp.Body); err != nil {
		return "", fmt.Errorf("write file: %w", err)
	}

	checksum := hex.EncodeToString(hasher.Sum(nil))
	if checksum != file.SHA256 {
		_ = os.Remove(destPath)
		return "", fmt.Errorf("checksum mismatch: got %s, want %s", checksum, file.SHA256)
	}

	fmt.Println("Download complete. Checksum verified.")

	return destPath, nil
}
