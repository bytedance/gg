// Package goption provides a generic representation of optional value.
//
// Every [O] contains a value or nothing.
//
// # Simplifying the "if v, ok := ...; ok {...}" pattern
//
// Use [os.LookupEnv] as example:
//
// The trivial way:
//
//	sh, ok := os.LookupEnv("SHELL")
//	if !ok {
//	    // Do something.
//	}
//	return sh
//
// Use optional value:
//
//	// Return zero value when the env is not present.
//	return Of(os.LookupEnv("SHELL")).Value()
//
//	// Return a custom value when the env is not present.
//	return Of(os.LookupEnv("SHELL")).ValueOr("/bin/sh")
//
// # JSON
//
// [O] implements [encoding/json.Marshaler] and [encoding/json.Ummarshaler], so
// you can use it in JSON marshaling/unmarshaling.
// See [goption.O.MarshalJSON] and [goption.O.UnmarshalJSON].
package goption

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/bytedance/gg/gvalue"
)

// O represents a generic optional value: Every O is a value T([OK]),
// or nothing([Nil]).
type O[T any] struct {
	val T
	ok  bool
}

// Of creates an optional value with type T from tuple (T, bool).
//
// Of is used to wrap result of "func () (T, bool)".
//
// 💡 NOTE: If the given bool is false, the value of T MUST be zero value of T,
// Otherwise this will be an undefined behavior.
func Of[T any](v T, ok bool) O[T] {
	return O[T]{v, ok}
}

// OK creates an optional value O containing value v.
func OK[T any](v T) O[T] {
	return Of(v, true)
}

// Nil creates an optional value O containing nothing.
func Nil[T any]() O[T] {
	return O[T]{}
}

// OfPtr is a variant of function [Of], creates an optional value from pointer v.
//
// If v != nil, returns value that the pointer points to, else returns nothing.
func OfPtr[T any](v *T) O[T] {
	if v == nil {
		return Nil[T]()
	}
	return OK(*v)
}

// Value returns internal value of O.
func (o O[T]) Value() T {
	return o.val
}

// ValueOr returns internal value of O.
// Custom value v is returned when O contains nothing.
func (o O[T]) ValueOr(v T) T {
	if o.ok {
		return o.val
	}
	return v
}

// ValueOrZero returns internal value of O.
// Zero value is returned when O contains nothing.
//
// 💡 HINT: Refer to function [github.com/bytedance/gg/gvalue.Zero]
// for explanation of zero value.
func (o O[T]) ValueOrZero() T {
	if o.ok {
		return o.val
	}
	return gvalue.Zero[T]()
}

// Ptr returns a pointer that points to the internal value of optional value O[T].
// Nil is returned when it contains nothing.
//
// 💡 NOTE: DON'T modify the internal value through the pointer,
// it won't work as you expect because the optional value is proposed to use as value,
// when you call method on it, it is copied.
func (o O[T]) Ptr() *T {
	if !o.ok {
		return nil
	}
	return &o.val
}

// Get returns the optional value in (value, ok) form.
func (o O[T]) Get() (T, bool) {
	return o.val, o.ok
}

// IsOK returns true when O contains value, otherwise false.
func (o O[T]) IsOK() bool {
	return o.ok
}

// IsNil returns true when O contains nothing, otherwise false.
func (o O[T]) IsNil() bool {
	return !o.ok
}

// IfOK calls function f when O contains value, otherwise do nothing.
func (o O[T]) IfOK(f func(T)) {
	if o.ok {
		f(o.val)
	}
}

// IfNil calls function f when O contains nil, otherwise do nothing.
func (o O[T]) IfNil(f func()) {
	if !o.ok {
		f()
	}
}

// typ returns the string representation of type of optional value.
func (o O[T]) typ() string {
	typ := reflect.TypeOf(gvalue.Zero[T]())
	if typ == nil {
		return "any"
	}
	return typ.String()
}

// String implements [fmt.Stringer].
func (o O[T]) String() string {
	if !o.ok {
		return fmt.Sprintf("goption.Nil[%s]()", o.typ())
	}
	return fmt.Sprintf("goption.OK[%s](%v)", o.typ(), o.val)
}

// MarshalJSON implements [encoding/json.Marshaler].
//
// Experimental: This API is experimental and may change in the future.
func (o O[T]) MarshalJSON() ([]byte, error) {
	if !o.ok {
		return []byte("null"), nil
	}
	return json.Marshal(o.val)
}

// UnmarshalJSON implements [encoding/json.Unmarshaler].
//
// Experimental: This API is experimental and may change in the future.
func (o *O[T]) UnmarshalJSON(data []byte) error {
	// Unmarshalers implement UnmarshalJSON([]byte("null")) as a no-op.
	if string(data) == "null" {
		return nil
	}
	if err := json.Unmarshal(data, &o.val); err != nil {
		return err
	}
	o.ok = true
	return nil
}

// Map applies function f to value of optional value O[F] if it contains value.
// Otherwise, Nil[T]() is returned.
func Map[F, T any](o O[F], f func(F) T) O[T] {
	if !o.ok {
		return Nil[T]()
	}
	return OK(f(o.val))
}

// Then calls function f and returns its result if O[F] contains value.
// Otherwise, Nil[T]() is returned.
//
// 💡 HINT: This function is similar to the Rust's std::option::Option.and_then
func Then[F, T any](o O[F], f func(F) O[T]) O[T] {
	if !o.ok {
		return Nil[T]()
	}
	return f(o.val)
}
