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

package gresult

import (
	"fmt"
	"io"
	"strconv"
)

func Example() {
	fmt.Println(Of(strconv.Atoi("1")).Value())        // 1
	fmt.Println(Err[int](io.EOF).IsErr())             // true
	fmt.Println(Err[int](io.EOF).ValueOr(10))         // 10
	fmt.Println(OK(1).IsOK())                         // true
	fmt.Println(OK(1).ValueOrZero())                  // 1
	fmt.Println(Of(strconv.Atoi("x")).Option().Get()) // 0 false
	fmt.Println(Map(OK(1), strconv.Itoa).Get())       // "1" nil

	// Output:
	// 1
	// true
	// 10
	// true
	// 1
	// 0 false
	// 1 <nil>
}
