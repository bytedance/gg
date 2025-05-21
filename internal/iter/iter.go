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

// Package iter provides definition of generic iterator Iter and high-order functions.
//
// Please refer to README.md for details.
package iter

const (
	ALL = -1
)

// Iter is a generic iterator interface, which helps us iterate various
// data structures in the same way.
//
// Users can apply various operations ([Map], [Filter], etc.) on custom data
// structures by implementing Iter for them.
// See ExampleIter_impl for details.
type Iter[T any] interface {
	// Next returns N items of iterator when it is not empty.
	// When the iterator is empty, nil is returned.
	// When n = [ALL] or n is greater than the number of remaining elements,
	// all remaining are returned.
	//
	// The returned slice is owned by caller. So implementer should return a
	// newly allocated slice if needed.
	//
	// Passing in a negative n (except [ALL]) is undefined behavior.
	Next(n int) []T
}

// emptyIter returns nil whenever its Next method is called.
// It can be as a default abnormal behavior when implements Iter.
// For example, in RangeWithStep, if the internal does not exist, it will return emptyIter,
// so the returned Iter works normally in silence in the subsequent iterator chain.
type emptyIter[T any] struct{}

func (i emptyIter[T]) Next(_ int) []T {
	return nil
}
