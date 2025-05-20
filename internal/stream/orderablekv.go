package stream

import (
	"github.com/bytedance/gg/collection/tuple"
	"github.com/bytedance/gg/internal/constraints"
	"github.com/bytedance/gg/iter"
)

// See function [github.com/bytedance/gg/iter.FromMap].
func FromOrderableMap[K constraints.Ordered, V any](m map[K]V) OrderableKV[K, V] {
	return OrderableKV[K, V]{FromMap(m)}
}

// See function [github.com/bytedance/gg/iter.Sort].
func (s OrderableKV[K, V]) Sort() OrderableKV[K, V] {
	less := func(x, y tuple.T2[K, V]) bool { return x.First < y.First }
	return FromOrderableKVIter(iter.SortBy(less, s.Iter))
}

// Keys returns stream of key.
func (s OrderableKV[K, V]) Keys() Orderable[K] {
	return FromOrderableIter(iter.Map(func(v tuple.T2[K, V]) K {
		return v.First
	}, s.Iter))
}
