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
	"encoding/json"
	"testing"

	"github.com/bytedance/gg/internal/assert"
	"github.com/bytedance/gg/internal/iter"
)

func TestLen(t *testing.T) {
	s := New[int]()
	assert.Zero(t, s.Len())
	s = New(1, 2, 3)
	assert.Equal(t, 3, s.Len())
	s = New(1, 1, 1, 1, 1, 1)
	assert.Equal(t, 1, s.Len())
}

func TestAdd(t *testing.T) {
	s := New[int]()
	assert.Zero(t, s.Len())

	assert.True(t, s.Add(1))
	assert.Equal(t, 1, s.Len())
	assert.True(t, s.Contains(1))

	assert.True(t, s.Add(2))
	assert.Equal(t, 2, s.Len())
	assert.True(t, s.Contains(2))

	assert.False(t, s.Add(1))
	assert.Equal(t, 2, s.Len())
	assert.True(t, s.Contains(1))
}

func TestAddN(t *testing.T) {
	s := New[int]()
	assert.Zero(t, s.Len())

	s.AddN(1, 2, 3, 4)
	assert.Equal(t, 4, s.Len())
	assert.True(t, s.ContainsAll(1, 2, 3, 4))

	s.AddN()
	assert.Equal(t, 4, s.Len())

	s.AddN(3, 4, 5)
	assert.Equal(t, 5, s.Len())
	assert.True(t, s.ContainsAll(1, 2, 3, 4, 5))
}

func TestRemove(t *testing.T) {
	s := New(1, 2, 3, 4)
	assert.Equal(t, 4, s.Len())

	assert.True(t, s.Remove(1))
	assert.Equal(t, 3, s.Len())
	assert.False(t, s.Contains(1))

	assert.False(t, s.Remove(1))
	assert.Equal(t, 3, s.Len())
	assert.False(t, s.Contains(1))

	assert.False(t, s.Remove(5))
	assert.Equal(t, 3, s.Len())
	assert.False(t, s.Contains(5))
}

func TestRemoveN(t *testing.T) {
	s := New(1, 2, 3, 4)
	assert.Equal(t, 4, s.Len())

	s.RemoveN()
	assert.Equal(t, 4, s.Len())

	s.RemoveN(1, 2)
	assert.Equal(t, 2, s.Len())

	s.RemoveN(1, 2, 3, 4)
	assert.Equal(t, 0, s.Len())

	s.RemoveN(1, 2, 3, 4)
	assert.Equal(t, 0, s.Len())

	s.RemoveN()
	assert.Equal(t, 0, s.Len())
}

func TestRange(t *testing.T) {
	var s []int
	New(1, 2, 3, 4).Range(func(v int) bool {
		s = append(s, v)
		return true
	})
	s = []int{3, 1, 2, 4}
	iter.Sort(iter.StealSlice(s))
	assert.Equal(t, []int{1, 2, 3, 4}, s)

	s = []int{}
	New(1, 2, 3, 4).Range(func(v int) bool {
		s = append(s, v)
		return len(s) != 3
	})
	assert.Equal(t, 3, len(s))
}

func TestUnion(t *testing.T) {
	assert.Equal(t,
		New[int](),
		New[int]().Union(New[int]()))
	assert.Equal(t,
		New(1, 2, 3),
		New[int]().Union(New(1, 2, 3)))
	assert.Equal(t,
		New(1, 2, 3),
		New(1, 2, 3).Union(New[int]()))
	assert.Equal(t,
		New(1, 2, 3, 4, 5, 6),
		New(1, 2, 3, 4).Union(New(3, 4, 5, 6)))
	assert.Equal(t,
		New(1, 2, 3, 4, 5, 6),
		New(1, 2, 3).Union(New(4, 5, 6)))
	assert.Equal(t,
		New(1, 2, 3),
		New(1, 2, 3).Union(New(1, 2, 3)))
}

func TestDiff(t *testing.T) {
	assert.Equal(t,
		New[int](),
		New[int]().Diff(New[int]()))
	assert.Equal(t,
		New[int](),
		New[int]().Diff(New(1, 2, 3)))
	assert.Equal(t,
		New(1, 2, 3),
		New(1, 2, 3).Diff(New[int]()))
	assert.Equal(t,
		New(1, 2),
		New(1, 2, 3, 4).Diff(New(3, 4, 5, 6)))
	assert.Equal(t,
		New(1, 2, 3),
		New(1, 2, 3).Diff(New(4, 5, 6)))
	assert.Equal(t,
		New[int](),
		New(1, 2, 3).Diff(New(1, 2, 3)))
}

func TestIntersect(t *testing.T) {
	assert.Equal(t,
		New[int](),
		New[int]().Intersect(New[int]()))
	assert.Equal(t,
		New[int](),
		New[int]().Intersect(New(1, 2, 3)))
	assert.Equal(t,
		New[int](),
		New(1, 2, 3).Intersect(New[int]()))
	assert.Equal(t,
		New(3, 4),
		New(1, 2, 3, 4).Intersect(New(3, 4, 5, 6)))
	assert.Equal(t,
		New[int](),
		New(1, 2, 3).Intersect(New(4, 5, 6)))
	assert.Equal(t,
		New(1, 2, 3),
		New(1, 2, 3).Intersect(New(1, 2, 3)))
}

func TestUpdate(t *testing.T) {
	s := New(1, 2, 3)
	s.Update(New(1, 2, 3))
	assert.Equal(t, New(1, 2, 3), s)
}

func TestUnionInplace(t *testing.T) {
	s := New[int]()
	s.UnionInplace(New[int]())
	assert.Equal(t, New[int](), s)

	s = New[int]()
	s.UnionInplace(New(1, 2, 3))
	assert.Equal(t, New(1, 2, 3), s)

	s = New(1, 2, 3)
	s.UnionInplace(New[int]())
	assert.Equal(t, New(1, 2, 3), s)

	s = New(1, 2, 3, 4)
	s.UnionInplace(New(3, 4, 5, 6))
	assert.Equal(t, New(1, 2, 3, 4, 5, 6), s)

	s = New(1, 2, 3)
	s.UnionInplace(New(1, 2, 3))
	assert.Equal(t, New(1, 2, 3), s)
}

func TestDiffInplace(t *testing.T) {
	s := New[int]()
	s.DiffInplace(New[int]())
	assert.Equal(t, New[int](), s)

	s = New[int]()
	s.DiffInplace(New(1, 2, 3))
	assert.Equal(t, New[int](), s)

	s = New(1, 2, 3)
	s.DiffInplace(New[int]())
	assert.Equal(t, New(1, 2, 3), s)

	s = New(1, 2, 3, 4)
	s.DiffInplace(New(3, 4, 5, 6))
	assert.Equal(t, New(1, 2), s)

	s = New(1, 2, 3)
	s.DiffInplace(New(1, 2, 3))
	assert.Equal(t, New[int](), s)
}

func TestIntersectInplace(t *testing.T) {
	s := New[int]()
	s.IntersectInplace(New[int]())
	assert.Equal(t, New[int](), s)

	s = New[int]()
	s.IntersectInplace(New(1, 2, 3))
	assert.Equal(t, New[int](), s)

	s = New(1, 2, 3)
	s.IntersectInplace(New[int]())
	assert.Equal(t, New[int](), s)

	s = New(1, 2, 3, 4)
	s.IntersectInplace(New(3, 4, 5, 6))
	assert.Equal(t, New(3, 4), s)

	s = New(1, 2, 3)
	s.IntersectInplace(New(1, 2, 3))
	assert.Equal(t, New(1, 2, 3), s)
}

func TestEqual(t *testing.T) {
	assert.True(t, New[int]().Equal(New[int]()))
	assert.True(t, New(1).Equal(New(1)))
	assert.True(t, New(1, 2, 3, 4).Equal(New(4, 3, 2, 1)))
	assert.False(t, New[int]().Equal(New(1)))
	assert.False(t, New(1).Equal(New(2)))
	assert.False(t, New(1, 2, 3, 4).Equal(New(5, 3, 2, 1)))
}

func TestIsSubset(t *testing.T) {
	assert.True(t, New[int]().IsSubset(New[int]()))
	assert.True(t, New(1).IsSubset(New(1)))
	assert.True(t, New(1, 2, 3, 4).IsSubset(New(4, 3, 2, 1)))

	assert.True(t, New[int]().IsSubset(New(1, 2, 3, 4)))
	assert.False(t, New(1, 2, 3, 4).IsSubset(New[int]()))

	assert.True(t, New(1, 2, 3).IsSubset(New(1, 2, 3, 4)))
	assert.False(t, New(1, 2, 3, 4).IsSubset(New(1, 2, 3)))

	assert.True(t, New[int]().IsSubset(New(1)))
	assert.False(t, New(1).IsSubset(New[int]()))

	assert.False(t, New(1).IsSubset(New(2)))
	assert.False(t, New(2).IsSubset(New(1)))

	assert.False(t, New(1, 2, 3, 4).IsSubset(New(5, 3, 2, 1)))
	assert.False(t, New(5, 3, 2, 1).IsSubset(New(1, 2, 3, 4)))
}

func TestIsSuperset(t *testing.T) {
	assert.True(t, New[int]().IsSuperset(New[int]()))
	assert.True(t, New(1).IsSuperset(New(1)))
	assert.True(t, New(1, 2, 3, 4).IsSuperset(New(4, 3, 2, 1)))

	assert.False(t, New[int]().IsSuperset(New(1, 2, 3, 4)))
	assert.True(t, New(1, 2, 3, 4).IsSuperset(New[int]()))

	assert.False(t, New(1, 2, 3).IsSuperset(New(1, 2, 3, 4)))
	assert.True(t, New(1, 2, 3, 4).IsSuperset(New(1, 2, 3)))

	assert.False(t, New[int]().IsSuperset(New(1)))
	assert.True(t, New(1).IsSuperset(New[int]()))

	assert.False(t, New(1).IsSuperset(New(2)))
	assert.False(t, New(2).IsSuperset(New(1)))

	assert.False(t, New(1, 2, 3, 4).IsSuperset(New(5, 3, 2, 1)))
	assert.False(t, New(5, 3, 2, 1).IsSuperset(New(1, 2, 3, 4)))
}

func TestContainsAny(t *testing.T) {
	s := New(1, 2, 3, 4)
	assert.True(t, s.ContainsAny(1))
	assert.False(t, s.ContainsAny(5))
	assert.True(t, s.ContainsAny(1, 5))
	assert.True(t, s.ContainsAny(1, 2))
	assert.False(t, s.ContainsAny(5, 6))
	assert.False(t, s.ContainsAny())
}

func TestContainsAll(t *testing.T) {
	s := New(1, 2, 3, 4)
	assert.True(t, s.ContainsAll(1))
	assert.False(t, s.ContainsAll(5))
	assert.False(t, s.ContainsAll(1, 5))
	assert.True(t, s.ContainsAll(1, 2))
	assert.False(t, s.ContainsAll(5, 6))
	assert.True(t, s.ContainsAll())
}

func TestToSlice(t *testing.T) {
	assert.Equal(t, []int{}, New[int]().ToSlice())
	assert.Equal(t, []int{1}, New(1).ToSlice())
}

func TestJSON(t *testing.T) {
	{
		// Test marshal.
		s1 := New(1, 2, 3, 4)
		bs, err := json.Marshal(s1)
		assert.Nil(t, err)
		assert.Equal(t, `[1,2,3,4]`, string(bs))

		// Test unmarshal.
		var s2 Set[int]
		err = json.Unmarshal(bs, &s2)
		assert.Nil(t, err)
		assert.True(t, s1.Equal(&s2))

		// Test overwrite.
		err = json.Unmarshal(bs, &s2)
		assert.Nil(t, err)
		assert.True(t, s1.Equal(&s2))
	}

	// Noop.
	assert.NotPanic(t, func() {
		var s Set[int]
		sp := &s
		err := json.Unmarshal([]byte("null"), sp)
		assert.Nil(t, err)
		assert.True(t, s.m == nil)
	})

	assert.NotPanic(t, func() {
		var s Set[int]
		err := json.Unmarshal([]byte("[]"), &s)
		assert.Nil(t, err)
		assert.True(t, s.m != nil && len(s.m) == 0)
	})

	// Test pointer as struct field
	{
		type Foo struct {
			Set *Set[string] `json:"set"`
		}

		f1 := Foo{New("foo", "bar")}
		bs, err := json.Marshal(f1)
		assert.Nil(t, err)
		assert.Equal(t, `{"set":["bar","foo"]}`, string(bs))

		f2 := Foo{}
		err = json.Unmarshal(bs, &f2)
		assert.Nil(t, err)
		assert.Equal(t, f2, f1)
	}
}

func TestClone(t *testing.T) {
	// { // Test nil
	// 	var s1 *Set[int]
	// 	s2 := s1.Clone()
	// 	assert.Equal(t, s1, s2)
	// }
	{ // Test empty
		s1 := New[int]()
		s2 := s1.Clone()
		assert.Equal(t, s1, s2)
		assert.Equal(t, New[int](), s1)
	}
	{
		s1 := New(1, 2, 3, 4, 5, 6)
		s2 := s1.Clone()
		assert.Equal(t, s1, s2)

		assert.Equal(t, s1, s2)
		assert.Equal(t, s1, New(1, 2, 3, 4, 5, 6))
	}
}
