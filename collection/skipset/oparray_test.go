package skipset

import (
	"testing"
	"unsafe"

	"github.com/bytedance/gg/internal/fastrand"
)

type dummy struct {
	data optionalArray
}

func TestOpArray(t *testing.T) {
	n := new(dummy)
	n.data.extra = new([op2]unsafe.Pointer)

	var array [maxLevel]unsafe.Pointer
	for i := 0; i < maxLevel; i++ {
		value := unsafe.Pointer(&dummy{})
		array[i] = value
		n.data.store(i, value)
	}

	for i := 0; i < maxLevel; i++ {
		if array[i] != n.data.load(i) || array[i] != n.data.atomicLoad(i) {
			t.Fatal(i, array[i], n.data.load(i))
		}
	}

	for i := 0; i < 1000; i++ {
		r := int(fastrand.Uint32n(maxLevel))
		value := unsafe.Pointer(&dummy{})
		if i%100 == 0 {
			value = nil
		}
		array[r] = value
		if fastrand.Uint32n(2) == 0 {
			n.data.store(r, value)
		} else {
			n.data.atomicStore(r, value)
		}
	}

	for i := 0; i < maxLevel; i++ {
		if array[i] != n.data.load(i) || array[i] != n.data.atomicLoad(i) {
			t.Fatal(i, array[i], n.data.load(i))
		}
	}
}
