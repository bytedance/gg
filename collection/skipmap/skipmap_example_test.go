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

package skipmap

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/bytedance/gg/stdwrap/gson"
)

func Example() {
	s := New[string, int]()
	s.Store("a", 0)
	s.Store("a", 1)
	s.Store("b", 2)
	s.Store("c", 3)
	fmt.Println(s.Len()) // 3

	fmt.Println(s.Load("a"))            // 1 true
	fmt.Println(s.LoadAndDelete("a"))   // 1 true
	fmt.Println(s.LoadOrStore("a", 11)) // 11 false

	fmt.Println(gson.ToString(s.ToMap())) // {"a":11, "b":2, "c": 3}

	s.Delete("a")
	s.Delete("b")
	s.Delete("c")
	var wg sync.WaitGroup
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		i := i
		go func() {
			defer wg.Done()
			s.Store(strconv.Itoa(i), i)
		}()
	}
	wg.Wait()
	fmt.Println(s.Len()) // 1000

	// Output:
	// 3
	// 1 true
	// 1 true
	// 11 false
	// {"a":11,"b":2,"c":3}
	// 1000
}
