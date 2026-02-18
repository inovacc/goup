package goversion

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/inovacc/goup/internal/platform"
)

const (
	releasesURL    = "https://go.dev/dl/?mode=json"
	allReleasesURL = "https://go.dev/dl/?mode=json&include=all"
)

// FetchReleases retrieves available Go releases from the official API (stable only by default).
func FetchReleases() ([]Release, error) {
	return FetchReleasesFromURL(releasesURL)
}

// FetchAllReleases retrieves all Go releases including old and unstable versions.
func FetchAllReleases() ([]Release, error) {
	return FetchReleasesFromURL(allReleasesURL)
}

// FetchReleasesFromURL retrieves Go releases from the given URL.
func FetchReleasesFromURL(url string) ([]Release, error) {
	client := &http.Client{Timeout: 30 * time.Second}

	resp, err := client.Get(url)
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

// LatestStable returns the latest stable release from the official API.
func LatestStable() (*Release, error) {
	return LatestStableFromURL(releasesURL)
}

// LatestStableFromURL returns the latest stable release from the given URL.
func LatestStableFromURL(url string) (*Release, error) {
	releases, err := FetchReleasesFromURL(url)
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

// FindRelease finds a specific version in the releases list.
func FindRelease(releases []Release, version string) (*Release, error) {
	for i := range releases {
		if releases[i].Version == version {
			return &releases[i], nil
		}
	}

	return nil, fmt.Errorf("release %s not found", version)
}

// FindArchiveFile finds the archive file (tar.gz or zip) for the given platform.
func FindArchiveFile(release *Release, p platform.Info) (*File, error) {
	for i := range release.Files {
		f := &release.Files[i]
		if f.OS == p.OS && f.Arch == p.Arch && f.Kind == "archive" {
			return f, nil
		}
	}

	return nil, fmt.Errorf("no archive file found for %s/%s", p.OS, p.Arch)
}
