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

package gresult

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"strconv"
	"testing"

	"github.com/bytedance/gg/goption"
	"github.com/bytedance/gg/internal/assert"
)

var (
	anErr = errors.New("E!R!R!O!R!")
)

func TestOf(t *testing.T) {
	assert.NotPanic(t, func() {
		_ = Of(1, nil)
	})
	assert.NotPanic(t, func() {
		_ = Of(0, nil)
	})
	assert.NotPanic(t, func() {
		_ = Of(0, anErr)
	})
	assert.Panic(t, func() {
		e := error((*fs.PathError)(nil))
		_ = Of(0, e)
	})
}

func TestErr(t *testing.T) {
	assert.Panic(t, func() {
		_ = Err[int](nil)
	})
	assert.Panic(t, func() {
		e := error((*fs.PathError)(nil))
		_ = Err[int](e)
	})
	assert.NotPanic(t, func() {
		_ = Err[int](anErr)
	})
}

func TestRValue(t *testing.T) {
	assert.Equal(t, 10, OK(10).Value())
	assert.Equal(t, 0, Err[int](anErr).Value())
	assert.Equal(t, 10, Of(10, nil).Value())
	assert.Equal(t, 10, Of(10, anErr).Value()) // 💡 NOTE: Undefined behavior
	assert.Equal(t, 0, Of(0, anErr).Value())
}

func TestRValueOr(t *testing.T) {
	assert.Equal(t, 10, OK(10).ValueOr(1))
	assert.Equal(t, 1, Err[int](anErr).ValueOr(1))
	assert.Equal(t, 10, Of(10, nil).ValueOr(1))
	assert.Equal(t, 1, Of(10, anErr).ValueOr(1)) // 💡 NOTE:
	assert.Equal(t, 1, Of(0, anErr).ValueOr(1))
}

func TestRValueOrZero(t *testing.T) {
	assert.Equal(t, 10, OK(10).ValueOrZero())
	assert.Equal(t, 0, Err[int](anErr).ValueOrZero())
	assert.Equal(t, 10, Of(10, nil).ValueOrZero())
	assert.Equal(t, 0, Of(10, anErr).ValueOrZero()) // 💡 NOTE:
	assert.Equal(t, 0, Of(0, anErr).ValueOrZero())
}

func TestRErr(t *testing.T) {
	assert.Nil(t, OK(10).Err())
	assert.Nil(t, OK(0).Err())
	assert.NotNil(t, Err[int](anErr).Err())
	assert.Nil(t, Of(10, nil).Err())
	assert.NotNil(t, Of(10, anErr).Err()) // 💡 NOTE: Undefined behavior
	assert.NotNil(t, Of(0, anErr).Err())
}

func TestRIfOK(t *testing.T) {
	assert.Panic(t, func() { OK(10).IfOK(func(int) { panic(0) }) })
	assert.NotPanic(t, func() { Err[int](anErr).IfOK(func(int) { panic(0) }) })
}

func TestRIfErr(t *testing.T) {
	assert.NotPanic(t, func() { OK(10).IfErr(func(error) { panic(0) }) })
	assert.Panic(t, func() { Err[int](anErr).IfErr(func(error) { panic(0) }) })
}

func foo1() (int, error) {
	return 1, nil
}

func foo2() R[int] {
	return OK(1)
}

func Benchmark(b *testing.B) {
	b.Run("(int,error)", func(b *testing.B) {
		for i := 0; i <= b.N; i++ {
			v, err := foo1()
			if err != nil || v != 1 {
				b.FailNow()
			}
		}
	})
	b.Run("gresult", func(b *testing.B) {
		for i := 0; i <= b.N; i++ {
			r := foo2()
			if r.Err() != nil || r.Value() != 1 {
				b.FailNow()
			}
		}
	})
}

func TestRtyp(t *testing.T) {
	assert.Equal(t, "any", Err[any](anErr).typ())
	assert.Equal(t, "int", Err[int](anErr).typ())
	assert.Equal(t, "int", OK(11).typ())
	assert.Equal(t, "int8", OK(int8(11)).typ())
	assert.Equal(t, "any", OK(any(11)).typ())
	assert.Equal(t, "any", OK[any](11).typ())
	assert.Equal(t, "any", OK[interface{}](11).typ())
	assert.Equal(t, "any", OK((interface{})(11)).typ())
}

func TestOString(t *testing.T) {

	assert.Equal(t, "gresult.Err[int](E!R!R!O!R!)", Err[int](anErr).String())
	assert.Equal(t, "gresult.OK[int](11)", OK(11).String())
	assert.Equal(t, "gresult.OK[any](11)", OK(any(11)).String())
	assert.Equal(t, "gresult.OK[int](11)", fmt.Sprintf("%s", OK(11)))
}

func TestJSON(t *testing.T) {
	// Simple
	{
		bs, err := json.Marshal(OK("test"))
		assert.Nil(t, err)
		assert.Equal(t, `{"value":"test"}`, string(bs))
	}
	{
		bs, err := json.Marshal(Err[string](errors.New("test")))
		assert.Nil(t, err)
		assert.Equal(t, `{"error":"test"}`, string(bs))
	}

	{ // Bidirect
		before := OK("test")
		bs, err := json.Marshal(before)
		assert.Nil(t, err)

		var after1 R[int]
		err = json.Unmarshal(bs, &after1)
		assert.NotNil(t, err)
		assert.Equal(t, OK(0), after1)

		var after2 R[float64]
		err = json.Unmarshal(bs, &after2)
		assert.NotNil(t, err)
		assert.Equal(t, OK(0.0), after2)

		var after3 R[string]
		err = json.Unmarshal(bs, &after3)
		assert.Nil(t, err)
		assert.Equal(t, before, after3)
	}

	{
		// Bidirect with ptr
		before := OK[*int](nil)
		bs, err := json.Marshal(before)
		assert.Nil(t, err)

		var after R[*int]
		err = json.Unmarshal(bs, &after)
		assert.Nil(t, err)
		assert.Equal(t, before, after)
	}

	{ // Unmarshal
		var r R[string]
		err := json.Unmarshal([]byte(`{"value":"test"}`), &r)
		assert.Nil(t, err)
		assert.Equal(t, OK("test"), r)
	}
	{ // Unmarshal empty: `{}`
		var r R[string]
		err := json.Unmarshal([]byte(`{}`), &r)
		assert.Nil(t, err)
		assert.Equal(t, OK(""), r)
	}
	{ // Unmarshal empty: `null`
		var r R[string]
		err := json.Unmarshal([]byte(`null`), &r)
		assert.Nil(t, err)
		assert.Equal(t, OK(""), r)
	}
	{ // Unmarshal error
		var r R[string]
		err := json.Unmarshal([]byte(`{"error":"test"}`), &r)
		assert.Nil(t, err)
		assert.Equal(t, Err[string](errors.New("test")), r)
	}
	{ // Unmarshal illegal
		var r R[string]
		err := json.Unmarshal([]byte(`{"value":"test","error":"test"}`), &r)
		assert.NotNil(t, err)
		t.Log(err)
		assert.Equal(t, OK(""), r)
	}

	{
		// Unmarshal illegal: `false`
		// Although `false` is a valid JSON,but golang does not support this.
		// FYI: https://github.com/golang/go/issues/22518
		var r R[bool]
		err := json.Unmarshal([]byte(`false`), &r)
		assert.NotNil(t, err)
		assert.Equal(t, OK(false), r)
	}
	{
		// Unmarshal illegal: ``
		// Although `` is a valid JSON,but golang does not support this.
		// FYI: https://github.com/golang/go/issues/22518
		var r R[string]
		err := json.Unmarshal([]byte(``), &r)
		assert.NotNil(t, err)
		assert.Equal(t, OK(""), r)
	}
	// Struct field
	{
		// Unmarshal JSON.Number which is larger than Number.MAX_SAFE_INTEGER(9007199254740991)
		type Foo struct {
			Bar R[json.Number] `json:"bar"`
		}

		foo := Foo{OK(json.Number("9007199254740992"))}
		bs, err := json.Marshal(foo)
		assert.Nil(t, err)
		assert.Equal(t, `{"bar":{"value":9007199254740992}}`, string(bs))

		foo1 := Foo{}
		err = json.Unmarshal(bs, &foo1)
		assert.Nil(t, err)
		assert.Equal(t, foo1, foo1)
	}
	{
		e := errors.New("test")

		type Foo struct {
			Bar R[int] `json:"bar"`
		}

		foo1 := Foo{OK(0)}
		bs1, err := json.Marshal(foo1)
		assert.Nil(t, err)
		assert.Equal(t, `{"bar":{"value":0}}`, string(bs1))

		foo2 := Foo{Err[int](e)}
		bs2, err := json.Marshal(foo2)
		assert.Nil(t, err)
		assert.Equal(t, `{"bar":{"error":"test"}}`, string(bs2))

		foo3 := Foo{}
		err = json.Unmarshal(bs1, &foo3)
		assert.Nil(t, err)
		assert.Equal(t, foo1, foo3)

		foo4 := Foo{}
		err = json.Unmarshal(bs2, &foo4)
		assert.Nil(t, err)
		assert.Equal(t, foo2, foo4)
		assert.False(t, foo2 == foo4) // different error instances

		type Fooo struct {
			Bar *R[int] `json:"bar"`
		}

		foo5 := Fooo{}
		err = json.Unmarshal(bs1, &foo5)
		assert.Nil(t, err)
		assert.Equal(t, foo1.Bar, *foo5.Bar)
	}
}

func TestRIsOK(t *testing.T) {
	assert.True(t, OK(10).IsOK())
	assert.True(t, OK(0).IsOK())
	assert.False(t, Err[int](anErr).IsOK())
	assert.True(t, Of(10, nil).IsOK())
	assert.False(t, Of(10, anErr).IsOK()) // 💡 NOTE: not recommend
	assert.False(t, Of(0, anErr).IsOK())
}

func TestRIsErr(t *testing.T) {
	assert.False(t, OK(10).IsErr())
	assert.False(t, OK(0).IsErr())
	assert.True(t, Err[int](anErr).IsErr())
	assert.False(t, Of(10, nil).IsErr())
	assert.True(t, Of(10, anErr).IsErr()) // 💡 NOTE: not recommend
	assert.True(t, Of(0, anErr).IsErr())
}

func TestRMap(t *testing.T) {
	assert.Equal(t, OK("1"), Map(OK(1), strconv.Itoa))
	assert.Equal(t, Err[string](anErr), Map(Err[int](anErr), strconv.Itoa))

	assert.NotPanic(t, func() {
		f := func(v int) string { panic("function should not be called") }
		assert.Equal(t, Err[string](anErr), Map(Err[int](anErr), f))
	})
}

func TestRMapErr(t *testing.T) {
	anotherErr := errors.New("another error")
	mapper := func(error) error {
		return anotherErr
	}
	assert.Equal(t, OK(1), MapErr(OK(1), mapper))
	assert.Equal(t, Err[int](anotherErr), MapErr(Err[int](anErr), mapper))
}

func TestRMapThen(t *testing.T) {
	anotherErr := errors.New("another error")
	do := func(v int) R[string] {
		return OK(strconv.Itoa(v))
	}
	doErr := func(v int) R[string] {
		return Err[string](anotherErr)
	}
	assert.Equal(t, OK("1"), Then(OK(1), do))
	assert.Equal(t, Err[string](anotherErr), Then(OK(1), doErr))
	assert.Equal(t, Err[string](anErr), Then(Err[int](anErr), do))
	assert.Equal(t, Err[string](anErr), Then(Err[int](anErr), doErr))

	assert.NotPanic(t, func() {
		f := func(v int) R[string] { panic("function should not be called") }
		assert.Equal(t, Err[string](anErr), Then(Err[int](anErr), f))
	})
}

func TestROptional(t *testing.T) {
	assert.Equal(t, goption.OK(1), OK(1).Option())
	assert.Equal(t, goption.Nil[int](), Err[int](anErr).Option())
}
