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

	"github.com/bytedance/gg/iter"
)

func Example() {
	s := New(10, 10, 12, 15)
	fmt.Println("here are", s.Len(), "members")

	if !s.Add(10) {
		fmt.Println("member 10 already exists")
	}

	if s.Add(11) {
		fmt.Println("11 is added as member")
	}

	if s.Remove(11) && s.Remove(12) {
		fmt.Println("member 11 and 12 are removed")
	}

	var members []int
	s.Range(func(v int) bool {
		members = append(members, v)
		return true
	})
	fmt.Println("here are", len(members), "members")
	// Output:
	// here are 3 members
	// member 10 already exists
	// 11 is added as member
	// member 11 and 12 are removed
	// here are 2 members
}

func ExampleIter() {
	s := New(5, 3, 2, 1, 4)

	for _, v := range iter.ToSlice(iter.Sort(s.Iter())) {
		fmt.Println(v)
	}

	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
}

func ExampleSet_ContainsAny() {
	s := New(1, 2, 3, 4)

	fmt.Println(s.ContainsAny(1, 5))
	fmt.Println(s.ContainsAny(5, 6))
	fmt.Println(s.ContainsAny())

	// Output:
	// true
	// false
	// false
}

func ExampleSet_ContainsAll() {
	s := New(1, 2, 3, 4)

	fmt.Println(s.ContainsAll(1, 5))
	fmt.Println(s.ContainsAll(1, 2))
	fmt.Println(s.ContainsAll())

	// Output:
	// false
	// true
	// true
}
