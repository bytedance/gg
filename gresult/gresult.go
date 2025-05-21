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

// Package gresult provides a generic union type for value or error.
//
// Every [R] contains a value, representing success and containing a value,
// or an error, representing failure and containing an error.
//
// # Simplifying the "if v, err := ...; err != nil {...}" pattern
//
// Use [os.Open] as example:
//
// The trivial way:
//
//	fd, err := os.Open("/tmp/error.log")
//	if err != nil {
//	    // Do something.
//	}
//	return fd
//
// Use result:
//
//	// Return zero value when the file is not present.
//	return Of(os.Open("/tmp/error.log")).Value()
//
//	// Return a custom file object when the file is not present.
//	return Of(os.Open("/tmp/error.log")).ValueOr(os.Stderr)
//
// # JSON
//
// [R] implements [encoding/json.Marshaler] and [encoding/json.Unmarshaler], so
// you can use it in JSON marshaling/unmarshaling.
// See [gresult.R.MarshalJSON] and [gresult.R.UnmarshalJSON].
package gresult

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/bytedance/gg/goption"
	"github.com/bytedance/gg/gptr"
	"github.com/bytedance/gg/gvalue"
)

// R represents a generic result: Every result is a value,
// representing success and containing a value ([OK]),
// or an error, representing error and containing an error ([Err]).
type R[T any] struct {
	val T
	err error
}

// Of creates a result with type T from tuple (T, error).
//
// Of is often used to wrap result of func () (T, error)`,
// see README.md for detailed usage.
//
// 💡 NOTE: If the given error is not nil, the value of T MUST be zero value,
// Otherwise this will be an undefined behavior.
//
// ⚠️ WARNING: Passing a nil error with type (such as (*fs.PathError)(nil)) will
// cause ❌PANIC❌!
func Of[T any](v T, e error) R[T] {
	checkErr(e)
	return R[T]{v, e}
}

// OK creates a result containing value v.
func OK[T any](v T) R[T] {
	return R[T]{val: v}
}

// Err creates a result containing error e.
//
// ⚠️ WARNING: Passing a nil error will cause ❌PANIC❌!
func Err[T any](e error) R[T] {
	if gvalue.IsNil(e) {
		panic(fmt.Errorf("expected a non-nil error, but found nil error with type %T", e))
	}
	return R[T]{err: e}
}

// Value returns the internal value of R.
func (r R[T]) Value() T {
	return r.val
}

// ValueOr returns the internal value of R.
// Custom value v is returned when the result contains error.
func (r R[T]) ValueOr(v T) T {
	if r.err == nil {
		return r.val
	}
	return v
}

// ValueOrZero returns the internal value of R.
// Custom value v is returned when the result contains error.
//
// 💡 HINT: Refer to function [github.com/bytedance/gg/gvalue.Zero]
// for explanation of zero value.
func (r R[T]) ValueOrZero() T {
	if r.err == nil {
		return r.val
	}
	return gvalue.Zero[T]()
}

// Err returns the internal error of R.
func (r R[T]) Err() error {
	return r.err
}

// Get returns the result in (value, error) form.
func (r R[T]) Get() (T, error) {
	return r.val, r.err
}

// IsOK returns true when R contains value, otherwise false.
func (r R[T]) IsOK() bool {
	return r.err == nil
}

// IsErr true when R contains error, otherwise false.
func (r R[T]) IsErr() bool {
	return r.err != nil
}

// IfOK calls functions f when R contains value.
func (r R[T]) IfOK(f func(T)) {
	if r.err == nil {
		f(r.val)
	}
}

// IfErr calls functions f when R contains error.
func (r R[T]) IfErr(f func(error)) {
	if r.err != nil {
		f(r.err)
	}
}

func checkErr(e error) {
	if e != nil && gvalue.IsNil(e) {
		panic(fmt.Errorf("error with type %T is nil", e))
	}
}

// typ returns the string representation of type of optional value.
func (r R[T]) typ() string {
	typ := reflect.TypeOf(gvalue.Zero[T]())
	if typ == nil {
		return "any"
	}
	return typ.String()
}

// String implements fmt.Stringer.
func (r R[T]) String() string {
	if r.err != nil {
		return fmt.Sprintf("gresult.Err[%s](%s)", r.typ(), r.err)
	}
	return fmt.Sprintf("gresult.OK[%s](%v)", r.typ(), r.val)
}

type jsonR[T any] struct {
	Val *T      `json:"value,omitempty"`
	Err *string `json:"error,omitempty"`
}

// MarshalJSON implements [encoding/json.Marshaler].
//
// Experimental: This API is experimental and may change in the future.
func (r R[T]) MarshalJSON() ([]byte, error) {
	jr := jsonR[T]{}
	if r.err != nil {
		jr.Err = gptr.Of(r.err.Error())
	} else {
		jr.Val = &r.val
	}
	return json.Marshal(jr)
}

// UnmarshalJSON implements [encoding/json.Unmarshaler].
//
// ⚠️ WARNING: After unmarshaling, user MUST NOT make any assumptions about type
// type of [R.Err].
//
// Experimental: This API is experimental and may change in the future.
func (r *R[T]) UnmarshalJSON(data []byte) error {
	// Unmarshalers implement UnmarshalJSON([]byte("null")) as a no-op.
	if string(data) == "null" {
		return nil
	}
	jr := jsonR[T]{}
	if err := json.Unmarshal(data, &jr); err != nil {
		return err
	}

	// Deal with illegal JSON payload.
	if jr.Err != nil && jr.Val != nil {
		return errors.New("gresult: neither error nor value is nil")
	}

	if jr.Err == nil && jr.Val == nil {
		r.val = gvalue.Zero[T]()
		r.err = nil
	} else if jr.Err != nil {
		r.val = gvalue.Zero[T]()
		r.err = errors.New(*jr.Err)
	} else {
		r.val = *jr.Val
		r.err = nil
	}

	return nil
}

// Map applies function f to value of result R[F] if it contains value.
// Otherwise, passes the error of result R[F] to R[T].
func Map[F, T any](r R[F], f func(F) T) R[T] {
	if r.err != nil {
		return Err[T](r.err)
	}
	return OK(f(r.val))
}

// MapErr applies function f to error of result R[F] if it contains error.
// Otherwise, passes the value of result R[F] to R[T].
func MapErr[T any](r R[T], f func(error) error) R[T] {
	if r.err == nil {
		return r
	}
	return Err[T](f(r.err))
}

// Then calls function f and returns its result if R[F] contains value.
// Otherwise, passes the error of result R[F] to R[T].
//
// 💡 HINT: This function is similar to the Rust's std::result::Result.and_then
func Then[F, T any](r R[F], f func(F) R[T]) R[T] {
	if r.err != nil {
		return Err[T](r.err)
	}
	return f(r.val)
}

// Option converts result to an optional value (a.k.a. [goption.O]).
// So user can user methods provided by [goption] package.
//
//   - If the R[T] is a value, goption.OK(v) is returned.
//   - If the R[T] is an error, it will be dropped, and goption.Nil[T]() is returned.
//
// 🚀 EXAMPLE:
//
//	OK(1).Option()         ⏩ goption.OK(1)
//	err := errors.New("woof!")
//	Err[int](err).Option() ⏩ goption.Nil[int]()
func (r R[T]) Option() goption.O[T] {
	if r.err != nil {
		return goption.Nil[T]()
	}
	return goption.OK(r.val)
}
