# Feature Requests

## Completed Features

- **Check for updates** — Compares installed Go version against latest stable from go.dev API
- **Install latest Go** — Downloads, verifies SHA256, and runs platform-specific installer
- **List versions** — Lists available Go releases with stable/unstable filtering
- **Platform detection** — Automatically detects OS and architecture via runtime
- **Force reinstall** — `--force` flag to reinstall even when up to date

## Proposed Features

### Download Progress Bar
- **Priority:** P2
- **Description:** Show download progress with bytes transferred and ETA
- **Motivation:** Large downloads (~150MB) with no feedback is poor UX

### Proxy Support
- **Priority:** P2
- **Description:** Support HTTP/SOCKS proxy for corporate environments
- **Motivation:** PowerShell script had proxy fallback; Go version should too

### Install Specific Version
- **Priority:** P2
- **Description:** `goup install go1.22.5` to install a specific version
- **Motivation:** Some projects require pinned Go versions

### Go Tools Installation
- **Priority:** P3
- **Description:** Optional install of gopls, dlv, staticcheck, gotests, gomodifytags
- **Motivation:** PowerShell script installed these; useful for fresh setups
