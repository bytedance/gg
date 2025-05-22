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

package gsync

import (
	"testing"

	"github.com/bytedance/gg/goption"
	"github.com/bytedance/gg/internal/assert"
)

func TestStoreLoad(t *testing.T) {
	sm := Map[string, string]{}
	assert.False(t, sm.LoadO("k").IsOK())
	sm.Store("k", "")
	assert.Equal(t, goption.OK(""), sm.LoadO("k"))
	sm.Store("k", "v")
	sm.Store("k1", "v1")
	assert.Equal(t, goption.OK("v"), sm.LoadO("k"))
	assert.Equal(t, goption.OK("v1"), sm.LoadO("k1"))
	assert.False(t, sm.LoadO("k2").IsOK())
}

func TestDelete(t *testing.T) {
	sm := Map[string, string]{}
	sm.Store("k", "v")
	assert.Equal(t, goption.OK("v"), sm.LoadO("k"))
	sm.Delete("k")
	assert.False(t, sm.LoadO("k").IsOK())
}

func TestLoadOrStore(t *testing.T) {
	sm := Map[string, string]{}
	sm.Store("k", "v")
	v, ok := sm.LoadOrStore("k", "v1")
	assert.Equal(t, "v", v)
	assert.Equal(t, true, ok)
	v, ok = sm.LoadOrStore("k1", "v1")
	assert.Equal(t, "v1", v)
	assert.Equal(t, false, ok)
	assert.Equal(t, goption.OK("v1"), sm.LoadO("k1"))
}

func TestLoadAndDelete(t *testing.T) {
	sm := Map[string, string]{}
	sm.Store("k", "v")
	v, ok := sm.LoadAndDelete("k")
	assert.Equal(t, "v", v)
	assert.Equal(t, true, ok)
	assert.False(t, sm.LoadO("k").IsOK())
	v, ok = sm.LoadAndDelete("k1")
	assert.Equal(t, "", v)
	assert.Equal(t, false, ok)
}

func TestRange(t *testing.T) {
	sm := Map[string, string]{}
	sm.Store("k", "v")
	sm.Store("k1", "v1")
	sm.Store("k2", "v2")
	str := ""
	count := 0
	f := func(k, v string) bool {
		str = str + k + v
		count++
		return true
	}
	sm.Range(f)
	assert.Equal(t, len("kvk1v1k2v2"), len(str))
	assert.Equal(t, 3, count)
}

func TestToMap(t *testing.T) {
	sm := Map[string, string]{}
	sm.Store("k", "v")
	sm.Store("k1", "v1")
	sm.Store("k2", "v2")
	assert.Equal(t, map[string]string{"k": "v", "k1": "v1", "k2": "v2"}, sm.ToMap())
}
