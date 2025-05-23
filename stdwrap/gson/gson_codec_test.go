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

func TestValidWithCodec(t *testing.T) {

	for name, codec := range codecs {
		t.Run(name, func(t *testing.T) {
			assert.True(t, ValidWithCodec(validJSONString, codec))
			assert.True(t, ValidWithCodec(validJSONBytes, codec))
			assert.True(t, ValidWithCodec(validJSONMyString, codec))
			assert.True(t, ValidWithCodec(validJSONMyBytes, codec))

			assert.False(t, ValidWithCodec(invalidJSONString, codec))
			assert.False(t, ValidWithCodec(invalidJSONBytes, codec))
			assert.False(t, ValidWithCodec(invalidJSONMyString, codec))
			assert.False(t, ValidWithCodec(invalidJSONMyBytes, codec))
		})
	}

}

func TestMarshalWithCodec(t *testing.T) {
	expected := []byte(`{"name":"test","age":10}`)
	for name, codec := range codecs {
		t.Run(name, func(t *testing.T) {
			got, err := MarshalWithCodec(testcase, codec)
			assert.Nil(t, err)
			assert.Equal(t, expected, got)
		})
	}
}

func TestMarshalIndentWithCodec(t *testing.T) {
	expected := []byte("{\n  \"name\": \"test\",\n  \"age\": 10\n}")
	for name, codec := range codecs {
		t.Run(name, func(t *testing.T) {
			got, err := MarshalIndentWithCodec(testcase, "", "  ", codec)
			assert.Nil(t, err)
			assert.Equal(t, expected, got)
		})
	}
}

func TestMarshalStringWithCodec(t *testing.T) {
	expected := `{"name":"test","age":10}`
	for name, codec := range codecs {
		t.Run(name, func(t *testing.T) {
			got, err := MarshalStringWithCodec(testStruct{Name: "test", Age: 10}, codec)
			assert.Nil(t, err)
			assert.Equal(t, expected, got)
		})
	}
}

func TestToStringWithCodec(t *testing.T) {
	expected := `{"name":"test","age":10}`
	for name, codec := range codecs {
		t.Run(name, func(t *testing.T) {
			got := ToStringWithCodec(testStruct{Name: "test", Age: 10}, codec)
			assert.Equal(t, expected, got)
		})
	}
}

func TestToStringIndentWithCodec(t *testing.T) {
	expected := `{
  "name": "test",
  "age": 10
}`
	for name, codec := range codecs {
		t.Run(name, func(t *testing.T) {
			got := ToStringIndentWithCodec(testStruct{Name: "test", Age: 10}, "", "  ", codec)
			assert.Equal(t, expected, got)
		})
	}
}

func TestUnmarshalWithCodec(t *testing.T) {

	//UnmarshalWithCodec[testStruct](``,  sonic.ConfigFastest)

	for _, codec := range codecs {
		{
			got, err := UnmarshalWithCodec[testStruct](``, codec)
			expected := testStruct{}
			assert.NotNil(t, err)
			assert.Equal(t, expected, got)
		}
		{
			got, err := UnmarshalWithCodec[testStruct]([]byte(``), codec)
			expected := testStruct{}
			assert.NotNil(t, err)
			assert.Equal(t, expected, got)
		}
		{
			got, err := UnmarshalWithCodec[testStruct](`{"name":"test","age":10}`, codec)
			expected := testStruct{Name: "test", Age: 10}
			assert.Nil(t, err)
			assert.Equal(t, expected, got)
		}
		{
			got, err := UnmarshalWithCodec[testStruct]([]byte(`{"name":"test","age":10}`), codec)
			expected := testStruct{Name: "test", Age: 10}
			assert.Nil(t, err)
			assert.Equal(t, expected, got)
		}
		{
			got, err := UnmarshalWithCodec[*testStruct](`{"name":"test","age":10}`, codec)
			expected := &testStruct{Name: "test", Age: 10}
			assert.Nil(t, err)
			assert.Equal(t, expected, got)
		}
		{
			got, err := UnmarshalWithCodec[*testStruct]([]byte(`{"name":"test","age":10}`), codec)
			expected := &testStruct{Name: "test", Age: 10}
			assert.Nil(t, err)
			assert.Equal(t, expected, got)
		}
		{
			got, err := UnmarshalWithCodec[map[string]any](`{"name":"test","age":10}`, codec)
			expected := map[string]any{"name": "test", "age": float64(10)}
			assert.Nil(t, err)
			assert.Equal(t, expected, got)
		}
		{
			got, err := UnmarshalWithCodec[[]int32](`[1,2, 3]`, codec)
			expected := []int32{1, 2, 3}
			assert.Nil(t, err)
			assert.Equal(t, expected, got)
		}
		{
			got, err := UnmarshalWithCodec[*set.Set[int32]](`[1,2, 3]`, codec)
			expected := set.New[int32](1, 2, 3)
			assert.Nil(t, err)
			assert.Equal(t, expected, got)
		}
	}

}
