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

package gslice

import (
	"math"
	"reflect"
	"testing"
)

type wideUint uint64
type wideInt int64

func TestWideIntegerIndexesDoNotWrap(t *testing.T) {
	s := []int{10, 20, 30}
	maxUint := uint64(math.MaxUint64)
	minInt := int64(math.MinInt64)

	if got := Get(s, maxUint); !got.IsNil() {
		t.Fatalf("Get(MaxUint64) = %v, want nil", got)
	}
	if got := Get(s, minInt); !got.IsNil() {
		t.Fatalf("Get(MinInt64) = %v, want nil", got)
	}

	if got := Take(s, maxUint); !reflect.DeepEqual(got, s) {
		t.Fatalf("Take(MaxUint64) = %v, want %v", got, s)
	}
	if got := Take(s, minInt); !reflect.DeepEqual(got, s) {
		t.Fatalf("Take(MinInt64) = %v, want %v", got, s)
	}

	if got := Slice(s, uint64(0), maxUint); !reflect.DeepEqual(got, s) {
		t.Fatalf("Slice(0, MaxUint64) = %v, want %v", got, s)
	}
	if got := Slice(s, minInt, int64(math.MaxInt64)); !reflect.DeepEqual(got, s) {
		t.Fatalf("Slice(MinInt64, MaxInt64) = %v, want %v", got, s)
	}

	if got := Insert(s, maxUint, 99); !reflect.DeepEqual(got, []int{10, 20, 30, 99}) {
		t.Fatalf("Insert(MaxUint64) = %v", got)
	}
	if got := Insert(s, minInt, 99); !reflect.DeepEqual(got, []int{99, 10, 20, 30}) {
		t.Fatalf("Insert(MinInt64) = %v", got)
	}

	if got := RemoveIndex(s, maxUint); !reflect.DeepEqual(got, s) {
		t.Fatalf("RemoveIndex(MaxUint64) = %v, want clone of %v", got, s)
	}
	if got := RemoveIndex(s, minInt); !reflect.DeepEqual(got, s) {
		t.Fatalf("RemoveIndex(MinInt64) = %v, want clone of %v", got, s)
	}
}

func TestNamedWideIntegerIndexesDoNotWrap(t *testing.T) {
	s := []int{10, 20, 30}
	maxUint := wideUint(math.MaxUint64)
	minInt := wideInt(math.MinInt64)

	if !Get(s, maxUint).IsNil() || !Get(s, minInt).IsNil() {
		t.Fatal("named out-of-range integer unexpectedly resolved to an element")
	}
	if got := Take(s, maxUint); !reflect.DeepEqual(got, s) {
		t.Fatalf("Take(named MaxUint64) = %v", got)
	}
	if got := Slice(s, wideUint(0), maxUint); !reflect.DeepEqual(got, s) {
		t.Fatalf("Slice(0, named MaxUint64) = %v", got)
	}
	if got := Insert(s, maxUint, 99); !reflect.DeepEqual(got, []int{10, 20, 30, 99}) {
		t.Fatalf("Insert(named MaxUint64) = %v", got)
	}
	if got := RemoveIndex(s, maxUint); !reflect.DeepEqual(got, s) {
		t.Fatalf("RemoveIndex(named MaxUint64) = %v", got)
	}
}
