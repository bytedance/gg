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

// This file contains the additional operations of [List].
//
// Before Go 1.27, a method can not declare its own type parameters, so the
// operations that change the element type (such as [List.Map]) were not
// available as methods.
package list

import "github.com/bytedance/gg/goption"

// ToSlice returns the values of the list as a slice.
//
// 🚀 EXAMPLE:
//
//	l := New[int]()
//	l.PushBack(1)
//	l.PushBack(2)
//	l.ToSlice()  ⏩ [1, 2]
func (l *List[T]) ToSlice() []T {
	s := make([]T, 0, l.Len())
	for e := l.Front(); e != nil; e = e.Next() {
		s = append(s, e.Value)
	}
	return s
}

// Map applies function f to each value of the list and returns a new list with
// the results.
//
// 🚀 EXAMPLE:
//
//	l := New[int]()
//	l.PushBack(1)
//	l.Map(strconv.Itoa).ToSlice()  ⏩ ["1"]
func (l *List[T]) Map[R any](f func(T) R) *List[R] {
	nl := New[R]()
	for e := l.Front(); e != nil; e = e.Next() {
		nl.PushBack(f(e.Value))
	}
	return nl
}

// Filter returns a new list with the values that f returns true for.
//
// 🚀 EXAMPLE:
//
//	l := New[int]()
//	l.PushBack(1)
//	l.PushBack(2)
//	l.Filter(func(i int) bool { return i%2 == 0 }).ToSlice()  ⏩ [2]
func (l *List[T]) Filter(f func(T) bool) *List[T] {
	nl := New[T]()
	for e := l.Front(); e != nil; e = e.Next() {
		if f(e.Value) {
			nl.PushBack(e.Value)
		}
	}
	return nl
}

// Any returns whether any value of the list satisfies f.
func (l *List[T]) Any(f func(T) bool) bool {
	for e := l.Front(); e != nil; e = e.Next() {
		if f(e.Value) {
			return true
		}
	}
	return false
}

// All returns whether all values of the list satisfy f.
func (l *List[T]) All(f func(T) bool) bool {
	for e := l.Front(); e != nil; e = e.Next() {
		if !f(e.Value) {
			return false
		}
	}
	return true
}

// Find returns the first value of the list that satisfies f.
func (l *List[T]) Find(f func(T) bool) goption.O[T] {
	for e := l.Front(); e != nil; e = e.Next() {
		if f(e.Value) {
			return goption.OK(e.Value)
		}
	}
	return goption.Nil[T]()
}

// ForEach calls function f for each value of the list.
func (l *List[T]) ForEach(f func(T)) {
	for e := l.Front(); e != nil; e = e.Next() {
		f(e.Value)
	}
}

// Reduce reduces the list to a single value by function f.
// It returns nothing when the list is empty.
func (l *List[T]) Reduce(f func(T, T) T) goption.O[T] {
	e := l.Front()
	if e == nil {
		return goption.Nil[T]()
	}
	acc := e.Value
	for e = e.Next(); e != nil; e = e.Next() {
		acc = f(acc, e.Value)
	}
	return goption.OK(acc)
}

// Partition splits the list into two lists: the values that f returns true for,
// and the rest.
func (l *List[T]) Partition(f func(T) bool) (*List[T], *List[T]) {
	matched, unmatched := New[T](), New[T]()
	for e := l.Front(); e != nil; e = e.Next() {
		if f(e.Value) {
			matched.PushBack(e.Value)
		} else {
			unmatched.PushBack(e.Value)
		}
	}
	return matched, unmatched
}
