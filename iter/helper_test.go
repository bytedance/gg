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
