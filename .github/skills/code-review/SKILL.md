---
name: code-review
description: Review changes, pull requests, or commits in dovecot-xaps-daemon for Go correctness, APNS and HTTP safety, configuration compatibility, Debian packaging, systemd lifecycle behavior, tests, and release automation.
---

# Code Review — dovecot-xaps-daemon

Apply this checklist when reviewing this repository.

## Review standard

- Establish the target before reviewing: working tree, staged changes, commit range, or PR against its merge base. State the target and base.
- Review only changed behavior. Do not present pre-existing issues as findings unless the change materially worsens them.
- Report only demonstrable functional, security, reliability, compatibility, packaging, or maintainability defects. Give the triggering conditions, impact, and smallest useful `file:line` range.
- Do not invent findings or elevate style preferences. Keep optional improvements separate.
- Rank findings as P0 catastrophic/release-blocking, P1 likely serious defect, P2 real defect under plausible conditions, or P3 worthwhile non-blocking defect.

## Project invariants

- The canonical Go module is `github.com/timlaing/dovecot-xaps-daemon`; internal imports must use it.
- The default listener is IPv6 loopback `[::1]:11619`, matching the plugin's `http://[::1]:11619` endpoint.
- Configuration defaults and YAML keys are compatibility-sensitive. Preserve existing installations and fail clearly on invalid required APNS fields.
- The standard Apple token key filename derived by installation is `AuthKey_<KeyID>.p8`.
- Do not expose APNS private keys, key/team IDs, device tokens, account identifiers, registration data, email data, or private topology in logs, tests, fixtures, or errors.
- Review HTTP request limits, authentication, JSON escaping, concurrency, persistence atomicity, APNS response handling, and cleanup of inactive registrations.
- Goroutines and timers must have a bounded lifecycle or an intentional process-lifetime owner. Avoid races around registration data and delayed notifications.

## Tests and Go quality

- Run `gofmt`; changed Go files must produce no diff.
- Run `go test -race ./...` and `go vet ./...`. Add focused tests for changed parsing, persistence, concurrency, or error behavior.
- Avoid deprecated APIs when a supported equivalent exists. Handle returned errors unless intentionally safe and documented.
- Keep production package boundaries intact: entry point in `cmd/xapsd`, APNS/socket logic in `internal`, configuration under `internal/config`, and persistence under `internal/database`.

## Debian and systemd packaging

- `debian/` is the packaging source of truth. Build through `dpkg-buildpackage`, not a hand-assembled archive.
- Installation may create `/etc/xapsd/xapsd.yaml` from debconf answers only when it does not exist. Upgrades must preserve administrator edits.
- Debhelper-generated systemd lifecycle snippets must enable/start on install, restart appropriately on upgrade, and stop/disable on removal.
- Purge may remove generated YAML and debconf state, but must not delete APNS keys or registration data without explicit policy and documentation.
- Ensure service hardening still permits reading `/etc/xapsd` and writing `/var/lib/xapsd`.
- Run `lintian --fail-on error` against the resulting `.deb`; report warnings separately.

## Workflows and releases

- Lint, test, SonarCloud, Release Drafter, and release workflows must remain separate and use least-privilege permissions.
- Use the Go version declared in `go.mod` throughout CI.
- Tagged release versions must match `vMAJOR.MINOR.PATCH` and the Debian changelog/package version.
- Releases must publish the `.deb`, `SHA256SUMS`, and SPDX JSON SBOM, using the existing Release Drafter release notes.
- Do not expose `SONAR_TOKEN` or other secrets to untrusted fork pull requests.

## Verification

Run checks proportionate to the change:

```sh
test -z "$(gofmt -l cmd internal)"
go test -race ./...
go vet ./...
dpkg-buildpackage --build=binary --no-sign
lintian --fail-on error ../xapsd_*.deb
```

For maintainer-script changes, inspect generated package scripts and test fresh install, upgrade with edited YAML, removal, purge, and noninteractive installation in a disposable Debian/Ubuntu environment.

## Review output

1. List findings first in severity order with tight file and line references.
2. Explain why each issue is reachable and what fails.
3. Separate non-blocking suggestions from findings.
4. End with verification evidence and limitations.
5. If nothing qualifies, say explicitly that no actionable findings were found.

A review request is read-only. Do not edit, commit, push, resolve discussions, or merge unless the user explicitly asks.
