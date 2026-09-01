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

// This file contains the method form of the slice operations.
//
// A slice []T can not be a method receiver, so the operations used to be
// package-level functions only. With Go 1.27 generic methods, a slice can be
// wrapped into [S] and the operations can be chained, where the methods that
// change the element type (such as [S.Map] and [S.FilterMap]) declare their own
// type parameters:
//
//	Wrap([]int{1, 2, 3}).
//	    Filter(func(i int) bool { return i%2 == 1 }).
//	    Map(strconv.Itoa).
//	    ToMapValues(func(s string) string { return s })
//
// # Element constraints
//
// A method can not add constraints on the type parameters of its receiver, so
// the operations are spread over several wrapper types, each of them declares
// the constraint it needs on its own type parameter:
//
//   - [S]: no constraint on the element (the default choice).
//   - [C]: the element must be comparable, e.g. [C.Contains], [C.Uniq].
//   - [O]: the element must be ordered, e.g. [O.Sort], [O.Max].
//   - [N]: the element must be a number, e.g. [N.Sum], [N.Avg].
//
// The constrained wrappers are converted back to [S] by [C.S], [O.S] and [N.S],
// and [O] (and [N]) also provide [O.C] since ordered (and number) elements are
// comparable as well.
package gslice

import (
	"github.com/bytedance/gg/collection/tuple"
	"github.com/bytedance/gg/goption"
	"github.com/bytedance/gg/gresult"
	"github.com/bytedance/gg/internal/constraints"
)

// S is a wrapper of slice []T, it provides the slice operations as methods.
//
// Use [Wrap] to wrap an existing slice into S, [S.Unwrap] to unwrap it.
//
// 🚀 EXAMPLE:
//
//	Wrap([]int{1, 2, 3}).Map(strconv.Itoa).Unwrap()  ⏩ ["1", "2", "3"]
type S[T any] []T

// Wrap wraps an existing slice s into [S], so that the slice operations can be
// chained as methods.
//
// 🚀 EXAMPLE:
//
//	Wrap([]int{1, 2, 3}).Len()  ⏩ 3
func Wrap[T any](s []T) S[T] {
	return S[T](s)
}

// New creates a slice with the given elements.
//
// It is the [S] variant of function [Of], which is named differently because
// [Of] is taken by the package-level function returning a plain slice.
//
// 💡 HINT: The slice produced by the package-level constructors (such as [Of],
// [Range] and [Repeat]) is wrapped by [Wrap]:
//
//	Wrap(Range(1, 5)).Map(strconv.Itoa)
//
// 🚀 EXAMPLE:
//
//	New(1, 2, 3).Len()  ⏩ 3
func New[T any](vs ...T) S[T] {
	return S[T](vs)
}

// Unwrap returns the wrapped slice of S.
//
// 💡 AKA: Get, Value
func (s S[T]) Unwrap() []T {
	return []T(s)
}

// Len returns the length of the slice.
func (s S[T]) Len() int {
	return Len(s)
}

// Clone returns a copy of the slice.
//
// 💡 HINT: The element is copied using assignment (=), so this is a shallow clone.
func (s S[T]) Clone() S[T] {
	return Clone(s)
}

// CloneBy is variant of [S.Clone], element is copied using function f.
func (s S[T]) CloneBy(f func(T) T) S[T] {
	return CloneBy(s, f)
}

// Map applies function f to each element of the slice and returns a new slice
// with the results.
//
// 🚀 EXAMPLE:
//
//	Wrap([]int{1, 2, 3}).Map(strconv.Itoa).Unwrap()  ⏩ ["1", "2", "3"]
func (s S[T]) Map[R any](f func(T) R) S[R] {
	return Map(s, f)
}

// MapIndexed is variant of [S.Map], accepts the index of the element.
func (s S[T]) MapIndexed[R any](f func(T, int) R) S[R] {
	return MapIndexed(s, f)
}

// TryMap is variant of [S.Map], accepts a function that may return an error.
// The error is wrapped into [gresult.R].
func (s S[T]) TryMap[R any](f func(T) (R, error)) gresult.R[S[R]] {
	r := TryMap(s, f)
	if r.IsErr() {
		return gresult.Err[S[R]](r.Err())
	}
	return gresult.OK(S[R](r.Value()))
}

// Filter returns a new slice with the elements that f returns true for.
//
// 🚀 EXAMPLE:
//
//	Wrap([]int{1, 2, 3}).Filter(func(i int) bool { return i%2 == 1 }).Unwrap()  ⏩ [1, 3]
func (s S[T]) Filter(f func(T) bool) S[T] {
	return Filter(s, f)
}

// FilterMap is a combination of [S.Filter] and [S.Map], the function f returns
// the mapped value and whether it should be kept.
func (s S[T]) FilterMap[R any](f func(T) (R, bool)) S[R] {
	return FilterMap(s, f)
}

// TryFilterMap is variant of [S.FilterMap], the function f returns an error
// instead of a bool, the element is dropped when an error is returned.
func (s S[T]) TryFilterMap[R any](f func(T) (R, error)) S[R] {
	return TryFilterMap(s, f)
}

// Reject is negation of [S.Filter].
func (s S[T]) Reject(f func(T) bool) S[T] {
	return Reject(s, f)
}

// Partition splits the slice into two slices: the elements that f returns true
// for, and the rest.
//
// 🚀 EXAMPLE:
//
//	Wrap([]int{1, 2, 3, 4}).Partition(func(i int) bool { return i%2 == 0 })
//	// [2 4] [1 3]
func (s S[T]) Partition(f func(T) bool) (S[T], S[T]) {
	matched, unmatched := Partition(s, f)
	return matched, unmatched
}

// Reduce reduces the slice to a single value by function f.
// It returns nothing when the slice is empty.
//
// 🚀 EXAMPLE:
//
//	Wrap([]int{1, 2, 3}).Reduce(gvalue.Add[int]).Value()  ⏩ 6
func (s S[T]) Reduce(f func(T, T) T) goption.O[T] {
	return Reduce(s, f)
}

// Fold reduces the slice to a single value by function f, starting from init.
//
// 🚀 EXAMPLE:
//
//	Wrap([]int{1, 2, 3}).Fold(func(acc string, i int) string { return acc + strconv.Itoa(i) }, "")
//	// "123"
func (s S[T]) Fold[R any](f func(R, T) R, init R) R {
	return Fold(s, f, init)
}

// Any returns whether any element of the slice satisfies f.
func (s S[T]) Any(f func(T) bool) bool {
	return Any(s, f)
}

// All returns whether all elements of the slice satisfy f.
func (s S[T]) All(f func(T) bool) bool {
	return All(s, f)
}

// Find returns the first element that satisfies f.
func (s S[T]) Find(f func(T) bool) goption.O[T] {
	return Find(s, f)
}

// FindRev is variant of [S.Find], it searches from the end.
func (s S[T]) FindRev(f func(T) bool) goption.O[T] {
	return FindRev(s, f)
}

// First returns the first element of the slice.
func (s S[T]) First() goption.O[T] {
	return First(s)
}

// Last returns the last element of the slice.
func (s S[T]) Last() goption.O[T] {
	return Last(s)
}

// Get returns the element at index n, negative index counts from the end.
//
// 🚀 EXAMPLE:
//
//	Wrap([]int{1, 2, 3}).Get(1).Value()   ⏩ 2
//	Wrap([]int{1, 2, 3}).Get(-1).Value()  ⏩ 3
func (s S[T]) Get[I constraints.Integer](n I) goption.O[T] {
	return Get(s, n)
}

// IndexBy returns the index of the first element that satisfies f.
func (s S[T]) IndexBy(f func(T) bool) goption.O[int] {
	return IndexBy(s, f)
}

// IndexRevBy is variant of [S.IndexBy], it searches from the end.
func (s S[T]) IndexRevBy(f func(T) bool) goption.O[int] {
	return IndexRevBy(s, f)
}

// CountBy returns the number of elements that satisfy f.
func (s S[T]) CountBy(f func(T) bool) int {
	return CountBy(s, f)
}

// Chunk splits the slice into chunks of the given size.
//
// 🚀 EXAMPLE:
//
//	Wrap([]int{1, 2, 3, 4, 5}).Chunk(2)  ⏩ [[1, 2], [3, 4], [5]]
func (s S[T]) Chunk(size int) []S[T] {
	return Chunk(s, size)
}

// ChunkClone is variant of [S.Chunk], the chunks are copies of the elements.
func (s S[T]) ChunkClone(size int) []S[T] {
	return ChunkClone(s, size)
}

// Divide splits the slice into n chunks as evenly as possible.
//
// 🚀 EXAMPLE:
//
//	Wrap([]int{1, 2, 3, 4, 5}).Divide(2)  ⏩ [[1, 2, 3], [4, 5]]
func (s S[T]) Divide(n int) []S[T] {
	return Divide(s, n)
}

// DivideClone is variant of [S.Divide], the chunks are copies of the elements.
func (s S[T]) DivideClone(n int) []S[T] {
	return DivideClone(s, n)
}

// GroupBy groups the elements into a map by the key returned by f.
//
// 🚀 EXAMPLE:
//
//	Wrap([]int{1, 2, 3}).GroupBy(func(i int) string { return strconv.Itoa(i%2) })
//	// {"0": [2], "1": [1, 3]}
func (s S[T]) GroupBy[K comparable](f func(T) K) map[K]S[T] {
	return GroupBy(s, f)
}

// UniqBy returns a new slice with duplicate elements removed,
// the key of an element is returned by f.
func (s S[T]) UniqBy[K comparable](f func(T) K) S[T] {
	return UniqBy(s, f)
}

// DupBy returns the duplicate elements of the slice,
// the key of an element is returned by f.
func (s S[T]) DupBy[K comparable](f func(T) K) S[T] {
	return DupBy(s, f)
}

// CountValuesBy counts the elements by the key returned by f.
func (s S[T]) CountValuesBy[K comparable](f func(T) K) map[K]int {
	return CountValuesBy(s, f)
}

// MaxBy returns the maximum element of the slice, compared by function less.
func (s S[T]) MaxBy(less func(T, T) bool) goption.O[T] {
	return MaxBy(s, less)
}

// MinBy returns the minimum element of the slice, compared by function less.
func (s S[T]) MinBy(less func(T, T) bool) goption.O[T] {
	return MinBy(s, less)
}

// MinMaxBy returns the minimum and maximum element of the slice,
// compared by function less.
func (s S[T]) MinMaxBy(less func(T, T) bool) goption.O[tuple.T2[T, T]] {
	return MinMaxBy(s, less)
}

// FlatMap applies function f to each element of the slice and flattens the
// result.
func (s S[T]) FlatMap[R any](f func(T) []R) S[R] {
	return FlatMap(s, f)
}

// ForEach calls function f for each element of the slice.
func (s S[T]) ForEach(f func(T)) {
	ForEach(s, f)
}

// ForEachIndexed is variant of [S.ForEach], accepts the index of the element.
func (s S[T]) ForEachIndexed(f func(int, T)) {
	ForEachIndexed(s, f)
}

// EqualBy returns whether the two slices are equal, compared by function eq.
func (s S[T]) EqualBy(s2 []T, eq func(T, T) bool) bool {
	return EqualBy(s, s2, eq)
}

// Reverse reverses the order of the elements in place and returns the slice
// itself for chaining.
//
// 🚀 EXAMPLE:
//
//	Wrap([]int{1, 2, 3}).Reverse().Unwrap()  ⏩ [3, 2, 1]
func (s S[T]) Reverse() S[T] {
	Reverse(s)
	return s
}

// ReverseClone is variant of [S.Reverse], it returns a new slice and keeps the
// original one untouched.
func (s S[T]) ReverseClone() S[T] {
	return ReverseClone(s)
}

// SortBy sorts the elements in place by function less and returns the slice
// itself for chaining.
func (s S[T]) SortBy(less func(T, T) bool) S[T] {
	SortBy(s, less)
	return s
}

// StableSortBy is variant of [S.SortBy], it keeps the original order of equal
// elements.
func (s S[T]) StableSortBy(less func(T, T) bool) S[T] {
	StableSortBy(s, less)
	return s
}

// SortCloneBy is variant of [S.SortBy], it returns a new slice and keeps the
// original one untouched.
func (s S[T]) SortCloneBy(less func(T, T) bool) S[T] {
	return SortCloneBy(s, less)
}

// Shuffle pseudo-randomizes the order of the elements in place and returns the
// slice itself for chaining.
func (s S[T]) Shuffle() S[T] {
	Shuffle(s)
	return s
}

// ShuffleClone is variant of [S.Shuffle], it returns a new slice and keeps the
// original one untouched.
func (s S[T]) ShuffleClone() S[T] {
	return ShuffleClone(s)
}

// Take returns the first n elements of the slice, negative n takes from the end.
//
// 🚀 EXAMPLE:
//
//	Wrap([]int{1, 2, 3}).Take(2).Unwrap()   ⏩ [1, 2]
//	Wrap([]int{1, 2, 3}).Take(-2).Unwrap()  ⏩ [2, 3]
func (s S[T]) Take[I constraints.Integer](n I) S[T] {
	return Take(s, n)
}

// TakeClone is variant of [S.Take], it returns a new slice.
func (s S[T]) TakeClone[I constraints.Integer](n I) S[T] {
	return TakeClone(s, n)
}

// Slice returns the elements between start (inclusive) and end (exclusive),
// negative index counts from the end.
//
// 🚀 EXAMPLE:
//
//	Wrap([]int{1, 2, 3, 4}).Slice(1, 3).Unwrap()  ⏩ [2, 3]
func (s S[T]) Slice[I constraints.Integer](start, end I) S[T] {
	return Slice(s, start, end)
}

// SliceClone is variant of [S.Slice], it returns a new slice.
func (s S[T]) SliceClone[I constraints.Integer](start, end I) S[T] {
	return SliceClone(s, start, end)
}

// Drop removes the first n elements of the slice, negative n drops from the end.
func (s S[T]) Drop(n int) S[T] {
	return Drop(s, n)
}

// DropClone is variant of [S.Drop], it returns a new slice.
func (s S[T]) DropClone(n int) S[T] {
	return DropClone(s, n)
}

// RemoveIndex removes the element at the given index.
func (s S[T]) RemoveIndex[I constraints.Integer](index I) S[T] {
	return RemoveIndex(s, index)
}

// Insert inserts the values vs at the given position.
func (s S[T]) Insert[I constraints.Integer](pos I, vs ...T) S[T] {
	return Insert(s, pos, vs...)
}

// Concat concatenates the slice with the given ones and returns a new slice.
//
// 🚀 EXAMPLE:
//
//	Wrap([]int{1, 2}).Concat([]int{3}, []int{4}).Unwrap()  ⏩ [1, 2, 3, 4]
func (s S[T]) Concat(ss ...[]T) S[T] {
	return Concat(append([][]T{s}, ss...)...)
}

// Merge is an alias of [S.Concat].
func (s S[T]) Merge(ss ...[]T) S[T] {
	return Merge(append([][]T{s}, ss...)...)
}

// SumBy sums the values returned by function f.
func (s S[T]) SumBy[N constraints.Number](f func(T) N) N {
	return SumBy(s, f)
}

// AvgBy returns the average of the values returned by function f.
func (s S[T]) AvgBy[N constraints.Number](f func(T) N) float64 {
	return AvgBy(s, f)
}

// ToMap converts the slice to a map, the key and value of an element are
// returned by f.
//
// 🚀 EXAMPLE:
//
//	Wrap([]int{1, 2}).ToMap(func(i int) (string, int) { return strconv.Itoa(i), i })
//	// {"1": 1, "2": 2}
func (s S[T]) ToMap[K comparable, V any](f func(T) (K, V)) map[K]V {
	return ToMap(s, f)
}

// ToMapValues converts the slice to a map, the key of an element is returned by f.
func (s S[T]) ToMapValues[K comparable](f func(T) K) map[K]T {
	return ToMapValues(s, f)
}

// TypeAssert converts the elements from type T to type R by [type assertion].
//
// ⚠️ WARNING: It may ❌PANIC❌ when type assertion failed.
//
// [type assertion]: https://go.dev/tour/methods/15
func (s S[T]) TypeAssert[R any]() S[R] {
	return TypeAssert[R](s)
}
