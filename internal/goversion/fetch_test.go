package goversion

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/inovacc/goup/internal/platform"
)

func mockReleasesServer(t *testing.T, releases []Release) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(releases); err != nil {
			t.Fatalf("encode: %v", err)
		}
	}))
}

func TestFetchReleases(t *testing.T) {
	releases := []Release{
		{Version: "go1.23.0", Stable: true, Files: []File{{Filename: "go1.23.0.linux-amd64.tar.gz", OS: "linux", Arch: "amd64", Kind: "archive"}}},
		{Version: "go1.22.0", Stable: true},
	}

	srv := mockReleasesServer(t, releases)
	defer srv.Close()

	got, err := FetchReleasesFromURL(srv.URL)
	if err != nil {
		t.Fatalf("FetchReleasesFromURL: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if got[0].Version != "go1.23.0" {
		t.Errorf("Version = %q, want go1.23.0", got[0].Version)
	}
}

func TestFetchReleasesError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	_, err := FetchReleasesFromURL(srv.URL)
	if err == nil {
		t.Fatal("expected error for 500 response")
	}
}

func TestLatestStable(t *testing.T) {
	releases := []Release{
		{Version: "go1.24rc1", Stable: false},
		{Version: "go1.23.0", Stable: true},
		{Version: "go1.22.0", Stable: true},
	}

	srv := mockReleasesServer(t, releases)
	defer srv.Close()

	got, err := LatestStableFromURL(srv.URL)
	if err != nil {
		t.Fatalf("LatestStableFromURL: %v", err)
	}
	if got.Version != "go1.23.0" {
		t.Errorf("Version = %q, want go1.23.0", got.Version)
	}
}

func TestLatestStableNoStable(t *testing.T) {
	releases := []Release{
		{Version: "go1.24rc1", Stable: false},
		{Version: "go1.24beta1", Stable: false},
	}

	srv := mockReleasesServer(t, releases)
	defer srv.Close()

	_, err := LatestStableFromURL(srv.URL)
	if err == nil {
		t.Fatal("expected error when no stable release")
	}
}

func TestFindFile(t *testing.T) {
	release := &Release{
		Version: "go1.23.0",
		Files: []File{
			{Filename: "go1.23.0.linux-amd64.tar.gz", OS: "linux", Arch: "amd64", Kind: "archive"},
			{Filename: "go1.23.0.windows-amd64.msi", OS: "windows", Arch: "amd64", Kind: "installer"},
			{Filename: "go1.23.0.windows-amd64.zip", OS: "windows", Arch: "amd64", Kind: "archive"},
			{Filename: "go1.23.0.darwin-arm64.pkg", OS: "darwin", Arch: "arm64", Kind: "installer"},
		},
	}

	tests := []struct {
		name     string
		p        platform.Info
		wantFile string
	}{
		{"linux archive", platform.Info{OS: "linux", Arch: "amd64"}, "go1.23.0.linux-amd64.tar.gz"},
		{"windows installer", platform.Info{OS: "windows", Arch: "amd64"}, "go1.23.0.windows-amd64.msi"},
		{"darwin installer", platform.Info{OS: "darwin", Arch: "arm64"}, "go1.23.0.darwin-arm64.pkg"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := FindFile(release, tt.p)
			if err != nil {
				t.Fatalf("FindFile: %v", err)
			}
			if f.Filename != tt.wantFile {
				t.Errorf("Filename = %q, want %q", f.Filename, tt.wantFile)
			}
		})
	}
}

func TestFindFileFallbackToArchive(t *testing.T) {
	release := &Release{
		Version: "go1.23.0",
		Files: []File{
			{Filename: "go1.23.0.windows-amd64.zip", OS: "windows", Arch: "amd64", Kind: "archive"},
		},
	}

	// Windows wants "installer" but only archive exists — should fallback.
	f, err := FindFile(release, platform.Info{OS: "windows", Arch: "amd64"})
	if err != nil {
		t.Fatalf("FindFile: %v", err)
	}
	if f.Kind != "archive" {
		t.Errorf("Kind = %q, want archive", f.Kind)
	}
}

func TestFindFileNoMatch(t *testing.T) {
	release := &Release{
		Version: "go1.23.0",
		Files: []File{
			{Filename: "go1.23.0.linux-amd64.tar.gz", OS: "linux", Arch: "amd64", Kind: "archive"},
		},
	}

	_, err := FindFile(release, platform.Info{OS: "freebsd", Arch: "arm64"})
	if err == nil {
		t.Fatal("expected error for no matching file")
	}
}
