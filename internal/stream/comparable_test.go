package stream

import (
	"testing"

	"github.com/bytedance/gg/internal/assert"
)

func TestComparable_Contains(t *testing.T) {
	assert.True(t, FromComparableSlice([]int{1, 2, 3}).Contains(3))
	assert.False(t, FromComparableSlice([]int{1, 2, 3}).Contains(4))
}

func TestComparable_ContainsAny(t *testing.T) {
	assert.True(t, FromComparableSlice([]int{1, 2, 3}).ContainsAny(3))
	assert.False(t, FromComparableSlice([]int{1, 2, 3}).ContainsAny(4))
	assert.True(t, FromComparableSlice([]int{1, 2, 3}).ContainsAny(1, 2))
	assert.True(t, FromComparableSlice([]int{1, 2, 3}).ContainsAny(3, 4))
}

func TestComparable_ContainsAll(t *testing.T) {
	assert.True(t, FromComparableSlice([]int{1, 2, 3}).ContainsAll(3))
	assert.False(t, FromComparableSlice([]int{1, 2, 3}).ContainsAll(4))
	assert.True(t, FromComparableSlice([]int{1, 2, 3}).ContainsAll(1, 2))
	assert.False(t, FromComparableSlice([]int{1, 2, 3}).ContainsAll(3, 4))
}

func TestComparable_Uniq(t *testing.T) {
	assert.Equal(t,
		[]int{1, 2, 3},
		FromComparableSlice([]int{1, 2, 2, 3}).Uniq().ToSlice())
}

func TestComparable_Remove(t *testing.T) {
	assert.Equal(t,
		[]int{1, 3},
		FromComparableSlice([]int{1, 2, 2, 3}).Remove(2).ToSlice())
}

func TestComparable_RemoveN(t *testing.T) {
	assert.Equal(t,
		[]int{1, 2, 3},
		FromComparableSlice([]int{1, 2, 2, 3}).RemoveN(2, 1).ToSlice())
}
