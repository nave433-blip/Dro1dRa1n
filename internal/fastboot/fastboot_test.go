package fastboot

import "testing"

func TestParseDevices(t *testing.T) {
	out := "ABC123\tfastboot\nDEF456\tfastboot\n"
	devs := ParseDevices(out)
	if len(devs) != 2 {
		t.Fatalf("expected 2, got %d", len(devs))
	}
	if devs[0].Serial != "ABC123" || devs[0].State != "fastboot" {
		t.Errorf("bad device: %+v", devs[0])
	}
	if got := len(ParseDevices("\n")); got != 0 {
		t.Errorf("expected 0 for empty input, got %d", got)
	}
}
