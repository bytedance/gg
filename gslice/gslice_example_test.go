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

package gslice

import (
	"fmt"
	"strconv"

	"github.com/bytedance/gg/gson"
	"github.com/bytedance/gg/gvalue"
)

func Example() {
	// High-order function
	fmt.Println(Map([]int{1, 2, 3, 4, 5}, strconv.Itoa)) // ["1", "2", "3", "4", "5"]
	isEven := func(i int) bool { return i%2 == 0 }
	fmt.Println(Filter([]int{1, 2, 3, 4, 5}, isEven))                  // [2, 4]
	fmt.Println(Reduce([]int{1, 2, 3, 4, 5}, gvalue.Add[int]).Value()) // 15
	fmt.Println(Any([]int{1, 2, 3, 4, 5}, isEven))                     // true
	fmt.Println(All([]int{1, 2, 3, 4, 5}, isEven))                     // false

	// CRUD operation
	fmt.Println(Contains([]int{1, 2, 3, 4, 5}, 2))                                    // true
	fmt.Println(ContainsBy([]int{1, 2, 3, 4, 5}, func(v int) bool { return v == 2 })) // true
	fmt.Println(ContainsAny([]int{1, 2, 3, 4, 5}, 2, 6))                              // true
	fmt.Println(ContainsAll([]int{1, 2, 3, 4, 5}, 2, 6))                              // false
	fmt.Println(Index([]int{1, 2, 3, 4, 5}, 3).Value())                               // 2
	fmt.Println(Find([]int{1, 2, 3, 4, 5}, isEven).Value())                           // 2
	fmt.Println(First([]int{1, 2, 3, 4, 5}).Value())                                  // 1
	fmt.Println(Get([]int{1, 2, 3, 4, 5}, 1).Value())                                 // 2
	fmt.Println(Get([]int{1, 2, 3, 4, 5}, -1).Value())                                //  5

	// Partion operation
	fmt.Println(Range(1, 5))                             // [1, 2, 3, 4]
	fmt.Println(RangeWithStep(5, 1, -2))                 // [5, 3]
	fmt.Println(Take([]int{1, 2, 3, 4, 5}, 2))           // [1, 2]
	fmt.Println(Take([]int{1, 2, 3, 4, 5}, -2))          // [4, 5]
	fmt.Println(Slice([]int{1, 2, 3, 4, 5}, 1, 3))       // [2, 3]
	fmt.Println(Chunk([]int{1, 2, 3, 4, 5}, 2))          // [[1, 2] [3, 4] [5]]
	fmt.Println(Divide([]int{1, 2, 3, 4, 5}, 2))         // [[1, 2, 3] [4, 5]]
	fmt.Println(Concat([]int{1, 2}, []int{3, 4, 5}))     // [1, 2, 3, 4, 5]
	fmt.Println(Flatten([][]int{{1, 2}, {3, 4, 5}}))     // [1, 2, 3, 4, 5]
	fmt.Println(Partition([]int{1, 2, 3, 4, 5}, isEven)) // [2, 4], [1, 3, 5]

	// Math operation
	fmt.Println(Max([]int{1, 2, 3, 4, 5}).Value())             // 5
	fmt.Println(Min([]int{1, 2, 3, 4, 5}).Value())             // 1
	fmt.Println(MinMax([]int{1, 2, 3, 4, 5}).Value().Values()) // 1 5
	fmt.Println(Sum([]int{1, 2, 3, 4, 5}))                     // 15

	// Convert to Map
	fmt.Println(gson.ToString(ToMap([]int{1, 2, 3, 4, 5}, func(i int) (string, int) { return strconv.Itoa(i), i }))) // {"1":1,"2":2,"3":3,"4":4,"5":5}
	fmt.Println(gson.ToString(ToMapValues([]int{1, 2, 3, 4, 5}, strconv.Itoa)))                                      // {"1":1,"2":2,"3":3,"4":4,"5":5}
	fmt.Println(gson.ToString(ToBoolMap([]int{1, 2, 3, 3, 2})))                                                          // {"1":true,"2":true,"3":true}
	fmt.Println(gson.ToString(GroupBy([]int{1, 2, 3, 4, 5}, func(i int) string {
		if i%2 == 0 {
			return "even"
		} else {
			return "odd"
		}
	}))) // {"even":[2,4], "odd":[1,3,5]}

	// Set operation
	fmt.Println(Union([]int{1, 2, 3}, []int{3, 4, 5}))     // [1, 2, 3, 4, 5]
	fmt.Println(Intersect([]int{1, 2, 3}, []int{3, 4, 5})) // [3]
	fmt.Println(Diff([]int{1, 2, 3}, []int{3, 4, 5}))      // [1, 2]
	fmt.Println(Uniq([]int{1, 1, 2, 2, 3}))                // [1, 2, 3]
	fmt.Println(Dup([]int{1, 1, 2, 2, 3}))                 // [1, 2]

	// Re-order operation
	s1 := []int{5, 1, 2, 3, 4}
	s2, s3, s4 := Clone(s1), Clone(s1), Clone(s1)
	Sort(s1)
	SortBy(s2, func(i, j int) bool { return i > j })
	StableSortBy(s3, func(i, j int) bool { return i > j })
	Reverse(s4)
	fmt.Println(s1) // [1, 2, 3, 4, 5]
	fmt.Println(s2) // [5, 4, 3, 2, 1]
	fmt.Println(s3) // [5, 4, 3, 2, 1]
	fmt.Println(s4) // [4, 3, 2, 1, 5]

	// Output:
	// [1 2 3 4 5]
	// [2 4]
	// 15
	// true
	// false
	// true
	// true
	// true
	// false
	// 2
	// 2
	// 1
	// 2
	// 5
	// [1 2 3 4]
	// [5 3]
	// [1 2]
	// [4 5]
	// [2 3]
	// [[1 2] [3 4] [5]]
	// [[1 2 3] [4 5]]
	// [1 2 3 4 5]
	// [1 2 3 4 5]
	// [2 4] [1 3 5]
	// 5
	// 1
	// 1 5
	// 15
	// {"1":1,"2":2,"3":3,"4":4,"5":5}
	// {"1":1,"2":2,"3":3,"4":4,"5":5}
	// {"1":true,"2":true,"3":true}
	// {"even":[2,4],"odd":[1,3,5]}
	// [1 2 3 4 5]
	// [3]
	// [1 2]
	// [1 2 3]
	// [1 2]
	// [1 2 3 4 5]
	// [5 4 3 2 1]
	// [5 4 3 2 1]
	// [4 3 2 1 5]
}
