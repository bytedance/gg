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

package gptr

import (
	"strconv"
	"testing"

	"github.com/bytedance/gg/internal/assert"
)

func TestPNew(t *testing.T) {
	i := 1
	assert.Equal(t, 1, New(1).Indirect())
	assert.Equal(t, 1, Wrap(&i).Indirect())
	assert.Equal(t, 1, *Wrap(&i).Get())
	assert.Equal(t, 100, Wrap[int](nil).IndirectOr(100))
	assert.Equal(t, 100, Wrap[int](nil).IndirectOrLazy(func() int { return 100 }))
	assert.Equal(t, 1, New(1).IndirectOrLazy(func() int { return 100 }))
}

func TestPIsNil(t *testing.T) {
	assert.True(t, Wrap[int](nil).IsNil())
	assert.True(t, New(1).IsNotNil())
	assert.True(t, Wrap[int](nil).IsNilOrZero())
	assert.True(t, New(0).IsNilOrZero())
	assert.False(t, New(1).IsNilOrZero())
}

func TestPMap(t *testing.T) {
	i := 1
	assert.Equal(t, "1", Wrap(&i).Map(strconv.Itoa).Indirect())
	assert.True(t, Wrap[int](nil).Map(strconv.Itoa).IsNil())

	// The key benefit of generic methods: the pointee type can change in a chain.
	assert.Equal(t, 3, New(1).Map(func(i int) int { return i + 2 }).Indirect())
}

func TestPClone(t *testing.T) {
	assert.Equal(t, 1, New(1).Clone().Indirect())
	assert.Equal(t, 2, New(1).CloneBy(func(i int) int { return i * 2 }).Indirect())
	assert.True(t, Wrap[int](nil).Clone().IsNil())
}

func TestPEqual(t *testing.T) {
	assert.True(t, New(1).EqualTo(1))
	assert.False(t, New(1).EqualTo(2))
	assert.False(t, Wrap[int](nil).EqualTo(1))

	assert.True(t, New(1).Equal(New(1)))
	assert.True(t, New(1).Equal(New(1))) // same value, different address
	assert.False(t, New(1).Equal(New(2)))
	assert.True(t, Wrap[int](nil).Equal(Wrap[int](nil)))
	assert.False(t, New(1).Equal(Wrap[int](nil)))
}

func TestPEqualSameAddress(t *testing.T) {
	i := 1
	p := Wrap(&i)
	// The same address short-circuits the value comparison.
	assert.True(t, p.Equal(p))
	assert.True(t, p.Equal(Wrap(&i)))

	// One of them is nil.
	assert.False(t, p.Equal(Wrap[int](nil)))
	assert.False(t, Wrap[int](nil).Equal(p))
}

func TestPZeroValue(t *testing.T) {
	assert.True(t, New(0).HasZeroValue())
	assert.True(t, New(1).HasNonZeroValue())
	assert.False(t, Wrap[int](nil).HasZeroValue())
	assert.False(t, Wrap[int](nil).HasNonZeroValue())
}

func TestPString(t *testing.T) {
	assert.Equal(t, "1", New(1).String())
	assert.Equal(t, "<nil>", Wrap[int](nil).String())
}
