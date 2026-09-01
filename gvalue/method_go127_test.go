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

package gvalue

import (
	"strconv"
	"testing"

	"github.com/bytedance/gg/internal/assert"
)

func TestVCast(t *testing.T) {
	assert.Equal(t, 1, Of(any(1)).Cast[int]())

	r, ok := Of(any(1)).TryCast[string]()
	assert.Equal(t, "", r)
	assert.False(t, ok)

	r2, ok2 := Of(any(1)).TryCast[int]()
	assert.Equal(t, 1, r2)
	assert.True(t, ok2)
}

func TestVMap(t *testing.T) {
	assert.Equal(t, "1", Of(1).Map(strconv.Itoa).Get())
	assert.Equal(t, "1", Of(1).Then(func(i int) V[string] { return Of(strconv.Itoa(i)) }).Get())

	// The key benefit of generic methods: the value type can change in a chain.
	assert.Equal(t, "4", Of(1).Map(func(i int) int { return i + 1 }).Map(func(i int) int { return i * 2 }).Map(strconv.Itoa).Get())
}

func TestVZero(t *testing.T) {
	assert.True(t, Of(0).IsZero())
	assert.True(t, Of("").IsZero())
	assert.True(t, Of([]int(nil)).IsZero())
	assert.False(t, Of([]int{}).IsZero())
	assert.True(t, Of(1).IsNotZero())
}

func TestVNil(t *testing.T) {
	assert.True(t, Of([]int(nil)).IsNil())
	assert.False(t, Of(1).IsNil())
	assert.True(t, Of(1).IsNotNil())
}

func TestVEqual(t *testing.T) {
	assert.True(t, Of(1).Equal(1))
	assert.False(t, Of(1).Equal(2))
	// Non-comparable values fall back to reflect.DeepEqual.
	assert.True(t, Of([]int{1, 2}).Equal([]int{1, 2}))
	assert.False(t, Of([]int{1, 2}).Equal([]int{1, 3}))
}

func TestVString(t *testing.T) {
	assert.Equal(t, "1", Of(1).String())
	// A nil slice is formatted as "[]" by the fmt package.
	assert.Equal(t, "[]", Of([]int(nil)).String())
}

func TestOrd(t *testing.T) {
	assert.Equal(t, 3, OrderedOf(1).Max(2, 3))
	assert.Equal(t, 1, OrderedOf(2).Min(1, 3))

	min, max := OrderedOf(2).MinMax(1, 3)
	assert.Equal(t, 1, min)
	assert.Equal(t, 3, max)

	assert.Equal(t, 3, OrderedOf(5).Clamp(1, 3))
	assert.Equal(t, 1, OrderedOf(0).Clamp(1, 3))
	assert.True(t, OrderedOf(2).Between(1, 3))
	assert.False(t, OrderedOf(4).Between(1, 3))

	assert.True(t, OrderedOf(1).Less(2))
	assert.True(t, OrderedOf(2).LessEqual(2))
	assert.True(t, OrderedOf(3).Greater(2))
	assert.True(t, OrderedOf(2).GreaterEqual(2))

	assert.Equal(t, 2, OrderedOf(2).Get())
	assert.Equal(t, "2", OrderedOf(2).String())
}

func TestNum(t *testing.T) {
	assert.Equal(t, 3, NumericOf(1).Add(2))
	assert.Equal(t, "ab", NumericOf("a").Add("b"))
	assert.Equal(t, 1, NumericOf(1).Get())
	assert.Equal(t, "1", NumericOf(1).String())
}
