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
