# Dro1dRa1n release packaging

## Windows zip (`dro1drain-vX.Y.Z-windows.zip`) layout

```
dro1drain-windows/
  dro1drain.exe                 # GOOS=windows go build -o dro1drain.exe .
  install.ps1                   # one-command user-level installer (from tools/win/)
  Drivers/
    SAMSUNG-USB-Driver-for-Mobile-Phones.exe   # copy from android-hacking/Drivers/
  README.txt                    # copy of README.md
```

`install.ps1` detects `Drivers\SAMSUNG-USB-Driver-for-Mobile-Phones.exe` next to
itself and offers to install it (Galaxy devices need it; requires admin).
If the driver is omitted from the zip, the installer falls back to pointing the
user at `android-hacking/Drivers/` in the GitHub repo.

## Unix (Linux/macOS)

No binary release artifact: users run `sh tools/install.sh` from a clone, which
builds with Go 1.23+, installs to `~/.local/bin`, and checks platform-tools.
