package stream

import (
	"testing"

	"github.com/bytedance/gg/internal/assert"
)

func TestBool_And(t *testing.T) {
	assert.True(t, FromBoolSlice([]bool{true, true, true}).And())
}

func TestBool_Or(t *testing.T) {
	assert.True(t, FromBoolSlice([]bool{false, false, true}).Or())
}
