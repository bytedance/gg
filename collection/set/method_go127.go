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

//go:build go1.27

// This file contains the additional operations of [Set].
//
// A set is a container of comparable elements, the operations that change the
// element type (such as [Set.Map] and [Set.FlatMap]) need a method-level type
// parameter and therefore are only possible since Go 1.27 generic methods.
package set

import "github.com/bytedance/gg/goption"

// Map applies function f to each member of the set and returns a new set with
// the results.
//
// 💡 NOTE: The result may have fewer members than the original set,
// since the results of f may collide.
//
// 🚀 EXAMPLE:
//
//	New(1, 2).Map(strconv.Itoa).ToSlice()  ⏩ ["1", "2"]
func (s *Set[T]) Map[R comparable](f func(T) R) *Set[R] {
	ns := NewWithCap[R](s.Len())
	if s == nil {
		return ns
	}
	for v := range s.m {
		ns.Add(f(v))
	}
	return ns
}

// FlatMap applies function f to each member of the set and flattens the results
// into a new set.
//
// 🚀 EXAMPLE:
//
//	New(1, 2).FlatMap(func(i int) []string { return []string{strconv.Itoa(i), strconv.Itoa(-i)} }).Len()
//	// 4
func (s *Set[T]) FlatMap[R comparable](f func(T) []R) *Set[R] {
	ns := NewWithCap[R](s.Len())
	if s == nil {
		return ns
	}
	for v := range s.m {
		ns.AddN(f(v)...)
	}
	return ns
}

// Filter returns a new set with the members that f returns true for.
//
// 🚀 EXAMPLE:
//
//	New(1, 2, 3).Filter(func(i int) bool { return i%2 == 1 }).Len()  ⏩ 2
func (s *Set[T]) Filter(f func(T) bool) *Set[T] {
	ns := NewWithCap[T](s.Len())
	if s == nil {
		return ns
	}
	for v := range s.m {
		if f(v) {
			ns.Add(v)
		}
	}
	return ns
}

// Partition splits the set into two sets: the members that f returns true for,
// and the rest.
func (s *Set[T]) Partition(f func(T) bool) (*Set[T], *Set[T]) {
	matched, unmatched := NewWithCap[T](s.Len()), NewWithCap[T](s.Len())
	if s == nil {
		return matched, unmatched
	}
	for v := range s.m {
		if f(v) {
			matched.Add(v)
		} else {
			unmatched.Add(v)
		}
	}
	return matched, unmatched
}

// Any returns whether any member of the set satisfies f.
func (s *Set[T]) Any(f func(T) bool) bool {
	if s == nil {
		return false
	}
	for v := range s.m {
		if f(v) {
			return true
		}
	}
	return false
}

// All returns whether all members of the set satisfy f.
//
// 💡 NOTE: All returns true for an empty set.
func (s *Set[T]) All(f func(T) bool) bool {
	if s == nil {
		return true
	}
	for v := range s.m {
		if !f(v) {
			return false
		}
	}
	return true
}

// Find returns the first member of the set that satisfies f.
//
// 💡 HINT: The iteration order of a set is undefined, so the result is not
// deterministic when more than one member satisfies f.
func (s *Set[T]) Find(f func(T) bool) goption.O[T] {
	if s == nil {
		return goption.Nil[T]()
	}
	for v := range s.m {
		if f(v) {
			return goption.OK(v)
		}
	}
	return goption.Nil[T]()
}

// ForEach calls function f for each member of the set.
func (s *Set[T]) ForEach(f func(T)) {
	if s == nil {
		return
	}
	for v := range s.m {
		f(v)
	}
}

// Reduce reduces the set to a single value by function f.
// It returns nothing when the set is empty.
//
// 💡 HINT: The iteration order of a set is undefined, so f should be
// commutative and associative, such as the addition.
func (s *Set[T]) Reduce(f func(T, T) T) goption.O[T] {
	if s == nil || len(s.m) == 0 {
		return goption.Nil[T]()
	}
	var (
		acc  T
		init bool
	)
	for v := range s.m {
		if !init {
			acc, init = v, true
			continue
		}
		acc = f(acc, v)
	}
	return goption.OK(acc)
}

// MaxBy returns the maximum member of the set, compared by function less.
func (s *Set[T]) MaxBy(less func(T, T) bool) goption.O[T] {
	if s == nil || len(s.m) == 0 {
		return goption.Nil[T]()
	}
	var (
		max  T
		init bool
	)
	for v := range s.m {
		if !init || less(max, v) {
			max, init = v, true
		}
	}
	return goption.OK(max)
}

// MinBy returns the minimum member of the set, compared by function less.
func (s *Set[T]) MinBy(less func(T, T) bool) goption.O[T] {
	if s == nil || len(s.m) == 0 {
		return goption.Nil[T]()
	}
	var (
		min  T
		init bool
	)
	for v := range s.m {
		if !init || less(v, min) {
			min, init = v, true
		}
	}
	return goption.OK(min)
}

// GroupBy groups the members into a map of sets by the key returned by f.
//
// 🚀 EXAMPLE:
//
//	New(1, 2, 3).GroupBy(func(i int) string { return strconv.Itoa(i % 2) })
//	// {"0": {2}, "1": {1, 3}}
func (s *Set[T]) GroupBy[K comparable](f func(T) K) map[K]*Set[T] {
	m := make(map[K]*Set[T])
	if s == nil {
		return m
	}
	for v := range s.m {
		k := f(v)
		if _, ok := m[k]; !ok {
			m[k] = New[T]()
		}
		m[k].Add(v)
	}
	return m
}
