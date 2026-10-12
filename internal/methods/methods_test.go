package methods

import "testing"

func TestValidate(t *testing.T) {
	if errs := Validate(); len(errs) > 0 {
		t.Fatalf("DB validation errors: %v", errs)
	}
}

func TestMatchSamsung(t *testing.T) {
	ms := Match("samsung", 13)
	if len(ms) == 0 {
		t.Fatal("expected methods for samsung/13")
	}
	// brand-specific should sort before generic
	if isGeneric(ms[0].Brands) {
		t.Errorf("first match should be brand-specific, got %s", ms[0].ID)
	}
	found := false
	for _, m := range ms {
		if m.ID == "samsung-testmode-frp" {
			found = true
		}
	}
	if !found {
		t.Error("samsung-testmode-frp not matched for samsung/13")
	}
}

func TestMatchXiaomi(t *testing.T) {
	ms := Match("Xiaomi", 12) // case-insensitive
	found := false
	for _, m := range ms {
		if m.ID == "xiaomi-sideload-account" {
			found = true
		}
	}
	if !found {
		t.Error("xiaomi-sideload-account not matched for Xiaomi/12")
	}
}

func TestMatchVersionBounds(t *testing.T) {
	// talkback method caps at 12; android 13 samsung should exclude it
	for _, m := range Match("samsung", 13) {
		if m.ID == "samsung-talkback-frp" {
			t.Error("talkback method should not match android 13")
		}
	}
	// android 10 samsung should include it
	found := false
	for _, m := range Match("samsung", 10) {
		if m.ID == "samsung-talkback-frp" {
			found = true
		}
	}
	if !found {
		t.Error("talkback method should match android 10")
	}
}

func TestByID(t *testing.T) {
	if ByID("nope") != nil {
		t.Error("ByID should return nil for unknown id")
	}
	if m := ByID("fastboot-factory-reset"); m == nil || m.Category != "reset" {
		t.Error("ByID failed for fastboot-factory-reset")
	}
}
