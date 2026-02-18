package goversion

import "testing"

func TestNeedsUpdate(t *testing.T) {
	tests := []struct {
		name      string
		installed string
		target    string
		want      bool
	}{
		{"same version", "go1.23.0", "go1.23.0", false},
		{"different version", "go1.22.0", "go1.23.0", true},
		{"empty installed", "", "go1.23.0", true},
		{"empty target", "go1.23.0", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NeedsUpdate(tt.installed, tt.target); got != tt.want {
				t.Errorf("NeedsUpdate(%q, %q) = %v, want %v", tt.installed, tt.target, got, tt.want)
			}
		})
	}
}
