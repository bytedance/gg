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

// This file contains the method form of the map operations.
//
// A map[K]V can not be a method receiver, so the operations used to be
// package-level functions only. With Go 1.27 generic methods, a map can be
// wrapped into [M] and the operations can be chained, where the methods that
// change the key/value type (such as [M.Map] and [M.MapValues]) declare their
// own type parameters:
//
//	Wrap(map[int]int{1: 2}).MapValues(strconv.Itoa).Keys()
//
// # Element constraints
//
// A method can not add constraints on the type parameters of its receiver, so
// the operations are spread over several wrapper types, each of them declares
// the constraint it needs on its own type parameters:
//
//   - [M]: no constraint on the value (the default choice).
//   - [MK]: the key must be ordered, e.g. [MK.OrderedKeys], [MK.ToOrderedSlice].
//   - [MV]: the value must be comparable, e.g. [MV.Invert], [MV.Equal].
//   - [MO]: the value must be ordered, e.g. [MO.Max], [MO.Min].
//   - [MN]: the value must be a number, e.g. [MN.Sum], [MN.Avg].
//
// The constrained wrappers are converted back to [M] by [MK.M], [MV.M],
// [MO.M] and [MN.M].
package gmap

import (
	"github.com/bytedance/gg/collection/tuple"
	"github.com/bytedance/gg/goption"
	"github.com/bytedance/gg/gresult"
	"github.com/bytedance/gg/internal/constraints"
)

// M is a wrapper of map[K]V, it provides the map operations as methods.
//
// Use [Wrap] to wrap an existing map into M, [M.Unwrap] to unwrap it.
//
// 🚀 EXAMPLE:
//
//	Wrap(map[int]int{1: 2}).MapValues(strconv.Itoa).Unwrap()  ⏩ {1: "2"}
type M[K comparable, V any] map[K]V

// Wrap wraps an existing map m into [M], so that the map operations can be
// chained as methods.
//
// 🚀 EXAMPLE:
//
//	Wrap(map[int]int{1: 2}).Len()  ⏩ 1
func Wrap[K comparable, V any](m map[K]V) M[K, V] {
	return M[K, V](m)
}

// New creates an empty map with the given capacity.
//
// 🚀 EXAMPLE:
//
//	New[string, int](4)["a"] = 1
func New[K comparable, V any](cap int) M[K, V] {
	return M[K, V](make(map[K]V, cap))
}

// Unwrap returns the wrapped map of M.
//
// 💡 AKA: Get, Value
func (m M[K, V]) Unwrap() map[K]V {
	return map[K]V(m)
}

// Len returns the number of the items of the map.
func (m M[K, V]) Len() int {
	return Len(m)
}

// Clone returns a copy of the map.
func (m M[K, V]) Clone() M[K, V] {
	return Clone(m)
}

// CloneBy is variant of [M.Clone], the value is copied using function f.
func (m M[K, V]) CloneBy(f func(V) V) M[K, V] {
	return CloneBy(m, f)
}

// Map applies function f to each item of the map and returns a new map.
//
// 🚀 EXAMPLE:
//
//	Wrap(map[int]int{1: 2}).Map(func(k, v int) (string, string) {
//	    return strconv.Itoa(k), strconv.Itoa(v)
//	}).Unwrap()  ⏩ {"1": "2"}
func (m M[K, V]) Map[K2 comparable, V2 any](f func(K, V) (K2, V2)) M[K2, V2] {
	return Map(m, f)
}

// TryMap is variant of [M.Map], accepts a function that may return an error.
// The error is wrapped into [gresult.R].
func (m M[K, V]) TryMap[K2 comparable, V2 any](f func(K, V) (K2, V2, error)) gresult.R[M[K2, V2]] {
	r := TryMap(m, f)
	if r.IsErr() {
		return gresult.Err[M[K2, V2]](r.Err())
	}
	return gresult.OK(M[K2, V2](r.Value()))
}

// MapKeys applies function f to each key of the map, the values are kept.
func (m M[K, V]) MapKeys[K2 comparable](f func(K) K2) M[K2, V] {
	return MapKeys(m, f)
}

// TryMapKeys is variant of [M.MapKeys], accepts a function that may return an
// error. The error is wrapped into [gresult.R].
func (m M[K, V]) TryMapKeys[K2 comparable](f func(K) (K2, error)) gresult.R[M[K2, V]] {
	r := TryMapKeys(m, f)
	if r.IsErr() {
		return gresult.Err[M[K2, V]](r.Err())
	}
	return gresult.OK(M[K2, V](r.Value()))
}

// MapValues applies function f to each value of the map, the keys are kept.
//
// 🚀 EXAMPLE:
//
//	Wrap(map[int]int{1: 2}).MapValues(strconv.Itoa).Unwrap()  ⏩ {1: "2"}
func (m M[K, V]) MapValues[V2 any](f func(V) V2) M[K, V2] {
	return MapValues(m, f)
}

// TryMapValues is variant of [M.MapValues], accepts a function that may return
// an error. The error is wrapped into [gresult.R].
func (m M[K, V]) TryMapValues[V2 any](f func(V) (V2, error)) gresult.R[M[K, V2]] {
	r := TryMapValues(m, f)
	if r.IsErr() {
		return gresult.Err[M[K, V2]](r.Err())
	}
	return gresult.OK(M[K, V2](r.Value()))
}

// FilterMap is a combination of [M.Filter] and [M.Map], the function f returns
// the mapped item and whether it should be kept.
func (m M[K, V]) FilterMap[K2 comparable, V2 any](f func(K, V) (K2, V2, bool)) M[K2, V2] {
	return FilterMap(m, f)
}

// TryFilterMap is variant of [M.FilterMap], the function f returns an error
// instead of a bool, the item is dropped when an error is returned.
func (m M[K, V]) TryFilterMap[K2 comparable, V2 any](f func(K, V) (K2, V2, error)) M[K2, V2] {
	return TryFilterMap(m, f)
}

// FilterMapKeys is a combination of [M.Filter] and [M.MapKeys].
func (m M[K, V]) FilterMapKeys[K2 comparable](f func(K) (K2, bool)) M[K2, V] {
	return FilterMapKeys(m, f)
}

// TryFilterMapKeys is variant of [M.FilterMapKeys], accepts a function that may
// return an error.
func (m M[K, V]) TryFilterMapKeys[K2 comparable](f func(K) (K2, error)) M[K2, V] {
	return TryFilterMapKeys(m, f)
}

// FilterMapValues is a combination of [M.Filter] and [M.MapValues].
func (m M[K, V]) FilterMapValues[V2 any](f func(V) (V2, bool)) M[K, V2] {
	return FilterMapValues(m, f)
}

// TryFilterMapValues is variant of [M.FilterMapValues], accepts a function that
// may return an error.
func (m M[K, V]) TryFilterMapValues[V2 any](f func(V) (V2, error)) M[K, V2] {
	return TryFilterMapValues(m, f)
}

// Filter returns a new map with the items that f returns true for.
//
// 🚀 EXAMPLE:
//
//	Wrap(map[int]int{1: 2, 2: 3}).Filter(func(k, v int) bool { return k+v > 3 }).Unwrap()
//	// {2: 3}
func (m M[K, V]) Filter(f func(K, V) bool) M[K, V] {
	return Filter(m, f)
}

// FilterKeys is variant of [M.Filter], the predicate accepts the key only.
func (m M[K, V]) FilterKeys(f func(K) bool) M[K, V] {
	return FilterKeys(m, f)
}

// FilterByKeys returns a new map with the items whose key is in the given keys.
func (m M[K, V]) FilterByKeys(keys ...K) M[K, V] {
	return FilterByKeys(m, keys...)
}

// FilterValues is variant of [M.Filter], the predicate accepts the value only.
func (m M[K, V]) FilterValues(f func(V) bool) M[K, V] {
	return FilterValues(m, f)
}

// Reject is negation of [M.Filter].
func (m M[K, V]) Reject(f func(K, V) bool) M[K, V] {
	return Reject(m, f)
}

// RejectKeys is negation of [M.FilterKeys].
func (m M[K, V]) RejectKeys(f func(K) bool) M[K, V] {
	return RejectKeys(m, f)
}

// RejectByKeys is negation of [M.FilterByKeys].
func (m M[K, V]) RejectByKeys(keys ...K) M[K, V] {
	return RejectByKeys(m, keys...)
}

// RejectValues is negation of [M.FilterValues].
func (m M[K, V]) RejectValues(f func(V) bool) M[K, V] {
	return RejectValues(m, f)
}

// Keys returns the keys of the map.
func (m M[K, V]) Keys() []K {
	return Keys(m)
}

// Values returns the values of the map.
func (m M[K, V]) Values() []V {
	return Values(m)
}

// Items returns the items (key-value pairs) of the map.
func (m M[K, V]) Items() tuple.S2[K, V] {
	return Items(m)
}

// Merge merges the given maps into a new map, the value of the later map wins
// on conflict.
func (m M[K, V]) Merge(ms ...M[K, V]) M[K, V] {
	return Merge(append([]M[K, V]{m}, ms...)...)
}

// Union returns the union of the map and the given ones, the value of the later
// map wins on conflict.
//
// 🚀 EXAMPLE:
//
//	Wrap(map[int]int{1: 2}).Union(Wrap(map[int]int{1: 3, 2: 4})).Unwrap()  ⏩ {1: 3, 2: 4}
func (m M[K, V]) Union(ms ...M[K, V]) M[K, V] {
	return Union(append([]M[K, V]{m}, ms...)...)
}

// UnionBy is variant of [M.Union], the conflict is resolved by function
// onConflict.
func (m M[K, V]) UnionBy(ms []M[K, V], onConflict ConflictFunc[K, V]) M[K, V] {
	return UnionBy(append([]M[K, V]{m}, ms...), onConflict)
}

// Diff returns the items of the map whose key is not in the given ones.
//
// 🚀 EXAMPLE:
//
//	Wrap(map[int]int{1: 2, 2: 3}).Diff(Wrap(map[int]int{2: 4})).Unwrap()  ⏩ {1: 2}
func (m M[K, V]) Diff(againsts ...M[K, V]) M[K, V] {
	return Diff(m, againsts...)
}

// Intersect returns the items whose key is both in the map and in all of the
// given ones.
//
// 🚀 EXAMPLE:
//
//	Wrap(map[int]int{1: 2, 2: 3}).Intersect(Wrap(map[int]int{2: 4})).Unwrap()  ⏩ {2: 4}
func (m M[K, V]) Intersect(others ...M[K, V]) M[K, V] {
	return Intersect(append([]M[K, V]{m}, others...)...)
}

// IntersectBy is variant of [M.Intersect], the conflict is resolved by function
// onConflict.
func (m M[K, V]) IntersectBy(others []M[K, V], onConflict ConflictFunc[K, V]) M[K, V] {
	return IntersectBy(append([]M[K, V]{m}, others...), onConflict)
}

// Load returns the value of the given key.
//
// 🚀 EXAMPLE:
//
//	Wrap(map[int]int{1: 2}).Load(1).Value()  ⏩ 2
func (m M[K, V]) Load(k K) goption.O[V] {
	return Load(m, k)
}

// LoadOrStore returns the value of the given key, stores and returns defaultV
// when the key is not present.
func (m M[K, V]) LoadOrStore(k K, defaultV V) (V, bool) {
	return LoadOrStore(m, k, defaultV)
}

// LoadOrStoreLazy is variant of [M.LoadOrStore], the default value is lazily
// evaluated by function f.
func (m M[K, V]) LoadOrStoreLazy(k K, f func() V) (V, bool) {
	return LoadOrStoreLazy(m, k, f)
}

// LoadAndDelete returns the value of the given key and deletes it.
func (m M[K, V]) LoadAndDelete(k K) goption.O[V] {
	return LoadAndDelete(m, k)
}

// LoadBy returns the value of the first item that f returns true for.
func (m M[K, V]) LoadBy(f func(K, V) bool) goption.O[V] {
	return LoadBy(m, f)
}

// LoadKeyBy returns the key of the first item that f returns true for.
func (m M[K, V]) LoadKeyBy(f func(K, V) bool) goption.O[K] {
	return LoadKeyBy(m, f)
}

// LoadItemBy returns the first item that f returns true for.
func (m M[K, V]) LoadItemBy(f func(K, V) bool) goption.O[tuple.T2[K, V]] {
	return LoadItemBy(m, f)
}

// LoadAll returns the values of the given keys, it returns nothing when any of
// the keys is missing.
func (m M[K, V]) LoadAll(ks ...K) []V {
	return LoadAll(m, ks...)
}

// LoadAny returns the value of the first present key.
func (m M[K, V]) LoadAny(ks ...K) goption.O[V] {
	return LoadAny(m, ks...)
}

// LoadSome returns the values of the present keys.
func (m M[K, V]) LoadSome(ks ...K) []V {
	return LoadSome(m, ks...)
}

// Contains returns whether the map contains the given key.
func (m M[K, V]) Contains(k K) bool {
	return Contains(m, k)
}

// ContainsAny returns whether the map contains any of the given keys.
func (m M[K, V]) ContainsAny(ks ...K) bool {
	return ContainsAny(m, ks...)
}

// ContainsAll returns whether the map contains all of the given keys.
func (m M[K, V]) ContainsAll(ks ...K) bool {
	return ContainsAll(m, ks...)
}

// EqualBy returns whether the two maps are equal, the values are compared by
// function eq.
func (m M[K, V]) EqualBy(m2 map[K]V, eq func(V, V) bool) bool {
	return EqualBy(m, m2, eq)
}

// MaxBy returns the maximum value of the map, compared by function less.
func (m M[K, V]) MaxBy(less func(V, V) bool) goption.O[V] {
	return MaxBy(m, less)
}

// MinBy returns the minimum value of the map, compared by function less.
func (m M[K, V]) MinBy(less func(V, V) bool) goption.O[V] {
	return MinBy(m, less)
}

// MinMaxBy returns the minimum and maximum value of the map, compared by
// function less.
func (m M[K, V]) MinMaxBy(less func(V, V) bool) goption.O[tuple.T2[V, V]] {
	return MinMaxBy(m, less)
}

// SumBy sums the values returned by function f.
func (m M[K, V]) SumBy[N constraints.Number](f func(V) N) N {
	return SumBy(m, f)
}

// AvgBy returns the average of the values returned by function f.
func (m M[K, V]) AvgBy[N constraints.Number](f func(V) N) float64 {
	return AvgBy(m, f)
}

// Chunk splits the map into chunks of the given size.
func (m M[K, V]) Chunk(size int) []M[K, V] {
	return Chunk(m, size)
}

// Divide splits the map into n chunks as evenly as possible.
func (m M[K, V]) Divide(n int) []M[K, V] {
	return Divide(m, n)
}

// CountBy returns the number of items that f returns true for.
func (m M[K, V]) CountBy(f func(K, V) bool) int {
	return CountBy(m, f)
}

// CountValueBy returns the number of items whose value satisfies f.
func (m M[K, V]) CountValueBy(f func(V) bool) int {
	return CountValueBy(m, f)
}

// Pop returns and deletes an arbitrary item of the map.
func (m M[K, V]) Pop() goption.O[V] {
	return Pop(m)
}

// PopItem returns and deletes an arbitrary item (key-value pair) of the map.
func (m M[K, V]) PopItem() goption.O[tuple.T2[K, V]] {
	return PopItem(m)
}

// Peek returns an arbitrary value of the map without deleting it.
func (m M[K, V]) Peek() goption.O[V] {
	return Peek(m)
}

// PeekItem returns an arbitrary item (key-value pair) of the map without
// deleting it.
func (m M[K, V]) PeekItem() goption.O[tuple.T2[K, V]] {
	return PeekItem(m)
}

// ToSlice converts the map to a slice, the element of the slice is returned by f.
//
// 🚀 EXAMPLE:
//
//	Wrap(map[int]int{1: 2}).ToSlice(func(k, v int) string { return strconv.Itoa(k+v) })
//	// ["3"]
func (m M[K, V]) ToSlice[R any](f func(K, V) R) []R {
	return ToSlice(m, f)
}

// TypeAssert converts the values from type V to type R by [type assertion].
//
// ⚠️ WARNING: It may ❌PANIC❌ when type assertion failed.
//
// [type assertion]: https://go.dev/tour/methods/15
func (m M[K, V]) TypeAssert[R any]() M[K, R] {
	return TypeAssert[R](m)
}

// MK is a wrapper of map[K]V whose key is ordered, it provides the operations
// that need ordered keys as methods.
//
// Use [WrapMK] to wrap an existing map into MK, [MK.M] to convert it to [M].
//
// 🚀 EXAMPLE:
//
//	WrapMK(map[int]int{2: 3, 1: 2}).OrderedKeys()  ⏩ [1, 2]
type MK[K constraints.Ordered, V any] map[K]V

// WrapMK wraps an existing map m into [MK], so that the operations that need
// ordered keys can be chained as methods.
//
// 🚀 EXAMPLE:
//
//	WrapMK(map[int]int{2: 3, 1: 2}).OrderedKeys()  ⏩ [1, 2]
func WrapMK[K constraints.Ordered, V any](m map[K]V) MK[K, V] {
	return MK[K, V](m)
}

// Unwrap returns the wrapped map of MK.
func (m MK[K, V]) Unwrap() map[K]V {
	return map[K]V(m)
}

// M converts MK to [M], so that the operations without constraint on the key
// are available.
func (m MK[K, V]) M() M[K, V] {
	return M[K, V](m)
}

// OrderedKeys returns the keys of the map in ascending order.
//
// 🚀 EXAMPLE:
//
//	WrapMK(map[int]int{2: 3, 1: 2}).OrderedKeys()  ⏩ [1, 2]
func (m MK[K, V]) OrderedKeys() []K {
	return OrderedKeys(m)
}

// OrderedValues returns the values of the map, ordered by the keys in ascending
// order.
//
// 🚀 EXAMPLE:
//
//	WrapMK(map[int]int{2: 3, 1: 2}).OrderedValues()  ⏩ [2, 3]
func (m MK[K, V]) OrderedValues() []V {
	return OrderedValues(m)
}

// OrderedItems returns the items (key-value pairs) of the map, ordered by the
// keys in ascending order.
func (m MK[K, V]) OrderedItems() tuple.S2[K, V] {
	return OrderedItems(m)
}

// ToOrderedSlice converts the map to a slice ordered by the keys in ascending
// order, the element of the slice is returned by f.
//
// 🚀 EXAMPLE:
//
//	WrapMK(map[int]int{2: 3, 1: 2}).ToOrderedSlice(func(k, v int) string { return strconv.Itoa(k) })
//	// ["1", "2"]
func (m MK[K, V]) ToOrderedSlice[R any](f func(K, V) R) []R {
	return ToOrderedSlice(m, f)
}
