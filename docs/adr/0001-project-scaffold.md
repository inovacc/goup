# ADR-0001: Project Scaffold and Tooling Choices

## Status
Accepted

## Context
Replacing the `golang-setup.ps1` PowerShell script with a cross-platform Go CLI that detects OS, downloads, and installs Go updates.

## Decision
- **Structure:** Hexagonal/Clean Architecture (cmd/, internal/)
- **CLI Framework:** Cobra
- **API Source:** `https://go.dev/dl/?mode=json` for release metadata
- **Platform Install:** Build-tag separated files (`install_windows.go`, `install_darwin.go`, `install_linux.go`)
- **Checksum:** SHA256 verification on every download
- **Task Runner:** Taskfile over Makefile for cross-platform support

## Consequences

### Positive
- Works on all platforms (Windows, macOS, Linux) from a single binary
- SHA256 verification prevents corrupted installs
- No PowerShell/bash dependency — pure Go

### Negative
- Requires elevated privileges for installation (sudo on Linux/macOS, admin on Windows)
