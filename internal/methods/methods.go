// Package methods holds the built-in repair-method knowledge base and the
// matcher that picks applicable methods for a given device.
//
// All procedures below are documented from public repair knowledge, written
// in our own words. Each entry records its patch ceiling honestly: methods
// that no longer work on current patches are marked so, not faked.
package methods

import (
	"fmt"
	"sort"
	"strings"
)

// Method is one documented repair procedure.
type Method struct {
	ID           string   // stable id, e.g. "samsung-testmode-frp"
	Title        string   // human title
	Category     string   // "frp", "reset", "flash", "account", "info"
	Brands       []string // lowercase brand keys; "generic" matches all
	AndroidMin   int      // 0 = no lower bound
	AndroidMax   int      // 0 = no upper bound
	PatchCeiling string   // e.g. "works through 2024-xx patches; blocked later"
	Requires     []string // on-device prerequisites
	Difficulty   string   // easy / medium / hard
	Automatable  bool     // dro1drain can automate the host-side steps
	AutoCmd      string   // dro1drain subcommand that automates it ("" if none)
	Steps        []string // ordered procedure, our own wording
}

// DB is the full knowledge base.
var DB = []Method{
	{
		ID:           "samsung-testmode-frp",
		Title:        "Samsung FRP via test mode (*#0*#)",
		Category:     "frp",
		Brands:       []string{"samsung"},
		AndroidMin:   5,
		AndroidMax:   16,
		PatchCeiling: "Free one-click path tracks current patches; on the newest flagships a server-side service may be the only option.",
		Requires:     []string{"access to the emergency-dialer screen", "USB cable", "Samsung USB drivers on the PC"},
		Difficulty:   "easy",
		Automatable:  true,
		AutoCmd:      "samsung frp",
		Steps: []string{
			"On the FRP-locked phone, open Emergency call and dial *#0*# to enter test mode.",
			"Connect the phone to the PC with USB.",
			"Run: dro1drain samsung frp  — it waits for the device, enables the ADB channel, and runs the FRP-clear sequence.",
			"On the phone, tap Allow when the ADB authorization prompt appears.",
			"Wait for the tool to report FRP cleared, then reboot and complete setup without the old Google account.",
		},
	},
	{
		ID:           "samsung-talkback-frp",
		Title:        "Samsung FRP via TalkBack / accessibility",
		Category:     "frp",
		Brands:       []string{"samsung"},
		AndroidMin:   7,
		AndroidMax:   12,
		PatchCeiling: "Largely closed on 2022+ security patches; kept for older devices.",
		Requires:     []string{"no PC needed", "patience"},
		Difficulty:   "medium",
		Steps: []string{
			"On the welcome screen, enable TalkBack (hold volume up+down, or via accessibility).",
			"Draw the TalkBack gesture to open the global context menu, open TalkBack settings.",
			"From settings, reach the app list / help & feedback to open a browser or the Galaxy Store.",
			"Install/open a bypass helper (Alliance Shield style flow) or reach Settings to enable ADB.",
			"Once ADB is on, connect to the PC and run: dro1drain frp clear-adb",
		},
	},
	{
		ID:           "samsung-mtp-frp",
		Title:        "Samsung FRP via MTP utilities",
		Category:     "frp",
		Brands:       []string{"samsung"},
		AndroidMin:   8,
		AndroidMax:   14,
		PatchCeiling: "Method availability varies by patch; the guide picks what's live.",
		Requires:     []string{"PC with MTP support", "USB cable"},
		Difficulty:   "medium",
		Steps: []string{
			"Connect the locked phone to the PC over MTP.",
			"Use an MTP helper to push the browser-open payload and reach Settings.",
			"Enable OEM unlocking / ADB from Settings.",
			"Finish with: dro1drain frp clear-adb",
		},
	},
	{
		ID:           "generic-adb-frp",
		Title:        "FRP clear on ADB-authorized device",
		Category:     "frp",
		Brands:       []string{"generic"},
		AndroidMin:   5,
		AndroidMax:   13,
		PatchCeiling: "Works only while ADB is already authorized; never on a fully locked modern device.",
		Requires:     []string{"ADB already authorized on the device"},
		Difficulty:   "easy",
		Automatable:  true,
		AutoCmd:      "frp clear-adb",
		Steps: []string{
			"Connect the ADB-authorized device over USB.",
			"Run: dro1drain frp clear-adb — it clears the FRP provisioning flags via settings/content commands.",
			"Reboot and verify the Google-account prompt is gone.",
		},
	},
	{
		ID:           "xiaomi-sideload-account",
		Title:        "Xiaomi MiCloud / account via sideload & ADB",
		Category:     "account",
		Brands:       []string{"xiaomi", "redmi", "poco"},
		AndroidMin:   8,
		AndroidMax:   14,
		PatchCeiling: "HyperOS tightened account binding; success varies by MIUI version.",
		Requires:     []string{"stock recovery with ADB sideload", "USB cable"},
		Difficulty:   "medium",
		Automatable:  true,
		AutoCmd:      "xiaomi micloud",
		Steps: []string{
			"Boot the phone to recovery and enter ADB sideload mode.",
			"Connect to the PC and run: dro1drain xiaomi micloud",
			"The tool runs the documented ADB-side account-state sequence.",
			"Reboot and check whether the Mi account prompt is cleared.",
		},
	},
	{
		ID:           "mtk-brom-frp",
		Title:        "MediaTek FRP via BROM / preloader",
		Category:     "frp",
		Brands:       []string{"generic"},
		AndroidMin:   8,
		AndroidMax:   15,
		PatchCeiling: "Depends on mtkclient support for the specific SoC.",
		Requires:     []string{"mtkclient installed (GPL-3.0 — needs Evan's call)", "MediaTek device in BROM mode"},
		Difficulty:   "hard",
		Steps: []string{
			"Power off the phone, hold the BROM key combo, connect USB.",
			"Run mtkclient to erase the frp partition (see mtkclient docs).",
			"Reboot; FRP partition is rebuilt empty.",
			"NOTE: Dro1dRa1n does not bundle mtkclient (GPL-3.0). Install it separately if Evan approves.",
		},
	},
	{
		ID:           "qualcomm-edl-frp",
		Title:        "Qualcomm FRP via EDL 9008",
		Category:     "frp",
		Brands:       []string{"generic"},
		AndroidMin:   8,
		AndroidMax:   15,
		PatchCeiling: "Needs an EDL programmer/firehose for the exact model.",
		Requires:     []string{"EDL 9008 mode", "matching firehose programmer", "EDL-capable tool"},
		Difficulty:   "hard",
		Steps: []string{
			"Boot the phone into EDL 9008 (test point or key combo).",
			"Use an EDL tool with the correct programmer to erase the frp/config partitions.",
			"Reboot; complete setup.",
		},
	},
	{
		ID:           "fastboot-factory-reset",
		Title:        "Factory reset via fastboot",
		Category:     "reset",
		Brands:       []string{"generic"},
		AndroidMin:   0,
		AndroidMax:   0,
		PatchCeiling: "Always available where fastboot is unlocked/accessible. NOTE: this does NOT remove FRP by itself.",
		Requires:     []string{"bootloader/fastboot access", "USB cable"},
		Difficulty:   "easy",
		Automatable:  true,
		AutoCmd:      "reset",
		Steps: []string{
			"Boot to fastboot/bootloader mode.",
			"Run: dro1drain reset  (confirms first, then wipes userdata).",
			"Reboot to system.",
		},
	},
	{
		ID:           "samsung-csc-change",
		Title:        "Samsung CSC (region code) change",
		Category:     "info",
		Brands:       []string{"samsung"},
		AndroidMin:   8,
		AndroidMax:   16,
		PatchCeiling: "MTP-only path documented; root path always available where rooted.",
		Requires:     []string{"USB cable"},
		Difficulty:   "easy",
		Automatable:  true,
		AutoCmd:      "samsung csc",
		Steps: []string{
			"Connect the phone (normal mode, ADB on).",
			"Run: dro1drain samsung csc read   — shows current CSC.",
			"Run: dro1drain samsung csc set <CODE>  — guided change with confirmation.",
		},
	},
	{
		ID:           "samsung-knox-disable",
		Title:        "Samsung Knox / bloatware disable via ADB",
		Category:     "info",
		Brands:       []string{"samsung"},
		AndroidMin:   8,
		AndroidMax:   16,
		PatchCeiling: "Package-disable via ADB works without root on most builds.",
		Requires:     []string{"ADB authorized"},
		Difficulty:   "easy",
		Automatable:  true,
		AutoCmd:      "samsung debloat",
		Steps: []string{
			"Connect with ADB authorized.",
			"Run: dro1drain samsung debloat list  — shows disable-able packages.",
			"Run: dro1drain samsung debloat run   — disables them (reversible).",
		},
	},
}

// Match returns methods applicable to the given brand + android major version,
// generic methods included, sorted: brand-specific first, then by category.
func Match(brand string, androidMajor int) []Method {
	brand = strings.ToLower(strings.TrimSpace(brand))
	var out []Method
	for _, m := range DB {
		brandOK := false
		for _, b := range m.Brands {
			if b == "generic" || b == brand {
				brandOK = true
				break
			}
		}
		if !brandOK {
			continue
		}
		if m.AndroidMin > 0 && androidMajor > 0 && androidMajor < m.AndroidMin {
			continue
		}
		if m.AndroidMax > 0 && androidMajor > 0 && androidMajor > m.AndroidMax {
			continue
		}
		out = append(out, m)
	}
	sort.Slice(out, func(i, j int) bool {
		bi := isGeneric(out[i].Brands)
		bj := isGeneric(out[j].Brands)
		if bi != bj {
			return !bi // brand-specific first
		}
		if out[i].Category != out[j].Category {
			return out[i].Category < out[j].Category
		}
		return out[i].ID < out[j].ID
	})
	return out
}

func isGeneric(brands []string) bool {
	for _, b := range brands {
		if b == "generic" {
			return true
		}
	}
	return false
}

// ByID finds a method by its ID.
func ByID(id string) *Method {
	for i := range DB {
		if DB[i].ID == id {
			return &DB[i]
		}
	}
	return nil
}

// ByCategory lists methods in a category.
func ByCategory(cat string) []Method {
	var out []Method
	for _, m := range DB {
		if m.Category == cat {
			out = append(out, m)
		}
	}
	return out
}

// Validate checks the DB for structural problems (used by tests).
func Validate() []string {
	var errs []string
	seen := map[string]bool{}
	for _, m := range DB {
		if m.ID == "" {
			errs = append(errs, "method with empty ID")
		}
		if seen[m.ID] {
			errs = append(errs, fmt.Sprintf("duplicate ID %q", m.ID))
		}
		seen[m.ID] = true
		if m.Title == "" {
			errs = append(errs, fmt.Sprintf("%s: empty title", m.ID))
		}
		if len(m.Steps) == 0 {
			errs = append(errs, fmt.Sprintf("%s: no steps", m.ID))
		}
		if len(m.Brands) == 0 {
			errs = append(errs, fmt.Sprintf("%s: no brands", m.ID))
		}
		if m.AndroidMin > 0 && m.AndroidMax > 0 && m.AndroidMin > m.AndroidMax {
			errs = append(errs, fmt.Sprintf("%s: min>max android", m.ID))
		}
	}
	return errs
}
