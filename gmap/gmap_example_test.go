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

package gmap

import (
	"fmt"
	"strconv"

	"github.com/bytedance/gg/gslice"
	"github.com/bytedance/gg/stdwrap/gson"
)

func Example() {
	// Keys / Values getter
	fmt.Println(Keys(map[int]int{1: 2}))                             // [1]
	fmt.Println(Values(map[int]int{1: 2}))                           // [2]
	fmt.Println(Items(map[int]int{1: 2}).Unzip())                    // [1] [2]
	fmt.Println(OrderedKeys(map[int]int{1: 2, 2: 3, 3: 4}))          // [1, 2, 3]
	fmt.Println(OrderedValues(map[int]int{1: 2, 2: 3, 3: 4}))        // [2, 3, 4]
	fmt.Println(OrderedItems(map[int]int{1: 2, 2: 3, 3: 4}).Unzip()) // [1, 2, 3] [2, 3, 4]
	f := func(k, v int) string { return strconv.Itoa(k) + ":" + strconv.Itoa(v) }
	fmt.Println(ToSlice(map[int]int{1: 2}, f))                    // ["1:2"]
	fmt.Println(ToOrderedSlice(map[int]int{1: 2, 2: 3, 3: 4}, f)) // ["1:2", "2:3", "3:4"]

	// High-order function
	fmt.Println(gson.ToString(Map(map[int]int{1: 2, 2: 3, 3: 4}, func(k int, v int) (string, string) {
		return strconv.Itoa(k), strconv.Itoa(k + 1)
	}))) // {"1":"2", "2":"3", "3":"4"}
	fmt.Println(gson.ToString(Filter(map[int]int{1: 2, 2: 3, 3: 4}, func(k int, v int) bool {
		return k+v > 3
	}))) // {"2":2, "3":3}

	// CRUD operation
	fmt.Println(Contains(map[int]int{1: 2, 2: 3, 3: 4}, 1))           // true
	fmt.Println(ContainsAny(map[int]int{1: 2, 2: 3, 3: 4}, 1, 4))     // true
	fmt.Println(ContainsAll(map[int]int{1: 2, 2: 3, 3: 4}, 1, 4))     // false
	fmt.Println(Load(map[int]int{1: 2, 2: 3, 3: 4}, 1).Value())       // 2
	fmt.Println(LoadAny(map[int]int{1: 2, 2: 3, 3: 4}, 1, 4).Value()) // 2
	fmt.Println(LoadAll(map[int]int{1: 2, 2: 3, 3: 4}, 1, 4))         // []
	fmt.Println(LoadSome(map[int]int{1: 2, 2: 3, 3: 4}, 1, 4))        // [2]

	// Partion operation
	Chunk(map[int]int{1: 2, 2: 3, 3: 4, 4: 5, 5: 6}, 2)  // possible output: [{1:2, 2:3}, {3:4, 4:5}, {5:6}]
	Divide(map[int]int{1: 2, 2: 3, 3: 4, 4: 5, 5: 6}, 2) // possible output: [{1:2, 2:3, 3:4}, {4:5, 5:6}]

	// Math operation
	fmt.Println(Max(map[int]int{1: 2, 2: 3, 3: 4}).Value())             // 4
	fmt.Println(Min(map[int]int{1: 2, 2: 3, 3: 4}).Value())             // 2
	fmt.Println(MinMax(map[int]int{1: 2, 2: 3, 3: 4}).Value().Values()) // 2 4
	fmt.Println(Sum(map[int]int{1: 2, 2: 3, 3: 4}))                     // 9

	// Set operation
	fmt.Println(gson.ToString(Union(map[int]int{1: 2, 2: 3, 3: 4}, map[int]int{3: 14, 4: 15, 5: 16})))                                          // {1:2, 2:3, 3:14, 4:15, 5:16}
	fmt.Println(gson.ToString(Intersect(map[int]int{1: 2, 2: 3, 3: 4}, map[int]int{3: 14, 4: 15, 5: 16})))                                      // {3:14}
	fmt.Println(gson.ToString(Diff(map[int]int{1: 2, 2: 3, 3: 4}, map[int]int{3: 14, 4: 15, 5: 16})))                                           // {1:2, 2:3}
	fmt.Println(gson.ToString(UnionBy(gslice.Of(map[int]int{1: 2, 2: 3, 3: 4}, map[int]int{3: 14, 4: 15, 5: 16}), DiscardNew[int, int]())))     // {1:2, 2:3, 3:4, 4:15, 5:16}
	fmt.Println(gson.ToString(IntersectBy(gslice.Of(map[int]int{1: 2, 2: 3, 3: 4}, map[int]int{3: 14, 4: 15, 5: 16}), DiscardNew[int, int]()))) // {3:4}

	// Output:
	// [1]
	// [2]
	// [1] [2]
	// [1 2 3]
	// [2 3 4]
	// [1 2 3] [2 3 4]
	// [1:2]
	// [1:2 2:3 3:4]
	// {"1":"2","2":"3","3":"4"}
	// {"2":3,"3":4}
	// true
	// true
	// false
	// 2
	// 2
	// []
	// [2]
	// 4
	// 2
	// 2 4
	// 9
	// {"1":2,"2":3,"3":14,"4":15,"5":16}
	// {"3":14}
	// {"1":2,"2":3}
	// {"1":2,"2":3,"3":4,"4":15,"5":16}
	// {"3":4}
}
