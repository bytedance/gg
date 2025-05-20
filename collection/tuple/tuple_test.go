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
	p := Make2("red", 14)

	if p.First != "red" {
		t.Error()
	}
	if p.Second != 14 {
		t.Error()
	}
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
