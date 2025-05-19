package tuple

import (
	"testing"

	"github.com/bytedance/gg/internal/assert"
)

func TestT2(t *testing.T) {
	p := Make2("red", 14)

	if p.First != "red" {
		t.Error()
	}
	if p.Second != 14 {
		t.Error()
	}
}

func TestS2(t *testing.T) {
	{
		s := Zip2([]string{"red", "green", "blue"}, []int{14, 15, 16})
		s1, s2 := s.Unzip()
		assert.Equal(t, []string{"red", "green", "blue"}, s1)
		assert.Equal(t, []int{14, 15, 16}, s2)
	}
	{ // Test empty.
		s := Zip2([]string{}, []int{})
		s1, s2 := s.Unzip()
		assert.Equal(t, []string{}, s1)
		assert.Equal(t, []int{}, s2)
	}
	{ // Test nil.
		s := Zip2([]string(nil), []int(nil))
		s1, s2 := s.Unzip()
		assert.Equal(t, []string{}, s1)
		assert.Equal(t, []int{}, s2)
	}
}
