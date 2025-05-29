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

//go:build go1.20
// +build go1.20

package gsync

import (
	"testing"

	"github.com/bytedance/gg/goption"
	"github.com/bytedance/gg/internal/assert"
)

func TestSwap(t *testing.T) {
	sm := Map[string, int]{}
	sm.Store("k", 1)
	assert.Equal(t, goption.OK(1), goption.Of(sm.Swap("k", 2)))
	assert.Equal(t, goption.OK(2), sm.LoadO("k"))
	assert.Equal(t, goption.Nil[int](), goption.Of(sm.Swap("l", 3)))
	assert.Equal(t, goption.OK(3), sm.LoadO("l"))
}

func TestCompareAndSwap(t *testing.T) {
	sm := Map[string, int]{}
	sm.Store("k", 1)
	assert.False(t, sm.CompareAndSwap("k", 2, 3))
	assert.Equal(t, goption.OK(1), sm.LoadO("k"))
	assert.True(t, sm.CompareAndSwap("k", 1, 3))
	assert.Equal(t, goption.OK(3), sm.LoadO("k"))
}

func TestCompareAndDelete(t *testing.T) {
	sm := Map[string, int]{}
	sm.Store("k", 1)
	assert.False(t, sm.CompareAndDelete("k", 2))
	assert.Equal(t, goption.OK(1), sm.LoadO("k"))
	assert.True(t, sm.CompareAndDelete("k", 1))
	assert.False(t, sm.LoadO("k").IsOK())
}
