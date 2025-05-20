// Copyright 2025 Bytedance Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package stream

import (
	"github.com/bytedance/gg/collection/tuple"
	"github.com/bytedance/gg/internal/constraints"
	"github.com/bytedance/gg/internal/iter"
)

// See function [github.com/bytedance/gg/internal/iter.FromMap].
func FromOrderableMap[K constraints.Ordered, V any](m map[K]V) OrderableKV[K, V] {
	return OrderableKV[K, V]{FromMap(m)}
}

// See function [github.com/bytedance/gg/internal/iter.Sort].
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
