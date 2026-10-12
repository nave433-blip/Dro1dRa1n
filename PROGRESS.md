# Dro1dRa1n — build progress

Repo arrived as README-only ("frp all in one bypass", 1 commit).
Rebuilding as a real, functional Android FRP/repair toolkit in Go (stdlib-only).

## Completed
- 2026-10-12: v0.1.0 — full toolkit built from README-only stub. Go CLI
  (stdlib-only): doctor, devices, info, frp status/guide/clear-adb,
  samsung frp/csc/debloat, xiaomi micloud, reboot, reset, shell,
  methods list/show. 10-method knowledge base with matcher. 13 tests
  passing, `go vet` + gofmt clean. Committed as 293ef56.
- 2026-10-12: Windows port — GOOS=windows build verified as PE32+ x86-64
  (`dist/dro1drain.exe`); Linux binary in `dist/dro1drain`. Added frp +
  device unit tests (17 total, all passing). Committed as 731105a.

## Audit gaps (Phase 1)
- No code at all: every claimed feature ("frp all in one bypass") was unimplemented.
- No license, no docs beyond one line, no tests, no builds.

## Competitor R&D (Phase 2) — see docs/COMPETITORS.md
## Placeholders — see PLACEHOLDERS.md

## 2026-10-12 — Gemini-loop refinement (solo rounds; Gemini API quota exhausted)
- Round 1: context-based timeouts for adb/fastboot (fixed: fastboot never killed hung processes); WaitForSerial with backoff; SamsungFRP simplified. +3 tests.
- Round 2: methods DB GuidePath field wired to android-hacking/ bundle; +7 device-specific methods (Redmi 9A, A02s EDL, Samsung combination, Oppo A15 MTK, Redmi 10C Second Space, Vivo EDL, Realme C21Y); Validate() checks; GuidePath-existence test.
- Round 3: `frp guide` auto-detects brand/Android from the connected device; doctor reports tool versions + adb server state; device.AndroidMajor(); Summary() cosmetic fix.
- Round 4: fixed real bugs — `samsung csc read`/`debloat list` treated subcommands as serials; `xiaomi micloud` now accepts sideload-state devices and explains sideload limits instead of running doomed shell commands; pickSerial validates serials against attached devices.
- Binaries rebuilt (windows exe verified PE32+). Tests: 20 passing (was 13).
