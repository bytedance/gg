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

// This file contains the method form of the pointer operations.
//
// A pointer *T can not be a method receiver, so the operations used to be
// package-level functions only. Go 1.27 generic methods make it possible to
// wrap a pointer into [P] and chain the operations, where [P.Map] declares its
// own type parameter and returns a pointer of another type:
//
//	Wrap(&i).Map(strconv.Itoa).Indirect()  ⏩ "1"
package gptr

import (
	"fmt"

	"github.com/bytedance/gg/gvalue"
)

// P is a wrapper of pointer *T, it provides the pointer operations as methods.
//
// The zero value of P is a nil pointer.
//
// Use [New] to create a P from a value, or [Wrap] to wrap an existing pointer.
//
// 🚀 EXAMPLE:
//
//	New(1).Indirect()               ⏩ 1
//	Wrap(&i).Map(strconv.Itoa).Get()  ⏩ (*string)("1")
type P[T any] struct {
	p *T
}

// New returns a P that points to equivalent value of value v. (T → *T)
//
// It is useful when you want to "convert" an unaddressable value to pointer.
//
// 💡 HINT: use [P.Indirect] to dereference the pointer (*T → T).
//
// ⚠️  WARNING: The returned pointer does not point to the original value because
// Go is always pass by value, user CAN NOT modify the value by modifying the pointer.
func New[T any](v T) P[T] {
	return P[T]{p: &v}
}

// Wrap wraps an existing pointer p into [P].
//
// 🚀 EXAMPLE:
//
//	i := 1
//	Wrap(&i).Indirect()  ⏩ 1
func Wrap[T any](p *T) P[T] {
	return P[T]{p: p}
}

// Get returns the wrapped pointer of P.
//
// 💡 AKA: Unwrap, Ptr
func (p P[T]) Get() *T {
	return p.p
}

// Indirect returns the value pointed to by the pointer p.
// If the pointer is nil, returns the zero value of T instead.
//
// It is the method form of function [Indirect].
//
// 🚀 EXAMPLE:
//
//	v := 1
//	Wrap(&v).Indirect()       ⏩ 1
//	Wrap[int](nil).Indirect() ⏩ 0
//
// 💡 AKA: Unref, Unreference, Deref, Dereference
func (p P[T]) Indirect() T {
	return Indirect(p.p)
}

// IndirectOr is a variant of [P.Indirect],
// If the pointer is nil, returns the fallback value instead.
//
// 🚀 EXAMPLE:
//
//	v := 1
//	Wrap(&v).IndirectOr(100)       ⏩ 1
//	Wrap[int](nil).IndirectOr(100) ⏩ 100
func (p P[T]) IndirectOr(fallback T) T {
	return IndirectOr(p.p, fallback)
}

// IndirectOrLazy is a variant of [P.Indirect], the fallback value is lazily
// evaluated by function f.
func (p P[T]) IndirectOrLazy(f func() T) T {
	if p.p == nil {
		return f()
	}
	return *p.p
}

// IsNil returns whether the wrapped pointer is nil.
func (p P[T]) IsNil() bool {
	return p.p == nil
}

// IsNotNil is negation of [P.IsNil].
func (p P[T]) IsNotNil() bool {
	return p.p != nil
}

// IsNilOrZero returns whether the wrapped pointer is nil or the value it points
// to is zero.
//
// It is a variant of function [IsNilOrZero] without the comparable constraint.
func (p P[T]) IsNilOrZero() bool {
	return p.p == nil || gvalue.Of(*p.p).IsZero()
}

// Clone returns a shallow copy of the value that the pointer points to.
// If the pointer is nil, nil is returned.
//
// 💡 AKA: Copy
func (p P[T]) Clone() P[T] {
	return Wrap(Clone(p.p))
}

// CloneBy is variant of [P.Clone], the value is copied using function f.
// If the pointer is nil, nil is returned.
//
// 💡 AKA: CopyBy
func (p P[T]) CloneBy(f func(T) T) P[T] {
	return p.Map(f)
}

// Equal returns whether the wrapped pointer and o are equal.
//
// Pointers are equal when either condition is satisfied:
//
//   - Both of them are nil
//   - They point to same address
//   - They point to same value (compared by [gvalue.V.Equal])
//
// 🚀 EXAMPLE:
//
//	x, y, z := 1, 1, 2
//	Wrap(&x).Equal(Wrap(&x))        ⏩ true
//	Wrap(&x).Equal(Wrap(&y))        ⏩ true
//	Wrap(&x).Equal(Wrap(&z))        ⏩ false
//	Wrap(&x).Equal(P[int]{})        ⏩ false
//	P[int]{}.Equal(P[int]{})        ⏩ true
//
// 💡 HINT: use [P.EqualTo] to compare between pointer and value.
func (p P[T]) Equal(o P[T]) bool {
	if p.p == nil || o.p == nil {
		return p.p == nil && o.p == nil
	}
	if p.p == o.p {
		return true
	}
	return gvalue.Of(*p.p).Equal(*o.p)
}

// EqualTo returns whether the value of the wrapped pointer is equal to value v.
//
// It a shortcut of "p != nil && *p == v".
//
// 🚀 EXAMPLE:
//
//	x, y := 1, 2
//	Wrap(&x).EqualTo(1)          ⏩ true
//	Wrap(&y).EqualTo(1)          ⏩ false
//	Wrap[int](nil).EqualTo(1)    ⏩ false
func (p P[T]) EqualTo(v T) bool {
	if p.p == nil {
		return false
	}
	return gvalue.Of(*p.p).Equal(v)
}

// Map applies function f to element of the wrapped pointer.
// If it is nil, f will not be called and nil is returned, otherwise,
// result of f are returned as a new pointer.
//
// It is the method form of function [Map], the result type R is a method-level
// type parameter, which is only supported by Go 1.27 generic methods.
//
// 🚀 EXAMPLE:
//
//	i := 1
//	Wrap(&i).Map(strconv.Itoa).Indirect()  ⏩ "1"
//	Wrap[int](nil).Map(strconv.Itoa)       ⏩ P[string] of nil
func (p P[T]) Map[R any](f func(T) R) P[R] {
	return Wrap(Map(p.p, f))
}

// HasZeroValue returns whether the wrapped pointer is not nil and the value it
// points to is zero.
//
// 🚀 EXAMPLE:
//
//	var i int = 0
//	Wrap(&i).HasZeroValue()        ⏩ true
//	Wrap[int](nil).HasZeroValue()  ⏩ false
func (p P[T]) HasZeroValue() bool {
	return p.p != nil && gvalue.Of(*p.p).IsZero()
}

// HasNonZeroValue returns whether the wrapped pointer is not nil and the value
// it points to is not zero.
//
// 🚀 EXAMPLE:
//
//	var i int = 1
//	Wrap(&i).HasNonZeroValue()  ⏩ true
func (p P[T]) HasNonZeroValue() bool {
	return p.p != nil && gvalue.Of(*p.p).IsNotZero()
}

// String implements [fmt.Stringer].
func (p P[T]) String() string {
	if p.p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%v", *p.p)
}
