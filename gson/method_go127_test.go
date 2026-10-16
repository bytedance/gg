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

package gson

import (
	"testing"

	"github.com/bytedance/gg/internal/assert"
)

type go127Struct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestGo127V(t *testing.T) {
	v := Of(go127Struct{Name: "test", Age: 10})

	assert.Equal(t, `{"name":"test","age":10}`, v.ToString())

	b, err := v.Marshal()
	assert.Nil(t, err)
	assert.Equal(t, `{"name":"test","age":10}`, string(b))

	s, err := v.MarshalString()
	assert.Nil(t, err)
	assert.Equal(t, `{"name":"test","age":10}`, s)

	bi, err := v.MarshalIndent("", "  ")
	assert.Nil(t, err)
	assert.Equal(t, "{\n  \"name\": \"test\",\n  \"age\": 10\n}", string(bi))
	assert.Equal(t, "{\n  \"name\": \"test\",\n  \"age\": 10\n}", v.ToStringIndent("", "  "))

	// MarshalBy accepts any codec.
	b2, err := v.MarshalBy(stdJSON)
	assert.Nil(t, err)
	assert.Equal(t, `{"name":"test","age":10}`, string(b2))

	assert.Equal(t, 10, v.Get().Age)
}
