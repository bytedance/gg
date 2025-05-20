// Package gson provides operations of JSON encoding and decoding.
package gson

import (
	"encoding/json"

	"github.com/bytedance/gg/internal/conv"
)

// Valid reports whether data is a valid JSON encoding.
func Valid[V ~[]byte | ~string](data V) bool {
	return json.Valid([]byte(data))
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
