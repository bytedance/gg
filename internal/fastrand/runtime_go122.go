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
