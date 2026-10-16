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

//go:build go1.27

package gresult

import (
	"errors"
	"io"
	"strconv"
	"testing"

	"github.com/bytedance/gg/goption"
	"github.com/bytedance/gg/internal/assert"
)

func TestRMapMethod(t *testing.T) {
	assert.Equal(t, "1", OK(1).Map(strconv.Itoa).Value())

	errResult := Err[int](io.EOF).Map(strconv.Itoa)
	assert.True(t, errResult.IsErr())
	assert.Equal(t, io.EOF, errResult.Err())
}

func TestRThenMethod(t *testing.T) {
	toString := func(i int) R[string] { return OK(strconv.Itoa(i)) }
	assert.Equal(t, "1", Of(strconv.Atoi("1")).Then(toString).Value())

	errResult := Err[int](io.EOF).Then(toString)
	assert.True(t, errResult.IsErr())
	assert.Equal(t, io.EOF, errResult.Err())
}

func TestRMapErrMethod(t *testing.T) {
	assert.Equal(t, 1, OK(1).MapErr(func(err error) error { return io.ErrUnexpectedEOF }).Value())
	assert.Equal(t, io.ErrUnexpectedEOF, Err[int](io.EOF).MapErr(func(err error) error {
		return io.ErrUnexpectedEOF
	}).Err())
}

func TestRValueOrLazy(t *testing.T) {
	assert.Equal(t, 3, Err[int](io.EOF).ValueOrLazy(func(err error) int { return len(err.Error()) }))
	assert.Equal(t, 1, OK(1).ValueOrLazy(func(err error) int { return 10 }))
}

func TestROr(t *testing.T) {
	assert.Equal(t, 2, Err[int](io.EOF).Or(OK(2)).Value())
	assert.Equal(t, 1, OK(1).Or(OK(2)).Value())
	assert.Equal(t, 2, Err[int](io.EOF).OrElse(func(err error) R[int] { return OK(2) }).Value())
	assert.Equal(t, 1, OK(1).OrElse(func(err error) R[int] { return OK(2) }).Value())
}

func TestFromOption(t *testing.T) {
	assert.Equal(t, 1, FromOption(goption.OK(1), io.EOF).Value())
	assert.True(t, FromOption(goption.Nil[int](), io.EOF).IsErr())
}

func TestRChaining(t *testing.T) {
	// The key benefit of generic methods: the value type can change in a chain.
	got := Of(strconv.Atoi("21")).
		Map(func(i int) int { return i * 2 }).
		Map(strconv.Itoa).
		ValueOr("?")
	assert.Equal(t, "42", got)

	assert.Equal(t, "?", Of(strconv.Atoi("x")).Map(strconv.Itoa).ValueOr("?"))
}

func TestRChainingError(t *testing.T) {
	err := errors.New("woof!")
	got := Err[int](err).
		Map(strconv.Itoa).
		Map(func(s string) string { return s + "!" }).
		Err()
	assert.Equal(t, err, got)
}
