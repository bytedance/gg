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

package skipset

import (
	"testing"

	"github.com/bytedance/gg/internal/assert"
)

func TestGo127OrderedSet(t *testing.T) {
	s := New[int]()
	s.Add(1)
	s.Add(2)
	s.Add(3)

	assert.True(t, s.Any(func(i int) bool { return i == 2 }))
	assert.False(t, s.All(func(i int) bool { return i == 2 }))
	assert.True(t, New[int]().All(func(i int) bool { return false }))

	assert.True(t, s.Find(func(i int) bool { return i > 1 }).IsOK())
	assert.True(t, New[int]().Find(func(i int) bool { return true }).IsNil())

	assert.Equal(t, 6, s.Reduce(func(a, b int) int { return a + b }).Value())
	assert.True(t, New[int]().Reduce(func(a, b int) int { return a + b }).IsNil())

	sum := 0
	s.ForEach(func(i int) { sum += i })
	assert.Equal(t, 6, sum)
}

func TestGo127OrderedSetDesc(t *testing.T) {
	s := NewDesc[int]()
	s.Add(1)
	s.Add(2)

	assert.True(t, s.Any(func(i int) bool { return i == 2 }))
	assert.Equal(t, 3, s.Reduce(func(a, b int) int { return a + b }).Value())
	assert.Equal(t, 2, s.Find(func(i int) bool { return i > 0 }).Value())
}

func TestGo127FuncSet(t *testing.T) {
	s := NewFunc[int](func(a, b int) bool { return a < b })
	s.Add(1)
	s.Add(2)

	assert.True(t, s.Any(func(i int) bool { return i == 1 }))
	assert.False(t, s.All(func(i int) bool { return i == 1 }))
	assert.Equal(t, 2, s.Find(func(i int) bool { return i > 1 }).Value())
	assert.Equal(t, 3, s.Reduce(func(a, b int) int { return a + b }).Value())

	sum := 0
	s.ForEach(func(i int) { sum += i })
	assert.Equal(t, 3, sum)
}

func TestGo127SetFindOrder(t *testing.T) {
	// The set is traversed in ascending order, so Find returns the smallest
	// element that satisfies the predicate.
	s := New[int]()
	s.Add(3)
	s.Add(1)
	s.Add(2)
	assert.Equal(t, 2, s.Find(func(i int) bool { return i > 1 }).Value())
}
