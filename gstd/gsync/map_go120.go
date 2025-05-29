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

//go:build go1.20
// +build go1.20

package gsync

import (
	"github.com/bytedance/gg/gvalue"
)

// Swap wraps [sync.Map.Swap].
//
// 💡 NOTE: Newly added in go1.20
func (sm *Map[K, V]) Swap(key K, value V) (V, bool) {
	previous, loaded := sm.m.Swap(key, value)
	if loaded {
		return previous.(V), true
	}
	return gvalue.Zero[V](), false
}

// CompareAndSwap wraps [sync.Map.CompareAndSwap].
//
// 💡 NOTE: Newly added in go1.20
func (sm *Map[K, V]) CompareAndSwap(key K, old, new V) bool {
	return sm.m.CompareAndSwap(key, old, new)
}

// CompareAndDelete wraps [sync.Map.CompareAndDelete].
//
// 💡 NOTE: Newly added in go1.20
func (sm *Map[K, V]) CompareAndDelete(key K, old V) bool {
	return sm.m.CompareAndDelete(key, old)
}
