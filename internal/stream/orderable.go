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
