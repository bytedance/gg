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

// Package skipmap is a high-performance, scalable, concurrent-safe map based on skip-list.
// In the typical pattern(100000 operations, 90%LOAD 9%STORE 1%DELETE, 8C16T), the skipmap
// up to 10x faster than the built-in sync.Map.
//
//go:generate go run gen.go
package skipmap

import "github.com/bytedance/gg/internal/constraints"

// NewFunc returns an empty skipmap in ascending order.
//
// Note that the less function requires a strict weak ordering,
// see https://en.wikipedia.org/wiki/Weak_ordering#Strict_weak_orderings,
// or undefined behavior will happen.
func NewFunc[keyT any, valueT any](less func(a, b keyT) bool) *FuncMap[keyT, valueT] {
	var (
		t1 keyT
		t2 valueT
	)
	h := newFuncNode(t1, t2, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &FuncMap[keyT, valueT]{
		header:       h,
		highestLevel: defaultHighestLevel,
		less:         less,
	}
}

// New returns an empty skipmap in ascending order.
func New[keyT constraints.Ordered, valueT any]() *OrderedMap[keyT, valueT] {
	var (
		t1 keyT
		t2 valueT
	)
	h := newOrderedNode(t1, t2, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &OrderedMap[keyT, valueT]{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewDesc returns an empty skipmap in descending order.
func NewDesc[keyT constraints.Ordered, valueT any]() *OrderedMapDesc[keyT, valueT] {
	var (
		t1 keyT
		t2 valueT
	)
	h := newOrderedNodeDesc(t1, t2, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &OrderedMapDesc[keyT, valueT]{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}
