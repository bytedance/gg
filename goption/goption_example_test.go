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

package goption

import (
	"fmt"
	"strconv"
)

func Example() {
	fmt.Println(Of(1, true).Value())            // 1
	fmt.Println(Nil[int]().IsNil())             // true
	fmt.Println(Nil[int]().ValueOr(10))         // 10
	fmt.Println(OK(1).IsOK())                   // true
	fmt.Println(OK(1).ValueOrZero())            // 1
	fmt.Println(OfPtr((*int)(nil)).Ptr())       // nil
	fmt.Println(Map(OK(1), strconv.Itoa).Get()) // "1" true

	// Output:
	// 1
	// true
	// 10
	// true
	// 1
	// <nil>
	// 1 true
}
