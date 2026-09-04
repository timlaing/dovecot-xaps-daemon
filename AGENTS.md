# Repository guidance

## Project

This repository provides the APNS daemon used by
[`timlaing/dovecot-xaps-plugin`](https://github.com/timlaing/dovecot-xaps-plugin). It is maintained from
[`freswa/dovecot-xaps-daemon`](https://github.com/freswa/dovecot-xaps-daemon), which derives from Stefan Arentz's
original project; preserve that attribution.

## Development rules

- Use the module path `github.com/timlaing/dovecot-xaps-daemon` for internal imports.
- Use the Go version declared by `go.mod`.
- Keep `[::1]:11619` aligned with the plugin's default endpoint.
- Preserve existing `/etc/xapsd/xapsd.yaml` values during upgrades.
- Derive the Apple token filename as `AuthKey_<KeyID>.p8` during first installation.
- Treat APNS credentials, device tokens, account IDs, registration data, email data, and private server details as secrets.
- Keep `debian/` as the package source of truth; do not reintroduce hand-built `DEBIAN/control` archives.
- Preserve debhelper-managed systemd enable/start/restart/stop behavior and service hardening.

## Verification

```sh
test -z "$(gofmt -l cmd internal)"
go test -race ./...
go vet ./...
dpkg-buildpackage --build=binary --no-sign
lintian --fail-on error ../xapsd_*.deb
```

Test maintainer-script changes in a disposable Debian or Ubuntu environment, including fresh installation, upgrade with
an edited YAML file, noninteractive installation, removal, and purge. Report every check as pass, fail, or not run.

## Reviews and changes

- Use `.github/skills/code-review/SKILL.md` for formal reviews.
- Keep review findings scoped to the requested diff and rank real defects P0–P3.
- Keep lint, tests, SonarCloud, Release Drafter, and release automation as separate workflows.
- Do not commit, push, create releases, or merge without explicit authorization.
