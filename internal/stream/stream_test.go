// The TestXXX functions only test the correctness of wrapper, please refer to
// package iter for detailed tests.
package stream

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/bytedance/gg/gcond"
	"github.com/bytedance/gg/gfunc"
	"github.com/bytedance/gg/gvalue"
	"github.com/bytedance/gg/internal/assert"
	"github.com/bytedance/gg/iter"
)

// Filter zero value.
func ExampleStream() {
	s := FromSlice([]int{0, 1, 2, 3, 4}).
		Filter(gvalue.IsNotZero[int]).
		ToSlice()

	fmt.Println(s)
	// Output:
	// [1 2 3 4]
}

func ExampleFromSlice() {
	add := gvalue.Add[int]                 // instantiate a int version of Add function
	add1 := gfunc.Partial2(add).Partial(1) // bind the first argument to 1
	s := FromSlice([]int{1, 2, 3, 4}).
		Map(add1).
		ToSlice()

	fmt.Println(s)
	// Output:
	// [2 3 4 5]
}

func ExampleFromIter() {
	s := FromIter(iter.Map(strconv.Itoa, iter.Range(1, 6))).
		Intersperse(", ").
		Prepend("[").
		Append("]").
		Fold(gvalue.Add[string], "")

	fmt.Println(s)
	// Output:
	// [1, 2, 3, 4, 5]
}

func BenchmarkStreamFilterP60N100(b *testing.B) {
	benchmarkStreamFilterPN(b, 60, 100)
}

func BenchmarkStreamFilterP60N10000(b *testing.B) {
	benchmarkStreamFilterPN(b, 60, 10000)
}

func BenchmarkStreamFilterP90N100(b *testing.B) {
	benchmarkStreamFilterPN(b, 90, 100)
}

func BenchmarkStreamFilterP90N10000(b *testing.B) {
	benchmarkStreamFilterPN(b, 90, 10000)
}

func BenchmarkStreamFilterP99N100(b *testing.B) {
	benchmarkStreamFilterPN(b, 99, 100)
}

func BenchmarkStreamFilterP99N10000(b *testing.B) {
	benchmarkStreamFilterPN(b, 99, 10000)
}

// P: Filter percentage
// N: Element count
func benchmarkStreamFilterPN(b *testing.B, p, n int) {
	rand.Seed(time.Now().UnixNano())
	var s []int
	for i := 0; i < n; i++ {
		s = append(s, rand.Intn(100))
	}
	b.ResetTimer()

	b.Run("Baseline", func(b *testing.B) {
		for i := 0; i <= b.N; i++ {
			var res []int
			for _, v := range s {
				if v > 100-p {
					res = append(res, v)
				}
			}
			res = nil
		}
	})

	pred := func(v int) bool { return v > 100-p }
	b.Run("Stream", func(b *testing.B) {
		for i := 0; i <= b.N; i++ {
			res := FromSlice(s).
				Filter(pred).
				ToSlice()
			_ = res
		}
	})

	b.Run("StreamSteal", func(b *testing.B) {
		b.StopTimer()
		ss := make([][]int, b.N+1)
		for i := range ss {
			ss[i] = make([]int, len(s))
			copy(ss[i], s)
		}
		b.StartTimer()

		for i := 0; i <= b.N; i++ {
			res := StealSlice(ss[i]).
				Filter(pred).
				ToSlice()
			_ = res
		}
	})
}

func TestStream_Map(t *testing.T) {
	assert.Equal(t,
		[]int{1, 2, 3},
		FromSlice([]int{0, 1, 2}).
			Map(func(v int) int { return v + 1 }).
			ToSlice())
}

func TestStream_MapToAny(t *testing.T) {
	assert.Equal(t,
		[]any{1, 2, 3},
		FromSlice([]int{0, 1, 2}).
			MapToAny(func(v int) any { return v + 1 }).
			ToSlice())
}

func TestStream_FlatMap(t *testing.T) {
	assert.Equal(t,
		[]int{1, 1, 2, 2, 3, 3},
		FromSlice([]int{1, 2, 3}).
			FlatMap(func(v int) []int { return []int{v, v} }).
			ToSlice())
}

func TestStream_FlatMapToAny(t *testing.T) {
	assert.Equal(t,
		[]any{1, 1, 2, 2, 3, 3},
		FromSlice([]int{1, 2, 3}).
			FlatMapToAny(func(v int) []any { return []any{v, v} }).
			ToSlice())
}

func TestStream_Fold(t *testing.T) {
	assert.Equal(t, 6, FromSlice([]int{1, 2, 3}).Fold(gvalue.Add[int], 0))
}

func TestStream_FoldWithToAny(t *testing.T) {
	assert.Equal(t,
		any(6),
		FromSlice([]int{1, 2, 3}).
			FoldToAnyWith(func(a any, v int) any { return a.(int) + v }, 0))
}

func TestStream_Reduce(t *testing.T) {
	assert.Equal(t, 6, FromSlice([]int{1, 2, 3}).Reduce(gvalue.Add[int]).Value())
}

func TestStream_Filter(t *testing.T) {
	assert.Equal(t,
		[]int{1, 2, 3},
		FromSlice([]int{0, 1, 2, 3}).
			Filter(gvalue.IsNotZero[int]).
			ToSlice())
}

func TestStream_ForEach(t *testing.T) {
	sum := 0
	FromSlice([]int{1, 2, 3}).ForEach(func(v int) { sum += v })
	assert.Equal(t, 6, sum)
}

func TestStream_ForEachIndexed(t *testing.T) {
	sum := 0
	FromSlice([]int{1, 2, 3}).ForEachIndexed(func(i, v int) { sum += i })
	assert.Equal(t, 3, sum)
}

func TestStream_Head(t *testing.T) {
	assert.Equal(t, 1, FromSlice([]int{1, 2, 3}).Head().Value())
}

func TestStream_Reverse(t *testing.T) {
	assert.Equal(t, []int{3, 2, 1}, FromSlice([]int{1, 2, 3}).Reverse().ToSlice())
}

func TestStream_Take(t *testing.T) {
	assert.Equal(t, []int{1, 2}, FromSlice([]int{1, 2, 3}).Take(2).ToSlice())
}

func TestStream_Drop(t *testing.T) {
	assert.Equal(t, []int{2, 3}, FromSlice([]int{1, 2, 3}).Drop(1).ToSlice())
}

func TestStream_All(t *testing.T) {
	assert.True(t, FromSlice([]int{1, 2, 3}).All(gvalue.IsNotZero[int]))
}

func TestStream_Any(t *testing.T) {
	assert.True(t, FromSlice([]int{1, 2, 3, 0}).Any(gvalue.IsNotZero[int]))
}

func TestStream_Concat(t *testing.T) {
	assert.Equal(t, []int{1, 2, 3, 4}, FromSlice([]int{1, 2}).Concat(FromSlice([]int{3, 4})).ToSlice())
}

func TestStream_Zip(t *testing.T) {
	assert.Equal(t,
		[]int{2, 4, 6},
		FromSlice([]int{1, 2, 3}).
			Zip(gvalue.Add[int], FromSlice([]int{1, 2, 3})).
			ToSlice())
}

func TestStream_Intersperse(t *testing.T) {
	assert.Equal(t,
		[]int{1, 0, 2, 0, 3},
		FromSlice([]int{1, 2, 3}).
			Intersperse(0).
			ToSlice())
}

func TestStream_Append(t *testing.T) {
	assert.Equal(t,
		[]int{1, 2, 3, 4},
		FromSlice([]int{1, 2, 3}).Append(4).ToSlice())
}

func TestStream_Prepend(t *testing.T) {
	assert.Equal(t,
		[]int{0, 1, 2, 3},
		FromSlice([]int{1, 2, 3}).Prepend(0).ToSlice())
}

func TestStream_Find(t *testing.T) {
	assert.Equal(t, 1, FromSlice([]int{1, 2, 3}).Find(gvalue.IsNotZero[int]).Value())
}

func TestStream_Count(t *testing.T) {
	assert.Equal(t, 3, FromSlice([]int{1, 2, 3}).Count())
}

func TestStream_TakeWhile(t *testing.T) {
	assert.Equal(t,
		[]int{1, 2},
		FromSlice([]int{1, 2, 0, 3}).
			TakeWhile(gvalue.IsNotZero[int]).
			ToSlice())
}

func TestStream_DropWhile(t *testing.T) {
	assert.Equal(t,
		[]int{0, 3},
		FromSlice([]int{1, 2, 0, 3}).
			DropWhile(gvalue.IsNotZero[int]).
			ToSlice())
}

func TestStream_SortBy(t *testing.T) {
	assert.Equal(t,
		[]int{0, 1, 2, 3},
		FromSlice([]int{1, 2, 0, 3}).
			SortBy(gvalue.Less[int]).
			ToSlice())
}

func TestStream_At(t *testing.T) {
	assert.Equal(t, 3, FromSlice([]int{1, 2, 3, 4}).At(2).Value())
}

func TestStream_UniqBy(t *testing.T) {
	assert.Equal(t,
		[]int{1, 4, 2, 3},
		FromSlice([]int{1, 4, 2, 3, 4}).
			UniqBy(func(v int) any { return v }).
			ToSlice())
}

func TestStream_Chunk(t *testing.T) {
	assert.Equal(t,
		[][]int{{1, 2}, {3, 3}, {4}},
		FromSlice([]int{1, 2, 3, 3, 4}).Chunk(2))
}

func TestStream_GroupBy(t *testing.T) {
	assert.Equal(t,
		map[any][]int{"odd": {1, 3}, "even": {2, 4}},
		FromSlice([]int{1, 2, 3, 4}).GroupBy(func(v int) any {
			return gcond.If(v%2 == 0, "even", "odd")
		}))
}
