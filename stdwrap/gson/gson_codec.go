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

// Package gson provides operations of JSON encoding and decoding, using the provided codec
// arbitrary values using pluggable serialization strategies (e.g., JSON, MsgPack).
//
// # Codec Abstraction
//
// The Codec interface unifies encoding and decoding logic for arbitrary types.
// You can implement this interface for various formats such as JSON, YAML, or MsgPack, Avro.
//
// For example:
//
//   - Use [encoding/json] for JSONCodec
//   - Use [github.com/bytedance/sonic] for high-performance JSONCodec
//   - Use [github.com/json-iterator/go] for customizable JSONCodec
//   - Use [github.com/vmihailenco/msgpack/v5] for MsgpackCodec
//
// # Supported Operations
//
//   - Validation: [ValidBy]
//   - Marshal to []byte: [MarshalBy], [MarshalIndentBy]
//   - Marshal to string: [MarshalStringBy], [ToStringBy], [ToStringIndentBy]
//   - Unmarshal to object: [UnmarshalBy]
package gson

import (
	"github.com/bytedance/gg/internal/conv"
)

// Marshaler defines the interface for serializing a value into a byte slice.
type Marshaler interface {
	Marshal(v any) ([]byte, error)
}

// Unmarshaler defines the interface for deserializing data into a Go value.
type Unmarshaler interface {
	Unmarshal(data []byte, out any) error
}

// PrettyMarshaler defines the interface for pretty-printing serialized output with indentation.
type PrettyMarshaler interface {
	MarshalIndent(v any, prefix, indent string) ([]byte, error)
}

// Validator defines the interface for validating whether a byte slice is a valid-encoded format.
type Validator interface {
	Valid(data []byte) bool
}

// Codec is the minimal interface for encoding and decoding values.
// It combines Marshaler and Unmarshaler.
type Codec interface {
	Marshaler
	Unmarshaler
}

// FullCodec is an extended interface for codecs that support indentation and validation.
// It combines Codec, PrettyMarshaler, and Validator.
type FullCodec interface {
	Codec
	PrettyMarshaler
	Validator
}

// ValidBy reports whether the input data is valid according to the given codec.
func ValidBy[V ~[]byte | ~string](codec FullCodec, data V) bool {
	switch v := any(data).(type) {
	case string: // support types like ~string
		return codec.Valid(conv.StringToBytes(v))
	case []byte: // for types like []byte, ~[]bytes
		return codec.Valid(v)
	default:
		// fallback for robustness: theoretically unreachable due to type constraint V ~[]byte | ~string
		return codec.Valid([]byte(data))
	}
}

// MarshalBy marshals the value v into bytes using the provided codec.
//
// 🚀 Example:
//
//	MarshalBy(codec, map[string]any{"name": "test", "age": 10}) ⏩  []byte("{\"name\":\"test\",\"age\":10}")
//
// 💡 HINT: For high-performance JSON serialization, see [github.com/json-iterator/go] or [github.com/bytedance/sonic].
// Common implementations include JsonStdCodec, jsoniter.ConfigDefault, and sonic.ConfigDefault.
func MarshalBy[T any](codec Codec, v T) ([]byte, error) {
	return codec.Marshal(v)
}

// MarshalIndentBy marshals the value v into indented bytes using the provided codec.
func MarshalIndentBy[T any](codec FullCodec, v T, prefix, indent string) ([]byte, error) {
	return codec.MarshalIndent(v, prefix, indent)
}

// MarshalStringBy marshals the value v into a JSON string using the provided codec.
//
// 🚀 Example:
//
//	MarshalStringBy(codec, map[string]any{"name": "test", "age": 10}) ⏩  "{\"name\":\"test\",\"age\":10}"
//
// 💡 HINT: For high-performance JSON serialization, see [github.com/json-iterator/go] or [github.com/bytedance/sonic].
func MarshalStringBy[T any](codec Codec, v T) (string, error) {
	data, err := codec.Marshal(v)
	return conv.BytesToString(data), err
}

// ToStringBy returns the marshaled string representation of v using the codec, ignoring errors.
func ToStringBy[T any](codec Codec, v T) string {
	data, _ := codec.Marshal(v)
	return conv.BytesToString(data)
}

// ToStringIndentBy returns the indented string representation of v using the codec, ignoring errors.
func ToStringIndentBy[T any](codec FullCodec, v T, prefix, indent string) string {
	data, _ := codec.MarshalIndent(v, prefix, indent)
	return conv.BytesToString(data)
}

// UnmarshalBy unmarshals the input data v into a value of type T using the provided codec.
//
// 🚀 Example:
//
//	UnmarshalBy[User](codec, `{"name":"test","age":10}`) ⏩  User{Name: "test", Age: 10}
//
// 💡 HINT: For high-performance JSON decoding, see [github.com/json-iterator/go] or [github.com/bytedance/sonic].
func UnmarshalBy[T any, V ~[]byte | ~string](codec Codec, v V) (T, error) {
	var t T
	err := codec.Unmarshal([]byte(v), &t)
	return t, err
}
