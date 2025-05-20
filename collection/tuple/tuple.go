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

// Package tuple provides definition of generic n-ary tuples, from [T2] to [T10].
//
// # Quick Start
//
// package main
//
//	import (
//		"fmt"
//		"github.com/bytedance/gg/collection/tuple"
//	)
//
//	func main() {
//		addr := tuple.Make2("localhost", 8080)
//		fmt.Printf("%s:%d\n", addr.First, addr.Second)
//		// Output:
//		// localhost:8080
//	 }
//
// If you have a need for n-ary (where n > 10) tuple, please file an issue.
//
// [Type Parameters Proposal]: https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#type-inference-for-composite-literals)
// [Type inference for functions is supported]: https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#type-inference
package tuple

import (
	"github.com/bytedance/gg/gvalue"
)

type Pair[V1, V2 any] T2[V1, V2]

// T2 is a 2-ary tuple.
type T2[V1, V2 any] struct {
	First  V1
	Second V2
}

// Values returns all elements of tuple.
func (t T2[V1, V2]) Values() (V1, V2) {
	return t.First, t.Second
}

// Make2 creates a tuple of 2 elements.
func Make2[V1, V2 any](first V1, second V2) T2[V1, V2] {
	return T2[V1, V2]{first, second}
}

// T3 is a 3-ary tuple.
type T3[V1, V2, V3 any] struct {
	First  V1
	Second V2
	Third  V3
}

// Values returns all elements of tuple.
func (t T3[V1, V2, V3]) Values() (V1, V2, V3) {
	return t.First, t.Second, t.Third
}

// Make3 creates a tuple of 3 elements.
func Make3[V1, V2, V3 any](first V1, second V2, third V3) T3[V1, V2, V3] {
	return T3[V1, V2, V3]{first, second, third}
}

// T4 is a 4-ary tuple.
type T4[V1, V2, V3, V4 any] struct {
	First  V1
	Second V2
	Third  V3
	Fourth V4
}

// Values returns all elements of tuple.
func (t T4[V1, V2, V3, V4]) Values() (V1, V2, V3, V4) {
	return t.First, t.Second, t.Third, t.Fourth
}

// Make4 creates a tuple of 4 elements.
func Make4[V1, V2, V3, V4 any](first V1, second V2, third V3, fourth V4) T4[V1, V2, V3, V4] {
	return T4[V1, V2, V3, V4]{first, second, third, fourth}
}

// T5 is a 5-ary tuple.
type T5[V1, V2, V3, V4, V5 any] struct {
	First  V1
	Second V2
	Third  V3
	Fourth V4
	Fifth  V5
}

// Values returns all elements of tuple.
func (t T5[V1, V2, V3, V4, V5]) Values() (V1, V2, V3, V4, V5) {
	return t.First, t.Second, t.Third, t.Fourth, t.Fifth
}

// Make5 creates a tuple of 5 elements.
func Make5[V1, V2, V3, V4, V5 any](first V1, second V2, third V3, fourth V4, fifth V5) T5[V1, V2, V3, V4, V5] {
	return T5[V1, V2, V3, V4, V5]{first, second, third, fourth, fifth}
}

// T6 is a 6-ary tuple.
type T6[V1, V2, V3, V4, V5, V6 any] struct {
	First  V1
	Second V2
	Third  V3
	Fourth V4
	Fifth  V5
	Sixth  V6
}

// Values returns all elements of tuple.
func (t T6[V1, V2, V3, V4, V5, V6]) Values() (V1, V2, V3, V4, V5, V6) {
	return t.First, t.Second, t.Third, t.Fourth, t.Fifth, t.Sixth
}

// Make6 creates a tuple of 6 elements.
func Make6[V1, V2, V3, V4, V5, V6 any](first V1, second V2, third V3, fourth V4, fifth V5, sixth V6) T6[V1, V2, V3, V4, V5, V6] {
	return T6[V1, V2, V3, V4, V5, V6]{first, second, third, fourth, fifth, sixth}
}

// T7 is a 7-ary tuple.
type T7[V1, V2, V3, V4, V5, V6, V7 any] struct {
	First   V1
	Second  V2
	Third   V3
	Fourth  V4
	Fifth   V5
	Sixth   V6
	Seventh V7
}

// Values returns all elements of tuple.
func (t T7[V1, V2, V3, V4, V5, V6, V7]) Values() (V1, V2, V3, V4, V5, V6, V7) {
	return t.First, t.Second, t.Third, t.Fourth, t.Fifth, t.Sixth, t.Seventh
}

// Make7 creates a tuple of 7 elements.
func Make7[V1, V2, V3, V4, V5, V6, V7 any](first V1, second V2, third V3, fourth V4, fifth V5, sixth V6, seventh V7) T7[V1, V2, V3, V4, V5, V6, V7] {
	return T7[V1, V2, V3, V4, V5, V6, V7]{first, second, third, fourth, fifth, sixth, seventh}
}

// T8 is a 8-ary tuple.
type T8[V1, V2, V3, V4, V5, V6, V7, V8 any] struct {
	First   V1
	Second  V2
	Third   V3
	Fourth  V4
	Fifth   V5
	Sixth   V6
	Seventh V7
	Eighth  V8
}

// Values returns all elements of tuple.
func (t T8[V1, V2, V3, V4, V5, V6, V7, V8]) Values() (V1, V2, V3, V4, V5, V6, V7, V8) {
	return t.First, t.Second, t.Third, t.Fourth, t.Fifth, t.Sixth, t.Seventh, t.Eighth
}

// Make8 creates a tuple of 8 elements.
func Make8[V1, V2, V3, V4, V5, V6, V7, V8 any](first V1, second V2, third V3, fourth V4, fifth V5, sixth V6, seventh V7, eighth V8) T8[V1, V2, V3, V4, V5, V6, V7, V8] {
	return T8[V1, V2, V3, V4, V5, V6, V7, V8]{first, second, third, fourth, fifth, sixth, seventh, eighth}
}

// T9 is a 9-ary tuple.
type T9[V1, V2, V3, V4, V5, V6, V7, V8, V9 any] struct {
	First   V1
	Second  V2
	Third   V3
	Fourth  V4
	Fifth   V5
	Sixth   V6
	Seventh V7
	Eighth  V8
	Ninth   V9
}

// Values returns all elements of tuple.
func (t T9[V1, V2, V3, V4, V5, V6, V7, V8, V9]) Values() (V1, V2, V3, V4, V5, V6, V7, V8, V9) {
	return t.First, t.Second, t.Third, t.Fourth, t.Fifth, t.Sixth, t.Seventh, t.Eighth, t.Ninth
}

// Make9 creates a tuple of 9 elements.
func Make9[V1, V2, V3, V4, V5, V6, V7, V8, V9 any](first V1, second V2, third V3, fourth V4, fifth V5, sixth V6, seventh V7, eighth V8, ninth V9) T9[V1, V2, V3, V4, V5, V6, V7, V8, V9] {
	return T9[V1, V2, V3, V4, V5, V6, V7, V8, V9]{first, second, third, fourth, fifth, sixth, seventh, eighth, ninth}
}

// T10 is a 10-ary tuple.
type T10[V1, V2, V3, V4, V5, V6, V7, V8, V9, V10 any] struct {
	First   V1
	Second  V2
	Third   V3
	Fourth  V4
	Fifth   V5
	Sixth   V6
	Seventh V7
	Eighth  V8
	Ninth   V9
	Tenth   V10
}

// Values returns all elements of tuple.
func (t T10[V1, V2, V3, V4, V5, V6, V7, V8, V9, V10]) Values() (V1, V2, V3, V4, V5, V6, V7, V8, V9, V10) {
	return t.First, t.Second, t.Third, t.Fourth, t.Fifth, t.Sixth, t.Seventh, t.Eighth, t.Ninth, t.Tenth
}

// Make10 creates a tuple of 10 elements.
func Make10[V1, V2, V3, V4, V5, V6, V7, V8, V9, V10 any](first V1, second V2, third V3, fourth V4, fifth V5, sixth V6, seventh V7, eighth V8, ninth V9, tenth V10) T10[V1, V2, V3, V4, V5, V6, V7, V8, V9, V10] {
	return T10[V1, V2, V3, V4, V5, V6, V7, V8, V9, V10]{first, second, third, fourth, fifth, sixth, seventh, eighth, ninth, tenth}
}

// S2 is a slice of 2-ary tuple.
type S2[V1, V2 any] []T2[V1, V2]

// Unzip unpacks elements of tuple to slice.
func (s S2[V1, V2]) Unzip() ([]V1, []V2) {
	s1 := make([]V1, len(s))
	s2 := make([]V2, len(s))
	for i := range s {
		s1[i], s2[i] = s[i].Values()
	}
	return s1, s2
}

func Zip2[V1, V2 any](s1 []V1, s2 []V2) S2[V1, V2] {
	size := gvalue.Min(len(s1), len(s2))
	s := make(S2[V1, V2], size)
	for i := 0; i < size; i++ {
		s[i] = Make2(s1[i], s2[i])
	}
	return s
}

// S3 is a slice of 3-ary tuple.
type S3[V1, V2, V3 any] []T3[V1, V2, V3]

// Unzip unpacks elements of tuple to slice.
func (s S3[V1, V2, V3]) Unzip() ([]V1, []V2, []V3) {
	s1 := make([]V1, len(s))
	s2 := make([]V2, len(s))
	s3 := make([]V3, len(s))
	for i := range s {
		s1[i], s2[i], s3[i] = s[i].Values()
	}
	return s1, s2, s3
}

func Zip3[V1, V2, V3 any](s1 []V1, s2 []V2, s3 []V3) S3[V1, V2, V3] {
	size := gvalue.Min(len(s1), len(s2), len(s3))
	s := make(S3[V1, V2, V3], size)
	for i := 0; i < size; i++ {
		s[i] = Make3(s1[i], s2[i], s3[i])
	}
	return s
}

// S4 is a slice of 4-ary tuple.
type S4[V1, V2, V3, V4 any] []T4[V1, V2, V3, V4]

// Unzip unpacks elements of tuple to slice.
func (s S4[V1, V2, V3, V4]) Unzip() ([]V1, []V2, []V3, []V4) {
	s1 := make([]V1, len(s))
	s2 := make([]V2, len(s))
	s3 := make([]V3, len(s))
	s4 := make([]V4, len(s))
	for i := range s {
		s1[i], s2[i], s3[i], s4[i] = s[i].Values()
	}
	return s1, s2, s3, s4
}

func Zip4[V1, V2, V3, V4 any](s1 []V1, s2 []V2, s3 []V3, s4 []V4) S4[V1, V2, V3, V4] {
	size := gvalue.Min(len(s1), len(s2), len(s3), len(s4))
	s := make(S4[V1, V2, V3, V4], size)
	for i := 0; i < size; i++ {
		s[i] = Make4(s1[i], s2[i], s3[i], s4[i])
	}
	return s
}

// S5 is a slice of 5-ary tuple.
type S5[V1, V2, V3, V4, V5 any] []T5[V1, V2, V3, V4, V5]

// Unzip unpacks elements of tuple to slice.
func (s S5[V1, V2, V3, V4, V5]) Unzip() ([]V1, []V2, []V3, []V4, []V5) {
	s1 := make([]V1, len(s))
	s2 := make([]V2, len(s))
	s3 := make([]V3, len(s))
	s4 := make([]V4, len(s))
	s5 := make([]V5, len(s))
	for i := range s {
		s1[i], s2[i], s3[i], s4[i], s5[i] = s[i].Values()
	}
	return s1, s2, s3, s4, s5
}

func Zip5[V1, V2, V3, V4, V5 any](s1 []V1, s2 []V2, s3 []V3, s4 []V4, s5 []V5) S5[V1, V2, V3, V4, V5] {
	size := gvalue.Min(len(s1), len(s2), len(s3), len(s4), len(s5))
	s := make(S5[V1, V2, V3, V4, V5], size)
	for i := 0; i < size; i++ {
		s[i] = Make5(s1[i], s2[i], s3[i], s4[i], s5[i])
	}
	return s
}

// S6 is a slice of 6-ary tuple.
type S6[V1, V2, V3, V4, V5, V6 any] []T6[V1, V2, V3, V4, V5, V6]

// Unzip unpacks elements of tuple to slice.
func (s S6[V1, V2, V3, V4, V5, V6]) Unzip() ([]V1, []V2, []V3, []V4, []V5, []V6) {
	s1 := make([]V1, len(s))
	s2 := make([]V2, len(s))
	s3 := make([]V3, len(s))
	s4 := make([]V4, len(s))
	s5 := make([]V5, len(s))
	s6 := make([]V6, len(s))
	for i := range s {
		s1[i], s2[i], s3[i], s4[i], s5[i], s6[i] = s[i].Values()
	}
	return s1, s2, s3, s4, s5, s6
}

func Zip6[V1, V2, V3, V4, V5, V6 any](s1 []V1, s2 []V2, s3 []V3, s4 []V4, s5 []V5, s6 []V6) S6[V1, V2, V3, V4, V5, V6] {
	size := gvalue.Min(len(s1), len(s2), len(s3), len(s4), len(s5), len(s6))
	s := make(S6[V1, V2, V3, V4, V5, V6], size)
	for i := 0; i < size; i++ {
		s[i] = Make6(s1[i], s2[i], s3[i], s4[i], s5[i], s6[i])
	}
	return s
}

// S7 is a slice of 7-ary tuple.
type S7[V1, V2, V3, V4, V5, V6, V7 any] []T7[V1, V2, V3, V4, V5, V6, V7]

// Unzip unpacks elements of tuple to slice.
func (s S7[V1, V2, V3, V4, V5, V6, V7]) Unzip() ([]V1, []V2, []V3, []V4, []V5, []V6, []V7) {
	s1 := make([]V1, len(s))
	s2 := make([]V2, len(s))
	s3 := make([]V3, len(s))
	s4 := make([]V4, len(s))
	s5 := make([]V5, len(s))
	s6 := make([]V6, len(s))
	s7 := make([]V7, len(s))
	for i := range s {
		s1[i], s2[i], s3[i], s4[i], s5[i], s6[i], s7[i] = s[i].Values()
	}
	return s1, s2, s3, s4, s5, s6, s7
}

func Zip7[V1, V2, V3, V4, V5, V6, V7 any](s1 []V1, s2 []V2, s3 []V3, s4 []V4, s5 []V5, s6 []V6, s7 []V7) S7[V1, V2, V3, V4, V5, V6, V7] {
	size := gvalue.Min(len(s1), len(s2), len(s3), len(s4), len(s5), len(s6), len(s7))
	s := make(S7[V1, V2, V3, V4, V5, V6, V7], size)
	for i := 0; i < size; i++ {
		s[i] = Make7(s1[i], s2[i], s3[i], s4[i], s5[i], s6[i], s7[i])
	}
	return s
}

// S8 is a slice of 8-ary tuple.
type S8[V1, V2, V3, V4, V5, V6, V7, V8 any] []T8[V1, V2, V3, V4, V5, V6, V7, V8]

// Unzip unpacks elements of tuple to slice.
func (s S8[V1, V2, V3, V4, V5, V6, V7, V8]) Unzip() ([]V1, []V2, []V3, []V4, []V5, []V6, []V7, []V8) {
	s1 := make([]V1, len(s))
	s2 := make([]V2, len(s))
	s3 := make([]V3, len(s))
	s4 := make([]V4, len(s))
	s5 := make([]V5, len(s))
	s6 := make([]V6, len(s))
	s7 := make([]V7, len(s))
	s8 := make([]V8, len(s))
	for i := range s {
		s1[i], s2[i], s3[i], s4[i], s5[i], s6[i], s7[i], s8[i] = s[i].Values()
	}
	return s1, s2, s3, s4, s5, s6, s7, s8
}

func Zip8[V1, V2, V3, V4, V5, V6, V7, V8 any](s1 []V1, s2 []V2, s3 []V3, s4 []V4, s5 []V5, s6 []V6, s7 []V7, s8 []V8) S8[V1, V2, V3, V4, V5, V6, V7, V8] {
	size := gvalue.Min(len(s1), len(s2), len(s3), len(s4), len(s5), len(s6), len(s7), len(s8))
	s := make(S8[V1, V2, V3, V4, V5, V6, V7, V8], size)
	for i := 0; i < size; i++ {
		s[i] = Make8(s1[i], s2[i], s3[i], s4[i], s5[i], s6[i], s7[i], s8[i])
	}
	return s
}

// S9 is a slice of 9-ary tuple.
type S9[V1, V2, V3, V4, V5, V6, V7, V8, V9 any] []T9[V1, V2, V3, V4, V5, V6, V7, V8, V9]

// Unzip unpacks elements of tuple to slice.
func (s S9[V1, V2, V3, V4, V5, V6, V7, V8, V9]) Unzip() ([]V1, []V2, []V3, []V4, []V5, []V6, []V7, []V8, []V9) {
	s1 := make([]V1, len(s))
	s2 := make([]V2, len(s))
	s3 := make([]V3, len(s))
	s4 := make([]V4, len(s))
	s5 := make([]V5, len(s))
	s6 := make([]V6, len(s))
	s7 := make([]V7, len(s))
	s8 := make([]V8, len(s))
	s9 := make([]V9, len(s))
	for i := range s {
		s1[i], s2[i], s3[i], s4[i], s5[i], s6[i], s7[i], s8[i], s9[i] = s[i].Values()
	}
	return s1, s2, s3, s4, s5, s6, s7, s8, s9
}

func Zip9[V1, V2, V3, V4, V5, V6, V7, V8, V9 any](s1 []V1, s2 []V2, s3 []V3, s4 []V4, s5 []V5, s6 []V6, s7 []V7, s8 []V8, s9 []V9) S9[V1, V2, V3, V4, V5, V6, V7, V8, V9] {
	size := gvalue.Min(len(s1), len(s2), len(s3), len(s4), len(s5), len(s6), len(s7), len(s8), len(s9))
	s := make(S9[V1, V2, V3, V4, V5, V6, V7, V8, V9], size)
	for i := 0; i < size; i++ {
		s[i] = Make9(s1[i], s2[i], s3[i], s4[i], s5[i], s6[i], s7[i], s8[i], s9[i])
	}
	return s
}

// S10 is a slice of 10-ary tuple.
type S10[V1, V2, V3, V4, V5, V6, V7, V8, V9, V10 any] []T10[V1, V2, V3, V4, V5, V6, V7, V8, V9, V10]

// Unzip unpacks elements of tuple to slice.
func (s S10[V1, V2, V3, V4, V5, V6, V7, V8, V9, V10]) Unzip() ([]V1, []V2, []V3, []V4, []V5, []V6, []V7, []V8, []V9, []V10) {
	s1 := make([]V1, len(s))
	s2 := make([]V2, len(s))
	s3 := make([]V3, len(s))
	s4 := make([]V4, len(s))
	s5 := make([]V5, len(s))
	s6 := make([]V6, len(s))
	s7 := make([]V7, len(s))
	s8 := make([]V8, len(s))
	s9 := make([]V9, len(s))
	s10 := make([]V10, len(s))
	for i := range s {
		s1[i], s2[i], s3[i], s4[i], s5[i], s6[i], s7[i], s8[i], s9[i], s10[i] = s[i].Values()
	}
	return s1, s2, s3, s4, s5, s6, s7, s8, s9, s10
}

func Zip10[V1, V2, V3, V4, V5, V6, V7, V8, V9, V10 any](s1 []V1, s2 []V2, s3 []V3, s4 []V4, s5 []V5, s6 []V6, s7 []V7, s8 []V8, s9 []V9, s10 []V10) S10[V1, V2, V3, V4, V5, V6, V7, V8, V9, V10] {
	size := gvalue.Min(len(s1), len(s2), len(s3), len(s4), len(s5), len(s6), len(s7), len(s8), len(s9), len(s10))
	s := make(S10[V1, V2, V3, V4, V5, V6, V7, V8, V9, V10], size)
	for i := 0; i < size; i++ {
		s[i] = Make10(s1[i], s2[i], s3[i], s4[i], s5[i], s6[i], s7[i], s8[i], s9[i], s10[i])
	}
	return s
}
