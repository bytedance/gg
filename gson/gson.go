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

// Package gson provides operations of JSON encoding and decoding.
package gson

import (
	"encoding/json"
)

// default json std lib.
type stdJSONCodec struct{}

func (stdJSONCodec) Valid(data []byte) bool {
	return json.Valid(data)
}

func (stdJSONCodec) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (stdJSONCodec) MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

func (stdJSONCodec) Unmarshal(data []byte, out any) error {
	return json.Unmarshal(data, out)
}

var stdJSON JSONCodec = stdJSONCodec{}

// Valid reports whether data is a valid JSON encoding.
func Valid[V ~[]byte | ~string](data V) bool {
	return ValidBy(stdJSON, data)
}

// Marshal returns the JSON-encoded bytes of v.
func Marshal[V any](v V) ([]byte, error) {
	return MarshalBy(stdJSON, v)
}

// MarshalIndent returns the JSON-encoded bytes with indent and prefix.
func MarshalIndent[V any](v V, prefix, indent string) ([]byte, error) {
	return MarshalIndentBy(stdJSON, v, prefix, indent)
}

// MarshalString returns the JSON-encoded string of v.
func MarshalString[V any](v V) (string, error) {
	return MarshalStringBy(stdJSON, v)
}

// ToString returns the JSON-encoded string of v and ignores error.
func ToString[V any](v V) string {
	return ToStringBy(stdJSON, v)
}

// ToStringIndent returns the JSON-encoded string with indent and prefix of v and ignores error.
func ToStringIndent[V any](v V, prefix, indent string) string {
	return ToStringIndentBy(stdJSON, v, prefix, indent)
}

// Unmarshal parses the JSON-encoded bytes and string and returns the result.
func Unmarshal[T any, V ~[]byte | ~string](v V) (T, error) {
	return UnmarshalBy[T](stdJSON, v)
}
