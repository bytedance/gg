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
	"github.com/bytedance/gg/collection/tuple"
	"github.com/bytedance/gg/goption"
	"github.com/bytedance/gg/internal/constraints"
)

// C is a wrapper of slice []T whose element is comparable,
// it provides the operations that need comparable elements as methods.
//
// Use [WrapC] to wrap an existing slice into C, [C.S] to convert it to [S] so
// that the operations without element constraint are available.
//
// 🚀 EXAMPLE:
//
//	WrapC([]int{1, 1, 2}).Uniq().S().Map(strconv.Itoa).Unwrap()  ⏩ ["1", "2"]
type C[T comparable] []T

// WrapC wraps an existing slice s into [C], so that the operations that need
// comparable elements can be chained as methods.
//
// 🚀 EXAMPLE:
//
//	WrapC([]int{1, 2, 3}).Contains(2)  ⏩ true
func WrapC[T comparable](s []T) C[T] {
	return C[T](s)
}

// Unwrap returns the wrapped slice of C.
func (c C[T]) Unwrap() []T {
	return []T(c)
}

// S converts C to [S], so that the operations without element constraint are
// available.
//
// 🚀 EXAMPLE:
//
//	WrapC([]int{1, 1, 2}).Uniq().S().Map(strconv.Itoa).Unwrap()  ⏩ ["1", "2"]
func (c C[T]) S() S[T] {
	return S[T](c)
}

// Len returns the length of the slice.
func (c C[T]) Len() int {
	return Len(c)
}

// Contains returns whether the slice contains element v.
//
// 🚀 EXAMPLE:
//
//	WrapC([]int{1, 2, 3}).Contains(2)  ⏩ true
func (c C[T]) Contains(v T) bool {
	return Contains(c, v)
}

// ContainsAny returns whether the slice contains any of the given elements.
//
// 🚀 EXAMPLE:
//
//	WrapC([]int{1, 2, 3}).ContainsAny(2, 6)  ⏩ true
func (c C[T]) ContainsAny(vs ...T) bool {
	return ContainsAny(c, vs...)
}

// ContainsAll returns whether the slice contains all of the given elements.
//
// 🚀 EXAMPLE:
//
//	WrapC([]int{1, 2, 3}).ContainsAll(2, 6)  ⏩ false
func (c C[T]) ContainsAll(vs ...T) bool {
	return ContainsAll(c, vs...)
}

// Index returns the index of element e, it returns nothing when the element is
// not found.
//
// 🚀 EXAMPLE:
//
//	WrapC([]int{1, 2, 3}).Index(3).Value()  ⏩ 2
func (c C[T]) Index(e T) goption.O[int] {
	return Index(c, e)
}

// IndexRev is variant of [C.Index], it searches from the end.
func (c C[T]) IndexRev(e T) goption.O[int] {
	return IndexRev(c, e)
}

// Remove removes all occurrences of element v.
//
// 🚀 EXAMPLE:
//
//	WrapC([]int{1, 2, 2}).Remove(2).Unwrap()  ⏩ [1]
func (c C[T]) Remove(v T) C[T] {
	return Remove(c, v)
}

// Uniq removes the duplicate elements and keeps the original order.
//
// 🚀 EXAMPLE:
//
//	WrapC([]int{1, 1, 2}).Uniq().Unwrap()  ⏩ [1, 2]
func (c C[T]) Uniq() C[T] {
	return Uniq(c)
}

// Dup returns the duplicate elements of the slice.
//
// 🚀 EXAMPLE:
//
//	WrapC([]int{1, 1, 2}).Dup().Unwrap()  ⏩ [1]
func (c C[T]) Dup() C[T] {
	return Dup(c)
}

// Compact removes the zero-value elements of the slice.
func (c C[T]) Compact() C[T] {
	return Compact(c)
}

// Union returns the union of the slice and the given ones.
//
// 🚀 EXAMPLE:
//
//	WrapC([]int{1, 2}).Union(WrapC([]int{2, 3})).Unwrap()  ⏩ [1, 2, 3]
func (c C[T]) Union(ss ...C[T]) C[T] {
	return Union(append([]C[T]{c}, ss...)...)
}

// Diff returns the elements of the slice that are not in the given ones.
//
// 🚀 EXAMPLE:
//
//	WrapC([]int{1, 2}).Diff(WrapC([]int{2, 3})).Unwrap()  ⏩ [1]
func (c C[T]) Diff(againsts ...C[T]) C[T] {
	return Diff(c, againsts...)
}

// Intersect returns the elements that are both in the slice and in all of the
// given ones.
//
// 🚀 EXAMPLE:
//
//	WrapC([]int{1, 2}).Intersect(WrapC([]int{2, 3})).Unwrap()  ⏩ [2]
func (c C[T]) Intersect(others ...C[T]) C[T] {
	return Intersect(append([]C[T]{c}, others...)...)
}

// Equal returns whether the two slices are equal.
func (c C[T]) Equal(s2 []T) bool {
	return Equal(c, s2)
}

// Count returns the number of occurrences of element v.
func (c C[T]) Count(v T) int {
	return Count(c, v)
}

// CountValues counts the occurrences of each element.
func (c C[T]) CountValues() map[T]int {
	return CountValues(c)
}

// ToBoolMap converts the slice to a map with the elements as keys and true as
// value.
//
// 🚀 EXAMPLE:
//
//	WrapC([]int{1, 2, 2}).ToBoolMap()  ⏩ {1: true, 2: true}
func (c C[T]) ToBoolMap() map[T]bool {
	return ToBoolMap(c)
}

// O is a wrapper of slice []T whose element is ordered,
// it provides the operations that need ordered elements as methods.
//
// Use [WrapO] to wrap an existing slice into O, [O.S] to convert it to [S] and
// [O.C] to convert it to [C] (ordered elements are comparable).
//
// 🚀 EXAMPLE:
//
//	WrapO([]int{3, 1, 2}).Sort().Unwrap()  ⏩ [1, 2, 3]
//	WrapO([]int{3, 1, 2}).Max().Value()    ⏩ 3
type O[T constraints.Ordered] []T

// WrapO wraps an existing slice s into [O], so that the operations that need
// ordered elements can be chained as methods.
//
// 🚀 EXAMPLE:
//
//	WrapO([]int{3, 1, 2}).Sort().Unwrap()  ⏩ [1, 2, 3]
func WrapO[T constraints.Ordered](s []T) O[T] {
	return O[T](s)
}

// Unwrap returns the wrapped slice of O.
func (o O[T]) Unwrap() []T {
	return []T(o)
}

// S converts O to [S], so that the operations without element constraint are
// available.
func (o O[T]) S() S[T] {
	return S[T](o)
}

// C converts O to [C], since ordered elements are comparable.
func (o O[T]) C() C[T] {
	return C[T](o)
}

// Len returns the length of the slice.
func (o O[T]) Len() int {
	return Len(o)
}

// Max returns the maximum element of the slice.
//
// 🚀 EXAMPLE:
//
//	WrapO([]int{1, 2, 3}).Max().Value()  ⏩ 3
func (o O[T]) Max() goption.O[T] {
	return Max(o)
}

// Min returns the minimum element of the slice.
//
// 🚀 EXAMPLE:
//
//	WrapO([]int{1, 2, 3}).Min().Value()  ⏩ 1
func (o O[T]) Min() goption.O[T] {
	return Min(o)
}

// MinMax returns the minimum and maximum element of the slice.
//
// 🚀 EXAMPLE:
//
//	WrapO([]int{1, 2, 3}).MinMax().Value().Values()  ⏩ 1 3
func (o O[T]) MinMax() goption.O[tuple.T2[T, T]] {
	return MinMax(o)
}

// Sort sorts the elements in ascending order in place and returns the slice
// itself for chaining.
//
// 🚀 EXAMPLE:
//
//	WrapO([]int{3, 1, 2}).Sort().Unwrap()  ⏩ [1, 2, 3]
func (o O[T]) Sort() O[T] {
	Sort(o)
	return o
}

// SortClone is variant of [O.Sort], it returns a new slice and keeps the
// original one untouched.
func (o O[T]) SortClone() O[T] {
	return SortClone(o)
}

// PartialSort sorts the slice so that the first k elements are the smallest ones
// in ascending order, it returns the slice itself for chaining.
func (o O[T]) PartialSort(k int) O[T] {
	PartialSort(o, k)
	return o
}

// PartialSortBy is variant of [O.PartialSort], the order is defined by function
// less.
func (o O[T]) PartialSortBy(k int, less func(T, T) bool) O[T] {
	PartialSortBy(o, k, less)
	return o
}

// N is a wrapper of slice []T whose element is a number,
// it provides the arithmetic operations as methods.
//
// Use [WrapN] to wrap an existing slice into N, [N.S] to convert it to [S] and
// [N.C] to convert it to [C] (numbers are comparable).
//
// 🚀 EXAMPLE:
//
//	WrapN([]int{1, 2, 3}).Sum()  ⏩ 6
//	WrapN([]int{1, 2, 3}).Avg()  ⏩ 2
type N[T constraints.Number] []T

// WrapN wraps an existing slice s into [N], so that the arithmetic operations
// can be chained as methods.
//
// 🚀 EXAMPLE:
//
//	WrapN([]int{1, 2, 3}).Sum()  ⏩ 6
func WrapN[T constraints.Number](s []T) N[T] {
	return N[T](s)
}

// Unwrap returns the wrapped slice of N.
func (n N[T]) Unwrap() []T {
	return []T(n)
}

// S converts N to [S], so that the operations without element constraint are
// available.
func (n N[T]) S() S[T] {
	return S[T](n)
}

// C converts N to [C], since numbers are comparable.
func (n N[T]) C() C[T] {
	return C[T](n)
}

// O converts N to [O], since numbers are ordered.
func (n N[T]) O() O[T] {
	return O[T](n)
}

// Len returns the length of the slice.
func (n N[T]) Len() int {
	return Len(n)
}

// Sum sums up the elements of the slice.
//
// 🚀 EXAMPLE:
//
//	WrapN([]int{1, 2, 3}).Sum()  ⏩ 6
func (n N[T]) Sum() T {
	return Sum(n)
}

// Avg returns the average of the elements of the slice.
//
// 🚀 EXAMPLE:
//
//	WrapN([]int{1, 2, 3}).Avg()  ⏩ 2
func (n N[T]) Avg() float64 {
	return Avg(n)
}
