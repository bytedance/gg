// Copyright 2025 Bytedance Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build !go1.19
// +build !go1.19

package fastrand

import (
	_ "unsafe"
)

//go:linkname runtimefastrand runtime.fastrand
func runtimefastrand() uint32

func runtimefastrand64() uint64 {
	return (uint64(runtimefastrand()) << 32) | uint64(runtimefastrand())
}

func runtimefastrandu() uint {
	// PtrSize is the size of a pointer in bytes - unsafe.Sizeof(uintptr(0)) but as an ideal constant.
	// It is also the size of the machine's native word size (that is, 4 on 32-bit systems, 8 on 64-bit).
	const PtrSize = 4 << (^uintptr(0) >> 63)
	if PtrSize == 4 {
		return uint(runtimefastrand())
	}
	return uint(runtimefastrand64())
}
