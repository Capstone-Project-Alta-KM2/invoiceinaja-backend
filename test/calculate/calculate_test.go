package calculate

import "testing"

func TestAdd(t *testing.T) {
	if Add(2, 3) != 5 {
		t.Error("expected 2 + 3 = 5")
	}
}

func TestSub(t *testing.T) {
	if Sub(2, 3) != -1 {
		t.Error("expected 2 - 3 = -1")
	}
}
