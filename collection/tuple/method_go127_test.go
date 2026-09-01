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
