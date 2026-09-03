# Contributing to dovecot-xaps-daemon

Thank you for contributing. This repository is maintained at
[github.com/timlaing/dovecot-xaps-daemon](https://github.com/timlaing/dovecot-xaps-daemon) and builds on the work of
[freswa/dovecot-xaps-daemon](https://github.com/freswa/dovecot-xaps-daemon),
[st3fan/dovecot-xaps-daemon](https://github.com/st3fan/dovecot-xaps-daemon), and their contributors.

## Before you start

Search the existing issues and pull requests before opening a new one. Bugs in the Dovecot plugin belong in the
[dovecot-xaps-plugin issue tracker](https://github.com/timlaing/dovecot-xaps-plugin/issues).

## Development setup

Install the Go version declared in `go.mod`, then clone and verify the repository:

```sh
git clone https://github.com/timlaing/dovecot-xaps-daemon.git
cd dovecot-xaps-daemon
go test ./...
go build ./cmd/xapsd
```

Use a test configuration and non-production APNS credentials during development.

## Making changes

1. Fork the repository and create a focused branch such as `feature/my-change` or `fix/my-bug`.
2. Keep each pull request limited to one self-contained change.
3. Format Go changes with `gofmt`.
4. Add or update tests for behavioral changes.
5. Update the sample configuration and documentation when configuration changes.
6. Do not commit APNS private keys, credentials, device tokens, account IDs, databases, or private server details.

## Verification

Before submitting a pull request, run:

```sh
gofmt -w <changed-go-files>
go test ./...
go build ./cmd/xapsd
sh -n debian/config debian/postinst debian/postrm
dpkg-buildpackage --build=binary --no-sign
lintian --fail-on error ../xapsd_*.deb
```

Review the formatted files before committing. If packaging changed, verify the tagged-release workflow or build and
inspect the Debian package on a Debian or Ubuntu system. Describe any checks you could not run and why.

## Pull requests

Complete the pull-request template, link related issues with `Fixes #123` where appropriate, and allow maintainer
edits. Pull-request titles may use conventional prefixes such as `feat:`, `fix:`, `docs:`, or `chore:` so
Release Drafter can categorize the change.
