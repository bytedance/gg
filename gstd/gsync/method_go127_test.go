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

package gsync

import (
	"sort"
	"strconv"
	"testing"

	"github.com/bytedance/gg/internal/assert"
)

func TestGo127MapExtras(t *testing.T) {
	sm := Map[string, int]{}
	sm.Store("a", 1)
	sm.Store("b", 2)

	assert.Equal(t, 2, sm.Len())

	keys := sm.Keys()
	sort.Strings(keys)
	assert.Equal(t, []string{"a", "b"}, keys)

	values := sm.Values()
	sort.Ints(values)
	assert.Equal(t, []int{1, 2}, values)

	items := sm.ToSlice(func(k string, v int) string { return k + strconv.Itoa(v) })
	sort.Strings(items)
	assert.Equal(t, []string{"a1", "b2"}, items)
}

func TestGo127PoolWith(t *testing.T) {
	pool := Pool[*int]{New: func() *int {
		i := 1
		return &i
	}}

	called := false
	pool.With(func(i *int) {
		called = true
		*i = 2
	})
	assert.True(t, called)

	// The value is put back to the pool.
	pool.With(func(i *int) {
		assert.NotNil(t, i)
	})
}
