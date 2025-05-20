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

// Package gmap provides generic operations for maps.
//
// 💡 HINT: We provide similar functionality for different types in different packages.
// For example, [github.com/bytedance/gg/gslice.Clone] for copying slice while
// [Clone] for copying map.
//
//   - Use [github.com/bytedance/gg/gslice] for slice operations.
//   - Use [github.com/bytedance/gg/gvalue] for value operations.
//   - Use [github.com/bytedance/gg/gptr] for pointer operations.
//   - …
//
// # Operations
//
// Keys / Values getter:
//
//   - [Keys], [Values]
//   - [OrderedKeys], [OrderedValues]
//   - [Items]
//   - [OrderedItems]
//   - [ToSlice]
//   - [ToOrderedSlice]
//
// CRUD operations:
//
//   - [Load], [LoadAll], [LoadSome], [LoadAny]
//   - [LoadOrStore], [LoadAndDelete]
//   - [Contains], [ContainsAny], [ContainsAll]
//
// Set operations:
//
//   - [Union], [Intersect], [Diff]
//
// Partition operations:
//
//   - [Chunk], [Divide]
//
// Math operations:
//
//   - [Max], [Min], [MinMax]
//   - [Sum], [Avg]
//
// Type casting/assertion/conversion:
//
//   - [TypeAssert]
//   - [PtrOf], [Indirect], [IndirectOr]
//
// Predicates:
//
//   - [Equal], [EqualStrict]
//
// High-order functions:
//
//   - [Map]
//   - [Filter], [Reject], [FilterMap]
//
// # Interface type satisfies comparable constraint after Go1.20 and later
//
// According [Go1.20 Language Change], Comparable types (such as ordinary interfaces)
// may now satisfy comparable constraints, even if the type arguments are not
// strictly comparable (comparison may PANIC at runtime).
//
// Which means operations of gmap can be used on more map types, for example:
//
//	var m map[io.Reader]string
//	readers = gmap.Keys (m)
//	// Go1.19 and earlier:	❌ io.Reader does not implement comparable
//	// Go1.20 and later:	✔️ running well
//
// It also means gmap operations may PANIC at runtime:
//
//	// Implement io.Reader for unhashable.
//	type unhashable []int
//	func (unhashable) Read([]byte) (_ int, _ error) { return }
//
//	m := make(map[io.Reader]string)
//	key := io.Reader(unhashable{})
//	_, _ = gmap.LoadOrStore(m, key, "")
//	// Go1.19 and earlier:	❌ io.Reader does not implement comparable
//	// Go1.20 and later:	❌ panic: runtime error: hash of unhashable type main.unhashable
//
// # Conflict resolution
//
// When operating on multiple maps, key conflicts often occur,
// the newer value replace the old by default ([DiscardOld]).
// These operations include but are not limited to [Invert], [Union],
// [Intersect] and so on…
//
// We provide [ConflictFunc] to help user to customize conflict resolution.
// All of above operations supports variant with ConflictFunc support:
//
//   - [Invert] ⏩ [InvertBy]
//   - [Union] ⏩ [UnionBy]
//   - [Intersect] ⏩ [IntersectBy]
//
// [Go1.20 Language Change]: https://tip.golang.org/doc/go1.20#language
package gmap

import (
	"sort"

	"github.com/bytedance/gg/collection/tuple"
	"github.com/bytedance/gg/goption"
	"github.com/bytedance/gg/gptr"
	"github.com/bytedance/gg/gresult"
	"github.com/bytedance/gg/gslice"
	"github.com/bytedance/gg/gvalue"
	"github.com/bytedance/gg/internal/constraints"
	"github.com/bytedance/gg/iter"
)

// Map applies function f to each key and value of map m.
// Results of f are returned as a new map.
//
// 🚀 EXAMPLE:
//
//	f := func(k, v int) (string, string) { return strconv.Itoa(k), strconv.Itoa(v) }
//	Map(map[int]int{1: 1}, f) ⏩ map[string]string{"1": "1"}
//	Map(map[int]int{}, f)     ⏩ map[string]string{}
//
// 💡 HINT:
//
//   - Use [MapKeys] if you only need to map the keys.
//   - Use [MapValues] if you only need to map the values.
//   - Use [FilterMap] if you also want to ignore keys/values during mapping.
//   - Use [ToSlice] if you want to "map" both key and value to single element
//   - Use [TryMap] if function f may fail (returns (K2, V2, error))
func Map[K1, K2 comparable, V1, V2 any](m map[K1]V1, f func(K1, V1) (K2, V2)) map[K2]V2 {
	r := make(map[K2]V2, len(m))
	for k, v := range m {
		k2, v2 := f(k, v)
		r[k2] = v2
	}
	return r
}

// TryMap is a variant of [Map] that allows function f to fail (return error).
//
// 🚀 EXAMPLE:
//
//	f := func(k, v int) (string, string, error) {
//		ki, kerr := strconv.Atoi(k)
//		vi, verr := strconv.Atoi(v)
//		return ki, vi, errors.Join(kerr, verr)
//	}
//	TryMap(map[string]string{"1": "1"}, f) ⏩ gresult.OK(map[int]int{1: 1})
//	TryMap(map[string]string{"1": "a"}, f) ⏩ gresult.Err("strconv.Atoi: parsing \"a\": invalid syntax")
//
// 💡 HINT:
//
//   - Use [TryFilterMap] if you want to ignore error during mapping.
//   - Use [TryMapKeys] if you only need to map the keys.
//   - Use [TryMapValues] if you only need to map the values.
func TryMap[K1, K2 comparable, V1, V2 any](m map[K1]V1, f func(K1, V1) (K2, V2, error)) gresult.R[map[K2]V2] {
	r := make(map[K2]V2, len(m))
	for k, v := range m {
		k2, v2, err := f(k, v)
		if err != nil {
			return gresult.Err[map[K2]V2](err)
		}
		r[k2] = v2
	}
	return gresult.OK(r)
}

// MapKeys is a variant of [Map], applies function f to each key of map m.
// Results of f and the corresponding values are returned as a new map.
//
// 🚀 EXAMPLE:
//
//	MapKeys(map[int]int{1: 1}, strconv.Itoa) ⏩ map[string]int{"1": 1}
//	MapKeys(map[int]int{}, strconv.Itoa)     ⏩ map[string]int{}
func MapKeys[K1, K2 comparable, V any](m map[K1]V, f func(K1) K2) map[K2]V {
	r := make(map[K2]V, len(m))
	for k, v := range m {
		r[f(k)] = v
	}
	return r
}

// TryMapKeys is a variant of [MapKeys] that allows function f to fail (return error).
//
// 🚀 EXAMPLE:
//
//	TryMapKeys(map[string]string{"1": "1"}, strconv.Atoi) ⏩ gresult.OK(map[int]string{1: "1"})
//	TryMapKeys(map[string]string{"a": "1"}, strconv.Atoi) ⏩ gresult.Err("strconv.Atoi: parsing \"a\": invalid syntax")
//	TryMapKeys(map[string]string{}, strconv.Itoa)         ⏩ gresult.OK(map[int]string{})
func TryMapKeys[K1, K2 comparable, V any](m map[K1]V, f func(K1) (K2, error)) gresult.R[map[K2]V] {
	r := make(map[K2]V, len(m))
	for k, v := range m {
		k2, err := f(k)
		if err != nil {
			return gresult.Err[map[K2]V](err)
		}
		r[k2] = v
	}
	return gresult.OK(r)
}

// MapValues is a variant of [Map], applies function f to each values of map m.
// Results of f and the corresponding keys are returned as a new map.
//
// 🚀 EXAMPLE:
//
//	MapValues(map[int]int{1: 1}, strconv.Itoa) ⏩ map[int]string{1: "1"}
//	MapValues(map[int]int{}, strconv.Itoa)     ⏩ map[int]string{}
func MapValues[K comparable, V1, V2 any](m map[K]V1, f func(V1) V2) map[K]V2 {
	r := make(map[K]V2, len(m))
	for k, v := range m {
		r[k] = f(v)
	}
	return r
}

// TryMapValues is a variant of [MapValues] that allows function f to fail (return error).
//
// 🚀 EXAMPLE:
//
//	TryMapValues(map[string]string{"1": "1"}, strconv.Atoi) ⏩ gresult.OK(map[string]int{"1": 1})
//	TryMapValues(map[string]string{"1": "a"}, strconv.Atoi) ⏩ gresult.Err("strconv.Atoi: parsing \"a\": invalid syntax")
//	TryMapValues(map[string]string{}, strconv.Itoa)         ⏩ gresult.OK(map[string]int{})
func TryMapValues[K comparable, V1, V2 any](m map[K]V1, f func(V1) (V2, error)) gresult.R[map[K]V2] {
	r := make(map[K]V2, len(m))
	for k, v := range m {
		v2, err := f(v)
		if err != nil {
			return gresult.Err[map[K]V2](err)
		}
		r[k] = v2
	}
	return gresult.OK(r)
}

// Filter applies predicate f to each key and value of map m,
// returns those keys and values that satisfy the predicate f as a new map.
//
// 🚀 EXAMPLE:
//
//	m := map[int]int{1: 1, 2: 2, 3: 2, 4: 3}
//	pred := func(k, v int) bool { return (k+v)%2 == 0 }
//	Filter(m, pred) ⏩ map[int]int{1: 1, 2: 2}
//
// 💡 HINT:
//
//   - Use [FilterKeys] if you only need to filter the keys.
//   - Use [FilterValues] if you only need to filter the values.
//   - Use [FilterMap] if you also want to modify the keys/values during filtering.
func Filter[K comparable, V any](m map[K]V, f func(K, V) bool) map[K]V {
	r := make(map[K]V, len(m)/2)
	for k, v := range m {
		if f(k, v) {
			r[k] = v
		}
	}
	return r
}

// FilterKeys is a variant of [Filter], applies predicate f to each key of map m,
// returns keys that satisfy the predicate f and the corresponding values as a new map.
//
// 🚀 EXAMPLE:
//
//	m := map[int]int{1: 1, 2: 2, 3: 2, 4: 3}
//	pred := func(k int) bool { return k%2 == 0 }
//	FilterKeys(m, pred) ⏩ map[int]int{2: 2, 4: 3}
func FilterKeys[K comparable, V any](m map[K]V, f func(K) bool) map[K]V {
	r := make(map[K]V, len(m)/2)
	for k, v := range m {
		if f(k) {
			r[k] = v
		}
	}
	return r
}

// FilterByKeys is a variant of [Filter], filters map m by given keys slice,
// returns a new map containing only the key-value pairs where the key exists in the keys slice.
//
// 🚀 EXAMPLE:
//
//	m := map[int]int{1: 1, 2: 2, 3: 3, 4: 4}
//	keys := []int{1, 3, 5}
//	FilterByKeys(m, keys) ⏩ map[int]int{1: 1, 3: 3}
func FilterByKeys[K comparable, V any](m map[K]V, keys ...K) map[K]V {
	r := make(map[K]V, gvalue.Min(len(keys), len(m)))
	for _, key := range keys {
		if v, ok := m[key]; ok {
			r[key] = v
		}
	}
	return r
}

// FilterValues is a variant of [Filter], applies predicate f to each value of map m,
// returns values that satisfy the predicate f and the corresponding keys as a new map.
//
// 🚀 EXAMPLE:
//
//	m := map[int]int{1: 1, 2: 2, 3: 2, 4: 3}
//	pred := func(v int) bool { return v%2 == 0 }
//	FilterValues(m, pred) ⏩ map[int]int{2: 2, 3: 2}
func FilterValues[K comparable, V any](m map[K]V, f func(V) bool) map[K]V {
	r := make(map[K]V, len(m)/2)
	for k, v := range m {
		if f(v) {
			r[k] = v
		}
	}
	return r
}

// FilterByValues is a variant of [Filter], filters map m by given values slice,
// returns a new map containing only the key-value pairs where the value exists in the values slice.
//
// 🚀 EXAMPLE:
//
//	m := map[int]int{1: 10, 2: 20, 3: 10, 4: 30}
//	values := []int{10, 30}
//	FilterByValues(m, values) ⏩ map[int]int{1: 10, 3: 10, 4: 30}
func FilterByValues[K, V comparable](m map[K]V, values ...V) map[K]V {
	r := make(map[K]V, gvalue.Min(len(values), len(m)))
	for k, v := range m {
		if gslice.Contains(values, v) {
			r[k] = v
		}
	}
	return r
}

// Reject applies predicate f to each key and value of map m,
// returns those keys and values that do not satisfy the predicate f as a new map.
//
// 🚀 EXAMPLE:
//
//	m := map[int]int{1: 1, 2: 2, 3: 2, 4: 3}
//	pred := func(k, v int) bool { return (k+v)%2 != 0 }
//	Reject(m, pred) ⏩ map[int]int{1: 1, 2: 2}
//
// 💡 HINT:
//
//   - Use [RejectKeys] if you only need to reject the keys.
//   - Use [RejectValues] if you only need to reject the values.
func Reject[K comparable, V any](m map[K]V, f func(K, V) bool) map[K]V {
	r := make(map[K]V, len(m)/2)
	for k, v := range m {
		if !f(k, v) {
			r[k] = v
		}
	}
	return r
}

// RejectKeys applies predicate f to each key of map m,
// returns keys that do not satisfy the predicate f and the corresponding values as a new map.
//
// 🚀 EXAMPLE:
//
//	m := map[int]int{1: 1, 2: 2, 3: 2, 4: 3}
//	pred := func(k int) bool { return k%2 != 0 }
//	RejectKeys(m, pred) ⏩ map[int]int{2: 2, 4: 3}
func RejectKeys[K comparable, V any](m map[K]V, f func(K) bool) map[K]V {
	r := make(map[K]V, len(m)/2)
	for k, v := range m {
		if !f(k) {
			r[k] = v
		}
	}
	return r
}

// RejectByKeys is the opposite of [FilterByKeys], removes entries from map m where the key exists in the keys slice,
// returns a new map containing only the key-value pairs where the key does not exist in the keys slice.
//
// 🚀 EXAMPLE:
//
//	m := map[int]int{1: 1, 2: 2, 3: 3, 4: 4}
//	keys := []int{1, 3}
//	RejectByKeys(m, keys) ⏩ map[int]int{2: 2, 4: 4}
func RejectByKeys[K comparable, V any](m map[K]V, keys ...K) map[K]V {
	r := Clone(m)
	for _, key := range keys {
		delete(r, key)
	}
	return r
}

// RejectValues applies predicate f to each value of map m,
// returns values that do not satisfy the predicate f and the corresponding keys as a new map.
//
// 🚀 EXAMPLE:
//
//	 m := map[int]int{1: 1, 2: 2, 3: 2, 4: 3}
//	 pred := func(v int) bool { return v%2 != 0 }
//		RejectValues(m, pred) ⏩ map[int]int{2: 2, 3: 2}
func RejectValues[K comparable, V any](m map[K]V, f func(V) bool) map[K]V {
	r := make(map[K]V, len(m)/2)
	for k, v := range m {
		if !f(v) {
			r[k] = v
		}
	}
	return r
}

// RejectByValues is the opposite of [FilterByValues], removes entries from map m where the value exists in the values slice,
// returns a new map containing only the key-value pairs where the value does not exist in the values slice.
//
// 🚀 EXAMPLE:
//
//	m := map[int]int{1: 10, 2: 20, 3: 10, 4: 30}
//	values := []int{10, 30}
//	RejectByValues(m, values) ⏩ map[int]int{2: 20}
func RejectByValues[K, V comparable](m map[K]V, values ...V) map[K]V {
	r := make(map[K]V, len(m)/2)
	for k, v := range m {
		if !gslice.Contains(values, v) {
			r[k] = v
		}
	}
	return r
}

// TODO: Unhidden Fold/Reduce funcs
//
// fold applies function f cumulatively to each key and value of map m,
// so as to fold the map to a single value.
//
//	fold(map[int]int{1: 1, 2: 2}, func(acc, k, v int) int { return acc + k + v }, 0) ⏩ 6
func fold[K comparable, V, T any](m map[K]V, f func(T, K, V) T, init T) T {
	acc := init
	for k, v := range m {
		acc = f(acc, k, v)
	}
	return acc
}

// foldKeys applies function f cumulatively to each key of map m,
// so as to fold the keys of map to a single value.
func foldKeys[K comparable, V, T any](m map[K]V, f func(T, K) T, init T) T {
	acc := init
	for k := range m {
		acc = f(acc, k)
	}
	return acc
}

// foldValues applies function f cumulatively to each value of map m,
// so as to fold the values of map to a single value.
func foldValues[K comparable, V, T any](m map[K]V, f func(T, V) T, init T) T {
	acc := init
	for _, v := range m {
		acc = f(acc, v)
	}
	return acc
}

// reduce is a variant of Fold, use possible first key value tuple of map as
// the initial value of accumulation.
func reduce[K comparable, V any, KV tuple.T2[K, V]](m map[K]V, f func(KV, K, V) KV) goption.O[KV] {
	var acc KV
	var inited bool
	for k, v := range m {
		if inited {
			acc = f(acc, k, v)
		} else {
			acc = KV(tuple.Make2(k, v))
			inited = true
		}
	}
	return goption.Of(acc, inited)
}

// reduceKeys is a variant of FoldKeys, use possible first key of map as the
// initial value of accumulation.
func reduceKeys[K comparable, V any](m map[K]V, f func(K, K) K) goption.O[K] {
	var acc K
	var inited bool
	for k := range m {
		if inited {
			acc = f(acc, k)
		} else {
			acc = k
			inited = true
		}
	}
	return goption.Of(acc, inited)
}

// reduceValues is a variant of FoldValues, use possible first value of map as
// the initial value of accumulation.
func reduceValues[K comparable, V any](m map[K]V, f func(V, V) V) goption.O[V] {
	var acc V
	var inited bool
	for _, v := range m {
		if inited {
			acc = f(acc, v)
		} else {
			acc = v
			inited = true
		}
	}
	return goption.Of(acc, inited)
}

// Keys returns the keys of the map m.
//
// 🚀 EXAMPLE:
//
//	m := map[int]string{1: "1", 2: "2", 3: "3", 4: "4"}
//	Keys(m) ⏩ []int{1, 3, 2, 4} //⚠️INDETERMINATE ORDER⚠️
//
// ⚠️ WARNING: The keys will be in an indeterminate order,
// use [OrderedKeys] to get them in fixed order.
//
// 💡 HINT: If you want to merge key and value to single element, use [ToSlice].
func Keys[K comparable, V any](m map[K]V) []K {
	return iter.ToSlice(iter.FromMapKeys(m))
}

// Values returns the values of the map m.
//
// 🚀 EXAMPLE:
//
//	m := map[int]string{1: "1", 2: "2", 3: "3", 4: "4"}
//	Values(m) ⏩ []string{"1", "4", "2", "3"} //⚠️INDETERMINATE ORDER⚠️
//
// ⚠️ WARNING: The keys values be in an indeterminate order,
// use [OrderedValues] to get them in fixed order.
//
// 💡 HINT: If you want to merge key and value to single element, use [ToSlice].
func Values[K comparable, V any](m map[K]V) []V {
	return iter.ToSlice(iter.FromMapValues(m))
}

// OrderedKeys is the variant of [Keys],
// returns the keys of the map m in fixed order.
//
// 🚀 EXAMPLE:
//
//	m := map[int]string{1: "1", 2: "2", 3: "3", 4: "4"}
//	OrderedKeys(m) ⏩ []int{1, 2, 3, 4}
//
// 💡 HINT: If you want to merge key and value to single element, use [ToOrderedSlice].
//
// 💡 AKA: SortedKey
func OrderedKeys[K constraints.Ordered, V any](m map[K]V) []K {
	return iter.ToSlice(iter.Sort(iter.FromMapKeys(m)))
}

// OrderedValues is variant of [Values],
// returns the values of the map m in fixed order.
//
// 🚀 EXAMPLE:
//
//	m := map[int]string{1: "1", 2: "2", 3: "3", 4: "4"}
//	OrderedValues(m) ⏩ []string{"1", "2", "3", "4"}
//
// 💡 HINT: If you want to merge key and value to single element, use [ToOrderedSlice].
//
// 💡 AKA: SortedValues
func OrderedValues[K constraints.Ordered, V any](m map[K]V) []V {
	return iter.ToSlice(
		iter.Map(func(kv tuple.T2[K, V]) V { return kv.Second },
			iter.SortBy(func(kv1, kv2 tuple.T2[K, V]) bool { return kv1.First < kv2.First },
				iter.FromMap(m))))
}

// Items returns both the keys and values of the map m.
//
// 🚀 EXAMPLE:
//
//	m := map[int]string{1: "1", 2: "2", 3: "3"}
//	Items(m)          ⏩ []tuple.S2{{2, "2"}, {1, "1"}, {3, "3"}} // ⚠️INDETERMINATE ORDER⚠️
//	Items(m).Values() ⏩ []int{2, 1, 3}, []string{"2", "1", "3"}  // ⚠️INDETERMINATE ORDER⚠️
//
// ⚠️ WARNING: The returned items will be in an indeterminate order,
// use [OrderedItems] to get them in fixed order.
//
// 💡 HINT: The keys and values are returned in the form of a slice of tuples,
// and the keys slice values slice can be obtained separately through the
// [github.com/bytedance/gg/collection/tuple.S2.Values] method.
//
// 💡 AKA: KeyValues, KeyAndValues
func Items[K comparable, V any](m map[K]V) tuple.S2[K, V] {
	items := make(tuple.S2[K, V], 0, len(m))
	for k, v := range m {
		items = append(items, tuple.Make2(k, v))
	}
	return items
}

// OrderedItems is variant of [Items],
// returns the keys and values of the map m in fixed order.
//
// 🚀 EXAMPLE:
//
//	m := map[int]string{1: "1", 2: "2", 3: "3"}
//	OrderedItems(m)          ⏩ []tuple.S2{{1, "1"}, {2, "2"}, {3, "3"}}
//	OrderedItems(m).Values() ⏩ []int{1, 2, 3}, []string{"1", "2", "3"}
//
// 💡 HINT: The keys and values are returned in the form of a slice of tuples,
// and the keys slice values slice can be obtained separately through the
// [github.com/bytedance/gg/collection/tuple.S2.Values] method.
//
// 💡 AKA: SortedItems, SortedKeyValues, SortedKeyAndValues
func OrderedItems[K constraints.Ordered, V any](m map[K]V) tuple.S2[K, V] {
	items := Items(m)
	sort.Slice(items, func(i, j int) bool {
		return items[i].First < items[j].First
	})
	return items
}

// Merge is alias of [Union].
func Merge[K comparable, V any](ms ...map[K]V) map[K]V {
	return Union(ms...)
}

// Union returns the unions of maps as a new map.
//
// 💡 NOTE:
//
//   - Once the key conflicts, the newer value always replace the older one ([DiscardOld]),
//     use [UnionBy] and [ConflictFunc] to customize conflict resolution.
//   - If the result is an empty set, always return an empty map instead of nil
//
// 🚀 EXAMPLE:
//
//	m := map[int]int{1: 1, 2: 2}
//	Union(m, nil)             ⏩ map[int]int{1: 1, 2: 2}
//	Union(m, map[int]{3: 3})  ⏩ map[int]int{1: 1, 2: 2, 3: 3}
//	Union(m, map[int]{2: -1}) ⏩ map[int]int{1: 1, 2: -1} // "2:2" is replaced by the newer "2:-1"
//
// 💡 HINT: Use [github.com/bytedance/gg/collection/set.Set] if you need a
// set data structure
//
// 💡 AKA: Merge, Concat, Combine
func Union[K comparable, V any](ms ...map[K]V) map[K]V {
	// Fastpath: no map or only one map given.
	if len(ms) == 0 {
		return make(map[K]V)
	}
	if len(ms) == 1 {
		return cloneWithoutNilCheck(ms[0])
	}

	var maxLen int
	for _, m := range ms {
		maxLen = gvalue.Max(maxLen, len(m))
	}
	ret := make(map[K]V, maxLen)
	// Fastpath: all maps are empty.
	if maxLen == 0 {
		return ret
	}

	// Union all maps.
	for _, m := range ms {
		for k, v := range m {
			ret[k] = v
		}
	}
	return ret
}

// UnionBy returns the unions of maps as a new map, conflicts are resolved by a
// custom [ConflictFunc].
//
// 🚀 EXAMPLE:
//
//	m := map[int]int{1: 1, 2: 2}
//	Union(m, map[int]{2: 0})                               ⏩ map[int]int{1: 1, 2: 0} // "2:2" is replaced by the newer "2:0"
//	UnionBy(gslice.Of(m, map[int]int{2: 0}), DiscardOld()) ⏩ map[int]int{1: 1, 2: 0} // same as above
//	UnionBy(gslice.Of(m, map[int]int{2: 0}), DiscardNew()) ⏩ map[int]int{1: 1, 2: 2} // "2:2" is kept because it is older
//
// For more examples, see [ConflictFunc].
func UnionBy[K comparable, V any, M ~map[K]V](ms []M, onConflict ConflictFunc[K, V]) M {
	// Fastpath: no map or only one map given.
	if len(ms) == 0 {
		return make(M)
	}
	if len(ms) == 1 {
		return cloneWithoutNilCheck(ms[0])
	}

	var maxLen int
	for _, m := range ms {
		maxLen = gvalue.Max(maxLen, len(m))
	}
	ret := make(map[K]V, maxLen)
	// Fastpath: all maps are empty.
	if maxLen == 0 {
		return ret
	}

	// Union all maps with ConflictFunc.
	for _, m := range ms {
		for k, newV := range m {
			if oldV, ok := ret[k]; ok {
				ret[k] = onConflict(k, oldV, newV)
			} else {
				ret[k] = newV
			}
		}
	}
	return ret
}

// Diff returns the difference of map m against other maps as a new map.
//
// 💡 NOTE: If the result is an empty set, always return an empty map instead of nil
//
// 🚀 EXAMPLE:
//
//	m := map[int]int{1: 1, 2: 2}
//	Diff(m, nil)             ⏩ map[int]int{1: 1, 2: 2}
//	Diff(m, map[int]{1: 1})  ⏩ map[int]int{2: 2}
//	Diff(m, map[int]{3: 3})  ⏩ map[int]int{1: 1, 2: 2}
//
// 💡 HINT: Use [github.com/bytedance/gg/collection/set.Set] if you need a
// set data structure
//
// TODO: Value type of againsts can be diff from m.
func Diff[K comparable, V any](m map[K]V, againsts ...map[K]V) map[K]V {
	if len(m) == 0 {
		return make(map[K]V)
	}
	if len(againsts) == 0 {
		return cloneWithoutNilCheck(m)
	}
	ret := make(map[K]V, len(m)/2)
	for k, v := range m {
		var found bool
		for _, a := range againsts {
			if _, found = a[k]; found {
				break
			}
		}
		if !found {
			ret[k] = v
		}
	}
	return ret
}

// Intersect returns the intersection of maps as a new map.
//
// 💡 NOTE:
//
//   - Once the key conflicts, the newer one will replace the older one ([DiscardOld]),
//     use [IntersectBy] and [ConflictFunc] to customize conflict resolution.
//   - If the result is an empty set, always return an empty map instead of nil
//
// 🚀 EXAMPLE:
//
//	m := map[int]int{1: 1, 2: 2}
//	Intersect(m, nil)             ⏩ map[int]int{}
//	Intersect(m, map[int]{3: 3})  ⏩ map[int]int{}
//	Intersect(m, map[int]{1: 1})  ⏩ map[int]int{1: 1}
//	Intersect(m, map[int]{1: -1}) ⏩ map[int]int{1: -1} // "1:1" is replaced by the newer "1:-1"
//
// 💡 HINT: Use [github.com/bytedance/gg/collection/set.Set] if you need a
// set data structure
func Intersect[K comparable, V any](ms ...map[K]V) map[K]V {
	// Fastpath: no map or only one map given.
	if len(ms) == 0 {
		return make(map[K]V)
	}
	if len(ms) == 1 {
		return cloneWithoutNilCheck(ms[0])
	}

	minLen := len(ms[0])
	for _, m := range ms[1:] {
		minLen = gvalue.Min(minLen, len(m))
	}
	ret := make(map[K]V, minLen)
	// Fastpath: all maps are empty.
	if minLen == 0 {
		return ret
	}

	// Intersect all maps.
	for k, v := range ms[0] {
		found := true // at least we found it in ms[0]
		for _, m := range ms[1:] {
			if v, found = m[k]; !found {
				break
			}
		}
		if found {
			ret[k] = v
		}
	}
	return ret
}

// IntersectBy returns the intersection of maps as a new map, conflicts are resolved by a
// custom [ConflictFunc].
//
// 🚀 EXAMPLE:
//
//	m := map[int]int{1: 1, 2: 2}
//	Intersect(m, map[int]{1: -1})                                     ⏩ map[int]int{1: -1} // "1:1" is replaced by the newer "1:-1"
//	IntersectBy(gslice.Of(m, map[int]{1: -1}), DiscardOld[int,int]()) ⏩ map[int]int{1: -1} // same as above
//	IntersectBy(gslice.Of(m, map[int]{1: -1}), DiscardNew[int,int]()) ⏩ map[int]int{1: 1} // "1:1" is kept because it is older
//
// For more examples, see [ConflictFunc].
func IntersectBy[K comparable, V any, M ~map[K]V](ms []M, onConflict ConflictFunc[K, V]) M {
	if len(ms) == 0 {
		return make(M)
	}
	if len(ms) == 1 {
		return cloneWithoutNilCheck(ms[0])
	}
	minLen := len(ms[0])
	for _, m := range ms[1:] {
		minLen = gvalue.Min(minLen, len(m))
	}
	ret := make(map[K]V, minLen)
	// Fastpath: all maps are empty.
	if minLen == 0 {
		return ret
	}
	for k, v := range ms[0] {
		found := true // at least we found it in ms[0]
		for _, m := range ms[1:] {
			var tmp V
			if tmp, found = m[k]; !found {
				break
			} else {
				v = onConflict(k, v, tmp)
			}
		}
		if found {
			ret[k] = v
		}
	}
	return ret
}

// Load returns the value stored in the map for a key.
//
// If the value was not found in the map. goption.Nil[V]() is returned.
//
// If the given map is nil, goption.Nil[V]() is returned.
//
// 💡 HINT: See also [LoadAny], [LoadAll], [LoadSome] if you have multiple values
// to load.
//
// 💡 AKA: Get
func Load[K comparable, V any](m map[K]V, k K) goption.O[V] {
	if m == nil || len(m) == 0 {
		return goption.Nil[V]()
	}
	v, ok := m[k]
	if !ok {
		return goption.Nil[V]()
	}
	return goption.OK(v)
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
//
// The loaded result is true if the value was loaded, false if stored.
//
// ⚠️ WARNING: LoadOrStore panics when a nil map is given.
//
// 💡 AKA: setdefault
func LoadOrStore[K comparable, V any](m map[K]V, k K, defaultV V) (v V, loaded bool) {
	assertNonNilMap(m)
	v, loaded = m[k]
	if !loaded {
		v = defaultV
		m[k] = v
	}
	return
}

// LoadOrStoreLazy returns the existing value for the key if present.
// Otherwise, it stores and returns the value that lazy computed by function f.
//
// The loaded result is true if the value was loaded, false if stored.
//
// ⚠️ WARNING: LoadOrStoreLazy panics when a nil map is given.
func LoadOrStoreLazy[K comparable, V any](m map[K]V, k K, f func() V) (v V, loaded bool) {
	assertNonNilMap(m)
	v, loaded = m[k]
	if !loaded {
		v = f()
		m[k] = v
	}
	return
}

// LoadAndDelete deletes the value for a key, returning the previous value if any.
//
// 🚀 EXAMPLE:
//
//	var m = map[string]int { "foo": 1 }
//	LoadAndDelete(m, "bar") ⏩ goption.Nil()
//	LoadAndDelete(m, "foo") ⏩ goption.OK(1)
//	LoadAndDelete(m, "foo") ⏩ goption.Nil()
//
// 💡 HINT: If you want to delete an element "randomly", use [Pop].
func LoadAndDelete[K comparable, V any](m map[K]V, k K) goption.O[V] {
	if m == nil || len(m) == 0 {
		return goption.Nil[V]()
	}
	v, ok := m[k]
	if !ok {
		return goption.Nil[V]()
	}
	delete(m, k)
	return goption.OK(v)
}

// LoadKey find the first key that mapped to the specified value.
//
// 💡 NOTE: LoadKey has O(N) time complexity.
//
// 🚀 EXAMPLE:
//
//	m := map[int]string{1: "1", 2: "2", 3: "3", 4: "4"}
//	LoadKey(m, "1") ⏩ goption.OK(1)
//	LoadKey(m, "a") ⏩ goption.Nil[int]()
//
// 💡 AKA: FindKey, FindByKey, GetKey, GetByKey
func LoadKey[K, V comparable](m map[K]V, v V) goption.O[K] {
	for k, vv := range m {
		if vv == v {
			return goption.OK(k)
		}
	}
	return goption.Nil[K]()
}

// LoadBy find the first value that satisfy the predicate f.
//
// 💡 NOTE: LoadBy has O(N) time complexity.
//
// 💡 AKA: FindBy, FindValueBy, GetBy, GetValueBy
func LoadBy[K comparable, V any](m map[K]V, f func(K, V) bool) goption.O[V] {
	if len(m) == 0 {
		return goption.Nil[V]()
	}
	for k, v := range m {
		if f(k, v) {
			return goption.OK(v)
		}
	}
	return goption.Nil[V]()
}

// LoadKeyBy find the first key that satisfy the predicate f.
//
// 💡 NOTE: LoadKeyBy has O(N) time complexity.
//
// 💡 AKA: FindKeyBy, GetKeyBy
func LoadKeyBy[K comparable, V any](m map[K]V, f func(K, V) bool) goption.O[K] {
	if len(m) == 0 {
		return goption.Nil[K]()
	}
	for k, v := range m {
		if f(k, v) {
			return goption.OK(k)
		}
	}
	return goption.Nil[K]()
}

// LoadItemBy find the first key-value pair that satisfy the predicate f.
//
// 💡 NOTE: LoadItemBy has O(N) time complexity.
//
// 💡 AKA: FindItemBy, GetItemBy
func LoadItemBy[K comparable, V any](m map[K]V, f func(K, V) bool) goption.O[tuple.T2[K, V]] {
	if len(m) == 0 {
		return goption.Nil[tuple.T2[K, V]]()
	}
	for k, v := range m {
		if f(k, v) {
			return goption.OK(tuple.Make2(k, v))
		}
	}
	return goption.Nil[tuple.T2[K, V]]()
}

// LoadAll returns the all value stored in the map for given keys.
//
// If not all keys are not found in the map, nil is returned.
// Otherwise, the length of returned values should equal the length of given keys.
//
// 🚀 EXAMPLE:
//
//	m := map[int]string{1: "1", 2: "2", 3: "3"}
//	LoadAll(m, 1, 2) ⏩ []string{"1", "2"}
//	LoadAll(m, 1, 4) ⏩ nil
func LoadAll[K comparable, V any](m map[K]V, ks ...K) []V {
	if m == nil || len(m) == 0 || len(ks) == 0 {
		return nil
	}
	vs := make([]V, 0, len(ks))
	for _, k := range ks {
		v, ok := m[k]
		if !ok {
			return nil
		}
		vs = append(vs, v)
	}
	return vs
}

// LoadAny returns the all value stored in the map for given keys.
//
// If no value is found in the map, goption.Nil[V]() is returned.
// Otherwise, the first found value is returned.
//
// 🚀 EXAMPLE:
//
//	m := map[int]string{1: "1", 2: "2", 3: "3"}
//	LoadAny(m, 1, 2) ⏩ goption.OK("1")
//	LoadAny(m, 5, 1) ⏩ goption.OK("1")
//	LoadAny(m, 5, 6) ⏩ goption.Nil[string]()
func LoadAny[K comparable, V any](m map[K]V, ks ...K) (r goption.O[V]) {
	if m == nil || len(m) == 0 || len(ks) == 0 {
		return
	}
	for _, k := range ks {
		if v, ok := m[k]; ok {
			return goption.OK(v)
		}
	}
	return
}

// LoadSome returns the some values stored in the map for given keys.
//
// 🚀 EXAMPLE:
//
//	m := map[int]string{1: "1", 2: "2", 3: "3"}
//	LoadSome(m, 1, 2) ⏩ []string{"1", "2"}
//	LoadSome(m, 1, 4) ⏩ []string{"1"}
func LoadSome[K comparable, V any](m map[K]V, ks ...K) []V {
	if m == nil || len(m) == 0 || len(ks) == 0 {
		return nil
	}
	var vs []V
	for _, k := range ks {
		if v, ok := m[k]; ok {
			vs = append(vs, v)
		}
	}
	return vs
}

// Invert inverts the keys and values of map, and returns a new map.
// (map[K]V] → map[V]K).
//
// ⚠️ WARNING: The iteration of the map is in an indeterminate order,
// for multiple KV-pairs with the same V, the K retained after inversion is UNSTABLE.
// If the length of the returned map is equal to the length of the given map,
// there are no key conflicts.
// Use [InvertBy] and [ConflictFunc] to customize conflict resolution.
// Use [InvertGroup] to avoid key loss when multiple keys mapped to the same value.
//
// 🚀 EXAMPLE:
//
//	Invert(map[string]int{"1": 1, "2": 2}) ⏩ map[int]string{1: "1", 2: "2"},
//	Invert(map[string]int{"1": 1, "2": 1}) ⏩ ⚠️ UNSTABLE: map[int]string{1: "1"} or map[int]string{1: "2"}
//
// 💡 AKA: Reverse
func Invert[K, V comparable](m map[K]V) map[V]K {
	r := make(map[V]K)
	for k, v := range m {
		r[v] = k
	}
	return r
}

// InvertBy inverts the keys and values of map, and returns a new map.
// (map[K]V] → map[V]K), conflicts are resolved by a custom [ConflictFunc].
//
// 💡 NOTE: the "oldVal", and "newVal" naming of [ConflictFunc] are meaningless
// because of the map's indeterminate iteration order. Further,
// [DiscardOld] and [DiscardNew] are also meaningless.
//
// 🚀 EXAMPLE:
//
//	InvertBy(map[string]int{"1": 1, "": 1}, DiscardZero(nil) ⏩ map[int]string{1: "1"},
func InvertBy[K, V comparable](m map[K]V, onConflict ConflictFunc[V, K]) map[V]K {
	r := make(map[V]K)
	for k, v := range m {
		if oldKey, ok := r[v]; ok {
			r[v] = onConflict(v, oldKey, k)
		} else {
			r[v] = k
		}
	}
	return r
}

// InvertGroup inverts the map by grouping keys that mapped to the same value into a slice.
// (map[K]V] → map[V][]K).
//
// ⚠️ WARNING: The iteration of the map is in an indeterminate order,
// for multiple KV-pairs with the same V, the order of K in the slice is UNSTABLE.
//
// 🚀 EXAMPLE:
//
//	InvertGroup(map[string]int{"1": 1, "2": 2}) ⏩ map[int][]string{1: {"1"}, 2: {"2"}},
//	InvertGroup(map[string]int{"1": 1, "2": 1}) ⏩ ⚠️ UNSTABLE: map[int][]string{1: {"1", "2"}} or map[int]string{1: {"2", "1"}}
func InvertGroup[K, V comparable](m map[K]V) map[V][]K {
	r := make(map[V][]K)
	for k, v := range m {
		r[v] = append(r[v], k)
	}
	return r
}

// Equal reports whether two maps contain the same key/value pairs.
// Values are compared using ==.
//
// 💡 NOTE: Equal does NOT distinguish between nil and empty maps
// (which means Equal(map[int]int{}, nil) returns true), use [EqualStrict] if necessary.
//
// 🚀 EXAMPLE:
//
//	Equal(map[int]int{1: 1, 2: 2}, map[int]int{1: 1, 2: 2}) ⏩ true
//	Equal(map[int]int{1: 1}, map[int]int{1: 1, 2: 2})       ⏩ false
//	Equal(map[int]int{}, map[int]int{})                     ⏩ true
//	Equal(map[int]int{}, nil)                               ⏩ true
func Equal[K, V comparable](m1, m2 map[K]V) bool {
	if len(m1) != len(m2) {
		return false
	}

	for k, v1 := range m1 {
		if v2, ok := m2[k]; !ok || v1 != v2 {
			return false
		}
	}
	return true
}

// EqualBy reports whether two maps contain the same key/value pairs.
// Values are compared using function eq.
//
// 💡 NOTE: EqualBy does NOT distinguish between nil and empty maps
// (which means Equal(map[int]int{}, nil, gvalue.Equal[int]) returns true),
// use [EqualStrictBy] if necessary.
//
// 🚀 EXAMPLE:
//
//	eq := gvalue.Equal[int]
//	EqualBy(map[int]int{1: 1, 2: 2}, map[int]int{1: 1, 2: 2}, eq) ⏩ true
//	EqualBy(map[int]int{1: 1}, map[int]int{1: 1, 2: 2}, eq)       ⏩ false
//	EqualBy(map[int]int{}, map[int]int{}, eq)                     ⏩ true
//	EqualBy(map[int]int{}, nil, eq)                               ⏩ true
func EqualBy[K comparable, V any](m1, m2 map[K]V, eq func(v1, v2 V) bool) bool {
	if len(m1) != len(m2) {
		return false
	}

	for k, v1 := range m1 {
		if v2, ok := m2[k]; !ok || !eq(v1, v2) {
			return false
		}
	}
	return true
}

// EqualStrict is a variant of [Equal], which can distinguish between nil and empty maps.
//
// 🚀 EXAMPLE:
//
//	EqualStrict(map[int]int{1: 1, 2: 2}, map[int]int{1: 1, 2: 2}) ⏩ true
//	EqualStrict(map[int]int{1: 1}, map[int]int{1: 1, 2: 2})       ⏩ false
//	EqualStrict(map[int]int{}, map[int]int{})                     ⏩ true
//	EqualStrict(map[int]int{}, nil)                               ⏩ false
func EqualStrict[K, V comparable](m1, m2 map[K]V) bool {
	if (m1 == nil && m2 != nil) || (m1 != nil && m2 == nil) {
		return false
	}
	return Equal(m1, m2)
}

// EqualStrictBy is a variant of [EqualBy], which can distinguish between nil and empty maps.
//
// 🚀 EXAMPLE:
//
//	eq := gvalue.Equal[int]
//	EqualStrictBy(map[int]int{1: 1, 2: 2}, map[int]int{1: 1, 2: 2}, eq) ⏩ true
//	EqualStrictBy(map[int]int{1: 1}, map[int]int{1: 1, 2: 2}, eq)       ⏩ false
//	EqualStrictBy(map[int]int{}, map[int]int{}, eq)                     ⏩ true
//	EqualStrictBy(map[int]int{}, nil, eq)                               ⏩ false
func EqualStrictBy[K comparable, V any](m1, m2 map[K]V, eq func(v1, v2 V) bool) bool {
	if m1 == nil && m2 != nil {
		return false
	} else if m1 != nil && m2 == nil {
		return false
	}
	return EqualBy(m1, m2, eq)
}

// Clone returns a shallow copy of map.
// If the given map is nil, nil is returned.
//
// 🚀 EXAMPLE:
//
//	Clone(map[int]int{1: 1, 2: 2}) ⏩ map[int]int{1: 1, 2: 2}
//	Clone(map[int]int{})           ⏩ map[int]int{}
//	Clone[int, int](nil)           ⏩ nil
//
// 💡 HINT: Both keys and values are copied using assignment (=), so this is a shallow clone.
// If you want to do a deep clone, use [CloneBy] with an appropriate value
// clone function.
//
// 💡 AKA: Copy
func Clone[K comparable, V any, M ~map[K]V](m M) M {
	if m == nil {
		return nil
	}
	return cloneWithoutNilCheck(m)
}

// CloneBy is variant of [Clone], it returns a copy of the map.
// Elements are copied using function f.
// If the given map is nil, nil is returned.
//
// TODO: Example
//
// 💡 AKA: CopyBy
func CloneBy[K comparable, V any, M ~map[K]V](m M, f func(V) V) M {
	if m == nil {
		return nil
	}
	return MapValues(m, f)
}

func cloneWithoutNilCheck[K comparable, V any, M ~map[K]V](m M) M {
	r := make(M, len(m))
	for k, v := range m {
		r[k] = v
	}
	return r
}

func assertNonNilMap[K comparable, V any](m map[K]V) {
	if m == nil {
		panic("nil map is not allowed")
	}
}

// Contains returns whether the key occur in map.
//
// 🚀 EXAMPLE:
//
//	m := map[int]string{1: ""}
//	Contains(m, 1)             ⏩ true
//	Contains(m, 0)             ⏩ false
//	var nilMap map[int]string
//	Contains(nilMap, 0)        ⏩ false
//
// 💡 HINT: See also [ContainsAll], [ContainsAny] if you have multiple values to
// query.
func Contains[K comparable, V any](m map[K]V, k K) bool {
	if m == nil || len(m) == 0 {
		return false
	}
	_, ok := m[k]
	return ok
}

// ContainsAny returns whether any of given keys occur in map.
//
// 🚀 EXAMPLE:
//
//	m := map[int]string{1: "", 2: ""}
//	ContainsAny(m, 1, 2) ⏩ true
//	ContainsAny(m, 1, 3) ⏩ true
//	ContainsAny(m, 3)    ⏩ false
func ContainsAny[K comparable, V any](m map[K]V, ks ...K) bool {
	if m == nil || len(m) == 0 {
		return false
	}
	for _, k := range ks {
		if _, ok := m[k]; ok {
			return true
		}
	}
	return false
}

// ContainsAll returns whether all of given keys occur in map.
//
// 🚀 EXAMPLE:
//
//	m := map[int]string{1: "", 2: "",}
//	ContainsAll(m, 1, 2) ⏩ true
//	ContainsAll(m, 1, 3) ⏩ false
//	ContainsAll(m, 3)    ⏩ false
func ContainsAll[K comparable, V any](m map[K]V, ks ...K) bool {
	if (m == nil || len(m) == 0) && len(ks) != 0 {
		return false
	}
	for _, k := range ks {
		if _, ok := m[k]; !ok {
			return false
		}
	}
	return true
}

// Sum returns the arithmetic sum of the values of map m.
//
// 🚀 EXAMPLE:
//
//	Sum(map[string]int{"1": 1, "2": 2, "3": 3}) ⏩ 6
//
// 💡 NOTE: The returned type is still T, it may overflow for smaller types
// (such as int8, uint8).
func Sum[K comparable, V constraints.Number](m map[K]V) V {
	return iter.Sum(iter.FromMapValues(m))
}

// SumBy applies function f to each value of map m,
// returns the arithmetic sum of function result.
func SumBy[K comparable, V any, N constraints.Number](m map[K]V, f func(V) N) N {
	return iter.SumBy(f, iter.FromMapValues(m))
}

// Avg returns the arithmetic mean of the values of map s.
//
// 🚀 EXAMPLE:
//
//	Avg(map[string]int{"1": 1, "2": 2, "3": 3}) ⏩ 2.0
//
// 💡 AKA: Mean, Average
func Avg[K comparable, V constraints.Number](m map[K]V) float64 {
	return iter.Avg(iter.FromMapValues(m))
}

// AvgBy applies function f to each values of map m,
// returns the arithmetic mean of function result.
//
// 💡 AKA: MeanBy, AverageBy
func AvgBy[K comparable, V any, N constraints.Number](m map[K]V, f func(V) N) float64 {
	return iter.AvgBy(f, iter.FromMapValues(m))
}

// Max returns the maximum value of map m.
//
// 🚀 EXAMPLE:
//
//	Max(map[string]int{"1": 1, "2": 2, "3": 3}) ⏩ goption.OK(3)
//
// 💡 NOTE: If the given map is empty, goption.Nil[T]() is returned.
func Max[K comparable, V constraints.Ordered](m map[K]V) goption.O[V] {
	return iter.Max(iter.FromMapValues(m))
}

// MaxBy returns the maximum value of map m determined by function less.
//
// 🚀 EXAMPLE:
//
//	type Foo struct { Value int }
//	less := func(x, y Foo) bool { return x.Value < y.Value }
//	MaxBy(map[string]Foo{"1": {1}, "2": {2}, "3": {3}}, less) ⏩ goption.OK(Foo{3})
//
// 💡 NOTE: If the given map is empty, goption.Nil[V]() is returned.
func MaxBy[K comparable, V any](m map[K]V, less func(V, V) bool) goption.O[V] {
	return iter.MaxBy(less, iter.FromMapValues(m))
}

// Min returns the minimum element of map m.
//
// 🚀 EXAMPLE:
//
//	Min(map[string]int{"1": 1, "2": 2, "3": 3}) ⏩ goption.OK(1)
//
// 💡 NOTE: If the given map is empty, goption.Nil[V]() is returned.
func Min[K comparable, V constraints.Ordered](m map[K]V) goption.O[V] {
	return iter.Min(iter.FromMapValues(m))
}

// MinBy returns the minimum value of map m determined by function less.
//
// 🚀 EXAMPLE:
//
//	type Foo struct { Value int }
//	less := func(x, y Foo) bool { return x.Value < y.Value }
//	MinBy(map[string]Foo{"1": {1}, "2": {2}, "3": {3}}, less) ⏩ goption.OK(Foo{1})
//
// 💡 NOTE: If the given map is empty, goption.Nil[V]() is returned.
func MinBy[K comparable, V any](m map[K]V, less func(V, V) bool) goption.O[V] {
	return iter.MinBy(less, iter.FromMapValues(m))
}

// MinMax returns both minimum and maximum elements of map m.
// If the given map is empty, goption.Nil[tuple.T2[V,V]]() is returned.
//
//	MinMax(map[string]int{"1": 1, "2": 2, "3": 3}) ⏩ goption.OK(tuple.T2{1, 3})
//
// 💡 AKA: Bound
func MinMax[K comparable, V constraints.Ordered](m map[K]V) goption.O[tuple.T2[V, V]] {
	return iter.MinMax(iter.FromMapValues(m))
}

// MinMaxBy returns both minimum and maximum elements of map m determined
// by function less.
// If the given map is empty, goption.Nil[tuple.T2[V,V]]() is returned.
//
// 🚀 EXAMPLE:
//
//	type Foo struct { Value int }
//	less := func(x, y Foo) bool { return x.Value < y.Value }
//	m := map[string]Foo{"1": {1}, "2": {2}, "3": {3}}
//	MinMaxBy(m, less) ⏩ goption.OK(tuple.T2{Foo{1}, Foo{3}})
//
// 💡 AKA: BoundBy
func MinMaxBy[K comparable, V any](m map[K]V, less func(V, V) bool) goption.O[tuple.T2[V, V]] {
	return iter.MinMaxBy(less, iter.FromMapValues(m))
}

// Chunk splits map into length-n chunks and returns chunks by a new slice.
//
// The last chunk will be shorter if n does not evenly divide the length of the map.
//
// ⚠️ WARNING: The values in chunks will be in an indeterminate order.
func Chunk[K comparable, V any](m map[K]V, size int) []map[K]V {
	return iter.ToSlice(
		iter.Map(func(s []tuple.T2[K, V]) map[K]V { return iter.KVToMap(iter.StealSlice(s)) },
			iter.Chunk(size,
				iter.FromMap(m))))
}

// Divide splits map into exactly n slices and returns chunks by a new slice.
//
// The length of chunks will be different if n does not evenly divide the length
// of the slice.
//
// ⚠️ WARNING: The values in chunks will be in an indeterminate order.
func Divide[K comparable, V any](m map[K]V, n int) []map[K]V {
	return iter.ToSlice(
		iter.Map(func(s []tuple.T2[K, V]) map[K]V { return iter.KVToMap(iter.StealSlice(s)) },
			iter.Divide(n,
				iter.FromMap(m))))
}

// PtrOf returns pointers that point to equivalent values of map m.
// (map[K]V → map[K]*V).
//
// 🚀 EXAMPLE:
//
//	PtrOf(map[int]string{1: "1", 2: "2"}) ⏩ map[int]*string{1: (*string)("1"), 2: (*string)("2")}
//
// ⚠️ WARNING: The returned pointers do not point to values of the original
// map, user CAN NOT modify the value by modifying the pointer.
func PtrOf[K comparable, V any](m map[K]V) map[K]*V {
	return MapValues(m, gptr.Of[V])
}

// Indirect returns the values pointed to by the pointers.
// If the pointer is nil, filter it out of the returned map.
//
// 🚀 EXAMPLE:
//
//		v1, v2 := "1", "2"
//	 m := map[int]*string{ 1: &v1, 2: &v2, 3: nil}
//	 Indirect(m) ⏩ map[int]string{1: "1", 2: "2"}
//
// 💡 HINT: If you want to replace nil pointer with default value,
// use [IndirectOr].
func Indirect[K comparable, V any](m map[K]*V) map[K]V {
	ret := make(map[K]V, len(m)/2)
	for k, v := range m {
		if v == nil {
			continue
		}
		ret[k] = *v
	}
	return ret
}

// IndirectOr is variant of [Indirect].
// If the pointer is nil, returns the fallback value instead.
//
// 🚀 EXAMPLE:
//
//		v1, v2 := "1", "2"
//	 m := map[int]*string{ 1: &v1, 2: &v2, 3: nil}
//	 IndirectOr(m, "nil") ⏩ map[int]string{1: "1", 2: "2", 3: "nil"}
func IndirectOr[K comparable, V any](m map[K]*V, fallback V) map[K]V {
	ret := make(map[K]V, len(m))
	for k, v := range m {
		if v == nil {
			ret[k] = fallback
		} else {
			ret[k] = *v
		}
	}
	return ret
}

// TypeAssert converts values of map from type From to type To by type assertion.
//
// 🚀 EXAMPLE:
//
//	TypeAssert[int](map[int]any{1: 1, 2: 2})   ⏩ map[int]int{1: 1, 2: 2}
//	TypeAssert[any](map[int]int{1: 1, 2: 2})   ⏩ map[int]any{1: 1, 2: 2}
//	TypeAssert[int64](map[int]int{1: 1, 2: 2}) ⏩ ❌PANIC❌
//
// ⚠️ WARNING:
//
//   - This function may ❌PANIC❌.
//     See [github.com/bytedance/gg/gvalue.TypeAssert] for more details
func TypeAssert[To any, K comparable, From any](m map[K]From) map[K]To {
	return MapValues(m, gvalue.TypeAssert[To, From])
}

// Len returns the length of map m.
//
// 💡 HINT: This function is designed for high-order function, because the builtin
// function can not be used as function pointer.
// For example, if you want to get the total length of a 2D slice:
//
//	var s []map[int]int
//	total1 := SumBy(s, len)          // ❌ERROR❌ len (built-in) must be called
//	total2 := SumBy(s, Len[int,int]) // OK
func Len[K comparable, V any](m map[K]V) int {
	return len(m)
}

// Compact removes all zero values from given map m, returns a new map.
//
// 🚀 EXAMPLE:
//
//	m := map[int]string{0: "", 1: "foo", 2: "", 3: "bar"}
//	Compact(m) ⏩ map[int]string{1: "foo", 3: "bar"}
//
// 💡 HINT: See [github.com/bytedance/gg/gvalue.Zero] for details of zero value.
func Compact[K, V comparable](m map[K]V) map[K]V {
	return FilterValues(m, gvalue.IsNotZero[V])
}

// ToSlice converts the map m to a slice by function f.
//
// ⚠️ WARNING: The returned slice will be in an indeterminate order,
// use [ToOrderedSlice] to get them in fixed order.
//
// 🚀 EXAMPLE:
//
//	f := func (k, v int) string {
//		return fmt.Sprintf("%d: %d", k, v)
//	}
//	m := map[int]int{1: 1, 2: 2, 3: 3}
//	ToSlice(m, f) ⏩ []string{"1: 1", "3: 3", "2: 2"} //⚠️INDETERMINATE ORDER⚠️
//
// 💡 HINT:
//
//   - If you only need the key slice or value slice, use [Keys] or [Values].
//   - If you need the key-value pair slice, use [Items].
//   - See also [github.com/bytedance/gg/gslice.ToMap].
func ToSlice[K comparable, V, T any](m map[K]V, f func(K, V) T) []T {
	return gslice.Map(Items(m), func(kv tuple.T2[K, V]) T { return f(kv.Values()) })
}

// ToOrderedSlice is variant of [ToSlice], the returned slice is in fixed order.
//
// 🚀 EXAMPLE:
//
//	f := func (k, v int) string {
//		return fmt.Sprintf("%d: %d", k, v)
//	}
//	m := map[int]int{1: 1, 2: 2, 3: 3}
//	ToOrderedSlice(m, f) ⏩ []string{"1: 1", "2: 2", "3: 3"}
func ToOrderedSlice[K constraints.Ordered, V, T any](m map[K]V, f func(K, V) T) []T {
	return gslice.Map(OrderedItems(m), func(kv tuple.T2[K, V]) T { return f(kv.Values()) })
}

// FilterMap does [Filter] and [Map] at the same time, applies function f to
// each key and value of map m. f returns (K2, V2, bool):
//
//   - If true ,the returned key and value will added to the result map[K2]V2.
//   - If false, the returned key and value will be dropped.
//
// 🚀 EXAMPLE:
//
//	f := func(k, v int) (string, string, bool) { return strconv.Itoa(k), strconv.Itoa(v), k != 0 && v != 0 }
//	FilterMap(map[int]int{1: 1, 2: 0, 0: 3}, f) ⏩ map[string]string{"1": "1"}
func FilterMap[K1, K2 comparable, V1, V2 any](m map[K1]V1, f func(K1, V1) (K2, V2, bool)) map[K2]V2 {
	r := make(map[K2]V2, len(m)/2)
	for k, v := range m {
		if kk, vv, ok := f(k, v); ok {
			r[kk] = vv
		}
	}
	return r
}

// TryFilterMap is a variant of [FilterMap] that allows function f to fail (return error).
//
// 🚀 EXAMPLE:
//
//	f := func(k, v int) (string, string, error) {
//		ki, kerr := strconv.Atoi(k)
//		vi, verr := strconv.Atoi(v)
//		return ki, vi, errors.Join(kerr, verr)
//	}
//	TryFilterMap(map[string]string{"1": "1", "2": "2"}, f) ⏩ map[int]int{1: 1, 2: 2}
//	TryFilterMap(map[string]string{"1": "a", "2": "2"}, f) ⏩ map[int]int{2: 2})
func TryFilterMap[K1, K2 comparable, V1, V2 any](m map[K1]V1, f func(K1, V1) (K2, V2, error)) map[K2]V2 {
	r := make(map[K2]V2, len(m)/2)
	for k, v := range m {
		if kk, vv, err := f(k, v); err == nil {
			r[kk] = vv
		}
	}
	return r
}

// FilterMapKeys is a variant of [FilterMap].
//
// 🚀 EXAMPLE:
//
//	f := func(v int) (string, bool) { return strconv.Itoa(v), v != 0 }
//	FilterMapKeys(map[int]int{1: 1, 2: 0, 0: 3}, f) ⏩ map[string]int{"1": 1, "2": 0}
func FilterMapKeys[K1, K2 comparable, V any](m map[K1]V, f func(K1) (K2, bool)) map[K2]V {
	r := make(map[K2]V, len(m)/2)
	for k, v := range m {
		if kk, ok := f(k); ok {
			r[kk] = v
		}
	}
	return r
}

// TryFilterMapKeys is a variant of [FilterMapKeys] that allows function f to fail (return error).
//
// 🚀 EXAMPLE:
//
//	FilterMapKeys(map[string]string{"1": "1", "2": "2"}, strconv.Atoi) ⏩ map[int]string{1: "1", 2: "2"}
//	FilterMapKeys(map[string]string{"1": "1", "a": "2"}, strconv.Atoi) ⏩ map[int]string{1: "1"}
func TryFilterMapKeys[K1, K2 comparable, V any](m map[K1]V, f func(K1) (K2, error)) map[K2]V {
	r := make(map[K2]V, len(m)/2)
	for k, v := range m {
		if kk, err := f(k); err == nil {
			r[kk] = v
		}
	}
	return r
}

// FilterMapValues is a variant of [FilterMap].
//
// 🚀 EXAMPLE:
//
//	f := func(v int) (string, bool) { return strconv.Itoa(v), v != 0 }
//	FilterMapValues(map[int]int{1: 1, 2: 0, 0: 3}, f) ⏩ map[int]string{1: "1", 0: "3"}
func FilterMapValues[K comparable, V1, V2 any](m map[K]V1, f func(V1) (V2, bool)) map[K]V2 {
	r := make(map[K]V2, len(m)/2)
	for k, v := range m {
		if vv, ok := f(v); ok {
			r[k] = vv
		}
	}
	return r
}

// TryFilterMapValues is a variant of [FilterMapValues] that allows function f to fail (return error).
//
// 🚀 EXAMPLE:
//
//	FilterMapValues(map[string]string{"1": "1", "2": "2"}, strconv.Atoi) ⏩ map[string]int{"1": 1, "2": 2}
//	FilterMapValues(map[string]string{"1": "1", "2": "a"}, strconv.Atoi) ⏩ map[string]int{"1": 1}
func TryFilterMapValues[K comparable, V1, V2 any](m map[K]V1, f func(V1) (V2, error)) map[K]V2 {
	r := make(map[K]V2, len(m)/2)
	for k, v := range m {
		if vv, err := f(v); err == nil {
			r[k] = vv
		}
	}
	return r
}

// ConflictFunc is used to merge the conflicting key-value of map operations.
//
// Once the key conflicts, the conflicting key and the corresponding values will
// be passed to this function, and the user can resolve the conflict by
// returning a new value. Here are some pre-defined ConflictFuncs:
//
//   - [DiscardOld]: this is the default behavior of most map operations
//   - [DiscardNew]
//   - [DiscardZero]
//   - [DiscardNil]
//
// 🚀 EXAMPLE:
//
//	Union(
//		map[int]int{1: 1, 2: 2},
//		map[int]int{      2: 0})        ⏩ map[int]int{1: 1, 2: 0} // "2:2" is replaced by the newer "2:0"
//
//	UnionBy(
//		gslice.Of(,
//			map[int]int{1: 1, 2: 2},
//			map[int]int{      2: 0}),
//		DiscardOld())                   ⏩ map[int]int{1: 1, 2: 0} // same as above, DiscardOld is the default behavior
//
//	UnionBy(
//		gslice.Of(,
//			map[int]int{1: 1, 2: 2},
//			map[int]int{      2: 0}),
//		DiscardNew())                   ⏩ map[int]int{1: 1, 2: 2} // "2:2" is kept because the newer is always discarded
//
//	UnionBy(
//		gslice.Of(,
//			map[int]int{      2: 0},
//			map[int]int{1: 1, 2: 1},
//			map[int]int{1: 1, 2: 2},
//			map[int]int{      2: 0}),
//		DiscardZero(nil))               ⏩ map[int]int{1: 1, 2: 2} // "2:2" is kept because 2 is the newest non-zero value
//
//	UnionBy(
//		gslice.Of(,
//			map[int]int{      2: 0},
//			map[int]int{1: 1, 2: 1},
//			map[int]int{1: 1, 2: 2},
//			map[int]int{      2: 0}),
//		DiscardZero(DiscardNew()))      ⏩ map[int]int{1: 1, 2: 1} // "2:1" is kept because 1 is the oldest non-zero value
type ConflictFunc[K comparable, V any] func(key K, oldVal, newVal V) V

// DiscardOld returns a [ConflictFunc] that always return the newer value.
func DiscardOld[K comparable, V any]() ConflictFunc[K, V] {
	return discardOld[K, V]
}

// discardOld is a internal implementation of [DiscardOld].
func discardOld[K comparable, V any](_ K, _, newVal V) V {
	return newVal
}

// DiscardNew returns a [ConflictFunc] that always return the older value.
func DiscardNew[K comparable, V any]() ConflictFunc[K, V] {
	return func(_ K, oldVal, _ V) V { return oldVal }
}

// DiscardZero returns a [ConflictFunc] that always return the non-zero value.
//
// 💡 NOTE: If both values are non-zero, the fallback function will be called.
// If the fallback function is nil, [DiscardOld] will be called.
//
// 💡 HINT: See [github.com/bytedance/gg/gvalue.Zero] for details of zero value.
func DiscardZero[K comparable, V comparable](fallback ConflictFunc[K, V]) ConflictFunc[K, V] {
	zeroVal := gvalue.Zero[V]()
	return func(key K, oldVal, newVal V) V {
		if oldVal == zeroVal {
			if newVal != zeroVal {
				return newVal
			} else if fallback != nil {
				return fallback(key, oldVal, newVal)
			} else {
				return discardOld(key, oldVal, newVal)
			}
		} else {
			if newVal == zeroVal {
				return oldVal
			} else if fallback != nil {
				return fallback(key, oldVal, newVal)
			} else {
				return discardOld(key, oldVal, newVal)
			}
		}
	}
}

// DiscardNil returns a [ConflictFunc] that always return the non-zero value.
//
// 💡 NOTE: If both values are non-nil, the fallback function will be called.
// If the fallback function is nil, [DiscardOld] will be called.
func DiscardNil[K comparable, V comparable](fallback ConflictFunc[K, *V]) ConflictFunc[K, *V] {
	return func(key K, oldVal, newVal *V) *V {
		if oldVal == nil {
			if newVal != nil {
				return newVal
			} else if fallback != nil {
				return fallback(key, oldVal, newVal)
			} else {
				return discardOld(key, oldVal, newVal)
			}
		} else {
			if newVal == nil {
				return oldVal
			} else if fallback != nil {
				return fallback(key, oldVal, newVal)
			} else {
				return discardOld(key, oldVal, newVal)
			}
		}
	}
}

// Count returns the times of value v that occur in map m.
//
// 🚀 EXAMPLE:
//
//	Count(map[int]string{1: "a", 2: "a", 3: "b"}, "a") ⏩ 2
//
// 💡 HINT:
//
//   - Use [CountValueBy] if type of v is non-comparable
//   - Use [CountBy] if you need to consider key when counting
func Count[K, V comparable](m map[K]V, v V) int {
	var count int
	for _, vv := range m {
		if vv == v {
			count++
		}
	}
	return count
}

// CountBy returns the times of pair (k, v) in map m that satisfy the predicate f.
//
// 🚀 EXAMPLE:
//
//	f := func (k int, v string) bool {
//		i, _ := strconv.Atoi(v)
//		return k%2 == 1 && i%2 == 1
//	}
//	CountBy(map[int]string{1: "1", 2: "2", 3: "3"}, f) ⏩ 0
//	CountBy(map[int]string{1: "1", 2: "2", 3: "4"}, f) ⏩ 1
func CountBy[K comparable, V any](m map[K]V, f func(K, V) bool) int {
	var count int
	for k, v := range m {
		if f(k, v) {
			count++
		}
	}
	return count
}

// CountValueBy returns the times of value v in map m that satisfy the predicate f.
//
// 🚀 EXAMPLE:
//
//	f := func (v string) bool {
//		i, _ := strconv.Atoi(v)
//		return i%2 == 1
//	}
//	CountValueBy(map[int]string{1: "1", 2: "2", 3: "3"}, f) ⏩ 2
//	CountValueBy(map[int]string{1: "1", 2: "2", 3: "4"}, f) ⏩ 1
func CountValueBy[K comparable, V any](m map[K]V, f func(V) bool) int {
	var count int
	for _, v := range m {
		if f(v) {
			count++
		}
	}
	return count
}

// Pop tries to load and DELETE a "random" element from map m. If m is empty,
// goption.Nil[V]() is returned.
//
// 🚀 EXAMPLE:
//
//	var m = map[string]int { "foo": 1 }
//	Pop(m) ⏩ goption.OK(1)
//	Pop(m) ⏩ goption.Nil()
//
// ⚠️ WARNING: As map iteration is indeterminate ordered, we said it is "random".
//
// 💡 HINT:
//
//   - If you don't want to delete the element, use [Peek]
//   - If you want to delete element by key, use [LoadAndDelete]
//   - If you want to know the key of poped value, use [PopItem]
func Pop[K comparable, V any](m map[K]V) goption.O[V] {
	for k, v := range m {
		delete(m, k)
		return goption.OK(v)
	}
	return goption.Nil[V]()
}

// PopItem is variant of [Pop], return key-value pair instead of a single value.
func PopItem[K comparable, V any](m map[K]V) goption.O[tuple.T2[K, V]] {
	for k, v := range m {
		delete(m, k)
		return goption.OK(tuple.Make2(k, v))
	}
	return goption.Nil[tuple.T2[K, V]]()
}

// Peek tries to load a "random" element from map m. If m is empty,
// goption.Nil[V]() is returned.
//
// 🚀 EXAMPLE:
//
//	var m = map[string]int { "foo": 1 }
//	Peek(m) ⏩ goption.OK(1)
//	var m2 = map[string]int {}
//	Peek(m2) ⏩ goption.Nil()
//
// ⚠️ WARNING: As map iteration is indeterminate ordered, we said it is "random".
//
// 💡 HINT:
//
//   - If you want to delete the returned value, use [Pop]
//   - If you also want to know the key of returned value, use [PeekItem]
func Peek[K comparable, V any](m map[K]V) goption.O[V] {
	for _, v := range m {
		return goption.OK(v)
	}
	return goption.Nil[V]()
}

// PeekItem is variant of [Peek], return key-value pair instead of a single value.
func PeekItem[K comparable, V any](m map[K]V) goption.O[tuple.T2[K, V]] {
	for k, v := range m {
		return goption.OK(tuple.Make2(k, v))
	}
	return goption.Nil[tuple.T2[K, V]]()
}
