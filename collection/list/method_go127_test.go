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

//go:build go1.27

package list

import (
	"strconv"
	"testing"

	"github.com/bytedance/gg/internal/assert"
)

func newList(vs ...int) *List[int] {
	l := New[int]()
	for _, v := range vs {
		l.PushBack(v)
	}
	return l
}

func TestGo127ListToSlice(t *testing.T) {
	assert.Equal(t, []int{1, 2, 3}, newList(1, 2, 3).ToSlice())
	assert.Equal(t, 0, len(New[int]().ToSlice()))
}

func TestGo127ListMap(t *testing.T) {
	assert.Equal(t, []string{"1", "2", "3"}, newList(1, 2, 3).Map(strconv.Itoa).ToSlice())
}

func TestGo127ListFilter(t *testing.T) {
	assert.Equal(t, []int{2}, newList(1, 2, 3).Filter(func(i int) bool { return i == 2 }).ToSlice())
}

func TestGo127ListPredicate(t *testing.T) {
	l := newList(1, 2, 3)
	assert.True(t, l.Any(func(i int) bool { return i == 2 }))
	assert.False(t, l.All(func(i int) bool { return i == 2 }))
	assert.Equal(t, 2, l.Find(func(i int) bool { return i > 1 }).Value())
	assert.True(t, New[int]().Find(func(i int) bool { return true }).IsNil())
}

func TestGo127ListReduce(t *testing.T) {
	assert.Equal(t, 6, newList(1, 2, 3).Reduce(func(a, b int) int { return a + b }).Value())
	assert.True(t, New[int]().Reduce(func(a, b int) int { return a + b }).IsNil())
	// Single element list returns the element itself.
	assert.Equal(t, 1, newList(1).Reduce(func(a, b int) int { return a + b }).Value())
}

func TestGo127ListPartition(t *testing.T) {
	matched, unmatched := newList(1, 2, 3).Partition(func(i int) bool { return i%2 == 1 })
	assert.Equal(t, []int{1, 3}, matched.ToSlice())
	assert.Equal(t, []int{2}, unmatched.ToSlice())
}

func TestGo127ListForEach(t *testing.T) {
	sum := 0
	newList(1, 2, 3).ForEach(func(i int) { sum += i })
	assert.Equal(t, 6, sum)
}
