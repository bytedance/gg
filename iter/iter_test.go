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

package iter

import (
	"container/list"
	"fmt"
	"strconv"

	"github.com/bytedance/gg/gvalue"
)

type listIter[T any] struct {
	e *list.Element
}

func FromList[T any](l *list.List) Iter[T] {
	return &listIter[T]{l.Front()}
}

func (i *listIter[T]) Next(n int) []T {
	var next []T
	j := 0
	for i.e != nil {
		next = append(next, i.e.Value.(T))
		i.e = i.e.Next()
		j++
		if n != ALL && j >= n {
			break
		}
	}
	return next
}

// Iter for container list.
func ExampleIter_impl() {
	l := list.New()
	l.PushBack(0)
	l.PushBack(1)
	l.PushBack(2)

	s := ToSlice(
		Filter(gvalue.IsNotZero[int],
			FromList[int](l)))
	fmt.Println(s)

	// Output:
	// [1 2]
}

// Convert an int slice to string slice.
func Example() {
	s := ToSlice(
		Map(strconv.Itoa,
			Filter(gvalue.IsNotZero[int],
				FromSlice([]int{0, 1, 2, 3, 4}))))
	fmt.Printf("%q\n", s)

	// Output:
	// ["1" "2" "3" "4"]
}
