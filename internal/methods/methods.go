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
	GuidePath    string   // optional path into android-hacking/ bundle (repo-root-relative), e.g. "android-hacking/Bypassing_FRP_Lock/Redmi_9A/readme.md"
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
		GuidePath:    "android-hacking/Useful_ADB_Commands/readme.md",
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
		GuidePath:    "android-hacking/Bypassing_FRP_Lock/Mediatek.md",
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
	{
		ID:           "redmi-9a-frp",
		Title:        "Redmi 9A FRP via SD card + QuickShortcutMaker",
		Category:     "frp",
		Brands:       []string{"xiaomi", "redmi"},
		AndroidMin:   10,
		AndroidMax:   12,
		PatchCeiling: "Documented on MIUI 12-era builds; newer patches may close the Share-to-Notes path.",
		Requires:     []string{"microSD card", "second phone or PC with card reader", "QuickShortcutMaker APK"},
		Difficulty:   "medium",
		GuidePath:    "android-hacking/Bypassing_FRP_Lock/Redmi_9A/readme.md",
		Steps: []string{
			"Copy the QuickShortcutMaker APK onto a microSD card and insert it into the phone.",
			"Power on, connect to Wi-Fi, tap Add network.",
			"In the SSID field type random text, long-press, Share to Notes, then share as image and save — this opens Gallery.",
			"Take any photo, open its details, tap the file path to launch the file manager.",
			"Browse to the SD card, install and open QuickShortcutMaker.",
			"Search for Settings, tap Try, then set a new lock screen (PIN/pattern/password) under Security.",
			"Reboot: at the Google-account prompt, enter the lock screen you just created.",
		},
	},
	{
		ID:           "samsung-a02s-edl",
		Title:        "Samsung Galaxy A02s FRP via EDL 9008 + QPST",
		Category:     "frp",
		Brands:       []string{"samsung"},
		AndroidMin:   10,
		AndroidMax:   12,
		PatchCeiling: "Hardware-level method; works regardless of patch but needs disassembly.",
		Requires:     []string{"EDL test-point access (back cover off)", "matching firehose programmer (bundled in guide)", "QPST", "QDLoader HS-USB driver"},
		Difficulty:   "hard",
		GuidePath:    "android-hacking/Bypassing_FRP_Lock/Galaxy_A02S/readme.md",
		Steps: []string{
			"Factory-reset the phone first if you have not already.",
			"Power off, remove the back cover and the black plastic shield over the mainboard.",
			"Short the two EDL test points near the camera module while connecting USB to enter 9008 mode.",
			"Install the QDLoader HS-USB driver; confirm the 9008 port appears in Device Manager.",
			"Use QPST with the bundled prog-emmc-firehose-8953-ddr.mbn programmer per the guide to erase the FRP partition.",
			"Reboot and complete setup.",
		},
	},
	{
		ID:           "samsung-combination-frp",
		Title:        "Samsung FRP via combination firmware (Odin)",
		Category:     "frp",
		Brands:       []string{"samsung"},
		AndroidMin:   7,
		AndroidMax:   12,
		PatchCeiling: "Needs a combination file whose bootloader version is >= the device's; newer devices lack public combination files.",
		Requires:     []string{"combination file for the exact model", "Odin", "Samsung USB driver", "full stock firmware for reflash"},
		Difficulty:   "hard",
		GuidePath:    "android-hacking/Bypassing_FRP_Lock/Combination-File-Method-(Samsung).md",
		Steps: []string{
			"Boot to recovery (volume up + power) and note the bootloader version from the build string.",
			"Find a combination file for the exact model with bootloader version >= the device's.",
			"Boot to download mode (volume down + power) and flash the combination file with Odin.",
			"In the combination firmware: Settings > Developer options > enable OEM unlocking.",
			"Back in download mode, flash the full original stock firmware.",
			"Complete setup without the old Google account.",
		},
	},
	{
		ID:           "oppo-a15-mtk-frp",
		Title:        "Oppo A15 FRP via MTK bypass + SP Flash Tool",
		Category:     "frp",
		Brands:       []string{"oppo"},
		AndroidMin:   10,
		AndroidMax:   12,
		PatchCeiling: "Depends on MTK auth-bypass support for the SoC in use.",
		Requires:     []string{"MediaTek USB drivers", "SP Flash Tool", "MTK bypass utility"},
		Difficulty:   "hard",
		GuidePath:    "android-hacking/Bypassing_FRP_Lock/Oppo_A15_with_MTK_Bypass/readme.md",
		Steps: []string{
			"Install the MediaTek USB drivers on the PC.",
			"Power off the phone; run the MTK auth-bypass utility and connect USB to disable SLA/DAA.",
			"Open SP Flash Tool and use the format path from the guide screenshots to wipe the FRP partition.",
			"Reboot; complete setup without the old Google account.",
		},
	},
	{
		ID:           "redmi-10c-secondspace-frp",
		Title:        "Redmi 10C FRP via Second Space",
		Category:     "frp",
		Brands:       []string{"xiaomi", "redmi"},
		AndroidMin:   11,
		AndroidMax:   13,
		PatchCeiling: "Second Space path documented on MIUI 13-era builds; verify against current patch.",
		Requires:     []string{"no PC needed"},
		Difficulty:   "medium",
		GuidePath:    "android-hacking/Bypassing_FRP_Lock/Redmi_10C_Second_Space/readme.md",
		Steps: []string{
			"Follow the guide's Second Space walkthrough to reach a working launcher from the setup wizard.",
			"From Second Space, open Settings and enable the options the guide lists.",
			"Switch back / reboot and complete setup without the old Mi/Google account.",
		},
	},
	{
		ID:           "vivo-edl-unlock",
		Title:        "Vivo lockscreen removal via EDL mode",
		Category:     "reset",
		Brands:       []string{"vivo"},
		AndroidMin:   10,
		AndroidMax:   13,
		PatchCeiling: "Needs the correct EDL programmer for the model.",
		Requires:     []string{"EDL 9008 access", "matching firehose programmer"},
		Difficulty:   "hard",
		GuidePath:    "android-hacking/Forgot_Password/EDL-Mode-Method-(Vivo).md",
		Steps: []string{
			"Boot the Vivo phone into EDL 9008 mode.",
			"Use an EDL tool with the model's programmer per the guide to wipe the lockscreen/config data.",
			"Reboot; set up the phone fresh.",
		},
	},
	{
		ID:           "realme-c21y-lockscreen",
		Title:        "Realme C21Y lockscreen removal (.pac flash)",
		Category:     "reset",
		Brands:       []string{"realme"},
		AndroidMin:   10,
		AndroidMax:   12,
		PatchCeiling: "Model-specific PAC file bundled in the guide.",
		Requires:     []string{"ResearchDownload/SPD tool", "Unisoc drivers", "the C21Y .pac from the guide"},
		Difficulty:   "hard",
		GuidePath:    "android-hacking/Forgot_Password/Remove-Lockscreen-Realme-C21Y.md",
		Steps: []string{
			"Install the Unisoc/SPD drivers from the bundled Drivers folder.",
			"Load the guide's C21Y .pac file in ResearchDownload.",
			"Boot the phone into download mode and flash per the guide screenshots.",
			"Reboot; the lockscreen is cleared.",
		},
	},
	{
		ID:           "oppo-a16k-clone-frp",
		Title:        "Oppo A16K FRP via Clone Phone (TalkBack)",
		Category:     "frp",
		Brands:       []string{"oppo"},
		AndroidMin:   11,
		AndroidMax:   13,
		PatchCeiling: "TalkBack gesture path; verify against current patch.",
		Requires:     []string{"second phone", "Wi-Fi connection", "no PC needed"},
		Difficulty:   "medium",
		GuidePath:    "android-hacking/Bypassing_FRP_Lock/Oppo_A16K_Clone_Phone/readme.md",
		Steps: []string{
			"Connect to Wi-Fi; enable TalkBack (hold volume up + volume down).",
			"Use the TalkBack voice-command gesture to reach system settings.",
			"Open the Clone Phone path from the guide to reach a working environment.",
			"From there enable the options the guide lists, then reboot and complete setup.",
		},
	},
	{
		ID:           "redmi-a5-unisoc-frp",
		Title:        "Redmi A5 FRP via Unisoc bootloader unlock (CVE-2022-38694)",
		Category:     "frp",
		Brands:       []string{"xiaomi", "redmi"},
		AndroidMin:   13,
		AndroidMax:   15,
		PatchCeiling: "Hardware-level Unisoc exploit; patch-independent but model-specific.",
		Requires:     []string{"Unisoc SPD drivers", "the ums9230e unlock tool from the guide", "USB cable"},
		Difficulty:   "hard",
		GuidePath:    "android-hacking/Bypassing_FRP_Lock/Redmi_A5/readme.md",
		Steps: []string{
			"Install the Unisoc drivers from the bundled Drivers folder.",
			"Power off the phone; run the guide's unlock tool and connect USB while holding volume down.",
			"Follow the tool prompts until the phone reboots with the bootloader unlocked.",
			"Boot to fastboot/recovery per the guide and wipe, then complete setup without the old account.",
		},
	},
	{
		ID:           "redmi-13c-simlock-frp",
		Title:        "Redmi 13C FRP via SIM lock + Second Space (HyperOS/MIUI 14)",
		Category:     "frp",
		Brands:       []string{"xiaomi", "redmi"},
		AndroidMin:   13,
		AndroidMax:   14,
		PatchCeiling: "Documented on HyperOS/MIUI 14; newer patches may close the notification-shade path.",
		Requires:     []string{"SIM card with SIM lock enabled", "second Android phone", "Wi-Fi"},
		Difficulty:   "medium",
		GuidePath:    "android-hacking/Bypassing_FRP_Lock/SIMLock-SecondSpace-Redmi13C.md",
		Steps: []string{
			"Connect to Wi-Fi and insert the SIM-locked SIM; when the PIN screen appears, lock the phone and remove the SIM.",
			"Turn the screen on, pull down the notification shade, and tap the Settings icon.",
			"Follow the guide's Google-help/Share path to reach a working settings environment.",
			"Enable Second Space per the guide, then reboot and complete setup without the old account.",
		},
	},
	{
		ID:           "xperia-mxtp-remove",
		Title:        "Sony Xperia theft-protection removal (Flashtool)",
		Category:     "reset",
		Brands:       []string{"sony"},
		AndroidMin:   6,
		AndroidMax:   9,
		PatchCeiling: "Legacy Xperia devices; needs the exact FTF firmware + lock-remove file.",
		Requires:     []string{"Flashtool + flashmode drivers", "matching FTF firmware", "matching lock-remove FTF"},
		Difficulty:   "hard",
		GuidePath:    "android-hacking/Forgot_Password/Bypassing-My-Xperia-Theft-Protection.md",
		Steps: []string{
			"Find the FTF firmware and lock-remove file for the exact model; place both under the Flashtool firmwares folder.",
			"Open Flashtool > Flash Device > Flashmode and flash the firmware first.",
			"Flash the lock-remove file next.",
			"Reboot; the theft-protection lock is cleared.",
		},
	},
	{
		ID:           "samsung-oneui3-recovery",
		Title:        "Samsung OneUI 3.x recovery-mode access",
		Category:     "info",
		Brands:       []string{"samsung"},
		AndroidMin:   11,
		AndroidMax:   13,
		PatchCeiling: "Key-combo reference; works regardless of patch.",
		Requires:     []string{"USB cable connected to a PC"},
		Difficulty:   "easy",
		GuidePath:    "android-hacking/Forgot_Password/Factory-Reset-OneUI3.x.md",
		Steps: []string{
			"Power off the device completely.",
			"Connect the USB cable to the device (other end in the PC).",
			"Hold volume up + power until recovery starts.",
			"From recovery you can factory-reset or apply updates per the guide.",
		},
	},
	{
		ID:           "find-model-locked",
		Title:        "Identify model number on a locked phone",
		Category:     "info",
		Brands:       []string{"generic"},
		AndroidMin:   0,
		AndroidMax:   0,
		PatchCeiling: "Dialer-code reference; availability varies by brand.",
		Requires:     []string{"access to the emergency dialer"},
		Difficulty:   "easy",
		GuidePath:    "android-hacking/Forgot_Password/Find-device-model-number.md",
		Steps: []string{
			"Open the emergency call dialer on the locked phone.",
			"Dial *#899# and open the SoftwareVersion / version info entry.",
			"Note the model number, then look up the matching procedure: dro1drain frp guide",
		},
	},
	{
		ID:           "buildprop-tweaks",
		Title:        "build.prop performance tweaks (root)",
		Category:     "info",
		Brands:       []string{"generic"},
		AndroidMin:   8,
		AndroidMax:   14,
		PatchCeiling: "Reference list; only useful with root access to /system.",
		Requires:     []string{"root access", "a root-capable file manager or ADB root"},
		Difficulty:   "medium",
		GuidePath:    "android-hacking/Tweaks/Build.prop-Tweaks.md",
		Steps: []string{
			"Remount /system read-write (requires root).",
			"Back up the existing build.prop first.",
			"Apply the tweak lines from the bundled guide that fit your device.",
			"Reboot and verify; restore the backup if anything misbehaves.",
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
		if m.GuidePath != "" {
			if strings.HasPrefix(m.GuidePath, "/") || strings.Contains(m.GuidePath, "..") {
				errs = append(errs, fmt.Sprintf("%s: bad GuidePath %q", m.ID, m.GuidePath))
			}
			if m.Automatable && m.AutoCmd == "" {
				errs = append(errs, fmt.Sprintf("%s: automatable but no AutoCmd", m.ID))
			}
		}
	}
	return errs
}
