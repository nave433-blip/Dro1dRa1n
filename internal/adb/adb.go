// Package adb wraps the Android platform-tools adb binary and parses its output.
package adb

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Runner executes adb invocations. Binary is resolved once via PATH.
type Runner struct {
	Bin     string
	Timeout time.Duration
}

// New returns a Runner, or an error if adb is not on PATH.
func New() (*Runner, error) {
	p, err := exec.LookPath("adb")
	if err != nil {
		return nil, fmt.Errorf("adb not found on PATH (install Android platform-tools)")
	}
	return &Runner{Bin: p, Timeout: 30 * time.Second}, nil
}

// Run executes: adb [-s serial] args...
func (r *Runner) Run(serial string, args ...string) (string, error) {
	full := []string{}
	if serial != "" {
		full = append(full, "-s", serial)
	}
	full = append(full, args...)
	cmd := exec.Command(r.Bin, full...)
	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	// simple timeout via goroutine-free approach: exec with deadline is not in
	// stdlib pre-1.19 style; use a timer kill through Start/Wait.
	if err := cmd.Start(); err != nil {
		return "", err
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		if err != nil {
			return out.String(), fmt.Errorf("adb %s: %v: %s", strings.Join(args, " "), err, strings.TrimSpace(errb.String()))
		}
		return out.String(), nil
	case <-time.After(r.Timeout):
		cmd.Process.Kill()
		return "", fmt.Errorf("adb %s: timed out after %s", strings.Join(args, " "), r.Timeout)
	}
}

// Device is one entry of `adb devices -l`.
type Device struct {
	Serial  string
	State   string // device, unauthorized, offline, recovery, sideload...
	Model   string
	Product string
}

// ParseDevices parses `adb devices -l` output.
func ParseDevices(output string) []Device {
	var devs []Device
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "List of devices") || strings.HasPrefix(line, "*") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		d := Device{Serial: fields[0], State: fields[1]}
		for _, f := range fields[2:] {
			if strings.HasPrefix(f, "model:") {
				d.Model = strings.TrimPrefix(f, "model:")
			} else if strings.HasPrefix(f, "product:") {
				d.Product = strings.TrimPrefix(f, "product:")
			}
		}
		devs = append(devs, d)
	}
	return devs
}

// ParseGetprop parses `adb shell getprop` output into a map.
func ParseGetprop(output string) map[string]string {
	props := map[string]string{}
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "[") {
			continue
		}
		// format: [key]: [value]
		parts := strings.SplitN(line, "]:", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimPrefix(strings.TrimSpace(parts[0]), "[")
		val := strings.TrimSpace(parts[1])
		val = strings.TrimPrefix(val, "[")
		val = strings.TrimSuffix(val, "]")
		props[key] = val
	}
	return props
}

// WaitForDevice blocks until a device is visible to adb (or timeout).
func (r *Runner) WaitForDevice(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		out, err := r.Run("", "devices")
		if err == nil {
			for _, d := range ParseDevices(out) {
				if d.State == "device" {
					return nil
				}
			}
		}
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("no authorized adb device appeared within %s", timeout)
}
