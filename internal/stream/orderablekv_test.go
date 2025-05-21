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

package stream

import (
	"math/rand"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/bytedance/gg/internal/assert"
)

func TestOrderableKVSort(t *testing.T) {
	assert.Equal(t,
		[]string{"Alice", "Bob", "Zhang"},
		FromOrderableMap(map[string]int{"Alice": 99, "Bob": 100, "Zhang": 59}).
			Sort().
			Keys().
			ToSlice())
	assert.Equal(t,
		[]int{99, 100, 59},
		FromOrderableMap(map[string]int{"Alice": 99, "Bob": 100, "Zhang": 59}).
			Sort().
			Values().
			ToSlice())
}

func BenchmarkOrderedRangeString(b *testing.B) {
	rnd := rand.New(rand.NewSource(time.Now().Unix()))
	randString := func(n int) string {
		b := make([]byte, n)
		for i := range b {
			b[i] = byte(rnd.Intn(256))
		}
		return string(b)
	}

	n := 1000
	m := make(map[string]int)
	for i := 0; i < n; i++ {
		m[randString(32)] = i
	}
	b.ResetTimer()

	b.Run("Baseline", func(b *testing.B) {
		for i := 0; i <= b.N; i++ {
			var keys []string
			for k := range m {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				v := m[k]
				_ = v
			}
		}
	})

	b.Run("SortKey", func(b *testing.B) {
		for i := 0; i <= b.N; i++ {
			for _, k := range FromOrderableMapKeys(m).Sort().ToSlice() {
				v := m[k]
				_ = v
			}
		}
	})

	b.Run("SortKV", func(b *testing.B) {
		for i := 0; i <= b.N; i++ {
			for _, t := range FromOrderableMap(m).Sort().ToSlice() {
				_, _ = t.First, t.Second
			}
		}
	})
}

func BenchmarkOrderedRangeInt(b *testing.B) {
	n := 1000
	m := make(map[int]string)
	for i := 0; i < n; i++ {
		m[i] = strconv.Itoa(i)
	}
	b.ResetTimer()

	b.Run("Baseline", func(b *testing.B) {
		for i := 0; i <= b.N; i++ {
			var keys []int
			for k := range m {
				keys = append(keys, k)
			}
			sort.Ints(keys)
			for _, k := range keys {
				v := m[k]
				_ = v
			}
		}
	})

	b.Run("SortKey", func(b *testing.B) {
		for i := 0; i <= b.N; i++ {
			for _, k := range FromOrderableMapKeys(m).Sort().ToSlice() {
				v := m[k]
				_ = v
			}
		}
	})

	b.Run("SortKV", func(b *testing.B) {
		for i := 0; i <= b.N; i++ {
			for _, t := range FromOrderableMap(m).Sort().ToSlice() {
				_, _ = t.First, t.Second
			}
		}
	})
}

func TestOrderableKV_Keys(t *testing.T) {
	assert.Equal(t, []int{1, 2, 3},
		FromOrderableMap(map[int]int{1: 2, 2: 4, 3: 6}).
			Sort().
			Keys().
			ToSlice())
}
