[![Test](https://github.com/timlaing/dovecot-xaps-daemon/actions/workflows/test.yml/badge.svg)](https://github.com/timlaing/dovecot-xaps-daemon/actions/workflows/test.yml)
[![Lint](https://github.com/timlaing/dovecot-xaps-daemon/actions/workflows/lint.yml/badge.svg)](https://github.com/timlaing/dovecot-xaps-daemon/actions/workflows/lint.yml)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=timlaing_dovecot-xaps-daemon&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=timlaing_dovecot-xaps-daemon)
[![Release](https://img.shields.io/github/v/release/timlaing/dovecot-xaps-daemon)](https://github.com/timlaing/dovecot-xaps-daemon/releases)
[![License](https://img.shields.io/github/license/timlaing/dovecot-xaps-daemon)](LICENSE)

iOS Push Email daemon for Dovecot
=================================

This project provides `xapsd`, the daemon used with
[dovecot-xaps-plugin](https://github.com/timlaing/dovecot-xaps-plugin) to deliver native Apple Push Notification
Service (APNS) notifications for iOS Mail.

Both components are required:

1. The Dovecot plugin implements `XAPPLEPUSHSERVICE` and forwards registrations and mail events.
2. This daemon stores device registrations and sends notifications to APNS.

Apple did not publish an XAPPLEPUSHSERVICE specification. The implementation was derived from Apple's published
Dovecot patches and requires APNS credentials that you are legally entitled to use.

Installation
============

Debian package
--------------

Tagged releases build an `xapsd` Debian package. Download it from this repository's
[Releases](https://github.com/timlaing/dovecot-xaps-daemon/releases) page and install it with APT:

```sh
sudo apt install ./xapsd_<version>_<architecture>.deb
```

The installer prompts for:

* The APNS topic.
* The APNS key ID.
* The Apple Developer team ID.

It derives Apple's standard key filename `AuthKey_<KeyID>.p8` and writes these settings to
`/etc/xapsd/xapsd.yaml`. Copy the private key downloaded from Apple into that directory:

```sh
sudo install -o root -g xapsd -m 0640 \
  AuthKey_<KeyID>.p8 /etc/xapsd/AuthKey_<KeyID>.p8
```

The package also installs the binary, systemd service, system user definition, and state-directory definition.
Review `/etc/xapsd/xapsd.yaml`, then enable and start the daemon:

```sh
sudo systemctl enable --now xapsd
sudo systemctl status xapsd
```

By default, xapsd listens on IPv6 loopback at `[::1]:11619`. The plugin must use the matching endpoint:

```dovecot
xaps_url = http://[::1]:11619
```

Build from source
-----------------

The required Go version is declared in `go.mod` (currently Go 1.25):

```sh
git clone https://github.com/timlaing/dovecot-xaps-daemon.git
cd dovecot-xaps-daemon
go test ./...
go build -trimpath -o xapsd ./cmd/xapsd
```

For a manual system installation, install the binary as `/usr/bin/xapsd`, the sample configuration as
`/etc/xapsd/xapsd.yaml`, and the files under `configs/systemd/` in their corresponding systemd directories.
The Debian package is the recommended installation method on Debian and Ubuntu.

Development
-----------

Run formatting, analysis, and the test suite (with the race detector) using the Go version declared in `go.mod`:

```sh
test -z "$(gofmt -l cmd internal)"
go vet ./...
go test -race ./...
```

Measure test coverage with:

```sh
go test -cover ./...
```

The unit-test suite covers every reachable branch across the database, configuration, APNs, and socket packages, and
reaches 85.4% aggregate statement coverage. The APNS credential directory (`/etc/xapsd/` by default) is overridable
in tests to exercise `NewApns` with both PEM and token credentials. Only code paths that would block forever
(`main()`'s HTTP listener), exit the process (`log.Fatal`), or require system state (the package-managed
`/etc/xapsd/xapsd.yaml`) are not exercised by unit tests; those paths are covered by the packaging and release
workflows instead.

Build a Debian package
----------------------

The package uses standard debhelper metadata under `debian/`. On Debian or Ubuntu, install the build dependencies and
build the binary package with:

```sh
sudo apt-get install debhelper golang-any lintian
dpkg-buildpackage --build=binary --no-sign
lintian --fail-on error ../xapsd_*.deb
```

Debhelper generates the package checksums and systemd maintainer-script sections. Installation enables and starts
`xapsd`; upgrades restart it, while removal stops and disables it. Existing `/etc/xapsd/xapsd.yaml` values are
preserved during upgrades.

Configuration
=============

The main configuration is `/etc/xapsd/xapsd.yaml`. Important settings include:

* `listenAddr` and `port`: HTTP listener used by the Dovecot plugin.
* `databaseFile`: persistent device-registration database.
* `keyFileP8`: APNS private-key filename under `/etc/xapsd`.
* `keyFileTopic`: APNS topic.
* `keyFileKeyId`: APNS key ID.
* `keyFileTeamId`: Apple Developer team ID.

Never commit the APNS private key.

Setting up devices
==================

After both services are running, reconnect the iOS Mail account so it discovers `XAPPLEPUSHSERVICE`. Toggling
Airplane Mode or restarting the device forces a reconnect. Successful registration produces Dovecot log messages
ending with an HTTP `200 OK` response from xapsd.

Troubleshooting
===============

Check that xapsd is running and listening on the same address configured in Dovecot:

```sh
sudo systemctl status xapsd
sudo journalctl -u xapsd -n 100 --no-pager
sudo ss -ltnp | grep 11619
```

If xapsd listens on `[::1]:11619`, configuring the plugin with `http://localhost:11619` may make Dovecot attempt
`127.0.0.1` instead. Use `http://[::1]:11619` explicitly.

If two devices share the same account data after device migration and only one receives notifications, remove and
re-add the mail account on one device. The registrations are stored in `/var/lib/xapsd/database.json`.

Report daemon issues in this repository's [issue tracker](https://github.com/timlaing/dovecot-xaps-daemon/issues).

## Privacy

For each mail event, xapsd sends Apple a TLS-protected request containing the APNS device token, account ID, and
certificate or key topic. Apple can observe the sending server's IP address and notification timing. The email
message content is not sent to Apple by this software.

Acknowledgements
================

This repository is a maintained fork of [freswa/dovecot-xaps-daemon](https://github.com/freswa/dovecot-xaps-daemon),
maintained by Frederik Schwan, which is itself based on the original
[st3fan/dovecot-xaps-daemon](https://github.com/st3fan/dovecot-xaps-daemon) by Stefan Arentz. Their design,
implementation, maintenance, and the work of all contributors made this fork possible.
