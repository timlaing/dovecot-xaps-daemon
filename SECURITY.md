# Security Policy

## Supported versions

Only the latest release is supported with security fixes.

| Version | Supported |
| ------- | --------- |
| Latest  | Yes       |
| Older   | No        |

## Reporting a vulnerability

Do not report vulnerabilities in the public issue tracker. Use
[GitHub private vulnerability reporting](https://github.com/timlaing/dovecot-xaps-daemon/security/advisories/new)
so the report goes directly to the maintainers.

Include the affected version, impact, reproduction steps, and any suggested mitigation. Never include APNS private
keys, credentials, device tokens, account IDs, registration databases, or private server details in public reports.

Maintainers will acknowledge the report as soon as practical, coordinate remediation and disclosure with the reporter,
and provide credit if requested.

## Scope

This policy covers xapsd, its configuration handling, HTTP listener, APNS integration, registration database, systemd
service, and Debian packaging. Vulnerabilities in Dovecot, the companion plugin, Apple services, or the host operating
system should be reported to their respective maintainers.

## Security tooling

The repository uses Go race tests, `go vet`, Dependabot, SonarQube Cloud, and GitHub security scanning. These automated
checks complement, but do not replace, security review.
