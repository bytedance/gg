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

package jsonbuilder

import (
	"testing"

	"github.com/bytedance/gg/internal/assert"
)

func TestArrayBuild(t *testing.T) {
	{
		s := []int{1, 2, 3, 4, 5}
		a := NewArray()
		for _, v := range s {
			err := a.Append(v)
			assert.Nil(t, err)
		}
		bs, err := a.Build()
		assert.Nil(t, err)
		assert.Equal(t, []byte(`[1,2,3,4,5]`), bs)
	}

	{
		s := []int{}
		a := NewArray()
		for _, v := range s {
			err := a.Append(v)
			assert.Nil(t, err)
		}
		bs, err := a.Build()
		assert.Nil(t, err)
		assert.Equal(t, []byte(`[]`), bs)
	}

	{
		s := []string{"a"}
		a := NewArray()
		for _, v := range s {
			err := a.Append(v)
			assert.Nil(t, err)
		}
		bs, err := a.Build()
		assert.Nil(t, err)
		assert.Equal(t, []byte(`["a"]`), bs)
	}

	{
		var a *Array
		bs, err := a.Build()
		assert.Nil(t, err)
		assert.Equal(t, []byte(`null`), bs)
	}
}
