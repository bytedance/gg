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

type testStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type (
	MyString string
	MyBytes  []byte
)

var (
	testcase = testStruct{Name: "test", Age: 10}

	validJSONString   = `{"name":"test", "age": 10}`
	validJSONBytes    = []byte(validJSONString)
	invalidJSONString = `{"name":"test", "age": 10`
	invalidJSONBytes  = []byte(invalidJSONString)

	validJSONMyString   = MyString(validJSONString)
	validJSONMyBytes    = MyBytes(validJSONString)
	invalidJSONMyString = MyString(invalidJSONString)
	invalidJSONMyBytes  = MyBytes(invalidJSONString)
)

func TestValid(t *testing.T) {
	assert.True(t, Valid(validJSONString))
	assert.True(t, Valid(validJSONBytes))
	assert.True(t, Valid(validJSONMyString))
	assert.True(t, Valid(validJSONMyBytes))
	assert.False(t, Valid(invalidJSONString))
	assert.False(t, Valid(invalidJSONBytes))
	assert.False(t, Valid(invalidJSONMyString))
	assert.False(t, Valid(invalidJSONMyBytes))
}

func TestMarshal(t *testing.T) {
	got, err := Marshal(testcase)
	expected := []byte(`{"name":"test","age":10}`)
	assert.Nil(t, err)
	assert.Equal(t, expected, got)
}

func TestMarshalIndent(t *testing.T) {
	got, err := MarshalIndent(testcase, "", "  ")
	expected := []byte("{\n  \"name\": \"test\",\n  \"age\": 10\n}")
	assert.Nil(t, err)
	assert.Equal(t, expected, got)
}

func TestMarshalString(t *testing.T) {
	got, err := MarshalString(testStruct{Name: "test", Age: 10})
	expected := `{"name":"test","age":10}`
	assert.Nil(t, err)
	assert.Equal(t, expected, got)
}

func TestToString(t *testing.T) {
	got := ToString(testStruct{Name: "test", Age: 10})
	expected := `{"name":"test","age":10}`
	assert.Equal(t, expected, got)
}

func TestToStringIndent(t *testing.T) {
	got := ToStringIndent(testStruct{Name: "test", Age: 10}, "", "  ")
	expected := `{
  "name": "test",
  "age": 10
}`
	assert.Equal(t, expected, got)
}

func TestUnmarshal(t *testing.T) {
	{
		got, err := Unmarshal[testStruct](``)
		expected := testStruct{}
		assert.NotNil(t, err)
		assert.Equal(t, expected, got)
	}
	{
		got, err := Unmarshal[testStruct]([]byte(``))
		expected := testStruct{}
		assert.NotNil(t, err)
		assert.Equal(t, expected, got)
	}
	{
		got, err := Unmarshal[testStruct](`{"name":"test","age":10}`)
		expected := testStruct{Name: "test", Age: 10}
		assert.Nil(t, err)
		assert.Equal(t, expected, got)
	}
	{
		got, err := Unmarshal[testStruct]([]byte(`{"name":"test","age":10}`))
		expected := testStruct{Name: "test", Age: 10}
		assert.Nil(t, err)
		assert.Equal(t, expected, got)
	}
	{
		got, err := Unmarshal[*testStruct](`{"name":"test","age":10}`)
		expected := &testStruct{Name: "test", Age: 10}
		assert.Nil(t, err)
		assert.Equal(t, expected, got)
	}
	{
		got, err := Unmarshal[*testStruct]([]byte(`{"name":"test","age":10}`))
		expected := &testStruct{Name: "test", Age: 10}
		assert.Nil(t, err)
		assert.Equal(t, expected, got)
	}
	{
		got, err := Unmarshal[map[string]any](`{"name":"test","age":10}`)
		expected := map[string]any{"name": "test", "age": float64(10)}
		assert.Nil(t, err)
		assert.Equal(t, expected, got)
	}
	{
		got, err := Unmarshal[[]int32](`[1,2, 3]`)
		expected := []int32{1, 2, 3}
		assert.Nil(t, err)
		assert.Equal(t, expected, got)
	}
	{
		got, err := Unmarshal[*set.Set[int32]](`[1,2, 3]`)
		expected := set.New[int32](1, 2, 3)
		assert.Nil(t, err)
		assert.Equal(t, expected, got)
	}
}
