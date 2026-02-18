package installer

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
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
	return DownloadFromURL(file, downloadBaseURL)
}

// DownloadFromURL fetches a Go release file from the given base URL and verifies its checksum.
// If the file already exists locally with a matching SHA256, the download is skipped.
func DownloadFromURL(file *goversion.File, baseURL string) (string, error) {
	url := baseURL + file.Filename
	destDir := os.TempDir()
	destPath := filepath.Join(destDir, file.Filename)

	// Check if file already exists with correct checksum
	if existing, err := checksumFile(destPath); err == nil && existing == file.SHA256 {
		slog.Info("file already exists with valid checksum, skipping download", "file", file.Filename)
		return destPath, nil
	}

	slog.Info("downloading", "file", file.Filename, "size_mb", file.Size/(1024*1024))

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

	slog.Info("download complete, checksum verified")

	return destPath, nil
}

// checksumFile computes the SHA256 checksum of a local file.
func checksumFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer func() { _ = f.Close() }()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
