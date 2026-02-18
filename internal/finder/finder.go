package finder

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/inovacc/goup/internal/cache"
)

// Installation represents a discovered Go installation on the system.
type Installation struct {
	Version string // "go1.26.0"
	Path    string // path to go binary
	Source  string // "system", "goup cache", "PATH", "GOROOT", "scanned"
}

// Find discovers all Go installations on the system.
func Find() ([]Installation, error) {
	seen := make(map[string]bool)
	var results []Installation

	add := func(inst Installation) {
		abs, err := filepath.Abs(inst.Path)
		if err == nil {
			inst.Path = abs
		}
		// Resolve symlinks for dedup
		resolved, err := filepath.EvalSymlinks(inst.Path)
		if err != nil {
			resolved = inst.Path
		}
		if seen[resolved] {
			return
		}
		seen[resolved] = true
		results = append(results, inst)
	}

	// 1. GOROOT
	if goroot := os.Getenv("GOROOT"); goroot != "" {
		bin := goBinaryIn(goroot)
		if ver := goVersion(bin); ver != "" {
			add(Installation{Version: ver, Path: bin, Source: "GOROOT"})
		}
	}

	// 2. PATH
	if pathEnv := os.Getenv("PATH"); pathEnv != "" {
		for _, dir := range filepath.SplitList(pathEnv) {
			bin := filepath.Join(dir, goBinaryName())
			if ver := goVersion(bin); ver != "" {
				add(Installation{Version: ver, Path: bin, Source: "PATH"})
			}
		}
	}

	// 3. Goup cache
	versions, _ := cache.List()
	for _, v := range versions {
		bin, err := cache.GoBin(v)
		if err != nil {
			continue
		}
		if ver := goVersion(bin); ver != "" {
			add(Installation{Version: ver, Path: bin, Source: "goup cache"})
		}
	}

	// 4. Common system paths (platform-specific, run go version in parallel)
	candidates := scanCandidates()
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 8)

	for _, c := range candidates {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			if ver := goVersion(path); ver != "" {
				mu.Lock()
				add(Installation{Version: ver, Path: path, Source: "scanned"})
				mu.Unlock()
			}
		}(c)
	}

	wg.Wait()

	return results, nil
}

func goBinaryName() string {
	if runtime.GOOS == "windows" {
		return "go.exe"
	}
	return "go"
}

func goBinaryIn(dir string) string {
	return filepath.Join(dir, "bin", goBinaryName())
}

// goVersion runs `go version` on a binary and returns the version string, or "".
func goVersion(path string) string {
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return ""
	}

	out, err := exec.Command(path, "version").Output()
	if err != nil {
		return ""
	}

	return parseGoVersion(string(out))
}

// parseGoVersion extracts the version token (e.g. "go1.26.0") from `go version` output.
func parseGoVersion(output string) string {
	// Typical: "go version go1.26.0 windows/amd64"
	fields := strings.Fields(output)
	for _, f := range fields {
		if strings.HasPrefix(f, "go1") || strings.HasPrefix(f, "go2") {
			return f
		}
	}
	return ""
}

// expandGlob returns matching paths for a glob pattern, or nil.
func expandGlob(pattern string) []string {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil
	}
	return matches
}

// candidatesFromDirs returns go binary paths for directories that exist.
func candidatesFromDirs(dirs []string) []string {
	bin := goBinaryName()
	var out []string
	for _, d := range dirs {
		out = append(out, filepath.Join(d, "bin", bin))
	}
	return out
}

// scanCandidates returns platform-specific candidate binary paths (defined in paths_*.go).
func scanCandidates() []string {
	dirs := commonDirs()
	globs := commonGlobs()

	var candidates []string
	candidates = append(candidates, candidatesFromDirs(dirs)...)

	bin := goBinaryName()
	for _, g := range globs {
		for _, match := range expandGlob(g) {
			candidates = append(candidates, filepath.Join(match, "bin", bin))
		}
	}

	// Also check go/bin/ subdirectory pattern used by sdk downloads
	for _, g := range globs {
		for _, match := range expandGlob(g) {
			candidates = append(candidates, filepath.Join(match, "go", "bin", bin))
		}
	}

	return candidates
}

// FormatTable formats installations as an aligned table.
func FormatTable(installs []Installation) string {
	if len(installs) == 0 {
		return ""
	}

	// Find max widths
	maxVer, maxPath := 0, 0
	for _, inst := range installs {
		if len(inst.Version) > maxVer {
			maxVer = len(inst.Version)
		}
		if len(inst.Path) > maxPath {
			maxPath = len(inst.Path)
		}
	}

	var b strings.Builder
	for _, inst := range installs {
		_, _ = fmt.Fprintf(&b, "  %-*s    %-*s    (%s)\n", maxVer, inst.Version, maxPath, inst.Path, inst.Source)
	}
	return b.String()
}
