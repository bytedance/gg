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
	"testing"

	"github.com/bytedance/gg/internal/assert"
)

func TestOnceValue(t *testing.T) {
	{
		i := 2
		f := func() int {
			i++
			return i
		}
		assert.Equal(t, i, 2)
		once := OnceValue(f)
		assert.Equal(t, i, 2)
		assert.Equal(t, 3, once())
		assert.Equal(t, 3, once())
		assert.Equal(t, 3, once())
		assert.Equal(t, i, 3)
	}

	{ // Test concurrency.
		i := 2
		f := func() int {
			i++
			return i
		}
		assert.Equal(t, i, 2)
		once := OnceValue(f)
		assert.Equal(t, i, 2)
		var wg sync.WaitGroup
		wg.Add(100)
		for j := 0; j < 100; j++ {
			go func() {
				assert.Equal(t, 3, once())
				assert.Equal(t, i, 3)
				wg.Done()
			}()
		}
		wg.Wait()
		assert.Equal(t, i, 3)
	}
}
