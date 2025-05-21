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

package gptr

import (
	"fmt"
	"strconv"
)

func Example() {
	a := Of(1)
	fmt.Println(Indirect(a)) // 1

	b := OfNotZero(1)
	fmt.Println(IsNotNil(b))                    // true
	fmt.Println(IndirectOr(b, 2))               // 1
	fmt.Println(Indirect(Map(b, strconv.Itoa))) // "1"

	c := OfNotZero(0)
	fmt.Println(c)                // nil
	fmt.Println(IsNil(c))         // true
	fmt.Println(IndirectOr(c, 2)) // 2

	// Output:
	// 1
	// true
	// 1
	// 1
	// <nil>
	// true
	// 2
}
