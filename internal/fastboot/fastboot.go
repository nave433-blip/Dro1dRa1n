// Package fastboot wraps the platform-tools fastboot binary.
package fastboot

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Has reports whether fastboot is on PATH.
func Has() bool {
	_, err := exec.LookPath("fastboot")
	return err == nil
}

// Run executes: fastboot [-s serial] args...
func Run(serial string, args ...string) (string, error) {
	full := []string{}
	if serial != "" {
		full = append(full, "-s", serial)
	}
	full = append(full, args...)
	cmd := exec.Command("fastboot", full...)
	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	done := make(chan error, 1)
	go func() { done <- cmd.Run() }()
	select {
	case err := <-done:
		combined := out.String() + errb.String() // fastboot logs to stderr
		if err != nil {
			return combined, fmt.Errorf("fastboot %s: %v", strings.Join(args, " "), err)
		}
		return combined, nil
	case <-time.After(60 * time.Second):
		return "", fmt.Errorf("fastboot %s: timed out", strings.Join(args, " "))
	}
}

// Device is one entry of `fastboot devices`.
type Device struct {
	Serial string
	State  string
}

// ParseDevices parses `fastboot devices` output.
func ParseDevices(output string) []Device {
	var devs []Device
	for _, line := range strings.Split(output, "\n") {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) >= 2 {
			devs = append(devs, Device{Serial: fields[0], State: fields[1]})
		}
	}
	return devs
}

// List returns attached fastboot devices (empty when none / binary missing).
func List() []Device {
	if !Has() {
		return nil
	}
	out, _ := Run("", "devices")
	return ParseDevices(out)
}
