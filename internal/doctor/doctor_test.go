package doctor

import (
	"runtime"
	"strings"
	"testing"
)

// TestToolTableComplete mirrors the TestGuidePathsExist idea: every external
// dependency in the tool table must have a description, a valid required flag,
// and a working install hint for each supported OS.
func TestToolTableComplete(t *testing.T) {
	oses := []string{"linux", "darwin", "windows"}
	if len(tools) == 0 {
		t.Fatal("tool table is empty — nothing to audit")
	}
	seen := map[string]bool{}
	for _, tl := range tools {
		if tl.Binary == "" {
			t.Errorf("tool entry with empty Binary")
			continue
		}
		if seen[tl.Binary] {
			t.Errorf("duplicate tool entry %q", tl.Binary)
		}
		seen[tl.Binary] = true
		if tl.Description == "" {
			t.Errorf("%s: empty Description", tl.Binary)
		}
		// Required tools must have a hint for every supported OS.
		if tl.Required {
			for _, os := range oses {
				h, ok := tl.Hint[os]
				if !ok || strings.TrimSpace(h) == "" {
					t.Errorf("%s: required tool missing install hint for %s", tl.Binary, os)
				}
			}
		}
		// Hint keys must be valid GOOS values.
		for os := range tl.Hint {
			switch os {
			case "linux", "darwin", "windows":
			default:
				t.Errorf("%s: hint for unknown OS %q (want linux/darwin/windows)", tl.Binary, os)
			}
		}
	}
	// adb and fastboot are the hard requirements of the whole CLI.
	for _, req := range []string{"adb", "fastboot"} {
		if !seen[req] {
			t.Errorf("required tool %q missing from tool table", req)
		}
	}
	// The table must agree with HasADB on the current host.
	results := Check()
	if len(results) != len(tools) {
		t.Fatalf("Check returned %d results for %d tools", len(results), len(tools))
	}
	for _, r := range results {
		if r.Tool.Binary == "adb" && r.Present != HasADB() {
			t.Errorf("Check/HasADB disagree on adb presence")
		}
	}
}

// TestHintsPointSomewhere ensures every hint names an actionable installer
// mechanism (package manager name, winget/choco/brew/apt, or a URL/domain).
func TestHintsPointSomewhere(t *testing.T) {
	for _, tl := range tools {
		for os, h := range tl.Hint {
			l := strings.ToLower(h)
			actionable := strings.Contains(l, "apt") ||
				strings.Contains(l, "brew") ||
				strings.Contains(l, "winget") ||
				strings.Contains(l, "choco") ||
				strings.Contains(l, "pip") ||
				strings.Contains(l, "http") ||
				strings.Contains(l, ".com") ||
				strings.Contains(l, "developer.android.com")
			if !actionable {
				t.Errorf("%s [%s]: hint %q does not name an installer or source", tl.Binary, os, h)
			}
		}
	}
}

// TestDriverHint covers the platform driver one-liner used by frp flows.
func TestDriverHint(t *testing.T) {
	h := DriverHint()
	if strings.TrimSpace(h) == "" {
		t.Error("DriverHint returned empty string")
	}
	switch runtime.GOOS {
	case "windows":
		if !strings.Contains(strings.ToLower(h), "samsung") {
			t.Errorf("windows DriverHint should mention Samsung driver, got %q", h)
		}
	case "linux":
		if !strings.Contains(strings.ToLower(h), "udev") {
			t.Errorf("linux DriverHint should mention udev, got %q", h)
		}
	}
}
