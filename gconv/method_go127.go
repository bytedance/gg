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

// This file contains the method form of the type conversions.
//
// The target type of a conversion can not be inferred, so it must be given
// explicitly. Before Go 1.27, that forced the conversions to be package-level
// functions with the target type as (first) type parameter:
//
//	To[int]("1")
//
// With Go 1.27 generic methods, the value can be wrapped into [V] and the target
// type becomes a method-level type parameter, which reads better in a chain:
//
//	Of("1").To[int]()
//	gslice.Wrap([]string{"1"}).Map(Of[string]).To[int]()
package gconv

import "github.com/bytedance/gg/gresult"

// V is a wrapper of value T, it provides the type conversions as methods.
//
// Use [Of] to wrap a value into V, and [V.Get] to unwrap it.
//
// 🚀 EXAMPLE:
//
//	Of("1").To[int]()     ⏩ 1
//	Of("x").ToE[int]()    ⏩ 0 strconv.ParseInt: parsing "x": invalid syntax
type V[T any] struct {
	v T
}

// Of wraps value v into [V], so that the conversions can be chained as methods.
//
// 🚀 EXAMPLE:
//
//	Of("1").To[int]()  ⏩ 1
func Of[T any](v T) V[T] {
	return V[T]{v: v}
}

// Get returns the wrapped value of V.
//
// 💡 AKA: Unwrap, Value
func (v V[T]) Get() T {
	return v.v
}

// To converts the wrapped value to the convertible type R.
// If the conversion is not supported, a zero value is returned.
//
// It is the method form of function [To].
//
// 🚀 EXAMPLE:
//
//	Of("true").To[bool]()  ⏩ true
//	Of(1).To[string]()     ⏩ "1"
//	Of("x").To[int]()      ⏩ 0
func (v V[T]) To[R convertible]() R {
	return To[R](v.v)
}

// ToPtr converts the wrapped value to a pointer of the convertible type R.
// If the conversion is not supported, nil is returned.
//
// It is the method form of function [ToPtr].
//
// 🚀 EXAMPLE:
//
//	Of("1").ToPtr[int]()  ⏩ (*int)(1)
//	Of("x").ToPtr[int]()  ⏩ (*int)(nil)
func (v V[T]) ToPtr[R convertible]() *R {
	return ToPtr[R](v.v)
}

// ToR converts the wrapped value to the convertible type R and wraps the result
// into [gresult.R], so that the error is kept.
//
// It is the method form of function [ToR].
//
// 🚀 EXAMPLE:
//
//	Of("1").ToR[int]().Value()  ⏩ 1
//	Of("x").ToR[int]().IsErr()  ⏩ true
func (v V[T]) ToR[R convertible]() gresult.R[R] {
	return ToR[R](v.v)
}

// ToE converts the wrapped value to the convertible type R and returns the
// error when the conversion is not supported.
//
// It is the method form of function [ToE].
//
// 🚀 EXAMPLE:
//
//	Of("1").ToE[int]()  ⏩ 1 nil
//	Of("x").ToE[int]()  ⏩ 0 strconv.ParseInt: parsing "x": invalid syntax
func (v V[T]) ToE[R convertible]() (R, error) {
	return ToE[R](v.v)
}
