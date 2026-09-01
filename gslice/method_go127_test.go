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

package gslice

import (
	"errors"
	"strconv"
	"testing"

	"github.com/bytedance/gg/gvalue"
	"github.com/bytedance/gg/internal/assert"
)

func TestWrap(t *testing.T) {
	assert.Equal(t, []int{1, 2, 3}, Wrap([]int{1, 2, 3}).Unwrap())
	assert.Equal(t, 3, Wrap([]int{1, 2, 3}).Len())
	assert.Equal(t, []int{1, 2}, New(1, 2).Unwrap())
	assert.Equal(t, []int{1, 2, 3, 4}, Wrap(Range(1, 5)).Unwrap())
	assert.Equal(t, []int{5, 3}, Wrap(RangeWithStep(5, 1, -2)).Unwrap())
}

func TestSMap(t *testing.T) {
	assert.Equal(t, []string{"1", "2", "3"}, Wrap([]int{1, 2, 3}).Map(strconv.Itoa).Unwrap())
	assert.Equal(t, []string{"0:1", "1:2"}, Wrap([]int{1, 2}).MapIndexed(func(v, i int) string {
		return strconv.Itoa(i) + ":" + strconv.Itoa(v)
	}).Unwrap())

	// TryMap keeps the error.
	assert.Equal(t, []string{"1"}, Wrap([]int{1}).TryMap(func(i int) (string, error) {
		return strconv.Itoa(i), nil
	}).Value().Unwrap())
	assert.True(t, Wrap([]int{1}).TryMap(func(i int) (string, error) {
		return "", errors.New("woof!")
	}).IsErr())
}

func TestSFilter(t *testing.T) {
	isEven := func(i int) bool { return i%2 == 0 }
	assert.Equal(t, []int{2, 4}, Wrap([]int{1, 2, 3, 4}).Filter(isEven).Unwrap())
	assert.Equal(t, []int{1, 3}, Wrap([]int{1, 2, 3, 4}).Reject(isEven).Unwrap())
	assert.Equal(t, []string{"2", "4"}, Wrap([]int{1, 2, 3, 4}).FilterMap(func(i int) (string, bool) {
		return strconv.Itoa(i), i%2 == 0
	}).Unwrap())

	matched, unmatched := Wrap([]int{1, 2, 3, 4}).Partition(isEven)
	assert.Equal(t, []int{2, 4}, matched.Unwrap())
	assert.Equal(t, []int{1, 3}, unmatched.Unwrap())
}

func TestSReduce(t *testing.T) {
	assert.Equal(t, 6, Wrap([]int{1, 2, 3}).Reduce(gvalue.Add[int]).Value())
	assert.True(t, Wrap([]int{}).Reduce(gvalue.Add[int]).IsNil())
	assert.Equal(t, "123", Wrap([]int{1, 2, 3}).Fold(func(acc string, i int) string {
		return acc + strconv.Itoa(i)
	}, ""))
}

func TestSGet(t *testing.T) {
	s := Wrap([]int{1, 2, 3})
	assert.Equal(t, 1, s.First().Value())
	assert.Equal(t, 3, s.Last().Value())
	assert.Equal(t, 2, s.Get(1).Value())
	assert.Equal(t, 3, s.Get(-1).Value())
	assert.True(t, s.Get(9).IsNil())
	assert.Equal(t, 2, s.Find(func(i int) bool { return i > 1 }).Value())
	assert.Equal(t, 1, s.IndexBy(func(i int) bool { return i > 1 }).Value())
	assert.True(t, s.Any(func(i int) bool { return i == 2 }))
	assert.False(t, s.All(func(i int) bool { return i == 2 }))
	assert.Equal(t, 2, s.CountBy(func(i int) bool { return i > 1 }))
}

func TestSChunk(t *testing.T) {
	chunks := Wrap([]int{1, 2, 3, 4, 5}).Chunk(2)
	assert.Equal(t, 3, len(chunks))
	assert.Equal(t, []int{1, 2}, chunks[0].Unwrap())
	assert.Equal(t, []int{5}, chunks[2].Unwrap())

	divided := Wrap([]int{1, 2, 3, 4, 5}).Divide(2)
	assert.Equal(t, []int{1, 2, 3}, divided[0].Unwrap())
	assert.Equal(t, []int{4, 5}, divided[1].Unwrap())
}

func TestSGroupBy(t *testing.T) {
	grouped := Wrap([]int{1, 2}).GroupBy(strconv.Itoa)
	assert.Equal(t, 2, len(grouped))
	assert.Equal(t, []int{1}, grouped["1"].Unwrap())
	assert.Equal(t, []int{2}, grouped["2"].Unwrap())

	assert.Equal(t, []int{1, 2}, Wrap([]int{1, 1, 2}).UniqBy(strconv.Itoa).Unwrap())
	assert.Equal(t, []int{1}, Wrap([]int{1, 1, 2}).DupBy(strconv.Itoa).Unwrap())
	assert.Equal(t, map[string]int{"1": 2, "2": 1}, Wrap([]int{1, 1, 2}).CountValuesBy(strconv.Itoa))
}

func TestSTake(t *testing.T) {
	s := Wrap([]int{1, 2, 3, 4, 5})
	assert.Equal(t, []int{1, 2}, s.Take(2).Unwrap())
	assert.Equal(t, []int{4, 5}, s.Take(-2).Unwrap())
	assert.Equal(t, []int{2, 3}, s.Slice(1, 3).Unwrap())
	assert.Equal(t, []int{3, 4, 5}, s.Drop(2).Unwrap())
	assert.Equal(t, []int{1, 2, 4, 5}, s.RemoveIndex(2).Unwrap())
	assert.Equal(t, []int{1, 9, 2, 3, 4, 5}, s.Insert(1, 9).Unwrap())
	assert.Equal(t, []int{1, 2, 3}, Wrap([]int{1, 2}).Concat([]int{3}).Unwrap())
}

func TestSOrder(t *testing.T) {
	assert.Equal(t, []int{3, 2, 1}, Wrap([]int{1, 2, 3}).Reverse().Unwrap())
	assert.Equal(t, []int{3, 2, 1}, Wrap([]int{1, 2, 3}).ReverseClone().Unwrap())
	assert.Equal(t, []int{1, 2, 3}, Wrap([]int{3, 1, 2}).SortBy(gvalue.Less[int]).Unwrap())
	assert.Equal(t, []int{3, 2, 1}, Wrap([]int{1, 2, 3}).SortCloneBy(gvalue.Greater[int]).Unwrap())
	assert.Equal(t, []int{1, 2, 3}, Wrap([]int{3, 1, 2}).StableSortBy(gvalue.Less[int]).Unwrap())
	assert.Equal(t, 3, Wrap([]int{1, 2, 3}).Shuffle().Len())
	assert.Equal(t, 3, Wrap([]int{1, 2, 3}).ShuffleClone().Len())
	assert.Equal(t, []int{1, 2, 3}, Wrap([]int{1, 2, 3}).Clone().Unwrap())
	assert.Equal(t, []int{2, 4, 6}, Wrap([]int{1, 2, 3}).CloneBy(func(i int) int { return i * 2 }).Unwrap())
}

func TestSConvert(t *testing.T) {
	assert.Equal(t, map[string]int{"1": 1}, Wrap([]int{1}).ToMapValues(strconv.Itoa))
	assert.Equal(t, map[string]string{"1": "1"}, Wrap([]int{1}).ToMap(func(i int) (string, string) {
		return strconv.Itoa(i), strconv.Itoa(i)
	}))
	assert.Equal(t, []int{1}, Wrap([]any{1}).TypeAssert[int]().Unwrap())
	assert.Equal(t, []string{"1"}, Wrap([]int{1}).FlatMap(func(i int) []string {
		return []string{strconv.Itoa(i)}
	}).Unwrap())
	assert.Equal(t, 6, Wrap([]int{1, 2, 3}).SumBy(func(i int) int { return i }))
	assert.Equal(t, 2.0, Wrap([]int{1, 2, 3}).AvgBy(func(i int) int { return i }))
}

func TestSChaining(t *testing.T) {
	// The key benefit of generic methods: the element type can change in a chain.
	got := Wrap(Range(1, 6)).
		Filter(func(i int) bool { return i%2 == 1 }).
		Map(strconv.Itoa).
		Fold(func(acc string, s string) string { return acc + s }, "")
	assert.Equal(t, "135", got)

	// Constrained wrappers are converted back to S to keep chaining.
	assert.Equal(t,
		[]string{"2", "4"},
		WrapC([]int{1, 2, 2, 3, 4}).Uniq().S().Filter(func(i int) bool { return i%2 == 0 }).Map(strconv.Itoa).Unwrap(),
	)
}

func TestC(t *testing.T) {
	c := WrapC([]int{1, 2, 2, 3})
	assert.True(t, c.Contains(2))
	assert.False(t, c.Contains(9))
	assert.True(t, c.ContainsAny(3, 9))
	assert.False(t, c.ContainsAll(3, 9))
	assert.Equal(t, 1, c.Index(2).Value())
	assert.Equal(t, 2, c.IndexRev(2).Value())
	assert.Equal(t, []int{1, 2, 3}, c.Uniq().Unwrap())
	assert.Equal(t, []int{2}, c.Dup().Unwrap())
	assert.Equal(t, []int{1, 3}, c.Remove(2).Unwrap())
	assert.Equal(t, 2, c.Count(2))
	assert.Equal(t, map[int]int{1: 1, 2: 2, 3: 1}, c.CountValues())
	assert.Equal(t, map[int]bool{1: true, 2: true, 3: true}, c.ToBoolMap())
	assert.True(t, c.Equal([]int{1, 2, 2, 3}))
	assert.Equal(t, []int{1, 2, 2, 3}, WrapC([]int{0, 1, 0, 2, 2, 3}).Compact().Unwrap())
}

func TestCSetOp(t *testing.T) {
	// The order of the result of the set operations is not guaranteed,
	// so the results are sorted before comparing.
	sorted := func(s []int) []int { return WrapO(s).Sort().Unwrap() }

	assert.Equal(t, []int{1, 2, 3, 4}, sorted(WrapC([]int{1, 2, 2}).Union(WrapC([]int{3, 4})).Unwrap()))
	assert.Equal(t, []int{1, 3}, WrapC([]int{1, 2, 3}).Diff(WrapC([]int{2})).Unwrap())
	assert.Equal(t, []int{2, 3}, sorted(WrapC([]int{1, 2, 3}).Intersect(WrapC([]int{2, 3})).Unwrap()))
}

func TestO(t *testing.T) {
	assert.Equal(t, 3, WrapO([]int{3, 1, 2}).Max().Value())
	assert.Equal(t, 1, WrapO([]int{3, 1, 2}).Min().Value())

	min, max := WrapO([]int{3, 1, 2}).MinMax().Value().Values()
	assert.Equal(t, 1, min)
	assert.Equal(t, 3, max)

	assert.Equal(t, []int{1, 2, 3}, WrapO([]int{3, 1, 2}).Sort().Unwrap())
	assert.Equal(t, []int{1, 2, 3}, WrapO([]int{3, 1, 2}).SortClone().Unwrap())
	assert.Equal(t, 3, WrapO([]int{3, 1, 2}).PartialSort(3).Len())
	assert.Equal(t, 3, WrapO([]int{3, 1, 2}).Len())

	// Ordered elements are comparable and convertible to S and C.
	assert.Equal(t, []string{"1", "2", "3"}, WrapO([]int{3, 1, 2}).Sort().S().Map(strconv.Itoa).Unwrap())
	assert.True(t, WrapO([]int{3, 1, 2}).C().Contains(3))
}

func TestN(t *testing.T) {
	n := WrapN([]int{1, 2, 3})
	assert.Equal(t, 6, n.Sum())
	assert.Equal(t, 2.0, n.Avg())
	assert.Equal(t, []int{1, 2, 3}, n.S().Unwrap())
	assert.True(t, n.C().Contains(2))
	assert.Equal(t, 3, n.O().Max().Value())
	assert.Equal(t, 3, n.Len())
}
