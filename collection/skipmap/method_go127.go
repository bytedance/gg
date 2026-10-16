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

// This file contains the additional operations of the concurrent-safe
// containers.
//
// Before Go 1.27, a method can not declare its own type parameters, so the
// operations that were not tied to the type parameters of the container were
// still missing as methods.
package skipmap

import "github.com/bytedance/gg/goption"

// Keys returns the keys of the map.
func (s *FuncMap[keyT, valueT]) Keys() []keyT {
	keys := make([]keyT, 0, s.Len())
	s.Range(func(key keyT, value valueT) bool {
		keys = append(keys, key)
		return true
	})
	return keys
}

// Values returns the values of the map.
func (s *FuncMap[keyT, valueT]) Values() []valueT {
	values := make([]valueT, 0, s.Len())
	s.Range(func(key keyT, value valueT) bool {
		values = append(values, value)
		return true
	})
	return values
}

// Any returns whether any item of the map satisfies f.
func (s *FuncMap[keyT, valueT]) Any(f func(keyT, valueT) bool) bool {
	found := false
	s.Range(func(key keyT, value valueT) bool {
		if f(key, value) {
			found = true
			return false
		}
		return true
	})
	return found
}

// All returns whether all items of the map satisfy f.
//
// 💡 NOTE: All returns true for an empty map.
func (s *FuncMap[keyT, valueT]) All(f func(keyT, valueT) bool) bool {
	ok := true
	s.Range(func(key keyT, value valueT) bool {
		if !f(key, value) {
			ok = false
			return false
		}
		return true
	})
	return ok
}

// Find returns the value of the first item that satisfies f.
func (s *FuncMap[keyT, valueT]) Find(f func(keyT, valueT) bool) goption.O[valueT] {
	var (
		found valueT
		ok    bool
	)
	s.Range(func(key keyT, value valueT) bool {
		if f(key, value) {
			found, ok = value, true
			return false
		}
		return true
	})
	return goption.Of(found, ok)
}

// FindKey returns the key of the first item that satisfies f.
func (s *FuncMap[keyT, valueT]) FindKey(f func(keyT, valueT) bool) goption.O[keyT] {
	var (
		found keyT
		ok    bool
	)
	s.Range(func(key keyT, value valueT) bool {
		if f(key, value) {
			found, ok = key, true
			return false
		}
		return true
	})
	return goption.Of(found, ok)
}

// ForEach calls function f for each item of the map.
func (s *FuncMap[keyT, valueT]) ForEach(f func(keyT, valueT)) {
	s.Range(func(key keyT, value valueT) bool {
		f(key, value)
		return true
	})
}

// Keys returns the keys of the map.
func (s *OrderedMap[keyT, valueT]) Keys() []keyT {
	keys := make([]keyT, 0, s.Len())
	s.Range(func(key keyT, value valueT) bool {
		keys = append(keys, key)
		return true
	})
	return keys
}

// Values returns the values of the map.
func (s *OrderedMap[keyT, valueT]) Values() []valueT {
	values := make([]valueT, 0, s.Len())
	s.Range(func(key keyT, value valueT) bool {
		values = append(values, value)
		return true
	})
	return values
}

// Any returns whether any item of the map satisfies f.
func (s *OrderedMap[keyT, valueT]) Any(f func(keyT, valueT) bool) bool {
	found := false
	s.Range(func(key keyT, value valueT) bool {
		if f(key, value) {
			found = true
			return false
		}
		return true
	})
	return found
}

// All returns whether all items of the map satisfy f.
//
// 💡 NOTE: All returns true for an empty map.
func (s *OrderedMap[keyT, valueT]) All(f func(keyT, valueT) bool) bool {
	ok := true
	s.Range(func(key keyT, value valueT) bool {
		if !f(key, value) {
			ok = false
			return false
		}
		return true
	})
	return ok
}

// Find returns the value of the first item that satisfies f.
func (s *OrderedMap[keyT, valueT]) Find(f func(keyT, valueT) bool) goption.O[valueT] {
	var (
		found valueT
		ok    bool
	)
	s.Range(func(key keyT, value valueT) bool {
		if f(key, value) {
			found, ok = value, true
			return false
		}
		return true
	})
	return goption.Of(found, ok)
}

// FindKey returns the key of the first item that satisfies f.
func (s *OrderedMap[keyT, valueT]) FindKey(f func(keyT, valueT) bool) goption.O[keyT] {
	var (
		found keyT
		ok    bool
	)
	s.Range(func(key keyT, value valueT) bool {
		if f(key, value) {
			found, ok = key, true
			return false
		}
		return true
	})
	return goption.Of(found, ok)
}

// ForEach calls function f for each item of the map.
func (s *OrderedMap[keyT, valueT]) ForEach(f func(keyT, valueT)) {
	s.Range(func(key keyT, value valueT) bool {
		f(key, value)
		return true
	})
}

// Keys returns the keys of the map.
func (s *OrderedMapDesc[keyT, valueT]) Keys() []keyT {
	keys := make([]keyT, 0, s.Len())
	s.Range(func(key keyT, value valueT) bool {
		keys = append(keys, key)
		return true
	})
	return keys
}

// Values returns the values of the map.
func (s *OrderedMapDesc[keyT, valueT]) Values() []valueT {
	values := make([]valueT, 0, s.Len())
	s.Range(func(key keyT, value valueT) bool {
		values = append(values, value)
		return true
	})
	return values
}

// Any returns whether any item of the map satisfies f.
func (s *OrderedMapDesc[keyT, valueT]) Any(f func(keyT, valueT) bool) bool {
	found := false
	s.Range(func(key keyT, value valueT) bool {
		if f(key, value) {
			found = true
			return false
		}
		return true
	})
	return found
}

// All returns whether all items of the map satisfy f.
//
// 💡 NOTE: All returns true for an empty map.
func (s *OrderedMapDesc[keyT, valueT]) All(f func(keyT, valueT) bool) bool {
	ok := true
	s.Range(func(key keyT, value valueT) bool {
		if !f(key, value) {
			ok = false
			return false
		}
		return true
	})
	return ok
}

// Find returns the value of the first item that satisfies f.
func (s *OrderedMapDesc[keyT, valueT]) Find(f func(keyT, valueT) bool) goption.O[valueT] {
	var (
		found valueT
		ok    bool
	)
	s.Range(func(key keyT, value valueT) bool {
		if f(key, value) {
			found, ok = value, true
			return false
		}
		return true
	})
	return goption.Of(found, ok)
}

// FindKey returns the key of the first item that satisfies f.
func (s *OrderedMapDesc[keyT, valueT]) FindKey(f func(keyT, valueT) bool) goption.O[keyT] {
	var (
		found keyT
		ok    bool
	)
	s.Range(func(key keyT, value valueT) bool {
		if f(key, value) {
			found, ok = key, true
			return false
		}
		return true
	})
	return goption.Of(found, ok)
}

// ForEach calls function f for each item of the map.
func (s *OrderedMapDesc[keyT, valueT]) ForEach(f func(keyT, valueT)) {
	s.Range(func(key keyT, value valueT) bool {
		f(key, value)
		return true
	})
}
