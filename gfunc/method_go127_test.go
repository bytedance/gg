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

package gfunc

import (
	"strconv"
	"testing"

	"github.com/bytedance/gg/internal/assert"
)

func TestGo127Compose(t *testing.T) {
	add := Partial2(func(a, b int) int { return a + b })

	// Partial application is already a method, Compose continues the chain and
	// changes the result type.
	assert.Equal(t, "3", add.Partial(1).Compose(strconv.Itoa)(2))
	assert.Equal(t, "3", add.PartialR(2).Compose(strconv.Itoa)(1))
}

func TestGo127ComposeFunc1(t *testing.T) {
	double := Partial1(func(a int) int { return a * 2 })
	assert.Equal(t, "4", double.Compose(strconv.Itoa)(2))
}

func TestGo127ComposeFunc3(t *testing.T) {
	add3 := Partial3(func(a, b, c int) int { return a + b + c })
	assert.Equal(t, "6", add3.Partial(1).Partial(2).Compose(strconv.Itoa)(3))
	assert.Equal(t, "6", add3.PartialR(3).PartialR(2).Compose(strconv.Itoa)(1))
}

func TestGo127ComposeFunc4(t *testing.T) {
	add4 := Partial4(func(a, b, c, d int) int { return a + b + c + d })
	assert.Equal(t, "10", add4.Partial(1).Partial(2).Partial(3).Compose(strconv.Itoa)(4))
}
