# Changelog

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
