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

// This file contains the additional operations of the concurrent-safe
// containers.
//
// Before Go 1.27, a method can not declare its own type parameters, so the
// operations that were not tied to the type parameters of the container were
// still missing as methods.
package skipset

import "github.com/bytedance/gg/goption"

// Any returns whether any element of the set satisfies f.
//
// 💡 NOTE: The set is traversed in ascending order.
func (s *FuncSet[T]) Any(f func(T) bool) bool {
	found := false
	s.Range(func(v T) bool {
		if f(v) {
			found = true
			return false
		}
		return true
	})
	return found
}

// All returns whether all elements of the set satisfy f.
//
// 💡 NOTE: All returns true for an empty set.
func (s *FuncSet[T]) All(f func(T) bool) bool {
	ok := true
	s.Range(func(v T) bool {
		if !f(v) {
			ok = false
			return false
		}
		return true
	})
	return ok
}

// Find returns the first element of the set that satisfies f.
//
// 💡 NOTE: The set is traversed in ascending order, so the first element is
// the smallest one that satisfies f.
func (s *FuncSet[T]) Find(f func(T) bool) goption.O[T] {
	var (
		found T
		ok    bool
	)
	s.Range(func(v T) bool {
		if f(v) {
			found, ok = v, true
			return false
		}
		return true
	})
	return goption.Of(found, ok)
}

// ForEach calls function f for each element of the set.
func (s *FuncSet[T]) ForEach(f func(T)) {
	s.Range(func(v T) bool {
		f(v)
		return true
	})
}

// Reduce reduces the set to a single value by function f.
// It returns nothing when the set is empty.
func (s *FuncSet[T]) Reduce(f func(T, T) T) goption.O[T] {
	var (
		acc  T
		init bool
	)
	s.Range(func(v T) bool {
		if !init {
			acc, init = v, true
			return true
		}
		acc = f(acc, v)
		return true
	})
	if !init {
		return goption.Nil[T]()
	}
	return goption.OK(acc)
}

// Any returns whether any element of the set satisfies f.
//
// 💡 NOTE: The set is traversed in ascending order.
func (s *OrderedSet[T]) Any(f func(T) bool) bool {
	found := false
	s.Range(func(v T) bool {
		if f(v) {
			found = true
			return false
		}
		return true
	})
	return found
}

// All returns whether all elements of the set satisfy f.
//
// 💡 NOTE: All returns true for an empty set.
func (s *OrderedSet[T]) All(f func(T) bool) bool {
	ok := true
	s.Range(func(v T) bool {
		if !f(v) {
			ok = false
			return false
		}
		return true
	})
	return ok
}

// Find returns the first element of the set that satisfies f.
//
// 💡 NOTE: The set is traversed in ascending order, so the first element is
// the smallest one that satisfies f.
func (s *OrderedSet[T]) Find(f func(T) bool) goption.O[T] {
	var (
		found T
		ok    bool
	)
	s.Range(func(v T) bool {
		if f(v) {
			found, ok = v, true
			return false
		}
		return true
	})
	return goption.Of(found, ok)
}

// ForEach calls function f for each element of the set.
func (s *OrderedSet[T]) ForEach(f func(T)) {
	s.Range(func(v T) bool {
		f(v)
		return true
	})
}

// Reduce reduces the set to a single value by function f.
// It returns nothing when the set is empty.
func (s *OrderedSet[T]) Reduce(f func(T, T) T) goption.O[T] {
	var (
		acc  T
		init bool
	)
	s.Range(func(v T) bool {
		if !init {
			acc, init = v, true
			return true
		}
		acc = f(acc, v)
		return true
	})
	if !init {
		return goption.Nil[T]()
	}
	return goption.OK(acc)
}

// Any returns whether any element of the set satisfies f.
//
// 💡 NOTE: The set is traversed in ascending order.
func (s *OrderedSetDesc[T]) Any(f func(T) bool) bool {
	found := false
	s.Range(func(v T) bool {
		if f(v) {
			found = true
			return false
		}
		return true
	})
	return found
}

// All returns whether all elements of the set satisfy f.
//
// 💡 NOTE: All returns true for an empty set.
func (s *OrderedSetDesc[T]) All(f func(T) bool) bool {
	ok := true
	s.Range(func(v T) bool {
		if !f(v) {
			ok = false
			return false
		}
		return true
	})
	return ok
}

// Find returns the first element of the set that satisfies f.
//
// 💡 NOTE: The set is traversed in ascending order, so the first element is
// the smallest one that satisfies f.
func (s *OrderedSetDesc[T]) Find(f func(T) bool) goption.O[T] {
	var (
		found T
		ok    bool
	)
	s.Range(func(v T) bool {
		if f(v) {
			found, ok = v, true
			return false
		}
		return true
	})
	return goption.Of(found, ok)
}

// ForEach calls function f for each element of the set.
func (s *OrderedSetDesc[T]) ForEach(f func(T)) {
	s.Range(func(v T) bool {
		f(v)
		return true
	})
}

// Reduce reduces the set to a single value by function f.
// It returns nothing when the set is empty.
func (s *OrderedSetDesc[T]) Reduce(f func(T, T) T) goption.O[T] {
	var (
		acc  T
		init bool
	)
	s.Range(func(v T) bool {
		if !init {
			acc, init = v, true
			return true
		}
		acc = f(acc, v)
		return true
	})
	if !init {
		return goption.Nil[T]()
	}
	return goption.OK(acc)
}
