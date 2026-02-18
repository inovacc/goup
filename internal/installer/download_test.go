package installer

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/inovacc/goup/internal/goversion"
)

func TestDownload(t *testing.T) {
	content := []byte("hello goup test binary")
	hash := sha256.Sum256(content)
	checksum := hex.EncodeToString(hash[:])

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(content)
	}))
	defer srv.Close()

	file := &goversion.File{
		Filename: "testfile.tar.gz",
		SHA256:   checksum,
		Size:     int64(len(content)),
	}

	path, err := DownloadFromURL(file, srv.URL+"/")
	if err != nil {
		t.Fatalf("DownloadFromURL: %v", err)
	}

	defer func() { _ = os.Remove(path) }()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	if string(data) != string(content) {
		t.Errorf("content mismatch")
	}
}

func TestDownloadChecksumMismatch(t *testing.T) {
	content := []byte("some content")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(content)
	}))
	defer srv.Close()

	file := &goversion.File{
		Filename: "badchecksum.tar.gz",
		SHA256:   "0000000000000000000000000000000000000000000000000000000000000000",
		Size:     int64(len(content)),
	}

	_, err := DownloadFromURL(file, srv.URL+"/")
	if err == nil {
		t.Fatal("expected checksum mismatch error")
	}
}

func TestDownloadHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	file := &goversion.File{
		Filename: "notfound.tar.gz",
		SHA256:   "abc",
		Size:     100,
	}

	_, err := DownloadFromURL(file, srv.URL+"/")
	if err == nil {
		t.Fatal("expected error for 404")
	}
}

func TestDownloadSkipsWhenFileExists(t *testing.T) {
	content := []byte("cached binary content")
	hash := sha256.Sum256(content)
	checksum := hex.EncodeToString(hash[:])

	// Pre-create the file in temp with correct content
	destPath := filepath.Join(os.TempDir(), "skip-test-cached.tar.gz")
	if err := os.WriteFile(destPath, content, 0o644); err != nil {
		t.Fatal(err)
	}

	defer func() { _ = os.Remove(destPath) }()

	// Server should never be hit
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("server was called, download should have been skipped")
	}))
	defer srv.Close()

	file := &goversion.File{
		Filename: "skip-test-cached.tar.gz",
		SHA256:   checksum,
		Size:     int64(len(content)),
	}

	path, err := DownloadFromURL(file, srv.URL+"/")
	if err != nil {
		t.Fatalf("DownloadFromURL: %v", err)
	}

	if path != destPath {
		t.Errorf("path = %q, want %q", path, destPath)
	}
}

func TestDownloadRedownloadsOnChecksumMismatch(t *testing.T) {
	oldContent := []byte("old corrupted content")
	newContent := []byte("fresh correct content")
	hash := sha256.Sum256(newContent)
	checksum := hex.EncodeToString(hash[:])

	// Pre-create the file with wrong content
	destPath := filepath.Join(os.TempDir(), "redownload-test.tar.gz")
	if err := os.WriteFile(destPath, oldContent, 0o644); err != nil {
		t.Fatal(err)
	}

	defer func() { _ = os.Remove(destPath) }()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(newContent)
	}))
	defer srv.Close()

	file := &goversion.File{
		Filename: "redownload-test.tar.gz",
		SHA256:   checksum,
		Size:     int64(len(newContent)),
	}

	path, err := DownloadFromURL(file, srv.URL+"/")
	if err != nil {
		t.Fatalf("DownloadFromURL: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != string(newContent) {
		t.Error("file was not re-downloaded with correct content")
	}
}

func TestCleanup(t *testing.T) {
	f, err := os.CreateTemp("", "goup-test-*")
	if err != nil {
		t.Fatalf("CreateTemp: %v", err)
	}

	path := f.Name()
	_ = f.Close()

	Cleanup(path)

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("file still exists after Cleanup")
	}
}
