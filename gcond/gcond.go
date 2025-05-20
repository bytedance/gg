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

// Package gcond helps users choose values according to various conditions in one line.
package gcond

import (
	"github.com/bytedance/gg/gvalue"
)

// If returns onTrue when cond is true, otherwise returns onFalse.
// It is used as a replacement of ternary conditional operator (:?) in many other
// programming languages.
//
// ⚠️ WARNING: onTrue and onFalse always be evaluated regardless of the truth
// of cond. Use [IfLazy], [IfLazyL], and [IfLazyR] if you need lazy evaluation.
//
// 🚀 EXAMPLE:
//
//	If(true, 1, 2)                       ⏩ 1
//	If(false, 1, 2)                      ⏩ 2
//	If(p != nil, p.foo, nil)             ⏩ ❌PANIC❌
//	If(true, 1, default())               ⏩ 1 // ⚠️ but func default is always evaluated
func If[T any](cond bool, onTrue, onFalse T) T {
	if cond {
		return onTrue
	} else {
		return onFalse
	}
}

// Lazy is a value type that evaluates only when needed.
type Lazy[T any] func() T

// IfLazy is a variant of [If], accepts [Lazy] values.
//
// 🚀 EXAMPLE:
//
//	v1 := func() int {return 1}
//	v2 := func() int {return 2}
//	vp := func () int { panic("") }
//	IfLazy(true, v1, v2)   ⏩ 1
//	IfLazy(false, v1, v2)  ⏩ 2
//	IfLazy(true, v1, vp)   ⏩ 1 // won't panic
//	IfLazy(false, vp, v2)  ⏩ 2 // won't panic
func IfLazy[T any](cond bool, onTrue, onFalse Lazy[T]) T {
	if cond {
		return onTrue()
	} else {
		return onFalse()
	}
}

// IfLazyL is a variant of [If], accepts [Lazy] onTrue value.
func IfLazyL[T any](cond bool, onTrue Lazy[T], onFalse T) T {
	if cond {
		return onTrue()
	} else {
		return onFalse
	}
}

// IfLazyR is a variant of [If], accepts [Lazy] onFalse value.
func IfLazyR[T any](cond bool, onTrue T, onFalse Lazy[T]) T {
	if cond {
		return onTrue
	} else {
		return onFalse()
	}
}

type switchBuilder[R any, T comparable] struct {
	variable T
	matched  bool
	result   R
}

type whenClause[R any, T comparable] struct {
	parent  *switchBuilder[R, T]
	matched bool
}

// Switch initiates a new switchBuilder with the given variable.
// It is used as a more flexible alternative to the built-in switch statement.
//
// 🚀 EXAMPLE:
//
//	Switch[string](1).Case(1, "One").Default("Other")	⏩ "One"
//	Switch[string](2).Case(1, "One").Default("Other")	⏩ "Other"
func Switch[R any, T comparable](variable T) *switchBuilder[R, T] {
	return &switchBuilder[R, T]{
		variable: variable,
		matched:  false,
		result:   gvalue.Zero[R](),
	}
}

// Case adds a case to the switch statement. If the case matches and no previous
// case has matched, it sets the result.
//
// 🚀 EXAMPLE:
//
//	Switch[string](1).Case(1, "One").Default("Other")	⏩ "One"
//	Switch[string](2).Case(1, "One").Default("Other")	⏩ "Other"
func (s *switchBuilder[R, T]) Case(value T, result R) *switchBuilder[R, T] {
	if !s.matched && s.variable == value {
		s.matched = true
		s.result = result
	}
	return s
}

// CaseLazy is a variant of Case that accepts a Lazy result.
// The result function is only called if the case matches.
//
// 🚀 EXAMPLE:
//
//	Switch[string](1).CaseLazy(1, func() string { return "One" }).Default("Other")	⏩ "One"
//	Switch[string](2).CaseLazy(1, func() string { return "One" }).Default("Other")	⏩ "Other"
func (s *switchBuilder[R, T]) CaseLazy(value T, resultF Lazy[R]) *switchBuilder[R, T] {
	if !s.matched && s.variable == value {
		s.matched = true
		s.result = resultF()
	}
	return s
}

// When initiates a multi-value case statement. It returns a whenClause
// which must be followed by a Then or ThenLazy call.
//
// ⚠️ WARNING: When must be followed by a Then or ThenLazy call, otherwise
// the behavior is undefined.
//
// 🚀 EXAMPLE:
//
//	Switch[string](1).When(1, 2).Then("One").Default("Other")	⏩ "One"
//	Switch[string](2).When(1, 2).Then("One").Default("Other")	⏩ "One"
//	Switch[string](3).When(1, 2).Then("One").Default("Other")	⏩ "Other"
func (s *switchBuilder[R, T]) When(values ...T) *whenClause[R, T] {
	wc := &whenClause[R, T]{
		parent:  s,
		matched: false,
	}
	if !s.matched {
		for _, value := range values {
			if s.variable == value {
				wc.matched = true
				break
			}
		}
	}
	return wc
}

// Then sets the result for a When clause if it matches and no previous
// case has matched.
//
// 🚀 EXAMPLE:
//
//	Switch[string](1).When(1, 2).Then("One").Default("Other")	⏩ "One"
//	Switch[string](2).When(1, 2).Then("One").Default("Other")	⏩ "One"
//	Switch[string](3).When(1, 2).Then("One").Default("Other")	⏩ "Other"
func (wc *whenClause[R, T]) Then(result R) *switchBuilder[R, T] {
	if !wc.parent.matched && wc.matched {
		wc.parent.matched = true
		wc.parent.result = result
	}
	return wc.parent
}

// ThenLazy is a variant of Then that accepts a lazy result function.
// The function is only called if the When clause matches and no previous
// case has matched.
//
// 🚀 EXAMPLE:
//
//	Switch[string](1).When(1, 2).ThenLazy(func() string { return "One" }).Default("Other")	⏩ "One"
//	Switch[string](2).When(1, 2).ThenLazy(func() string { return "One" }).Default("Other")	⏩ "One"
//	Switch[string](3).When(1, 2).ThenLazy(func() string { return "One" }).Default("Other")	⏩ "Other"
func (wc *whenClause[R, T]) ThenLazy(resultF func() R) *switchBuilder[R, T] {
	if !wc.parent.matched && wc.matched {
		wc.parent.matched = true
		wc.parent.result = resultF()
	}
	return wc.parent
}

// Default sets the default result and returns the final result of the switch statement.
// It should be called at the end of the switch chain.
// The function is only called if no previous case has matched.
//
// 🚀 EXAMPLE:
//
//	Switch[string](1).Default("Other")	⏩ "Other"
func (s *switchBuilder[R, T]) Default(result R) R {
	if !s.matched {
		s.result = result
	}
	return s.result
}

// DefaultLazy is a variant of Default that accepts a lazy result function.
// It should be called at the end of the switch chain.
// The function is only called if no previous case has matched.
//
// 🚀 EXAMPLE:
//
//	Switch[string](1).DefaultLazy(func() string{ return "Other" })	⏩ "Other"
func (s *switchBuilder[R, T]) DefaultLazy(resultF Lazy[R]) R {
	if !s.matched {
		s.result = resultF()
	}
	return s.result
}
