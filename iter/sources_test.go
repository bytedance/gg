package iter

import (
	"context"
	"sort"
	"testing"

	"github.com/bytedance/gg/internal/assert"
)

func TestFromSlice(t *testing.T) {
	// Nil slice.
	assert.Equal(t,
		[]int{},
		ToSlice(FromSlice[int](nil)))
	// Empty slice.
	assert.Equal(t,
		[]int{},
		ToSlice(FromSlice([]int{})))
	assert.Equal(t,
		[]int{1, 2, 3, 4},
		ToSlice(FromSlice([]int{1, 2, 3, 4})))
	assert.Equal(t,
		[]int{},
		ToSlice(FromSlice([]int{})))
	assert.Equal(t,
		[]int{},
		ToSlice(FromSlice[int](nil)))
}

func TestFromMap(t *testing.T) {
	// Nil map.
	assert.Equal(t,
		map[int]string{},
		KVToMap(FromMap[int, string](nil)))
	// Empty map.
	assert.Equal(t,
		map[int]string{},
		KVToMap(FromMap(map[int]string{})))
	assert.Equal(t,
		map[int]string{1: "1", 2: "2", 3: "3"},
		KVToMap(FromMap(map[int]string{1: "1", 2: "2", 3: "3"})))
	assert.Equal(t,
		map[int]string{},
		KVToMap(FromMap(map[int]string{})))
	assert.Equal(t,
		map[int]string{},
		KVToMap(FromMap[int, string](nil)))
	assert.Equal(t,
		0,
		Count(
			Take(0,
				FromMap(map[int]string{1: "1", 2: "2", 3: "3"}))))
	assert.Equal(t,
		2,
		Count(
			Take(2,
				FromMap(map[int]string{1: "1", 2: "2", 3: "3"}))))
	assert.Equal(t,
		map[int]string{1: "1", 2: "2", 3: "3"},
		KVToMap(
			Take(100,
				FromMap(map[int]string{1: "1", 2: "2", 3: "3"}))))
}

func TestFromMapKeys_PartialRead(t *testing.T) {
	m := map[int]string{1: "1", 2: "2", 3: "3"}
	i := FromMapKeys(m)
	var s []int
	s = append(s, ToSlice(Take(1, i))...)
	s = append(s, ToSlice(Take(1, i))...)
	s = append(s, ToSlice(Take(1, i))...)
	sort.Ints(s)
	assert.Equal(t, []int{1, 2, 3}, s)
	assert.Equal(t, []int{}, ToSlice(Take(1, i)))
}

func TestFromMapKeys(t *testing.T) {
	assert.Equal(t,
		[]int{1, 2, 3},
		ToSlice(
			Sort(
				FromMapKeys(map[int]string{1: "1", 2: "2", 3: "3"}))))
	assert.Equal(t,
		[]int{},
		ToSlice(FromMapKeys(map[int]string{})))
	assert.Equal(t,
		[]int{},
		ToSlice(FromMapKeys[int, string](nil)))

	assert.Equal(t,
		0,
		Count(
			Take(0,
				FromMapKeys(map[int]string{1: "1", 2: "2", 3: "3"}))))
	assert.Equal(t,
		2,
		Count(
			Take(2,
				FromMapKeys(map[int]string{1: "1", 2: "2", 3: "3"}))))
	assert.Equal(t,
		[]int{1, 2, 3},
		ToSlice(
			Sort(
				Take(100,
					FromMapKeys(map[int]string{1: "1", 2: "2", 3: "3"})))))
}

func TestFromMapValues(t *testing.T) {
	assert.Equal(t,
		[]string{"1", "2", "3"},
		ToSlice(
			Sort(
				FromMapValues(map[int]string{1: "1", 2: "2", 3: "3"}))))
	assert.Equal(t,
		[]string{},
		ToSlice(FromMapValues(map[int]string{})))
	assert.Equal(t,
		[]string{},
		ToSlice(FromMapValues[int, string](nil)))

	assert.Equal(t,
		0,
		Count(
			Take(0,
				FromMapValues(map[int]string{1: "1", 2: "2", 3: "3"}))))
	assert.Equal(t,
		2,
		Count(
			Take(2,
				FromMapValues(map[int]string{1: "1", 2: "2", 3: "3"}))))
	assert.Equal(t,
		[]string{"1", "2", "3"},
		ToSlice(
			Sort(
				Take(100,
					FromMapValues(map[int]string{1: "1", 2: "2", 3: "3"})))))
}

func TestFromChan(t *testing.T) {
	ch1 := make(chan int)
	s1 := []int{1, 2, 3}
	go func() {
		for _, v := range s1 {
			ch1 <- v
		}
		close(ch1)
	}()
	assert.Equal(t, s1, ToSlice(FromChan(context.Background(), ch1)))

	// From closed channel
	ch2 := make(chan int)
	close(ch2)
	assert.Equal(t, []int{}, ToSlice(FromChan(context.Background(), ch2)))

	// Context
	ch3 := make(chan int)
	ctx3, cancel3 := context.WithCancel(context.Background())
	s3 := []int{1, 2, 3}
	go func() {
		for _, v := range s3 {
			ch3 <- v
		}
		cancel3()
	}()
	assert.Equal(t, s3, ToSlice(FromChan(ctx3, ch3)))

	// Take part of channel
	ch4 := make(chan int)
	s4 := []int{1, 2, 3}
	go func() {
		for _, v := range s4 {
			ch4 <- v
		}
		close(ch4)
	}()
	assert.Equal(t,
		s4[:2],
		ToSlice(
			Take(2,
				FromChan(context.Background(), ch4))))

	ch5 := make(chan int)
	s5 := []int{1, 2, 3}
	go func() {
		for _, v := range s5 {
			ch5 <- v
		}
		close(ch5)
	}()
	assert.Equal(t,
		s5,
		ToSlice(
			Take(100,
				FromChan(context.Background(), ch5))))

	// Take part of channel
	ch6 := make(chan int)
	s6 := []int{1, 2, 3}
	go func() {
		for _, v := range s6 {
			ch6 <- v
		}
		close(ch6)
	}()
	assert.Equal(t,
		[]int{},
		ToSlice(
			Take(0,
				FromChan(context.Background(), ch6))))
}

func TestRange(t *testing.T) {
	assertSliceEqual(t,
		[]int{},
		func() Iter[int] {
			return Range(3, 1)
		})
	assertSliceEqual(t,
		[]int{0, 1, 2, 3, 4},
		func() Iter[int] {
			return Range(0, 5)
		})
	assertSliceEqual(t,
		[]int{-2, -1, 0, 1, 2},
		func() Iter[int] {
			return Range(-2, 3)
		})
	assertSliceEqual(t,
		[]int{0},
		func() Iter[int] {
			return Range(0, 1)
		})
	assertSliceEqual(t,
		[]int{},
		func() Iter[int] {
			return Range(0, 0)
		})
	assertSliceEqual(t,
		[]int{},
		func() Iter[int] {
			return RangeWithStep(1, 3, -1)
		})
	assertSliceEqual(t,
		[]int{100, 89, 78, 67, 56},
		func() Iter[int] {
			return RangeWithStep(100, 50, -11)
		})
	assertSliceEqual(t,
		[]int{0, 3, 6, 9},
		func() Iter[int] {
			return RangeWithStep(0, 10, 3)
		})
	assertSliceEqual(t,
		[]float64{1.3, 1.6, 1.9, 2.2, 2.5, 2.8, 3.1, 3.4, 3.7},
		func() Iter[float64] {
			return RangeWithStep(1.3, 3.7, 0.3)
		})
	assertSliceEqual(t,
		[]float64{1.2, 0.7, 0.2, -0.3, -0.8},
		func() Iter[float64] {
			return RangeWithStep(1.2, -1.0, -0.5)
		})
}

func TestRepeat(t *testing.T) {
	assertSliceEqual(t,
		[]int{1, 1, 1, 1},
		func() Iter[int] {
			return Take(4, Repeat(1))
		})
	assert.Panic(t, func() {
		ToSlice(Repeat(1))
	})
}
