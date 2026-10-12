# Competitor R&D — Android FRP / repair tooling

Researched 2026-10-12. Only public product knowledge is used below; no
competitor code was copied.

## 1. SamFW FRP Tool (proprietary freeware, Windows-only)
The reference free tool for Samsung FRP work.
Killer features adopted into Dro1dRa1n's roadmap:
- One-click Samsung FRP remove via test mode (*#0*#) + MTP + ADB methods
  → implemented as `samsung frp` (automates the ADB side once the user
  enables test mode on-device).
- ADB/MTP device info read → `info`, `devices`.
- Factory reset via MTP/test mode → `reset` (fastboot/recovery paths).
- Change CSC → `samsung csc` guidance.
- Reboot to download/normal/recovery → `reboot`.
- Disable Knox / bloatware removal → documented in methods DB; ADB-side
  package disable automated where ADB allows.
- MediaTek module, Qualcomm EDL 9008 partition manager → NOT implemented;
  the standard backends (mtkclient=GPL-3.0, EDL tools) need Evan's licensing
  call first (see PLACEHOLDERS.md). `doctor` detects them if installed.
- Driver management utilities → `doctor` install hints.

## 2. mtkclient — https://github.com/bkerler/mtkclient (GPL-3.0)
Open-source MediaTek auth-bypass / partition read-write / bootloader unlock.
FLAGGED: GPL-3.0 — not bundled, not linked, not shelled-out-to without
Evan's explicit call. `doctor` lists it as an optional external tool.

## 3. Heimdall — https://gitlab.com/BenjaminDobell/Heimdall (LGPL)
Open-source Samsung Odin-protocol flasher.
FLAGGED: LGPL — same treatment as mtkclient. `doctor` lists it as optional.

## 4. XiaoMiToolV2 — https://github.com/francescotescari/XiaoMiToolV2 (Apache-2.0)
Xiaomi firmware management (abandoned upstream, community forks active).
Permissive license; concepts borrowed (firmware/sideload workflow guidance
in the methods DB), no code copied.

## 5. Paid boxes: UnlockTool, Chimera, Octoplus (proprietary)
Server-credit EDL/BROM flashing for latest patches. Not reproducible as
free software (server-side signing); Dro1dRa1n documents where a paid
service is genuinely the only path instead of faking it.

## What Dro1dRa1n deliberately does NOT do
- No fake "one-click" buttons for methods that require on-device steps or
  server-side signing. Every command reports honestly what it automated and
  what the user must still do by hand.
- No bundled exploits for patched vulnerabilities. The methods DB records
  patch ceilings per method so users don't waste time on dead paths.
