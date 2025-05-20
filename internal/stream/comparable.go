package stream

import (
	"github.com/bytedance/gg/iter"
)

// See function [github.com/bytedance/gg/iter.FromMapKeys].
func FromMapKeys[T comparable, I any](m map[T]I) Comparable[T] {
	return Comparable[T]{FromIter(iter.FromMapKeys(m))}
}

// See function [github.com/bytedance/gg/iter.Contains].
func (s Comparable[T]) Contains(v T) bool {
	return iter.Contains(v, s.Iter)
}

// See function [github.com/bytedance/gg/iter.ContainsAny].
func (s Comparable[T]) ContainsAny(vs ...T) bool {
	return iter.ContainsAny(vs, s.Iter)
}

// See function [github.com/bytedance/gg/iter.ContainsAll].
func (s Comparable[T]) ContainsAll(vs ...T) bool {
	return iter.ContainsAll(vs, s.Iter)
}

// See function [github.com/bytedance/gg/iter.Uniq].
func (s Comparable[T]) Uniq() Comparable[T] {
	return FromComparableIter(iter.Uniq(s.Iter))
}

// See function [github.com/bytedance/gg/iter.Remove].
func (s Comparable[T]) Remove(v T) Comparable[T] {
	return FromComparableIter(iter.Remove(v, s.Iter))
}

// See function [github.com/bytedance/gg/iter.RemoveN].
func (s Comparable[T]) RemoveN(v T, n int) Comparable[T] {
	return FromComparableIter(iter.RemoveN(v, n, s.Iter))
}
