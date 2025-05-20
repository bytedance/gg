//go:build go1.19 && !go1.22
// +build go1.19,!go1.22

package fastrand

import (
	_ "unsafe"
)

//go:linkname runtimefastrand runtime.fastrand
func runtimefastrand() uint32

//go:linkname runtimefastrand64 runtime.fastrand64
func runtimefastrand64() uint64

//go:linkname runtimefastrandu runtime.fastrandu
func runtimefastrandu() uint
