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

// This file contains the additional operations of the [sync] wrappers.
//
// [Map] and [Pool] are already method based, Go 1.27 generic methods add the
// operations that change the element type, such as [Map.ToSlice].
package gsync

// Len returns the number of the items of the map.
func (sm *Map[K, V]) Len() int {
	length := 0
	sm.Range(func(K, V) bool {
		length++
		return true
	})
	return length
}

// Keys returns the keys of the map.
func (sm *Map[K, V]) Keys() []K {
	m := sm.ToMap()
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values returns the values of the map.
func (sm *Map[K, V]) Values() []V {
	m := sm.ToMap()
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// ToSlice converts the map to a slice, the element of the slice is produced by
// function f.
//
// 🚀 EXAMPLE:
//
//	sm.ToSlice(func(k string, v int) string { return k + strconv.Itoa(v) })
//	// ["a1", "b2"]
func (sm *Map[K, V]) ToSlice[R any](f func(K, V) R) []R {
	m := sm.ToMap()
	r := make([]R, 0, len(m))
	for k, v := range m {
		r = append(r, f(k, v))
	}
	return r
}

// With gets a value from the pool, calls f with it and puts it back to the pool.
//
// 🚀 EXAMPLE:
//
//	pool.With(func(i *int) { *i = 10 })
func (p *Pool[T]) With(f func(T)) {
	v := p.Get()
	defer p.Put(v)
	f(v)
}
