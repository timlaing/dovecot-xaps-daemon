# Pull Request

<!-- Thank you for contributing to dovecot-xaps-daemon. Please complete this template. -->

## Summary

<!-- Briefly describe the change and the problem it solves. -->

## Type of change

- [ ] 🚀 Feature (`feature` / `enhancement`)
- [ ] 🐛 Bug fix (`bug` / `fix`)
- [ ] 🛠 Maintenance (`maintenance` / `dependencies`)
- [ ] 📚 Documentation (`documentation`)

## Related issues / pull requests

<!-- Link related work, for example: Fixes #123. -->

## Changes

<!-- Explain the implementation and any configuration, APNS, or compatibility impact. -->

## Verification

- [ ] This pull request contains one self-contained change.
- [ ] Maintainers are allowed to edit this pull request.
- [ ] Changed Go files were formatted with `gofmt`.
- [ ] `go test ./...` passes.
- [ ] `go build ./cmd/xapsd` succeeds.
- [ ] `dpkg-buildpackage --build=binary --no-sign` and `lintian --fail-on error ../xapsd_*.deb` pass, if packaging changed.
- [ ] Documentation and sample configuration were updated where needed.
- [ ] No credentials, private keys, device tokens, or private mail data are included.

<!-- List any checks that were not run and explain why. -->
