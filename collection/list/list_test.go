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
	"testing"

	"github.com/bytedance/gg/internal/assert"
)

func TestNew(t *testing.T) {
	l := New[int]()
	assert.NotNil(t, l)
	assert.Equal(t, 0, l.Len())
	assert.Nil(t, l.Front())
	assert.Nil(t, l.Back())
}

func TestListPush(t *testing.T) {
	l := New[int]()
	l.PushBack(3)  // 3
	l.PushBack(4)  // 3->4
	l.PushFront(1) // 1->3->4
	l.PushFront(2) // 2->1->3->4

	assert.Equal(t, 4, l.Len())
	assert.Equal(t, 2, l.Front().Value)
	assert.Equal(t, 1, l.Front().Next().Value)
	assert.Equal(t, 3, l.Front().Next().Next().Value)
	assert.Equal(t, 4, l.Front().Next().Next().Next().Value)
	assert.Nil(t, l.Front().Next().Next().Next().Next())
	assert.Equal(t, 4, l.Back().Value)
	assert.Equal(t, 3, l.Back().Prev().Value)
	assert.Equal(t, 1, l.Back().Prev().Prev().Value)
	assert.Equal(t, 2, l.Back().Prev().Prev().Prev().Value)
	assert.Nil(t, l.Back().Prev().Prev().Prev().Prev())
}

func TestListRemove(t *testing.T) {
	l := New[int]()
	e := l.PushBack(1)
	l.Remove(e)

	assert.Equal(t, 0, l.Len())
	assert.Nil(t, l.Front())
	assert.Nil(t, l.Back())
}

func TestListInsert(t *testing.T) {
	l := New[int]()
	e1 := l.PushBack(2)         // 2
	e2 := l.InsertBefore(1, e1) // 1->2
	l.InsertAfter(3, e1)        // 1->2->3
	l.InsertBefore(4, e2)       // 4->1->2->3

	assert.Equal(t, 4, l.Len())
	assert.Equal(t, 4, l.Front().Value)
	assert.Equal(t, 1, l.Front().Next().Value)
	assert.Equal(t, 2, l.Back().Prev().Value)
	assert.Equal(t, 3, l.Back().Value)
}

func TestListMoveToFront(t *testing.T) {
	l := New[int]()
	l.PushBack(1)
	l.PushBack(2)
	e := l.PushBack(3)
	l.MoveToFront(e)

	assert.Equal(t, 3, l.Front().Value)
	assert.Equal(t, 2, l.Back().Value)
}

func TestListMove(t *testing.T) {
	l := New[int]()
	e1 := l.PushBack(1) // 1
	e2 := l.PushBack(2) // 1->2
	e3 := l.PushBack(3) // 1->2->3
	e4 := l.PushBack(4) // 1->2->3->4
	e5 := l.PushBack(5) // 1->2->3->4->5

	l.MoveToBack(e1)     // 2->3->4->5->1
	l.MoveToFront(e5)    // 5->2->3->4->1
	l.MoveAfter(e2, e3)  // 5->3->2->4->1
	l.MoveBefore(e4, e3) // 5->4->3->2->1

	assert.Equal(t, 5, l.Front().Value)
	assert.Equal(t, 4, l.Front().Next().Value)
	assert.Equal(t, 3, l.Front().Next().Next().Value)
	assert.Equal(t, 2, l.Back().Prev().Value)
	assert.Equal(t, 1, l.Back().Value)
}

func TestListPushBackList(t *testing.T) {
	l1 := New[int]()
	l1.PushBack(1) // 1
	l1.PushBack(2) // 1->2

	l2 := New[int]()
	l2.PushBack(3)      // 3
	l2.PushBackList(l1) // 3->1->2

	assert.Equal(t, 3, l2.Len())
	assert.Equal(t, 3, l2.Front().Value)
	assert.Equal(t, 2, l2.Back().Value)
}

func TestListPushFrontList(t *testing.T) {
	l1 := New[int]()
	l1.PushBack(1) // 1
	l1.PushBack(2) // 1->2

	l2 := New[int]()
	l2.PushBack(3) // 3
	l2.PushFrontList(l1)

	assert.Equal(t, 3, l2.Len())
	assert.Equal(t, 1, l2.Front().Value)
	assert.Equal(t, 3, l2.Back().Value)
}
