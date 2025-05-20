package skipset

import (
	"testing"
)

func TestFlag(t *testing.T) {
	// Correctness.
	const (
		f0 = 1 << iota
		f1
		f2
		f3
		f4
		f5
		f6
		f7
	)
	x := &bitflag{}

	x.SetTrue(f1 | f3)
	if x.Get(f0) || !x.Get(f1) || x.Get(f2) || !x.Get(f3) || !x.MGet(f0|f1|f2|f3, f1|f3) {
		t.Fatal("invalid")
	}
	x.SetTrue(f1)
	x.SetTrue(f1 | f3)
	if x.data != f1+f3 {
		t.Fatal("invalid")
	}

	x.SetFalse(f1 | f2)
	if x.Get(f0) || x.Get(f1) || x.Get(f2) || !x.Get(f3) || !x.MGet(f0|f1|f2|f3, f3) {
		t.Fatal("invalid")
	}
	x.SetFalse(f1 | f2)
	if x.data != f3 {
		t.Fatal("invalid")
	}
}
