package adb

import (
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestParseDevices(t *testing.T) {
	out := `List of devices attached
R5CR31XXXX	device product:a52sxq model:SM_A525F device:a52sxq transport_id:1
emulator-5554	unauthorized
HT12345678	offline
`
	devs := ParseDevices(out)
	if len(devs) != 3 {
		t.Fatalf("expected 3 devices, got %d", len(devs))
	}
	if devs[0].Serial != "R5CR31XXXX" || devs[0].State != "device" {
		t.Errorf("bad device[0]: %+v", devs[0])
	}
	if devs[0].Model != "SM_A525F" || devs[0].Product != "a52sxq" {
		t.Errorf("bad model/product: %+v", devs[0])
	}
	if devs[1].State != "unauthorized" {
		t.Errorf("bad device[1]: %+v", devs[1])
	}
	if devs[2].State != "offline" {
		t.Errorf("bad device[2]: %+v", devs[2])
	}
}

func TestParseDevicesEmpty(t *testing.T) {
	devs := ParseDevices("List of devices attached\n\n")
	if len(devs) != 0 {
		t.Fatalf("expected 0 devices, got %d", len(devs))
	}
}

func TestParseGetprop(t *testing.T) {
	out := `[ro.product.brand]: [samsung]
[ro.product.model]: [SM-A525F]
[ro.build.version.release]: [13]
[ro.build.version.security_patch]: [2024-11-01]
[ro.build.fingerprint]: [samsung/a52sxqxx/a52sxq:13/TP1A.220624.014/A525FXXS8EXK1:user/release-keys]
not-a-prop-line
`
	props := ParseGetprop(out)
	want := map[string]string{
		"ro.product.brand":                "samsung",
		"ro.product.model":                "SM-A525F",
		"ro.build.version.release":        "13",
		"ro.build.version.security_patch": "2024-11-01",
	}
	for k, v := range want {
		if props[k] != v {
			t.Errorf("prop %q: got %q, want %q", k, props[k], v)
		}
	}
	if !reflect.DeepEqual(len(props), 5) {
		t.Errorf("expected 5 props, got %d", len(props))
	}
}

func TestRunTimeout(t *testing.T) {
	dir := t.TempDir()
	script := dir + "/slowadb"
	os.WriteFile(script, []byte("#!/bin/sh\nexec sleep 30\n"), 0755)
	r := &Runner{Bin: script, Timeout: 200 * time.Millisecond}
	_, err := r.Run("", "devices")
	if err == nil || !strings.Contains(err.Error(), "timed out") {
		t.Fatalf("expected timeout error, got %v", err)
	}
}

func TestWaitForSerialFinds(t *testing.T) {
	dir := t.TempDir()
	script := dir + "/fakeadb"
	os.WriteFile(script, []byte("#!/bin/sh\necho 'List of devices attached'\necho 'FAKE123\tdevice'\n"), 0755)
	r := &Runner{Bin: script, Timeout: 5 * time.Second}
	serial, err := r.WaitForSerial("", 3*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if serial != "FAKE123" {
		t.Fatalf("expected FAKE123, got %q", serial)
	}
}

func TestWaitForSerialTimeout(t *testing.T) {
	dir := t.TempDir()
	script := dir + "/emptyadb"
	os.WriteFile(script, []byte("#!/bin/sh\necho 'List of devices attached'\n"), 0755)
	r := &Runner{Bin: script, Timeout: 5 * time.Second}
	_, err := r.WaitForSerial("NOPE", 400*time.Millisecond)
	if err == nil || !strings.Contains(err.Error(), "NOPE") {
		t.Fatalf("expected serial-specific timeout error, got %v", err)
	}
}
