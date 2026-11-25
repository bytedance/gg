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

// This file contains the method form of the optional value operations.
//
// Before Go 1.27, a method can not declare its own type parameters, so the
// operations that change the value type (such as [Map] and [Then]) had to be
// package-level functions. Go 1.27 generic methods remove that limitation, and
// the operations can be chained naturally:
//
//	OK(1).Map(strconv.Itoa).Then(parse).ValueOr("?")
//
// The methods are the preferred form, the package-level functions are kept for
// backward compatibility and delegate to them.
package goption

// Map applies function f to value of optional value O[T] if it contains value.
// Otherwise, Nil[R]() is returned.
//
// It is the method form of function [Map].
//
// 🚀 EXAMPLE:
//
//	OK(1).Map(strconv.Itoa).Value()       ⏩ "1"
//	Nil[int]().Map(strconv.Itoa).IsNil()  ⏩ true
func (o O[T]) Map[R any](f func(T) R) O[R] {
	return Map(o, f)
}

// MapOr applies function f to value of O[T] if it contains value.
// Otherwise, the fallback value is returned.
//
// 🚀 EXAMPLE:
//
//	OK(1).MapOr(strconv.Itoa, "?")       ⏩ "1"
//	Nil[int]().MapOr(strconv.Itoa, "?")  ⏩ "?"
func (o O[T]) MapOr[R any](f func(T) R, fallback R) R {
	if !o.ok {
		return fallback
	}
	return f(o.val)
}

// Then calls function f and returns its result if O[T] contains value.
// Otherwise, Nil[R]() is returned.
//
// It is the method form of function [Then].
//
// 💡 HINT: This method is similar to the Rust's std::option::Option.and_then
//
// 🚀 EXAMPLE:
//
//	OK(1).Then(func(i int) O[string] { return OK(strconv.Itoa(i)) }).Value()  ⏩ "1"
func (o O[T]) Then[R any](f func(T) O[R]) O[R] {
	return Then(o, f)
}

// ValueOrLazy is a variant of [O.ValueOr], the fallback value is lazily
// evaluated by function f.
//
// 🚀 EXAMPLE:
//
//	OK(1).ValueOrLazy(func() int { return 10 })  ⏩ 1
//	Nil[int]().ValueOrLazy(func() int { return 10 })  ⏩ 10
func (o O[T]) ValueOrLazy(f func() T) T {
	if o.ok {
		return o.val
	}
	return f()
}

// Filter returns the O itself when it contains value and f returns true for the
// value, otherwise Nil[T]() is returned.
//
// 🚀 EXAMPLE:
//
//	OK(1).Filter(func(i int) bool { return i%2 == 0 }).IsNil()  ⏩ true
//	OK(2).Filter(func(i int) bool { return i%2 == 0 }).Value()  ⏩ 2
func (o O[T]) Filter(f func(T) bool) O[T] {
	if !o.ok || !f(o.val) {
		return Nil[T]()
	}
	return o
}

// Or returns the O itself when it contains value, otherwise the alternative is
// returned.
//
// 🚀 EXAMPLE:
//
//	Nil[int]().Or(OK(2)).Value()  ⏩ 2
//	OK(1).Or(OK(2)).Value()       ⏩ 1
func (o O[T]) Or(alternative O[T]) O[T] {
	if o.ok {
		return o
	}
	return alternative
}

// OrElse is a variant of [O.Or], the alternative is lazily evaluated by
// function f.
//
// 🚀 EXAMPLE:
//
//	Nil[int]().OrElse(func() O[int] { return OK(2) }).Value()  ⏩ 2
func (o O[T]) OrElse(f func() O[T]) O[T] {
	if o.ok {
		return o
	}
	return f()
}

// ToSlice converts O to a slice: it contains the value when O contains value,
// otherwise it is empty.
//
// 🚀 EXAMPLE:
//
//	OK(1).ToSlice()       ⏩ [1]
//	Nil[int]().ToSlice()  ⏩ []
func (o O[T]) ToSlice() []T {
	if !o.ok {
		return nil
	}
	return []T{o.val}
}
