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

// Package gslice provides generic operations for slices.
//
// 💡 HINT: We provide similar functionality for different types in different packages.
// For example, [Clone] for copying slice while
// [github.com/bytedance/gg/gmap.Clone] for copying map.
//
//   - Use [github.com/bytedance/gg/gmap] for map operations.
//   - Use [github.com/bytedance/gg/gvalue] for value operations.
//   - Use [github.com/bytedance/gg/gptr] for pointer operations.
//   - …
//
// # Operations
//
// High-order functions:
//
//   - [Map], [FlatMap]
//   - [Filter], [Reject], [FilterMap]
//   - [Reduce], [Fold]
//   - [All], [Any]
//
// CRUD operations:
//
//   - [Insert], [Remove]
//   - [Get], [First], [Last]
//   - [Index], [Find]
//   - [Contains], [ContainsAny], [ContainsAll]
//
// Partition operations:
//
//   - [Slice]
//   - [Take], [Drop]
//   - [Chunk], [Divide]
//   - [Concat], [Flatten]
//   - [Partition]
//
// Math operations:
//
//   - [Max], [Min], [MinMax]
//   - [Sum], [Avg]
//
// Convert to Map:
//
//   - [ToMap], [ToMapValues]
//   - [GroupBy]
//
// Set operations:
//
//   - [Union], [Intersect], [Diff]
//   - [Uniq], [Dup]
//
// Re-order operations:
//
//   - [Sort], [StableSortBy]
//   - [Reverse]
//   - [Shuffle]
//
// Type casting/assertion/conversion:
//
//   - [TypeAssert]
//   - [PtrOf], [Indirect], [IndirectOr]
//
// Predicates:
//
//   - [Equal]
//
// # Negative index
//
// Some of operations (such as [Get], [Insert] and [Slice]) support negative index like Python.
// The diagram is helpful too for remembering how negative index work:
//
//	         ┌───┬───┬───┬───┬───┬───┐
//	         │ g │ o │ l │ a │ n │ g │
//	         └───┴───┴───┴───┴───┴───┘
//	positive:  0   1   2   3   4   5  (0-based index from start)
//	negative: -6  -5  -4  -3  -2  -1  (negative index from end)
package gslice

import (
	"github.com/bytedance/gg/collection/set"
	"github.com/bytedance/gg/collection/tuple"
	"github.com/bytedance/gg/gfunc"
	"github.com/bytedance/gg/goption"
	"github.com/bytedance/gg/gptr"
	"github.com/bytedance/gg/gresult"
	"github.com/bytedance/gg/gvalue"
	"github.com/bytedance/gg/internal/constraints"
	"github.com/bytedance/gg/internal/iter"
	"github.com/bytedance/gg/internal/rtassert"
)

// Map applies function f to each element of slice s with type F.
// Results of f are returned as a newly allocated slice with type T.
//
// 🚀 EXAMPLE:
//
//	Map([]int{1, 2, 3}, strconv.Itoa) ⏩ []string{"1", "2", "3"}
//	Map([]int{}, strconv.Itoa)        ⏩ []string{}
//	Map(nil, strconv.Itoa)            ⏩ []string{}
//
// 💡 HINT:
//
//   - Use [FilterMap] if you also want to ignore some element during mapping.
//   - Use [TryMap] if function f may fail (return (T, error))
func Map[F, T any](s []F, f func(F) T) []T {
	ret := make([]T, 0, len(s))
	for _, v := range s {
		ret = append(ret, f(v))
	}
	return ret
}

// TryMap is a variant of [Map] that allows function f to fail (return error).
//
// 🚀 EXAMPLE:
//
//	TryMap([]string{"1", "2", "3"}, strconv.Atoi) ⏩ gresult.OK([]int{1, 2, 3})
//	TryMap([]string{"1", "2", "a"}, strconv.Atoi) ⏩ gresult.Err("strconv.Atoi: parsing \"a\": invalid syntax")
//	TryMap([]string{}, strconv.Atoi)              ⏩ gresult.OK([]int{})
//
// 💡 HINT: Use [TryFilterMap] if you want to ignore error during mapping.
func TryMap[F, T any](s []F, f func(F) (T, error)) gresult.R[[]T] {
	ret := make([]T, 0, len(s))
	for _, v := range s {
		r, err := f(v)
		if err != nil {
			return gresult.Err[[]T](err)
		}
		ret = append(ret, r)
	}
	return gresult.OK(ret)
}

// Filter applies predicate f to each element of slice s,
// returns those elements that satisfy the predicate f as a newly allocated slice.
//
// 🚀 EXAMPLE:
//
//	Filter([]int{0, 1, 2, 3}, gvalue.IsNotZero[int]) ⏩ []int{1, 2, 3}
//
// 💡 HINT:
//
//   - Use [FilterMap] if you also want to change the element during filtering.
//   - If you need elements that do not satisfy f, use [Reject]
//   - If you need both elements, use [Partition]
func Filter[S ~[]T, T any](s S, f func(T) bool) S {
	ret := make(S, 0, len(s)/2)
	for _, v := range s {
		if f(v) {
			ret = append(ret, v)
		}
	}
	return ret
}

// FilterMap does [Filter] and [Map] at the same time, applies function f to
// each element of slice s. f returns (T, bool):
//
//   - If true ,the return value with type T will added to
//     the result slice []T.
//   - If false, the return value with type T will be dropped.
//
// 🚀 EXAMPLE:
//
//	f := func(i int) (string, bool) { return strconv.Itoa(i), i != 0 }
//	FilterMap([]int{1, 2, 3, 0, 0}, f) ⏩ []string{"1", "2", "3"}
//
// 💡 HINT: Use [TryFilterMap] if function f returns (T, error).
func FilterMap[F, T any](s []F, f func(F) (T, bool)) []T {
	return iter.ToSlice(iter.FilterMap(f, iter.StealSlice(s)))
}

// TryFilterMap is a variant of [FilterMap] that allows function f to fail (return error).
//
// 🚀 EXAMPLE:
//
//	TryFilterMap([]string{"1", "2", "3"}, strconv.Atoi) ⏩ []int{1, 2, 3}
//	TryFilterMap([]string{"1", "2", "a"}, strconv.Atoi) ⏩ []int{1, 2}
func TryFilterMap[F, T any](s []F, f func(F) (T, error)) []T {
	ret := make([]T, 0, len(s)/2)
	for _, v := range s {
		r, err := f(v)
		if err != nil {
			continue // ignored
		}
		ret = append(ret, r)
	}
	return ret
}

// Reject applies predicate f to each element of slice s,
// returns those elements that do not satisfy the predicate f as a newly allocated slice.
//
// 🚀 EXAMPLE:
//
//	Reject([]int{0, 1, 2, 3}, gvalue.IsZero[int]) ⏩ []int{1, 2, 3}
//
// 💡 HINT:
//
//   - If you need elements that satisfy f, use [Filter]
//   - If you need both elements, use [Partition]
func Reject[S ~[]T, T any](s S, f func(T) bool) S {
	ret := make(S, 0, len(s)/2)
	for _, v := range s {
		if !f(v) {
			ret = append(ret, v)
		}
	}
	return ret
}

// Partition applies predicate f to each element of slice s,
// divides elements into 2 parts: satisfy f and do not satisfy f.
//
// 🚀 EXAMPLE:
//
//	Partition([]int{0, 1, 2, 3}, gvalue.IsNotZero[int]) ⏩ []int{1, 2, 3}, []int{0}
//
// 💡 HINT:
//
//   - Use [Filter] or [Reject] if you need only one of the return values
//   - Use [Chunk] or [Divide] if you want to divide elements by index
func Partition[S ~[]T, T any](s S, f func(T) bool) (S, S) {
	var (
		retTrue  = make(S, 0, len(s)/2)
		retFalse = make(S, 0, len(s)/2)
	)
	for _, v := range s {
		if f(v) {
			retTrue = append(retTrue, v)
		} else {
			retFalse = append(retFalse, v)
		}
	}
	return retTrue, retFalse
}

// Reduce is a variant of Fold, use possible first element of slice as the
// initial value of accumulation.
// If the given slice is empty, goption.Nil[T]() is returned.
//
// 🚀 EXAMPLE:
//
//	Reduce([]int{0, 1, 2, 3}, gvalue.Add[int]) ⏩ goption.OK(6)
//	Reduce([]int{}, gvalue.Add[int])           ⏩ goption.Nil[int]()
//
// 💡 HINT: Calculate the maximum value is only for example, you can directly use
// function [Max].
func Reduce[T any](s []T, f func(T, T) T) goption.O[T] {
	return iter.Reduce(f, iter.StealSlice(s))
}

// Fold applies function f cumulatively to each element of slice s,
// so as to fold the slice to a single value.
// An init element is needed as the initial value of accumulation.
// If the given slice is empty, the init element is returned.
//
// 🚀 EXAMPLE:
//
//	s := []int{0, 1, 2, 3}
//	Fold(s, gvalue.Add[int], 4)       ⏩ 10
//	Fold(s, gvalue.Add[int], 2)       ⏩ 8
//	Fold([]int{}, gvalue.Add[int], 1) ⏩ 1
func Fold[T1, T2 any](s []T1, f func(T2, T1) T2, init T2) T2 {
	return iter.Fold(f, init, iter.StealSlice(s))
}

// Contains returns whether the element occur in slice.
//
// 🚀 EXAMPLE:
//
//	Contains([]int{0, 1, 2, 3, 4}, 1) ⏩ true
//	Contains([]int{0, 1, 2, 3, 4}, 5) ⏩ false
//	Contains([]int{}, 5)              ⏩ false
//
// 💡 HINT:
//
//   - Use [ContainsAll], [ContainsAny] if you have multiple values to query
//   - Use [Index] if you also want to know index of the found value
//   - Use [Any] or [Find] if type of v is non-comparable
func Contains[T comparable](s []T, v T) bool {
	for _, vv := range s {
		if v == vv {
			return true
		}
	}
	return false
}

// ContainsBy returns whether any element of the slice s satisfies the predicate f.
//
// 🚀 EXAMPLE:
//
//	ContainsBy([]int{1, 2, 3}, func(x int) bool { return x > 2 })) ⏩ true
//	ContainsBy([]int{1, 2, 3}, func(x int) bool { return x > 3 })) ⏩ false
//	ContainsBy([]int{1, 2, 3}, gvalue.IsZero[int]))                ⏩ false
//	ContainsBy([]int{1, 2, 3}, gvalue.IsNotZero[int]))             ⏩ true
//
// 💡 AKA: Any
func ContainsBy[T any](s []T, f func(T) bool) bool {
	return Any(s, f)
}

// ContainsAny returns whether any of given elements occur in slice.
//
// 🚀 EXAMPLE:
//
//	s := []int{0, 1, 2, 3, 4}
//	ContainsAny(s, 0)    ⏩ true
//	ContainsAny(s, 5)    ⏩ false
//	ContainsAny(s, 0, 1) ⏩ true
//	ContainsAny(s, 0, 5) ⏩ true
func ContainsAny[T comparable](s []T, vs ...T) bool {
	return iter.ContainsAny(vs, iter.StealSlice(s))
}

// ContainsAll returns whether all of given elements occur in slice.
//
// 🚀 EXAMPLE:
//
//	s := []int{0, 1, 2, 3, 4}
//	ContainsAll(s, 0)    ⏩ true
//	ContainsAll(s, 5)    ⏩ false
//	ContainsAll(s, 0, 1) ⏩ true
//	ContainsAll(s, 0, 5) ⏩ false
func ContainsAll[T comparable](s []T, vs ...T) bool {
	return iter.ContainsAll(vs, iter.StealSlice(s))
}

// Any determines whether any (at least one) element of the slice s
// satisfies the predicate f.
//
// Any supports short-circuit evaluation.
//
// 🚀 EXAMPLE:
//
//	Any([]int{1, 2, 3}, func(x int) bool { return x > 2 }) ⏩ true
//
// 💡 HINT:
//   - Use [All] to known whether all elements satisfies the predicate f
//   - Use [CountBy] to known how many elements satisfies the predicate f
func Any[T any](s []T, f func(T) bool) bool {
	return iter.Any(f, iter.StealSlice(s))
}

// All determines whether all elements of the slice s satisfy the predicate f.
//
// 🚀 EXAMPLE:
//
//	All([]int{1, 2, 3}, func(x int) bool { return x > 0 }) ⏩ true
func All[T any](s []T, f func(T) bool) bool {
	return iter.All(f, iter.StealSlice(s))
}

// Find returns the possible first element of slice that satisfies predicate f.
//
// 🚀 EXAMPLE:
//
//	s := []int{0, 1, 2, 3, 4}
//	Find(s, func(v int) bool { return v > 0 }) ⏩ goption.OK(1)
//	Find(s, func(v int) bool { return v < 0 }) ⏩ goption.Nil[int]()
//
// 💡 HINT:
//
//   - Use [Contains] if you just want to know whether the value exists
//   - Use [IndexBy] if you want to know the index of value
//   - Use [FindRev] if you want to find in reverse order
//   - Use [Count] if you want to count the occurrences of element
//
// 💡 AKA: Search
func Find[T any](s []T, f func(T) bool) goption.O[T] {
	for _, v := range s {
		if f(v) {
			return goption.OK(v)
		}
	}
	return goption.Nil[T]()
}

// FindRev is a variant of [Find] in reverse order.
//
// 🚀 EXAMPLE:
//
//	s := []int{0, 1, 2, 3, 4}
//	FindRev(s, func(v int) bool { return v > 0 }) ⏩ goption.OK(4)
//	FindRev(s, func(v int) bool { return v < 0 }) ⏩ goption.Nil[int]()
func FindRev[T any](s []T, f func(T) bool) goption.O[T] {
	for i := len(s) - 1; i >= 0; i-- {
		if f(s[i]) {
			return goption.OK(s[i])
		}
	}
	return goption.Nil[T]()
}

// Remove removes all element v from the slice s, returns a newly allocated slice.
//
// 🚀 EXAMPLE:
//
//	Remove([]int{0, 1, 2, 3, 4}, 3)    ⏩ []int{0, 1, 2, 4}
//	Remove([]int{0, 1, 3, 2, 3, 4}, 3) ⏩ []int{0, 1, 2, 4}
//
// 💡 HINT:
//
//   - Use [Compact] if you just want to remove zero value.
//   - Use [RemoveIndex] if you want to remove value by index
//
// 💡 AKA: Delete
func Remove[S ~[]T, T comparable](s S, v T) S {
	return iter.ToSlice(iter.Remove(v, iter.FromSlice(s)))
}

// Chunk splits a slice into length-n chunks and returns chunks by a newly allocated slice.
//
// The last chunk will be shorter if n does not evenly divide the length of the list.
//
// 🚀 EXAMPLE:
//
//	Chunk([]int{0, 1, 2, 3, 4}, 2) ⏩ [][]int{{0, 1}, {2, 3}, {4}}
//
// 💡 HINT:
//
//   - If you want to split list into n chunks, use function [Divide].
//   - This function returns sub-slices of original slice,
//     if you modify the sub-slices, the original slice is modified too.
//     Use [ChunkClone] to prevent this.
//   - Use [Flatten] to restore chunks to flat slice.
//
// 💡 AKA: Page, Pagination
func Chunk[S ~[]T, T any](s S, size int) []S {
	rtassert.MustLessThan(size, 1)
	ret := make([]S, 0, (len(s)+size-1)/size)
	for i := 0; i < len(s); i += size {
		ret = append(ret, Slice(s, i, i+size))
	}
	return ret
}

// ChunkClone is variant of function [Chunk].
// It clones the original slice before chunking it.
func ChunkClone[S ~[]T, T any](s S, size int) []S {
	rtassert.MustLessThan(size, 1)
	ret := make([]S, 0, (len(s)+size-1)/size)
	for i := 0; i < len(s); i += size {
		ret = append(ret, SliceClone(s, i, i+size))
	}
	return ret
}

// Divide splits a list into exactly n slices and returns chunks by a newly allocated slice.
//
// The length of chunks will be different if n does not evenly divide the length
// of the slice.
//
// 🚀 EXAMPLE:
//
//	s := []int{0, 1, 2, 3, 4}
//	Divide(s, 2)       ⏩ [][]int{{0, 1, 2},  {3, 4}}
//	Divide(s, 3)       ⏩ [][]int{{0, 1}, {2, 3}, {4}}
//	Divide([]int{}, 2) ⏩ [][]int{{}, {}}
//
// 💡 HINT:
//
//   - If you want to split list into length-n chunks, use [Chunk].
//   - This function returns sub-slices of original slice,
//     if you modify the sub-slices, the original slice is modified too.
//     Use [DivideClone] to prevent this.
//   - Use [Flatten] to restore chunks to flat slice.
//
// 💡 AKA: Page, Pagination
func Divide[S ~[]T, T any](s S, n int) []S {
	rtassert.MustLessThan(n, 1)
	k := len(s) / n // Every chunk have at least k elements
	m := len(s) % n // The first m chunks have an extra element
	ret := make([]S, 0, n)
	for i := 0; i < n; i++ {
		ret = append(ret, Slice(s, i*k+gvalue.Min(i, m), (i+1)*k+gvalue.Min(i+1, m)))
	}
	return ret
}

// DivideClone is variant of function Divide.
// It clones the original slice before dividing it.
func DivideClone[S ~[]T, T any](s S, n int) []S {
	rtassert.MustLessThan(n, 1)
	k := len(s) / n // Every chunk have at least k elements
	m := len(s) % n // The first m chunks have an extra element
	ret := make([]S, 0, n)
	for i := 0; i < n; i++ {
		ret = append(ret, SliceClone(s, i*k+gvalue.Min(i, m), (i+1)*k+gvalue.Min(i+1, m)))
	}
	return ret
}

// GroupBy adjacent elements according to key returned by function f.
//
// 🚀 EXAMPLE:
//
//	GroupBy([]int{1, 2, 3, 4},
//	func(v int) string {
//	    return gcond.If(v%2 == 0, "even", "odd")
//	})
//
//	⏩
//
//	map[string][]int{
//	    "odd": {1, 3},
//	    "even": {2, 4},
//	}
//
// 💡 HINT: If function f returns bool, use [Partition] instead.
func GroupBy[S ~[]T, K comparable, T any](s S, f func(T) K) map[K]S {
	m := make(map[K]S)
	for i := range s {
		k := f(s[i])
		m[k] = append(m[k], s[i])
	}
	return m
}

// Uniq returns the distinct elements of slice.
// Elements are ordered by their first occurrence.
//
// 🚀 EXAMPLE:
//
//	Uniq([]int{0, 1, 4, 3, 1, 4}) ⏩ []int{0, 1, 4, 3}
//
// 💡 HINT:
//
//   - If type is not comparable, use [UniqBy].
//   - If you need  duplicate elements, use [Dup].
//
// 💡 AKA: Distinct, Dedup, Unique
func Uniq[S ~[]T, T comparable](s S) S {
	return iter.ToSlice(iter.Uniq(iter.FromSlice(s)))
}

// UniqBy returns the distinct elements of slice with key function f.
// The result is a newly allocated slice without duplicate elements.
//
// 🚀 EXAMPLE:
//
//	type Foo struct{ Value int }
//	s := []Foo{{0}, {1}, {4}, {3}, {1}, {4}}
//	UniqBy(s, func(v Foo) int { return v.Value }) ⏩ []Foo{{0}, {1}, {4}, {3}}
//
// 💡 AKA: DistinctBy, DedupBy.
func UniqBy[S ~[]T, K comparable, T any](s S, f func(T) K) S {
	return iter.ToSlice(iter.UniqBy(f, iter.FromSlice(s)))
}

// Dup returns the repeated elements of slice.
// The result are sorted in order of recurrence.
//
// 🚀 EXAMPLE:
//
//	Dup([]int{0, 1, 1, 1})    ⏩ []int{1}
//	Dup([]int{3, 2, 2, 3, 3}) ⏩ []int{2, 3} // in order of recurrence
//
// 💡 HINT:
//
//   - If type is not comparable, use [DupBy].
//   - If you need distinct elements, use [Uniq].
//
// 💡 AKA: Duplicate.
func Dup[S ~[]T, T comparable](s S) S {
	return iter.ToSlice(iter.Dup(iter.FromSlice(s)))
}

// DupBy returns the repeated elements of slice with key function f.
// The result is a newly allocated slice contains duplicate elements.
// The result are sorted in order of recurrence.
//
// 🚀 EXAMPLE:
//
//	type Foo struct{ Value int }
//	s := []Foo{{3}, {2}, {2}, {3}, {3}}
//	DupBy(s, func(v Foo) int { return v.Value }) ⏩ []Foo{{2}, {3}}
//
// 💡 AKA: DuplicateBy.
func DupBy[S ~[]T, K comparable, T any](s S, f func(T) K) S {
	return iter.ToSlice(iter.DupBy(f, iter.FromSlice(s)))
}

// Repeat returns a slice with value v repeating exactly n times.
// The result is an empty slice if n is 0.
//
// ⚠️ WARNING: The function panics if n is negative.
//
// 🚀 EXAMPLE:
//
//	Repeat(123, -1) ⏩ ❌PANIC❌
//	Repeat(123, 0)  ⏩ []int{}
//	Repeat(123, 3)  ⏩ []int{123, 123, 123}
//
// 💡 HINT: The result slice contains shallow copy of element v. Use [RepeatBy] with a copier if deep copy is necessary.
func Repeat[T any](v T, n int) []T {
	if n < 0 {
		panic("repeat count is negative")
	}
	return iter.Repeat(v).Next(n)
}

// RepeatBy returns a slice with elements generated by calling fn exactly n times.
// The result is an empty slice if n is 0.
//
// ⚠️ WARNING:
//   - The function panics if n is negative.
//
// 🚀 EXAMPLE:
//
//	fn := func() *int { v := 123; return &v }
//	RepeatBy(fn, -1) ⏩ ❌PANIC❌
//	RepeatBy(fn, 0)  ⏩ []*int{}
//	RepeatBy(fn, 3)  ⏩ []*int{ &int(123), &int(123), &int(123) } // different addresses!
func RepeatBy[T any](fn func() T, n int) []T {
	if n < 0 {
		panic("repeat count is negative")
	}
	result := make([]T, n)
	for i := 0; i < n; i++ {
		result[i] = fn()
	}
	return result
}

// Max returns maximum element of slice s.
// If the given slice is empty, goption.Nil[T]() is returned.
//
// 🚀 EXAMPLE:
//
//	Max([]int{0, 1, 4, 3, 1, 4}) ⏩ goption.OK(4)
//	Max([]int{})                 ⏩ goption.Nil[int]()
//
// 💡 HINT: If type is not orderable, use [MaxBy].
func Max[T constraints.Ordered](s []T) goption.O[T] {
	return iter.Max(iter.StealSlice(s))
}

// MaxBy returns the maximum element of slice s determined by function less.
// If the given slice is empty, goption.Nil[T]() is returned.
//
// 🚀 EXAMPLE:
//
//	type Foo struct { Value int }
//	less := func(x, y Foo) bool { return x.Value < y.Value }
//	s := []Foo{{10}, {1}, {-1}, {100}, {3}}
//	MaxBy(s, less) ⏩ goption.OK(Foo{100})
func MaxBy[T any](s []T, less func(T, T) bool) goption.O[T] {
	return iter.MaxBy(less, iter.StealSlice(s))
}

// Min returns the minimum element of slices s.
// If the given slice is empty, goption.Nil[T]() is returned.
//
// 🚀 EXAMPLE:
//
//	Min([]int{1, 4, 3, 1, 4}) ⏩ goption.OK(1)
//	Min([]int{})              ⏩ goption.Nil[int]()
//
// 💡 HINT: If type is not orderable, use [MinBy].
func Min[T constraints.Ordered](s []T) goption.O[T] {
	return iter.Min(iter.StealSlice(s))
}

// MinBy returns the minimum element of slices s determined by function less.
// If the given slice is empty, goption.Nil[T]() is returned.
//
// 🚀 EXAMPLE:
//
//	type Foo struct { Value int }
//	less := func(x, y Foo) bool { return x.Value < y.Value }
//	MinBy([]Foo{{10}, {1}, {-1}, {100}, {3}}, less) ⏩ goption.OK(Foo{-1})
func MinBy[T any](s []T, less func(T, T) bool) goption.O[T] {
	return iter.MinBy(less, iter.StealSlice(s))
}

// MinMax returns both minimum and maximum elements of slice s.
// If the given slice is empty, goption.Nil[tuple.T2[T, T]]() is returned.
//
// 🚀 EXAMPLE:
//
//	MinMax([]int{})                 ⏩ goption.Nil[int]()
//	MinMax([]int{1})                ⏩ goption.OK(tuple.T2{1, 1})
//	MinMax([]int{0, 1, 4, 3, 1, 4}) ⏩ goption.OK(tuple.T2{0, 4})
//
// 💡 HINT: If type is not orderable, use [MinMaxBy].
//
// 💡 AKA: Bound
func MinMax[T constraints.Ordered](s []T) goption.O[tuple.T2[T, T]] {
	return iter.MinMax(iter.StealSlice(s))
}

// MinMaxBy returns both minimum and maximum elements of slice s.
// If the given slice is empty, goption.Nil[tuple.T2[T, T]]() is returned.
//
// 🚀 EXAMPLE:
//
//	type Foo struct { Value int }
//	less := func(x, y Foo) bool { return x.Value < y.Value }
//	MinMaxBy([]Foo{{10}, {1}, {-1}, {100}, {3}}, less) ⏩ goption.OK(tuple.T2{Foo{-1}, Foo{100}})
//
// 💡 NOTE: The returned min and max elements may be the same object when each
// element of the slice is equal
//
// 💡 AKA: BoundBy
func MinMaxBy[T any](s []T, less func(T, T) bool) goption.O[tuple.T2[T, T]] {
	return iter.MinMaxBy(less, iter.StealSlice(s))
}

// Clone returns a shallow copy of the slice.
// If the given slice is nil, nil is returned.
//
// 🚀 EXAMPLE:
//
//	Clone([]int{1, 2, 3}) ⏩ []int{1, 2, 3}
//	Clone([]int{})        ⏩ []int{}
//	Clone[int](nil)       ⏩ nil
//
// 💡 HINT: The elements are copied using assignment (=), so this is a shallow clone.
// If you want to do a deep clone, use [CloneBy] with an appropriate element
// clone function.
//
// 💡 AKA: Copy
func Clone[S ~[]T, T any](s S) S {
	if s == nil {
		return nil
	}
	return iter.ToSlice(iter.FromSlice(s))
}

// CloneBy is variant of [Clone], it returns a copy of the slice.
// Elements are copied using function clone.
// If the given slice is nil, nil is returned.
//
// 💡 AKA: CopyBy
func CloneBy[S ~[]T, T any](s S, f func(T) T) S {
	if s == nil {
		return nil
	}
	return Map(s, f)
}

// Flatten collapses a tow-dimension slice to one dimension.
//
// 🚀 EXAMPLE:
//
//	Flatten([][]int{{0}, {1, 2}, {3, 4}}) ⏩ []int{0, 1, 2, 3, 4}
//
// 💡 HINT: Use [FlatMap] if you want to flatten non-slice elements.
func Flatten[S ~[]T, T any](s []S) S {
	return iter.ToSlice(
		iter.FlatMap(func(v []T) []T { return v },
			iter.Map(func(s S) []T { return s },
				iter.StealSlice(s))))
}

// FlatMap applies function f to each element of slice s with type F.
// Results of f are flatten and returned as a newly allocated slice with type T.
//
// 🚀 EXAMPLE:
//
//	type Site struct{ urls []string }
//	func (s Site) URLs() []string { return s.urls }
//
//	sites := []Site{
//		{[]string{"url1", "url2"}},
//		{[]string{"url3", "url4"}},
//	}
//
//	FlatMap(sites, Site.URLs) ⏩ []string{"url1", "url2", "url3", "url4"}
//
// 💡 HINT:
//
//   - Use [Flatten] if the element of given slice is also slice.
//   - Use [FilterMap] if you want to ignore some element during mapping
func FlatMap[F, T any](s []F, f func(F) []T) []T {
	return iter.ToSlice(iter.FlatMap(f, iter.StealSlice(s)))
}

// First returns the possible first element of slice s.
// If the given slice is empty, goption.Nil[T]() is returned.
//
// 🚀 EXAMPLE:
//
//	First([]int{4, 3, 1, 4}) ⏩ goption.OK(4)
//	First([]int{})           ⏩ goption.Nil[int]()
//
// 💡 HINT: Use [Get] to access element at any index.
//
// 💡 AKA: Head
func First[T any](s []T) goption.O[T] {
	if len(s) == 0 {
		return goption.Nil[T]()
	}
	return goption.OK(s[0])
}

// Get returns the possible element at index n.
//
// [Negative index] is supported. For example:
//
//   - Get(s, 0) returns the [First] element
//   - Get(s, -1) returns the [Last] element
//
// 🚀 EXAMPLE:
//
//	s := []int{1, 2, 3, 4}
//	Get(s, 0)  ⏩ goption.OK(1)
//	Get(s, 1)  ⏩ goption.OK(2)
//	Get(s, -1) ⏩ goption.OK(4)
//	Get(s, -2) ⏩ goption.OK(3)
//
// 💡 AKA: Nth, At, Access, ByIndex, Load
func Get[T any, I constraints.Integer](s []T, n I) goption.O[T] {
	index, ok := normalizeIndex(s, n)
	if !ok {
		return goption.Nil[T]()
	}
	return goption.OK(s[index])
}

// Last returns the possible last element of slice s.
// If the given slice is empty, goption.Nil[T]() is returned.
//
// 🚀 EXAMPLE:
//
//	Last([]int{4, 3, 1, 5}) ⏩ goption.OK(5)
//	Last([]int{})           ⏩ goption.Nil[int]()
//
// 💡 HINT: Use [Get] to access element at any index.
//
// 💡 AKA: Tail
func Last[T any](s []T) goption.O[T] {
	if len(s) == 0 {
		return goption.Nil[T]()
	}
	return goption.OK(s[len(s)-1])
}

// Union returns the unions of slices as a newly allocated slices.
//
// 💡 NOTE: If the result is an empty set, always return an empty slice instead of nil
//
// 🚀 EXAMPLE:
//
//	Union([]int{1, 2, 3}, []int{3, 4, 5}) ⏩ []int{1, 2, 3, 4, 5}
//	Union([]int{1, 2, 3}, []int{})        ⏩ []int{1, 2, 3}
//	Union([]int{}, []int{3, 4, 5})        ⏩ []int{3, 4, 5}
//
// 💡 HINT: if you need a set data structure,
// use [github.com/bytedance/gg/collection/set].
func Union[S ~[]T, T comparable](ss ...S) S {
	if len(ss) == 0 {
		return S{}
	}
	if len(ss) == 1 {
		return Uniq(ss[0])
	}
	members := set.New[T]()
	ret := S{} // TODO: Guess a cap.
	for _, s := range ss {
		for _, v := range s {
			if members.Add(v) {
				ret = append(ret, v)
			}
		}
	}
	return ret
}

// Diff returns the difference of slice s against other slices as a newly allocated slice.
//
// 💡 NOTE: If the result is an empty set, always return an empty slice instead of nil
//
// 🚀 EXAMPLE:
//
//	Diff([]int{1, 2, 3}, []int{3, 4, 5}) ⏩ []int{1, 2}
//	Diff([]int{1, 2, 3}, []int{4, 5, 6}) ⏩ []int{1, 2, 3}
//	Diff([]int{1, 2, 3}, []int{1, 2, 3}) ⏩ []int{}
//
// 💡 HINT: if you need a set data structure,
// use [github.com/bytedance/gg/collection/set].
func Diff[S ~[]T, T comparable](s S, againsts ...S) S {
	if len(s) == 0 {
		return []T{}
	}
	if len(againsts) == 0 {
		return Uniq(s)
	}
	members := set.New(s...)
	for _, s := range againsts {
		for _, v := range s {
			members.Remove(v)
		}
	}
	if members.Len() == 0 {
		return S{}
	}
	ret := make(S, 0, members.Len())
	for _, v := range s {
		if members.Remove(v) {
			ret = append(ret, v)
			if members.Len() == 0 {
				return ret
			}
		}
	}
	return ret // must not reach
}

// Intersect returns the intersection of slices as a newly allocated slice.
//
// 💡 NOTE: If the result is an empty set, always return an empty slice instead of nil
//
// 🚀 EXAMPLE:
//
//	Intersect([]int{1, 2, 3}, []int{2, 3, 4}) ⏩ []int{2, 3}
//	Intersect([]int{1, 2, 3}, []int{4, 5, 6}) ⏩ []int{}
//	Intersect([]int{1, 2, 3}, []int{1, 2, 3}) ⏩ []int{1, 2, 3}
//
// 💡 HINT: if you need a set data structure,
// use [github.com/bytedance/gg/collection/set].
func Intersect[S ~[]T, T comparable](ss ...S) S {
	if len(ss) == 0 {
		return S{}
	}
	if len(ss) == 1 {
		return Uniq(ss[0])
	}
	if len(ss[0]) == 0 {
		return S{}
	}
	members := set.New(ss[0]...)
	for _, s := range ss[1:] {
		if len(s) == 0 {
			return S{}
		}
		members.IntersectInplace(set.New(s...))
	}
	if members.Len() == 0 {
		return S{}
	}
	ret := make(S, 0, members.Len())
	for _, s := range ss {
		for _, v := range s {
			if members.Remove(v) {
				ret = append(ret, v)
				if members.Len() == 0 {
					return ret
				}
			}
		}
	}
	return ret // must not reach
}

// Reverse reverses the elements of slices.
//
// 💡 HINT: If you want to reverse in a newly allocated slice, use [ReverseClone].
func Reverse[T any](s []T) {
	_ = iter.ToSlice(iter.Reverse(iter.StealSlice(s)))
}

// ReverseClone is variant of [Reverse].
// It clones the original slice before reversing it.
func ReverseClone[S ~[]T, T any](s S) S {
	return iter.ToSlice(iter.Reverse(iter.FromSlice(s)))
}

// Sort sorts elements of slice in ascending order (from small to large).
//
// 🚀 EXAMPLE:
//
//	s := []int{1, 3, 2, 4}
//	Sort(s) ⏩ []int{1, 2, 3, 4}
//
// 💡 HINT:
//
//   - Sort in a newly allocated slice, use [SortClone]
//   - Sort by a custom comparison function, use [SortBy]
//   - Sort in descending order,
//     use [SortBy] + [github.com/bytedance/gg/gvalue.Greater]
//
// 💡 AKA: Order
func Sort[T constraints.Ordered](s []T) {
	_ = iter.ToSlice(iter.Sort(iter.StealSlice(s)))
}

// SortClone is variant of [Sort].
// It clones the original slice before sorting it.
func SortClone[S ~[]T, T constraints.Ordered](s S) S {
	return iter.ToSlice(iter.Sort(iter.FromSlice(s)))
}

// SortBy sorts elements of slices i with function less.
//
// 💡 AKA: OrderBy
func SortBy[T any](s []T, less func(T, T) bool) {
	_ = iter.ToSlice(iter.SortBy(less, iter.StealSlice(s)))
}

// SortCloneBy is variant of function [SortBy].
// It clones the original slice before sorting it.
func SortCloneBy[S ~[]T, T any](s S, less func(T, T) bool) S {
	return iter.ToSlice(iter.SortBy(less, iter.FromSlice(s)))
}

// StableSortBy is variant of [SortBy], it keeps the original order of equal elements
// when sorting.
func StableSortBy[T any](s []T, less func(T, T) bool) {
	_ = iter.ToSlice(iter.StableSortBy(less, iter.StealSlice(s)))
}

// TypeAssert converts a slice from type From to type To by type assertion.
//
// 🚀 EXAMPLE:
//
//	TypeAssert[int]([]any{1, 2, 3, 4})   ⏩ []int{1, 2, 3, 4}
//	TypeAssert[any]([]int{1, 2, 3, 4})   ⏩ []any{1, 2, 3, 4}
//	TypeAssert[int64]([]int{1, 2, 3, 4}) ⏩ ❌PANIC❌
//
// ⚠️ WARNING:
//
//   - This function may ❌PANIC❌.
//     See [github.com/bytedance/gg/gvalue.TypeAssert] for more details
func TypeAssert[To, From any](s []From) []To {
	return Map(s, gvalue.TypeAssert[To, From])
}

// ForEach applies function f to each element of slice s.
//
// 💡 HINT: Use [ForEachIndexed] If you want to get element with index.
func ForEach[T any](s []T, f func(v T)) {
	iter.ForEach(f, iter.StealSlice(s))
}

// ForEachIndexed applies function f to each element of slice s.
// The argument i of function f represents the zero-based index of that element
// of slice.
func ForEachIndexed[T any](s []T, f func(i int, v T)) {
	iter.ForEachIndexed(f, iter.StealSlice(s))
}

// Equal returns whether two slices are equal.
//
// 🚀 EXAMPLE:
//
//	Equal([]int{1, 2, 3}, []int{1, 2, 3})    ⏩ true
//	Equal([]int{1, 2, 3}, []int{1, 2, 3, 4}) ⏩ false
//	Equal([]int{}, []int{})                  ⏩ true
//	Equal([]int{}, nil)                      ⏩ true
func Equal[T comparable](s1, s2 []T) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

// EqualBy returns whether two slices are equal by function eq.
//
// 🚀 EXAMPLE:
//
//	eq := gvalue.Equal[int]
//	EqualBy([]int{1, 2, 3}, []int{1, 2, 3}, eq)    ⏩ true
//	EqualBy([]int{1, 2, 3}, []int{1, 2, 3, 4}, eq) ⏩ false
//	EqualBy([]int{}, []int{}, eq)                  ⏩ true
//	EqualBy([]int{}, nil, eq)                      ⏩ true
func EqualBy[T any](s1, s2 []T, eq func(T, T) bool) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if !eq(s1[i], s2[i]) {
			return false
		}
	}
	return true
}

// ToMap collects elements of slice to map, both map keys and values are produced
// by mapping function f.
//
// 🚀 EXAMPLE:
//
//	type Foo struct {
//		ID   int
//		Name string
//	}
//	mapper := func(f Foo) (int, string) { return f.ID, f.Name }
//	ToMap([]Foo{}, mapper) ⏩ map[int]string{}
//	s := []Foo{{1, "one"}, {2, "two"}, {3, "three"}}
//	ToMap(s, mapper)       ⏩ map[int]string{1: "one", 2: "two", 3: "three"}
func ToMap[T, V any, K comparable](s []T, f func(T) (K, V)) map[K]V {
	return iter.ToMap(f, iter.StealSlice(s))
}

// ToMapValues collects elements of slice to values of map, the map keys are
// produced by mapping function f.
//
// 🚀 EXAMPLE:
//
//	type Foo struct {
//	    ID int
//	}
//	id := func(f Foo) int { return f.ID }
//	ToMapValues([]Foo{}, id)                    ⏩ map[int]Foo{}
//	ToMapValues([]Foo{ {1}, {2}, {1}, {3}}, id) ⏩ map[int]Foo{1: {1}, 2: {2}, 3: {3}}
//
// 💡 AKA: Kotlin's associateBy
func ToMapValues[T any, K comparable](s []T, f func(T) K) map[K]T {
	return iter.ToMapValues(f, iter.StealSlice(s))
}

// PtrOf returns pointers that point to equivalent elements of slice s.
// ([]T → []*T).
//
// 🚀 EXAMPLE:
//
//	PtrOf([]int{1, 2, 3}) ⏩ []*int{ (*int)(1), (*int)(2), (*int)(3) },
//
// ⚠️  WARNING: The returned pointers do not point to elements of the original
// slice, user CAN NOT modify the element by modifying the pointer.
func PtrOf[T any](s []T) []*T {
	return Map(s, gptr.Of[T])
}

// Indirect returns the values pointed to by the pointers.
// If the element is nil, filter it out of the returned slice.
//
// 🚀 EXAMPLE:
//
//	v1, v2 := 1, 2
//	Indirect([]*int{ &v1, &v2, nil})  ⏩ []int{1, 2}
//
// 💡 HINT: If you want to replace nil pointer with default value,
// use [IndirectOr].
func Indirect[T any](s []*T) []T {
	return iter.ToSlice(
		iter.Map(gptr.Indirect[T],
			iter.Filter(gptr.IsNotNil[T],
				iter.FromSlice(s))))
}

// IndirectOr safely dereferences slice of pointers.
// If the pointer is nil, returns the value fallback instead.
//
// 🚀 EXAMPLE:
//
//	v1, v2 := 1, 2
//	IndirectOr([]*int{ &v1, &v2, nil}, -1)  ⏩ []int{1, 2, -1}
func IndirectOr[T any](s []*T, fallback T) []T {
	return Map(s, gfunc.Partial2(gptr.IndirectOr[T]).PartialR(fallback))
}

// Shuffle pseudo-randomizes the order of elements.
//
// Shuffle is 2x ~ 40x(parallel) faster than [math/rand.Shuffle].
//
// 💡 HINT: If you want to shuffle in a newly allocated slice, use [ShuffleClone] .
func Shuffle[T any](s []T) {
	_ = iter.ToSlice(iter.Shuffle(iter.StealSlice(s)))
}

// ShuffleClone is variant of [Shuffle].
// It clones the original slice before shuffling it.
func ShuffleClone[S ~[]T, T any](s S) S {
	return iter.ToSlice(iter.Shuffle(iter.FromSlice(s)))
}

// Index returns the index of the first occurrence of element in slice s,
// or nil if not present.
//
// 🚀 EXAMPLE:
//
//	s := []string{"a", "b", "b", "d"}
//	Index(s, "b") ⏩ goption.OK(1)
//	Index(s, "e") ⏩ goption.Nil[int]()
//
// 💡 HINT:
//
//   - Use [IndexBy] if complex comparison logic is required (instead of just ==)
//   - Use [Contains] if you just want to know whether the value exists
//   - Use [IndexRev] if you want to index element in reverse order.
func Index[T comparable](s []T, e T) goption.O[int] {
	for i := range s {
		if e == s[i] {
			return goption.OK(i)
		}
	}
	return goption.Nil[int]()
}

// IndexRev is a variant of [Index] in reverse order.
//
// 🚀 EXAMPLE:
//
//	s := []string{"a", "b", "b", "d"}
//	IndexRev(s, "b") ⏩ goption.OK(2)
//	IndexRev(s, "e") ⏩ goption.Nil[int]()
func IndexRev[T comparable](s []T, e T) goption.O[int] {
	for i := len(s) - 1; i >= 0; i-- {
		if e == s[i] {
			return goption.OK(i)
		}
	}
	return goption.Nil[int]()
}

// IndexBy is variant of [Index], returns the first index of element that
// satisfying predicate f, or nil if none do.
func IndexBy[T any](s []T, f func(T) bool) goption.O[int] {
	for i := range s {
		if f(s[i]) {
			return goption.OK(i)
		}
	}
	return goption.Nil[int]()
}

// IndexRevBy is variant of [IndexRev], returns the first index of element that
// satisfying predicate f, or nil if none do.
func IndexRevBy[T any](s []T, f func(T) bool) goption.O[int] {
	for i := len(s) - 1; i >= 0; i-- {
		if f(s[i]) {
			return goption.OK(i)
		}
	}
	return goption.Nil[int]()
}

// Take returns the first n elements of slices s if 0 <= n <= len(s), or slice itself if n > len(s).
// If -len(s) <= n < 0, returns the last -n elements of slice s, or slice itself if n < -len(s).
//
// 🚀 EXAMPLE:
//
//	s := []int{1, 2, 3, 4, 5}
//	Take(s, 0)   ⏩ []int{}
//	Take(s, 3)   ⏩ []int{1, 2, 3}
//	Take(s, 10)  ⏩ []int{1, 2, 3, 4, 5}
//	Take(s, -1)  ⏩ []int{5}
//	Take(s, -3)  ⏩ []int{3, 4, 5}
//	Take(s, -10) ⏩ []int{1, 2, 3, 4, 5}
//
// 💡 HINT: This function returns sub-slices of original slice,
// if you modify the sub-slices, the original slice is modified too.
// Use [TakeClone] to prevent this.
func Take[S ~[]T, I constraints.Integer, T any](s S, n I) S {
	startIdx, endIdx := 0, int(n)
	if n < 0 {
		endIdx = len(s)
		startIdx, _ = normalizeIndex(s, n)
		if startIdx < 0 {
			startIdx = 0
		}
	} else if endIdx > len(s) {
		endIdx = len(s)
	}
	return s[startIdx:endIdx]
}

// TakeClone is variant of [Take].
func TakeClone[S ~[]T, I constraints.Integer, T any](s S, n I) S {
	return Clone(Take(s, n))
}

// Slice returns a sub-slice of the slice S that contains the elements starting
// from the start-th element up to but not including the end-th element "[start:end)".
// In other words, it is safer replacement of [Slice Expression].
//
//   - Slice(s, 0, 3) 🟰 s[:3]
//   - Slice(s, 1, 3) 🟰 s[1:3]
//
// [Negative index] is supported:
//
//   - Slice(s, -3, -1) 🟰 s[len(s)-3:len(s)-1]
//   - Slice(s, -3, 0)  🟰 s[len(s)-3:] specially, the 0 at the end implies the end slice.
//
// 🚀 EXAMPLE:
//
//	s := []int{1, 2, 3, 4, 5}
//	Slice(s, 0, 3)     ⏩ []int{1, 2, 3}
//	Slice(s, 1, 3)     ⏩ []int{2, 3}
//	Slice(s, 0, 0)     ⏩ []int{}
//	Slice(s, 0, 100)   ⏩ []int{1, 2, 3, 4, 5}  // won't PANIC even out of range
//	Slice(s, 100, 99)  ⏩ []int{}               // won't PANIC even out of range
//	Slice(s, -3, -1)   ⏩ []int{3, 4}           // equal to Slice(s, 2, 4) and Slice(s, -3, 4)
//	Slice(s, -1, 0)    ⏩ []int{5}              // specially, the 0 at the end implies the end slice
//
// 💡 HINT: This function returns sub-slices of original slice,
// if you modify the sub-slices, the original slice is modified too.
// Use [SliceClone] to prevent this.
//
// [Slice Expression]: https://tip.golang.org/ref/spec#Slice_expressions
func Slice[S ~[]T, I constraints.Integer, T any](s S, start, end I) S {
	// Handle the negative index
	startIdx, _ := normalizeIndex(s, start)
	// Particularly, 0 in the right endpoint and the light endpoint is negative
	// implies the 0 is equal the last slice.
	var endIdx int
	if start < 0 && end == 0 {
		endIdx = len(s)
	} else {
		endIdx, _ = normalizeIndex(s, end)
	}

	if startIdx < 0 {
		startIdx = 0
	}
	if endIdx > len(s) {
		endIdx = len(s)
	}
	if startIdx >= endIdx {
		return S{}
	}

	return s[startIdx:endIdx]
}

// SliceClone is variant of [Slice].
func SliceClone[S ~[]T, I constraints.Integer, T any](s S, start, end I) S {
	return Clone(Slice(s, start, end))
}

// Drop drops the first n elements of slices s, returns the remaining part of
// slice, or empty slice if n > len(s).
//
// 🚀 EXAMPLE:
//
//	s := []int{1, 2, 3, 4, 5}
//	Drop(s, 0)  ⏩ []int{1, 2, 3, 4, 5}
//	Drop(s, 3)  ⏩ []int{4, 5}
//	Drop(s, 10) ⏩ []int{}
//
// ⚠️ WARNING: Panic when n < 0.
//
// 💡 NOTE: This function returns sub-slices of original slice,
// if you modify the sub-slices, the original slice is modified too.
// Use [DropClone] to prevent this.
func Drop[S ~[]T, T any](s S, n int) S {
	rtassert.MustNotNeg(n)
	if n > len(s) {
		n = len(s)
	}
	return s[n:]
}

// DropClone is variant of [Drop].
func DropClone[S ~[]T, T any](s S, n int) S {
	return Clone(Drop(s, n))
}

// Sum returns the arithmetic sum of the elements of slice s.
//
// 🚀 EXAMPLE:
//
//	Sum([]int{1, 2, 3, 4, 5})     ⏩ 15
//	Sum([]float64{1, 2, 3, 4, 5}) ⏩ 15.0
//
// ⚠️ WARNING: The returned type is still T, it may overflow for smaller types
// (such as int8, uint8).
func Sum[T constraints.Number](s []T) T {
	return iter.Sum(iter.StealSlice(s))
}

// SumBy applies function f to each element of slice s,
// returns the arithmetic sum of function result.
func SumBy[T any, N constraints.Number](s []T, f func(T) N) N {
	return iter.SumBy(f, iter.StealSlice(s))
}

// Avg returns the arithmetic mean of the elements of slice s.
//
// 🚀 EXAMPLE:
//
//	Avg([]int{1, 2, 3, 4, 5})      ⏩ 3.0
//	Avg([]float64{1, 2, 3, 4, 5})  ⏩ 3.0
//
// 💡 AKA: Mean, Average
func Avg[T constraints.Number](s []T) float64 {
	return iter.Avg(iter.StealSlice(s))
}

// AvgBy applies function f to each element of slice s,
// returns the arithmetic mean of function result.
//
// 💡 AKA: MeanBy, AverageBy
func AvgBy[T any, N constraints.Number](s []T, f func(T) N) float64 {
	return iter.AvgBy(f, iter.StealSlice(s))
}

// Len returns the length of slice s.
//
// 💡 HINT: This function is designed for high-order function, because the builtin
// function can not be used as function pointer.
// For example, if you want to get the total length of a 2D slice:
//
//	var s [][]int
//	total1 := SumBy(s, len)      // ❌ERROR❌ len (built-in) must be called
//	total2 := SumBy(s, Len[int]) // OK
func Len[T any](s []T) int {
	return len(s)
}

// Concat concatenates slices in order.
//
// 🚀 EXAMPLE:
//
//	Concat([]int{0}, []int{1, 2}, []int{3, 4}) ⏩ []int{0, 1, 2, 3, 4}
//
// 💡 AKA: Merge, Connect
func Concat[S ~[]T, T any](ss ...S) S {
	return Flatten(ss)
}

// Merge is alias of [Concat].
func Merge[S ~[]T, T any](ss ...S) S {
	return Flatten(ss)
}

// Compact removes all zero values from given slice s, returns a newly allocated slice.
//
// 🚀 EXAMPLE:
//
//	Compact([]int{0, 1, 2, 0, 3, 0, 0})     ⏩ []int{1, 2, 3}
//	Compact([]string{"", "foo", "", "bar"}) ⏩ []string{"foo", "bar"}
//
// 💡 HINT: See [github.com/bytedance/gg/gvalue.Zero] for details of zero value.
func Compact[S ~[]T, T comparable](s S) S {
	return Filter(s, gvalue.IsNotZero[T])
}

// Insert inserts elements vs before position pos, returns a newly allocated slice.
// [Negative index] is supported.
//
//   - Insert(x, 0, ...) inserts at the front of the slice
//   - Insert(x, len(x), ...) is equivalent to append(x, ...)
//   - Insert(x, -1, ...) is equivalent to Insert(x, len(x)-1, ...)
//
// 🚀 EXAMPLE:
//
//	s := []int{0, 1, 2, 3}
//	Insert(s, 0, 99)      ⏩ []int{99, 0, 1, 2, 3}
//	Insert(s, 0, 98, 99)  ⏩ []int{98, 99, 0, 1, 2, 3}
//	Insert(s, 4, 99)      ⏩ []int{0, 1, 2, 3, 99}
//	Insert(s, 1, 99)      ⏩ []int{0, 99, 1, 2, 3}
//	Insert(s, -1, 99)     ⏩ []int{0, 1, 2, 99, 3}
func Insert[S ~[]T, T any, I constraints.Integer](s S, pos I, vs ...T) S {
	if len(vs) == 0 {
		return Clone(s)
	}
	index, _ := normalizeIndex(s, pos)
	if index >= len(s) {
		index = len(s)
	} else if index < 0 {
		index = 0
	}

	dst := make(S, len(s)+len(vs))
	copy(dst, s[:index])
	copy(dst[index:], vs)
	copy(dst[index+len(vs):], s[index:])
	return dst
}

// insertInplace is a variant of [Insert], if the remaining capacity of the
// given slice is sufficient, the slice will be modified in place and returned.
func insertInplace[T any, I constraints.Integer](s []T, pos I, vs ...T) []T {
	l := len(s) + len(vs)
	// Slowpath: no enough capacity, allocate a new slice.
	if l > cap(s) {
		return Insert(s, pos, vs...)
	}
	if len(vs) == 0 {
		return s
	}
	index, _ := normalizeIndex(s, pos)
	if index >= len(s) {
		return append(s, vs...)
	}
	if index < 0 {
		index = 0
	}

	// Extend capacity to l, see https://silverrainz.me/notes/go/slice-expr.html#extend-capacity
	s = s[:l]
	copy(s[index+len(vs):], s[index:])
	copy(s[index:], vs)
	return s
}

// normalizeIndex normalizes possible [Negative index] to positive index.
// the returned bool indicate whether the normalized index is in range [0, len(s)).
func normalizeIndex[T any, I constraints.Integer](s []T, n I) (int, bool) {
	m := int(n)
	if m < 0 {
		m += len(s)
	}
	return m, m >= 0 && m < len(s)
}

// Of creates a slice from variadic arguments.
// If no argument given, an empty (non-nil) slice []T{} is returned.
//
// 💡 HINT: This function is used to omit verbose types like "[]LooooongTypeName{}"
// when constructing slices.
//
// 🚀 EXAMPLE:
//
//	Of(1, 2, 3) ⏩ []int{1, 2, 3}
//	Of(1)       ⏩ []int{1}
//	Of[int]()   ⏩ []int{}
func Of[T any](v ...T) []T {
	if len(v) == 0 {
		return []T{} // never return nil
	}
	return v
}

// RangeWithStep returns a slice of numbers from start (inclusive) to stop (exclusive)
// by step.
// If the interval does not exist, RangeWithStep returns an empty slice.
// If the step is positive, the returned slice is in ascending order.
// If the step is negative, the returned slice is in descending order.
//
// 🚀 EXAMPLE:
//
//	RangeWithStep(0, 0, 2)     ⏩ []int{}
//	RangeWithStep(0, -5, -1)   ⏩ []int{0, -1, -2, -3, -4}
//	RangeWithStep(0, 5, 2)     ⏩ []int{0, 2, 4}
//	RangeWithStep(0, 5, 3)     ⏩ []int{0, 3}
//	RangeWithStep(0.5, 2, 0.5) ⏩ []float64{0.5, 1, 1.5}
func RangeWithStep[I constraints.Number](start, stop, step I) []I {
	return iter.ToSlice(iter.RangeWithStep(start, stop, step))
}

// Range is a variant of RangeWithStep, with predefined step 1.
//
// 🚀 EXAMPLE:
//
//	Range(0, 0)    ⏩ []int{}
//	Range(0, -5)   ⏩ []int{}
//	Range(0, 5)    ⏩ []int{0, 1, 2, 3, 4}
func Range[I constraints.Number](start, stop I) []I {
	return iter.ToSlice(iter.Range(start, stop))
}

// RemoveIndex removes the element at index i from slice s and returns a newly allocated slice.
// If s[i] does not exist or is invalid, this function just clone the original slice.
// [Negative index] is supported.
//
//   - RemoveIndex(x, 0) 🟰 [Clone](s[1:])
//   - RemoveIndex(x, -1) 🟰 [Clone](s[0:len(x)-1])
//   - RemoveIndex(x, len(x)) 🟰 [Clone](s)
//
// 🚀 EXAMPLE:
//
//	RemoveIndex([]int{0, 1, 2, 3, 4}, 3)    ⏩ []int{0, 1, 2, 4}
//	RemoveIndex([]int{0, 1, 2, 3, 4}, -1)   ⏩ []int{0, 1, 2, 3}
//	RemoveIndex([]int{0, 1, 2, 3, 4}, 0)    ⏩ []int{1, 2, 3, 4}
//	RemoveIndex([]int{0, 1, 2, 3, 4}, 100)  ⏩ []int{0, 1, 2, 3, 4}
//
// 💡 Hint: This function has O(n) time complexity and ALWAYS returns a newly allocated slice.
//
// 💡 HINT: Use [Remove] if you want to remove elements by value
//
// 💡 AKA: DeleteIndex
func RemoveIndex[S ~[]T, I constraints.Integer, T any](s S, index I) S {
	idx, ok := normalizeIndex(s, int(index)) // conventionalize Index
	if !ok {
		return Clone(s) // fast path, not valid index. return the original slice
	}
	sLen := len(s) // delete from front
	if idx == 0 {
		return Clone(s[1:])
	} else if idx == sLen-1 {
		return Clone(s[0:idx])
	} else {
		return Concat(s[0:idx], s[idx+1:sLen])
	}
}

// Count returns the times of value v that occur in slice s.
//
// 🚀 EXAMPLE:
//
//	Count([]string{"a", "b", "c"}, "a") ⏩ 1
//	Count([]int{0, 1, 2, 0, 5, 3}, 0)   ⏩ 2
//
// 💡 HINT:
//
//   - Use [Contains] if you just want to know whether the element exitss or not
//   - Use [CountBy] if type of v is non-comparable
func Count[T comparable](s []T, v T) int {
	var count int
	for i := range s {
		if s[i] == v {
			count++
		}
	}
	return count
}

// CountBy returns the times of element in slice s that satisfy the predicate f.
//
// 🚀 EXAMPLE:
//
//	CountBy([]string{"a", "b", "c"}, func (v string) bool { return v < "b" }) ⏩ 1
//	CountBy([]int{0, 1, 2, 3, 4}, func (v int) bool { return v % i == 0 })    ⏩ 3
//
// 💡 HINT: Use [Any] if you just want to know whether at least one element satisfies predicate f.
func CountBy[T any](s []T, f func(T) bool) int {
	var count int
	for i := range s {
		if f(s[i]) {
			count++
		}
	}
	return count
}

// CountValues returns the occurrences of each element in slice s.
//
// 🚀 EXAMPLE:
//
//	CountValues([]string{"a", "b", "b"}) ⏩ map[string]int{"a": 1, "b": 2}
//	CountValues([]int{0, 1, 2, 0, 1, 1}) ⏩ map[int]int{0: 2, 1: 3, 2: 1}
//
// 💡 HINT:
//
//   - Use [CountValuesBy] if the element in slice s is non-comparable
func CountValues[T comparable](s []T) map[T]int {
	ret := make(map[T]int, len(s)/2)
	for i := range s {
		ret[s[i]]++
	}
	return ret
}

// CountValuesBy returns the times of each element in slice s that satisfy the predicate f.
//
// 🚀 EXAMPLE:
//
//	CountValuesBy([]int{0, 1, 2, 3, 4}, func(v int) bool { return v%2 == 0 }) ⏩ map[bool]int{true: 3, false: 2}
//	type Foo struct{ v int }
//	foos := []Foo{{1}, {2}, {3}}
//	CountValuesBy(foos, func(v Foo) bool { return v.v%2 == 0 }) ⏩ map[bool]int{true: 1, false: 2}
func CountValuesBy[K comparable, T any](s []T, f func(T) K) map[K]int {
	ret := make(map[K]int, len(s)/2)
	for i := range s {
		ret[f(s[i])]++
	}
	return ret
}
