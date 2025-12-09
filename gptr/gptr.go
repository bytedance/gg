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

// Package gptr provides generic operations for pointers.
//
// 💡 HINT: We provide similar functionality for different types in different packages.
// For example, [github.com/bytedance/gg/gslice.Clone] for copying slice while
// [github.com/bytedance/gg/gmap.Clone] for copying map.
//
//   - Use [github.com/bytedance/gg/gslice] for slice operations.
//   - Use [github.com/bytedance/gg/gmap] for map operations.
//   - Use [github.com/bytedance/gg/gvalue] for value operations.
//   - …
//
// # Operations
//
//   - Reference (T → *T): [Of], [OfNotZero], …
//   - Dereference (*T → T): [Indirect], [IndirectOr], …
//   - Predicate: (*T → bool): [Equal], [EqualTo], [IsNil], …
package gptr

import (
	"github.com/bytedance/gg/gvalue"
	"github.com/bytedance/gg/internal/constraints"
)

// Of returns a pointer that points to equivalent value of value v.
// (T → *T).
// It is useful when you want to "convert" a unaddressable value to pointer.
//
// If you need to assign the address of a literal to a pointer:
//
//	 payload := struct {
//		    Name *string
//	 }
//
// The practice without generic:
//
//	x := "name"
//	payload.Name = &x
//
// Use generic:
//
//	payload.Name = Of("name")
//
// 💡 HINT: use [Indirect] to dereference pointer (*T → T).
//
// ⚠️  WARNING: The returned pointer does not point to the original value because
// Go is always pass by value, user CAN NOT modify the value by modifying the pointer.
func Of[T any](v T) *T {
	return &v
}

// OfNotZero is variant of [Of], returns nil for zero value.
//
// 🚀 EXAMPLE:
//
//	OfNotZero(1)  ⏩ (*int)(1)
//	OfNotZero(0)  ⏩ (*int)(nil)
//
// 💡 HINT: Refer [github.com/bytedance/gg/gvalue.Zero] for definition of zero value.
func OfNotZero[T comparable](v T) *T {
	if gvalue.IsZero(v) {
		return nil
	}
	return &v
}

// OfPositive is variant of [Of], returns nil for non-positive number.
//
// 🚀 EXAMPLE:
//
//	OfPositive(1)   ⏩ (*int)(1)
//	OfPositive(0)   ⏩ (*int)(nil)
//	OfPositive(-1)  ⏩ (*int)(nil)
func OfPositive[T constraints.Number](v T) *T {
	if v <= 0 {
		return nil
	}
	return &v
}

// Indirect returns the value pointed to by the pointer p.
// If the pointer is nil, returns the zero value of T instead.
//
// 🚀 EXAMPLE:
//
//	v := 1
//	var ptrV *int = &v
//	var ptrNil *int
//	Indirect(ptrV)    ⏩ 1
//	Indirect(ptrNil)  ⏩ 0
//
// 💡 HINT: Refer [github.com/bytedance/gg/gvalue.Zero] for definition of zero value.
//
// 💡 AKA: Unref, Unreference, Deref, Dereference
func Indirect[T any](p *T) (v T) {
	if p == nil {
		// Explicitly return gvalue.Zero causes an extra copy.
		// return gvalue.Zero[T]()
		return // the initial value is zero value, see also [Indirect_gvalueZero].
	}
	return *p
}

// IndirectOr is a variant of [Indirect],
// If the pointer is nil, returns the fallback value instead.
//
// 🚀 EXAMPLE:
//
//	v := 1
//	IndirectOr(&v, 100)   ⏩ 1
//	IndirectOr(nil, 100)  ⏩ 100
func IndirectOr[T any](p *T, fallback T) T {
	if p == nil {
		return fallback
	}
	return *p
}

// IsNil returns whether the given pointer v is nil.
func IsNil[T any](p *T) bool {
	return p == nil
}

// IsNotNil is negation of [IsNil].
func IsNotNil[T any](p *T) bool {
	return p != nil
}

// Clone returns a shallow copy of the slice.
// If the given pointer is nil, nil is returned.
//
// 💡 HINT: The element is copied using assignment (=), so this is a shallow clone.
// If you want to do a deep clone, use [CloneBy] with an appropriate element
// clone function.
//
// 💡 AKA: Copy
func Clone[T any](p *T) *T {
	if p == nil {
		return nil
	}
	clone := *p
	return &clone
}

// CloneBy is variant of [Clone], it returns a copy of the map.
// Element is copied using function f.
// If the given pointer is nil, nil is returned.
//
// 💡 AKA: CopyBy
func CloneBy[T any](p *T, f func(T) T) *T {
	return Map(p, f)
}

// Equal returns whether the given pointer x and y are equal.
//
// Pointers x y are equal when either condition is satisfied:
//
//   - Both x and y is nil (x == nil && y == nil)
//   - x and y point to same address  (x == y)
//   - x and y point to same value  (*x == *y)
//
// 🚀 EXAMPLE:
//
//	x, y, z := 1, 1, 2
//	Equal(&x, &x)          ⏩ true
//	Equal(&x, &y)          ⏩ true
//	Equal(&x, &z)          ⏩ false
//	Equal(&x, nil)         ⏩ false
//	Equal[int](nil, nil)   ⏩ true
//
// 💡 HINT: use [EqualTo] to compare between pointer and value.
func Equal[T comparable](x, y *T) bool {
	if x == y {
		return true
	}
	if x == nil || y == nil {
		return false
	}
	return *x == *y
}

// EqualTo returns whether the value of pointer p is equal to value v.
//
// It a shortcut of "x != nil && *x == y".
//
// 🚀 EXAMPLE:
//
//	x, y := 1, 2
//	Equal(&x, 1)   ⏩ true
//	Equal(&y, 1)   ⏩ false
//	Equal(nil, 1)  ⏩ false
func EqualTo[T comparable](p *T, v T) bool {
	return p != nil && *p == v
}

// Map applies function f to element of pointer p.
// If p is nil, f will not be called and nil is returned, otherwise,
// result of f are returned as a new pointer.
//
// 🚀 EXAMPLE:
//
//	i := 1
//	Map(&i, strconv.Itoa)       ⏩ (*string)("1")
//	Map[int](nil, strconv.Itoa) ⏩ (*string)(nil)
func Map[F, T any](p *F, f func(F) T) *T {
	if p == nil {
		return nil
	}
	return Of(f(*p))
}

// HasZeroValue returns whether the pointer p is not nil and the value it points to is zero.
//
// 🚀 EXAMPLE:
//
//	var i int = 0
//	HasZeroValue(&i)   ⏩ true
//	HasZeroValue(nil)  ⏩ false
func HasZeroValue[T comparable](p *T) bool {
	return p != nil && gvalue.IsZero(*p)
}

// HasNonZeroValue returns whether the pointer p is not nil and the value it points to is not zero.
//
// 🚀 EXAMPLE:
//
//	var i int = 1
//	HasNonZeroValue(&i)   ⏩ true
//	var j int = 0
//	HasNonZeroValue(&j)   ⏩ false
//	HasNonZeroValue(nil)  ⏩ false
func HasNonZeroValue[T comparable](p *T) bool {
	return p != nil && gvalue.IsNotZero(*p)
}
