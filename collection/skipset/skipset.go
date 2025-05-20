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

// Package skipset is a high-performance, scalable, concurrent-safe set based on skip-list.
// In the typical pattern(100000 operations, 90%CONTAINS 9%Add 1%Remove, 8C16T), the skipset
// up to 15x faster than the built-in sync.Map.
//
//go:generate go run gen.go
package skipset

import "github.com/bytedance/gg/internal/constraints"

// New returns an empty skip set in ascending order.
func New[T constraints.Ordered]() *OrderedSet[T] {
	var t T
	h := newOrderedNode(t, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &OrderedSet[T]{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewDesc returns an empty skip set in descending order.
func NewDesc[T constraints.Ordered]() *OrderedSetDesc[T] {
	var t T
	h := newOrderedNodeDesc(t, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &OrderedSetDesc[T]{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewFunc returns an empty skip set in ascending order.
//
// Note that the less function requires a strict weak ordering,
// see https://en.wikipedia.org/wiki/Weak_ordering#Strict_weak_orderings,
// or undefined behavior will happen.
func NewFunc[T any](less func(a, b T) bool) *FuncSet[T] {
	var t T
	h := newFuncNode(t, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &FuncSet[T]{
		header:       h,
		highestLevel: defaultHighestLevel,
		less:         less,
	}
}
