package jsonbuilder

import (
	"testing"

	"github.com/bytedance/gg/internal/assert"
)

func TestArrayBuild(t *testing.T) {
	{
		s := []int{1, 2, 3, 4, 5}
		a := NewArray()
		for _, v := range s {
			err := a.Append(v)
			assert.Nil(t, err)
		}
		bs, err := a.Build()
		assert.Nil(t, err)
		assert.Equal(t, []byte(`[1,2,3,4,5]`), bs)
	}

	{
		s := []int{}
		a := NewArray()
		for _, v := range s {
			err := a.Append(v)
			assert.Nil(t, err)
		}
		bs, err := a.Build()
		assert.Nil(t, err)
		assert.Equal(t, []byte(`[]`), bs)
	}

	{
		s := []string{"a"}
		a := NewArray()
		for _, v := range s {
			err := a.Append(v)
			assert.Nil(t, err)
		}
		bs, err := a.Build()
		assert.Nil(t, err)
		assert.Equal(t, []byte(`["a"]`), bs)
	}

	{
		var a *Array
		bs, err := a.Build()
		assert.Nil(t, err)
		assert.Equal(t, []byte(`null`), bs)
	}
}
