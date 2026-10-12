# Changelog

## Unreleased — dependency audit

- Dependency audit: Go side verified fully stdlib-only — clean build with empty
  GOCACHE and GOPROXY=off, `go vet` and `gofmt` clean.
- `doctor` Windows hints now say `winget install Google.PlatformTools` (matches
  the installer); Samsung/USB-driver guidance kept per-OS.
- `tools/win/install.ps1` now detects a bundled `Drivers\` folder (see
  `docs/PACKAGING.md`) and offers to install the Samsung USB driver (admin),
  falling back to a repo/URL pointer when the zip omits it.
- New `tools/install.sh`: Linux/macOS installer — builds, installs to
  `~/.local/bin`, checks platform-tools with per-OS hints. Tested end to end.
- New `internal/doctor/doctor_test.go`: tool-table completeness (every required
  tool has hints for linux/darwin/windows; hint keys are valid GOOS values;
  adb+fastboot present), hint actionability, driver-hint coverage.
- New `docs/PACKAGING.md`: release-zip layout including the Samsung driver.
- License calls still pending Evan: mtkclient (GPL-3.0) and Heimdall (LGPL) are
  detected-only, never bundled or shelled out to.

## Unreleased — v0.2.0 (draft release notes)
Refinement wave: 19 tests green (was 13), Windows installer, deeper diagnostics.

- Exec plumbing: adb/fastboot timeouts rebuilt on `context` (fixes a real bug —
  fastboot never killed hung processes); `WaitForSerial` with backoff replaces
  the manual re-scan in the Samsung FRP flow.
- Methods DB: new `GuidePath` field links 11 procedures to the bundled
  `android-hacking/` guides; 7 new device-specific methods (Redmi 9A, Galaxy A02s
  EDL, Samsung combination-file, Oppo A15 MTK, Redmi 10C Second Space, Vivo EDL,
  Realme C21Y, Oppo A16K clone-phone, Redmi A5 Unisoc, Redmi 13C SIM-lock,
  Xperia MXTP, OneUI3 recovery, locked-phone model lookup, build.prop tweaks).
- `frp guide` auto-detects brand + Android version from the connected device.
- `info` now reports OEM-unlock permission, bootloader lock state, Knox warranty
  bit, and Xiaomi anti-rollback index.
- `doctor` reports adb/fastboot versions and adb-server state.
- Fixed: `samsung csc read` / `debloat list` treated the subcommand as a serial;
  `xiaomi micloud` now accepts sideload-state devices and explains sideload
  limits honestly instead of running doomed shell commands; `pickSerial`
  validates serials against attached devices.
- Windows: `tools/win/install.ps1` — one-command user-level installer
  (no admin; adds to user PATH; platform-tools/driver hints).

## v0.1.0 — 2026-10-12
Initial real release. Repo previously contained only a README stub.

- Go CLI (`dro1drain`), stdlib-only, builds on Linux/macOS/Windows.
- `doctor`: dependency check (adb, fastboot, drivers) with per-OS install hints.
- `devices`: list ADB + fastboot devices.
- `info`: full device info dump (brand, model, Android version, security patch,
  build fingerprint, bootloader/verified-boot state).
- `frp status`: FRP lock-state assessment from observable device state.
- `frp guide`: interactive method picker — matches brand/Android version/patch
  against the built-in methods database and prints step-by-step procedures.
- `samsung frp`: automates the ADB side of the documented Samsung test-mode
  (*#0*#) FRP method once the user enables it on-device.
- `samsung csc`: read/change CSC guidance (MTP and root paths).
- `xiaomi micloud`: Xiaomi ADB-side MiCloud/Account steps.
- `reboot`: reboot to system/recovery/bootloader/download/edl.
- `reset`: guided factory reset via fastboot/recovery (confirms first).
- `shell`: adb shell passthrough.
- `methods list/show`: browse the built-in repair-methods database.
- Windows: `GOOS=windows go build` produces `dro1drain.exe`.
