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

package gfunc

import (
	"testing"

	"github.com/bytedance/gg/internal/assert"
)

func TestPartial1(t *testing.T) {
	add := Partial1(func(a int) int {
		return a
	})
	assert.Equal(t, 1, add.Partial(1)())
	assert.Equal(t, 1, add.PartialR(1)())
}

func TestPartial2(t *testing.T) {
	add := Partial2(func(a, b int) int {
		return a + b
	})
	assert.Equal(t, 3, add.Partial(1).Partial(2)())
	assert.Equal(t, 3, add.PartialR(1).PartialR(2)())
}

func TestPartial3(t *testing.T) {
	add := Partial3(func(a, b, c int) int {
		return a + b + c
	})
	assert.Equal(t, 6, add.Partial(1)(2, 3))
	assert.Equal(t, 6, add.PartialR(1)(2, 3))
}

func TestPartial4(t *testing.T) {
	add := Partial4(func(a, b, c, d int) int {
		return a + b + c + d
	})
	assert.Equal(t, 10, add.Partial(1)(2, 3, 4))
	assert.Equal(t, 10, add.PartialR(1)(2, 3, 4))
}

func TestPartial5(t *testing.T) {
	add := Partial5(func(a, b, c, d, e int) int {
		return a + b + c + d + e
	})
	assert.Equal(t, 15, add.Partial(1)(2, 3, 4, 5))
	assert.Equal(t, 15, add.PartialR(1)(2, 3, 4, 5))
}

func TestPartial6(t *testing.T) {
	add := Partial6(func(a, b, c, d, e, f int) int {
		return a + b + c + d + e + f
	})
	assert.Equal(t, 21, add.Partial(1)(2, 3, 4, 5, 6))
	assert.Equal(t, 21, add.PartialR(1)(2, 3, 4, 5, 6))
}

func TestPartial7(t *testing.T) {
	add := Partial7(func(a, b, c, d, e, f, g int) int {
		return a + b + c + d + e + f + g
	})
	assert.Equal(t, 28, add.Partial(1)(2, 3, 4, 5, 6, 7))
	assert.Equal(t, 28, add.PartialR(1)(2, 3, 4, 5, 6, 7))
}

func TestPartial8(t *testing.T) {
	add := Partial8(func(a, b, c, d, e, f, g, h int) int {
		return a + b + c + d + e + f + g + h
	})
	assert.Equal(t, 36, add.Partial(1)(2, 3, 4, 5, 6, 7, 8))
	assert.Equal(t, 36, add.PartialR(1)(2, 3, 4, 5, 6, 7, 8))
}

func TestPartial9(t *testing.T) {
	add := Partial9(func(a, b, c, d, e, f, g, h, i int) int {
		return a + b + c + d + e + f + g + h + i
	})
	assert.Equal(t, 45, add.Partial(1)(2, 3, 4, 5, 6, 7, 8, 9))
	assert.Equal(t, 45, add.PartialR(1)(2, 3, 4, 5, 6, 7, 8, 9))
}

func TestPartial10(t *testing.T) {
	type myInt1 int
	type myInt2 int
	type myInt3 int
	type myInt4 int
	type myInt5 int
	type myInt6 int
	type myInt7 int
	type myInt8 int
	type myInt9 int
	type myInt10 int

	add := Partial10(func(a myInt1, b myInt2, c myInt3, d myInt4, e myInt5, f myInt6, g myInt7, h myInt8, i myInt9, j myInt10) int {
		return int(a) + int(b) + int(c) + int(d) + int(e) + int(f) + int(g) + int(h) + int(i) + int(j)
	})
	assert.Equal(t,
		55,
		add.
			Partial(1).
			Partial(2).
			Partial(3).
			Partial(4).
			Partial(5).
			Partial(6).
			Partial(7).
			Partial(8).
			Partial(9).
			Partial(10)())
	assert.Equal(t,
		55,
		add.
			PartialR(1).
			PartialR(2).
			PartialR(3).
			PartialR(4).
			PartialR(5).
			PartialR(6).
			PartialR(7).
			PartialR(8).
			PartialR(9).
			PartialR(10)())
}
