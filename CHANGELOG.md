# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project follows
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.4.0] - 2026-09-04

### Added

- Tagged-release workflow that builds and publishes an `xapsd` Debian package for amd64 and arm64.
- Interactive Debian configuration for the APNS topic, key ID, and Apple Developer team ID.
- Debian system user, state directory, service lifecycle, copyright, changelog, and removal handling.
- Lint, SonarQube Cloud, Dependabot, and Release Drafter automation.
- Contribution, security, conduct, funding, issue, and pull-request community files.
- Lintian override for the statically-linked Go binary.
- CI notification to `dovecot-xaps-apt` after publishing a release.
- White-box unit tests across all packages (database, config, APNs, socket, and the daemon entry point),
  covering every reachable branch; the suite passes with `go test -race`, `gofmt`, and `go vet`.
- A "Development" section to the README describing the test, coverage, and verification commands.

### Changed

- Moved the canonical Go module and internal imports to `github.com/timlaing/dovecot-xaps-daemon`.
- Updated tests to use the Go version declared in `go.mod` on Linux and macOS.
- Debian upgrades now preserve an existing `/etc/xapsd/xapsd.yaml`.
- Connected tagged releases to the categorized draft produced by Release Drafter.
- Updated documentation to reference the maintained `timlaing` repositories.
- Made the APNS credential directory and the HTTP router construction testable without changing runtime behavior:
  the daemon still loads credentials from `/etc/xapsd/` by default and `NewApns`/`NewHttpSocket` keep their
  signatures, only now the unit tests can drive them from a temporary directory. This raised aggregate statement
  coverage to 85.4%.
- Introduced the `config.XapsdConfigDir` constant so the `/etc/xapsd/` path is defined once rather than duplicated
  across the config and APNS packages.
- Ignored the local `.env` file and the regenerated `database_workingcpy.json` test artifact.
- Updated the SonarQube analysis to include both `cmd` and `internal` test directories and to report the aggregated
  Go coverage, forcing a fresh coverage profile on every run.
- Used Production environment for the APNs token client.
- Bumped `golang.org/x/net` from 0.47.0 to 0.55.0.
- Bumped `golang.org/x/crypto` from 0.40.0 to 0.45.0.

### Fixed

- Do not panic if APNs certificate renewal fails.
- Increased the limit for the HTTP response reader to 2^32.
- Fixed multiple upstream security vulnerabilities: CVE-2023-45288, CVE-2024-45338, CVE-2024-51744,
  CVE-2025-22870, and GHSA-fv92-fjc5-jj9h.

### Removed

- Deprecated Apple API for self-service push certificates.

[Unreleased]: https://github.com/timlaing/dovecot-xaps-daemon/compare/v1.4.0...HEAD
[1.4.0]: https://github.com/timlaing/dovecot-xaps-daemon/releases/tag/v1.4.0
