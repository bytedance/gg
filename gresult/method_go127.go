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

// This file contains the method form of the result operations.
//
// Before Go 1.27, a method can not declare its own type parameters, so the
// operations that change the value type (such as [Map] and [Then]) had to be
// package-level functions. Go 1.27 generic methods remove that limitation, and
// the operations can be chained naturally:
//
//	Of(strconv.Atoi("1")).Map(strconv.Itoa).ValueOr("?")
//
// The methods are the preferred form, the package-level functions are kept for
// backward compatibility and delegate to them.
package gresult

import "github.com/bytedance/gg/goption"

// Map applies function f to value of result R[T] if it contains value.
// Otherwise, the error is passed to R[R2].
//
// It is the method form of function [Map].
//
// 🚀 EXAMPLE:
//
//	OK(1).Map(strconv.Itoa).Value()          ⏩ "1"
//	Err[int](io.EOF).Map(strconv.Itoa).IsErr()  ⏩ true
func (r R[T]) Map[R2 any](f func(T) R2) R[R2] {
	return Map(r, f)
}

// Then calls function f and returns its result if R[T] contains value.
// Otherwise, the error is passed to R[R2].
//
// It is the method form of function [Then].
//
// 💡 HINT: This method is similar to the Rust's std::result::Result.and_then
//
// 🚀 EXAMPLE:
//
//	Of(strconv.Atoi("1")).Then(func(i int) R[string] { return OK(strconv.Itoa(i)) }).Value()  ⏩ "1"
func (r R[T]) Then[R2 any](f func(T) R[R2]) R[R2] {
	return Then(r, f)
}

// MapErr applies function f to error of result R[T] if it contains error.
// Otherwise, the value is passed through.
//
// It is the method form of function [MapErr].
func (r R[T]) MapErr(f func(error) error) R[T] {
	return MapErr(r, f)
}

// ValueOrLazy is a variant of [R.ValueOr], the fallback value is lazily
// evaluated by function f, which accepts the internal error.
//
// 🚀 EXAMPLE:
//
//	Err[int](io.EOF).ValueOrLazy(func(err error) int { return len(err.Error()) })  ⏩ 3
func (r R[T]) ValueOrLazy(f func(error) T) T {
	if r.err == nil {
		return r.val
	}
	return f(r.err)
}

// Or returns the R itself when it contains value, otherwise the alternative is
// returned.
//
// 🚀 EXAMPLE:
//
//	Err[int](io.EOF).Or(OK(2)).Value()  ⏩ 2
//	OK(1).Or(OK(2)).Value()             ⏩ 1
func (r R[T]) Or(alternative R[T]) R[T] {
	if r.err == nil {
		return r
	}
	return alternative
}

// OrElse is a variant of [R.Or], the alternative is lazily evaluated by
// function f, which accepts the internal error.
//
// 🚀 EXAMPLE:
//
//	Err[int](io.EOF).OrElse(func(err error) R[int] { return OK(2) }).Value()  ⏩ 2
func (r R[T]) OrElse(f func(error) R[T]) R[T] {
	if r.err == nil {
		return r
	}
	return f(r.err)
}

// FromOption creates a result from an optional value (a.k.a. [goption.O]).
//
//   - If the O[T] contains value, OK(v) is returned.
//   - If the O[T] contains nothing, Err[T](err) is returned.
//
// ⚠️ WARNING: Passing a nil error will cause ❌PANIC❌!
//
// 🚀 EXAMPLE:
//
//	FromOption(goption.OK(1), io.EOF).Value()          ⏩ 1
//	FromOption(goption.Nil[int](), io.EOF).IsErr()     ⏩ true
func FromOption[T any](o goption.O[T], err error) R[T] {
	if o.IsOK() {
		return OK(o.Value())
	}
	return Err[T](err)
}
