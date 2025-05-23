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
	//"json_iter.Compat":  jsoniter.ConfigCompatibleWithStandardLibrary,
	//"json_iter.Fastest": jsoniter.ConfigFastest,
	//"sonic.Fastest":     sonic.ConfigFastest,
}

func TestValidWith(t *testing.T) {
	for name, codec := range codecs {
		t.Run(name, func(t *testing.T) {
			assert.True(t, ValidWith(codec, validJSONString))
			assert.True(t, ValidWith(codec, validJSONBytes))
			assert.True(t, ValidWith(codec, validJSONMyString))
			assert.True(t, ValidWith(codec, validJSONMyBytes))

			assert.False(t, ValidWith(codec, invalidJSONString))
			assert.False(t, ValidWith(codec, invalidJSONBytes))
			assert.False(t, ValidWith(codec, invalidJSONMyString))
			assert.False(t, ValidWith(codec, invalidJSONMyBytes))
		})
	}
}

func TestMarshalWith(t *testing.T) {
	expected := []byte(`{"name":"test","age":10}`)
	for name, codec := range codecs {
		t.Run(name, func(t *testing.T) {
			got, err := MarshalWith(codec, testcase)
			assert.Nil(t, err)
			assert.Equal(t, expected, got)
		})
	}
}

func TestMarshalIndentWith(t *testing.T) {
	expected := []byte("{\n  \"name\": \"test\",\n  \"age\": 10\n}")
	for name, codec := range codecs {
		t.Run(name, func(t *testing.T) {
			got, err := MarshalIndentWith(codec, testcase, "", "  ")
			assert.Nil(t, err)
			assert.Equal(t, expected, got)
		})
	}
}

func TestMarshalToStringWith(t *testing.T) {
	expected := `{"name":"test","age":10}`
	for name, codec := range codecs {
		t.Run(name, func(t *testing.T) {
			got, err := MarshalToStringWith(codec, testStruct{Name: "test", Age: 10})
			assert.Nil(t, err)
			assert.Equal(t, expected, got)
		})
	}
}

func TestToStringWith(t *testing.T) {
	expected := `{"name":"test","age":10}`
	for name, codec := range codecs {
		t.Run(name, func(t *testing.T) {
			got := ToStringWith(codec, testStruct{Name: "test", Age: 10})
			assert.Equal(t, expected, got)
		})
	}
}

func TestToStringIndentWith(t *testing.T) {
	expected := `{
  "name": "test",
  "age": 10
}`
	for name, codec := range codecs {
		t.Run(name, func(t *testing.T) {
			got := ToStringIndentWith(codec, testStruct{Name: "test", Age: 10}, "", "  ")
			assert.Equal(t, expected, got)
		})
	}
}

func TestUnmarshalWith(t *testing.T) {
	for _, codec := range codecs {
		{
			got, err := UnmarshalWith[testStruct](codec, ``)
			expected := testStruct{}
			assert.NotNil(t, err)
			assert.Equal(t, expected, got)
		}
		{
			got, err := UnmarshalWith[testStruct](codec, []byte(``))
			expected := testStruct{}
			assert.NotNil(t, err)
			assert.Equal(t, expected, got)
		}
		{
			got, err := UnmarshalWith[testStruct](codec, `{"name":"test","age":10}`)
			expected := testStruct{Name: "test", Age: 10}
			assert.Nil(t, err)
			assert.Equal(t, expected, got)
		}
		{
			got, err := UnmarshalWith[testStruct](codec, []byte(`{"name":"test","age":10}`))
			expected := testStruct{Name: "test", Age: 10}
			assert.Nil(t, err)
			assert.Equal(t, expected, got)
		}
		{
			got, err := UnmarshalWith[*testStruct](codec, `{"name":"test","age":10}`)
			expected := &testStruct{Name: "test", Age: 10}
			assert.Nil(t, err)
			assert.Equal(t, expected, got)
		}
		{
			got, err := UnmarshalWith[*testStruct](codec, []byte(`{"name":"test","age":10}`))
			expected := &testStruct{Name: "test", Age: 10}
			assert.Nil(t, err)
			assert.Equal(t, expected, got)
		}
		{
			got, err := UnmarshalWith[map[string]any](codec, `{"name":"test","age":10}`)
			expected := map[string]any{"name": "test", "age": float64(10)}
			assert.Nil(t, err)
			assert.Equal(t, expected, got)
		}
		{
			got, err := UnmarshalWith[[]int32](codec, `[1,2, 3]`)
			expected := []int32{1, 2, 3}
			assert.Nil(t, err)
			assert.Equal(t, expected, got)
		}
		{
			got, err := UnmarshalWith[*set.Set[int32]](codec, `[1,2, 3]`)
			expected := set.New[int32](1, 2, 3)
			assert.Nil(t, err)
			assert.Equal(t, expected, got)
		}
	}
}
