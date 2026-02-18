package goversion

import (
	"encoding/json"
	"testing"
)

func TestReleaseJSON(t *testing.T) {
	original := Release{
		Version: "go1.23.0",
		Stable:  true,
		Files: []File{
			{
				Filename: "go1.23.0.windows-amd64.msi",
				OS:       "windows",
				Arch:     "amd64",
				Version:  "go1.23.0",
				SHA256:   "abc123",
				Size:     123456789,
				Kind:     "installer",
			},
		},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var decoded Release
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if decoded.Version != original.Version {
		t.Errorf("Version = %q, want %q", decoded.Version, original.Version)
	}
	if decoded.Stable != original.Stable {
		t.Errorf("Stable = %v, want %v", decoded.Stable, original.Stable)
	}
	if len(decoded.Files) != 1 {
		t.Fatalf("Files len = %d, want 1", len(decoded.Files))
	}

	f := decoded.Files[0]
	if f.Filename != "go1.23.0.windows-amd64.msi" {
		t.Errorf("Filename = %q", f.Filename)
	}
	if f.SHA256 != "abc123" {
		t.Errorf("SHA256 = %q", f.SHA256)
	}
	if f.Size != 123456789 {
		t.Errorf("Size = %d", f.Size)
	}
}
