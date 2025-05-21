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

package set

import (
	"fmt"
)

func Example() {
	s := New(10, 10, 12, 15)
	fmt.Println(s.Len())                      // 3
	fmt.Println(s.Add(10))                    // false
	fmt.Println(s.Add(11))                    // true
	fmt.Println(s.Remove(11) && s.Remove(12)) // true

	fmt.Println(s.ContainsAny(10, 15)) // true
	fmt.Println(s.ContainsAny(11, 12)) // false
	fmt.Println(s.ContainsAny())       // false
	fmt.Println(s.ContainsAll(10, 15)) // true
	fmt.Println(s.ContainsAll(10, 11)) // false
	fmt.Println(s.ContainsAll())       // true

	fmt.Println(len(s.ToSlice())) // 2

	// Output:
	// 3
	// false
	// true
	// true
	// true
	// false
	// false
	// true
	// false
	// true
	// 2
}
