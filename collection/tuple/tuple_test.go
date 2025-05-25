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

package tuple

import (
	"testing"

	"github.com/bytedance/gg/internal/assert"
)

func TestT2(t *testing.T) {
	t2 := Make2("red", 14)
	a, b := t2.Values()
	assert.Equal(t, "red", a)
	assert.Equal(t, 14, b)
}

func TestT3(t *testing.T) {
	t3 := Make3("red", 14, 15)
	a, b, c := t3.Values()
	assert.Equal(t, "red", a)
	assert.Equal(t, 14, b)
	assert.Equal(t, 15, c)
}

func TestT4(t *testing.T) {
	t4 := Make4("red", 14, 15, 16)
	a, b, c, d := t4.Values()
	assert.Equal(t, "red", a)
	assert.Equal(t, 14, b)
	assert.Equal(t, 15, c)
	assert.Equal(t, 16, d)
}

func TestT5(t *testing.T) {
	t5 := Make5("red", 14, 15, 16, 17)
	a, b, c, d, e := t5.Values()
	assert.Equal(t, "red", a)
	assert.Equal(t, 14, b)
	assert.Equal(t, 15, c)
	assert.Equal(t, 16, d)
	assert.Equal(t, 17, e)
}

func TestT6(t *testing.T) {
	t6 := Make6("red", 14, 15, 16, 17, 18)
	a, b, c, d, e, f := t6.Values()
	assert.Equal(t, "red", a)
	assert.Equal(t, 14, b)
	assert.Equal(t, 15, c)
	assert.Equal(t, 16, d)
	assert.Equal(t, 17, e)
	assert.Equal(t, 18, f)
}

func TestT7(t *testing.T) {
	t7 := Make7("red", 14, 15, 16, 17, 18, 19)
	a, b, c, d, e, f, g := t7.Values()
	assert.Equal(t, "red", a)
	assert.Equal(t, 14, b)
	assert.Equal(t, 15, c)
	assert.Equal(t, 16, d)
	assert.Equal(t, 17, e)
	assert.Equal(t, 18, f)
	assert.Equal(t, 19, g)
}

func TestT8(t *testing.T) {
	t8 := Make8("red", 14, 15, 16, 17, 18, 19, 20)
	a, b, c, d, e, f, g, h := t8.Values()
	assert.Equal(t, "red", a)
	assert.Equal(t, 14, b)
	assert.Equal(t, 15, c)
	assert.Equal(t, 16, d)
	assert.Equal(t, 17, e)
	assert.Equal(t, 18, f)
	assert.Equal(t, 19, g)
	assert.Equal(t, 20, h)
}

func TestT9(t *testing.T) {
	t9 := Make9("red", 14, 15, 16, 17, 18, 19, 20, 21)
	a, b, c, d, e, f, g, h, i := t9.Values()
	assert.Equal(t, "red", a)
	assert.Equal(t, 14, b)
	assert.Equal(t, 15, c)
	assert.Equal(t, 16, d)
	assert.Equal(t, 17, e)
	assert.Equal(t, 18, f)
	assert.Equal(t, 19, g)
	assert.Equal(t, 20, h)
	assert.Equal(t, 21, i)
}

func TestT10(t *testing.T) {
	t10 := Make10("red", 14, 15, 16, 17, 18, 19, 20, 21, 22)
	a, b, c, d, e, f, g, h, i, j := t10.Values()
	assert.Equal(t, "red", a)
	assert.Equal(t, 14, b)
	assert.Equal(t, 15, c)
	assert.Equal(t, 16, d)
	assert.Equal(t, 17, e)
	assert.Equal(t, 18, f)
	assert.Equal(t, 19, g)
	assert.Equal(t, 20, h)
	assert.Equal(t, 21, i)
	assert.Equal(t, 22, j)
}

func TestS2(t *testing.T) {
	{
		s := Zip2([]string{"red", "green", "blue"}, []int{14, 15, 16})
		s1, s2 := s.Unzip()
		assert.Equal(t, []string{"red", "green", "blue"}, s1)
		assert.Equal(t, []int{14, 15, 16}, s2)
	}
	{ // Test empty.
		s := Zip2([]string{}, []int{})
		s1, s2 := s.Unzip()
		assert.Equal(t, []string{}, s1)
		assert.Equal(t, []int{}, s2)
	}
	{ // Test nil.
		s := Zip2([]string(nil), []int(nil))
		s1, s2 := s.Unzip()
		assert.Equal(t, []string{}, s1)
		assert.Equal(t, []int{}, s2)
	}
}

func TestS3(t *testing.T) {
	s := Zip3([]string{"red", "green", "blue"}, []int{14, 15, 16}, []float64{1.4, 1.5, 1.6})
	s1, s2, s3 := s.Unzip()
	assert.Equal(t, []string{"red", "green", "blue"}, s1)
	assert.Equal(t, []int{14, 15, 16}, s2)
	assert.Equal(t, []float64{1.4, 1.5, 1.6}, s3)
}

func TestS4(t *testing.T) {
	s := Zip4([]string{"red", "green", "blue"}, []int{14, 15, 16}, []float64{1.4, 1.5, 1.6}, []bool{true, false, true})
	s1, s2, s3, s4 := s.Unzip()
	assert.Equal(t, []string{"red", "green", "blue"}, s1)
	assert.Equal(t, []int{14, 15, 16}, s2)
	assert.Equal(t, []float64{1.4, 1.5, 1.6}, s3)
	assert.Equal(t, []bool{true, false, true}, s4)
}

func TestS5(t *testing.T) {
	s := Zip5([]string{"red", "green", "blue"}, []int{14, 15, 16}, []float64{1.4, 1.5, 1.6}, []bool{true, false, true}, []string{"red", "green", "blue"})
	s1, s2, s3, s4, s5 := s.Unzip()
	assert.Equal(t, []string{"red", "green", "blue"}, s1)
	assert.Equal(t, []int{14, 15, 16}, s2)
	assert.Equal(t, []float64{1.4, 1.5, 1.6}, s3)
	assert.Equal(t, []bool{true, false, true}, s4)
	assert.Equal(t, []string{"red", "green", "blue"}, s5)
}

func TestS6(t *testing.T) {
	s := Zip6([]string{"red", "green", "blue"}, []int{14, 15, 16}, []float64{1.4, 1.5, 1.6}, []bool{true, false, true}, []string{"red", "green", "blue"}, []int{14, 15, 16})
	s1, s2, s3, s4, s5, s6 := s.Unzip()
	assert.Equal(t, []string{"red", "green", "blue"}, s1)
	assert.Equal(t, []int{14, 15, 16}, s2)
	assert.Equal(t, []float64{1.4, 1.5, 1.6}, s3)
	assert.Equal(t, []bool{true, false, true}, s4)
	assert.Equal(t, []string{"red", "green", "blue"}, s5)
	assert.Equal(t, []int{14, 15, 16}, s6)
}

func TestS7(t *testing.T) {
	s := Zip7([]string{"red", "green", "blue"}, []int{14, 15, 16}, []float64{1.4, 1.5, 1.6}, []bool{true, false, true}, []string{"red", "green", "blue"}, []int{14, 15, 16}, []float64{1.4, 1.5, 1.6})
	s1, s2, s3, s4, s5, s6, s7 := s.Unzip()
	assert.Equal(t, []string{"red", "green", "blue"}, s1)
	assert.Equal(t, []int{14, 15, 16}, s2)
	assert.Equal(t, []float64{1.4, 1.5, 1.6}, s3)
	assert.Equal(t, []bool{true, false, true}, s4)
	assert.Equal(t, []string{"red", "green", "blue"}, s5)
	assert.Equal(t, []int{14, 15, 16}, s6)
	assert.Equal(t, []float64{1.4, 1.5, 1.6}, s7)
}

func TestS8(t *testing.T) {
	s := Zip8([]string{"red", "green", "blue"}, []int{14, 15, 16}, []float64{1.4, 1.5, 1.6}, []bool{true, false, true}, []string{"red", "green", "blue"}, []int{14, 15, 16}, []float64{1.4, 1.5, 1.6}, []bool{true, false, true})
	s1, s2, s3, s4, s5, s6, s7, s8 := s.Unzip()
	assert.Equal(t, []string{"red", "green", "blue"}, s1)
	assert.Equal(t, []int{14, 15, 16}, s2)
	assert.Equal(t, []float64{1.4, 1.5, 1.6}, s3)
	assert.Equal(t, []bool{true, false, true}, s4)
	assert.Equal(t, []string{"red", "green", "blue"}, s5)
	assert.Equal(t, []int{14, 15, 16}, s6)
	assert.Equal(t, []float64{1.4, 1.5, 1.6}, s7)
	assert.Equal(t, []bool{true, false, true}, s8)
}

func TestS9(t *testing.T) {
	s := Zip9([]string{"red", "green", "blue"}, []int{14, 15, 16}, []float64{1.4, 1.5, 1.6}, []bool{true, false, true}, []string{"red", "green", "blue"}, []int{14, 15, 16}, []float64{1.4, 1.5, 1.6}, []bool{true, false, true}, []string{"red", "green", "blue"})
	s1, s2, s3, s4, s5, s6, s7, s8, s9 := s.Unzip()
	assert.Equal(t, []string{"red", "green", "blue"}, s1)
	assert.Equal(t, []int{14, 15, 16}, s2)
	assert.Equal(t, []float64{1.4, 1.5, 1.6}, s3)
	assert.Equal(t, []bool{true, false, true}, s4)
	assert.Equal(t, []string{"red", "green", "blue"}, s5)
	assert.Equal(t, []int{14, 15, 16}, s6)
	assert.Equal(t, []float64{1.4, 1.5, 1.6}, s7)
	assert.Equal(t, []bool{true, false, true}, s8)
	assert.Equal(t, []string{"red", "green", "blue"}, s9)
}

func TestS10(t *testing.T) {
	s := Zip10([]string{"red", "green", "blue"}, []int{14, 15, 16}, []float64{1.4, 1.5, 1.6}, []bool{true, false, true}, []string{"red", "green", "blue"}, []int{14, 15, 16}, []float64{1.4, 1.5, 1.6}, []bool{true, false, true}, []string{"red", "green", "blue"}, []int{14, 15, 16})
	s1, s2, s3, s4, s5, s6, s7, s8, s9, s10 := s.Unzip()
	assert.Equal(t, []string{"red", "green", "blue"}, s1)
	assert.Equal(t, []int{14, 15, 16}, s2)
	assert.Equal(t, []float64{1.4, 1.5, 1.6}, s3)
	assert.Equal(t, []bool{true, false, true}, s4)
	assert.Equal(t, []string{"red", "green", "blue"}, s5)
	assert.Equal(t, []int{14, 15, 16}, s6)
	assert.Equal(t, []float64{1.4, 1.5, 1.6}, s7)
	assert.Equal(t, []bool{true, false, true}, s8)
	assert.Equal(t, []string{"red", "green", "blue"}, s9)
	assert.Equal(t, []int{14, 15, 16}, s10)
}
