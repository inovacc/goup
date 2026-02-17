# Roadmap

## Current Status
**Overall Progress:** 70% - Core functionality implemented

## Phases

### Phase 1: Foundation [COMPLETE]
- [x] Project scaffold (structure, tooling, CI config)
- [x] Go release API client (`go.dev/dl/?mode=json`)
- [x] Platform detection (OS + arch via runtime)
- [x] CLI commands (check, install, list)

### Phase 2: Core Features [COMPLETE]
- [x] Download with SHA256 checksum verification
- [x] Platform-specific installers (Windows MSI, macOS pkg, Linux tar.gz)
- [x] Version comparison (installed vs latest)
- [x] Force reinstall flag

### Phase 3: Polish & Release [IN PROGRESS]
- [ ] Unit tests (80%+ coverage)
- [ ] Progress bar for downloads
- [ ] Proxy support
- [ ] v1.0.0 release

## Test Coverage
**Current:** 0.0%  |  **Target:** 80%

| Package | Coverage | Status |
|---------|----------|--------|
| cmd | 0.0% | No tests |
| internal/goversion | 0.0% | No tests |
| internal/installer | 0.0% | No tests |
| internal/platform | 0.0% | No tests |
