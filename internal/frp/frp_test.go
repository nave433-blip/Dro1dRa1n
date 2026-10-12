package frp

import "testing"

func TestAllOK(t *testing.T) {
	if !AllOK([]StepResult{{Name: "a", OK: true}, {Name: "b", OK: true}}) {
		t.Error("AllOK should be true when all steps ok")
	}
	if AllOK([]StepResult{{Name: "a", OK: true}, {Name: "b", OK: false}}) {
		t.Error("AllOK should be false when any step fails")
	}
	if !AllOK(nil) {
		t.Error("AllOK should be true for empty results")
	}
}

func TestClearSequenceWellFormed(t *testing.T) {
	if len(clearSequence) == 0 {
		t.Fatal("clearSequence must not be empty")
	}
	for _, s := range clearSequence {
		if s.name == "" {
			t.Error("sequence step with empty name")
		}
		if len(s.args) < 2 {
			t.Errorf("step %q has too few args", s.name)
		}
	}
}
