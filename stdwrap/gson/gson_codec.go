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
//   - Use [encoding/json] for JSONCodec
//   - Use [github.com/bytedance/sonic] for JSONCodec
//   - Use [github.com/json-iterator/go] for JSONCodec
//   - Use [github.com/vmihailenco/msgpack/v5] for MsgpackCodec
//
// # Supported Operations
//
//   - Validation: [ValidWithCodec]
//   - Marshal to []byte: [MarshalWithCodec], [MarshalIndentWithCodec]
//   - Marshal to string: [MarshalStringWithCodec], [ToStringWithCodec], [ToStringIndentWithCodec]
//   - Unmarshal to object: [UnmarshalWithCodec]
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

// ValidWithCodec reports whether data is a valid JSON encoding using the provided codec
func ValidWithCodec[V ~[]byte | ~string](v V, codec Codec) bool {
	return codec.Valid([]byte(v))
}

// MarshalWithCodec marshals the value v using the provided codec.
//
// 🚀 Example:
//
//	MarshalWithCodec(`{"name":"test","age":10}`, codec)  ⏩  []byte("{\"name\":\"test\",\"age\":10}")
//
// 💡 HINT: For high-performance JSON serialization, see [github.com/json-iterator/go] or [github.com/bytedance/sonic] implementations,
// such as JsonStdCodec, jsoniter.ConfigDefault, or sonic.ConfigDefault.
func MarshalWithCodec[T any](v T, codec Codec) ([]byte, error) {
	return codec.Marshal(v)
}

// MarshalIndentWithCodec marshals v with indent and prefix, using the provided codec
func MarshalIndentWithCodec[T any](v T, prefix, indent string, codec Codec) ([]byte, error) {
	return codec.MarshalIndent(v, prefix, indent)
}

// MarshalStringWithCodec marshals the value v using the provided codec.
//
// 🚀 Example:
//
//	MarshalStringWithCodec(`{"name":"test","age":10}`, codec)  ⏩  "{\"name\":\"test\",\"age\":10}"
//
// 💡 HINT: For high-performance JSON serialization, see [github.com/json-iterator/go] or [github.com/bytedance/sonic] implementations,
// such as JsonStdCodec, jsoniter.ConfigDefault, or sonic.ConfigDefault.
func MarshalStringWithCodec[V any](v V, codec Codec) (string, error) {
	data, err := codec.Marshal(v)
	return conv.BytesToString(data), err
}

// ToStringWithCodec returns the JSON-encoded string of v using the provided codec, and ignores error.
//
// 🚀 Example:
//
//	ToStringWithCodec(`{"name":"test","age":10}`, codec)  ⏩  "{\"name\":\"test\",\"age\":10}"
//
// 💡 HINT: For high-performance JSON serialization, see [github.com/json-iterator/go] or [github.com/bytedance/sonic] implementations,
// such as JsonStdCodec, jsoniter.ConfigDefault, or sonic.ConfigDefault.
func ToStringWithCodec[V any](v V, codec Codec) string {
	data, _ := codec.Marshal(v)
	return conv.BytesToString(data)
}

// ToStringIndentWithCodec returns the JSON-encoded string with indent and prefix of v using the provided codec, and ignores error.
func ToStringIndentWithCodec[V any](v V, prefix, indent string, codec Codec) string {
	data, _ := codec.MarshalIndent(v, prefix, indent)
	return conv.BytesToString(data)
}

// UnmarshalWithCodec unmarshals the input data v into a value of type T using the provided codec.
//
// 🚀 Example:
//
//	UnmarshalWithCodec[User](`{"name":"test","age":10}`, codec) ⏩  User{Name: "test", Age: 10}
//
// 💡 HINT: For high-performance JSON decoding, see [github.com/json-iterator/go] or [github.com/bytedance/sonic].
// Compatible implementations include JsonStdCodec, jsoniter.ConfigDefault, and sonic.ConfigDefault.
func UnmarshalWithCodec[T any, V ~[]byte | ~string](v V, codec Codec) (T, error) {
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
