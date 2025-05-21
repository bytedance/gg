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

package skipset

import (
	"fmt"
	"sync"
)

func Example() {
	s := New[int]()
	fmt.Println(s.Add(10)) // true
	fmt.Println(s.Add(10)) // false
	fmt.Println(s.Add(11)) // true
	fmt.Println(s.Add(12)) // true
	fmt.Println(s.Len())   // 3

	fmt.Println(s.Contains(10)) // true
	fmt.Println(s.Remove(10))   // true
	fmt.Println(s.Contains(10)) // false

	fmt.Println(s.ToSlice()) // [11, 12]

	var wg sync.WaitGroup
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		i := i
		go func() {
			defer wg.Done()
			s.Add(i)
		}()
	}
	wg.Wait()
	fmt.Println(s.Len()) // 1000

	// Output:
	// true
	// false
	// true
	// true
	// 3
	// true
	// true
	// false
	// [11 12]
	// 1000
}
