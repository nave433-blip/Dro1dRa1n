// Package doctor checks host dependencies and prints install hints.
package doctor

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/nave433-blip/Dro1dRa1n/internal/ui"
)

// Tool describes one external dependency.
type Tool struct {
	Binary      string
	Description string
	Required    bool
	Hint        map[string]string // os -> install hint
}

var tools = []Tool{
	{
		Binary:      "adb",
		Description: "device comms (Android platform-tools)",
		Required:    true,
		Hint: map[string]string{
			"linux":   "sudo apt install android-tools-adb  OR  download platform-tools from developer.android.com",
			"darwin":  "brew install android-platform-tools",
			"windows": "winget install Google.PlatformTools  (adds adb to PATH)",
		},
	},
	{
		Binary:      "fastboot",
		Description: "bootloader comms (Android platform-tools)",
		Required:    true,
		Hint: map[string]string{
			"linux":   "sudo apt install android-tools-fastboot  OR  platform-tools from developer.android.com",
			"darwin":  "brew install android-platform-tools",
			"windows": "winget install Google.PlatformTools  (adds fastboot to PATH)",
		},
	},
	{
		Binary:      "heimdall",
		Description: "Samsung Odin-protocol flashing (optional, LGPL — needs your call)",
		Hint: map[string]string{
			"linux":   "sudo apt install heimdall-flash",
			"darwin":  "brew install heimdall",
			"windows": "installer from glassechidna.com.au/heimdall, then bind driver via Zadig",
		},
	},
	{
		Binary:      "mtkclient",
		Description: "MediaTek auth-bypass/flashing (optional, GPL-3.0 — needs your call)",
		Hint: map[string]string{
			"linux":   "pip install --user mtkclient",
			"darwin":  "pip install --user mtkclient",
			"windows": "pip install mtkclient  (install UsbDk driver first)",
		},
	},
}

// Result is the outcome of a dependency check.
type Result struct {
	Tool    Tool
	Present bool
	Path    string
}

// Check runs the dependency scan.
func Check() []Result {
	var out []Result
	for _, t := range tools {
		p, err := exec.LookPath(t.Binary)
		out = append(out, Result{Tool: t, Present: err == nil, Path: p})
	}
	return out
}

// have is a tiny helper used by report code.
func have(binary string) bool {
	_, err := exec.LookPath(binary)
	return err == nil
}

// HasADB is exported for command gating.
func HasADB() bool { return have("adb") }

// Report prints the check results and returns false when a required tool is missing.
func Report() bool {
	ui.Header("dro1drain doctor")
	fmt.Printf("host: %s/%s\n\n", runtime.GOOS, runtime.GOARCH)
	ok := true
	for _, r := range Check() {
		tag := "required"
		if !r.Tool.Required {
			tag = "optional"
		}
		if r.Present {
			ver := toolVersion(r.Tool.Binary)
			ui.Found(fmt.Sprintf("%-12s %s  [%s]%s", r.Tool.Binary, r.Tool.Description, tag, ver))
		} else {
			ui.Missing(fmt.Sprintf("%-12s %s  [%s]", r.Tool.Binary, r.Tool.Description, tag))
			if hint, found := r.Tool.Hint[runtime.GOOS]; found {
				fmt.Printf("             -> %s\n", hint)
			}
			if r.Tool.Required {
				ok = false
			}
		}
	}
	if have("adb") {
		reportADBServer()
	}
	// Windows driver note
	if runtime.GOOS == "windows" {
		fmt.Println()
		ui.Warn("Windows needs device drivers: Samsung USB driver for Galaxy devices,")
		ui.Warn("Google USB driver for Pixels (via Android Studio SDK Manager).")
	}
	// udev note on linux
	if runtime.GOOS == "linux" {
		fmt.Println()
		ui.Warn("Linux: add udev rules for your device or run with sudo, else adb may show 'no permissions'.")
	}
	fmt.Println()
	if ok {
		ui.Ok("all required dependencies present.")
	} else {
		ui.Warn("missing required dependencies — install the hints above, then re-run.")
	}
	return ok
}

// toolVersion returns a short " (v...)" suffix from "<bin> --version".
func toolVersion(bin string) string {
	p, err := exec.LookPath(bin)
	if err != nil {
		return ""
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, p, "--version").Output()
	if err != nil {
		return ""
	}
	line := strings.TrimSpace(strings.SplitN(string(out), "\n", 2)[0])
	// adb prints "Android Debug Bridge version 1.0.41"; fastboot prints a version line too.
	if line == "" {
		return ""
	}
	if len(line) > 60 {
		line = line[:60]
	}
	return " — " + line
}

// reportADBServer checks whether the adb server answers and sees devices.
func reportADBServer() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "adb", "devices").Output()
	if err != nil {
		ui.Warn("adb server not responding — try: adb kill-server && adb start-server")
		return
	}
	n := 0
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "List of") || strings.HasPrefix(line, "*") {
			continue
		}
		n++
	}
	if n == 0 {
		ui.Info("adb server is up, no devices attached.")
	} else {
		ui.Ok(fmt.Sprintf("adb server is up, %d device(s) visible.", n))
	}
}

// DriverHint returns a one-line driver hint for the current OS.
func DriverHint() string {
	switch runtime.GOOS {
	case "windows":
		return "install Samsung USB driver / Google USB driver, then replug the device"
	case "darwin":
		return "no driver needed on macOS"
	default:
		return "check udev rules / replug; try sudo adb kill-server; sudo adb start-server"
	}
}
