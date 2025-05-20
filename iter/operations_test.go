package iter

import (
	"strconv"
	"strings"
	"testing"

	"github.com/bytedance/gg/collection/tuple"
	"github.com/bytedance/gg/gcond"
	"github.com/bytedance/gg/gfunc"
	"github.com/bytedance/gg/gvalue"
	"github.com/bytedance/gg/internal/assert"
)

func TestMap(t *testing.T) {
	assertSliceEqual(t,
		[]string{"1", "2", "3", "4"},
		func() Iter[string] {
			return Map(strconv.Itoa, FromSlice([]int{1, 2, 3, 4}))
		})
	assert.Equal(t,
		[]int{2, 3, 4, 5},
		ToSlice(
			Map(gfunc.Partial2(gvalue.Add[int]).Partial(1),
				FromSlice([]int{1, 2, 3, 4}))))
	// Test empty slice
	assert.Equal(t,
		[]string{},
		ToSlice(
			Map(strconv.Itoa,
				FromSlice([]int{}))))
}

func TestFilterMap(t *testing.T) {

	fn := func(i int) (string, bool) {
		return strconv.Itoa(i), i != 0
	}

	assertSliceEqual(t,
		[]string{"1", "2", "3", "4"},
		func() Iter[string] {
			return FilterMap(fn, FromSlice([]int{1, 2, 3, 4}))
		})

	assert.Equal(t,
		[]int{2, 3, 4, 5},
		ToSlice(
			FilterMap(func(i int) (int, bool) {
				return i + 1, i != 0
			}, FromSlice([]int{1, 2, 3, 4, 0, 0}))))
	// Test empty slice
	assert.Equal(t,
		[]string{},
		ToSlice(
			FilterMap(fn, FromSlice([]int{0, 0}))))
	assert.Equal(t,
		[]string{},
		ToSlice(
			FilterMap(fn, FromSlice([]int{}))))

	m := FilterMap(fn, FromSlice([]int{1, 2, 3, 0, 4}))
	assert.Equal(t, []string{"1", "2"}, m.Next(2))
	assert.Equal(t, []string{"3"}, m.Next(2))
	assert.Equal(t, []string{"4"}, m.Next(2))
	assert.Equal(t, nil, m.Next(2))
	assert.Equal(t, nil, m.Next(100))
	assert.Equal(t, nil, m.Next(ALL))

}

func TestMapInplace(t *testing.T) {
	assertSliceEqual(t,
		[]int{2, 3, 4, 5},
		func() Iter[int] {
			return MapInplace(gfunc.Partial2(gvalue.Add[int]).Partial(1),
				FromSlice([]int{1, 2, 3, 4}))
		})
	// Test empty slice
	assert.Equal(t,
		[]int{},
		ToSlice(
			MapInplace(gfunc.Partial2(gvalue.Add[int]).Partial(1),
				FromSlice([]int{}))))
	// Test in place.
	s1 := []int{1, 2, 3, 4}
	assert.Equal(t,
		[]int{2, 3, 4, 5},
		ToSlice(
			MapInplace(gfunc.Partial2(gvalue.Add[int]).Partial(1),
				StealSlice(s1))))
	assert.Equal(t, []int{2, 3, 4, 5}, s1)
}

func TestFlatMap(t *testing.T) {
	splitSpace := gfunc.Partial2(strings.Split).PartialR(" ")
	assert.Equal(t,
		[]string{},
		ToSlice(
			FlatMap(splitSpace,
				FromSlice([]string{}))))
	assertSliceEqual(t,
		[]string{"1", "2", "3", "4"},
		func() Iter[string] {
			return FlatMap(splitSpace,
				FromSlice([]string{"1 2", "3 4"}))
		})

	splitSpace2 := func(v string) []string {
		s := ToSlice(
			Filter(gvalue.IsNotZero[string],
				FromSlice(strings.Split(v, " "))))
		return gcond.If(len(s) == 0, nil, s)
	}
	assertSliceEqual(t,
		[]string{"1", "2", "3", "4"},
		func() Iter[string] {
			return FlatMap(splitSpace2,
				FromSlice([]string{"   ", "1", "2", "3 4"}))
		})
}

func TestFold(t *testing.T) {
	assert.Equal(t,
		"124",
		Fold(func(a string, b int) string { return a + strconv.Itoa(b) }, "",
			FromSlice([]int{1, 2, 4})))
}

func TestReduce(t *testing.T) {
	assert.Equal(t,
		7,
		Reduce(func(a, b int) int { return a + b },
			FromSlice([]int{1, 2, 4})).
			Value())
}

func TestFilter(t *testing.T) {
	assertSliceEqual(t,
		[]int{1, 2, 3, 4},
		func() Iter[int] {
			return Filter(gvalue.IsNotZero[int], FromSlice([]int{1, 2, 0, 3, 0, 4}))
		})

	// Test stateful predicate.
	skipNth := func(n int) func(int) bool {
		return func(int) bool {
			n--
			return n+1 != 0
		}
	}
	assertSliceEqual(t,
		[]int{2, 3, 4},
		func() Iter[int] {
			return Filter(skipNth(0), FromSlice([]int{1, 2, 3, 4}))
		})
}

func TestHead(t *testing.T) {
	assert.Equal(t,
		1,
		Head(FromSlice([]int{1, 2, 0, 3, 0, 4})).Value())
	assert.True(t,
		Head(FromSlice([]int{})).IsNil())
}
func TestReverse(t *testing.T) {
	// Even.
	assertSliceEqual(t,
		[]int{4, 3, 2, 1},
		func() Iter[int] {
			return Reverse(
				FromSlice([]int{1, 2, 3, 4}))
		})
	// Odd.
	assertSliceEqual(t,
		[]int{5, 4, 3, 2, 1},
		func() Iter[int] {
			return Reverse(
				FromSlice([]int{1, 2, 3, 4, 5}))
		})
	// Empty.
	assert.Equal(t,
		[]int{},
		ToSlice(
			Reverse(
				FromSlice([]int{}))))

	// Check internal state.
	s := []int{1, 2, 3, 4, 5, 6, 7, 8}
	assert.Equal(t,
		[]int{8, 7},
		ToSlice(
			Take(2,
				Reverse(
					StealSlice(s)))))
	assert.Equal(t, []int{8, 7, 3, 4, 5, 6, 2, 1}, s)

	// Check internal state.
	s = []int{1, 2, 3, 4, 5, 6, 7, 8}
	assert.Equal(t,
		[]int{8, 7, 6, 5, 4},
		ToSlice(
			Take(5,
				Reverse(
					StealSlice(s)))))
	assert.Equal(t, []int{8, 7, 6, 5, 4, 3, 2, 1}, s)
}

func TestMax(t *testing.T) {
	assert.False(t, Max(FromSlice([]int{})).IsOK())
	assert.Equal(t,
		-1,
		Max(FromSlice([]int{-1})).Value())
	assert.Equal(t,
		100,
		Max(FromSlice([]int{10, 1, -1, 100, 3})).Value())
}

func TestMaxBy(t *testing.T) {
	type Foo struct {
		Value int
	}
	// FIXME: use greflect
	less := func(foo1, foo2 Foo) bool {
		return foo1.Value < foo2.Value
	}
	assert.False(t, MaxBy(less, FromSlice([]Foo{})).IsOK())
	assert.Equal(t,
		Foo{-1},
		MaxBy(less, FromSlice([]Foo{{-1}})).Value())
	assert.Equal(t,
		Foo{100},
		MaxBy(less, FromSlice([]Foo{{10}, {1}, {-1}, {100}, {3}})).Value())
}

func TestMin(t *testing.T) {
	assert.False(t, Min(FromSlice([]int{})).IsOK())
	assert.Equal(t,
		100,
		Min(FromSlice([]int{100})).Value())
	assert.Equal(t,
		-1,
		Min(FromSlice([]int{10, 1, -1, 100, 3})).Value())
}

func TestMinBy(t *testing.T) {
	type Foo struct {
		Value int
	}
	// FIXME: use greflect
	less := func(foo1, foo2 Foo) bool {
		return foo1.Value < foo2.Value
	}
	assert.False(t, MinBy(less, FromSlice([]Foo{})).IsOK())
	assert.Equal(t,
		Foo{-1},
		MinBy(less, FromSlice([]Foo{{-1}})).Value())
	assert.Equal(t,
		Foo{-1},
		MinBy(less, FromSlice([]Foo{{10}, {1}, {-1}, {100}, {3}})).Value())
}

func TestMinMax(t *testing.T) {
	assert.False(t, MinMax(FromSlice([]int{})).IsOK())
	assert.Equal(t,
		tuple.Make2(100, 100),
		MinMax(FromSlice([]int{100})).Value())
	assert.Equal(t,
		tuple.Make2(-1, -1),
		MinMax(FromSlice([]int{-1})).Value())
	assert.Equal(t,
		tuple.Make2(-1, 100),
		MinMax(FromSlice([]int{10, 1, -1, 100, 3})).Value())
}

func TestMinMaxBy(t *testing.T) {
	type Foo struct {
		Value int
	}
	// FIXME: use greflect
	less := func(foo1, foo2 Foo) bool {
		return foo1.Value < foo2.Value
	}
	assert.False(t, MinMaxBy(less, FromSlice([]Foo{})).IsOK())
	assert.Equal(t,
		tuple.Make2(Foo{-1}, Foo{-1}),
		MinMaxBy(less, FromSlice([]Foo{{-1}})).Value())
	assert.Equal(t,
		tuple.Make2(Foo{-1}, Foo{100}),
		MinMaxBy(less, FromSlice([]Foo{{10}, {1}, {-1}, {100}, {3}})).Value())
}

func TestTake(t *testing.T) {
	assert.Equal(t,
		[]int{1, 2},
		ToSlice(
			Take(2,
				FromSlice([]int{1, 2, 3, 4}))))
	assert.Equal(t,
		[]int{1},
		ToSlice(
			Take(1,
				FromSlice([]int{1, 2, 3, 4}))))
	assert.Equal(t,
		[]int{},
		ToSlice(
			Take(0,
				FromSlice([]int{1, 2, 3, 4}))))
	assert.Equal(t,
		[]int{1, 2, 3, 4},
		ToSlice(
			Take(10,
				FromSlice([]int{1, 2, 3, 4}))))
	assert.Equal(t,
		[]int{1, 2, 3},
		ToSlice(
			Take(3,
				Take(10,
					FromSlice([]int{1, 2, 3, 4})))))
	assert.Panic(t, func() {
		ToSlice(
			Take(-1,
				FromSlice([]int{1, 2, 3, 4})))
	})
}

func TestDrop(t *testing.T) {
	assert.Equal(t,
		[]int{3, 4},
		ToSlice(
			Drop(2,
				FromSlice([]int{1, 2, 3, 4}))))
	assert.Equal(t,
		[]int{2, 3, 4},
		ToSlice(
			Drop(1,
				FromSlice([]int{1, 2, 3, 4}))))
	assert.Equal(t,
		[]int{1, 2, 3, 4},
		ToSlice(
			Drop(0,
				FromSlice([]int{1, 2, 3, 4}))))
	assert.Equal(t,
		[]int{},
		ToSlice(
			Drop(10,
				FromSlice([]int{1, 2, 3, 4}))))
	assert.Panic(t, func() {
		ToSlice(
			Drop(-1,
				FromSlice([]int{1, 2, 3, 4})))
	})
}

func TestAll(t *testing.T) {
	assert.False(t,
		All(gvalue.IsZero[int],
			FromSlice([]int{1, 2, 3, 4})))
	assert.True(t,
		All(gvalue.IsZero[int],
			FromSlice([]int{0, 0, 0, 0})))
	assert.True(t,
		All(gvalue.IsZero[int],
			FromSlice([]int{})))
}

func TestAny(t *testing.T) {
	assert.False(t,
		Any(gvalue.IsZero[int],
			FromSlice([]int{1, 2, 3, 4})))
	assert.True(t,
		Any(gvalue.IsZero[int],
			FromSlice([]int{0, 0, 0, 0})))
	assert.True(t,
		Any(gvalue.IsZero[int],
			FromSlice([]int{1, 2, 0, 4})))
	assert.False(t,
		Any(gvalue.IsZero[int],
			FromSlice([]int{})))
}

func TestAnd(t *testing.T) {
	assert.True(t, And(FromSlice([]bool{true, true, true})))
	assert.False(t, And(FromSlice([]bool{true, true, false})))
	assert.True(t, And(FromSlice([]bool{})))
}

func TestOr(t *testing.T) {
	assert.True(t, Or(FromSlice([]bool{true, true, true})))
	assert.False(t, Or(FromSlice([]bool{false, false, false})))
	assert.True(t, Or(FromSlice([]bool{false, false, true})))
	assert.False(t, Or(FromSlice([]bool{})))
}

func TestConcat(t *testing.T) {
	assert.Equal(t, []int{}, ToSlice(Concat[int]()))
	assertSliceEqual(t,
		[]int{0, 1, 2, 3, 4, 5},
		func() Iter[int] {
			return Concat(
				FromSlice([]int{0, 1, 2}),
				FromSlice([]int{3, 4, 5}))
		})
	assertSliceEqual(t,
		[]int{0, 1, 2},
		func() Iter[int] {
			return Concat(
				FromSlice([]int{0, 1, 2}),
				FromSlice([]int{}))
		})
	assertSliceEqual(t,
		[]int{3, 4, 5},
		func() Iter[int] {
			return Concat(
				FromSlice([]int{}),
				FromSlice([]int{3, 4, 5}))
		})
	assertSliceEqual(t,
		[]int{},
		func() Iter[int] {
			return Concat(
				FromSlice([]int{}),
				FromSlice([]int{}))
		})
}

func TestZip(t *testing.T) {
	assertSliceEqual(t,
		[]int{2, 4, 6},
		func() Iter[int] {
			return Zip(gvalue.Add[int],
				ToPeeker(FromSlice([]int{1, 2, 3})),
				ToPeeker(FromSlice([]int{1, 2, 3})))
		})

	// Test leftover elements
	i1 := ToPeeker(FromSlice([]int{1, 2, 3}))
	i2 := ToPeeker(FromSlice([]int{1, 2, 3, 4}))
	assert.Equal(t,
		[]int{2, 4, 6},
		ToSlice(
			Zip(gvalue.Add[int], i1, i2)))
	assert.Equal(t, []int{}, ToSlice[int](i1))
	assert.Equal(t, []int{4}, ToSlice[int](i2))

	// Test leftover elements
	i1 = ToPeeker(FromSlice([]int{1, 2, 3, 4, 5}))
	i2 = ToPeeker(FromSlice([]int{1, 2, 3}))
	assert.Equal(t,
		[]int{2, 4, 6},
		ToSlice(
			Zip(gvalue.Add[int], i1, i2)))
	assert.Equal(t, []int{4, 5}, ToSlice[int](i1))
	assert.Equal(t, []int{}, ToSlice[int](i2))

	// Test infinite elements
	assertSliceEqual(t,
		[]int{2, 2, 2},
		func() Iter[int] {
			return Zip(gvalue.Add[int],
				ToPeeker(FromSlice([]int{1, 1, 1})),
				ToPeeker(Repeat(1)))
		})
}

func TestIntersperse(t *testing.T) {
	assert.Equal(t,
		[]int{1},
		ToSlice(
			Intersperse(0,
				FromSlice([]int{1}))))
	assert.Equal(t,
		[]int{},
		ToSlice(
			Intersperse(0,
				FromSlice([]int{}))))
	assertSliceEqual(t,
		[]int{1, 0, 2, 0, 3, 0, 4, 0, 5},
		func() Iter[int] {
			return Intersperse(0, FromSlice([]int{1, 2, 3, 4, 5}))
		})
}

func TestAppend(t *testing.T) {
	assertSliceEqual(t,
		[]int{1, 2, 3, 4},
		func() Iter[int] {
			return Append(4, FromSlice([]int{1, 2, 3}))
		})
	assertSliceEqual(t,
		[]int{1},
		func() Iter[int] {
			return Append(1, FromSlice([]int{}))
		})
}

func TestPrepend(t *testing.T) {
	assertSliceEqual(t,
		[]int{0, 1, 2, 3},
		func() Iter[int] {
			return Prepend(0, FromSlice([]int{1, 2, 3}))
		})
	assert.Equal(t,
		[]int{1},
		ToSlice(
			Prepend(1,
				FromSlice([]int{}))))
}

// func TestCycle(t *testing.T) {
// 	assert.Equal(t,
// 		[]int{1, 2, 1, 2, 1, 2},
// 		ToSlice(
// 			Take(6,
// 				Cycle(
// 					FromSlice([]int{1, 2})))))
// }

func TestJoin(t *testing.T) {
	assert.Equal(t,
		"1, 2, 3",
		Join(", ", FromSlice([]string{"1", "2", "3"})))
	assert.Equal(t,
		"1",
		Join(", ", FromSlice([]string{"1"})))
	assert.Equal(t,
		"",
		Join(", ", FromSlice([]string{})))
}

func TestTypeAssert(t *testing.T) {
	assert.Equal(t,
		[]int{1, 2, 3},
		ToSlice(
			TypeAssert[int, any](
				FromSlice([]any{1, 2, 3}))))
	assert.Equal(t,
		[]any{1, 2, 3},
		ToSlice(
			TypeAssert[any, int](
				FromSlice([]int{1, 2, 3}))))

	assert.Panic(t, func() {
		ToSlice(
			TypeAssert[string, int](
				FromSlice([]int{1, 2, 3})))
	})

	// Omit original type.
	assert.Equal(t,
		[]int{1, 2, 3},
		ToSlice(
			TypeAssert[int](
				FromSlice([]any{1, 2, 3}))))
	assert.Equal(t,
		[]any{1, 2, 3},
		ToSlice(
			TypeAssert[any](
				FromSlice([]int{1, 2, 3}))))

	assert.Panic(t, func() {
		ToSlice(
			TypeAssert[string](
				FromSlice([]int{1, 2, 3})))
	})
}

func TestCount(t *testing.T) {
	assert.Equal(t, 0, Count(FromSlice([]int{})))
	assert.Equal(t, 1, Count(FromSlice([]int{1})))
	assert.Equal(t, 1000, Count(Take(1000, (Repeat(1)))))
	assert.Equal(t, 1000, Count(Range(0, 1000)))
}

func TestFind(t *testing.T) {
	assert.False(t,
		Find(
			gfunc.Partial2(gvalue.Greater[int]).PartialR(3),
			FromSlice([]int{1, 2, 3})).IsOK())
	assert.True(t,
		Find(
			gfunc.Partial2(gvalue.Less[int]).PartialR(3),
			FromSlice([]int{1, 2, 3})).IsOK())
	assert.Equal(t,
		1,
		Find(
			gfunc.Partial2(gvalue.Less[int]).PartialR(3),
			FromSlice([]int{1, 2, 3})).Value())
	assert.Equal(t,
		3,
		Find(
			gfunc.Partial2(gvalue.GreaterEqual[int]).PartialR(3),
			FromSlice([]int{1, 2, 3})).Value())
}

func TestTakeWhile(t *testing.T) {
	less3 := gfunc.Partial2(gvalue.Less[int]).PartialR(3)
	i1 := ToPeeker(FromSlice([]int{1, 2, 3, 4, 5, 1, 2, 3}))
	assert.Equal(t,
		[]int{1, 2},
		ToSlice(
			TakeWhile(less3, i1)))
	assert.Equal(t,
		[]int{3, 4, 5, 1, 2, 3},
		ToSlice(Iter[int](i1)))

	less9 := gfunc.Partial2(gvalue.Less[int]).PartialR(9)
	i2 := ToPeeker(FromSlice([]int{1, 2, 3}))
	assert.Equal(t,
		[]int{1, 2, 3},
		ToSlice(
			TakeWhile(less9, i2)))
	assert.Equal(t, []int{}, ToSlice(Iter[int](i2)))

	less0 := gfunc.Partial2(gvalue.Less[int]).PartialR(0)
	i3 := ToPeeker(FromSlice([]int{1, 2, 3}))
	assert.Equal(t,
		[]int{},
		ToSlice(
			TakeWhile(less0, i3)))
	assert.Equal(t, []int{1, 2, 3}, ToSlice(Iter[int](i3)))
}

func TestDropWhile(t *testing.T) {
	less3 := gfunc.Partial2(gvalue.Less[int]).PartialR(3)
	assert.Equal(t,
		[]int{3, 4, 5, 1, 2, 3},
		ToSlice(
			DropWhile(less3,
				FromSlice([]int{1, 2, 3, 4, 5, 1, 2, 3}))))

	less9 := gfunc.Partial2(gvalue.Less[int]).PartialR(9)
	assert.Equal(t,
		[]int{},
		ToSlice(
			DropWhile(less9,
				FromSlice([]int{1, 2, 3}))))

	less0 := gfunc.Partial2(gvalue.Less[int]).PartialR(0)
	assert.Equal(t,
		[]int{1, 2, 3},
		ToSlice(
			DropWhile(less0,
				FromSlice([]int{1, 2, 3}))))
}

func TestSort(t *testing.T) {
	assertSliceEqual(t,
		[]int{1, 1, 2, 2, 3, 3, 4, 5},
		func() Iter[int] {
			return Sort(
				FromSlice([]int{1, 2, 3, 4, 5, 1, 2, 3}))
		})
	assertSliceEqual(t,
		[]int{},
		func() Iter[int] {
			return Sort(
				FromSlice([]int{}))
		})
	assertSliceEqual(t,
		[]int{1, 2, 3, 4},
		func() Iter[int] {
			return Sort(
				FromSlice([]int{1, 2, 3, 4}))
		})

	{ // Inplace.
		s := []int{1, 3, 4, 2}
		_ = ToSlice(Sort(StealSlice(s)))
		assert.Equal(t, []int{1, 2, 3, 4}, s)
	}

	// { // Stable.
	// 	stable := true
	// 	for i := 0; i < 10000; i++ {
	// 		v1, v2, v3, v4 := 1, 2, 2, 4
	// 		s := []*int{&v2, &v3, &v4, &v1}
	// 		_ =
	// 			ToSlice(
	// 				SortBy(func(a, b *int) bool { return *a < *b },
	// 					StealSlice(s)))
	// 		if s[1] != &v2 || s[2] != &v3 {
	// 			stable = false
	// 		}
	// 	}
	// 	assert.False(t, stable)
	// }
}

func TestSortBy(t *testing.T) {
	assertSliceEqual(t,
		[]int{1, 1, 2, 2, 3, 3, 4, 5},
		func() Iter[int] {
			return SortBy(gvalue.Less[int],
				FromSlice([]int{1, 2, 3, 4, 5, 1, 2, 3}))
		})
	assertSliceEqual(t,
		[]int{},
		func() Iter[int] {
			return SortBy(gvalue.Less[int],
				FromSlice([]int{}))
		})
	assertSliceEqual(t,
		[]int{1, 2, 3, 4},
		func() Iter[int] {
			return SortBy(gvalue.Less[int],
				FromSlice([]int{1, 2, 3, 4}))
		})
	assertSliceEqual(t,
		[]int{4, 3, 2, 1},
		func() Iter[int] {
			return SortBy(gvalue.Greater[int],
				FromSlice([]int{1, 2, 3, 4}))
		})

	{ // Inplace.
		s := []int{1, 2, 3, 4}
		_ = ToSlice(SortBy(gvalue.Greater[int], StealSlice(s)))
		assert.Equal(t, []int{4, 3, 2, 1}, s)
	}
}

func TestStableSortBy(t *testing.T) {
	assertSliceEqual(t,
		[]int{1, 1, 2, 2, 3, 3, 4, 5},
		func() Iter[int] {
			return StableSortBy(gvalue.Less[int],
				FromSlice([]int{1, 2, 3, 4, 5, 1, 2, 3}))
		})
	assertSliceEqual(t,
		[]int{},
		func() Iter[int] {
			return StableSortBy(gvalue.Less[int],
				FromSlice([]int{}))
		})
	assertSliceEqual(t,
		[]int{1, 2, 3, 4},
		func() Iter[int] {
			return StableSortBy(gvalue.Less[int],
				FromSlice([]int{1, 2, 3, 4}))
		})
	assertSliceEqual(t,
		[]int{4, 3, 2, 1},
		func() Iter[int] {
			return StableSortBy(gvalue.Greater[int],
				FromSlice([]int{1, 2, 3, 4}))
		})

	{ // Inplace.
		s := []int{1, 2, 3, 4}
		_ = ToSlice(StableSortBy(gvalue.Greater[int], StealSlice(s)))
		assert.Equal(t, []int{4, 3, 2, 1}, s)
	}

	for i := 0; i < 1000; i++ { // Stable.
		v1, v2, v3, v4 := 1, 2, 2, 4
		s := []*int{&v2, &v3, &v4, &v1}
		_ =
			ToSlice(
				StableSortBy(func(a, b *int) bool { return *a < *b },
					StealSlice(s)))
		expect := []*int{&v1, &v2, &v3, &v4}
		for i := range expect {
			assert.True(t, expect[i] == s[i])
		}
	}
}

func TestContains(t *testing.T) {
	assert.True(t, Contains(5, FromSlice([]int{1, 2, 3, 4, 5, 1, 2, 3})))
	assert.False(t, Contains(5, FromSlice([]int{})))
	assert.False(t, Contains(-1, Range(1, 10)))
	assert.False(t, Contains(10, Range(1, 10)))
}

func TestContainsAny(t *testing.T) {
	assert.True(t, ContainsAny([]int{1, 2, 6}, FromSlice([]int{1, 2, 3, 4, 5, 1, 2, 3})))
	assert.True(t, ContainsAny([]int{5}, FromSlice([]int{1, 2, 3, 4, 5, 1, 2, 3})))
	assert.False(t, ContainsAny([]int{5}, FromSlice([]int{})))
	assert.False(t, ContainsAny([]int{-1}, Range(1, 10)))
	assert.False(t, ContainsAny([]int{10}, Range(1, 10)))
	assert.False(t, ContainsAny([]int{}, Range(1, 10)))
	assert.False(t, ContainsAny([]int{}, FromSlice([]int{})))
	assert.False(t, ContainsAny(nil, Range(1, 10)))
}

func TestContainsAll(t *testing.T) {
	assert.False(t, ContainsAll([]int{1, 2, 6}, FromSlice([]int{1, 2, 3, 4, 5, 1, 2, 3})))
	assert.True(t, ContainsAll([]int{1, 2, 5}, FromSlice([]int{1, 2, 3, 4, 5, 1, 2, 3})))
	assert.True(t, ContainsAll([]int{5}, FromSlice([]int{1, 2, 3, 4, 5, 1, 2, 3})))
	assert.False(t, ContainsAll([]int{5}, FromSlice([]int{})))
	assert.False(t, ContainsAll([]int{-1}, Range(1, 10)))
	assert.False(t, ContainsAll([]int{10}, Range(1, 10)))
	assert.True(t, ContainsAll([]int{}, Range(1, 10)))
	assert.True(t, ContainsAll(nil, Range(1, 10)))
	assert.True(t, ContainsAll([]int{}, FromSlice([]int{})))

}

func TestUniq(t *testing.T) {
	assertSliceEqual(t,
		[]int{3, 1, 2, 5},
		func() Iter[int] {
			return Uniq(FromSlice([]int{3, 1, 2, 3, 3, 5}))
		})
	assertSliceEqual(t,
		[]int{1},
		func() Iter[int] {
			return Uniq(FromSlice([]int{1}))
		})
	assertSliceEqual(t,
		[]int{},
		func() Iter[int] {
			return Uniq(FromSlice([]int{}))
		})
	assertSliceEqual(t,
		[]int{1},
		func() Iter[int] {
			return Uniq(FromSlice([]int{1, 1, 1, 1, 1}))
		})
}

func TestUniqBy(t *testing.T) {
	type Elem struct {
		Key int
	}
	keyOf := func(e Elem) int {
		return e.Key
	}

	assertSliceEqual(t,
		[]Elem{{3}, {1}, {2}, {5}},
		func() Iter[Elem] {
			return UniqBy(keyOf,
				FromSlice([]Elem{{3}, {1}, {2}, {3}, {3}, {5}}))
		})
	assertSliceEqual(t,
		[]Elem{{1}},
		func() Iter[Elem] {
			return UniqBy(keyOf,
				FromSlice([]Elem{{1}}))
		})
	assertSliceEqual(t,
		[]Elem{},
		func() Iter[Elem] {
			return UniqBy(keyOf,
				FromSlice([]Elem{}))
		})
	assertSliceEqual(t,
		[]Elem{{1}},
		func() Iter[Elem] {
			return UniqBy(keyOf,
				FromSlice([]Elem{{1}, {1}, {1}, {1}, {1}, {1}}))
		})
}

func TestDup(t *testing.T) {
	assertSliceEqual(t,
		[]int{2, 3},
		func() Iter[int] {
			return Dup(FromSlice([]int{3, 1, 2, 2, 2, 3, 5}))
		})
	assertSliceEqual(t,
		[]int{},
		func() Iter[int] {
			return Dup(FromSlice([]int{1}))
		})
	assertSliceEqual(t,
		[]int{},
		func() Iter[int] {
			return Dup(FromSlice([]int{}))
		})
	assertSliceEqual(t,
		[]int{1},
		func() Iter[int] {
			return Dup(FromSlice([]int{1, 1, 1, 1, 1}))
		})
}

func TestDupBy(t *testing.T) {
	type Elem struct {
		Key int
	}
	keyOf := func(e Elem) int {
		return e.Key
	}

	assertSliceEqual(t,
		[]Elem{{2}, {3}},
		func() Iter[Elem] {
			return DupBy(keyOf,
				FromSlice([]Elem{{3}, {2}, {2}, {3}, {3}, {2}}))
		})
	assertSliceEqual(t,
		[]Elem{},
		func() Iter[Elem] {
			return DupBy(keyOf,
				FromSlice([]Elem{{1}}))
		})
	assertSliceEqual(t,
		[]Elem{},
		func() Iter[Elem] {
			return DupBy(keyOf,
				FromSlice([]Elem{}))
		})
	assertSliceEqual(t,
		[]Elem{{1}},
		func() Iter[Elem] {
			return DupBy(keyOf,
				FromSlice([]Elem{{1}, {1}, {1}, {1}, {1}, {1}}))
		})
}

func TestForEach(t *testing.T) {
	var sum1 int
	ForEach(func(v int) { sum1 += v }, FromSlice([]int{1, 2, 3, 4}))
	assert.Equal(t, 10, sum1)
	var sum2 int
	ForEach(func(v int) { sum2 += v }, FromSlice([]int{}))
	assert.Zero(t, sum2)
}

func TestForEachIndexed(t *testing.T) {
	var sum1 int
	ForEachIndexed(func(_, v int) { sum1 += v }, FromSlice([]int{1, 2, 3, 4}))
	assert.Equal(t, 10, sum1)
	var sum2 int
	ForEachIndexed(func(_, v int) { sum2 += v }, FromSlice([]int{}))
	assert.Zero(t, sum2)
	var sum3 int
	ForEachIndexed(func(i, _ int) { sum3 += i }, FromSlice([]int{1, 2, 3, 4}))
	assert.Equal(t, 6, sum3)
}

func TestSum(t *testing.T) {
	assert.Equal(t, 0, Sum(FromSlice([]int{})))
	assert.Equal(t, 1, Sum(FromSlice([]int{1})))
	assert.Equal(t, 1000, Sum(Take(1000, (Repeat(1)))))
	assert.Equal(t, 5050, Sum(Range(0, 100+1)))
}

func TestSumBy(t *testing.T) {
	type Foo struct {
		Value int
	}
	// FIXME: use greflect
	getValue := func(foo Foo) int {
		return foo.Value
	}

	assert.Equal(t, 0, SumBy(getValue, FromSlice([]Foo{})))
	assert.Equal(t, 1, SumBy(getValue, FromSlice([]Foo{{1}})))

	var foos1 []Foo
	for i := 0; i < 1000; i++ {
		foos1 = append(foos1, Foo{1})
	}
	assert.Equal(t, 1000, SumBy(getValue, FromSlice(foos1)))
	var foos2 []Foo
	for i := 0; i <= 100; i++ {
		foos2 = append(foos2, Foo{i})
	}
	assert.Equal(t, 5050, SumBy(getValue, FromSlice(foos2)))
}

func TestAvg(t *testing.T) {
	assert.Equal(t, 0.0, Avg(FromSlice([]int{})))
	assert.Equal(t, 1.0, Avg(FromSlice([]int{1})))
	assert.Equal(t, 1.0, Avg(Take(1000, (Repeat(1)))))
	assert.Equal(t, 50.0, Avg(Range(0, 100+1)))

	assert.Equal(t, 1.0, Avg(FromSlice([]float64{1})))
	assert.Equal(t, 0.5, Avg(Range(0.0, 2.0)))

	// Test overflow
	upBound := []int8{127, 127, 127}
	assert.Equal(t, 127.0, Avg(FromSlice(upBound)))
	assert.NotEqual(t, float64(Sum(FromSlice(upBound)))/3.0, Avg(FromSlice(upBound)))
}

func TestAvgBy(t *testing.T) {
	type Foo struct {
		Value int
	}
	// FIXME: use greflect
	getValue := func(foo Foo) int {
		return foo.Value
	}
	assert.Equal(t, 0.0, AvgBy(getValue, FromSlice([]Foo{})))
	assert.Equal(t, 1.0, AvgBy(getValue, FromSlice([]Foo{{1}})))
	var foos1 []Foo
	for i := 0; i < 1000; i++ {
		foos1 = append(foos1, Foo{1})
	}
	assert.Equal(t, 1.0, AvgBy(getValue, FromSlice(foos1)))
	var foos2 []Foo
	for i := 0; i <= 100; i++ {
		foos2 = append(foos2, Foo{i})
	}
	assert.Equal(t, 50.0, AvgBy(getValue, FromSlice(foos2)))
}

func TestAt(t *testing.T) {
	assert.Panic(t, func() {
		At(-1, Range(1, 1000))
	})
	assert.Equal(t, 11, At(10, Range(1, 1000)).Value())
	assert.False(t, At(10, Range(1, 10)).IsOK())
	assert.Equal(t, 1, At(0, FromSlice([]int{1, 2, 3, 4})).Value())
	assert.Equal(t, 2, At(1, FromSlice([]int{1, 2, 3, 4})).Value())
	assert.Equal(t, 3, At(2, FromSlice([]int{1, 2, 3, 4})).Value())
	assert.Equal(t, 4, At(3, FromSlice([]int{1, 2, 3, 4})).Value())
	assert.False(t, At(4, FromSlice([]int{1, 2, 3, 4})).IsOK())
	assert.False(t, At(10000, FromSlice([]int{1, 2, 3, 4})).IsOK())
}

func assertSliceEqual[T any](t *testing.T, s []T, f func() Iter[T]) {
	// Mark as test helper.
	t.Helper()

	// Read nothing.
	if !assert.Zero(t, len(f().Next(0))) {
		t.Log("Read 0 element")
	}

	// Read all.
	if !assert.Equal(t, s, ToSlice(f())) {
		t.Log("Read ALL element")
	}

	// Read a part.
	for i := 1; i <= len(s); i++ {
		n := i - 1
		if !assert.Equal(t, s[:n], ToSlice(Take(n, f()))) {
			t.Logf("Read %d elements", n)
			t.FailNow()
		}
	}

	// Read multiple times.
	for i := 0; i < len(s); i++ {
		for times := 1; times < len(s); times++ {
			it := f()
			curTimes := 0
			ptr := 0
			size := len(s) / times
			// Read multiple times.
			for ptr+size < len(s) {
				curTimes++
				if !assert.Equal(t, s[ptr:ptr+size], ToSlice(Take(size, it))) {
					t.Logf("Read %d/%d times, range: %d:%d", curTimes, times, ptr, ptr+size)
					t.FailNow()
				}
				ptr += size
			}
			if ptr < len(s) {
				curTimes++
				// Read remaining.
				if !assert.Equal(t, s[ptr:], ToSlice(it)) {
					t.Logf("Read %d/%d times, range: %d:ALL", curTimes, times, ptr)
					t.FailNow()
				}
			}
		}
	}

	// Read overflow.
	for i := 0; i < 100; i++ {
		n := len(s) + i
		if !assert.Equal(t, s, ToSlice(Take(n, f()))) {
			t.Logf("Read overflow, size: %d", n)
			t.FailNow()
		}
	}
}

func TestGroupBy(t *testing.T) {
	assert.Equal(t,
		map[string][]int{"odd": {1, 3, 5}, "even": {2, 4, 6}},
		GroupBy(func(i int) string { return gcond.If(i%2 == 0, "even", "odd") },
			FromSlice([]int{1, 2, 3, 4, 5, 6})))

	// Test empty slice
	assert.Equal(t,
		map[string][]int{},
		GroupBy(func(i int) string { return gcond.If(i%2 == 0, "even", "odd") },
			FromSlice([]int{})))
}

func TestRemove(t *testing.T) {
	// Empty.
	assertSliceEqual(t,
		[]int{},
		func() Iter[int] {
			return Remove(0, FromSlice([]int{}))
		})
	// All removed.
	assertSliceEqual(t,
		[]int{},
		func() Iter[int] {
			return Remove(1, FromSlice([]int{1, 1, 1, 1, 1}))
		})
	// Nothing removed.
	assertSliceEqual(t,
		[]int{1, 1, 1, 1, 1},
		func() Iter[int] {
			return Remove(0, FromSlice([]int{1, 1, 1, 1, 1}))
		})
	assertSliceEqual(t,
		[]int{2, 3, 4},
		func() Iter[int] {
			return Remove(1, FromSlice([]int{1, 2, 3, 4, 1}))
		})
}

func TestRemoveN(t *testing.T) {
	// Empty.
	assertSliceEqual(t,
		[]int{},
		func() Iter[int] {
			return RemoveN(0, 1, FromSlice([]int{}))
		})
	// All removed.
	assertSliceEqual(t,
		[]int{},
		func() Iter[int] {
			return RemoveN(1, 100, FromSlice([]int{1, 1, 1, 1, 1}))
		})
	// 5 removed.
	assertSliceEqual(t,
		[]int{},
		func() Iter[int] {
			return RemoveN(1, 5, FromSlice([]int{1, 1, 1, 1, 1}))
		})
	// 3 removed.
	assertSliceEqual(t,
		[]int{1, 1},
		func() Iter[int] {
			return RemoveN(1, 3, FromSlice([]int{1, 1, 1, 1, 1}))
		})
	// Nothing removed.
	assertSliceEqual(t,
		[]int{1, 1, 1, 1, 1},
		func() Iter[int] {
			return RemoveN(0, 1, FromSlice([]int{1, 1, 1, 1, 1}))
		})
	// Nothing removed.
	assertSliceEqual(t,
		[]int{1, 1, 1, 1, 1},
		func() Iter[int] {
			return RemoveN(1, 0, FromSlice([]int{1, 1, 1, 1, 1}))
		})
	assertSliceEqual(t,
		[]int{2, 3, 4},
		func() Iter[int] {
			return RemoveN(1, 2, FromSlice([]int{1, 2, 3, 4, 1}))
		})
	assertSliceEqual(t,
		[]int{2, 3, 4, 1},
		func() Iter[int] {
			return RemoveN(1, 1, FromSlice([]int{1, 2, 3, 4, 1}))
		})
}

func TestChunk(t *testing.T) {
	// Empty.
	assertSliceEqual(t,
		[][]int{},
		func() Iter[[]int] {
			return Chunk(3, FromSlice([]int{}))
		})
	assertSliceEqual(t,
		[][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
		func() Iter[[]int] {
			return Chunk(3, FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}))
		})
	assertSliceEqual(t,
		[][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10}},
		func() Iter[[]int] {
			return Chunk(3, FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}))
		})
	assertSliceEqual(t,
		[][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10}},
		func() Iter[[]int] {
			return Chunk(8, FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}))
		})
	assertSliceEqual(t,
		[][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		func() Iter[[]int] {
			return Chunk(100, FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}))
		})
	assertSliceEqual(t,
		[][]int{},
		func() Iter[[]int] {
			return Chunk(0, FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}))
		})
}

func TestDivide(t *testing.T) {
	// Invalid N.
	assert.Panic(t, func() {
		_ = Divide(0, FromSlice([]int{}))
		_ = Divide(-1, FromSlice([]int{}))
	})
	// Empty.
	assertSliceEqual(t,
		[][]int{{}, {}, {}},
		func() Iter[[]int] {
			return Divide(3, FromSlice([]int{}))
		})

	// Divide 10 elements.
	assertSliceEqual(t,
		[][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		func() Iter[[]int] {
			return Divide(1, FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}))
		})
	assertSliceEqual(t,
		[][]int{{1, 2, 3, 4, 5}, {6, 7, 8, 9, 10}},
		func() Iter[[]int] {
			return Divide(2, FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}))
		})
	assertSliceEqual(t,
		[][]int{{1, 2, 3, 4}, {5, 6, 7}, {8, 9, 10}},
		func() Iter[[]int] {
			return Divide(3, FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}))
		})
	assertSliceEqual(t,
		[][]int{{1, 2, 3}, {4, 5, 6}, {7, 8}, {9, 10}},
		func() Iter[[]int] {
			return Divide(4, FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}))
		})
	assertSliceEqual(t,
		[][]int{{1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {9}, {10}},
		func() Iter[[]int] {
			return Divide(10, FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}))
		})
	assertSliceEqual(t,
		[][]int{{1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {9}, {10}, {}},
		func() Iter[[]int] {
			return Divide(11, FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}))
		})
}

func TestShuffle(t *testing.T) {
	for i := 0; i < 100; i++ {
		s := ToSlice(Shuffle(Range(0, i)))
		assert.Equal(t, i, len(s))
		if i > 10 {
			// Should not equal in most time?
			assert.NotEqual(t, ToSlice(Range(0, i)), s)
		}
		for j := 0; j < i; j++ {
			var found bool
			for _, v := range s {
				if v == j {
					found = true
					break
				}
			}
			assert.True(t, found)
		}
	}
}

func TestCompact(t *testing.T) {
	assertSliceEqual(t,
		[]int{1, 2, 3, 4},
		func() Iter[int] {
			return Compact(FromSlice([]int{1, 2, 0, 3, 0, 4}))
		})

	// All zero.
	assertSliceEqual(t,
		[]string{},
		func() Iter[string] {
			return Compact(FromSlice(make([]string, 1000)))
		})

	// Empty
	assertSliceEqual(t,
		[]float64{},
		func() Iter[float64] {
			return Compact(FromSlice([]float64{}))
		})

	// All non-zero.
	assertSliceEqual(t,
		ToSlice(Range(1, 100)),
		func() Iter[int] {
			return Compact(Range(1, 100))
		})
}

func BenchmarkContainsAll_Delete(b *testing.B) {

	b.Run("int_inspect_if", func(b *testing.B) {
		mp := make(map[int]int, b.N)
		for i := 0; i < b.N; i++ {
			mp[i] = i
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			deleted := (i%2 == 0)
			if deleted {
				if _, ok := mp[i]; ok {
					delete(mp, i)
				}
			}
		}
	})
	b.Run("int_no_if", func(b *testing.B) {
		mp := make(map[int]int, b.N)
		for i := 0; i < b.N; i++ {
			mp[i] = i
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			deleted := (i%2 == 0)
			if deleted {
				delete(mp, i)
			}
		}
	})

}
