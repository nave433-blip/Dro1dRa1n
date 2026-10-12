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
