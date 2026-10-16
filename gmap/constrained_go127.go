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

package gmap

import (
	"github.com/bytedance/gg/collection/tuple"
	"github.com/bytedance/gg/goption"
	"github.com/bytedance/gg/internal/constraints"
)

// MV is a wrapper of map[K]V whose value is comparable, it provides the
// operations that need comparable values as methods.
//
// Use [WrapMV] to wrap an existing map into MV, [MV.M] to convert it to [M].
//
// 🚀 EXAMPLE:
//
//	WrapMV(map[string]int{"a": 1, "b": 2}).Invert().Unwrap()  ⏩ {1: "a", 2: "b"}
type MV[K comparable, V comparable] map[K]V

// WrapMV wraps an existing map m into [MV], so that the operations that need
// comparable values can be chained as methods.
//
// 🚀 EXAMPLE:
//
//	WrapMV(map[string]int{"a": 1}).Equal(map[string]int{"a": 1})  ⏩ true
func WrapMV[K comparable, V comparable](m map[K]V) MV[K, V] {
	return MV[K, V](m)
}

// Unwrap returns the wrapped map of MV.
func (m MV[K, V]) Unwrap() map[K]V {
	return map[K]V(m)
}

// M converts MV to [M], so that the operations without constraint on the value
// are available.
func (m MV[K, V]) M() M[K, V] {
	return M[K, V](m)
}

// FilterByValues returns a new map with the items whose value is in the given
// values.
func (m MV[K, V]) FilterByValues(values ...V) MV[K, V] {
	return FilterByValues(m, values...)
}

// RejectByValues is negation of [MV.FilterByValues].
func (m MV[K, V]) RejectByValues(values ...V) MV[K, V] {
	return RejectByValues(m, values...)
}

// LoadKey returns the key of the first item whose value is v.
func (m MV[K, V]) LoadKey(v V) goption.O[K] {
	return LoadKey(m, v)
}

// Equal returns whether the two maps are equal.
func (m MV[K, V]) Equal(m2 map[K]V) bool {
	return Equal(m, m2)
}

// Compact removes the items whose value is zero.
func (m MV[K, V]) Compact() MV[K, V] {
	return Compact(m)
}

// Count returns the number of items whose value is v.
func (m MV[K, V]) Count(v V) int {
	return Count(m, v)
}

// Invert swaps the keys and the values of the map.
//
// 🚀 EXAMPLE:
//
//	WrapMV(map[string]int{"a": 1}).Invert().Unwrap()  ⏩ {1: "a"}
func (m MV[K, V]) Invert() MV[V, K] {
	return Invert(m)
}

// InvertBy is variant of [MV.Invert], the conflict is resolved by function
// onConflict.
func (m MV[K, V]) InvertBy(onConflict ConflictFunc[V, K]) MV[V, K] {
	return InvertBy(m, onConflict)
}

// InvertGroup swaps the keys and the values of the map, the keys with the same
// value are grouped into a slice.
//
// 🚀 EXAMPLE:
//
//	WrapMV(map[string]int{"a": 1, "b": 1}).InvertGroup()  ⏩ {1: ["a", "b"]}
func (m MV[K, V]) InvertGroup() map[V][]K {
	return InvertGroup(m)
}

// MO is a wrapper of map[K]V whose value is ordered, it provides the operations
// that need ordered values as methods.
//
// Use [WrapMO] to wrap an existing map into MO, [MO.M] to convert it to [M] and
// [MO.MV] to convert it to [MV] (ordered values are comparable).
//
// 🚀 EXAMPLE:
//
//	WrapMO(map[string]int{"a": 1, "b": 2}).Max().Value()  ⏩ 2
type MO[K comparable, V constraints.Ordered] map[K]V

// WrapMO wraps an existing map m into [MO], so that the operations that need
// ordered values can be chained as methods.
//
// 🚀 EXAMPLE:
//
//	WrapMO(map[string]int{"a": 1, "b": 2}).Max().Value()  ⏩ 2
func WrapMO[K comparable, V constraints.Ordered](m map[K]V) MO[K, V] {
	return MO[K, V](m)
}

// Unwrap returns the wrapped map of MO.
func (m MO[K, V]) Unwrap() map[K]V {
	return map[K]V(m)
}

// M converts MO to [M], so that the operations without constraint on the value
// are available.
func (m MO[K, V]) M() M[K, V] {
	return M[K, V](m)
}

// MV converts MO to [MV], since ordered values are comparable.
func (m MO[K, V]) MV() MV[K, V] {
	return MV[K, V](m)
}

// Max returns the maximum value of the map.
//
// 🚀 EXAMPLE:
//
//	WrapMO(map[string]int{"a": 1, "b": 2}).Max().Value()  ⏩ 2
func (m MO[K, V]) Max() goption.O[V] {
	return Max(m)
}

// Min returns the minimum value of the map.
//
// 🚀 EXAMPLE:
//
//	WrapMO(map[string]int{"a": 1, "b": 2}).Min().Value()  ⏩ 1
func (m MO[K, V]) Min() goption.O[V] {
	return Min(m)
}

// MinMax returns the minimum and maximum value of the map.
//
// 🚀 EXAMPLE:
//
//	WrapMO(map[string]int{"a": 1, "b": 2}).MinMax().Value().Values()  ⏩ 1 2
func (m MO[K, V]) MinMax() goption.O[tuple.T2[V, V]] {
	return MinMax(m)
}

// MN is a wrapper of map[K]V whose value is a number, it provides the
// arithmetic operations as methods.
//
// Use [WrapMN] to wrap an existing map into MN, [MN.M] to convert it to [M],
// [MN.MV] to convert it to [MV] and [MN.MO] to convert it to [MO].
//
// 🚀 EXAMPLE:
//
//	WrapMN(map[string]int{"a": 1, "b": 2}).Sum()  ⏩ 3
//	WrapMN(map[string]int{"a": 1, "b": 2}).Avg()  ⏩ 1.5
type MN[K comparable, V constraints.Number] map[K]V

// WrapMN wraps an existing map m into [MN], so that the arithmetic operations
// can be chained as methods.
//
// 🚀 EXAMPLE:
//
//	WrapMN(map[string]int{"a": 1, "b": 2}).Sum()  ⏩ 3
func WrapMN[K comparable, V constraints.Number](m map[K]V) MN[K, V] {
	return MN[K, V](m)
}

// Unwrap returns the wrapped map of MN.
func (m MN[K, V]) Unwrap() map[K]V {
	return map[K]V(m)
}

// M converts MN to [M], so that the operations without constraint on the value
// are available.
func (m MN[K, V]) M() M[K, V] {
	return M[K, V](m)
}

// MV converts MN to [MV], since numbers are comparable.
func (m MN[K, V]) MV() MV[K, V] {
	return MV[K, V](m)
}

// MO converts MN to [MO], since numbers are ordered.
func (m MN[K, V]) MO() MO[K, V] {
	return MO[K, V](m)
}

// Sum sums up the values of the map.
//
// 🚀 EXAMPLE:
//
//	WrapMN(map[string]int{"a": 1, "b": 2}).Sum()  ⏩ 3
func (m MN[K, V]) Sum() V {
	return Sum(m)
}

// Avg returns the average of the values of the map.
//
// 🚀 EXAMPLE:
//
//	WrapMN(map[string]int{"a": 1, "b": 2}).Avg()  ⏩ 1.5
func (m MN[K, V]) Avg() float64 {
	return Avg(m)
}
