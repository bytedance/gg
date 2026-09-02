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

package skipmap

import (
	"testing"

	"github.com/bytedance/gg/internal/assert"
)

func TestGo127OrderedMap(t *testing.T) {
	m := New[string, int]()
	m.Store("a", 1)
	m.Store("b", 2)

	assert.Equal(t, []string{"a", "b"}, m.Keys())
	assert.Equal(t, []int{1, 2}, m.Values())
	assert.Equal(t, map[string]int{"a": 1, "b": 2}, m.ToMap())

	assert.True(t, m.Any(func(k string, v int) bool { return v == 2 }))
	assert.False(t, m.All(func(k string, v int) bool { return v == 2 }))
	assert.True(t, New[string, int]().All(func(k string, v int) bool { return false }))

	assert.Equal(t, 2, m.Find(func(k string, v int) bool { return v == 2 }).Value())
	assert.Equal(t, "b", m.FindKey(func(k string, v int) bool { return v == 2 }).Value())
	assert.True(t, New[string, int]().Find(func(k string, v int) bool { return true }).IsNil())

	sum := 0
	m.ForEach(func(k string, v int) { sum += v })
	assert.Equal(t, 3, sum)
}

func TestGo127OrderedMapDesc(t *testing.T) {
	m := NewDesc[string, int]()
	m.Store("a", 1)
	m.Store("b", 2)

	// The map is traversed in descending order.
	assert.Equal(t, []string{"b", "a"}, m.Keys())
	assert.Equal(t, []int{2, 1}, m.Values())
	assert.True(t, m.Any(func(k string, v int) bool { return v == 1 }))
}

func TestGo127SkipMapDescExtras(t *testing.T) {
	m := NewDesc[string, int]()
	m.Store("a", 1)
	m.Store("b", 2)

	assert.True(t, m.All(func(k string, v int) bool { return v > 0 }))
	assert.False(t, m.All(func(k string, v int) bool { return v == 1 }))
	assert.Equal(t, 1, m.Find(func(k string, v int) bool { return v == 1 }).Value())
	assert.Equal(t, "a", m.FindKey(func(k string, v int) bool { return v == 1 }).Value())

	sum := 0
	m.ForEach(func(k string, v int) { sum += v })
	assert.Equal(t, 3, sum)
}

func TestGo127SkipMapNegative(t *testing.T) {
	// Empty maps hit the "nothing found" branches.
	empty := New[string, int]()
	assert.False(t, empty.Any(func(k string, v int) bool { return true }))
	assert.True(t, empty.All(func(k string, v int) bool { return false }))
	assert.True(t, empty.Find(func(k string, v int) bool { return true }).IsNil())
	assert.True(t, empty.FindKey(func(k string, v int) bool { return true }).IsNil())
	empty.ForEach(func(k string, v int) { t.Fatal("should not be called") })

	// A predicate that never matches.
	m := New[string, int]()
	m.Store("a", 1)
	assert.False(t, m.Any(func(k string, v int) bool { return v > 9 }))
	assert.False(t, m.All(func(k string, v int) bool { return v > 9 }))
	assert.True(t, m.Find(func(k string, v int) bool { return v > 9 }).IsNil())
	assert.True(t, m.FindKey(func(k string, v int) bool { return v > 9 }).IsNil())
}

func TestGo127FuncMapAny(t *testing.T) {
	m := NewFunc[string, int](func(a, b string) bool { return a < b })
	m.Store("a", 1)

	assert.True(t, m.Any(func(k string, v int) bool { return v == 1 }))
	assert.False(t, m.Any(func(k string, v int) bool { return v > 9 }))
}

func TestGo127SkipMapScanBranches(t *testing.T) {
	// The Range callback has three paths: match on the first item
	// (short-circuit), match on a later one (keep scanning), and never match.
	fm := NewFunc[string, int](func(a, b string) bool { return a < b })
	fm.Store("a", 1)
	fm.Store("b", 2)
	assert.True(t, fm.All(func(k string, v int) bool { return v > 0 }))
	assert.False(t, fm.All(func(k string, v int) bool { return v == 1 }))
	assert.True(t, fm.Any(func(k string, v int) bool { return v == 2 }))
	assert.False(t, fm.Any(func(k string, v int) bool { return v > 9 }))
	assert.Equal(t, 2, fm.Find(func(k string, v int) bool { return v == 2 }).Value())

	m := New[string, int]()
	m.Store("a", 1)
	m.Store("b", 2)
	assert.False(t, m.All(func(k string, v int) bool { return v == 1 }))
	assert.True(t, m.Any(func(k string, v int) bool { return v == 2 }))
	assert.False(t, m.Any(func(k string, v int) bool { return v > 9 }))

	dm := NewDesc[string, int]()
	dm.Store("a", 1)
	dm.Store("b", 2)
	assert.False(t, dm.All(func(k string, v int) bool { return v == 2 }))
	assert.True(t, dm.Any(func(k string, v int) bool { return v == 1 }))
	assert.False(t, dm.Any(func(k string, v int) bool { return v > 9 }))
}

func TestGo127FuncMap(t *testing.T) {
	m := NewFunc[string, int](func(a, b string) bool { return a < b })
	m.Store("a", 1)
	m.Store("b", 2)

	assert.Equal(t, []string{"a", "b"}, m.Keys())
	assert.Equal(t, []int{1, 2}, m.Values())
	assert.Equal(t, 2, m.Find(func(k string, v int) bool { return v == 2 }).Value())
	assert.Equal(t, "b", m.FindKey(func(k string, v int) bool { return v == 2 }).Value())
	assert.True(t, m.All(func(k string, v int) bool { return v > 0 }))

	sum := 0
	m.ForEach(func(k string, v int) { sum += v })
	assert.Equal(t, 3, sum)
}
