package gvalue

import (
	"fmt"
	"math"
	"net"
	"sync"
	"testing"

	"github.com/bytedance/gg/internal/assert"
)

func TestZero(t *testing.T) {
	assert.Zero(t, Zero[bool]())
	assert.Zero(t, Zero[int]())
	assert.Zero(t, Zero[*int]())
	assert.Zero(t, Zero[string]())
	assert.Zero(t, Zero[interface{}]())
	assert.Zero(t, Zero[*interface{}]())
}

func TestMin(t *testing.T) {
	assert.Equal(t, 1, Min(1, 2))
	assert.Equal(t, 2, Min(2))
	assert.Equal(t, 1, Min(2, 1, 3))
	assert.Equal(t, 1, Min(1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3))
	assert.Equal(t, -2147483648, Min(math.MaxInt32, 0, math.MaxInt64, math.MinInt32))
	assert.Equal(t, -1.0, Min[float32](2, -1))
	assert.Equal(t, math.E, Min(math.E, 3.0, 2.8))
	assert.Equal(t, math.E, Min(3.0, math.E, 2.8))
	assert.Equal(t, "1", Min("1"))
	assert.Equal(t, "", Min("    ", "", "  "))
	assert.Equal(t, "1099", Min("1999", "2", "1099"))
	assert.Equal(t, "\nzzz", Min("a", "1999", "2", "\nzzz"))
}

func TestMax(t *testing.T) {
	assert.Equal(t, 1, Max(1))
	assert.Equal(t, 1, Max(0, 1, 0, -1))
	assert.Equal(t, 3, Max(1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3))
	assert.Equal(t, 2147483647, Max[int64](0, math.MinInt64, math.MaxInt32, 0))
	assert.Equal(t, -1.0, Max[float32](-1, -2))
	assert.Equal(t, math.E, Max(2.0, math.E, 2.718))
	assert.Equal(t, "1", Max("1"))
	assert.Equal(t, "    ", Max("    ", "", "  "))
	assert.Equal(t, "2", Max("1999", "2", "1099"))
	assert.Equal(t, "a", Max("a", "1999", "2  ", "\nzzz"))
}

type Pair[T1, T2 any] struct {
	First  T1
	Second T2
}

func MakePair[T1, T2 any](first T1, second T2) Pair[T1, T2] {
	return Pair[T1, T2]{first, second}
}

func TestMinMax(t *testing.T) {
	assert.Equal(t, MakePair(1, 1), MakePair(MinMax(1)))
	assert.Equal(t, MakePair(-1, 1), MakePair(MinMax(0, 1, 0, -1)))
	assert.Equal(t, MakePair(1, 3), MakePair(MinMax(1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3)))
	assert.Equal(t, MakePair[int64, int64](-9223372036854775808, 2147483647), MakePair(MinMax[int64](0, math.MinInt64, math.MaxInt32, 0)))
	assert.Equal(t, MakePair[float32, float32](-2.0, -1.0), MakePair(MinMax[float32](-1, -2)))
	assert.Equal(t, MakePair(2.0, math.E), MakePair(MinMax(2.0, math.E, 2.718)))
	assert.Equal(t, MakePair("1", "1"), MakePair(MinMax("1")))
	assert.Equal(t, MakePair("", "    "), MakePair(MinMax("    ", "", "  ")))
	assert.Equal(t, MakePair("1099", "2"), MakePair(MinMax("1999", "2", "1099")))
	assert.Equal(t, MakePair("\nzzz", "a"), MakePair(MinMax("a", "1999", "2  ", "\nzzz")))
}

func TestClamp(t *testing.T) {
	assert.Equal(t, 2, Clamp(1, 2, 3))
	assert.Equal(t, 2, Clamp(2, 1, 3))
	assert.Equal(t, 2, Clamp(3, 1, 2))
	assert.Equal(t, "11", Clamp("2", "10", "11"))
	assert.Equal(t, 0, Clamp[int64](0, math.MinInt64, math.MaxInt64))
	assert.Equal(t, math.MinInt64, Clamp[int64](math.MinInt64, math.MinInt64, math.MaxInt64))
	assert.Equal(t, math.MaxInt64, Clamp[int64](math.MaxInt64, math.MinInt64, math.MaxInt64))
	assert.Equal(t, -1.0, Clamp[float64](-1e9, -1.0, 1.0))
	assert.Equal(t, "   ", Clamp[string]("", "   ", "     "))
}

func TestIsNil(t *testing.T) {
	assert.False(t, IsNil(1))
	ii := 1
	assert.False(t, IsNil(&ii))
	assert.False(t, &ii == nil)
	assert.True(t, IsNil(nil))
	// assert.True(t, nil == nil)

	// Nil
	var i *int
	assert.True(t, IsNil(i))
	assert.True(t, i == nil)
	assert.True(t, IsNil(Zero[*int]()))
	assert.True(t, IsNil((*int)(nil)))
	assert.True(t, (*int)(nil) == nil)

	// Interface
	var ip *net.IP
	assert.True(t, IsNil(fmt.Stringer(ip)))
	assert.True(t, ip == nil)
	assert.True(t, IsNil(fmt.Stringer((*net.IP)(nil))))
	assert.False(t, fmt.Stringer((*net.IP)(nil)) == nil)
	var s fmt.Stringer
	assert.True(t, IsNil(s))
	assert.True(t, s == nil)
	s = ip
	assert.True(t, IsNil(s))
	assert.False(t, s == nil)
	s = &net.IP{}
	assert.False(t, IsNil(s))
	assert.False(t, s == nil)
}

func TestIsNotNil(t *testing.T) {
	assert.True(t, IsNotNil(1))
	var i *int
	assert.False(t, IsNotNil(i))
}

func TestIsZero(t *testing.T) {
	assert.True(t, IsZero(0))
	assert.False(t, IsZero(1))

	assert.True(t, IsZero(""))
	assert.False(t, IsZero("0"))

	assert.True(t, IsZero[*int](nil))
	i := 0
	assert.False(t, IsZero(&i))
}

func TestIsNotZero(t *testing.T) {
	assert.False(t, IsNotZero(0))
	assert.True(t, IsNotZero(1))
}

func TestEqual(t *testing.T) {
	assert.True(t, Equal(1, 1))
	assert.False(t, Equal(1, 0))

	assert.True(t, Equal("a", "a"))
	assert.False(t, Equal("a", "A"))
}

func TestAdd(t *testing.T) {
	assert.Equal(t, 2, Add(1, 1))
	assert.Equal(t, "Alice", Add("Ali", "ce"))
}

func TestLess(t *testing.T) {
	assert.True(t, Less(1, 2))
	assert.False(t, Less(2, 1))
	assert.False(t, Less(1, 1))

	assert.True(t, Less("1", "2"))
	assert.False(t, Less("2", "1"))
	assert.False(t, Less("1", "1"))
}

func TestLessEqual(t *testing.T) {
	assert.True(t, LessEqual(1, 2))
	assert.False(t, LessEqual(2, 1))
	assert.True(t, LessEqual(1, 1))

	assert.True(t, LessEqual("1", "2"))
	assert.False(t, LessEqual("2", "1"))
	assert.True(t, LessEqual("1", "1"))
}

func TestGreater(t *testing.T) {
	assert.True(t, Greater(2, 1))
	assert.False(t, Greater(1, 2))
	assert.False(t, Greater(1, 1))

	assert.True(t, Greater("2", "1"))
	assert.False(t, Greater("1", "2"))
	assert.False(t, Greater("1", "1"))
}

func TestGreaterEqual(t *testing.T) {
	assert.True(t, GreaterEqual(2, 1))
	assert.False(t, GreaterEqual(1, 2))
	assert.True(t, GreaterEqual(1, 1))

	assert.True(t, GreaterEqual("2", "1"))
	assert.False(t, GreaterEqual("1", "2"))
	assert.True(t, GreaterEqual("1", "1"))
}

func TestBetween(t *testing.T) {
	assert.True(t, Between(2, 1, 2))
	assert.False(t, Between(1, 2, 3))
	assert.True(t, Between(1, 1, 1))

	assert.True(t, Between("2", "1", "2"))
	assert.False(t, Between("1", "2", "3"))
	assert.True(t, Between("1", "1", "1"))
}

func TestTypeAssert(t *testing.T) {
	assert.Equal(t, any(1), TypeAssert[any, int](1))
	assert.Equal(t, 1, TypeAssert[int, any](any(1)))

	// Omit original type.
	assert.Equal(t, any(1), TypeAssert[any](1))
	assert.Equal(t, 1, TypeAssert[int](any(1)))

	assert.Panic(t, func() {
		TypeAssert[float64](any(1))
	})
}

func TestTryAssert(t *testing.T) {
	assert.Equal(t, MakePair(any(1), true), MakePair(TryAssert[any, int](1)))
	assert.Equal(t, MakePair(1, true), MakePair(TryAssert[int, any](any(1))))

	// Omit original type.
	assert.Equal(t, MakePair(any(1), true), MakePair(TryAssert[any](1)))
	assert.Equal(t, MakePair(1, true), MakePair(TryAssert[int](any(1))))

	// Assert failed.
	assert.Equal(t, MakePair(float64(0), false), MakePair(TryAssert[float64](any(1))))

	assert.NotPanic(t, func() {
		TryAssert[float64](any(1))
	})
}

func TestOnce(t *testing.T) {
	{
		i := 2
		f := func() int {
			i++
			return i
		}
		assert.Equal(t, i, 2)
		once := Once(f)
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
		once := Once(f)
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
