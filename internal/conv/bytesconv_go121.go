//go:build go1.21
// +build go1.21

package conv

import (
	"unsafe"
)

// StringToBytes converts string to []byte without a memory allocation.
func StringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// BytesToString converts []byte to string without a memory allocation.
func BytesToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
