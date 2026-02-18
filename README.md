# goup

> Cross-platform CLI tool to detect, download, and install the latest Go version from the official API.

## Features

- Fetches available Go releases from `https://go.dev/dl/?mode=json`
- Detects current OS and architecture automatically
- Compares installed Go version against latest stable
- Downloads with SHA256 checksum verification
- Platform-specific installation (MSI on Windows, pkg on macOS, tar.gz on Linux)
- Force reinstall option

## Prerequisites

- Go 1.25+
- [Task](https://taskfile.dev) (task runner)

## Install

```bash
go install github.com/inovacc/goup@latest
```

### From Source

```bash
git clone https://github.com/inovacc/goup.git
cd goup
task deps
task build
```

## Usage

```bash
# Check for updates
goup check

# List available Go versions
goup list
goup list --all    # Include unstable releases

# Install latest stable Go
goup install
goup install --force   # Force reinstall
```

## Architecture

```
goup/
├── cmd/                        # CLI commands (Cobra)
│   ├── root.go                 # Root command
│   ├── check.go                # Check for updates
│   ├── install.go              # Download and install
│   └── list.go                 # List available versions
├── internal/
│   ├── goversion/              # Go release API client
│   │   ├── types.go            # Release/File structs
│   │   ├── fetch.go            # API fetch + file matching
│   │   └── installed.go        # Detect installed version
│   ├── installer/              # Download and install
│   │   ├── download.go         # Download with checksum
│   │   ├── install_windows.go  # MSI installer
│   │   ├── install_darwin.go   # pkg installer
│   │   └── install_linux.go    # tar.gz extraction
│   └── platform/               # OS/arch detection
│       └── detect.go
├── main.go
├── Taskfile.yml
└── go.mod
```

## Development

```bash
task --list          # List all tasks
task check           # Run fmt, vet, lint, test
task test:coverage   # Generate coverage report
```

## License

BSD 3-Clause License - see [LICENSE](LICENSE) for details.
