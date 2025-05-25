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

package gcond

import (
	"testing"

	"github.com/bytedance/gg/internal/assert"
)

func TestIf(t *testing.T) {
	assert.Equal(t, 1, If(true, 1, 2))
	assert.Equal(t, 2, If(false, 1, 2))
	assert.Equal(t, "2", If(false, "1", "2"))
	assert.Equal(t, "1", If(true, "1", "2"))

	assert.Panic(t, func() {
		var tt *testing.T
		_ = If(tt != nil, tt.Name(), "")
	})
	assert.Panic(t, func() {
		var tt *testing.T
		_ = If(tt == nil, "", tt.Name())
	})
}

func lazy[T any](v T) Lazy[T] {
	return func() T {
		return v
	}
}

func TestIfLazy(t *testing.T) {
	assert.Equal(t, 1, IfLazy(true, lazy(1), lazy(2)))
	assert.Equal(t, 2, IfLazy(false, lazy(1), lazy(2)))
	assert.Equal(t, "1", IfLazy(true, lazy("1"), lazy("2")))
	assert.Equal(t, "2", IfLazy(false, lazy("1"), lazy("2")))

	assert.NotPanic(t, func() {
		var tt *testing.T
		assert.Equal(t, "", IfLazy(tt != nil, func() string { return tt.Name() }, lazy("")))
		assert.Equal(t, "", IfLazy(tt == nil, lazy(""), func() string { return tt.Name() }))
	})
}

func TestIfLazyL(t *testing.T) {
	assert.Equal(t, 1, IfLazyL(true, lazy(1), 2))
	assert.Equal(t, 2, IfLazyL(false, lazy(1), 2))
	assert.Equal(t, "1", IfLazyL(true, lazy("1"), "2"))
	assert.Equal(t, "2", IfLazyL(false, lazy("1"), "2"))

	assert.NotPanic(t, func() {
		var tt *testing.T
		assert.Equal(t, "", IfLazyL(tt != nil, func() string { return tt.Name() }, ""))
	})
}

func TestIfLazyR(t *testing.T) {
	assert.Equal(t, 1, IfLazyR(true, 1, lazy(2)))
	assert.Equal(t, 2, IfLazyR(false, 1, lazy(2)))
	assert.Equal(t, "1", IfLazyR(true, "1", lazy("2")))
	assert.Equal(t, "2", IfLazyR(false, "1", lazy("2")))

	assert.NotPanic(t, func() {
		var tt *testing.T
		assert.Equal(t, "", IfLazyR(tt == nil, "", func() string { return tt.Name() }))
	})
}

func TestSwitch(t *testing.T) {
	v1 := Switch[string](1).
		Case(1, "1").
		Case(2, "2").
		CaseLazy(3, func() string { return "3" }).
		CaseLazy(4, func() string { return "4" }).
		Default("5")
	assert.Equal(t, v1, "1")

	v2 := Switch[string](3).
		Case(1, "1").
		Case(2, "2").
		CaseLazy(3, func() string { return "3" }).
		CaseLazy(4, func() string { return "4" }).
		Default("5")
	assert.Equal(t, v2, "3")

	v3 := Switch[string](10).
		Case(1, "1").
		Case(2, "2").
		CaseLazy(3, func() string { return "3" }).
		CaseLazy(4, func() string { return "4" }).
		Default("5")
	assert.Equal(t, v3, "5")
}

func TestSwitchWhen(t *testing.T) {
	v1 := Switch[string](1).
		When(1, 2).Then("1").
		When(3, 4).ThenLazy(func() string { return "3" }).
		DefaultLazy(func() string {
			return "5"
		})
	assert.Equal(t, v1, "1")

	v2 := Switch[string](4).
		When(1, 2).Then("1").
		When(3, 4).ThenLazy(func() string { return "3" }).
		DefaultLazy(func() string {
			return "5"
		})
	assert.Equal(t, v2, "3")

	v3 := Switch[string](10).
		When(1, 2).Then("1").
		When(3, 4).ThenLazy(func() string { return "3" }).
		DefaultLazy(func() string {
			return "5"
		})
	assert.Equal(t, v3, "5")
}
