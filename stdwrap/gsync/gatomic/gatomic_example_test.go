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

package gatomic

import (
	"fmt"
)

func Example() {
	var v Value[int]
	fmt.Println(v.Load()) // 0
	v.Store(1)
	fmt.Println(v.Load())               // 1
	fmt.Println(v.Swap(2))              // 1
	fmt.Println(v.Load())               // 2
	fmt.Println(v.CompareAndSwap(1, 3)) // false
	fmt.Println(v.Load())               // 2
	fmt.Println(v.CompareAndSwap(2, 3)) // true
	fmt.Println(v.Load())               // 3

	// Output:
	// 0
	// 1
	// 1
	// 2
	// false
	// 2
	// true
	// 3
}
