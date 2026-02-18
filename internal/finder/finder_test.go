package finder

import (
	"testing"
)

func TestParseGoVersion(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		want   string
	}{
		{"standard", "go version go1.26.0 windows/amd64", "go1.26.0"},
		{"linux", "go version go1.22.5 linux/amd64", "go1.22.5"},
		{"darwin arm", "go version go1.24.3 darwin/arm64", "go1.24.3"},
		{"devel", "go version devel go1.23-abcdef linux/amd64", "go1.23-abcdef"},
		{"empty", "", ""},
		{"garbage", "not a go version string", ""},
		{"go2 future", "go version go2.0.0 linux/amd64", "go2.0.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseGoVersion(tt.input)
			if got != tt.want {
				t.Errorf("parseGoVersion(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestFormatTable(t *testing.T) {
	installs := []Installation{
		{Version: "go1.26.0", Path: "/usr/local/go/bin/go", Source: "system"},
		{Version: "go1.22.0", Path: "/home/user/.goup/versions/go1.22.0/go/bin/go", Source: "goup cache"},
	}

	out := FormatTable(installs)
	if out == "" {
		t.Fatal("expected non-empty output")
	}

	// Should contain both versions
	if !contains(out, "go1.26.0") || !contains(out, "go1.22.0") {
		t.Errorf("table missing versions: %s", out)
	}

	// Should contain sources
	if !contains(out, "system") || !contains(out, "goup cache") {
		t.Errorf("table missing sources: %s", out)
	}
}

func TestFormatTableEmpty(t *testing.T) {
	if got := FormatTable(nil); got != "" {
		t.Errorf("expected empty string for nil input, got %q", got)
	}
}

func TestDeduplication(t *testing.T) {
	// Test that the dedup logic in Find works by verifying the seen map approach
	// We can't easily test Find() itself without real binaries, but we can test
	// that parseGoVersion handles the version extraction correctly for dedup
	v1 := parseGoVersion("go version go1.26.0 windows/amd64")
	v2 := parseGoVersion("go version go1.26.0 linux/amd64")
	if v1 != v2 {
		t.Errorf("same version from different platforms should parse equal: %q vs %q", v1, v2)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
