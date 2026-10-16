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

// This file contains the method form of the JSON operations.
//
// A value of any type can not be a method receiver, so the JSON operations used
// to be package-level functions only. Go 1.27 generic methods make it possible
// to wrap the value into [V] and chain the operations:
//
//	Of(testStruct{Name: "test"}).ToString()  ⏩ `{"name":"test","age":0}`
package gson

// V is a wrapper of value T, it provides the JSON operations as methods.
//
// Use [Of] to wrap a value into V, and [V.Get] to unwrap it.
//
// 🚀 EXAMPLE:
//
//	Of(testStruct{Name: "test", Age: 10}).ToString()  ⏩ `{"name":"test","age":10}`
type V[T any] struct {
	v T
}

// Of wraps value v into [V], so that the JSON operations can be chained as methods.
//
// 🚀 EXAMPLE:
//
//	Of(testStruct{Name: "test", Age: 10}).ToString()  ⏩ `{"name":"test","age":10}`
func Of[T any](v T) V[T] {
	return V[T]{v: v}
}

// Get returns the wrapped value of V.
//
// 💡 AKA: Unwrap, Value
func (v V[T]) Get() T {
	return v.v
}

// Marshal returns the JSON-encoded bytes of the wrapped value.
//
// It is the method form of function [Marshal].
//
// 🚀 EXAMPLE:
//
//	Of(testStruct{Name: "test", Age: 10}).Marshal()  ⏩ []byte(`{"name":"test","age":10}`) nil
func (v V[T]) Marshal() ([]byte, error) {
	return Marshal(v.v)
}

// MarshalBy returns the JSON-encoded bytes of the wrapped value with the given
// codec, such as [sonic] or [json-iterator].
//
// It is the method form of function [MarshalBy].
//
// [sonic]: https://github.com/bytedance/sonic
// [json-iterator]: https://github.com/json-iterator/go
func (v V[T]) MarshalBy(codec Marshaler) ([]byte, error) {
	return MarshalBy(codec, v.v)
}

// MarshalIndent returns the JSON-encoded bytes with indent and prefix of the
// wrapped value.
//
// It is the method form of function [MarshalIndent].
//
// 🚀 EXAMPLE:
//
//	Of(testStruct{Name: "test", Age: 10}).MarshalIndent("", "  ")
//	// "{\n  \"name\": \"test\",\n  \"age\": 10\n}" nil
func (v V[T]) MarshalIndent(prefix, indent string) ([]byte, error) {
	return MarshalIndent(v.v, prefix, indent)
}

// MarshalString returns the JSON-encoded string of the wrapped value.
//
// It is the method form of function [MarshalString].
//
// 🚀 EXAMPLE:
//
//	Of(testStruct{Name: "test", Age: 10}).MarshalString()  ⏩ `{"name":"test","age":10}` nil
func (v V[T]) MarshalString() (string, error) {
	return MarshalString(v.v)
}

// ToString returns the JSON-encoded string of the wrapped value and ignores error.
//
// It is the method form of function [ToString].
//
// 🚀 EXAMPLE:
//
//	Of(testStruct{Name: "test", Age: 10}).ToString()  ⏩ `{"name":"test","age":10}`
func (v V[T]) ToString() string {
	return ToString(v.v)
}

// ToStringIndent returns the JSON-encoded string with indent and prefix of the
// wrapped value and ignores error.
//
// It is the method form of function [ToStringIndent].
//
// 🚀 EXAMPLE:
//
//	Of(testStruct{Name: "test", Age: 10}).ToStringIndent("", "  ")
//	// "{\n  \"name\": \"test\",\n  \"age\": 10\n}"
func (v V[T]) ToStringIndent(prefix, indent string) string {
	return ToStringIndent(v.v, prefix, indent)
}
