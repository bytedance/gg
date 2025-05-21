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

package goption

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/bytedance/gg/gptr"
	"github.com/bytedance/gg/internal/assert"
)

func TestOf(t *testing.T) {
	_ = Of(1, true)
	_ = OfPtr((*int)(nil))
	_ = OfPtr(gptr.Of(1))
}

func TestOValue(t *testing.T) {
	assert.Equal(t, 10, OK(10).Value())
	assert.Equal(t, 0, Nil[int]().Value())
	assert.Equal(t, 10, Of(10, true).Value())
	assert.Equal(t, 10, Of(10, false).Value()) // 💡 NOTE: not recommend
	assert.Equal(t, 0, Of(0, false).Value())

	assert.Equal(t, 1, OfPtr(gptr.Of(1)).Value())
	assert.Equal(t, 0, OfPtr((*int)(nil)).Value())
}

func TestOValueOr(t *testing.T) {
	assert.Equal(t, 10, OK(10).ValueOr(1))
	assert.Equal(t, 1, Nil[int]().ValueOr(1))
	assert.Equal(t, 10, Of(10, true).ValueOr(1))
	assert.Equal(t, 1, Of(10, false).ValueOr(1)) // 💡 NOTE: not recommend
	assert.Equal(t, 1, Of(0, false).ValueOr(1))
}

func TestOValueOrZero(t *testing.T) {
	assert.Equal(t, 10, OK(10).ValueOrZero())
	assert.Equal(t, 0, Nil[int]().ValueOrZero())
	assert.Equal(t, 10, Of(10, true).ValueOrZero())
	assert.Equal(t, 0, Of(10, false).ValueOrZero()) // 💡 NOTE: not recommend
	assert.Equal(t, 0, Of(0, false).ValueOrZero())
}

func TestOOK(t *testing.T) {
	assert.True(t, OK(10).IsOK())
	assert.True(t, OK(0).IsOK())
	assert.False(t, Nil[int]().IsOK())
	assert.True(t, Of(10, true).IsOK())
	assert.False(t, Of(10, false).IsOK()) // 💡 NOTE: not recommend
	assert.False(t, Of(0, false).IsOK())
}

func TestOIfOK(t *testing.T) {
	assert.Panic(t, func() { OK(10).IfOK(func(int) { panic(0) }) })
	assert.NotPanic(t, func() { Nil[int]().IfOK(func(int) { panic(0) }) })
}

func TestOIfNil(t *testing.T) {
	assert.NotPanic(t, func() { OK(10).IfNil(func() { panic(0) }) })
	assert.Panic(t, func() { Nil[int]().IfNil(func() { panic(0) }) })
}

func foo1() (int, bool) {
	return 1, true
}

func foo2() O[int] {
	return OK(1)
}

func Benchmark(b *testing.B) {
	b.Run("(int,bool)", func(b *testing.B) {
		for i := 0; i <= b.N; i++ {
			v, ok := foo1()
			if !ok || v != 1 {
				b.FailNow()
			}
		}
	})
	b.Run("goption", func(b *testing.B) {
		for i := 0; i <= b.N; i++ {
			o := foo2()
			if !o.IsOK() || o.Value() != 1 {
				b.FailNow()
			}
		}
	})
}

func TestOtyp(t *testing.T) {
	assert.Equal(t, "any", Nil[any]().typ())
	assert.Equal(t, "int", Nil[int]().typ())
	assert.Equal(t, "int", OK(11).typ())
	assert.Equal(t, "int8", OK(int8(11)).typ())
	assert.Equal(t, "any", OK(any(11)).typ())
	assert.Equal(t, "any", OK[any](11).typ())
	assert.Equal(t, "any", OK[interface{}](11).typ())
	assert.Equal(t, "any", OK((interface{})(11)).typ())
}

func TestOString(t *testing.T) {
	assert.Equal(t, "goption.Nil[int]()", O[int]{}.String())
	assert.Equal(t, "goption.Nil[int]()", Nil[int]().String())
	assert.Equal(t, "goption.OK[int](11)", OK(11).String())
	assert.Equal(t, "goption.OK[any](11)", OK(any(11)).String())
	assert.Equal(t, "goption.OK[int](11)", fmt.Sprintf("%s", OK(11)))
}

func TestJSON(t *testing.T) {
	{
		var v *int
		expect, _ := json.Marshal(v)
		actual, _ := json.Marshal(OfPtr(v))
		assert.Equal(t, string(expect), string(actual))
	}
	{
		v := gptr.Of(1)
		expect, _ := json.Marshal(v)
		actual, _ := json.Marshal(OfPtr(v))
		assert.Equal(t, string(expect), string(actual))
	}
	{
		v := gptr.Of("test")
		expect, _ := json.Marshal(v)
		actual, _ := json.Marshal(OfPtr(v))
		assert.Equal(t, string(expect), string(actual))
	}

	// Simple.
	{
		bs, err := json.Marshal(OK("test"))
		assert.Nil(t, err)
		assert.Equal(t, `"test"`, string(bs))
	}
	{
		bs, err := json.Marshal(Nil[string]())
		assert.Nil(t, err)
		assert.Equal(t, `null`, string(bs))
	}

	{ // Bidirect
		before := OK("test")
		bs, err := json.Marshal(before)
		assert.Nil(t, err)

		var after1 O[int]
		err = json.Unmarshal(bs, &after1)
		assert.NotNil(t, err)
		assert.Equal(t, Nil[int](), after1)

		var after2 O[float64]
		err = json.Unmarshal(bs, &after2)
		assert.NotNil(t, err)
		assert.Equal(t, Nil[float64](), after2)

		var after3 O[string]
		err = json.Unmarshal(bs, &after3)
		assert.Nil(t, err)
		assert.Equal(t, before, after3)
	}

	{ // Unmarshal
		var o O[string]
		err := json.Unmarshal([]byte(`"test"`), &o)
		assert.Nil(t, err)
		assert.Equal(t, OK("test"), o)
	}
	{ // Unmarshal nil
		var o O[string]
		err := json.Unmarshal([]byte(`null`), &o)
		assert.Nil(t, err)
		assert.Equal(t, Nil[string](), o)
	}

	// Struct field
	{
		type Foo struct {
			Bar O[int] `json:"bar"`
		}

		foo1 := Foo{OK(0)}
		bs1, err := json.Marshal(foo1)
		assert.Nil(t, err)
		assert.Equal(t, `{"bar":0}`, string(bs1))

		foo2 := Foo{}
		bs2, err := json.Marshal(foo2)
		assert.Nil(t, err)
		assert.Equal(t, `{"bar":null}`, string(bs2))

		foo3 := Foo{}
		err = json.Unmarshal(bs1, &foo3)
		assert.Nil(t, err)
		assert.Equal(t, foo1, foo3)

		foo4 := Foo{}
		err = json.Unmarshal(bs2, &foo4)
		assert.Nil(t, err)
		assert.Equal(t, foo2, foo4)

		type Fooo struct {
			Bar *O[int] `json:"bar"`
		}

		foo5 := Fooo{}
		err = json.Unmarshal(bs1, &foo5)
		assert.Nil(t, err)
		assert.Equal(t, foo1.Bar, *foo5.Bar)

		foo6 := Fooo{}
		err = json.Unmarshal(bs2, &foo6)
		assert.Nil(t, err)
		assert.True(t, foo6.Bar == nil)
	}
}

func TestOIsOK(t *testing.T) {
	assert.True(t, OK(0).IsOK())
	assert.False(t, Nil[int]().IsOK())
	assert.True(t, Of(10, true).IsOK())
	assert.False(t, Of(10, false).IsOK()) // 💡 NOTE: not recommend
	assert.False(t, Of(0, false).IsOK())
}

func TestOIsNil(t *testing.T) {
	assert.False(t, OK(10).IsNil())
	assert.False(t, OK(0).IsNil())
	assert.True(t, Nil[int]().IsNil())
	assert.False(t, Of(10, true).IsNil())
	assert.True(t, Of(10, false).IsNil()) // 💡 NOTE: not recommend
	assert.True(t, Of(0, false).IsNil())
}

func TestO_Alias(t *testing.T) {
	assert.True(t, OK(10).IsOK())
	assert.Panic(t, func() { OK(10).IfOK(func(int) { panic(0) }) })
	assert.NotPanic(t, func() { Nil[int]().IfOK(func(int) { panic(0) }) })
}

func TestOMap(t *testing.T) {
	assert.Equal(t, OK("1"), Map(OK(1), strconv.Itoa))
	assert.Equal(t, Nil[string](), Map(Nil[int](), strconv.Itoa))

	assert.NotPanic(t, func() {
		f := func(v int) string { panic("function should not be called") }
		assert.Equal(t, Nil[string](), Map(Nil[int](), f))
	})
}

func TestOThen(t *testing.T) {
	do := func(v int) O[string] {
		return OK(strconv.Itoa(v))
	}
	doNil := func(v int) O[string] {
		return Nil[string]()
	}
	assert.Equal(t, OK("1"), Then(OK(1), do))
	assert.Equal(t, Nil[string](), Then(OK(1), doNil))
	assert.Equal(t, Nil[string](), Then(Nil[int](), do))
	assert.Equal(t, Nil[string](), Then(Nil[int](), doNil))

	assert.NotPanic(t, func() {
		f := func(v int) O[string] { panic("function should not be called") }
		assert.Equal(t, Nil[string](), Then(Nil[int](), f))
	})
}

func TestOPtr(t *testing.T) {
	assert.Equal(t, gptr.Of(10), OK(10).Ptr())
	assert.Equal(t, nil, Nil[int]().Ptr())

	// Test modify.
	{
		o := OK(10)
		ptr := o.Ptr()
		*ptr = 1
		assert.Equal(t, OK(10), o)
		assert.True(t, o.Ptr() != o.Ptr()) // o is copied
	}
}
