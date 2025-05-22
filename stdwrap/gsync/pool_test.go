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

package gsync

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/bytedance/gg/gslice"
	"github.com/bytedance/gg/internal/assert"
)

func TestPool(t *testing.T) {
	var numAlloc int64

	p := Pool[*int]{
		New: func() *int {
			atomic.AddInt64(&numAlloc, 1)
			var i int
			return &i
		},
	}

	var (
		n    = 10
		wg   sync.WaitGroup
		vals = make([]*int, n)
	)

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			vals[i] = p.Get()
		}(i)
	}
	wg.Wait()

	assert.Equal(t, n, int(numAlloc))
	assert.Equal(t, n, len(vals))
	assert.Equal(t, n, len(gslice.Uniq(vals)))

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			p.Put(vals[i])
			vals[i] = nil
		}(i)
	}
	wg.Wait()
	assert.Equal(t, n, int(numAlloc))
}
