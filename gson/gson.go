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

	"github.com/bytedance/gg/internal/conv"
)

// Valid reports whether data is a valid JSON encoding.
func Valid[V ~[]byte | ~string](data V) bool {
	switch v := any(data).(type) {
	case string: // support types like ~string
		return json.Valid(conv.StringToBytes(v))
	case []byte: // for types like []byte, ~[]bytes
		return json.Valid(v)
	default:
		// fallback for robustness: theoretically unreachable due to type constraint V ~[]byte | ~string
		return json.Valid([]byte(data))
	}
}

// Marshal returns the JSON-encoded bytes of v.
func Marshal[V any](v V) ([]byte, error) {
	return json.Marshal(v)
}

// MarshalIndent returns the JSON-encoded bytes with indent and prefix.
func MarshalIndent[V any](v V, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

// MarshalString returns the JSON-encoded string of v.
func MarshalString[V any](v V) (string, error) {
	data, err := json.Marshal(v)
	return conv.BytesToString(data), err
}

// ToString returns the JSON-encoded string of v and ignores error.
func ToString[V any](v V) string {
	data, _ := json.Marshal(v)
	return conv.BytesToString(data)
}

// ToStringIndent returns the JSON-encoded string with indent and prefix of v and ignores error.
func ToStringIndent[V any](v V, prefix, indent string) string {
	data, _ := json.MarshalIndent(v, prefix, indent)
	return conv.BytesToString(data)
}

// Unmarshal parses the JSON-encoded bytes and string and returns the result.
func Unmarshal[T any, V ~[]byte | ~string](v V) (T, error) {
	var t T
	err := json.Unmarshal([]byte(v), &t)
	return t, err
}
