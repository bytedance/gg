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

package gcond

import (
	"fmt"
)

func Example() {
	fmt.Println(If(true, 1, 2)) // 1
	var a *struct{ A int }
	getA := func() int { return a.A }
	get1 := func() int { return 1 }
	fmt.Println(IfLazy(a != nil, getA, get1)) // 1
	fmt.Println(IfLazyL(a != nil, getA, 1))   // 1
	fmt.Println(IfLazyR(a == nil, 1, getA))   // 1

	fmt.Println(Switch[string](3).
		Case(1, "1").
		CaseLazy(2, func() string { return "3" }).
		When(3, 4).Then("3/4").
		When(5, 6).ThenLazy(func() string { return "5/6" }).
		Default("other")) // 3/4

	// Output:
	// 1
	// 1
	// 1
	// 1
	// 3/4
}
