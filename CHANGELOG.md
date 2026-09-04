# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project follows
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Tagged-release workflow that builds and publishes an `xapsd` Debian package.
- Interactive Debian configuration for the APNS topic, key ID, and Apple Developer team ID.
- Debian system user, state directory, service lifecycle, copyright, changelog, and removal handling.
- Lint, SonarQube Cloud, Dependabot, and Release Drafter automation.
- Contribution, security, conduct, funding, issue, and pull-request community files.

### Changed

- Moved the canonical Go module and internal imports to `github.com/timlaing/dovecot-xaps-daemon`.
- Updated tests to use the Go version declared in `go.mod` on Linux and macOS.
- Debian upgrades now preserve an existing `/etc/xapsd/xapsd.yaml`.
- Connected tagged releases to the categorized draft produced by Release Drafter.
- Updated documentation to reference the maintained `timlaing` repositories.
- Added white-box unit tests across all packages (database, config, APNs, socket, and the daemon entry point),
  covering every reachable branch; the suite passes with `go test -race`, `gofmt`, and `go vet`.
- Ignored the local `.env` file and the regenerated `database_workingcpy.json` test artifact.
- Added a "Development" section to the README describing the test, coverage, and verification commands.
