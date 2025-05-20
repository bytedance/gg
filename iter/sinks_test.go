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
