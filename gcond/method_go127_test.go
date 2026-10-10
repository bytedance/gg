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

package gcond

import (
	"strconv"
	"testing"

	"github.com/bytedance/gg/internal/assert"
)

func TestGo127LazyGet(t *testing.T) {
	l := Lazy[int](func() int { return 1 })
	assert.Equal(t, 1, l.Get())
}

func TestGo127LazyMap(t *testing.T) {
	l := Lazy[int](func() int { return 1 })
	assert.Equal(t, "1", l.Map(strconv.Itoa).Get())

	// The mapping is lazy: the function is not called until Get is called.
	notCalled := true
	lazy := l.Map(func(i int) string {
		notCalled = false
		return strconv.Itoa(i)
	})
	assert.True(t, notCalled)
	assert.Equal(t, "1", lazy.Get())
	assert.False(t, notCalled)
}

func TestGo127LazyThen(t *testing.T) {
	l := Lazy[int](func() int { return 1 })
	got := l.Then(func(i int) Lazy[string] {
		return func() string { return strconv.Itoa(i) }
	}).Get()
	assert.Equal(t, "1", got)
}

func TestGo127LazyOrElse(t *testing.T) {
	l := Lazy[int](func() int { return 1 })
	assert.Equal(t, 1, l.OrElse(true, func() int { return 2 }).Get())
	assert.Equal(t, 2, l.OrElse(false, func() int { return 2 }).Get())
}
