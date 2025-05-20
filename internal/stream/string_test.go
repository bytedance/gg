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
	"strconv"
	"strings"
	"testing"

	"github.com/bytedance/gg/internal/assert"
	"github.com/bytedance/gg/iter"
)

func BenchmarkStringJoin(b *testing.B) {
	n := 10000
	var s []int
	for i := 0; i < n; i++ {
		s = append(s, rand.Int())
	}
	b.ResetTimer()

	b.Run("Baseline", func(b *testing.B) {
		for i := 0; i <= b.N; i++ {
			var strs []string
			for _, v := range s {
				strs = append(strs, strconv.Itoa(v))
			}
			strings.Join(strs, ", ")
		}
	})
	b.Run("Stream", func(b *testing.B) {
		for i := 0; i <= b.N; i++ {
			FromStringIter(
				iter.Map(strconv.Itoa, iter.StealSlice(s)),
			).Join(", ")
		}
	})
}

func TestString_Join(t *testing.T) {
	assert.Equal(t, "1,2,3", FromStringSlice([]string{"1", "2", "3"}).Join(","))
}
