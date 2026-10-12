# Dro1dRa1n

FRP all-in-one bypass toolkit for Android repair — device detection, FRP
status assessment, guided repair methods, and host-side automation of
documented procedures. Samsung, Xiaomi, MediaTek, Qualcomm, and generic
Android.

> **Repair use only.** Use only on devices you own or have explicit permission
> to work on. See `docs/COMPETITORS.md` for the honesty policy: no fake
> one-click buttons, no bundled exploits for patched vulnerabilities.

## Install

Requires Go 1.23+ to build, plus Android platform-tools to talk to devices.

```sh
git clone https://github.com/nave433-blip/Dro1dRa1n.git
cd Dro1dRa1n
sh tools/install.sh     # builds, installs to ~/.local/bin, checks platform-tools
./dro1drain doctor     # check dependencies, prints install hints
```

Windows: `GOOS=windows go build -o dro1drain.exe .` (see `dist/`).

## Usage

```sh
dro1drain doctor                  # dependency check
dro1drain devices                 # list ADB + fastboot devices
dro1drain info [serial]           # full device info dump
dro1drain frp status [serial]     # assess FRP lock state
dro1drain frp guide               # interactive method picker
dro1drain frp clear-adb [serial]  # ADB-side FRP clear (authorized ADB only)
dro1drain samsung frp             # Samsung test-mode (*#0*#) FRP flow
dro1drain samsung csc read        # read Samsung CSC region code
dro1drain samsung debloat list    # list disable-able Samsung packages
dro1drain xiaomi micloud          # Xiaomi ADB-side account sequence
dro1drain reboot recovery         # reboot to system|recovery|bootloader|download|edl
dro1drain reset [serial]          # factory reset via fastboot (asks first)
dro1drain methods list            # browse the repair-methods database
dro1drain methods show <id>       # full procedure for one method
```

## How the Samsung FRP flow works

1. On the locked phone: Emergency call → dial `*#0*#` (test mode).
2. Plug in USB.
3. Run `dro1drain samsung frp` — it waits for the device, brings up the ADB
   channel, and runs the documented FRP-clear ADB sequence, reporting every
   step. Tap **Allow** on the phone if prompted.

## Layout

- `main.go` — CLI dispatcher and commands
- `internal/adb` — adb wrapper + output parsers
- `internal/fastboot` — fastboot wrapper
- `internal/device` — structured device info collection
- `internal/doctor` — dependency checks with per-OS install hints
- `internal/methods` — repair-methods knowledge base + matcher
- `internal/frp` — FRP assessment + host-side automation
- `internal/ui` — terminal output helpers
- `docs/COMPETITORS.md` — competitor research + licensing flags

## Tests

```sh
go test ./...
```

## License

MIT — see `LICENSE`.

## Bundled: android-hacking collection

The `android-hacking/` directory is Evan's Android_Hacking repo merged in whole: device-specific FRP bypass guides (Redmi, Oppo, Realme, Vivo, Xperia...), lockscreen-removal guides, USB driver packs (Samsung, QDLoader, Unisoc/SPRD), Windows flash tools (Odin, SP Flash Tool, QPST, ResearchDownload), Magisk modules, build.prop tweaks, modded APKs, and a useful-ADB-commands reference. It backs the `methods` knowledge base with real-world procedures and files.
