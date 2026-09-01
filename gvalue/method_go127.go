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

// This file contains the method form of the value operations.
//
// Operations of a value T can not be methods of T itself (T may be any type,
// including the ones from other packages), so they used to be package-level
// functions only. Go 1.27 generic methods make it possible to wrap a value into
// [V] (or its constrained variants [Ord], [Num]) and chain the operations,
// where the methods that change the type (such as [V.Cast] and [V.Map]) declare
// their own type parameters:
//
//	Of(any(1)).Cast[int]()        ⏩ 1
//	Of(1).Map(strconv.Itoa).Get() ⏩ "1"
//	OrderedOf(5).Clamp(1, 3)      ⏩ 3
package gvalue

import (
	"fmt"
	"reflect"

	"github.com/bytedance/gg/internal/constraints"
)

// V is a wrapper of value T, it provides the value operations as methods.
//
// Use [Of] to wrap a value into V, and [V.Get] to unwrap it.
//
// 🚀 EXAMPLE:
//
//	Of(any(1)).Cast[int]()      ⏩ 1
//	Of(1).Map(strconv.Itoa).Get()  ⏩ "1"
//	Of([]int(nil)).IsZero()     ⏩ true
type V[T any] struct {
	v T
}

// Of wraps value v into [V], so that the value operations can be chained as methods.
//
// 🚀 EXAMPLE:
//
//	Of(1).Map(strconv.Itoa).Get()  ⏩ "1"
func Of[T any](v T) V[T] {
	return V[T]{v: v}
}

// Get returns the wrapped value of V.
//
// 💡 AKA: Unwrap, Value
func (v V[T]) Get() T {
	return v.v
}

// IsNil returns whether the wrapped value is nil.
//
// 💡 HINT: Refer to function [IsNil] for the definition of nil.
func (v V[T]) IsNil() bool {
	return IsNil(v.v)
}

// IsNotNil is negation of [V.IsNil].
func (v V[T]) IsNotNil() bool {
	return !v.IsNil()
}

// IsZero returns whether the wrapped value is zero value.
//
// It is a variant of function [IsZero] without the comparable constraint:
// non-comparable values (such as slice) are checked by reflection.
//
// 🚀 EXAMPLE:
//
//	Of(0).IsZero()          ⏩ true
//	Of([]int(nil)).IsZero() ⏩ true
//	Of([]int{}).IsZero()    ⏩ false
func (v V[T]) IsZero() bool {
	av := any(v.v)
	if av == nil {
		return true
	}
	return reflect.ValueOf(av).IsZero()
}

// IsNotZero is negation of [V.IsZero].
func (v V[T]) IsNotZero() bool {
	return !v.IsZero()
}

// Equal returns whether the wrapped value is equal to o.
//
// Comparable values are compared by the == operator (so the result is the same
// as function [Equal]), non-comparable values (such as slice) fall back to
// [reflect.DeepEqual].
func (v V[T]) Equal(o T) bool {
	av, ao := any(v.v), any(o)
	if av == nil || ao == nil {
		return av == ao
	}
	if rv := reflect.TypeOf(av); rv != nil && !rv.Comparable() {
		return reflect.DeepEqual(av, ao)
	}
	return av == ao
}

// Cast converts the wrapped value from type T to type R by [type assertion].
//
// It is the method form of function [TypeAssert].
//
// ⚠️ WARNING: It may ❌PANIC❌ when type assertion failed, use [V.TryCast] for
// a panic-free variant.
//
// 🚀 EXAMPLE:
//
//	Of(any(1)).Cast[int]()  ⏩ 1
//
// [type assertion]: https://go.dev/tour/methods/15
func (v V[T]) Cast[R any]() R {
	return any(v.v).(R)
}

// TryCast is a variant of [V.Cast], it returns false instead of panicking when
// the type assertion failed.
//
// It is the method form of function [TryAssert].
//
// 🚀 EXAMPLE:
//
//	Of(any(1)).TryCast[int]()     ⏩ 1, true
//	Of(any(1)).TryCast[string]()  ⏩ "", false
func (v V[T]) TryCast[R any]() (R, bool) {
	r, ok := any(v.v).(R)
	return r, ok
}

// Map applies function f to the wrapped value and wraps the result into V[R].
//
// 🚀 EXAMPLE:
//
//	Of(1).Map(strconv.Itoa).Get()  ⏩ "1"
func (v V[T]) Map[R any](f func(T) R) V[R] {
	return Of(f(v.v))
}

// Then is a variant of [V.Map], the function f returns a V[R] instead of R.
//
// 🚀 EXAMPLE:
//
//	Of(1).Then(func(i int) V[string] { return Of(strconv.Itoa(i)) }).Get()  ⏩ "1"
func (v V[T]) Then[R any](f func(T) V[R]) V[R] {
	return f(v.v)
}

// String implements [fmt.Stringer].
func (v V[T]) String() string {
	av := any(v.v)
	if av == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%v", av)
}

// Ord is a wrapper of an ordered value T, it provides the ordering operations
// (such as [Ord.Max], [Ord.Clamp]) as methods.
//
// Use [OrderedOf] to wrap an ordered value into Ord.
//
// 🚀 EXAMPLE:
//
//	OrderedOf(5).Clamp(1, 3)    ⏩ 3
//	OrderedOf(2).Between(1, 3)  ⏩ true
type Ord[T constraints.Ordered] struct {
	v T
}

// OrderedOf wraps an ordered value v into [Ord].
func OrderedOf[T constraints.Ordered](v T) Ord[T] {
	return Ord[T]{v: v}
}

// Get returns the wrapped value of Ord.
func (o Ord[T]) Get() T {
	return o.v
}

// V casts Ord to [V].
func (o Ord[T]) V() V[T] {
	return Of(o.v)
}

// Max returns the maximum value of the wrapped value and inputs.
//
// 🚀 EXAMPLE:
//
//	OrderedOf(1).Max(2, 3)  ⏩ 3
func (o Ord[T]) Max(vals ...T) T {
	return Max(o.v, vals...)
}

// Min returns the minimum value of the wrapped value and inputs.
//
// 🚀 EXAMPLE:
//
//	OrderedOf(2).Min(1, 3)  ⏩ 1
func (o Ord[T]) Min(vals ...T) T {
	return Min(o.v, vals...)
}

// MinMax returns the minimum and maximum value of the wrapped value and inputs.
//
// 🚀 EXAMPLE:
//
//	OrderedOf(2).MinMax(1, 3)  ⏩ 1, 3
func (o Ord[T]) MinMax(vals ...T) (T, T) {
	return MinMax(o.v, vals...)
}

// Clamp returns the value if the wrapped value is within [min, max];
// otherwise returns the nearest boundary.
//
// 🚀 EXAMPLE:
//
//	OrderedOf(5).Clamp(1, 3)  ⏩ 3
func (o Ord[T]) Clamp(min, max T) T {
	return Clamp(o.v, min, max)
}

// Between returns true when the wrapped value is within [min, max].
//
// 🚀 EXAMPLE:
//
//	OrderedOf(2).Between(1, 3)  ⏩ true
func (o Ord[T]) Between(min, max T) bool {
	return Between(o.v, min, max)
}

// Less returns true when the wrapped value is less than v.
func (o Ord[T]) Less(v T) bool {
	return Less(o.v, v)
}

// LessEqual returns true when the wrapped value is less than or equal to v.
func (o Ord[T]) LessEqual(v T) bool {
	return LessEqual(o.v, v)
}

// Greater returns true when the wrapped value is greater than v.
func (o Ord[T]) Greater(v T) bool {
	return Greater(o.v, v)
}

// GreaterEqual returns true when the wrapped value is greater than or equal to v.
func (o Ord[T]) GreaterEqual(v T) bool {
	return GreaterEqual(o.v, v)
}

// String implements [fmt.Stringer].
func (o Ord[T]) String() string {
	return fmt.Sprintf("%v", any(o.v))
}

// Num is a wrapper of a numeric (or string) value T, it provides the arithmetic
// operations (such as [Num.Add]) as methods.
//
// Use [NumericOf] to wrap a value into Num.
//
// 🚀 EXAMPLE:
//
//	NumericOf(1).Add(2)      ⏩ 3
//	NumericOf("a").Add("b")  ⏩ "ab"
type Num[T constraints.Number | constraints.Complex | ~string] struct {
	v T
}

// NumericOf wraps a numeric (or string) value v into [Num].
func NumericOf[T constraints.Number | constraints.Complex | ~string](v T) Num[T] {
	return Num[T]{v: v}
}

// Get returns the wrapped value of Num.
func (n Num[T]) Get() T {
	return n.v
}

// V casts Num to [V].
func (n Num[T]) V() V[T] {
	return Of(n.v)
}

// Add adds v to the wrapped value and returns the sum.
// For string, Add performs concatenation.
//
// It is the method form of function [Add].
//
// 🚀 EXAMPLE:
//
//	NumericOf(1).Add(2)      ⏩ 3
//	NumericOf("a").Add("b")  ⏩ "ab"
func (n Num[T]) Add(v T) T {
	return Add(n.v, v)
}

// String implements [fmt.Stringer].
func (n Num[T]) String() string {
	return fmt.Sprintf("%v", any(n.v))
}
