package plus

import "testing"

func TestPlus(t *testing.T) {
	res := Plus(1, 2)
	if res != 3 {
		t.Errorf("mismatch answer %d", res)
	}
}
