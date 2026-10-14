// Copyright 2026 Bytedance Ltd.
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
	"reflect"
	"testing"
)

// sameSeq treats nil and empty slices as equal.
func sameSeq(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	if len(a) == 0 {
		return true
	}
	return reflect.DeepEqual(a, b)
}

// collectChunked pulls with fixed chunk n until exhaustion and returns the
// concatenation.
func collectChunked[T any](i Iter[T], n int) []T {
	var got []T
	for {
		part := i.Next(n)
		if len(part) == 0 {
			break
		}
		got = append(got, part...)
	}
	return got
}

// Regression for bytedance/gg#47: Intersperse wrapped around a FilterMap
// dropped accepted elements when consumed in chunks. FilterMap used to report
// an empty result as soon as a single upstream batch was fully filtered out,
// which the Peeker feeding Intersperse mistook for exhaustion. Here the source
// [9, 0, 3] filters its way down to [9, 3]; chunked reads must match the
// Next(ALL) oracle [9, 0, 3].
func TestIntersperseOverFilterMapChunked(t *testing.T) {
	nonZero := func(x int) (int, bool) { return x, x != 0 }
	input := []int{9, 0, 3}
	oracle := ToSlice(Intersperse(0, FilterMap(nonZero, FromSlice(input))))
	if want := []int{9, 0, 3}; !sameSeq(oracle, want) {
		t.Fatalf("oracle = %v, want %v", oracle, want)
	}
	for _, chunk := range []int{1, 2, 3} {
		got := collectChunked(Intersperse(0, FilterMap(nonZero, FromSlice(input))), chunk)
		if !sameSeq(got, oracle) {
			t.Fatalf("chunk=%d: got %v, want %v", chunk, got, oracle)
		}
	}
}

// FilterMap over Intersperse must also keep dropping separators without
// prematurely terminating: chunked reads must equal the Next(ALL) oracle.
func TestFilterMapOverIntersperseChunked(t *testing.T) {
	dropSep := func(x int) (int, bool) { return x, x != 0 }
	for length := 0; length <= 24; length++ {
		input := make([]int, length)
		for j := range input {
			input[j] = j + 1
		}
		oracle := ToSlice(FilterMap(dropSep, Intersperse(0, FromSlice(input))))
		for chunk := 1; chunk <= 6; chunk++ {
			got := collectChunked(FilterMap(dropSep, Intersperse(0, FromSlice(input))), chunk)
			if !sameSeq(got, oracle) {
				t.Fatalf("length=%d chunk=%d: got %v, want %v", length, chunk, got, oracle)
			}
		}
	}
}

// Differential sweep: chunked consumption of Intersperse(FilterMap(...)) must
// match the Next(ALL) oracle across lengths and chunk sizes.
func TestIntersperseFilterMapChunkedMatchesAllAtOnce(t *testing.T) {
	keepOdd := func(x int) (int, bool) { return x, x%2 == 1 }
	for length := 0; length <= 40; length++ {
		input := make([]int, length)
		for j := range input {
			input[j] = j + 1
		}
		oracle := ToSlice(Intersperse(0, FilterMap(keepOdd, FromSlice(input))))
		for chunk := 1; chunk <= 10; chunk++ {
			got := collectChunked(Intersperse(0, FilterMap(keepOdd, FromSlice(input))), chunk)
			if !sameSeq(got, oracle) {
				t.Fatalf("length=%d chunk=%d: got %v, want %v", length, chunk, got, oracle)
			}
		}
	}
}
