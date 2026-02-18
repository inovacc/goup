package cache

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDir(t *testing.T) {
	dir, err := Dir()
	if err != nil {
		t.Fatal(err)
	}

	home, _ := os.UserHomeDir()

	want := filepath.Join(home, ".goup")
	if dir != want {
		t.Errorf("Dir() = %q, want %q", dir, want)
	}
}

func TestVersionDir(t *testing.T) {
	vdir, err := VersionDir("go1.22.0")
	if err != nil {
		t.Fatal(err)
	}

	home, _ := os.UserHomeDir()

	want := filepath.Join(home, ".goup", "versions", "go1.22.0")
	if vdir != want {
		t.Errorf("VersionDir() = %q, want %q", vdir, want)
	}
}

func TestHas(t *testing.T) {
	// Non-existent version should return false
	if Has("go99.99.99") {
		t.Error("Has() returned true for non-existent version")
	}
}

func TestList(t *testing.T) {
	dir, err := Dir()
	if err != nil {
		t.Fatal(err)
	}

	versionsDir := filepath.Join(dir, "versions")
	testDirs := []string{"go1.21.0", "go1.22.0"}

	for _, d := range testDirs {
		p := filepath.Join(versionsDir, d)
		if err := os.MkdirAll(p, 0o755); err != nil {
			t.Fatal(err)
		}
	}

	defer func() { _ = os.RemoveAll(versionsDir) }()

	versions, err := List()
	if err != nil {
		t.Fatal(err)
	}

	if len(versions) < 2 {
		t.Fatalf("List() returned %d versions, want at least 2", len(versions))
	}

	found := map[string]bool{}
	for _, v := range versions {
		found[v] = true
	}

	for _, d := range testDirs {
		if !found[d] {
			t.Errorf("List() missing %q", d)
		}
	}
}
