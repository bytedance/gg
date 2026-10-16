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
	"reflect"
	"testing"
)

func consumeIntersperse(input []int, sep int, chunks []int) []int {
	i := Intersperse(sep, FromSlice(input))
	var got []int
	for _, n := range chunks {
		got = append(got, i.Next(n)...)
	}
	got = append(got, i.Next(ALL)...)
	return got
}

func TestIntersperseChunkedReadsHaveNoTrailingSeparator(t *testing.T) {
	input := []int{1, 2, 3}
	want := []int{1, 9, 2, 9, 3}
	for _, chunks := range [][]int{
		{1},
		{2},
		{3},
		{1, 1},
		{2, 1},
		{3, 2, 1}, // used to return [1 9 2 9 3 9]
		{4, 1, 1},
		{5, 1},
	} {
		if got := consumeIntersperse(input, 9, chunks); !reflect.DeepEqual(got, want) {
			t.Errorf("chunks %v: got %v, want %v", chunks, got, want)
		}
	}
}

func TestIntersperseChunkedMatchesAllAtOnce(t *testing.T) {
	for length := 0; length <= 64; length++ {
		input := make([]int, length)
		for j := range input {
			input[j] = j
		}
		want := ToSlice(Intersperse(-1, FromSlice(input)))
		for chunkSize := 1; chunkSize <= 16; chunkSize++ {
			i := Intersperse(-1, FromSlice(input))
			var got []int
			for {
				part := i.Next(chunkSize)
				if len(part) == 0 {
					break
				}
				got = append(got, part...)
			}
			if len(got) != len(want) || (len(got) != 0 && !reflect.DeepEqual(got, want)) {
				t.Fatalf("length=%d chunkSize=%d: got %v, want %v", length, chunkSize, got, want)
			}
		}
	}
}
