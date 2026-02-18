package goversion

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/inovacc/goup/internal/platform"
)

const releasesURL = "https://go.dev/dl/?mode=json"

// FetchReleases retrieves all available Go releases from the official API.
func FetchReleases() ([]Release, error) {
	client := &http.Client{Timeout: 30 * time.Second}

	resp, err := client.Get(releasesURL)
	if err != nil {
		return nil, fmt.Errorf("fetch releases: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var releases []Release
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, fmt.Errorf("decode releases: %w", err)
	}

	return releases, nil
}

// LatestStable returns the latest stable release.
func LatestStable() (*Release, error) {
	releases, err := FetchReleases()
	if err != nil {
		return nil, err
	}

	for i := range releases {
		if releases[i].Stable {
			return &releases[i], nil
		}
	}

	return nil, fmt.Errorf("no stable release found")
}

// FindFile finds the matching file for the given platform in a release.
func FindFile(release *Release, p platform.Info) (*File, error) {
	kind := p.InstallerKind()

	for i := range release.Files {
		f := &release.Files[i]
		if f.OS == p.OS && f.Arch == p.Arch && f.Kind == kind {
			return f, nil
		}
	}

	// Fallback to archive if installer not found.
	for i := range release.Files {
		f := &release.Files[i]
		if f.OS == p.OS && f.Arch == p.Arch && f.Kind == "archive" {
			return f, nil
		}
	}

	return nil, fmt.Errorf("no file found for %s/%s", p.OS, p.Arch)
}
