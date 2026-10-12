#!/bin/sh
# Dro1dRa1n installer for Linux and macOS — no root needed.
# Builds dro1drain and installs it to ~/.local/bin, then checks platform-tools.
# Run from the repo root:  sh tools/install.sh
set -e

cd "$(dirname "$0")/.."

if ! command -v go >/dev/null 2>&1; then
    echo "Go is not installed. Install Go 1.23+ first:"
    echo "  Linux:  download from https://go.dev/dl/"
    echo "  macOS:  brew install go"
    exit 1
fi

echo "Building dro1drain..."
go build -o dro1drain .

DEST="$HOME/.local/bin"
mkdir -p "$DEST"
cp dro1drain "$DEST/dro1drain"
echo "Installed to $DEST/dro1drain"

case ":$PATH:" in
    *":$DEST:"*) ;;
    *) echo "NOTE: $DEST is not on your PATH. Add this to your shell profile:"
       echo "  export PATH=\"\$HOME/.local/bin:\$PATH\"" ;;
esac

echo ""
for bin in adb fastboot; do
    if command -v "$bin" >/dev/null 2>&1; then
        echo "$bin found: $(command -v "$bin")"
    else
        echo "$bin MISSING — install it:"
        case "$(uname -s)" in
            Linux)  echo "  sudo apt install android-tools-adb android-tools-fastboot" ;;
            Darwin) echo "  brew install android-platform-tools" ;;
        esac
        echo "  or download platform-tools from https://developer.android.com/tools/releases/platform-tools"
    fi
done

if [ "$(uname -s)" = "Linux" ]; then
    echo ""
    echo "Linux note: add udev rules for your device (or use sudo) so adb can see it."
fi

echo ""
echo "Next: run  dro1drain doctor"
