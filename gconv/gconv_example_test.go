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

package gconv

import (
	"fmt"

	"github.com/bytedance/gg/gptr"
)

func Example() {
	fmt.Println(To[string](1))                           // "1"
	fmt.Println(To[int]("1"))                            // 1
	fmt.Println(To[int]("x"))                            // 0
	fmt.Println(To[bool]("true"))                        // true
	fmt.Println(To[bool]("x"))                           // false
	fmt.Println(To[int](gptr.Of(gptr.Of(gptr.Of("1"))))) // 1
	type myInt int
	type myString string
	fmt.Println(To[myInt](myString("1"))) // 1
	fmt.Println(To[myString](myInt(1)))   // "1"

	fmt.Println(ToE[int]("x")) // 0 strconv.ParseInt: parsing "x": invalid syntax

	// Output:
	// 1
	// 1
	// 0
	// true
	// false
	// 1
	// 1
	// 1
	// 0 strconv.ParseInt: parsing "x": invalid syntax
}
