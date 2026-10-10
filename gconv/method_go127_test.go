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

package gconv

import (
	"testing"

	"github.com/bytedance/gg/internal/assert"
)

func TestVTo(t *testing.T) {
	assert.Equal(t, 1, Of("1").To[int]())
	assert.Equal(t, "1", Of(1).To[string]())
	assert.Equal(t, true, Of("true").To[bool]())
	assert.Equal(t, 0, Of("x").To[int]())
	assert.Equal(t, "1", Of("1").Get())
}

func TestVToPtr(t *testing.T) {
	assert.Equal(t, 1, *Of("1").ToPtr[int]())
	assert.Nil(t, Of("x").ToPtr[int]())
}

func TestVToR(t *testing.T) {
	assert.Equal(t, 1, Of("1").ToR[int]().Value())
	assert.True(t, Of("x").ToR[int]().IsErr())
}

func TestVToE(t *testing.T) {
	v, err := Of("1").ToE[int]()
	assert.Equal(t, 1, v)
	assert.Nil(t, err)

	v, err = Of("x").ToE[int]()
	assert.Equal(t, 0, v)
	assert.NotNil(t, err)
}
