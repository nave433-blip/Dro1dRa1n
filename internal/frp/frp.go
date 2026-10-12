// Package frp implements FRP assessment and the host-side automation of
// documented FRP procedures.
//
// Honesty contract: nothing here magically unlocks a fully locked modern
// device. Each function automates the host (PC) side of a documented method
// and reports exactly what succeeded, what failed, and what the user must
// still do on the device.
package frp

import (
	"fmt"
	"strings"
	"time"

	"github.com/nave433-blip/Dro1dRa1n/internal/adb"
	"github.com/nave433-blip/Dro1dRa1n/internal/ui"
)

// StepResult records one automation step's outcome.
type StepResult struct {
	Name    string
	Command string
	OK      bool
	Output  string
}

// Status is the observable FRP-relevant state of a device.
type Status struct {
	Serial            string
	DeviceProvisioned string
	UserSetupComplete string
	ADBState          string // device / unauthorized / offline / none
	FRPPst            string
}

// Assess reads the FRP-relevant state that is observable over ADB.
// On a fully FRP-locked device ADB is usually unauthorized/absent — that
// itself is reported (it means the guide path, not the automated path).
func Assess(r *adb.Runner, serial string) (*Status, error) {
	st := &Status{Serial: serial}
	devs := adb.ParseDevices(must(r.Run("", "devices")))
	for _, d := range devs {
		if serial == "" || d.Serial == serial {
			st.ADBState = d.State
			if serial == "" {
				st.Serial = d.Serial
			}
			break
		}
	}
	if st.ADBState == "" {
		st.ADBState = "none"
	}
	if st.ADBState != "device" {
		return st, nil // nothing more observable
	}
	if out, err := r.Run(st.Serial, "shell", "settings", "get", "global", "device_provisioned"); err == nil {
		st.DeviceProvisioned = strings.TrimSpace(out)
	}
	if out, err := r.Run(st.Serial, "shell", "settings", "get", "secure", "user_setup_complete"); err == nil {
		st.UserSetupComplete = strings.TrimSpace(out)
	}
	if out, err := r.Run(st.Serial, "shell", "getprop", "ro.frp.pst"); err == nil {
		st.FRPPst = strings.TrimSpace(out)
	}
	return st, nil
}

func must(out string, err error) string {
	if err != nil {
		return ""
	}
	return out
}

// Report prints a human-readable assessment.
func (s *Status) Report() {
	ui.Header("FRP status")
	fmt.Printf("  serial:            %s\n", s.Serial)
	fmt.Printf("  adb state:         %s\n", s.ADBState)
	switch s.ADBState {
	case "device":
		fmt.Printf("  device_provisioned: %s\n", orUnknown(s.DeviceProvisioned))
		fmt.Printf("  user_setup_complete: %s\n", orUnknown(s.UserSetupComplete))
		fmt.Printf("  ro.frp.pst:        %s\n", orUnknown(s.FRPPst))
		if s.DeviceProvisioned == "0" || s.UserSetupComplete == "0" {
			ui.Warn("setup looks incomplete — FRP/setup wizard may still be active.")
		} else {
			ui.Ok("device looks provisioned (no FRP prompt expected).")
		}
	case "unauthorized":
		ui.Warn("ADB is unauthorized — tap 'Allow' on the device, then re-run.")
		ui.Info("If the device is FRP-locked with no Allow prompt, use: dro1drain frp guide")
	case "none", "offline":
		ui.Warn("no ADB device visible. Check USB/drivers, then run: dro1drain devices")
		ui.Info("driver hint: " + driverHint())
	}
}

func orUnknown(v string) string {
	if v == "" {
		return "(unknown)"
	}
	return v
}

// clearSequence is the documented ADB-side FRP-clear command sequence.
// Each entry: description + adb args.
var clearSequence = []struct {
	name string
	args []string
}{
	{"mark device provisioned", []string{"shell", "settings", "put", "global", "device_provisioned", "1"}},
	{"mark setup complete", []string{"shell", "settings", "put", "secure", "user_setup_complete", "1"}},
	{"dismiss setup wizard", []string{"shell", "settings", "put", "global", "setup_wizard_has_run", "1"}},
	{"open settings", []string{"shell", "am", "start", "-n", "com.android.settings/.Settings"}},
}

// ClearADB runs the ADB-side FRP-clear sequence against an authorized device.
func ClearADB(r *adb.Runner, serial string) []StepResult {
	var results []StepResult
	for _, s := range clearSequence {
		out, err := r.Run(serial, s.args...)
		res := StepResult{
			Name:    s.name,
			Command: "adb " + strings.Join(s.args, " "),
			OK:      err == nil,
			Output:  strings.TrimSpace(out),
		}
		if err != nil {
			res.Output = err.Error()
		}
		results = append(results, res)
		if res.OK {
			ui.Ok(s.name)
		} else {
			ui.Warn(s.name + " FAILED: " + res.Output)
		}
	}
	return results
}

// AllOK reports whether every step succeeded.
func AllOK(rs []StepResult) bool {
	for _, r := range rs {
		if !r.OK {
			return false
		}
	}
	return true
}

// SamsungFRP automates the host side of the Samsung test-mode method:
// waits for the device (user must have dialed *#0*# on it), verifies the
// ADB channel, then runs ClearADB.
func SamsungFRP(r *adb.Runner, waitTimeout time.Duration) error {
	ui.Header("samsung frp — test-mode method")
	ui.Info("On the phone: open Emergency call and dial *#0*# (test mode), then plug in USB.")
	ui.Info("Waiting for the device...")
	if err := r.WaitForDevice(waitTimeout); err != nil {
		return err
	}
	ui.Ok("device detected, ADB channel up.")

	// Identify which serial came up.
	out, _ := r.Run("", "devices")
	devs := adb.ParseDevices(out)
	var serial string
	for _, d := range devs {
		if d.State == "device" {
			serial = d.Serial
			break
		}
	}
	if serial == "" {
		return fmt.Errorf("no authorized device — tap Allow on the phone if prompted")
	}
	ui.Info("target: " + serial)
	results := ClearADB(r, serial)
	fmt.Println()
	if AllOK(results) {
		ui.Ok("ADB-side FRP sequence completed.")
		ui.Info("Reboot the phone and finish setup — the Google-account prompt should be gone.")
		ui.Info("If the prompt persists, your patch level needs a different path: dro1drain frp guide")
		return nil
	}
	return fmt.Errorf("some steps failed — see above; try: dro1drain frp guide")
}

// XiaomiMiCloud runs the documented ADB-side sequence for Xiaomi account
// state work (device in sideload/recovery ADB mode).
func XiaomiMiCloud(r *adb.Runner, serial string) error {
	ui.Header("xiaomi micloud — adb side")
	ui.Info("Put the phone in recovery ADB sideload mode and connect USB.")
	results := ClearADB(r, serial)
	fmt.Println()
	if AllOK(results) {
		ui.Ok("ADB-side sequence completed. Reboot and check the account prompt.")
		return nil
	}
	return fmt.Errorf("some steps failed — see above")
}

func driverHint() string {
	return "windows: Samsung/Google USB driver · linux: udev rules or sudo · macos: none needed"
}
