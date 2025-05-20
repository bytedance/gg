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

	"github.com/bytedance/gg/internal/assert"
)

func TestPeekerPeek(t *testing.T) {
	p := ToPeeker(FromSlice([]int{}))
	for i := 0; i < 10; i++ {
		assert.Zero(t, len(p.Peek(i)))
	}

	s := []int{1, 2, 3, 4}
	p = ToPeeker(FromSlice(s))
	for i := 0; i < 10; i++ {
		if i < len(s) {
			assert.NotZero(t, len(p.Peek(1)))
			assert.NotZero(t, len(p.Next(1)))
		} else {
			assert.Zero(t, len(p.Peek(1)))
			assert.Zero(t, len(p.Next(1)))
		}
	}
}
