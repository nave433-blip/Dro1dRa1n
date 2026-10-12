package device

import (
	"strings"
	"testing"
)

func TestBrandKey(t *testing.T) {
	i := &Info{Brand: " Samsung "}
	if got := i.BrandKey(); got != "samsung" {
		t.Errorf("BrandKey = %q, want samsung", got)
	}
}

func TestSummary(t *testing.T) {
	i := &Info{Serial: "ABC", Brand: "samsung", Model: "SM-A525F", Android: "13", SDK: "33"}
	s := i.Summary()
	for _, want := range []string{"ABC", "samsung", "SM-A525F", "13"} {
		if !strings.Contains(s, want) {
			t.Errorf("Summary missing %q:\n%s", want, s)
		}
	}
	// empty fields render as (unknown), never blank-panic
	empty := (&Info{}).Summary()
	if !strings.Contains(empty, "(unknown)") {
		t.Error("Summary should mark unknown fields")
	}
}

func TestAndroidMajor(t *testing.T) {
	cases := map[string]int{"13": 13, "12.0": 12, "": 0, "abc": 0, " 14 ": 14}
	for in, want := range cases {
		if got := (&Info{Android: in}).AndroidMajor(); got != want {
			t.Errorf("Android %q: got %d, want %d", in, got, want)
		}
	}
}

func TestSummarySecurityRows(t *testing.T) {
	i := &Info{Serial: "X", Brand: "samsung", OEMUnlock: "1", FlashLocked: "0", KnoxWarranty: "1", AntiRollback: "2"}
	s := i.Summary()
	for _, want := range []string{"OEM unlocking allowed", "bootloader UNLOCKED", "Knox TRIPPED", "anti-rollback:  2"} {
		if !strings.Contains(s, want) {
			t.Errorf("summary missing %q\n%s", want, s)
		}
	}
	// empty security props -> no rows, no "(unknown)" noise for them
	s2 := (&Info{Serial: "Y"}).Summary()
	for _, notWant := range []string{"oem unlock", "knox warranty", "anti-rollback"} {
		if strings.Contains(s2, notWant) {
			t.Errorf("summary should omit empty %q", notWant)
		}
	}
}
