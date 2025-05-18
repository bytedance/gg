package gptr

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/bytedance/gg/gvalue"
	"github.com/bytedance/gg/internal/assert"
)

func TestOf(t *testing.T) {
	assert.Equal(t, 543, *Of(543))
	assert.Equal(t, "Alice", *Of("Alice"))
	assert.Equal(t, "Alice", **Of(Of("Alice")))
	assert.Equal(t, "Alice", ***Of(Of(Of("Alice"))))
	assert.False(t, IsNil(Of[*int](nil)))
	assert.True(t, IsNil(*Of[*int](nil)))
	// assert.Nil(t, *Of[*interface{}](nil))

	// Test modifying pointer.
	{
		v := 1
		p := Of(v)
		assert.False(t, p == &v)
		*p = 2
		assert.Equal(t, 1, v)
		assert.Equal(t, 2, *p)
	}
}

func TestOfNotZero(t *testing.T) {
	assert.Equal(t, 543, *OfNotZero(543))
	assert.Equal(t, "Alice", *OfNotZero("Alice"))

	// Test zero.
	assert.True(t, IsNil(OfNotZero(0)))
	assert.True(t, IsNil(OfNotZero("")))
	assert.True(t, IsNil(OfNotZero[*int](nil)))
}

func TestOfPositive(t *testing.T) {
	assert.Equal(t, 543, *OfPositive(543))
	assert.Equal(t, 1.23, *OfPositive(1.23))

	// Test non-positive number.
	assert.True(t, IsNil(OfPositive(0)))
	assert.True(t, IsNil(OfPositive(-1)))
	assert.True(t, IsNil(OfPositive(-1.23)))
}

func TestIndirect(t *testing.T) {
	assert.Equal(t, 543, Indirect(Of(543)))
	assert.Equal(t, "Alice", Indirect(Of("Alice")))
	assert.Zero(t, Indirect[int](nil))
	assert.Nil(t, Indirect[interface{}](nil))
	assert.Nil(t, Indirect(Of[fmt.Stringer](nil)))
}

func TestIndirectOr(t *testing.T) {
	assert.Equal(t, "Alice", IndirectOr(Of("Alice"), "Bob"))
	assert.Equal(t, "Bob", IndirectOr(nil, "Bob"))
}

func TestIsNil(t *testing.T) {
	assert.False(t, IsNil(Of(1)))
	assert.True(t, IsNil[int](nil))
}

func TestEqual(t *testing.T) {
	ptr := Of(1)
	assert.True(t, Equal(ptr, ptr))
	assert.True(t, Equal(Of(1), Of(1)))
	assert.False(t, Equal(Of(1), Of(2)))
	assert.False(t, Equal(Of(1), nil))
	assert.False(t, Equal(nil, Of(1)))
	assert.True(t, Equal[string](nil, nil))
}

func TestEqualTo(t *testing.T) {
	assert.True(t, EqualTo(Of(1), 1))
	assert.False(t, EqualTo(Of(2), 1))
	assert.False(t, EqualTo(nil, 0))
}

func TestClone(t *testing.T) {
	assert.True(t, IsNil(Clone(((*int)(nil)))))

	v := 1
	assert.True(t, Clone(&v) != &v)
	assert.True(t, Equal(Clone(&v), &v))

	src := Of(1)
	dst := Clone(&src)
	assert.Equal(t, &src, dst)
	assert.True(t, src == *dst)
}

func TestCloneBy(t *testing.T) {
	assert.True(t, IsNil(CloneBy(((**int)(nil)), Clone[int])))

	src := Of(1)
	dst := CloneBy(&src, Clone[int])
	assert.Equal(t, &src, dst)
	assert.False(t, src == *dst)
}

func TestMap(t *testing.T) {
	i := 1
	assert.Equal(t, Of("1"), Map(&i, strconv.Itoa))
	assert.True(t, Map(nil, strconv.Itoa) == nil)

	assert.NotPanic(t, func() {
		_ = Map(nil, func(int) string {
			panic("Q_Q")
		})
	})

	assert.Panic(t, func() {
		_ = Map(&i, func(int) string {
			panic("Q_Q")
		})
	})
}

func Indirect_gvalueZero[T any](p *T) (v T) {
	if p == nil {
		return gvalue.Zero[T]()
	}
	return *p
}

func BenchmarkIndirect(b *testing.B) {
	type Big struct {
		Foo [200]string
		Bar int
	}

	var big *Big
	b.Run("Named", func(b *testing.B) {
		var v Big
		for i := 0; i <= b.N; i++ {
			v = Indirect(big)
		}
		_ = v
	})
	b.Run("gvalue.Zero", func(b *testing.B) {
		var v Big
		for i := 0; i <= b.N; i++ {
			v = Indirect_gvalueZero(big)
		}
		_ = v
	})
}
