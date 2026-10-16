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

package goption

import (
	"strconv"
	"testing"

	"github.com/bytedance/gg/internal/assert"
)

func TestOMapMethod(t *testing.T) {
	assert.Equal(t, "1", OK(1).Map(strconv.Itoa).Value())
	assert.True(t, Nil[int]().Map(strconv.Itoa).IsNil())
	assert.Equal(t, "1", OK(1).MapOr(strconv.Itoa, "?"))
	assert.Equal(t, "?", Nil[int]().MapOr(strconv.Itoa, "?"))
}

func TestOThenMethod(t *testing.T) {
	toString := func(i int) O[string] { return OK(strconv.Itoa(i)) }
	assert.Equal(t, "1", OK(1).Then(toString).Value())
	assert.True(t, Nil[int]().Then(toString).IsNil())
}

func TestOValueOrLazy(t *testing.T) {
	assert.Equal(t, 10, Nil[int]().ValueOrLazy(func() int { return 10 }))
	assert.Equal(t, 1, OK(1).ValueOrLazy(func() int { return 10 }))
}

func TestOFilter(t *testing.T) {
	isEven := func(i int) bool { return i%2 == 0 }
	assert.Equal(t, 2, OK(2).Filter(isEven).Value())
	assert.True(t, OK(1).Filter(isEven).IsNil())
	assert.True(t, Nil[int]().Filter(isEven).IsNil())
}

func TestOOr(t *testing.T) {
	assert.Equal(t, 2, Nil[int]().Or(OK(2)).Value())
	assert.Equal(t, 1, OK(1).Or(OK(2)).Value())
	assert.Equal(t, 2, Nil[int]().OrElse(func() O[int] { return OK(2) }).Value())
	assert.Equal(t, 1, OK(1).OrElse(func() O[int] { return OK(2) }).Value())
}

func TestOToSlice(t *testing.T) {
	assert.Equal(t, []int{1}, OK(1).ToSlice())
	assert.Equal(t, 0, len(Nil[int]().ToSlice()))
}

func TestOChaining(t *testing.T) {
	// The key benefit of generic methods: the value type can change in a chain.
	got := OK(1).
		Map(func(i int) int { return i * 2 }).
		Filter(func(i int) bool { return i > 1 }).
		Map(strconv.Itoa).
		ValueOr("?")
	assert.Equal(t, "2", got)

	assert.Equal(t, "?", OK(1).Map(strconv.Itoa).Filter(func(s string) bool { return s == "2" }).MapOr(
		func(s string) string { return s }, "?"))
}
