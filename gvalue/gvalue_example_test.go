// Copyright 2025 Bytedance Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gvalue

import (
	"fmt"
)

var once = Once(func() int {
	fmt.Println("once")
	return 0
})

func Example() {
	// Zero value
	a := Zero[int]()
	fmt.Println(a)         // 0
	fmt.Println(IsZero(a)) // true
	b := Zero[*int]()
	fmt.Println(b)        // nil
	fmt.Println(IsNil(b)) // true

	// Math operation
	fmt.Println(Max(1, 2, 3))    // 3
	fmt.Println(Min(1, 2, 3))    // 1
	fmt.Println(MinMax(1, 2, 3)) // 1 3
	fmt.Println(Clamp(5, 1, 10)) // 5
	fmt.Println(Add(1, 2))       // 3

	// Comparison
	fmt.Println(Equal(1, 1))      // true
	fmt.Println(Between(2, 1, 3)) // true

	// Type assertion
	fmt.Println(TypeAssert[int](any(1))) // 1
	fmt.Println(TryAssert[int](any(1)))  // 1 true

	// Once
	once() // "once"
	once() // (no output)
	once() // (no output)

	// Output:
	// 0
	// true
	// <nil>
	// true
	// 3
	// 1
	// 1 3
	// 5
	// 3
	// true
	// true
	// 1
	// 1 true
	// once
}
