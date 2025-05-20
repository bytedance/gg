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

package iter

import (
	"testing"

	"github.com/bytedance/gg/internal/fastrand"
)

func BenchmarkUniq(b *testing.B) {
	const M = 1000
	nums := make([]int, M)
	verify := make(map[int]struct{})
	for i := 0; i < M; i++ {
		nums[i] = fastrand.Intn(100)
		verify[nums[i]] = struct{}{}
	}

	b.Run("fast", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			l := len(ToSlice(Uniq(FromSlice(nums))))
			if l != len(verify) {
				b.Errorf("mismatched len, expect %d, found %d", len(verify), l)
			}
		}
	})

	b.Run("slow", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			l := len(ToSlice(Take(M, Uniq(FromSlice(nums)))))
			if l != len(verify) {
				b.Errorf("mismatched len, expect %d, found %d", len(verify), l)
			}
		}
	})
}
