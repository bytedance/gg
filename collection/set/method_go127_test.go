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

package set

import (
	"strconv"
	"testing"

	"github.com/bytedance/gg/internal/assert"
)

func TestGo127SetMap(t *testing.T) {
	mapped := New(1, 2, 3).Map(strconv.Itoa)
	assert.Equal(t, 3, mapped.Len())
	assert.True(t, mapped.Contains("1"))
	assert.True(t, mapped.Contains("3"))
}

func TestGo127SetFlatMap(t *testing.T) {
	flat := New(1, 2).FlatMap(func(i int) []string {
		return []string{strconv.Itoa(i), strconv.Itoa(-i)}
	})
	assert.Equal(t, 4, flat.Len())
	assert.True(t, flat.Contains("-1"))
}

func TestGo127SetFilter(t *testing.T) {
	s := New(1, 2, 3)
	assert.Equal(t, 2, s.Filter(func(i int) bool { return i%2 == 1 }).Len())

	matched, unmatched := s.Partition(func(i int) bool { return i%2 == 1 })
	assert.Equal(t, 2, matched.Len())
	assert.Equal(t, 1, unmatched.Len())
}

func TestGo127SetPredicate(t *testing.T) {
	s := New(1, 2, 3)
	assert.True(t, s.Any(func(i int) bool { return i == 2 }))
	assert.False(t, s.All(func(i int) bool { return i == 2 }))
	assert.True(t, New[int]().All(func(i int) bool { return false }))
	assert.True(t, s.Find(func(i int) bool { return i > 1 }).IsOK())
	assert.True(t, New[int]().Find(func(i int) bool { return true }).IsNil())
}

func TestGo127SetReduce(t *testing.T) {
	assert.Equal(t, 6, New(1, 2, 3).Reduce(func(a, b int) int { return a + b }).Value())
	assert.True(t, New[int]().Reduce(func(a, b int) int { return a + b }).IsNil())
}

func TestGo127SetMinMax(t *testing.T) {
	less := func(a, b int) bool { return a < b }
	assert.Equal(t, 3, New(1, 2, 3).MaxBy(less).Value())
	assert.Equal(t, 1, New(1, 2, 3).MinBy(less).Value())
	assert.True(t, New[int]().MaxBy(less).IsNil())
	assert.True(t, New[int]().MinBy(less).IsNil())
}

func TestGo127SetGroupBy(t *testing.T) {
	grouped := New(1, 2, 3).GroupBy(func(i int) string { return strconv.Itoa(i % 2) })
	assert.Equal(t, 2, len(grouped))
	assert.Equal(t, 2, grouped["1"].Len())
	assert.Equal(t, 1, grouped["0"].Len())
}

func TestGo127SetForEach(t *testing.T) {
	sum := 0
	New(1, 2, 3).ForEach(func(i int) { sum += i })
	assert.Equal(t, 6, sum)
}

func TestGo127SetNegative(t *testing.T) {
	// The predicate never matches, so Any/All return false.
	assert.False(t, New(1, 2, 3).Any(func(i int) bool { return i > 9 }))
	assert.False(t, New(1, 2, 3).All(func(i int) bool { return i == 1 }))
	assert.True(t, New(1, 2, 3).Find(func(i int) bool { return i > 9 }).IsNil())
}

func TestGo127SetNilSafety(t *testing.T) {
	var s *Set[int]
	assert.Equal(t, 0, s.Map(strconv.Itoa).Len())
	assert.Equal(t, 0, s.Filter(func(i int) bool { return true }).Len())
	assert.False(t, s.Any(func(i int) bool { return true }))
	assert.True(t, s.All(func(i int) bool { return false }))
	assert.True(t, s.Find(func(i int) bool { return true }).IsNil())
	assert.True(t, s.Reduce(func(a, b int) int { return a + b }).IsNil())
	s.ForEach(func(i int) { t.Fatal("should not be called") })
	assert.Equal(t, 0, len(s.GroupBy(strconv.Itoa)))

	// FlatMap and Partition are nil-safe as well.
	assert.Equal(t, 0, s.FlatMap(func(i int) []string { return []string{strconv.Itoa(i)} }).Len())
	matched, unmatched := s.Partition(func(i int) bool { return true })
	assert.Equal(t, 0, matched.Len())
	assert.Equal(t, 0, unmatched.Len())
	assert.True(t, s.MaxBy(func(a, b int) bool { return a < b }).IsNil())
	assert.True(t, s.MinBy(func(a, b int) bool { return a < b }).IsNil())
}
