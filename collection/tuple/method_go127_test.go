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

package tuple

import (
	"strconv"
	"testing"

	"github.com/bytedance/gg/internal/assert"
)

func TestGo127T2Map(t *testing.T) {
	got := Make2(1, "a").Map(func(i int, s string) (string, int) { return s, i })
	assert.Equal(t, Make2("a", 1), got)
}

func TestGo127T2Swap(t *testing.T) {
	assert.Equal(t, Make2("a", 1), Make2(1, "a").Swap())
}

func TestGo127T3Map(t *testing.T) {
	got := Make3(1, "a", true).Map(func(i int, s string, b bool) (string, int, bool) {
		return s, i, b
	})
	assert.Equal(t, Make3("a", 1, true), got)
}

func TestGo127T4Map(t *testing.T) {
	got := Make4(1, 2, 3, 4).Map(func(a, b, c, d int) (string, string, string, string) {
		return strconv.Itoa(a), strconv.Itoa(b), strconv.Itoa(c), strconv.Itoa(d)
	})
	assert.Equal(t, Make4("1", "2", "3", "4"), got)
}

func TestGo127T10Map(t *testing.T) {
	got := Make10(1, 2, 3, 4, 5, 6, 7, 8, 9, 10).Map(
		func(a, b, c, d, e, f, g, h, i, j int) (int, int, int, int, int, int, int, int, int, int) {
			return a * 2, b * 2, c * 2, d * 2, e * 2, f * 2, g * 2, h * 2, i * 2, j * 2
		})
	assert.Equal(t, Make10(2, 4, 6, 8, 10, 12, 14, 16, 18, 20), got)
}

func TestGo127PairMap(t *testing.T) {
	p := Pair[int, string](Make2(1, "a"))
	assert.Equal(t, Pair[string, int](Make2("a", 1)), p.Map(func(i int, s string) (string, int) { return s, i }))
	assert.Equal(t, Pair[string, int](Make2("a", 1)), p.Swap())
}

func TestGo127S2Map(t *testing.T) {
	s := Zip2([]int{1, 2}, []string{"a", "b"})
	assert.Equal(t, []string{"a", "b"}, s.Map(func(t T2[int, string]) string { return t.Second }))
}

func TestGo127S2Filter(t *testing.T) {
	s := Zip2([]int{1, 2}, []string{"a", "b"})
	filtered := s.Filter(func(t T2[int, string]) bool { return t.First == 1 })
	assert.Equal(t, 1, len(filtered))

	first, second := filtered.Unzip()
	assert.Equal(t, []int{1}, first)
	assert.Equal(t, []string{"a"}, second)
}

func TestGo127TupleMapT5ToT9(t *testing.T) {
	assert.Equal(t, Make5(2, 4, 6, 8, 10), Make5(1, 2, 3, 4, 5).Map(
		func(a, b, c, d, e int) (int, int, int, int, int) {
			return a * 2, b * 2, c * 2, d * 2, e * 2
		}))

	assert.Equal(t, Make6(2, 4, 6, 8, 10, 12), Make6(1, 2, 3, 4, 5, 6).Map(
		func(a, b, c, d, e, f int) (int, int, int, int, int, int) {
			return a * 2, b * 2, c * 2, d * 2, e * 2, f * 2
		}))

	assert.Equal(t, Make7(2, 4, 6, 8, 10, 12, 14), Make7(1, 2, 3, 4, 5, 6, 7).Map(
		func(a, b, c, d, e, f, g int) (int, int, int, int, int, int, int) {
			return a * 2, b * 2, c * 2, d * 2, e * 2, f * 2, g * 2
		}))

	assert.Equal(t, Make8(2, 4, 6, 8, 10, 12, 14, 16), Make8(1, 2, 3, 4, 5, 6, 7, 8).Map(
		func(a, b, c, d, e, f, g, h int) (int, int, int, int, int, int, int, int) {
			return a * 2, b * 2, c * 2, d * 2, e * 2, f * 2, g * 2, h * 2
		}))

	assert.Equal(t, Make9(2, 4, 6, 8, 10, 12, 14, 16, 18), Make9(1, 2, 3, 4, 5, 6, 7, 8, 9).Map(
		func(a, b, c, d, e, f, g, h, i int) (int, int, int, int, int, int, int, int, int) {
			return a * 2, b * 2, c * 2, d * 2, e * 2, f * 2, g * 2, h * 2, i * 2
		}))
}

func intSlice() []int { return []int{1, 2} }

func TestGo127SMapFilterS3ToS10(t *testing.T) {
	s3 := Zip3(intSlice(), intSlice(), intSlice())
	assert.Equal(t, []int{1, 2}, s3.Map(func(t T3[int, int, int]) int { return t.First }))
	assert.Equal(t, 1, len(s3.Filter(func(t T3[int, int, int]) bool { return t.First == 1 })))

	s4 := Zip4(intSlice(), intSlice(), intSlice(), intSlice())
	assert.Equal(t, []int{1, 2}, s4.Map(func(t T4[int, int, int, int]) int { return t.First }))
	assert.Equal(t, 1, len(s4.Filter(func(t T4[int, int, int, int]) bool { return t.First == 1 })))

	s5 := Zip5(intSlice(), intSlice(), intSlice(), intSlice(), intSlice())
	assert.Equal(t, []int{1, 2}, s5.Map(func(t T5[int, int, int, int, int]) int { return t.First }))
	assert.Equal(t, 1, len(s5.Filter(func(t T5[int, int, int, int, int]) bool { return t.First == 1 })))

	s6 := Zip6(intSlice(), intSlice(), intSlice(), intSlice(), intSlice(), intSlice())
	assert.Equal(t, []int{1, 2}, s6.Map(func(t T6[int, int, int, int, int, int]) int { return t.First }))
	assert.Equal(t, 1, len(s6.Filter(func(t T6[int, int, int, int, int, int]) bool { return t.First == 1 })))

	s7 := Zip7(intSlice(), intSlice(), intSlice(), intSlice(), intSlice(), intSlice(), intSlice())
	assert.Equal(t, []int{1, 2}, s7.Map(func(t T7[int, int, int, int, int, int, int]) int { return t.First }))
	assert.Equal(t, 1, len(s7.Filter(func(t T7[int, int, int, int, int, int, int]) bool {
		return t.First == 1
	})))

	s8 := Zip8(intSlice(), intSlice(), intSlice(), intSlice(), intSlice(), intSlice(), intSlice(), intSlice())
	assert.Equal(t, []int{1, 2}, s8.Map(func(t T8[int, int, int, int, int, int, int, int]) int {
		return t.First
	}))
	assert.Equal(t, 1, len(s8.Filter(func(t T8[int, int, int, int, int, int, int, int]) bool {
		return t.First == 1
	})))

	s9 := Zip9(intSlice(), intSlice(), intSlice(), intSlice(), intSlice(), intSlice(), intSlice(), intSlice(), intSlice())
	assert.Equal(t, []int{1, 2}, s9.Map(func(t T9[int, int, int, int, int, int, int, int, int]) int {
		return t.First
	}))
	assert.Equal(t, 1, len(s9.Filter(func(t T9[int, int, int, int, int, int, int, int, int]) bool {
		return t.First == 1
	})))

	s10 := Zip10(intSlice(), intSlice(), intSlice(), intSlice(), intSlice(), intSlice(), intSlice(), intSlice(), intSlice(), intSlice())
	assert.Equal(t, []int{1, 2}, s10.Map(func(t T10[int, int, int, int, int, int, int, int, int, int]) int {
		return t.First
	}))
	assert.Equal(t, 1, len(s10.Filter(func(t T10[int, int, int, int, int, int, int, int, int, int]) bool {
		return t.First == 1
	})))
}
