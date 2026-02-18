package platform

import (
	"runtime"
	"testing"
)

func TestDetect(t *testing.T) {
	info := Detect()

	if info.OS == "" {
		t.Fatal("OS should not be empty")
	}
	if info.Arch == "" {
		t.Fatal("Arch should not be empty")
	}
	if info.OS != runtime.GOOS {
		t.Errorf("OS = %q, want %q", info.OS, runtime.GOOS)
	}
	if info.Arch != runtime.GOARCH {
		t.Errorf("Arch = %q, want %q", info.Arch, runtime.GOARCH)
	}
}

func TestInstallerKind(t *testing.T) {
	tests := []struct {
		os   string
		want string
	}{
		{"windows", "installer"},
		{"darwin", "installer"},
		{"linux", "archive"},
		{"freebsd", "archive"},
	}

	for _, tt := range tests {
		t.Run(tt.os, func(t *testing.T) {
			info := Info{OS: tt.os, Arch: "amd64"}
			if got := info.InstallerKind(); got != tt.want {
				t.Errorf("InstallerKind() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFileExtension(t *testing.T) {
	tests := []struct {
		os   string
		want string
	}{
		{"windows", ".msi"},
		{"darwin", ".pkg"},
		{"linux", ".tar.gz"},
		{"freebsd", ".tar.gz"},
	}

	for _, tt := range tests {
		t.Run(tt.os, func(t *testing.T) {
			info := Info{OS: tt.os, Arch: "amd64"}
			if got := info.FileExtension(); got != tt.want {
				t.Errorf("FileExtension() = %q, want %q", got, tt.want)
			}
		})
	}
}
