// Dro1dRa1n — FRP all-in-one bypass toolkit for Android repair.
//
// Device repair, development, and educational use only. Use only on devices
// you own or have explicit permission to work on.
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nave433-blip/Dro1dRa1n/internal/adb"
	"github.com/nave433-blip/Dro1dRa1n/internal/device"
	"github.com/nave433-blip/Dro1dRa1n/internal/doctor"
	"github.com/nave433-blip/Dro1dRa1n/internal/fastboot"
	"github.com/nave433-blip/Dro1dRa1n/internal/frp"
	"github.com/nave433-blip/Dro1dRa1n/internal/methods"
	"github.com/nave433-blip/Dro1dRa1n/internal/ui"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	cmd, args := os.Args[1], os.Args[2:]
	var err error
	switch cmd {
	case "doctor":
		ok := doctor.Report()
		if !ok {
			os.Exit(2)
		}
	case "devices":
		err = cmdDevices()
	case "info":
		err = cmdInfo(args)
	case "frp":
		err = cmdFRP(args)
	case "samsung":
		err = cmdSamsung(args)
	case "xiaomi":
		err = cmdXiaomi(args)
	case "reboot":
		err = cmdReboot(args)
	case "reset":
		err = cmdReset(args)
	case "shell":
		err = cmdShell(args)
	case "methods":
		err = cmdMethods(args)
	case "version", "--version", "-v":
		fmt.Println("dro1drain", version)
	case "help", "--help", "-h":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n", cmd)
		usage()
		os.Exit(1)
	}
	if err != nil {
		ui.Fatal(err.Error())
	}
}

func usage() {
	fmt.Print(`dro1drain — FRP all-in-one bypass toolkit (repair use only)

Usage:
  dro1drain <command> [args]

Commands:
  doctor                  check dependencies (adb/fastboot/drivers)
  devices                 list ADB + fastboot devices
  info [serial]           full device info dump
  frp status [serial]     assess FRP lock state
  frp guide               interactive method picker for your device
  frp clear-adb [serial]  run ADB-side FRP clear (needs authorized ADB)
  samsung frp             automated Samsung test-mode (*#0*#) FRP flow
  samsung csc read        read Samsung CSC region code
  samsung debloat list    list disable-able Samsung packages
  samsung debloat run     disable Samsung bloatware (reversible)
  xiaomi micloud          Xiaomi ADB-side account sequence
  reboot <mode>           system|recovery|bootloader|download|edl
  reset [serial]          factory reset via fastboot (asks first)
  shell [serial] [-- cmd] adb shell passthrough
  methods list            browse the repair-methods database
  methods show <id>       show one method's full procedure
  version                 print version
`)
}

func needADB() *adb.Runner {
	r, err := adb.New()
	if err != nil {
		ui.Fatal(err.Error() + " — run: dro1drain doctor")
	}
	return r
}

func pickSerial(r *adb.Runner, args []string) string {
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		return args[0]
	}
	out, err := r.Run("", "devices")
	if err != nil {
		ui.Fatal(err.Error())
	}
	devs := adb.ParseDevices(out)
	var ok []adb.Device
	for _, d := range devs {
		if d.State == "device" {
			ok = append(ok, d)
		}
	}
	if len(ok) == 0 {
		ui.Fatal("no authorized adb devices — run: dro1drain devices")
	}
	if len(ok) == 1 {
		return ok[0].Serial
	}
	ui.Fatal("multiple devices; pass a serial explicitly")
	return ""
}

func cmdDevices() error {
	r, err := adb.New()
	hasADB := err == nil
	n := 0
	if hasADB {
		out, err := r.Run("", "devices")
		if err != nil {
			return err
		}
		ui.Header("adb devices")
		for _, d := range adb.ParseDevices(out) {
			fmt.Printf("  %-22s %-12s model=%s\n", d.Serial, d.State, d.Model)
			n++
		}
		if n == 0 {
			fmt.Println("  (none)")
		}
	} else {
		ui.Warn("adb not installed — run: dro1drain doctor")
	}
	ui.Header("fastboot devices")
	fdevs := fastboot.List()
	for _, d := range fdevs {
		fmt.Printf("  %-22s %s\n", d.Serial, d.State)
		n++
	}
	if len(fdevs) == 0 {
		fmt.Println("  (none)")
	}
	if n == 0 {
		ui.Info("driver hint: " + "windows: Samsung/Google USB driver · linux: udev rules or sudo · macos: none needed")
	}
	return nil
}

func cmdInfo(args []string) error {
	r := needADB()
	serial := pickSerial(r, args)
	info, err := device.Collect(r, serial)
	if err != nil {
		return err
	}
	ui.Header("device info")
	fmt.Print(info.Summary())
	return nil
}

func cmdFRP(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: dro1drain frp <status|guide|clear-adb> [serial]")
	}
	switch args[0] {
	case "status":
		r := needADB()
		serial := ""
		if len(args) > 1 {
			serial = args[1]
		}
		st, err := frp.Assess(r, serial)
		if err != nil {
			return err
		}
		st.Report()
		return nil
	case "guide":
		return cmdGuide()
	case "clear-adb":
		r := needADB()
		serial := pickSerial(r, args[1:])
		ui.Header("frp clear-adb")
		results := frp.ClearADB(r, serial)
		if frp.AllOK(results) {
			ui.Ok("done — reboot and verify.")
			return nil
		}
		return fmt.Errorf("some steps failed")
	default:
		return fmt.Errorf("unknown frp subcommand %q", args[0])
	}
}

// detectDevice returns brand + android major from the single connected
// authorized device, or ("", 0, false) when detection isn't possible.
func detectDevice() (string, int, bool) {
	r, err := adb.New()
	if err != nil {
		return "", 0, false
	}
	out, err := r.Run("", "devices")
	if err != nil {
		return "", 0, false
	}
	var serial string
	n := 0
	for _, d := range adb.ParseDevices(out) {
		if d.State == "device" {
			n++
			serial = d.Serial
		}
	}
	if n != 1 {
		return "", 0, false
	}
	info, err := device.Collect(r, serial)
	if err != nil {
		return "", 0, false
	}
	if info.BrandKey() == "" && info.AndroidMajor() == 0 {
		return "", 0, false
	}
	return info.BrandKey(), info.AndroidMajor(), true
}

func itoaOrEmpty(v int) string {
	if v == 0 {
		return ""
	}
	return strconv.Itoa(v)
}

func cmdGuide() error {
	ui.Header("frp guide — method picker")
	brand, major, auto := detectDevice()
	var verStr string
	if auto {
		verStr = itoaOrEmpty(major)
		fmt.Printf("auto-detected device: brand=%s android=%s\n\n", orDash(brand), orDash(verStr))
	} else {
		brand = strings.ToLower(ui.Ask("device brand (e.g. samsung, xiaomi, motorola)"))
		verStr = ui.Ask("android major version (e.g. 13; empty if unknown)")
		if v, err := strconv.Atoi(strings.TrimSpace(verStr)); err == nil {
			major = v
		}
	}
	ms := methods.Match(brand, major)
	if len(ms) == 0 {
		ui.Warn("no methods matched — try the generic list: dro1drain methods list")
		return nil
	}
	fmt.Printf("\n%d method(s) match %s / android %s:\n", len(ms), brand, verStr)
	for i, m := range ms {
		auto := ""
		if m.Automatable {
			auto = "  [automatable: dro1drain " + m.AutoCmd + "]"
		}
		fmt.Printf("\n  %d. %s%s\n     difficulty: %s · patch notes: %s\n", i+1, m.Title, auto, m.Difficulty, m.PatchCeiling)
	}
	choice := ui.Ask("\nshow full steps for number (empty to quit)")
	if choice == "" {
		return nil
	}
	idx, err := strconv.Atoi(strings.TrimSpace(choice))
	if err != nil || idx < 1 || idx > len(ms) {
		return fmt.Errorf("bad choice")
	}
	printMethod(&ms[idx-1])
	return nil
}

func printMethod(m *methods.Method) {
	ui.Header(m.Title)
	fmt.Printf("id: %s · category: %s · difficulty: %s\n", m.ID, m.Category, m.Difficulty)
	fmt.Printf("patch notes: %s\n", m.PatchCeiling)
	if len(m.Requires) > 0 {
		fmt.Println("requires:")
		for _, r := range m.Requires {
			fmt.Println("  -", r)
		}
	}
	if m.GuidePath != "" {
		fmt.Printf("bundled guide: %s\n", m.GuidePath)
	}
	for i, s := range m.Steps {
		ui.Step(i+1, s)
	}
	if m.Automatable {
		fmt.Printf("\nThis method's host side can be automated: dro1drain %s\n", m.AutoCmd)
	}
}

func cmdSamsung(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: dro1drain samsung <frp|csc|debloat> ...")
	}
	switch args[0] {
	case "frp":
		r := needADB()
		return frp.SamsungFRP(r, 120*time.Second)
	case "csc":
		return cmdSamsungCSC(args[1:])
	case "debloat":
		return cmdSamsungDebloat(args[1:])
	default:
		return fmt.Errorf("unknown samsung subcommand %q", args[0])
	}
}

func cmdSamsungCSC(args []string) error {
	r := needADB()
	serial := pickSerial(r, args)
	if len(args) == 0 || args[0] == "read" {
		out, err := r.Run(serial, "shell", "getprop", "ro.csc.sales_code")
		if err != nil {
			return err
		}
		code := strings.TrimSpace(out)
		if code == "" {
			out2, err2 := r.Run(serial, "shell", "getprop", "ro.csc.country_code")
			if err2 == nil {
				code = strings.TrimSpace(out2)
			}
		}
		ui.Header("samsung csc")
		fmt.Printf("  current CSC: %s\n", orDash(code))
		ui.Info("to change: dro1drain samsung csc set <CODE>  (guided, confirms first)")
		return nil
	}
	if args[0] == "set" {
		if len(args) < 2 {
			return fmt.Errorf("usage: dro1drain samsung csc set <CODE>")
		}
		code := strings.ToUpper(strings.TrimSpace(args[1]))
		ui.Warn("CSC change rewrites region config — the phone will factory-reset.")
		if !ui.Confirm("change CSC to " + code + "?") {
			ui.Info("cancelled.")
			return nil
		}
		// Documented MTP-path CSC change is device-UI driven; the host side
		// stages the change request where the build permits it.
		out, err := r.Run(serial, "shell", "am", "broadcast", "-a",
			"com.samsung.android.action.CHANGE_CSC", "--es", "csc", code)
		fmt.Printf("  broadcast result: %s\n", strings.TrimSpace(out))
		if err != nil {
			ui.Warn("broadcast not accepted on this build — use the MTP method: dro1drain methods show samsung-csc-change")
			return nil
		}
		ui.Ok("CSC change staged — reboot and confirm in Settings > About.")
		return nil
	}
	return fmt.Errorf("usage: dro1drain samsung csc <read|set CODE>")
}

func cmdSamsungDebloat(args []string) error {
	r := needADB()
	serial := pickSerial(r, args)
	// Curated, safe-to-disable Samsung packages (reversible via `cmd package install-existing`).
	pkgs := []string{
		"com.samsung.android.bixby.agent",
		"com.samsung.android.bixby.wakeup",
		"com.samsung.android.app.spage",
		"com.samsung.android.mateagent",
		"com.samsung.android.da.daagent",
		"com.samsung.android.arzone",
		"com.samsung.android.ipsgeofence",
		"com.samsung.android.game.gamehome",
		"com.samsung.android.game.gametools",
		"com.samsung.android.kidsinstaller",
		"com.samsung.android.app.simplesharing",
		"com.samsung.android.themestore",
		"com.samsung.android.samsungpass",
		"com.samsung.android.spay",
	}
	sub := "list"
	if len(args) > 0 {
		sub = args[0]
	}
	ui.Header("samsung debloat")
	switch sub {
	case "list":
		for _, p := range pkgs {
			fmt.Println("  " + p)
		}
		ui.Info("run: dro1drain samsung debloat run   (reversible)")
		return nil
	case "run":
		if !ui.Confirm(fmt.Sprintf("disable %d packages for user 0?", len(pkgs))) {
			ui.Info("cancelled.")
			return nil
		}
		done, failed := 0, 0
		for _, p := range pkgs {
			_, err := r.Run(serial, "shell", "pm", "disable-user", "--user", "0", p)
			if err != nil {
				ui.Warn("failed: " + p)
				failed++
			} else {
				done++
			}
		}
		ui.Ok(fmt.Sprintf("disabled %d, failed %d (restore: pm install-existing --user 0 <pkg>)", done, failed))
		return nil
	default:
		return fmt.Errorf("usage: dro1drain samsung debloat <list|run>")
	}
}

func cmdXiaomi(args []string) error {
	if len(args) > 0 && args[0] != "micloud" {
		return fmt.Errorf("unknown xiaomi subcommand %q", args[0])
	}
	rest := args
	if len(rest) > 0 && rest[0] == "micloud" {
		rest = rest[1:]
	}
	r := needADB()
	serial := pickSerial(r, rest)
	return frp.XiaomiMiCloud(r, serial)
}

func cmdReboot(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: dro1drain reboot <system|recovery|bootloader|download|edl>")
	}
	mode := strings.ToLower(args[0])
	serial := ""
	rest := args[1:]
	if len(rest) > 0 && !strings.HasPrefix(rest[0], "-") {
		serial = rest[0]
	}
	r := needADB()
	var adbArgs []string
	switch mode {
	case "system", "normal":
		adbArgs = []string{"reboot"}
	case "recovery":
		adbArgs = []string{"reboot", "recovery"}
	case "bootloader", "fastboot":
		adbArgs = []string{"reboot", "bootloader"}
	case "download", "odin":
		adbArgs = []string{"reboot", "download"}
	case "edl", "9008":
		adbArgs = []string{"reboot", "edl"}
	default:
		return fmt.Errorf("unknown mode %q", mode)
	}
	ui.Info("rebooting to " + mode + " ...")
	out, err := r.Run(serial, adbArgs...)
	fmt.Printf("  %s\n", strings.TrimSpace(out))
	return err
}

func cmdReset(args []string) error {
	serial := ""
	if len(args) > 0 {
		serial = args[0]
	}
	ui.Header("factory reset")
	ui.Warn("This wipes userdata via fastboot. It does NOT remove FRP by itself.")
	ui.Warn("Use only on devices you own or have permission to wipe.")
	if !ui.Confirm("wipe userdata on " + orDash(serial) + "?") {
		ui.Info("cancelled.")
		return nil
	}
	out, err := fastboot.Run(serial, "-w")
	fmt.Printf("  %s\n", strings.TrimSpace(out))
	if err != nil {
		return fmt.Errorf("fastboot wipe failed: %v (is the device in fastboot mode?)", err)
	}
	ui.Ok("wiped. Rebooting...")
	_, _ = fastboot.Run(serial, "reboot")
	return nil
}

func cmdShell(args []string) error {
	r := needADB()
	serial := ""
	rest := args
	// allow: dro1drain shell [serial] [--] [command...]
	var shArgs []string
	for i, a := range rest {
		if a == "--" {
			shArgs = rest[i+1:]
			rest = rest[:i]
			break
		}
	}
	if len(rest) > 0 && !strings.HasPrefix(rest[0], "-") {
		// first token could be serial or start of command; treat as serial only
		// if a device with that serial exists
		out, _ := r.Run("", "devices")
		for _, d := range adb.ParseDevices(out) {
			if d.Serial == rest[0] && d.State == "device" {
				serial = rest[0]
				rest = rest[1:]
				break
			}
		}
		shArgs = append(shArgs, rest...)
	} else {
		shArgs = rest
	}
	full := append([]string{"shell"}, shArgs...)
	out, err := r.Run(serial, full...)
	fmt.Print(out)
	return err
}

func cmdMethods(args []string) error {
	if len(args) == 0 || args[0] == "list" {
		cat := ""
		fs := flag.NewFlagSet("methods list", flag.ContinueOnError)
		fs.StringVar(&cat, "category", "", "filter: frp|reset|flash|account|info")
		_ = fs.Parse(args[1:])
		ui.Header("methods database")
		n := 0
		for _, m := range methods.DB {
			if cat != "" && m.Category != cat {
				continue
			}
			auto := ""
			if m.Automatable {
				auto = " [auto]"
			}
			fmt.Printf("  %-24s %-8s %s%s\n", m.ID, m.Category, m.Title, auto)
			n++
		}
		fmt.Printf("\n%d method(s). Show one: dro1drain methods show <id>\n", n)
		return nil
	}
	if args[0] == "show" {
		if len(args) < 2 {
			return fmt.Errorf("usage: dro1drain methods show <id>")
		}
		m := methods.ByID(args[1])
		if m == nil {
			return fmt.Errorf("unknown method %q", args[1])
		}
		printMethod(m)
		return nil
	}
	return fmt.Errorf("usage: dro1drain methods <list|show ID>")
}

func orDash(s string) string {
	if strings.TrimSpace(s) == "" {
		return "(any)"
	}
	return s
}
