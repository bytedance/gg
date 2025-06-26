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

// Package set provides a collection that contains no duplicate elements for comparable type.
//
// 💡 NOTE: Set is not concurrent-safe.
// If you need a high-performance, scalable, concurrent-safe set,
// use [github.com/bytedance/gg/collection/skipset].
//
// # Structures
//
//   - [Set]
//
// # Operations
//
//   - Constructor: [New], …
//   - CRUD operations: [set.Set.Add], [set.Set.Remove], [set.Set.Contains], [set.Set.ContainsAll], [set.Set.ContainsAny], …
//   - Set operations: [set.Set.Union], [set.Set.Intersect], [set.Set.Diff], … and its variants [set.Set.UnionInplace], [set.Set.IntersectInplace], …
//   - Predicates: [set.Set.Equal], [set.Set.IsSubset], [set.Set.IsSuperset], …
//   - Conversion: [set.Set.String], [set.Set.ToSlice], …
//   - Range operations: [set.Set.Range], [set.Set.Iter], …
//
// # JSON
//
// [set.Set] implements [encoding/json.Marshaler] and [encoding/json.Unmarshaler], so
// you can use it in JSON marshaling/unmarshaling.
// See [set.Set.MarshalJSON] and [set.Set.UnmarshalJSON].
//
// # Unspecified iteration order
//
// As [set.Set.Range] said, the iteration order over sets is not specified.
//
// If you need fixed order iteration, you can:
//
//   - Use [set.Set.ToSlice] and use [github.com/bytedance/gg/gslice.Sort] before iteration
package set

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bytedance/gg/gcond"
	"github.com/bytedance/gg/gvalue"
	"github.com/bytedance/gg/internal/heapsort"
	"github.com/bytedance/gg/internal/jsonbuilder"
)

const (
	initSize = 32
)

// Set is a set for comparable type.
//
// 💡 NOTE: Set is not concurrent-safe.
//
//	If you need a high-performance, scalable, concurrent-safe set,
//	use [github.com/bytedance/gg/collection/skipset].
type Set[T comparable] struct {
	m map[T]struct{} // Internal map for storing members.
}

// New creates a new set with initial members.
func New[T comparable](members ...T) *Set[T] {
	s := &Set[T]{}

	s.m = make(map[T]struct{}, gcond.If(len(members) == 0, initSize, len(members)))
	for _, v := range members {
		s.m[v] = struct{}{}
	}

	return s
}

// NewWithCap creates a new set with capacity.
func NewWithCap[T comparable](capacity int) *Set[T] {
	s := &Set[T]{}
	s.m = make(map[T]struct{}, capacity)
	return s
}

// Len returns the number of elements of set s.
// The complexity is O(1).
func (s *Set[T]) Len() int {
	if s == nil {
		return 0
	}
	return len(s.m)
}

// lazyInit lazily initializes a zero Set value.
func (s *Set[T]) lazyInit() {
	if s.m == nil {
		s.m = make(map[T]struct{}, initSize)
	}
}

// Add adds element v to set.
// If element is already member of set, return false.
func (s *Set[T]) Add(v T) bool {
	s.lazyInit()
	if _, ok := s.m[v]; ok {
		return false
	}
	s.m[v] = struct{}{}
	return true
}

// AddN is a variant of [set.Set.Add], adds multiple elements to set.
// It will not tell you which elements have been successfully added.
func (s *Set[T]) AddN(vs ...T) {
	s.lazyInit()
	for i := range vs {
		s.m[vs[i]] = struct{}{}
	}
}

// Remove removes element v from set.
// If element is not member of set, return false.
func (s *Set[T]) Remove(v T) bool {
	if s == nil {
		return false
	}
	_, ok := s.m[v]
	if ok {
		delete(s.m, v)
	}
	return ok
}

// RemoveN is a variant of [set.Set.Remove], removes multiple elements from set.
// It will not tell you which elements have been successfully removed.
func (s *Set[T]) RemoveN(vs ...T) {
	if s == nil {
		return
	}
	for i := range vs {
		delete(s.m, vs[i])
	}
}

// Contains returns true if element v is member of set.
func (s *Set[T]) Contains(v T) bool {
	if s == nil {
		return false
	}
	_, ok := s.m[v]
	return ok
}

// ContainsAny returns true if one of elements is member of set.
//
// 💡 NOTE: If no element given, ContainsAny always return false.
func (s *Set[T]) ContainsAny(vs ...T) bool {
	if s == nil {
		return false
	}
	for _, v := range vs {
		if _, ok := s.m[v]; ok {
			return true
		}
	}
	return false
}

// ContainsAll returns true if all elements are member of set.
//
// 💡 NOTE: If no element given, ContainsAll always return true.
func (s *Set[T]) ContainsAll(vs ...T) bool {
	if s == nil && len(vs) > 0 {
		return false
	}
	for _, v := range vs {
		if _, ok := s.m[v]; !ok {
			return false
		}
	}
	return true
}

// Range calls f sequentially for each member in the set.
// If f returns false, range stops the iteration.
//
// 💡 NOTE: The iteration order over sets is not specified and is not guaranteed
// to be the same from one iteration to the next.
func (s *Set[T]) Range(f func(T) bool) {
	if s == nil {
		return
	}
	for v := range s.m {
		if !f(v) {
			return
		}
	}
}

func (s *Set[T]) forEach(f func(T)) {
	if s == nil {
		return
	}
	for v := range s.m {
		f(v)
	}
}

// Union returns the unions of sets as a new set.
func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	res := NewWithCap[T](s.Len() + other.Len())
	s.forEach(func(v T) {
		res.m[v] = struct{}{}
	})
	other.forEach(func(v T) {
		res.m[v] = struct{}{}
	})
	return res
}

// Diff returns the difference of sets as a new set.
func (s *Set[T]) Diff(other *Set[T]) *Set[T] {
	res := NewWithCap[T](s.Len())
	s.forEach(func(v T) {
		if !other.Contains(v) {
			res.m[v] = struct{}{}
		}
	})
	return res
}

// Intersect returns the intersection of sets as a new set.
func (s *Set[T]) Intersect(other *Set[T]) *Set[T] {
	res := NewWithCap[T](gvalue.Min(s.Len(), other.Len()))
	s.forEach(func(v T) {
		if other.Contains(v) {
			res.m[v] = struct{}{}
		}
	})
	return res
}

// Update is alias of [Set.UnionInplace].
func (s *Set[T]) Update(other *Set[T]) {
	s.UnionInplace(other)
}

// UnionInplace updates set s with union itself and set other.
func (s *Set[T]) UnionInplace(other *Set[T]) {
	s.lazyInit()
	other.forEach(func(v T) {
		s.m[v] = struct{}{}
	})
}

// DiffInplace removes all elements of set other from set s.
func (s *Set[T]) DiffInplace(other *Set[T]) {
	s.forEach(func(v T) {
		if other.Contains(v) {
			delete(s.m, v)
		}
	})
}

// IntersectInplace updates set s with the intersection of itself and set other.
func (s *Set[T]) IntersectInplace(other *Set[T]) {
	s.forEach(func(v T) {
		if !other.Contains(v) {
			delete(s.m, v)
		}
	})
}

// Equal returns whether set s and other are equal.
func (s *Set[T]) Equal(other *Set[T]) bool {
	if s.Len() != other.Len() {
		return false
	}
	if s.Len() == 0 {
		return true
	}
	for v := range s.m {
		_, ok := other.m[v]
		if !ok {
			return false
		}
	}
	return true
}

// IsSubset returns whether another set contains this set.
func (s *Set[T]) IsSubset(other *Set[T]) bool {
	if s.Len() == 0 {
		return true
	}
	if s.Len() > other.Len() {
		return false
	}
	for v := range s.m {
		_, ok := other.m[v]
		if !ok {
			return false
		}
	}
	return true
}

// IsSuperset returns whether this set contains another set.
func (s *Set[T]) IsSuperset(other *Set[T]) bool {
	if other.Len() == 0 {
		return true
	}
	if s.Len() < other.Len() {
		return false
	}
	for v := range other.m {
		_, ok := s.m[v]
		if !ok {
			return false
		}
	}
	return true
}

// String implements [fmt.Stringer].
//
// Experimental: This API is experimental and may change in the future.
func (s *Set[T]) String() string {
	if s == nil {
		return "set[]"
	}
	members := make([]string, 0, s.Len())
	for m := range s.m {
		members = append(members, fmt.Sprintf("%v", m))
	}
	heapsort.Sort(members)
	return fmt.Sprintf("set[%s]", strings.Join(members, " "))
}

// ToSlice collects all members to slice.
//
// 💡 NOTE: The order of returned slice is not specified and is not guaranteed
// to be the same from another ToSlice call.
func (s *Set[T]) ToSlice() []T {
	members := make([]T, 0, s.Len())
	s.forEach(func(v T) {
		members = append(members, v)
	})
	return members
}

// MarshalJSON implements [encoding/json.Marshaler].
//
// NOTE: The returned bytes is null or JSON array. Elements of array are
// sorted lexicographically.
//
// Experimental: This API is experimental and may change in the future.
func (s *Set[T]) MarshalJSON() ([]byte, error) {
	if s == nil {
		return []byte("null"), nil
	}
	b := jsonbuilder.NewArray()
	for m := range s.m {
		if err := b.Append(m); err != nil {
			return nil, err
		}
	}
	b.Sort()
	return b.Build()
}

// UnmarshalJSON implements [encoding/json.Unmarshaler].
//
// Experimental: This API is experimental and may change in the future.
func (s *Set[T]) UnmarshalJSON(data []byte) error {
	// Unmarshalers implement UnmarshalJSON([]byte("null")) as a no-op.
	if string(data) == "null" {
		return nil
	}

	var members []T
	if err := json.Unmarshal(data, &members); err != nil {
		return err
	}
	// Always override original members.
	*s = *New(members...)
	return nil
}

// Clone returns a copy of the set.
//
// 💡 NOTE: Members are copied using assignment (=).
func (s *Set[T]) Clone() *Set[T] {
	ns := NewWithCap[T](s.Len())
	s.forEach(func(v T) {
		ns.m[v] = struct{}{}
	})
	return ns
}
