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

// Package gsync provides generics wrappers of [sync] package.
//
// Currently, we provide these wrappers: [Map], [Pool].
// If you want to initialize value with [sync.Once],
// we recommend [github.com/bytedance/gg/gvalue.Once].
package gsync

import (
	"sync"

	"github.com/bytedance/gg/goption"
	"github.com/bytedance/gg/gvalue"
)

// Map wraps [sync.Map].
type Map[K comparable, V any] struct {
	m sync.Map
}

// Load wraps [sync.Map.Load].
func (sm *Map[K, V]) Load(key K) (V, bool) {
	v, ok := sm.m.Load(key)
	if !ok {
		return gvalue.Zero[V](), false
	}
	return v.(V), true
}

// LoadO wraps [Load], returns a goption value instead of (V, bool).
func (sm *Map[K, V]) LoadO(key K) goption.O[V] {
	return goption.Of(sm.Load(key))
}

// Store wraps [sync.Map.Store].
func (sm *Map[K, V]) Store(key K, value V) {
	sm.m.Store(key, value)
}

// LoadOrStore wraps [sync.Map.LoadOrStore].
func (sm *Map[K, V]) LoadOrStore(key K, value V) (V, bool) {
	v, loaded := sm.m.LoadOrStore(key, value)
	if loaded {
		return v.(V), true
	}
	return value, false
}

// LoadAndDelete wraps [sync.Map.LoadAndDelete].
func (sm *Map[K, V]) LoadAndDelete(key K) (V, bool) {
	v, loaded := sm.m.LoadAndDelete(key)
	if loaded {
		return v.(V), true
	}
	return gvalue.Zero[V](), false
}

// Delete wraps [sync.Map.Delete].
func (sm *Map[K, V]) Delete(key K) {
	sm.m.Delete(key)
}

// Range wraps [sync.Map.Range].
func (sm *Map[K, V]) Range(f func(K, V) bool) {
	sm.m.Range(func(key, value any) bool {
		return f(key.(K), value.(V))
	})
}

// ToMap collects all keys and values to a go builtin map.
func (sm *Map[K, V]) ToMap() map[K]V {
	m := make(map[K]V)
	sm.Range(func(k K, v V) bool {
		m[k] = v
		return true
	})
	return m
}
