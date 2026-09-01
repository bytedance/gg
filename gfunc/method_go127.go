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

// This file contains the composition of the function types.
//
// [Func1] to [Func10] already provide partial application as methods
// ([Func2.Partial], [Func2.PartialR], …), but composing a function with another
// one changes the result type, which requires a method-level type parameter and
// therefore is only possible since Go 1.27 generic methods.
package gfunc

// Compose returns a new function that applies f and then g to the result of f.
//
// 🚀 EXAMPLE:
//
//	parse := Func1[string, int](func(s string) int { return len(s) })
//	parse.Compose(strconv.Itoa)("abc")  ⏩ "3"
func (f Func1[T1, R]) Compose[R2 any](g func(R) R2) Func1[T1, R2] {
	return func(t1 T1) R2 {
		return g(f(t1))
	}
}

// Compose returns a new function that applies f and then g to the result of f.
func (f Func2[T1, T2, R]) Compose[R2 any](g func(R) R2) Func2[T1, T2, R2] {
	return func(t1 T1, t2 T2) R2 {
		return g(f(t1, t2))
	}
}

// Compose returns a new function that applies f and then g to the result of f.
func (f Func3[T1, T2, T3, R]) Compose[R2 any](g func(R) R2) Func3[T1, T2, T3, R2] {
	return func(t1 T1, t2 T2, t3 T3) R2 {
		return g(f(t1, t2, t3))
	}
}

// Compose returns a new function that applies f and then g to the result of f.
func (f Func4[T1, T2, T3, T4, R]) Compose[R2 any](g func(R) R2) Func4[T1, T2, T3, T4, R2] {
	return func(t1 T1, t2 T2, t3 T3, t4 T4) R2 {
		return g(f(t1, t2, t3, t4))
	}
}

// Compose returns a new function that applies f and then g to the result of f.
func (f Func5[T1, T2, T3, T4, T5, R]) Compose[R2 any](g func(R) R2) Func5[T1, T2, T3, T4, T5, R2] {
	return func(t1 T1, t2 T2, t3 T3, t4 T4, t5 T5) R2 {
		return g(f(t1, t2, t3, t4, t5))
	}
}

// Compose returns a new function that applies f and then g to the result of f.
func (f Func6[T1, T2, T3, T4, T5, T6, R]) Compose[R2 any](g func(R) R2) Func6[T1, T2, T3, T4, T5, T6, R2] {
	return func(t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6) R2 {
		return g(f(t1, t2, t3, t4, t5, t6))
	}
}

// Compose returns a new function that applies f and then g to the result of f.
func (f Func7[T1, T2, T3, T4, T5, T6, T7, R]) Compose[R2 any](g func(R) R2) Func7[T1, T2, T3, T4, T5, T6, T7, R2] {
	return func(t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7) R2 {
		return g(f(t1, t2, t3, t4, t5, t6, t7))
	}
}

// Compose returns a new function that applies f and then g to the result of f.
func (f Func8[T1, T2, T3, T4, T5, T6, T7, T8, R]) Compose[R2 any](g func(R) R2) Func8[T1, T2, T3, T4, T5, T6, T7, T8, R2] {
	return func(t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8) R2 {
		return g(f(t1, t2, t3, t4, t5, t6, t7, t8))
	}
}

// Compose returns a new function that applies f and then g to the result of f.
func (f Func9[T1, T2, T3, T4, T5, T6, T7, T8, T9, R]) Compose[R2 any](g func(R) R2) Func9[T1, T2, T3, T4, T5, T6, T7, T8, T9, R2] {
	return func(t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8, t9 T9) R2 {
		return g(f(t1, t2, t3, t4, t5, t6, t7, t8, t9))
	}
}

// Compose returns a new function that applies f and then g to the result of f.
func (f Func10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, R]) Compose[R2 any](g func(R) R2) Func10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, R2] {
	return func(t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8, t9 T9, t10 T10) R2 {
		return g(f(t1, t2, t3, t4, t5, t6, t7, t8, t9, t10))
	}
}
