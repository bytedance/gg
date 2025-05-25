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

package gatomic

import (
	"sync/atomic"

	"github.com/bytedance/gg/gvalue"
)

// Value wraps [sync/atomic.Value].
type Value[T any] struct {
	v atomic.Value
}

// Load wraps [sync/atomic.Value.Load].
func (av *Value[T]) Load() T {
	v := av.v.Load()
	if v == nil {
		return gvalue.Zero[T]()
	}
	return v.(T)
}

// Store wraps [sync/atomic.Value.Load].
func (av *Value[T]) Store(v T) {
	av.v.Store(v)
}

// Swap wraps [sync/atomic.Value.Swap].
func (av *Value[T]) Swap(new T) T {
	old := av.v.Swap(new)
	if old == nil {
		return gvalue.Zero[T]()
	}
	return old.(T)
}

// CompareAndSwap wraps [sync/atomic.Value.CompareAndSwap].
func (av *Value[T]) CompareAndSwap(old, new T) (swapped bool) {
	return av.v.CompareAndSwap(old, new)
}
