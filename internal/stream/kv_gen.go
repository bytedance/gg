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

// code generated by go run gen.go; DO NOT EDIT.

package stream

import (
	"context"
	"github.com/bytedance/gg/collection/tuple"
	"github.com/bytedance/gg/internal/iter"
)

// KV is a tuple.T2[K, V] variant of Stream.
type KV[K comparable, V any] struct {
	Stream[tuple.T2[K, V]]
}

// FromIter wraps an [github.com/bytedance/gg/internal/iter.Iter] to [Stream].
func FromKVIter[K comparable, V any](i iter.Iter[tuple.T2[K, V]]) KV[K, V] {
	return KV[K, V]{FromIter(i)}
}

// See function [github.com/bytedance/gg/internal/iter.FromSlice].
func FromKVSlice[K comparable, V any](s []tuple.T2[K, V]) KV[K, V] {
	return KV[K, V]{FromSlice(s)}
}

// See function [github.com/bytedance/gg/internal/iter.StealSlice].
func StealKVSlice[K comparable, V any](s []tuple.T2[K, V]) KV[K, V] {
	return KV[K, V]{StealSlice(s)}
}

// See function [github.com/bytedance/gg/internal/iter.FromChan].
func FromKVChan[K comparable, V any](ctx context.Context, ch <-chan tuple.T2[K, V]) KV[K, V] {
	return KV[K, V]{FromChan(ctx, ch)}
}

// See function [github.com/bytedance/gg/internal/iter.FlatMap].
func (s KV[K, V]) FlatMap(f func(tuple.T2[K, V]) []tuple.T2[K, V]) KV[K, V] {
	return KV[K, V]{s.Stream.FlatMap(f)}
}

// See function [github.com/bytedance/gg/internal/iter.Reverse].
func (s KV[K, V]) Reverse() KV[K, V] {
	return KV[K, V]{s.Stream.Reverse()}
}

// See function [github.com/bytedance/gg/internal/iter.Take].
func (s KV[K, V]) Take(n int) KV[K, V] {
	return KV[K, V]{s.Stream.Take(n)}
}

// See function [github.com/bytedance/gg/internal/iter.Drop].
func (s KV[K, V]) Drop(n int) KV[K, V] {
	return KV[K, V]{s.Stream.Drop(n)}
}

// See function [github.com/bytedance/gg/internal/iter.Concat].
func (s KV[K, V]) Concat(ss ...KV[K, V]) KV[K, V] {
	conv := func(c KV[K, V]) Stream[tuple.T2[K, V]] {
		return c.Stream
	}
	tmp := iter.ToSlice(iter.Map(conv, iter.FromSlice(ss)))
	return KV[K, V]{s.Stream.Concat(tmp...)}
}

// See function [github.com/bytedance/gg/internal/iter.Shuffle].
func (s KV[K, V]) Shuffle() KV[K, V] {
	return KV[K, V]{s.Stream.Shuffle()}
}
