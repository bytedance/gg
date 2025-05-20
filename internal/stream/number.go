package stream

import (
	"github.com/bytedance/gg/internal/constraints"
	"github.com/bytedance/gg/iter"
)

// See function [github.com/bytedance/gg/iter.Range].
func Range[T constraints.Number](start, stop T) Number[T] {
	return FromNumberIter(iter.Range(start, stop))
}

// See function [github.com/bytedance/gg/iter.RangeWithStep].
func RangeWithStep[T constraints.Number](start, stop, step T) Number[T] {
	return FromNumberIter(iter.RangeWithStep(start, stop, step))
}

// See function [github.com/bytedance/gg/iter.Sum].
func (s Number[T]) Sum() T {
	return iter.Sum(s.Iter)
}

// See function [github.com/bytedance/gg/iter.Avg].
func (s Number[T]) Avg() float64 {
	return iter.Avg(s.Iter)
}
