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
