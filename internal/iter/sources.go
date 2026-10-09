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

package iter

import (
	"context"
	"math"

	"github.com/bytedance/gg/collection/tuple"
	"github.com/bytedance/gg/internal/constraints"
)

var (
	_ Iter[int]                   = &sliceIter[int]{}
	_ Iter[tuple.T2[string, int]] = &mapKeyValueIter[string, int]{}
	_ Iter[string]                = &mapKeyIter[string, int]{}
	_ Iter[int]                   = &mapValueIter[string, int]{}
	_ Iter[int]                   = &chanIter[int]{}
	_ Iter[int]                   = &rangeIter[int]{}
	_ Iter[int]                   = &repeatIter[int]{}
)

type sliceIter[T any] struct {
	s []T
}

func (i *sliceIter[T]) Next(n int) []T {
	if n == ALL || n > len(i.s) {
		n = len(i.s)
	}
	next := make([]T, n)
	copy(next, i.s[:n])
	i.s = i.s[n:]
	return next
}

// FromSlice constructs an Iter from slice s, in order from left to right.
// An empty Iter (without element) is returned if the given slice is empty
// or nil.
func FromSlice[T any](s []T) Iter[T] {
	return &sliceIter[T]{s}
}

type stealSliceIter[T any] struct {
	s []T
}

func (i *stealSliceIter[T]) Next(n int) []T {
	if n == ALL || n > len(i.s) {
		n = len(i.s)
	}
	next := i.s[:n]
	i.s = i.s[n:]
	return next
}

func StealSlice[T any](s []T) Iter[T] {
	return &stealSliceIter[T]{s}
}

type mapKeyValueIter[K comparable, V any] struct {
	i *unsafeMapIter[K, V]
}

func (i *mapKeyValueIter[K, V]) Next(n int) []tuple.T2[K, V] {
	keys, values := i.i.Next(n, true, true)
	// len(keys) must equal to len(value)
	if len(keys) == 0 {
		return nil
	}
	next := make([]tuple.T2[K, V], len(keys))
	for i := 0; i < len(next); i++ {
		next[i] = tuple.Make2(keys[i], values[i])
	}
	return next
}

// FromMap constructs an Iter of (key, value) pair from map m.
//
// 💡 NOTE: Function follows the same iteration semantics as a range statement.
// See https://go.dev/blog/maps#iteration-order for details.
func FromMap[K comparable, V any](m map[K]V) Iter[tuple.T2[K, V]] {
	return &mapKeyValueIter[K, V]{newUnsafeMapIter(m)}
}

type mapKeyIter[K comparable, V any] struct {
	i *unsafeMapIter[K, V]
}

func (i *mapKeyIter[K, V]) Next(n int) []K {
	next, _ := i.i.Next(n, true, false)
	return next
}

// FromMapKeys constructs an Iter of map's key from map m.
//
// 💡 NOTE: Function follows the same iteration semantics as a range statement.
// See https://go.dev/blog/maps#iteration-order for details.
func FromMapKeys[K comparable, V any](m map[K]V) Iter[K] {
	return &mapKeyIter[K, V]{newUnsafeMapIter(m)}
}

type mapValueIter[K comparable, V any] struct {
	i *unsafeMapIter[K, V]
}

func (i *mapValueIter[K, V]) Next(n int) []V {
	_, next := i.i.Next(n, false, true)
	return next
}

// FromMapValues constructs an Iter of map's value from map m.
//
// 💡 NOTE: Function follows the same iteration semantics as a range statement.
// See https://go.dev/blog/maps#iteration-order for details.
func FromMapValues[K comparable, V any](m map[K]V) Iter[V] {
	return &mapValueIter[K, V]{newUnsafeMapIter(m)}
}

type chanIter[T any] struct {
	ctx context.Context
	ch  <-chan T
}

func (i *chanIter[T]) Next(n int) (r []T) {
	if n == 0 {
		return
	}

	if n == ALL {
		for {
			select {
			case <-i.ctx.Done():
				return
			case v, ok := <-i.ch:
				if !ok {
					return
				}
				r = append(r, v)
			}
		}
	}

	r = make([]T, 0, n)
	for j := 0; j < n; j++ {
		select {
		case <-i.ctx.Done():
			return
		case v, ok := <-i.ch:
			if !ok {
				return
			}
			r = append(r, v)
		}
	}
	return
}

// FromChan constructs an Iter from channel ch.
// Elements in Iter are exhausted when the context is done or given channel is closed.
// FIXME: better doc
func FromChan[T any](ctx context.Context, ch <-chan T) Iter[T] {
	return &chanIter[T]{ctx, ch}
}

type rangeIter[T constraints.Number] struct {
	cur       T
	stop      T
	step      T
	exhausted bool
}

func (i *rangeIter[T]) Next(n int) (r []T) {
	if n == 0 || i.exhausted {
		return nil
	}
	// Estimate capacity only after converting the operands to float64. Doing
	// the subtraction in T first is exactly what caused integer ranges that
	// cross zero to overflow. Cap the estimate to avoid a huge speculative
	// allocation for a short or precision-limited range.
	const maxRangePrealloc = 1 << 20
	capacity := int(math.Ceil(math.Abs(float64(i.stop)-float64(i.cur)) / math.Abs(float64(i.step))))
	if capacity < 0 || capacity > maxRangePrealloc {
		capacity = maxRangePrealloc
	}
	if n != ALL && n < capacity {
		capacity = n
	}
	r = make([]T, 0, capacity)
	for n == ALL || len(r) < n {
		if (i.step > 0 && i.cur >= i.stop) || (i.step < 0 && i.cur <= i.stop) {
			i.exhausted = true
			break
		}
		r = append(r, i.cur)
		next := i.cur + i.step
		// The addition must move strictly toward stop. Otherwise an integer
		// overflow or floating-point precision loss would make the iterator
		// wrap around or repeat the same value forever.
		if (i.step > 0 && next <= i.cur) || (i.step < 0 && next >= i.cur) {
			i.exhausted = true
			break
		}
		i.cur = next
	}
	return r
}

// Range is a variant of RangeWithStep, with predefined step 1.
func Range[T constraints.Number](start, stop T) Iter[T] {
	return RangeWithStep(start, stop, 1)
}

// RangeWithStep constructs an Iter of number from start (inclusive) to stop (exclusive)
// by step.
// If the interval does not exist, RangeWithStep returns an emptyIter.
func RangeWithStep[T constraints.Number](start, stop, step T) Iter[T] {
	// A NaN compares unequal to itself. Reject it explicitly; otherwise all
	// direction and termination comparisons would remain false forever.
	if start != start || stop != stop || step != step || step == 0 ||
		(step > 0 && start >= stop) || (step < 0 && start <= stop) {
		return emptyIter[T]{}
	}
	return &rangeIter[T]{cur: start, stop: stop, step: step}
}

type repeatIter[T any] struct {
	v T
}

func (i *repeatIter[T]) Next(n int) (r []T) {
	if n == ALL {
		panic("infinite elements")
	}
	r = make([]T, 0, n)
	for j := 0; j < n; j++ {
		r = append(r, i.v)
	}
	return
}

// Repeat constructs an infinite Iter, with v the value of every element.
func Repeat[T any](v T) Iter[T] {
	return &repeatIter[T]{v}
}
