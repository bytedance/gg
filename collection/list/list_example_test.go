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

package list

import (
	"fmt"
)

func Example() {
	l := New[int]()
	e1 := l.PushFront(1)        // 1
	e2 := l.PushBack(2)         // 1->2
	e3 := l.InsertBefore(3, e2) // 1->3->2
	e4 := l.InsertAfter(4, e1)  // 1->4->3->2

	l.MoveToFront(e4)    // 4->1->3->2
	l.MoveToBack(e1)     // 4->3->2->1
	l.MoveAfter(e3, e2)  // 4->2->3->1
	l.MoveBefore(e4, e1) // 2->3->4->1

	fmt.Println(l.Len())         // 4
	fmt.Println(l.Front().Value) // 2
	fmt.Println(l.Back().Value)  // 1

	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}

	// Output:
	// 4
	// 2
	// 1
	// 2
	// 3
	// 4
	// 1
}
