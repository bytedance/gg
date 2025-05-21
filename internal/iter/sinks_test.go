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
	"context"
	"testing"

	"github.com/bytedance/gg/internal/assert"
)

func TestToChan(t *testing.T) {
	assertSliceEqual(t,
		[]int{1, 2, 3},
		func() Iter[int] {
			return FromChan(context.Background(),
				ToChan(context.Background(),
					FromSlice([]int{1, 2, 3})))
		})

	assert.NotPanic(t, func() {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel it immediately.

		_ = ToSlice(
			FromChan(context.Background(),
				ToChan(ctx,
					Iter[int](Range(1, 100000)))))
	})
}

func TestToBufferedChan(t *testing.T) {
	assertSliceEqual(t,
		[]int{1, 2, 3},
		func() Iter[int] {
			return FromChan(context.Background(),
				ToBufferedChan(context.Background(), 10,
					FromSlice([]int{1, 2, 3})))
		})

	assert.NotPanic(t, func() {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel it immediately.

		_ = ToSlice(
			FromChan(context.Background(),
				ToBufferedChan(ctx, 100,
					Iter[int](Range(1, 100000)))))
	})
}
