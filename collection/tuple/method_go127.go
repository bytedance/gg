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

// This file contains the transformations of the tuples.
//
// Transforming a tuple into a tuple of other types (such as [T2.Map]) changes
// the element types, which requires method-level type parameters and therefore
// is only possible since Go 1.27 generic methods.
package tuple

// Map applies function f to all elements of the tuple and returns a new
// tuple with the results.
//
// 🚀 EXAMPLE:
//
//	Make2(1, "a").Map(func(i int, s string) (string, int) { return s, i })
//	// T2[string, int]{First: "a", Second: 1}
func (t T2[V1, V2]) Map[R1, R2 any](f func(V1, V2) (R1, R2)) T2[R1, R2] {
	return Make2(f(t.First, t.Second))
}

// Map applies function f to all elements of the tuple and returns a new
// tuple with the results.
func (t T3[V1, V2, V3]) Map[R1, R2, R3 any](f func(V1, V2, V3) (R1, R2, R3)) T3[R1, R2, R3] {
	return Make3(f(t.First, t.Second, t.Third))
}

// Map applies function f to all elements of the tuple and returns a new
// tuple with the results.
func (t T4[V1, V2, V3, V4]) Map[R1, R2, R3, R4 any](f func(V1, V2, V3, V4) (R1, R2, R3, R4)) T4[R1, R2, R3, R4] {
	return Make4(f(t.First, t.Second, t.Third, t.Fourth))
}

// Map applies function f to all elements of the tuple and returns a new
// tuple with the results.
func (t T5[V1, V2, V3, V4, V5]) Map[R1, R2, R3, R4, R5 any](f func(V1, V2, V3, V4, V5) (R1, R2, R3, R4, R5)) T5[R1, R2, R3, R4, R5] {
	return Make5(f(t.First, t.Second, t.Third, t.Fourth, t.Fifth))
}

// Map applies function f to all elements of the tuple and returns a new
// tuple with the results.
func (t T6[V1, V2, V3, V4, V5, V6]) Map[R1, R2, R3, R4, R5, R6 any](f func(V1, V2, V3, V4, V5, V6) (R1, R2, R3, R4, R5, R6)) T6[R1, R2, R3, R4, R5, R6] {
	return Make6(f(t.First, t.Second, t.Third, t.Fourth, t.Fifth, t.Sixth))
}

// Map applies function f to all elements of the tuple and returns a new
// tuple with the results.
func (t T7[V1, V2, V3, V4, V5, V6, V7]) Map[R1, R2, R3, R4, R5, R6, R7 any](f func(V1, V2, V3, V4, V5, V6, V7) (R1, R2, R3, R4, R5, R6, R7)) T7[R1, R2, R3, R4, R5, R6, R7] {
	return Make7(f(t.First, t.Second, t.Third, t.Fourth, t.Fifth, t.Sixth, t.Seventh))
}

// Map applies function f to all elements of the tuple and returns a new
// tuple with the results.
func (t T8[V1, V2, V3, V4, V5, V6, V7, V8]) Map[R1, R2, R3, R4, R5, R6, R7, R8 any](f func(V1, V2, V3, V4, V5, V6, V7, V8) (R1, R2, R3, R4, R5, R6, R7, R8)) T8[R1, R2, R3, R4, R5, R6, R7, R8] {
	return Make8(f(t.First, t.Second, t.Third, t.Fourth, t.Fifth, t.Sixth, t.Seventh, t.Eighth))
}

// Map applies function f to all elements of the tuple and returns a new
// tuple with the results.
func (t T9[V1, V2, V3, V4, V5, V6, V7, V8, V9]) Map[R1, R2, R3, R4, R5, R6, R7, R8, R9 any](f func(V1, V2, V3, V4, V5, V6, V7, V8, V9) (R1, R2, R3, R4, R5, R6, R7, R8, R9)) T9[R1, R2, R3, R4, R5, R6, R7, R8, R9] {
	return Make9(f(t.First, t.Second, t.Third, t.Fourth, t.Fifth, t.Sixth, t.Seventh, t.Eighth, t.Ninth))
}

// Map applies function f to all elements of the tuple and returns a new
// tuple with the results.
func (t T10[V1, V2, V3, V4, V5, V6, V7, V8, V9, V10]) Map[R1, R2, R3, R4, R5, R6, R7, R8, R9, R10 any](f func(V1, V2, V3, V4, V5, V6, V7, V8, V9, V10) (R1, R2, R3, R4, R5, R6, R7, R8, R9, R10)) T10[R1, R2, R3, R4, R5, R6, R7, R8, R9, R10] {
	return Make10(f(t.First, t.Second, t.Third, t.Fourth, t.Fifth, t.Sixth, t.Seventh, t.Eighth, t.Ninth, t.Tenth))
}

// Swap returns a new tuple with the elements swapped.
//
// 🚀 EXAMPLE:
//
//	Make2(1, "a").Swap()  ⏩ T2[string, int]{First: "a", Second: 1}
func (t T2[V1, V2]) Swap() T2[V2, V1] {
	return Make2(t.Second, t.First)
}

// Map applies function f to all elements of the pair and returns a new pair
// with the results.
func (p Pair[V1, V2]) Map[R1, R2 any](f func(V1, V2) (R1, R2)) Pair[R1, R2] {
	return Pair[R1, R2](T2[V1, V2](p).Map(f))
}

// Swap returns a new pair with the elements swapped.
func (p Pair[V1, V2]) Swap() Pair[V2, V1] {
	return Pair[V2, V1](T2[V1, V2](p).Swap())
}

// Map applies function f to every tuple of the slice and returns a new
// slice with the results.
func (s S2[V1, V2]) Map[R any](f func(T2[V1, V2]) R) []R {
	r := make([]R, 0, len(s))
	for _, v := range s {
		r = append(r, f(v))
	}
	return r
}

// Filter returns a new slice of tuples with the tuples that f returns true
// for.
func (s S2[V1, V2]) Filter(f func(T2[V1, V2]) bool) S2[V1, V2] {
	r := make(S2[V1, V2], 0, len(s))
	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}
	return r
}

// Map applies function f to every tuple of the slice and returns a new
// slice with the results.
func (s S3[V1, V2, V3]) Map[R any](f func(T3[V1, V2, V3]) R) []R {
	r := make([]R, 0, len(s))
	for _, v := range s {
		r = append(r, f(v))
	}
	return r
}

// Filter returns a new slice of tuples with the tuples that f returns true
// for.
func (s S3[V1, V2, V3]) Filter(f func(T3[V1, V2, V3]) bool) S3[V1, V2, V3] {
	r := make(S3[V1, V2, V3], 0, len(s))
	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}
	return r
}

// Map applies function f to every tuple of the slice and returns a new
// slice with the results.
func (s S4[V1, V2, V3, V4]) Map[R any](f func(T4[V1, V2, V3, V4]) R) []R {
	r := make([]R, 0, len(s))
	for _, v := range s {
		r = append(r, f(v))
	}
	return r
}

// Filter returns a new slice of tuples with the tuples that f returns true
// for.
func (s S4[V1, V2, V3, V4]) Filter(f func(T4[V1, V2, V3, V4]) bool) S4[V1, V2, V3, V4] {
	r := make(S4[V1, V2, V3, V4], 0, len(s))
	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}
	return r
}

// Map applies function f to every tuple of the slice and returns a new
// slice with the results.
func (s S5[V1, V2, V3, V4, V5]) Map[R any](f func(T5[V1, V2, V3, V4, V5]) R) []R {
	r := make([]R, 0, len(s))
	for _, v := range s {
		r = append(r, f(v))
	}
	return r
}

// Filter returns a new slice of tuples with the tuples that f returns true
// for.
func (s S5[V1, V2, V3, V4, V5]) Filter(f func(T5[V1, V2, V3, V4, V5]) bool) S5[V1, V2, V3, V4, V5] {
	r := make(S5[V1, V2, V3, V4, V5], 0, len(s))
	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}
	return r
}

// Map applies function f to every tuple of the slice and returns a new
// slice with the results.
func (s S6[V1, V2, V3, V4, V5, V6]) Map[R any](f func(T6[V1, V2, V3, V4, V5, V6]) R) []R {
	r := make([]R, 0, len(s))
	for _, v := range s {
		r = append(r, f(v))
	}
	return r
}

// Filter returns a new slice of tuples with the tuples that f returns true
// for.
func (s S6[V1, V2, V3, V4, V5, V6]) Filter(f func(T6[V1, V2, V3, V4, V5, V6]) bool) S6[V1, V2, V3, V4, V5, V6] {
	r := make(S6[V1, V2, V3, V4, V5, V6], 0, len(s))
	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}
	return r
}

// Map applies function f to every tuple of the slice and returns a new
// slice with the results.
func (s S7[V1, V2, V3, V4, V5, V6, V7]) Map[R any](f func(T7[V1, V2, V3, V4, V5, V6, V7]) R) []R {
	r := make([]R, 0, len(s))
	for _, v := range s {
		r = append(r, f(v))
	}
	return r
}

// Filter returns a new slice of tuples with the tuples that f returns true
// for.
func (s S7[V1, V2, V3, V4, V5, V6, V7]) Filter(f func(T7[V1, V2, V3, V4, V5, V6, V7]) bool) S7[V1, V2, V3, V4, V5, V6, V7] {
	r := make(S7[V1, V2, V3, V4, V5, V6, V7], 0, len(s))
	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}
	return r
}

// Map applies function f to every tuple of the slice and returns a new
// slice with the results.
func (s S8[V1, V2, V3, V4, V5, V6, V7, V8]) Map[R any](f func(T8[V1, V2, V3, V4, V5, V6, V7, V8]) R) []R {
	r := make([]R, 0, len(s))
	for _, v := range s {
		r = append(r, f(v))
	}
	return r
}

// Filter returns a new slice of tuples with the tuples that f returns true
// for.
func (s S8[V1, V2, V3, V4, V5, V6, V7, V8]) Filter(f func(T8[V1, V2, V3, V4, V5, V6, V7, V8]) bool) S8[V1, V2, V3, V4, V5, V6, V7, V8] {
	r := make(S8[V1, V2, V3, V4, V5, V6, V7, V8], 0, len(s))
	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}
	return r
}

// Map applies function f to every tuple of the slice and returns a new
// slice with the results.
func (s S9[V1, V2, V3, V4, V5, V6, V7, V8, V9]) Map[R any](f func(T9[V1, V2, V3, V4, V5, V6, V7, V8, V9]) R) []R {
	r := make([]R, 0, len(s))
	for _, v := range s {
		r = append(r, f(v))
	}
	return r
}

// Filter returns a new slice of tuples with the tuples that f returns true
// for.
func (s S9[V1, V2, V3, V4, V5, V6, V7, V8, V9]) Filter(f func(T9[V1, V2, V3, V4, V5, V6, V7, V8, V9]) bool) S9[V1, V2, V3, V4, V5, V6, V7, V8, V9] {
	r := make(S9[V1, V2, V3, V4, V5, V6, V7, V8, V9], 0, len(s))
	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}
	return r
}

// Map applies function f to every tuple of the slice and returns a new
// slice with the results.
func (s S10[V1, V2, V3, V4, V5, V6, V7, V8, V9, V10]) Map[R any](f func(T10[V1, V2, V3, V4, V5, V6, V7, V8, V9, V10]) R) []R {
	r := make([]R, 0, len(s))
	for _, v := range s {
		r = append(r, f(v))
	}
	return r
}

// Filter returns a new slice of tuples with the tuples that f returns true
// for.
func (s S10[V1, V2, V3, V4, V5, V6, V7, V8, V9, V10]) Filter(f func(T10[V1, V2, V3, V4, V5, V6, V7, V8, V9, V10]) bool) S10[V1, V2, V3, V4, V5, V6, V7, V8, V9, V10] {
	r := make(S10[V1, V2, V3, V4, V5, V6, V7, V8, V9, V10], 0, len(s))
	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}
	return r
}
