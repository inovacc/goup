# CLAUDE.md - goupdater

## Project Overview

Cross-platform Go CLI tool that detects, downloads, and installs the latest Go version. Uses the official Go download API (`https://go.dev/dl/?mode=json`).

## Architecture

- **Structure:** Hexagonal/Clean (`cmd/`, `internal/`)
- **CLI Framework:** Cobra
- **Module:** `github.com/inovacc/goupdater`

### Key Packages

| Package | Purpose |
|---------|---------|
| `cmd/` | Cobra CLI commands (check, install, list) |
| `internal/goversion/` | Go release API client, version comparison |
| `internal/installer/` | Download with SHA256 verification, platform-specific install |
| `internal/platform/` | OS/arch detection via `runtime.GOOS`/`runtime.GOARCH` |

### Platform-Specific Files

- `install_windows.go` — MSI via `msiexec`
- `install_darwin.go` — pkg via `sudo installer`
- `install_linux.go` — tar.gz extraction to `/usr/local`

## Build & Test

```bash
task build           # Build to dist/
task test            # Run tests with coverage
task check           # fmt + vet + lint + test
task run             # go run .
```

## Conventions

- Use `go run .` not `go build && ./app`
- Table-driven tests, 80%+ coverage target
- `errors.Is`/`errors.As` for error comparison
- `defer func() { _ = file.Close() }()` pattern
- Structured logging with `log/slog` when needed
