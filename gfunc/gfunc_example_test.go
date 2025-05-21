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

package gfunc

import (
	"fmt"

	"github.com/bytedance/gg/gvalue"
)

func Example() {
	add := Partial2(gvalue.Add[int]) // Cast f to "partial application"-able function
	add1 := add.Partial(1)           // Bind argument a to 1
	fmt.Println(add1(0))             // 0 + 1 = 1
	fmt.Println(add1(1))             // add1 can be reused, 1 + 1 = 2
	add1n2 := add1.PartialR(2)       // Bind argument b to 2, all arguments are fixed
	fmt.Println(add1n2())            // 1 + 2 = 3

	// Output:
	// 1
	// 2
	// 3
}
