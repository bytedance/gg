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
