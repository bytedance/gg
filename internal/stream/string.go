package stream

import (
	"github.com/bytedance/gg/iter"
)

// See function [github.com/bytedance/gg/iter.Join].
func (s String[T]) Join(sep T) T {
	return iter.Join(sep, s.Iter)
}
