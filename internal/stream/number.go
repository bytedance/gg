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
	"github.com/bytedance/gg/internal/constraints"
	"github.com/bytedance/gg/internal/iter"
)

// See function [github.com/bytedance/gg/internal/iter.Range].
func Range[T constraints.Number](start, stop T) Number[T] {
	return FromNumberIter(iter.Range(start, stop))
}

// See function [github.com/bytedance/gg/internal/iter.RangeWithStep].
func RangeWithStep[T constraints.Number](start, stop, step T) Number[T] {
	return FromNumberIter(iter.RangeWithStep(start, stop, step))
}

// See function [github.com/bytedance/gg/internal/iter.Sum].
func (s Number[T]) Sum() T {
	return iter.Sum(s.Iter)
}

// See function [github.com/bytedance/gg/internal/iter.Avg].
func (s Number[T]) Avg() float64 {
	return iter.Avg(s.Iter)
}
