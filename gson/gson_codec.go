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
//   - Validation: [ValidWith]
//   - Marshal to []byte: [MarshalWith], [MarshalIndentWith]
//   - Marshal to string: [MarshalToStringWith], [ToStringWith], [ToStringIndentWith]
//   - Unmarshal to object: [UnmarshalWith]
package gson

import (
	"encoding/json"

	"github.com/bytedance/gg/internal/conv"
)

// Codec abstracts encoding and decoding of values.
type Codec interface {
	Marshal(v any) ([]byte, error)
	MarshalIndent(v any, prefix, indent string) ([]byte, error)
	Unmarshal(data []byte, out any) error
	Valid(data []byte) bool
}

// ValidWith reports whether the input data is valid according to the given codec.
func ValidWith[V ~[]byte | ~string](codec Codec, v V) bool {
	return codec.Valid([]byte(v))
}

// MarshalWith marshals the value v into bytes using the provided codec.
//
// 🚀 Example:
//
//	MarshalWith(codec, map[string]any{"name": "test", "age": 10}) ⏩  []byte("{\"name\":\"test\",\"age\":10}")
//
// 💡 HINT: For high-performance JSON serialization, see [github.com/json-iterator/go] or [github.com/bytedance/sonic].
// Common implementations include JsonStdCodec, jsoniter.ConfigDefault, and sonic.ConfigDefault.
func MarshalWith[T any](codec Codec, v T) ([]byte, error) {
	return codec.Marshal(v)
}

// MarshalIndentWith marshals the value v into indented bytes using the provided codec.
func MarshalIndentWith[T any](codec Codec, v T, prefix, indent string) ([]byte, error) {
	return codec.MarshalIndent(v, prefix, indent)
}

// MarshalToStringWith marshals the value v into a JSON string using the provided codec.
//
// 🚀 Example:
//
//	MarshalToStringWith(codec, map[string]any{"name": "test", "age": 10}) ⏩  "{\"name\":\"test\",\"age\":10}"
//
// 💡 HINT: For high-performance JSON serialization, see [github.com/json-iterator/go] or [github.com/bytedance/sonic].
func MarshalToStringWith[T any](codec Codec, v T) (string, error) {
	data, err := codec.Marshal(v)
	return conv.BytesToString(data), err
}

// ToStringWith returns the marshaled string representation of v using the codec, ignoring errors.
func ToStringWith[T any](codec Codec, v T) string {
	data, _ := codec.Marshal(v)
	return conv.BytesToString(data)
}

// ToStringIndentWith returns the indented string representation of v using the codec, ignoring errors.
func ToStringIndentWith[T any](codec Codec, v T, prefix, indent string) string {
	data, _ := codec.MarshalIndent(v, prefix, indent)
	return conv.BytesToString(data)
}

// UnmarshalWith unmarshals the input data v into a value of type T using the provided codec.
//
// 🚀 Example:
//
//	UnmarshalWith[User](codec, `{"name":"test","age":10}`) ⏩  User{Name: "test", Age: 10}
//
// 💡 HINT: For high-performance JSON decoding, see [github.com/json-iterator/go] or [github.com/bytedance/sonic].
func UnmarshalWith[T any, V ~[]byte | ~string](codec Codec, v V) (T, error) {
	var t T
	err := codec.Unmarshal([]byte(v), &t)
	return t, err
}

// default json std lib
type stdCodec struct{}

func (stdCodec) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (stdCodec) MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

func (stdCodec) Unmarshal(data []byte, out any) error {
	return json.Unmarshal(data, out)
}

func (stdCodec) Valid(data []byte) bool {
	return json.Valid(data)
}

var JsonStdCodec Codec = stdCodec{}
