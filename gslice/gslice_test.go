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

package gslice

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"unsafe"

	"github.com/bytedance/gg/collection/tuple"
	"github.com/bytedance/gg/goption"
	"github.com/bytedance/gg/gptr"
	"github.com/bytedance/gg/gresult"
	"github.com/bytedance/gg/gvalue"
	"github.com/bytedance/gg/internal/assert"
)

func TestMap(t *testing.T) {
	assert.Equal(t,
		[]string{"1", "2", "3"},
		Map([]int{1, 2, 3}, strconv.Itoa))
}

func TestTryMap(t *testing.T) {
	assert.Equal(t, TryMap(Of("1", "2", "3"), strconv.Atoi), gresult.OK(Of(1, 2, 3)))
	assert.Equal(t, TryMap(nil, strconv.Atoi), gresult.OK(([]int{})))
	assert.Equal(t, TryMap(Of("1", "2", "a"), strconv.Atoi).Err().Error(), "strconv.Atoi: parsing \"a\": invalid syntax")
}

func TestFilter(t *testing.T) {
	assert.Equal(t,
		[]int{},
		Filter([]int(nil), gvalue.IsZero[int]))
	assert.Equal(t,
		[]int{},
		Filter([]int{}, gvalue.IsZero[int]))
	assert.Equal(t,
		[]int{1, 2, 3},
		Filter([]int{0, 1, 2, 3}, gvalue.IsNotZero[int]))
	assert.Equal(t,
		[]int{0},
		Filter([]int{0, 1, 2, 3}, gvalue.IsZero[int]))
}

func TestFilterMap(t *testing.T) {
	assert.Equal(t,
		[]string{"1", "2", "3"},
		FilterMap([]int{1, 2, 3, 0, 0}, func(i int) (string, bool) {
			return strconv.Itoa(i), i != 0
		}),
	)

	assert.Equal(t,
		[]string{},
		FilterMap([]int{0, 0}, func(i int) (string, bool) {
			return strconv.Itoa(i), i != 0
		}),
	)

	assert.Equal(t,
		[]string{},
		FilterMap([]int{}, func(i int) (string, bool) {
			return strconv.Itoa(i), i != 0
		}),
	)
}

func TestTryFilterMap(t *testing.T) {
	assert.Equal(t, TryFilterMap(Of("1", "2", "3"), strconv.Atoi), Of(1, 2, 3))
	assert.Equal(t, TryFilterMap(Of("1", "2", "a"), strconv.Atoi), Of(1, 2))
	assert.Equal(t, TryFilterMap(Of("1", "a", "3"), strconv.Atoi), Of(1, 3))
	assert.Equal(t, TryFilterMap(Of("a", "2", "3"), strconv.Atoi), Of(2, 3))
	assert.Equal(t, TryFilterMap(Of("a", "a", "a"), strconv.Atoi), []int{})
	assert.Equal(t, TryFilterMap(nil, strconv.Atoi), []int{})
}

func TestReject(t *testing.T) {
	assert.Equal(t,
		[]int{},
		Reject([]int(nil), gvalue.IsZero[int]))
	assert.Equal(t,
		[]int{},
		Reject([]int{}, gvalue.IsZero[int]))
	assert.Equal(t,
		[]int{1, 2, 3},
		Reject([]int{0, 1, 2, 3}, gvalue.IsZero[int]))
	assert.Equal(t,
		[]int{0},
		Reject([]int{0, 1, 2, 3}, gvalue.IsNotZero[int]))
}

func TestPartition(t *testing.T) {
	{
		filter, reject := Partition([]int(nil), gvalue.IsZero[int])
		assert.Equal(t, []int{}, filter)
		assert.Equal(t, []int{}, reject)
	}
	{
		filter, reject := Partition([]int{}, gvalue.IsZero[int])
		assert.Equal(t, []int{}, filter)
		assert.Equal(t, []int{}, reject)
	}
	{
		filter, reject := Partition([]int{0, 1, 2, 3}, gvalue.IsNotZero[int])
		assert.Equal(t, []int{1, 2, 3}, filter)
		assert.Equal(t, []int{0}, reject)
	}
	{
		filter, reject := Partition([]int{0, 1, 2, 3}, gvalue.IsZero[int])
		assert.Equal(t, []int{0}, filter)
		assert.Equal(t, []int{1, 2, 3}, reject)
	}
}

func TestReduce(t *testing.T) {
	assert.Equal(t, 6, Reduce([]int{0, 1, 2, 3}, gvalue.Add[int]).Value())
	assert.False(t, Reduce([]int{}, gvalue.Add[int]).IsOK())
}

func TestFold(t *testing.T) {
	assert.Equal(t, 10, Fold([]int{0, 1, 2, 3}, gvalue.Add[int], 4))
	assert.Equal(t, 1, Fold([]int{}, gvalue.Add[int], 1))
}

func TestChunk(t *testing.T) {
	{
		s := []int{0, 1, 2, 3, 4}
		chunks := Chunk(s, 2)
		assert.Equal(t, [][]int{{0, 1}, {2, 3}, {4}}, chunks)
		chunks[1][1] = 9 // Modify original slice
		assert.Equal(t, []int{0, 1, 2, 9, 4}, s)
	}
}

func TestChunkClone(t *testing.T) {
	{
		s := []int{0, 1, 2, 3, 4}
		chunks := ChunkClone(s, 2)
		assert.Equal(t, [][]int{{0, 1}, {2, 3}, {4}}, chunks)
		chunks[1][1] = 9
		assert.Equal(t, []int{0, 1, 2, 3, 4}, s)
	}
}

func TestGroupBy(t *testing.T) {
	assert.Equal(t,
		map[string][]int{"odd": {1, 3}, "even": {2, 4}},
		GroupBy([]int{1, 2, 3, 4},
			func(v int) string {
				if v%2 == 0 {
					return "even"
				} else {
					return "odd"
				}
			}))

	// Test custom type.
	type IntSlice []int
	assert.Equal(t,
		map[string]IntSlice{},
		GroupBy(IntSlice{}, func(n int) string { return strconv.Itoa(n) }))

}

func TestContains(t *testing.T) {
	assert.True(t, Contains([]int{0, 1, 2, 3, 4}, 0))
	assert.False(t, Contains([]int{0, 1, 2, 3, 4}, 5))
	assert.False(t, Contains([]int{}, 5))
}

func TestContainsAny(t *testing.T) {
	assert.True(t, ContainsAny([]int{0, 1, 2, 3, 4}, 0))
	assert.False(t, Contains([]int{0, 1, 2, 3, 4}, 5))
	assert.True(t, ContainsAny([]int{0, 1, 2, 3, 4}, 0, 1))
	assert.True(t, ContainsAny([]int{0, 1, 2, 3, 4}, 0, 5))
}

func TestContainsAll(t *testing.T) {
	assert.True(t, ContainsAll([]int{0, 1, 2, 3, 4}, 0))
	assert.False(t, Contains([]int{0, 1, 2, 3, 4}, 5))
	assert.True(t, ContainsAll([]int{0, 1, 2, 3, 4}, 0, 1))
	assert.False(t, ContainsAll([]int{0, 1, 2, 3, 4}, 0, 5))
}

func TestRemove(t *testing.T) {
	assert.Equal(t, []int{0, 1, 2, 4}, Remove([]int{0, 1, 2, 3, 4}, 3))
	assert.Equal(t, []int{0, 1, 2, 4}, Remove([]int{0, 1, 3, 2, 3, 4}, 3))
}

func TestUniq(t *testing.T) {
	assert.Equal(t, []int{0, 1, 4, 3},
		Uniq([]int{0, 1, 4, 3, 1, 4}))
}

func TestUniqBy(t *testing.T) {
	type Foo struct{ Value int }
	assert.Equal(t, []Foo{{0}, {1}, {4}, {3}},
		UniqBy([]Foo{{0}, {1}, {4}, {3}, {1}, {4}},
			func(v Foo) int { return v.Value }))
}

func TestDup(t *testing.T) {
	assert.Equal(t, []int{1},
		Dup([]int{0, 1, 1, 1, 1}))

	assert.Equal(t, []int{2, 3},
		Dup([]int{3, 2, 2, 3, 3}))

	assert.Equal(t, []int{1, 4},
		Dup([]int{0, 1, 4, 3, 1, 4}))
}

func TestDupBy(t *testing.T) {
	type Foo struct{ Value int }
	assert.Equal(t, []Foo{{1}},
		DupBy([]Foo{{0}, {1}, {1}, {1}, {1}},
			func(v Foo) int { return v.Value }))

	assert.Equal(t, []Foo{{2}, {3}},
		DupBy([]Foo{{3}, {2}, {2}, {3}, {3}},
			func(v Foo) int { return v.Value }))

	assert.Equal(t, []Foo{{1}, {4}},
		DupBy([]Foo{{0}, {1}, {4}, {3}, {1}, {4}},
			func(v Foo) int { return v.Value }))

}

func TestRepeat(t *testing.T) {
	{
		assert.Panic(t, func() { _ = Repeat(123, -1) })
		assert.Equal(t, Repeat(123, 0), []int{})
		assert.Equal(t, Repeat(123, 3), []int{123, 123, 123})
	}
	// test shallow copy
	{
		type testStruct struct {
			Value int
		}
		v := &testStruct{Value: 123}
		repeat := Repeat(v, 3)
		assert.Equal(t, repeat, []*testStruct{{Value: 123}, {Value: 123}, {Value: 123}})
		repeat[1].Value = 456
		assert.Equal(t, repeat, []*testStruct{{Value: 456}, {Value: 456}, {Value: 456}})
	}
}

func TestRepeatBy(t *testing.T) {
	{
		fn := func() int { return 123 }
		assert.Panic(t, func() { _ = RepeatBy(fn, -1) })
		assert.Equal(t, RepeatBy(fn, 0), []int{})
		assert.Equal(t, RepeatBy(fn, 3), []int{123, 123, 123})
	}
	// test deep copy
	{
		addrs := func(s []*int) []unsafe.Pointer {
			r := make([]unsafe.Pointer, 0, len(s))
			for _, elem := range s {
				r = append(r, unsafe.Pointer(elem))
			}
			return r
		}
		fnAddr := func() *int { return gptr.Of(123) }
		assert.Panic(t, func() { _ = RepeatBy(fnAddr, -1) })
		assert.Equal(t, RepeatBy(fnAddr, 0), []*int{})
		lhs, rhs := RepeatBy(fnAddr, 3), RepeatBy(fnAddr, 3)
		// have same value
		assert.Equal(t, lhs, rhs)
		// but with different addresses
		assert.NotEqual(t, addrs(lhs), addrs(rhs))
	}
}

func TestMax(t *testing.T) {
	assert.Equal(t, 4, Max([]int{0, 1, 4, 3, 1, 4}).Value())
}

func TestMaxBy(t *testing.T) {
	type Foo struct {
		Value int
	}
	// FIXME: use greflect
	less := func(foo1, foo2 Foo) bool {
		return foo1.Value < foo2.Value
	}
	assert.Equal(t,
		Foo{100},
		MaxBy([]Foo{{10}, {1}, {-1}, {100}, {3}}, less).Value())
}

func TestMin(t *testing.T) {
	assert.Equal(t, 0, Min([]int{0, 1, 4, 3, 1, 4}).Value())
}

func TestMinBy(t *testing.T) {
	type Foo struct {
		Value int
	}
	// FIXME: use greflect
	less := func(foo1, foo2 Foo) bool {
		return foo1.Value < foo2.Value
	}
	assert.Equal(t,
		Foo{-1},
		MinBy([]Foo{{10}, {1}, {-1}, {100}, {3}}, less).Value())
}

func TestMinMax(t *testing.T) {
	assert.Equal(t, tuple.Make2(1, 1), MinMax([]int{1}).Value())
	assert.Equal(t, tuple.Make2(0, 4), MinMax([]int{0, 1, 4, 3, 1, 4}).Value())
}

func TestMinMaxBy(t *testing.T) {
	type Foo struct {
		Value int
	}
	// FIXME: use greflect
	less := func(foo1, foo2 Foo) bool {
		return foo1.Value < foo2.Value
	}
	assert.Equal(t,
		tuple.Make2(Foo{-1}, Foo{100}),
		MinMaxBy([]Foo{{10}, {1}, {-1}, {100}, {3}}, less).Value())
}

func TestClone(t *testing.T) {
	assert.Equal(t, []int{0, 1, 4, 3, 1, 4}, Clone([]int{0, 1, 4, 3, 1, 4}))
	assert.True(t, Clone[[]int](nil) == nil)

	// Test new type.
	type Ints []int
	assert.Equal(t, Ints{0, 1, 4, 3, 1, 4}, Clone(Ints{0, 1, 4, 3, 1, 4}))
	assert.Equal(t, "gslice.Ints", fmt.Sprintf("%T", Clone(Ints{0, 1, 4, 3, 1, 4})))

	// Test shallow clone.
	src := []*int{gptr.Of(1), gptr.Of(2)}
	dst := Clone(src)
	assert.Equal(t, src, dst)
	assert.False(t, overlaps(src, dst))
	assert.True(t, src[0] == dst[0])
	assert.True(t, src[1] == dst[1])
}

func TestCloneBy(t *testing.T) {
	id := func(v int) int { return v }
	assert.Equal(t, []int{0, 1, 4, 3, 1, 4},
		CloneBy([]int{0, 1, 4, 3, 1, 4}, id))

	type Ints []int
	assert.Equal(t, Ints{0, 1, 4, 3, 1, 4},
		CloneBy(Ints{0, 1, 4, 3, 1, 4}, id))
	assert.Equal(t, "gslice.Ints", fmt.Sprintf("%T", CloneBy(Ints{0, 1, 4, 3, 1, 4}, id)))

	assert.True(t, CloneBy[[]int](nil, nil) == nil)

	// Test deep clone.
	src := []*int{gptr.Of(1), gptr.Of(2)}
	dst := CloneBy(src, gptr.Clone[int])
	assert.Equal(t, src, dst)
	assert.False(t, src[0] == dst[0])
	assert.False(t, src[1] == dst[1])
}

func TestFlatMap(t *testing.T) {
	assert.Equal(t, []int{0, 1, 2, 3, 4},
		FlatMap([][]int{{0}, {1, 2}, {3, 4}},
			func(v []int) []int { return v }))
}

// Query the months of some quarters under given fiscal quarter definition.
func TestFlatMap2(t *testing.T) {
	assert.Equal(t, []int{1, 2, 3, 7, 8, 9},
		FlatMap([]string{"Q1", "Q3"},
			func(s string) []int {
				return map[string][]int{
					"Q1": {1, 2, 3},
					"Q2": {4, 5, 6},
					"Q3": {7, 8, 9},
					"Q4": {10, 11, 12},
				}[s]
			}),
	)
}

// make a funcs.FromMap? or gmap.ToFunc
func fromMap[K comparable, V any](m map[K]V) func(K) V {
	return func(x K) V {
		return m[x]
	}
}

// Query the chessboard squares that a knight can reach in 1 step, starting from given squares.
func TestFlatMap4(t *testing.T) {
	knightReach := map[string][]string{
		"a1": {"b3", "c2"},
		"a2": {"b4", "c1", "c3"},
		"a3": {"b1", "b5", "c2", "c4"},
	}
	start := []string{"a1", "a2"}
	next := []string{"b3", "c2", "b4", "c1", "c3"}
	assert.Equal(t, next, FlatMap(start, fromMap(knightReach)))
}

// Query the score you can reach after doing (+1) or (-1) for 1,2,... times, starting from score = 0.
func TestFlatMap5(t *testing.T) {
	start := []int{0}
	next1 := []int{-1, 1}
	next2 := []int{-2, 0, 0, 2}
	move := func(x int) []int { return []int{x - 1, x + 1} }
	assert.Equal(t, next1, FlatMap(start, move))
	assert.Equal(t, next2, FlatMap(FlatMap(start, move), move))
}

// Query the score you can reach after doing (+1) or (-1) for 1,2,... times, starting from score = 0.
func TestFlatMap5s(t *testing.T) {
	move := func(x int) []int { return []int{x - 1, x + 1} }
	assert.Equal(t, []int{-1, 1}, FlatMap([]int{0}, move))
	assert.Equal(t, []int{-2, 0, 0, 2}, FlatMap([]int{-1, 1}, move))
}

// Build the description of one's parents, grandparents, grand-grandparents, ...
func TestFlatMap6(t *testing.T) {
	parents := func(i string) []string { return []string{i + " 's mom", i + " 's dad"} }
	assert.Equal(t, []string{"L 's mom 's mom", "L 's mom 's dad", "L 's dad 's mom", "L 's dad 's dad"},
		FlatMap([]string{"L 's mom", "L 's dad"}, parents))
}

func TestFlatten(t *testing.T) {
	assert.Equal(t, []int{0, 1, 2, 3, 4}, Flatten([][]int{{0}, {1, 2}, {3, 4}}))
}

func TestAny(t *testing.T) {
	{
		sequence := []int{1, 2, 3}
		predicate := func(x int) bool { return x > 2 }
		assert.True(t, Any(sequence, predicate))
	}
	{
		sequence := []int{1, 2, 3}
		predicate := func(x int) bool { return x > 3 }
		assert.False(t, Any(sequence, predicate))
	}
}

func TestAll(t *testing.T) {
	{
		sequence := []int{1, 2, 3}
		predicate := func(x int) bool { return x > 0 }
		assert.True(t, All(sequence, predicate))

	}
	{
		sequence := []int{1, 2, 3}
		predicate := func(x int) bool { return x > 1 }
		assert.False(t, All(sequence, predicate))
	}
}

func TestFirst(t *testing.T) {
	assert.Equal(t, 4, First([]int{4, 3, 1, 4}).Value())
	assert.False(t, First([]int{}).IsOK())
}

func TestLast(t *testing.T) {
	assert.Equal(t, 4, Last([]int{4, 3, 1, 4}).Value())
	assert.False(t, Last([]int{}).IsOK())
}

func TestUnion(t *testing.T) {
	assert.Equal(t, []int{1, 2, 3, 4, 5}, Union([]int{1, 2, 3}, []int{3, 4, 5}))
	assert.Equal(t, []int{1, 2, 3}, Union([]int{1, 2, 3}, []int{}))
	assert.Equal(t, []int{3, 4, 5}, Union([]int{}, []int{3, 4, 5}))

	// Test duplicate elems.
	assert.Equal(t, []int{1, 2, 3, 4, 5}, Union([]int{1, 1, 2, 3}, []int{1, 3, 1, 4, 5}))
	assert.Equal(t, []int{1, 2, 3, 4, 5}, Union([]int{1, 1, 2, 3, 2, 4}, []int{1, 3, 1, 4, 5}))

	// Test multiple slices.
	assert.Equal(t, []int{}, Union[[]int]())
	assert.Equal(t, []int{1, 2}, Union([]int{1, 2, 1}))
	assert.Equal(t, []int{1, 2, 3, 4, 5}, Union([]int{1, 2}, []int{2}, []int{1, 3}, []int{4, 4, 4}, []int{5}))
	assert.Equal(t, []int{}, Union([]int{}, []int{}))
}

func TestDiff(t *testing.T) {
	assert.Equal(t, []int{1, 2}, Diff([]int{1, 2, 3}, []int{3, 4, 5}))
	assert.Equal(t, []int{1, 2, 3}, Diff([]int{1, 2, 3}, []int{4, 5, 6}))
	assert.Equal(t, []int{}, Diff([]int{1, 2, 3}, []int{1, 2, 3}))

	// Test duplicate elems.
	assert.Equal(t, []int{}, Diff([]int{1, 2, 3, 2, 3}, []int{1, 2, 3}))
	assert.Equal(t, []int{1, 2}, Diff([]int{1, 2, 2, 3}, []int{3, 4, 5}))

	// Test multiple slices.
	assert.Equal(t, []int{1, 2}, Diff([]int{1, 2, 1}))
	assert.Equal(t, []int{3}, Diff([]int{1, 2, 3, 3}, []int{1}, []int{2}))
	assert.Equal(t, []int{3}, Diff([]int{1, 2, 3, 3}, []int{1}, []int{2}))
	assert.Equal(t, []int{}, Diff([]int{}, []int{}))
}

func TestIntersect(t *testing.T) {
	assert.Equal(t, []int{2, 3}, Intersect([]int{1, 2, 3}, []int{2, 3, 4}))
	assert.Equal(t, []int{}, Intersect([]int{1, 2, 3}, []int{4, 5, 6}))
	assert.Equal(t, []int{1, 2, 3}, Intersect([]int{1, 2, 3}, []int{1, 2, 3}))

	// Test duplicate elems.
	assert.Equal(t, []int{1, 2, 3}, Intersect([]int{1, 2, 2, 3}, []int{1, 2, 3}))
	assert.Equal(t, []int{1, 2, 3}, Intersect([]int{1, 2, 2, 3}, []int{1, 2, 3, 3}))

	// Test multiple slices.
	assert.Equal(t, []int{1, 2}, Intersect([]int{1, 2, 1}))
	assert.Equal(t, []int{}, Intersect([]int{1, 2, 2}, []int{5}, []int{1, 3}, []int{4, 4, 4}, []int{}))
	assert.Equal(t, []int{2}, Intersect([]int{1, 2, 2}, []int{5, 2}, []int{1, 2, 3}))
	assert.Equal(t, []int{}, Intersect([]int{}, []int{}))
	assert.Equal(t, []int{}, Intersect([]int{}))
	assert.Equal(t, []int{}, Intersect[[]int]())
	assert.Equal(t, []int{1, 2}, Intersect([]int{1, 2, 2, 3}, []int{1, 1, 2, 3, 5, 5}, []int{1, 2, 4}))
}

func TestReverse(t *testing.T) {
	{
		s := []int{1, 2, 3, 4}
		Reverse(s)
		assert.Equal(t, []int{4, 3, 2, 1}, s)
	}

	// Test any type.
	{
		s := []any{1, 2, 3, 4}
		Reverse(s)
		assert.Equal(t, []any{4, 3, 2, 1}, s)
	}
}

func TestReverseClone(t *testing.T) {
	{
		s := []int{1, 2, 3, 4}
		assert.Equal(t, []int{4, 3, 2, 1}, ReverseClone(s))
		assert.Equal(t, []int{1, 2, 3, 4}, s)
	}

	// Test any type.
	{
		s := []any{1, 2, 3, 4}
		assert.Equal(t, []any{4, 3, 2, 1}, ReverseClone(s))
		assert.Equal(t, []any{1, 2, 3, 4}, s)
	}
}

func TestSort(t *testing.T) {
	{
		s := []int{1, 3, 2, 4}
		Sort(s)
		assert.Equal(t, []int{1, 2, 3, 4}, s)
	}
}

func TestSortClone(t *testing.T) {
	{
		s := []int{1, 3, 2, 4}
		assert.Equal(t, []int{1, 2, 3, 4}, SortClone(s))
		assert.Equal(t, []int{1, 3, 2, 4}, s)
	}
}

func TestSortBy(t *testing.T) {
	{
		s := []int{1, 3, 2, 4}
		SortBy(s, gvalue.Less[int])
		assert.Equal(t, []int{1, 2, 3, 4}, s)
	}

	// Test any type.
	{
		s := []any{1, 3, 2, 4}
		SortBy(s, func(a, b any) bool { return a.(int) < b.(int) })
		assert.Equal(t, []any{1, 2, 3, 4}, s)
	}
}

func TestSortCloneBy(t *testing.T) {
	{
		s := []int{1, 3, 2, 4}
		assert.Equal(t, []int{1, 2, 3, 4}, SortCloneBy(s, gvalue.Less[int]))
		assert.Equal(t, []int{1, 3, 2, 4}, s)
	}

	// Test any type.
	{
		s := []any{1, 3, 2, 4}
		assert.Equal(t, []any{1, 2, 3, 4}, SortCloneBy(s, func(a, b any) bool { return a.(int) < b.(int) }))
		assert.Equal(t, []any{1, 3, 2, 4}, s)
	}
}

func TestStableSortBy(t *testing.T) {
	{
		s := []int{1, 3, 2, 4}
		StableSortBy(s, gvalue.Less[int])
		assert.Equal(t, []int{1, 2, 3, 4}, s)
	}

	// Test any type.
	{
		s := []any{1, 3, 2, 4}
		StableSortBy(s, func(a, b any) bool { return a.(int) < b.(int) })
		assert.Equal(t, []any{1, 2, 3, 4}, s)
	}
}

func TestTypeAssert(t *testing.T) {
	assert.Equal(t, []int{1, 2, 3, 4}, TypeAssert[int, any]([]any{1, 2, 3, 4}))
	assert.Equal(t, []any{1, 2, 3, 4}, TypeAssert[any, int]([]int{1, 2, 3, 4}))

	// Omit original type.
	assert.Equal(t, []int{1, 2, 3, 4}, TypeAssert[int]([]interface{}{1, 2, 3, 4}))
	assert.Equal(t, []any{1, 2, 3, 4}, TypeAssert[any]([]int{1, 2, 3, 4}))

	assert.Panic(t, func() {
		TypeAssert[float64]([]any{1, 2, 3, 4})
	})
}

func TestEqual(t *testing.T) {
	assert.True(t, Equal[int](nil, nil))
	assert.True(t, Equal([]int{}, []int{}))
	assert.True(t, Equal([]int{}, nil))
	assert.True(t, Equal(nil, []int{}))
	assert.True(t, Equal([]int{1, 2, 3}, []int{1, 2, 3}))
	assert.False(t, Equal([]int{1, 2, 3}, []int{1, 2, 3, 4}))
	assert.False(t, Equal([]int{1, 2, 3}, []int{1, 2, 4}))
}

func TestEqualBy(t *testing.T) {
	eq := gvalue.Equal[int]
	assert.True(t, EqualBy(nil, nil, eq))
	assert.True(t, EqualBy([]int{}, []int{}, eq))
	assert.True(t, EqualBy([]int{}, nil, eq))
	assert.True(t, EqualBy(nil, []int{}, eq))
	assert.True(t, EqualBy([]int{1, 2, 3}, []int{1, 2, 3}, eq))
	assert.False(t, EqualBy([]int{1, 2, 3}, []int{1, 2, 3, 4}, eq))
	assert.False(t, EqualBy([]int{1, 2, 3}, []int{1, 2, 4}, eq))

	// Test any type.
	eqAny := func(a, b any) bool { return a == b }
	assert.True(t, EqualBy(nil, nil, eqAny))
	assert.True(t, EqualBy([]any{}, []any{}, eqAny))
	assert.True(t, EqualBy([]any{}, nil, eqAny))
	assert.True(t, EqualBy(nil, []any{}, eqAny))
	assert.True(t, EqualBy([]any{1, 2, 3}, []any{1, 2, 3}, eqAny))
	assert.False(t, EqualBy([]any{1, 2, 3}, []any{1, 2, 3, 4}, eqAny))
	assert.False(t, EqualBy([]any{1, 2, 3}, []any{1, 2, 4}, eqAny))
}

func TestToMapValues(t *testing.T) {
	type Foo struct {
		ID int
	}
	mapper := func(f Foo) int { return f.ID }
	assert.Equal(t, map[int]Foo{}, ToMapValues([]Foo{}, mapper))
	assert.Equal(t, map[int]Foo{}, ToMapValues(nil, mapper))
	assert.Equal(t, map[int]Foo{1: {1}, 2: {2}, 3: {3}}, ToMapValues([]Foo{{1}, {2}, {1}, {3}}, mapper))
}

func TestToMap(t *testing.T) {
	type Foo struct {
		ID   int
		Name string
	}
	mapper := func(f Foo) (int, string) { return f.ID, f.Name }
	assert.Equal(t, map[int]string{}, ToMap([]Foo{}, mapper))
	assert.Equal(t, map[int]string{}, ToMap(nil, mapper))
	assert.Equal(t,
		map[int]string{1: "one", 2: "two", 3: "three"},
		ToMap([]Foo{{1, "one"}, {2, "two"}, {3, "three"}}, mapper))
}

func TestToSet(t *testing.T) {
	assert.Equal(t, map[int]bool{}, ToSet([]int{}))
	assert.Equal(t, map[int]bool{1: true, 2: true, 3: true}, ToSet([]int{1, 2, 2, 3}))
	assert.Equal(t,
		map[string]bool{"a": true, "b": true},
		ToSet([]string{"a", "b", "a", "a", "b"}))
}

func TestDivide(t *testing.T) {
	{
		s := []int{0, 1, 2, 3, 4}
		chunks := Divide(s, 2)
		assert.Equal(t, [][]int{{0, 1, 2}, {3, 4}}, chunks)
		chunks[1][1] = 9 // Modify original slice
		assert.Equal(t, []int{0, 1, 2, 3, 9}, s)
		assert.Equal(t, [][]int{{0, 1, 2}, {3, 9}}, chunks)
	}
	{
		s := []int{0, 1, 2, 3, 4, 5, 6}
		chunks := Divide(s, 3)
		assert.Equal(t, [][]int{{0, 1, 2}, {3, 4}, {5, 6}}, chunks)
	}
}

func TestDivideClone(t *testing.T) {
	{
		s := []int{0, 1, 2, 3, 4}
		chunks := DivideClone(s, 2)
		assert.Equal(t, [][]int{{0, 1, 2}, {3, 4}}, chunks)
		chunks[1][1] = 9 // Modify original slice
		assert.Equal(t, []int{0, 1, 2, 3, 4}, s)
		assert.Equal(t, [][]int{{0, 1, 2}, {3, 9}}, chunks)
	}
	{
		s := []int{0, 1, 2, 3, 4, 5, 6}
		chunks := DivideClone(s, 3)
		assert.Equal(t, [][]int{{0, 1, 2}, {3, 4}, {5, 6}}, chunks)
	}
}

func TestFind(t *testing.T) {
	assert.Equal(t, goption.OK(1), Find([]int{0, 1, 2, 3, 4}, func(v int) bool { return v > 0 }))
	assert.Equal(t, goption.Nil[int](), Find([]int{0, 1, 2, 3, 4}, func(v int) bool { return v < 0 }))
	assert.Equal(t, goption.OK(0), Find([]int{0, 1, 2, 3, 4}, func(v int) bool { return v <= 0 }))
	assert.Equal(t, goption.Nil[int](), Find([]int{}, func(v int) bool { return v > 0 }))
	assert.Equal(t, goption.Nil[int](), Find([]int{0, 1, 2, 3, 4}, func(v int) bool { return v < 0 }))
}

func TestFindRev(t *testing.T) {
	assert.Equal(t, goption.OK(4), FindRev([]int{0, 1, 2, 3, 4}, func(v int) bool { return v > 0 }))
	assert.Equal(t, goption.Nil[int](), FindRev([]int{0, 1, 2, 3, 4}, func(v int) bool { return v < 0 }))
	assert.Equal(t, goption.OK(0), FindRev([]int{0, 1, 2, 3, 4}, func(v int) bool { return v <= 0 }))
	assert.Equal(t, goption.Nil[int](), FindRev([]int{}, func(v int) bool { return v > 0 }))
	assert.Equal(t, goption.Nil[int](), FindRev([]int{0, 1, 2, 3, 4}, func(v int) bool { return v < 0 }))
}

func reflectMap(s any, f func(i any) any) any {
	src := reflect.ValueOf(s)
	var dst *reflect.Value
	for i := 0; i < src.Len(); i++ {
		v := reflect.ValueOf(f(src.Index(i).Interface()))
		if dst == nil {
			dst = gptr.Of(reflect.MakeSlice(reflect.SliceOf(v.Type()), src.Len(), src.Len()))
		}
		dst.Index(i).Set(v)
	}
	return dst.Interface()
}

func TestReflectMap(t *testing.T) {
	assert.Equal(t,
		[]string{"1", "2", "3"},
		reflectMap([]int{1, 2, 3}, func(i any) any { return strconv.Itoa(i.(int)) }).([]string))
}

func TestPtrOf(t *testing.T) {
	{
		v1, v2, v3 := 1, 2, 3
		assert.Equal(t, []*int{&v1, &v2, &v3}, PtrOf([]int{1, 2, 3}))
	}

	// Test modifying pointer.
	{
		v1, v2, v3 := 1, 2, 3
		ptrs := PtrOf([]int{v1, v2, v3})
		assert.False(t, ptrs[0] == &v1)
		assert.False(t, ptrs[1] == &v2)
		assert.False(t, ptrs[2] == &v3)
		*ptrs[0] = 4
		assert.Equal(t, 1, v1)
		assert.Equal(t, 4, *ptrs[0])
	}

}
func TestIndirect(t *testing.T) {
	s1 := []*int{nil, nil, nil, gptr.Of(102), gptr.Of(103), gptr.Of(104)}
	s2 := Clone(s1)
	assert.Equal(t, []int{102, 103, 104}, Indirect(s1))
	assert.Equal(t, s2, s1)
}

func TestIndirectOr(t *testing.T) {
	v1, v2, v3 := 1, 2, 3

	assert.Equal(t, []int{1, 2, -1, 3, -1}, IndirectOr([]*int{&v1, &v2, nil, &v3, nil}, -1))
}

func TestShuffle(t *testing.T) {
	{
		expect := []int{1, 2, 3, 4, 5, 6}
		actual := Clone(expect)
		for {
			Shuffle(actual)
			if !Equal(expect, actual) {
				break
			}
		}
		Sort(actual)
		assert.Equal(t, expect, actual)
		actual[0] = 9
		assert.NotEqual(t, expect, actual)
	}
}

func TestShuffleClone(t *testing.T) {
	{
		expect := []int{1, 2, 3, 4, 5, 6}
		var actual []int
		for {
			actual = ShuffleClone(expect)
			if !Equal(expect, actual) {
				break
			}
		}
		Sort(actual)
		assert.Equal(t, expect, actual)
		actual[0] = 9
		assert.NotEqual(t, expect, actual)
	}
}

func TestGet(t *testing.T) {
	assert.Equal(t, goption.OK(0), Get([]int{0, 1, 2, 3, 4}, 0))
	assert.Equal(t, goption.OK(1), Get([]int{0, 1, 2, 3, 4}, 1))
	assert.Equal(t, goption.OK(4), Get([]int{0, 1, 2, 3, 4}, 4))
	assert.Equal(t, goption.Nil[int](), Get([]int{0, 1, 2, 3, 4}, 5))
	assert.Equal(t, goption.Nil[int](), Get([]int{0, 1, 2, 3, 4}, 500))
	assert.Equal(t, goption.OK(4), Get([]int{0, 1, 2, 3, 4}, -1))
	assert.Equal(t, goption.OK(3), Get([]int{0, 1, 2, 3, 4}, -2))
	assert.Equal(t, goption.OK(0), Get([]int{0, 1, 2, 3, 4}, -5))
	assert.Equal(t, goption.Nil[int](), Get([]int{0, 1, 2, 3, 4}, -6))
	assert.Equal(t, goption.Nil[int](), Get([]int{0, 1, 2, 3, 4}, -500))

	// Test integer index.
	assert.Equal(t, goption.OK(1), Get([]int{0, 1, 2, 3, 4}, int8(1)))
	assert.Equal(t, goption.OK(4), Get([]int{0, 1, 2, 3, 4}, int8(-1)))
	assert.Equal(t, goption.OK(1), Get([]int{0, 1, 2, 3, 4}, uint8(1)))
	assert.Equal(t, goption.OK(1), Get([]int{0, 1, 2, 3, 4}, int64(1)))
	assert.Equal(t, goption.OK(4), Get([]int{0, 1, 2, 3, 4}, int64(-1)))
	assert.Equal(t, goption.OK(1), Get([]int{0, 1, 2, 3, 4}, uint64(1)))
}

func TestIndex(t *testing.T) {
	// nil or empty.
	assert.Equal(t, goption.Nil[int](), Index([]string{}, "0"))
	assert.Equal(t, goption.Nil[int](), Index([]int{}, 0))
	assert.Equal(t, goption.Nil[int](), Index(nil, "0"))
	assert.Equal(t, goption.Nil[int](), Index(nil, 0))

	// Smoke cases.
	assert.Equal(t, goption.Nil[int](), Index([]int{0, 1, 2, 3, 4}, -1))
	assert.Equal(t, goption.OK(0), Index([]int{0, 1, 2, 3, 4}, 0))
	assert.Equal(t, goption.OK(1), Index([]int{0, 1, 2, 3, 4}, 1))
	assert.Equal(t, goption.OK(2), Index([]int{0, 1, 2, 3, 4}, 2))
	assert.Equal(t, goption.OK(3), Index([]int{0, 1, 2, 3, 4}, 3))
	assert.Equal(t, goption.OK(4), Index([]int{0, 1, 2, 3, 4}, 4))
	assert.Equal(t, goption.Nil[int](), Index([]int{0, 1, 2, 3, 4}, 5))

	// Duplicate elements.
	assert.Equal(t, goption.OK(0), Index([]int{0, 1, 2, 2, 1}, 0))
	assert.Equal(t, goption.OK(1), Index([]int{0, 1, 2, 2, 1}, 1))
	assert.Equal(t, goption.OK(2), Index([]int{0, 1, 2, 2, 1}, 2))
	assert.Equal(t, goption.Nil[int](), Index([]int{0, 1, 2, 2, 1}, 3))

	assert.Equal(t, goption.OK(1), Index([]string{"a", "b", "b", "d"}, "b"))
	assert.Equal(t, goption.OK(2), Index([]string{"a", "c", "b", "d"}, "b"))
	assert.Equal(t, goption.Nil[int](), Index([]string{"a", "b", "c", "d"}, "e"))
}

func TestIndexRev(t *testing.T) {
	// nil or empty.
	assert.Equal(t, goption.Nil[int](), IndexRev([]string{}, "0"))
	assert.Equal(t, goption.Nil[int](), IndexRev([]int{}, 0))
	assert.Equal(t, goption.Nil[int](), IndexRev(nil, "0"))
	assert.Equal(t, goption.Nil[int](), IndexRev(nil, 0))

	// Smoke cases.
	assert.Equal(t, goption.Nil[int](), IndexRev([]int{0, 1, 2, 3, 4}, -1))
	assert.Equal(t, goption.OK(0), IndexRev([]int{0, 1, 2, 3, 4}, 0))
	assert.Equal(t, goption.OK(1), IndexRev([]int{0, 1, 2, 3, 4}, 1))
	assert.Equal(t, goption.OK(2), IndexRev([]int{0, 1, 2, 3, 4}, 2))
	assert.Equal(t, goption.OK(3), IndexRev([]int{0, 1, 2, 3, 4}, 3))
	assert.Equal(t, goption.OK(4), IndexRev([]int{0, 1, 2, 3, 4}, 4))
	assert.Equal(t, goption.Nil[int](), IndexRev([]int{0, 1, 2, 3, 4}, 5))

	// Duplicate elements.
	assert.Equal(t, goption.OK(0), IndexRev([]int{0, 1, 2, 2, 1}, 0))
	assert.Equal(t, goption.OK(4), IndexRev([]int{0, 1, 2, 2, 1}, 1))
	assert.Equal(t, goption.OK(3), IndexRev([]int{0, 1, 2, 2, 1}, 2))
	assert.Equal(t, goption.Nil[int](), IndexRev([]int{0, 1, 2, 2, 1}, 3))

	assert.Equal(t, goption.OK(2), IndexRev([]string{"a", "b", "b", "d"}, "b"))
	assert.Equal(t, goption.OK(2), IndexRev([]string{"a", "c", "b", "d"}, "b"))
	assert.Equal(t, goption.Nil[int](), IndexRev([]string{"a", "b", "c", "d"}, "e"))
}

func TestIndexBy(t *testing.T) {
	odd := func(v string) bool {
		i, _ := strconv.Atoi(v)
		return i%2 == 1
	}
	assert.Equal(t, goption.OK(1), IndexBy([]string{"0", "1", "2", "3", "4"}, odd))
	assert.Equal(t, goption.OK(3), IndexBy([]string{"0", "2", "2", "3", "4"}, odd))
	assert.Equal(t, goption.Nil[int](), IndexBy([]string{"0", "2", "4"}, odd))

	// Test non-comparable.
	oddAny := func(v any) bool {
		i, _ := strconv.Atoi(v.(string))
		return i%2 == 1
	}
	assert.Equal(t, goption.OK(1), IndexBy([]any{"0", "1", "2", "3", "4"}, oddAny))
}

func TestIndexRevBy(t *testing.T) {
	odd := func(v string) bool {
		i, _ := strconv.Atoi(v)
		return i%2 == 1
	}
	assert.Equal(t, goption.OK(3), IndexRevBy([]string{"0", "1", "2", "3", "4"}, odd))
	assert.Equal(t, goption.OK(3), IndexRevBy([]string{"0", "2", "2", "3", "4"}, odd))
	assert.Equal(t, goption.Nil[int](), IndexRevBy([]string{"0", "2", "4"}, odd))

	// Test non-comparable.
	oddAny := func(v any) bool {
		i, _ := strconv.Atoi(v.(string))
		return i%2 == 1
	}
	assert.Equal(t, goption.OK(3), IndexRevBy([]any{"0", "1", "2", "3", "4"}, oddAny))
}

func TestTake(t *testing.T) {
	assert.Equal(t, []int{}, Take([]int{1, 2, 3, 4, 5}, 0))
	assert.Equal(t, []int{1, 2, 3}, Take([]int{1, 2, 3, 4, 5}, 3))
	assert.Equal(t, []int{1, 2, 3, 4, 5}, Take([]int{1, 2, 3, 4, 5}, 5))
	assert.Equal(t, []int{1, 2, 3, 4, 5}, Take([]int{1, 2, 3, 4, 5}, 10))
	assert.Equal(t, []int{5}, Take([]int{1, 2, 3, 4, 5}, -1))
	assert.Equal(t, []int{1, 2, 3, 4, 5}, Take([]int{1, 2, 3, 4, 5}, -5))
	assert.Equal(t, []int{1, 2, 3, 4, 5}, Take([]int{1, 2, 3, 4, 5}, -10))
	{
		s1 := []int{1, 2, 3, 4, 5}
		s2 := Take(s1, 3)
		s2[0] = 5
		assert.Equal(t, []int{5, 2, 3, 4, 5}, s1)
	}
}

func TestTakeClone(t *testing.T) {
	assert.Equal(t, []int{}, TakeClone([]int{1, 2, 3, 4, 5}, 0))
	assert.Equal(t, []int{1, 2, 3}, TakeClone([]int{1, 2, 3, 4, 5}, 3))
	assert.Equal(t, []int{1, 2, 3, 4, 5}, TakeClone([]int{1, 2, 3, 4, 5}, 5))
	assert.Equal(t, []int{1, 2, 3, 4, 5}, TakeClone([]int{1, 2, 3, 4, 5}, 10))
	assert.Equal(t, []int{5}, TakeClone([]int{1, 2, 3, 4, 5}, -1))
	assert.Equal(t, []int{1, 2, 3, 4, 5}, TakeClone([]int{1, 2, 3, 4, 5}, -5))
	assert.Equal(t, []int{1, 2, 3, 4, 5}, TakeClone([]int{1, 2, 3, 4, 5}, -10))
	{
		s1 := []int{1, 2, 3, 4, 5}
		s2 := TakeClone(s1, 3)
		s2[0] = 5
		assert.Equal(t, []int{1, 2, 3, 4, 5}, s1)
	}
}

func TestDrop(t *testing.T) {
	assert.Equal(t, []int{1, 2, 3, 4, 5}, Drop([]int{1, 2, 3, 4, 5}, 0))
	assert.Equal(t, []int{4, 5}, Drop([]int{1, 2, 3, 4, 5}, 3))
	assert.Equal(t, []int{}, Drop([]int{1, 2, 3, 4, 5}, 5))
	assert.Equal(t, []int{}, Drop([]int{1, 2, 3, 4, 5}, 10))
	assert.Panic(t, func() {
		Drop([]int{1, 2, 3, 4, 5}, -1)
	})
	{
		s1 := []int{1, 2, 3, 4, 5}
		s2 := Drop(s1, 3)
		s2[0] = 5
		assert.Equal(t, []int{1, 2, 3, 5, 5}, s1)
	}
}

func TestDropClone(t *testing.T) {
	assert.Equal(t, []int{1, 2, 3, 4, 5}, DropClone([]int{1, 2, 3, 4, 5}, 0))
	assert.Equal(t, []int{4, 5}, DropClone([]int{1, 2, 3, 4, 5}, 3))
	assert.Equal(t, []int{}, DropClone([]int{1, 2, 3, 4, 5}, 5))
	assert.Equal(t, []int{}, DropClone([]int{1, 2, 3, 4, 5}, 10))
	assert.Panic(t, func() {
		DropClone([]int{1, 2, 3, 4, 5}, -1)
	})
	{
		s1 := []int{1, 2, 3, 4, 5}
		s2 := DropClone(s1, 3)
		s2[0] = 5
		assert.Equal(t, []int{1, 2, 3, 4, 5}, s1)
	}
}

func TestSum(t *testing.T) {
	assert.Equal(t, 0, Sum([]int{}))
	assert.Equal(t, 5, Sum([]int{5}))
	assert.Equal(t, 15, Sum([]int{1, 2, 3, 4, 5}))
	assert.Equal(t, 15.0, Sum([]float64{1, 2, 3, 4, 5}))
}

func TestSumBy(t *testing.T) {
	type Foo struct {
		Value int
	}
	getValue := func(foo Foo) int {
		return foo.Value
	}
	assert.Equal(t, 5, SumBy([]Foo{{5}}, getValue))
}

func TestAvg(t *testing.T) {
	assert.Equal(t, 0, Avg([]int{}))
	assert.Equal(t, 5, Avg([]int{5}))
	assert.Equal(t, 5.0, Avg([]float64{5}))
	assert.True(t, Avg([]int{1, 2, 3, 4, 5})-3.0 < 0.0001)
	assert.True(t, Avg([]float64{1, 2, 3, 4, 5})-3.0 < 0.0001)
}

func TestAvgBy(t *testing.T) {
	type Foo struct {
		Value int
	}
	getValue := func(foo Foo) int {
		return foo.Value
	}
	assert.Equal(t, 5.0, AvgBy([]Foo{{5}}, getValue))
}

func TestLen(t *testing.T) {
	assert.Equal(t, 5, Len([]int{0, 1, 2, 3, 4}))
	assert.Equal(t, 1, Len([]int{0}))
	assert.Equal(t, 0, Len([]int{}))
	assert.Equal(t, 0, Len[int](nil))
}

func TestForEach(t *testing.T) {
	{
		s := []int{0, 1, 2, 3, 4}
		clone := []int{}
		ForEach(s, func(v int) {
			clone = append(clone, v)
		})
		assert.Equal(t, s, clone)
	}
}

func TestForEachIndexed(t *testing.T) {
	{
		s := []string{"0", "1", "2", "3", "4"}
		clone := []int{}
		ForEachIndexed(s, func(i int, v string) {
			clone = append(clone, i)
		})
		assert.Equal(t, []int{0, 1, 2, 3, 4}, clone)
	}
}

func TestConcat(t *testing.T) {
	assert.Equal(t, []int{0, 1, 2, 3, 4}, Concat([]int{0}, []int{1, 2}, []int{3, 4}))
}

func TestMerge(t *testing.T) {
	assert.Equal(t, []int{0, 1, 2, 3, 4}, Merge([]int{0}, []int{1, 2}, []int{3, 4}))
}

func TestCompact(t *testing.T) {
	assert.Equal(t, []int{}, Compact([]int(nil)))
	assert.Equal(t, []int{}, Compact([]int{}))
	assert.Equal(t, []int{1, 2, 3, 4}, Compact([]int{0, 1, 2, 3, 4}))
	assert.Equal(t, []int{1, 2, 3, -1, 4}, Compact([]int{0, 1, 0, 0, 2, 3, 0, -1, 4}))
	assert.Equal(t, []int{}, Compact([]int{0, 0, 0}))
	assert.Equal(t, []string{"foo", "bar"}, Compact([]string{"", "foo", "", "bar"}))
}

func TestInsertInplace(t *testing.T) {
	// Test empty.
	assert.Equal(t, nil, insertInplace[int](nil, 0))
	assert.Equal(t, []int{1}, insertInplace(nil, 0, 1))
	assert.Equal(t, []int{1}, insertInplace(nil, 100, 1))
	assert.Equal(t, []int{1}, insertInplace(nil, -1, 1))
	assert.Equal(t, []int{1}, insertInplace(nil, -100, 1))
	assert.Equal(t, []int{1}, insertInplace([]int{}, 0, 1))
	assert.Equal(t, []int{1}, insertInplace([]int{}, 100, 1))
	assert.Equal(t, []int{1}, insertInplace([]int{}, -1, 1))
	assert.Equal(t, []int{1}, insertInplace([]int{}, -100, 1))

	for i := -100; i < -5; i++ {
		assert.Equal(t, []int{999, 1, 2, 3, 4}, insertInplace([]int{1, 2, 3, 4}, i, 999))
	}
	assert.Equal(t, []int{999, 1, 2, 3, 4}, insertInplace([]int{1, 2, 3, 4}, -4, 999))
	assert.Equal(t, []int{1, 999, 2, 3, 4}, insertInplace([]int{1, 2, 3, 4}, -3, 999))
	assert.Equal(t, []int{1, 2, 999, 3, 4}, insertInplace([]int{1, 2, 3, 4}, -2, 999))
	assert.Equal(t, []int{1, 2, 3, 999, 4}, insertInplace([]int{1, 2, 3, 4}, -1, 999))
	assert.Equal(t, []int{999, 1, 2, 3, 4}, insertInplace([]int{1, 2, 3, 4}, 0, 999))
	assert.Equal(t, []int{1, 999, 2, 3, 4}, insertInplace([]int{1, 2, 3, 4}, 1, 999))
	assert.Equal(t, []int{1, 2, 999, 3, 4}, insertInplace([]int{1, 2, 3, 4}, 2, 999))
	assert.Equal(t, []int{1, 2, 3, 999, 4}, insertInplace([]int{1, 2, 3, 4}, 3, 999))
	assert.Equal(t, []int{1, 2, 3, 4, 999}, insertInplace([]int{1, 2, 3, 4}, 4, 999))
	for i := 5; i < 100; i++ {
		assert.Equal(t, []int{1, 2, 3, 4, 999}, insertInplace([]int{1, 2, 3, 4}, i, 999))
	}

	// Test reuse
	{
		before := make([]int, 0, 5)
		before = append(before, 1, 2, 3, 4)
		after := insertInplace(before, 1, 999)
		assert.Equal(t, []int{1, 999, 2, 3}, before)
		assert.True(t, len(before) == 4 && cap(before) == 5)
		assert.Equal(t, []int{1, 999, 2, 3, 4}, after)
		assert.True(t, len(after) == 5 && cap(after) == 5)
		assert.Equal(t, fmt.Sprintf("%p", before), fmt.Sprintf("%p", after))
	}

	// Test multiple.
	assert.Equal(t, []int{1, 2, 997, 998, 999, 3, 4},
		insertInplace([]int{1, 2, 3, 4}, 2, 997, 998, 999))
	assert.Equal(t, []int{1, 2, 3, 4}, insertInplace([]int{1, 2, 3, 4}, 0))

	{
		// Test for big cap.
		arr := make([]int, 3, 100)
		arr[0], arr[1], arr[2] = 1, 2, 3
		assert.Equal(t, []int{1, 2, 3, 4, 5}, insertInplace(arr, 5, 4, 5))
		arr = make([]int, 3, 100)
		arr[0], arr[1], arr[2] = 1, 2, 3
		assert.Equal(t, []int{4, 5, 1, 2, 3}, insertInplace(arr, -100, 4, 5))
	}
}

func TestInsert(t *testing.T) {
	// Test empty.
	assert.Equal(t, nil, Insert([]int(nil), 0))
	assert.Equal(t, []int{1}, Insert([]int(nil), 0, 1))
	assert.Equal(t, []int{1}, Insert([]int(nil), 100, 1))
	assert.Equal(t, []int{1}, Insert([]int(nil), -1, 1))
	assert.Equal(t, []int{1}, Insert([]int(nil), -100, 1))
	assert.Equal(t, []int{1}, Insert([]int{}, 0, 1))
	assert.Equal(t, []int{1}, Insert([]int{}, 100, 1))
	assert.Equal(t, []int{1}, Insert([]int{}, -1, 1))
	assert.Equal(t, []int{1}, Insert([]int{}, -100, 1))

	for i := -100; i < -5; i++ {
		assert.Equal(t, []int{999, 1, 2, 3, 4}, Insert([]int{1, 2, 3, 4}, i, 999))
	}
	assert.Equal(t, []int{999, 1, 2, 3, 4}, Insert([]int{1, 2, 3, 4}, -4, 999))
	assert.Equal(t, []int{1, 999, 2, 3, 4}, Insert([]int{1, 2, 3, 4}, -3, 999))
	assert.Equal(t, []int{1, 2, 999, 3, 4}, Insert([]int{1, 2, 3, 4}, -2, 999))
	assert.Equal(t, []int{1, 2, 3, 999, 4}, Insert([]int{1, 2, 3, 4}, -1, 999))
	assert.Equal(t, []int{999, 1, 2, 3, 4}, Insert([]int{1, 2, 3, 4}, 0, 999))
	assert.Equal(t, []int{1, 999, 2, 3, 4}, Insert([]int{1, 2, 3, 4}, 1, 999))
	assert.Equal(t, []int{1, 2, 999, 3, 4}, Insert([]int{1, 2, 3, 4}, 2, 999))
	assert.Equal(t, []int{1, 2, 3, 999, 4}, Insert([]int{1, 2, 3, 4}, 3, 999))
	assert.Equal(t, []int{1, 2, 3, 4, 999}, Insert([]int{1, 2, 3, 4}, 4, 999))
	for i := 5; i < 100; i++ {
		assert.Equal(t, []int{1, 2, 3, 4, 999}, Insert([]int{1, 2, 3, 4}, i, 999))
	}

	// Test reuse
	{
		before := make([]int, 0, 5)
		before = append(before, 1, 2, 3, 4)
		after := Insert(before, 1, 999)
		assert.Equal(t, []int{1, 2, 3, 4}, before) // no change
		assert.True(t, len(before) == 4 && cap(before) == 5)
		assert.Equal(t, []int{1, 999, 2, 3, 4}, after)
		assert.True(t, len(after) == 5 && cap(after) == 5)
		assert.NotEqual(t, fmt.Sprintf("%p", before), fmt.Sprintf("%p", after))
	}

	// Test multiple.
	assert.Equal(t, []int{1, 2, 997, 998, 999, 3, 4},
		Insert([]int{1, 2, 3, 4}, 2, 997, 998, 999))
	assert.Equal(t, []int{1, 2, 3, 4}, Insert([]int{1, 2, 3, 4}, 0))

	// Test integer index.
	assert.Equal(t, []int{1, 2, 3, 999, 4}, Insert([]int{1, 2, 3, 4}, int8(-1), 999))
	assert.Equal(t, []int{1, 999, 2, 3, 4}, Insert([]int{1, 2, 3, 4}, int8(1), 999))
	assert.Equal(t, []int{1, 999, 2, 3, 4}, Insert([]int{1, 2, 3, 4}, uint8(1), 999))
	assert.Equal(t, []int{1, 2, 3, 999, 4}, Insert([]int{1, 2, 3, 4}, int64(-1), 999))
	assert.Equal(t, []int{1, 999, 2, 3, 4}, Insert([]int{1, 2, 3, 4}, int64(1), 999))
	assert.Equal(t, []int{1, 999, 2, 3, 4}, Insert([]int{1, 2, 3, 4}, uint64(1), 999))
}

func TestSlice(t *testing.T) {
	intTbl := []struct {
		in       []int
		start    int
		end      int
		expected []int
	}{
		{in: []int{1, 2, 3, 4, 5}, start: -1000, end: 0, expected: []int{1, 2, 3, 4, 5}},
		{in: []int{1, 2, 3, 4, 5}, start: 0, end: 0, expected: []int{}},
		{in: []int{1, 2, 3, 4, 5}, start: 1, end: 0, expected: []int{}},
		{in: []int{1, 2, 3, 4, 5}, start: 0, end: 1, expected: []int{1}},
		{in: []int{1, 2, 3, 4, 5}, start: 0, end: 2, expected: []int{1, 2}},
		{in: []int{1, 2, 3, 4, 5}, start: 0, end: 3, expected: []int{1, 2, 3}},
		{in: []int{1, 2, 3, 4, 5}, start: 0, end: 4, expected: []int{1, 2, 3, 4}},
		{in: []int{1, 2, 3, 4, 5}, start: 1, end: 0, expected: []int{}},
		{in: []int{1, 2, 3, 4, 5}, start: 1, end: 1, expected: []int{}},
		{in: []int{1, 2, 3, 4, 5}, start: 2, end: 3, expected: []int{3}},
		{in: []int{1, 2, 3, 4, 5}, start: 0, end: 100, expected: []int{1, 2, 3, 4, 5}},
		{in: []int{1, 2, 3, 4, 5}, start: 100, end: 99, expected: []int{}},
		{in: []int{1, 2, 3, 4, 5}, start: 100, end: 100, expected: []int{}},
		{in: []int{1, 2, 3, 4, 5}, start: 100, end: 0, expected: []int{}},
		{in: []int{1, 2, 3, 4, 5}, start: -1, end: 0, expected: []int{5}},
		{in: []int{1, 2, 3, 4, 5}, start: -3, end: 0, expected: []int{3, 4, 5}},
		{in: []int{1, 2, 3, 4, 5}, start: -3, end: 4, expected: []int{3, 4}},
		{in: []int{1, 2, 3, 4, 5}, start: -3, end: -5, expected: []int{}},
		{in: []int{1, 2, 3, 4, 5}, start: -100, end: -5, expected: []int{}},
		{in: []int{1, 2, 3, 4, 5}, start: -100, end: 5, expected: []int{1, 2, 3, 4, 5}},
		{in: []int{1, 2, 3, 4, 5}, start: -3, end: -7, expected: []int{}},
		{in: []int{1, 2, 3, 4, 5}, start: -3, end: -1, expected: []int{3, 4}},
		{in: []int{1, 2, 3, 4, 5}, start: -3, end: 4, expected: []int{3, 4}},
		{in: []int{1, 2, 3, 4, 5}, start: 2, end: -1, expected: []int{3, 4}},
		{in: []int{1, 2, 3, 4, 5}, start: 2, end: 4, expected: []int{3, 4}},
		{in: []int{1, 2, 3, 4, 5}, start: -3, end: -1, expected: []int{3, 4}},
	}
	stringTbl := []struct {
		in       []string
		start    int
		end      int
		expected []string
	}{
		{in: []string{"byte", "dance", "is", "best", "company"}, start: 0, end: 0, expected: []string{}},
		{in: []string{"byte", "dance", "is", "best", "company"}, start: 0, end: 1, expected: []string{"byte"}},
	}
	t.Run("SliceRange", func(t *testing.T) {
		t.Parallel()
		for _, col := range intTbl {
			out := Slice(col.in, col.start, col.end)
			assert.Equal(t, col.expected, out)
		}
		for _, col := range stringTbl {
			out := Slice(col.in, col.start, col.end)
			assert.Equal(t, col.expected, out)
		}
		// TestModify
		s := []int{1, 2, 3, 4, 5}
		out := Slice(s, 0, 5)
		out[0], out[4] = -1, -5
		assert.Equal(t, []int{-1, 2, 3, 4, -5}, s)
	})
	t.Run("SliceRangeClone", func(t *testing.T) {
		t.Parallel()
		for _, col := range intTbl {
			out := SliceClone(col.in, col.start, col.end)
			assert.Equal(t, col.expected, out)
			assert.False(t, overlaps(col.in, out))
		}
		for _, col := range stringTbl {
			out := SliceClone(col.in, col.start, col.end)
			assert.Equal(t, col.expected, out)
			assert.False(t, overlaps(col.in, out))
		}

		// TestModifyClone
		s := []int{1, 2, 3, 4, 5}
		out := SliceClone(s, 0, 5)
		out[0], out[4] = -1, -5
		assert.Equal(t, []int{1, 2, 3, 4, 5}, s)

	})

	// Test integer index.
	assert.Equal(t, []int{1, 2}, Slice([]int{1, 2, 3, 4, 5}, int8(0), int8(2)))
	assert.Equal(t, []int{3, 4}, Slice([]int{1, 2, 3, 4, 5}, int8(-3), int8(-1)))
	assert.Equal(t, []int{1, 2}, Slice([]int{1, 2, 3, 4, 5}, uint8(0), uint8(2)))
	assert.Equal(t, []int{1, 2}, Slice([]int{1, 2, 3, 4, 5}, int64(0), int64(2)))
	assert.Equal(t, []int{3, 4}, Slice([]int{1, 2, 3, 4, 5}, int64(-3), int64(-1)))
	assert.Equal(t, []int{1, 2}, Slice([]int{1, 2, 3, 4, 5}, uint64(0), uint64(2)))
}

func TestOf(t *testing.T) {
	assert.Equal(t, []int{}, Of[int]())
	assert.Equal(t, []int{1}, Of(1))
	assert.Equal(t, []int{1, 2, 3}, Of(1, 2, 3))
}

func TestRangeWithStep(t *testing.T) {
	assert.Equal(t, []int{}, RangeWithStep(0, 0, 2))
	assert.Equal(t, []int{}, RangeWithStep(0, -1, 2))
	assert.Equal(t, []int{0}, RangeWithStep(0, 1, 2))
	assert.Equal(t, []int{0, -1, -2, -3, -4}, RangeWithStep(0, -5, -1))
	assert.Equal(t, []int{0, 2, 4}, RangeWithStep(0, 5, 2))
	assert.Equal(t, []int{0, 3}, RangeWithStep(0, 5, 3))

	assert.Equal(t, []float64{0.5, 1, 1.5}, RangeWithStep(0.5, 2.0, 0.5))
}

func TestRange(t *testing.T) {
	assert.Equal(t, []int{}, Range(0, 0))
	assert.Equal(t, []int{}, Range(0, -1))
	assert.Equal(t, []int{0}, Range(0, 1))
	assert.Equal(t, []int{-2, -1}, Range(-2, 0))
	assert.Equal(t, []int{0, 1, 2, 3, 4}, Range(0, 5))
}

func TestRemoveIndex(t *testing.T) {
	// traditional cases
	tbl := []struct {
		input    []int
		index    int
		expected []int
	}{
		{input: []int{}, index: 0, expected: []int{}},
		{input: []int{}, index: 1, expected: []int{}},
		{input: []int{}, index: -1, expected: []int{}},
		{input: []int{}, index: 100, expected: []int{}},
		{input: []int{}, index: -100, expected: []int{}},
		{input: []int{1}, index: 0, expected: []int{}},
		{input: []int{1}, index: 0, expected: []int{}},
		{input: []int{1}, index: 1, expected: []int{1}},
		{input: []int{1, 2}, index: 1, expected: []int{1}},
		{input: []int{1, 2}, index: 2, expected: []int{1, 2}},
		{input: []int{1, 2, 3}, index: 1, expected: []int{1, 3}},
		{input: []int{1, 2, 3, 4, 5}, index: -1, expected: []int{1, 2, 3, 4}},
		{input: []int{1, 2, 3, 4, 5}, index: -2, expected: []int{1, 2, 3, 5}},
		{input: []int{1, 2, 3, 4, 5}, index: -1000, expected: []int{1, 2, 3, 4, 5}},
		{input: []int{1, 2, 3, 4, 5}, index: 1000, expected: []int{1, 2, 3, 4, 5}},
		{input: []int{1, 2, 3, 4, 5}, index: 3, expected: []int{1, 2, 3, 5}},
	}
	for _, col := range tbl {
		output := RemoveIndex(col.input, col.index)
		assert.Equal(t, col.expected, output)
		assert.False(t, overlaps(col.input, output))
	}

	// Different type cases.
	assert.Equal(t, []string{"hello"}, RemoveIndex([]string{"hello", "world"}, int64(1)))
	assert.Equal(t, []string{"hello"}, RemoveIndex([]string{"hello", "world"}, int32(1)))
	assert.Equal(t, []string{"hello"}, RemoveIndex([]string{"hello", "world"}, int8(1)))
	assert.Equal(t, []string{"hello"}, RemoveIndex([]string{"hello", "world"}, int16(1)))
	assert.Equal(t, []rune("acdefghijklmnop"), RemoveIndex([]rune("abcdefghijklmnop"), 1))

	// equivalence
	{
		// RemoveIndex(x, 0) remove at the front of the slice, is equivalent to s[1:]
		arr := []int{1, 2, 3, 4, 5}
		assert.Equal(t, arr[1:len(arr)], RemoveIndex([]int{1, 2, 3, 4, 5}, 0))

		// RemoveIndex(x, 0) remove at the front of the slice, is equivalent to s[1:]
		arr = []int{1, 2, 3, 4, 5}
		assert.Equal(t, arr[1:], RemoveIndex([]int{1, 2, 3, 4, 5}, 0))

		// RemoveIndex(x, len(x) - 1) is equivalent to s[0:len(x)-1]
		arr = []int{1, 2, 3, 4, 5}
		assert.Equal(t, arr[0:len(arr)-1], RemoveIndex([]int{1, 2, 3, 4, 5}, len(arr)-1))
		// RemoveIndex(x, len(x) - 1) is equivalent to s[0:len(x)-1]
		arr = []int{1, 2, 3, 4, 5}
		assert.Equal(t, arr[0:len(arr)-1], RemoveIndex([]int{1, 2, 3, 4, 5}, 4))

		// RemoveIndex(x, -1) is equivalent to RemoveIndex(x, len(x)-1)
		arr = []int{1, 2, 3, 4, 5}
		assert.Equal(t, arr[0:len(arr)-1], RemoveIndex([]int{1, 2, 3, 4, 5}, -1))

		// RemoveIndex(x, len(x)) does nothing.
		arr = []int{1, 2, 3, 4, 5}
		assert.Equal(t, arr, RemoveIndex([]int{1, 2, 3, 4, 5}, len(arr)))
	}

	// Overlap
	{
		tbl := []struct {
			oldArr      []int
			index       int
			expectedArr []int
		}{
			{oldArr: []int{1, 2, 3, 4, 5}, index: -1000, expectedArr: []int{1, 2, 3, 4, 5}},
			{oldArr: []int{1, 2, 3, 4, 5}, index: 0, expectedArr: []int{2, 3, 4, 5}},
			{oldArr: []int{1, 2, 3, 4, 5}, index: 4, expectedArr: []int{1, 2, 3, 4}},
			{oldArr: []int{1, 2, 3, 4, 5}, index: 2, expectedArr: []int{1, 2, 4, 5}},
		}
		for _, col := range tbl {
			getArr := RemoveIndex(col.oldArr, col.index)
			assert.Equal(t, col.expectedArr, getArr)
			// all the slice should be cloned, not overlap.
			assert.False(t, overlaps(getArr, col.oldArr))
		}
	}
}

func TestCount(t *testing.T) {
	assert.Equal(t, Count([]int{}, 0), 0)
	assert.Equal(t, Count([]int{}, 1), 0)

	assert.Equal(t, Count([]int{1, 2, 3, 3, 2, 3}, 0), 0)
	assert.Equal(t, Count([]int{2, 1, 2, 3, 3, 3}, 1), 1)
	assert.Equal(t, Count([]int{2, 1, 2, 3, 3, 3}, 2), 2)
	assert.Equal(t, Count([]int{2, 1, 2, 3, 3, 3}, 3), 3)
	assert.Equal(t, Count([]int{2, 1, 2, 3, 3, 3}, 4), 0)
}

func TestCountBy(t *testing.T) {
	f := func(v int) bool { return v%2 == 0 }
	assert.Equal(t, CountBy([]int{}, f), 0)
	assert.Equal(t, CountBy([]int{}, f), 0)

	assert.Equal(t, CountBy([]int{1, 3, 3, 3}, f), 0)
	assert.Equal(t, CountBy([]int{1, 3, 3, 2, 3}, f), 1)
	assert.Equal(t, CountBy([]int{2, 1, 2, 3, 3, 3}, f), 2)
	assert.Equal(t, CountBy([]int{2, 1, 2, 3, 3, 4, 3}, f), 3)
	assert.Equal(t, CountBy([]int{2, 1, 2, 4, 3, 3, 4, 3}, f), 4)

	// Test non-comparable type.

	type Foo struct{ v int }
	foos := []Foo{{1}, {2}, {3}}
	assert.Equal(t, CountBy(foos, func(v Foo) bool { return v.v%2 == 0 }), 1)
}

func TestCountValues(t *testing.T) {
	assert.Equal(t, CountValues([]int{}), map[int]int{})
	assert.Equal(t, CountValues([]string{"a", "b", "b"}), map[string]int{"a": 1, "b": 2})
	assert.Equal(t, CountValues([]int{0, 1, 2, 0, 1, 1}), map[int]int{0: 2, 1: 3, 2: 1})
}

func TestCountValuesBy(t *testing.T) {
	isEven := func(v int) bool { return v%2 == 0 }
	assert.Equal(t, CountValuesBy([]int{}, isEven), map[bool]int{})
	assert.Equal(t, CountValuesBy([]int{0, 1, 2, 3, 4}, isEven), map[bool]int{true: 3, false: 2})

	// Test non-comparable type.

	type Foo struct{ v int }
	foos := []Foo{{1}, {2}, {3}}
	assert.Equal(t, CountValuesBy(foos, func(v Foo) bool { return v.v%2 == 0 }), map[bool]int{true: 1, false: 2})
}

// overlaps reports whether the memory ranges a[0:len(a)] and b[0:len(b)] overlap.
// https://github.com/golang/go/blob/master/src/slices/slices.go#L466-L479
func overlaps[E any](a, b []E) bool {
	if len(a) == 0 || len(b) == 0 {
		return false
	}
	elemSize := unsafe.Sizeof(a[0])
	if elemSize == 0 {
		return false
	}
	return uintptr(unsafe.Pointer(&a[0])) <= uintptr(unsafe.Pointer(&b[len(b)-1]))+(elemSize-1) &&
		uintptr(unsafe.Pointer(&b[0])) <= uintptr(unsafe.Pointer(&a[len(a)-1]))+(elemSize-1)
}
