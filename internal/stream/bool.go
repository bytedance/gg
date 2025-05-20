package stream

import (
	"github.com/bytedance/gg/iter"
)

// See function [github.com/bytedance/gg/iter.And].
func (s Bool[T]) And() bool {
	return iter.And(s.Iter)
}

// See function [github.com/bytedance/gg/iter.Or].
func (s Bool[T]) Or() bool {
	return iter.Or(s.Iter)
}
