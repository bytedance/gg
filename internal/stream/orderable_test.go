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
	"testing"

	"github.com/bytedance/gg/internal/assert"
)

func BenchmarkOrderableMax(b *testing.B) {
	n := 10000
	var s []int
	for i := 0; i < n; i++ {
		s = append(s, rand.Int())
	}
	b.ResetTimer()

	b.Run("Baseline", func(b *testing.B) {
		for i := 0; i <= b.N; i++ {
			var res int
			if len(s) != 0 {
				res = s[0]
			}
			for _, v := range s {
				if res < v {
					res = v
				}
			}
			_ = res
		}
	})
	b.Run("Stream", func(b *testing.B) {
		var res int
		for i := 0; i <= b.N; i++ {
			res = StealOrderableSlice(s).Max().Value()
		}
		_ = res
	})
}

func TestOrderable_Max(t *testing.T) {
	assert.Equal(t, 99, FromOrderableSlice([]int{1, 3, 4, 99}).Max().Value())
}

func TestOrderable_Min(t *testing.T) {
	assert.Equal(t, 1, FromOrderableSlice([]int{1, 3, 4, 99}).Min().Value())
}

func TestOrderable_Sort(t *testing.T) {
	assert.Equal(t, []int{1, 2, 3, 4}, FromOrderableSlice([]int{4, 1, 3, 2}).Sort().ToSlice())
}
