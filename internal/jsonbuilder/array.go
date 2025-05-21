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

package jsonbuilder

import (
	"bytes"
	"encoding/json"

	"github.com/bytedance/gg/internal/stream"
)

// Array is a builder for building JSON array.
type Array struct {
	elems [][]byte
	size  int
}

func NewArray() *Array {
	return &Array{}
}

func (a *Array) Append(v any) error {
	bs, err := json.Marshal(v)
	if err != nil {
		return err
	}
	a.elems = append(a.elems, bs)
	a.size += len(bs)
	return nil
}

func (a *Array) Sort() {
	a.elems = stream.
		StealSlice(a.elems).
		SortBy(func(a, b []byte) bool { return bytes.Compare(a, b) == -1 }).
		ToSlice()
}

func (a *Array) Build() ([]byte, error) {
	if a == nil {
		return []byte("null"), nil
	}

	var buf bytes.Buffer

	size := a.size
	size += len(a.elems) - 1 // count of comma ","
	size += 2                //  "[" and "]"
	buf.Grow(size)

	buf.WriteByte('[')
	for _, bs := range a.elems {
		buf.Write(bs)
		buf.WriteByte(',')
	}

	var out []byte
	if len(a.elems) == 0 {
		buf.WriteByte(']')
		out = buf.Bytes()
	} else {
		extraComma := buf.Len() - 1
		out = buf.Bytes()
		// Replace extra `,` to  `]`.
		out[extraComma] = ']'
	}

	return out, nil
}
