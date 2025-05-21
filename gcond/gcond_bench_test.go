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
)

func BenchmarkIf(b *testing.B) {
	cond := true

	b.Run("Baseline", func(b *testing.B) {
		var v int
		for i := 0; i < b.N; i++ {
			if cond {
				v = 1
			} else {
				v = 2
			}
		}
		if v != 1 {
			b.FailNow()
		}
	})

	b.Run("If", func(b *testing.B) {
		var v int
		for i := 0; i < b.N; i++ {
			v = If(cond, 1, 2)
		}
		if v != 1 {
			b.FailNow()
		}
	})

	b.Run("IfLazy", func(b *testing.B) {
		var v int
		for i := 0; i < b.N; i++ {
			v = IfLazy(cond, func() int { return 1 }, func() int { return 2 })
		}
		if v != 1 {
			b.FailNow()
		}
	})
}
