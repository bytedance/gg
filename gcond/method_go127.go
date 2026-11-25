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

// This file contains the method form of the [Lazy] operations.
//
// Before Go 1.27, lazy values could only be evaluated by the package-level
// functions ([IfLazy], [IfLazyL], [IfLazyR]). With Go 1.27 generic methods, a
// [Lazy] can be transformed ([Lazy.Map], [Lazy.Then]) before it is evaluated.
package gcond

// Get evaluates the lazy value.
//
// 🚀 EXAMPLE:
//
//	l := Lazy[int](func() int { return 1 })
//	l.Get()  ⏩ 1
func (l Lazy[T]) Get() T {
	return l()
}

// Map returns a new lazy value that applies function f to the result of l.
//
// The function f is not called until the returned lazy value is evaluated.
//
// 🚀 EXAMPLE:
//
//	Lazy[int](func() int { return 1 }).Map(strconv.Itoa).Get()  ⏩ "1"
func (l Lazy[T]) Map[R any](f func(T) R) Lazy[R] {
	return func() R {
		return f(l())
	}
}

// Then is a variant of [Lazy.Map], the function f returns a [Lazy] instead of a
// value.
//
// 🚀 EXAMPLE:
//
//	l := Lazy[int](func() int { return 1 })
//	l.Then(func(i int) Lazy[string] { return func() string { return strconv.Itoa(i) } }).Get()  ⏩ "1"
func (l Lazy[T]) Then[R any](f func(T) Lazy[R]) Lazy[R] {
	return func() R {
		return f(l()).Get()
	}
}

// OrElse returns lazy value l when cond is true, otherwise returns the
// alternative lazy value.
//
// Neither l nor alternative is evaluated until the returned lazy value is
// evaluated.
//
// 🚀 EXAMPLE:
//
//	l := Lazy[int](func() int { return 1 })
//	l.OrElse(true, func() int { return 2 }).Get()   ⏩ 1
//	l.OrElse(false, func() int { return 2 }).Get()  ⏩ 2
func (l Lazy[T]) OrElse(cond bool, alternative Lazy[T]) Lazy[T] {
	if cond {
		return l
	}
	return alternative
}
