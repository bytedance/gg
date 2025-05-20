//go:build go1.22
// +build go1.22

package fastrand

import (
	"math/rand/v2"
)

func runtimefastrand() uint32 {
	return rand.Uint32()
}

func runtimefastrand64() uint64 {
	return rand.Uint64()
}

func runtimefastrandu() uint {
	return uint(rand.Uint64())
}
