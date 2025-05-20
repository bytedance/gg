package stream

import (
	"github.com/bytedance/gg/collection/tuple"
	"github.com/bytedance/gg/goption"
	"github.com/bytedance/gg/iter"
)

// See function [github.com/bytedance/gg/iter.Max].
func (s Orderable[T]) Max() goption.O[T] {
	return iter.Max(s.Iter)
}

// See function [github.com/bytedance/gg/iter.Min].
func (s Orderable[T]) Min() goption.O[T] {
	return iter.Min(s.Iter)
}

// See function [github.com/bytedance/gg/iter.MinMax].
func (s Orderable[T]) MinMax() goption.O[tuple.T2[T, T]] {
	return iter.MinMax(s.Iter)
}

// See function [github.com/bytedance/gg/iter.Sort].
func (s Orderable[T]) Sort() Orderable[T] {
	return FromOrderableIter(iter.Sort(s.Iter))
}
