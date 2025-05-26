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

package gson

import (
	"testing"

	"github.com/bytedance/gg/collection/set"
	"github.com/bytedance/gg/internal/assert"
)

var codecs = map[string]Codec{
	"stdlib": JsonStdCodec,
	//"sonic.Default":     sonic.ConfigDefault,
	//"sonic.Std":         sonic.ConfigStd,
	//"json_iter.Default": jsoniter.ConfigDefault,
	//"json_iter.Compat":  jsoniter.ConfigCompatibleByStandardLibrary,
	//"json_iter.Fastest": jsoniter.ConfigFastest,
	//"sonic.Fastest":     sonic.ConfigFastest,
}

func TestValidBy(t *testing.T) {
	for name, codec := range codecs {
		t.Run(name, func(t *testing.T) {
			assert.True(t, ValidBy(codec, validJSONString))
			assert.True(t, ValidBy(codec, validJSONBytes))
			assert.True(t, ValidBy(codec, validJSONMyString))
			assert.True(t, ValidBy(codec, validJSONMyBytes))

			assert.False(t, ValidBy(codec, invalidJSONString))
			assert.False(t, ValidBy(codec, invalidJSONBytes))
			assert.False(t, ValidBy(codec, invalidJSONMyString))
			assert.False(t, ValidBy(codec, invalidJSONMyBytes))
		})
	}
}

func TestMarshalBy(t *testing.T) {
	expected := []byte(`{"name":"test","age":10}`)
	for name, codec := range codecs {
		t.Run(name, func(t *testing.T) {
			got, err := MarshalBy(codec, testcase)
			assert.Nil(t, err)
			assert.Equal(t, expected, got)
		})
	}
}

func TestMarshalIndentBy(t *testing.T) {
	expected := []byte("{\n  \"name\": \"test\",\n  \"age\": 10\n}")
	for name, codec := range codecs {
		t.Run(name, func(t *testing.T) {
			got, err := MarshalIndentBy(codec, testcase, "", "  ")
			assert.Nil(t, err)
			assert.Equal(t, expected, got)
		})
	}
}

func TestMarshalToStringBy(t *testing.T) {
	expected := `{"name":"test","age":10}`
	for name, codec := range codecs {
		t.Run(name, func(t *testing.T) {
			got, err := MarshalToStringBy(codec, testStruct{Name: "test", Age: 10})
			assert.Nil(t, err)
			assert.Equal(t, expected, got)
		})
	}
}

func TestToStringBy(t *testing.T) {
	expected := `{"name":"test","age":10}`
	for name, codec := range codecs {
		t.Run(name, func(t *testing.T) {
			got := ToStringBy(codec, testStruct{Name: "test", Age: 10})
			assert.Equal(t, expected, got)
		})
	}
}

func TestToStringIndentBy(t *testing.T) {
	expected := `{
  "name": "test",
  "age": 10
}`
	for name, codec := range codecs {
		t.Run(name, func(t *testing.T) {
			got := ToStringIndentBy(codec, testStruct{Name: "test", Age: 10}, "", "  ")
			assert.Equal(t, expected, got)
		})
	}
}

func TestUnmarshalBy(t *testing.T) {
	for _, codec := range codecs {
		{
			got, err := UnmarshalBy[testStruct](codec, ``)
			expected := testStruct{}
			assert.NotNil(t, err)
			assert.Equal(t, expected, got)
		}
		{
			got, err := UnmarshalBy[testStruct](codec, []byte(``))
			expected := testStruct{}
			assert.NotNil(t, err)
			assert.Equal(t, expected, got)
		}
		{
			got, err := UnmarshalBy[testStruct](codec, `{"name":"test","age":10}`)
			expected := testStruct{Name: "test", Age: 10}
			assert.Nil(t, err)
			assert.Equal(t, expected, got)
		}
		{
			got, err := UnmarshalBy[testStruct](codec, []byte(`{"name":"test","age":10}`))
			expected := testStruct{Name: "test", Age: 10}
			assert.Nil(t, err)
			assert.Equal(t, expected, got)
		}
		{
			got, err := UnmarshalBy[*testStruct](codec, `{"name":"test","age":10}`)
			expected := &testStruct{Name: "test", Age: 10}
			assert.Nil(t, err)
			assert.Equal(t, expected, got)
		}
		{
			got, err := UnmarshalBy[*testStruct](codec, []byte(`{"name":"test","age":10}`))
			expected := &testStruct{Name: "test", Age: 10}
			assert.Nil(t, err)
			assert.Equal(t, expected, got)
		}
		{
			got, err := UnmarshalBy[map[string]any](codec, `{"name":"test","age":10}`)
			expected := map[string]any{"name": "test", "age": float64(10)}
			assert.Nil(t, err)
			assert.Equal(t, expected, got)
		}
		{
			got, err := UnmarshalBy[[]int32](codec, `[1,2, 3]`)
			expected := []int32{1, 2, 3}
			assert.Nil(t, err)
			assert.Equal(t, expected, got)
		}
		{
			got, err := UnmarshalBy[*set.Set[int32]](codec, `[1,2, 3]`)
			expected := set.New[int32](1, 2, 3)
			assert.Nil(t, err)
			assert.Equal(t, expected, got)
		}
	}
}
