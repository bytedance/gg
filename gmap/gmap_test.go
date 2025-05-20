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

package gmap

import (
	"fmt"
	"sort"
	"strconv"
	"testing"

	"github.com/bytedance/gg/collection/tuple"
	"github.com/bytedance/gg/goption"
	"github.com/bytedance/gg/gptr"
	"github.com/bytedance/gg/gresult"
	"github.com/bytedance/gg/gslice"
	"github.com/bytedance/gg/gvalue"
	"github.com/bytedance/gg/internal/assert"
)

func TestMap(t *testing.T) {
	assert.Equal(t,
		map[string]string{"1": "1", "2": "2"},
		Map(map[int]int{1: 1, 2: 2}, func(k, v int) (string, string) {
			return strconv.Itoa(k), strconv.Itoa(v)
		}))
	assert.Equal(t,
		map[string]string{},
		Map(map[int]int{}, func(k, v int) (string, string) {
			return strconv.Itoa(k), strconv.Itoa(v)
		}))
}

func TestMapKeys(t *testing.T) {
	assert.Equal(t,
		map[string]int{"1": 1, "2": 2},
		MapKeys(map[int]int{1: 1, 2: 2}, strconv.Itoa))
	assert.Equal(t,
		map[string]int{},
		MapKeys(map[int]int{}, strconv.Itoa))
}

func TestTryMapKeys(t *testing.T) {
	assert.Equal(t,
		gresult.OK(map[int]int{}),
		TryMapKeys(map[string]int{}, strconv.Atoi))
	assert.Equal(t,
		gresult.OK(map[int]int{1: 1, 2: 2}),
		TryMapKeys(map[string]int{"1": 1, "2": 2}, strconv.Atoi))
	assert.Equal(t,
		"strconv.Atoi: parsing \"a\": invalid syntax",
		TryMapKeys(map[string]int{"1": 1, "a": 2}, strconv.Atoi).Err().Error())
}

func TestMapValues(t *testing.T) {
	assert.Equal(t,
		map[int]string{1: "1", 2: "2"},
		MapValues(map[int]int{1: 1, 2: 2}, strconv.Itoa))
	assert.Equal(t,
		map[int]string{},
		MapValues(map[int]int{}, strconv.Itoa))
}

func TestTryMapValues(t *testing.T) {
	assert.Equal(t,
		gresult.OK(map[int]int{}),
		TryMapValues(map[int]string{}, strconv.Atoi))
	assert.Equal(t,
		gresult.OK(map[int]int{1: 1, 2: 2}),
		TryMapValues(map[int]string{1: "1", 2: "2"}, strconv.Atoi))
	assert.Equal(t,
		"strconv.Atoi: parsing \"a\": invalid syntax",
		TryMapValues(map[int]string{1: "1", 2: "a"}, strconv.Atoi).Err().Error())
}

func TestFilter(t *testing.T) {
	assert.Equal(t,
		map[int]int{1: 1, 2: 2},
		Filter(map[int]int{1: 1, 2: 2, 3: 2, 4: 3}, func(k, v int) bool { return (k+v)%2 == 0 }))
	assert.Equal(t,
		map[int]int{},
		Filter(map[int]int{}, func(k, v int) bool { return (k+v)%2 == 0 }))
	assert.Equal(t,
		map[int]int{},
		Filter(map[int]int{1: 1, 2: 2, 3: 2, 4: 3}, func(k, v int) bool { return k+v > 100 }))
	assert.Equal(t,
		map[int]int{1: 1, 2: 2, 3: 2, 4: 3},
		Filter(map[int]int{1: 1, 2: 2, 3: 2, 4: 3}, func(k, v int) bool { return k+v > 0 }))
}

func TestFilterKeys(t *testing.T) {
	assert.Equal(t,
		map[int]int{2: 2, 4: 3},
		FilterKeys(map[int]int{1: 1, 2: 2, 3: 2, 4: 3}, func(k int) bool { return k%2 == 0 }))
	assert.Equal(t,
		map[int]int{},
		FilterKeys(map[int]int{}, func(k int) bool { return k%2 == 0 }))
	assert.Equal(t,
		map[int]int{},
		FilterKeys(map[int]int{1: 1, 2: 2, 3: 2, 4: 3}, func(k int) bool { return k > 100 }))
	assert.Equal(t,
		map[int]int{1: 1, 2: 2, 3: 2, 4: 3},
		FilterKeys(map[int]int{1: 1, 2: 2, 3: 2, 4: 3}, func(k int) bool { return k > 0 }))
}

func TestFilterByKeys(t *testing.T) {
	tests := []struct {
		name   string
		input  map[int]string
		keys   []int
		expect map[int]string
	}{
		{
			name:   "basic filtering",
			input:  map[int]string{1: "a", 2: "b", 3: "c", 4: "d"},
			keys:   []int{1, 3},
			expect: map[int]string{1: "a", 3: "c"},
		},
		{
			name:   "empty keys",
			input:  map[int]string{1: "a", 2: "b"},
			keys:   []int{},
			expect: map[int]string{},
		},
		{
			name:   "non-existent keys",
			input:  map[int]string{1: "a", 2: "b"},
			keys:   []int{3, 4},
			expect: map[int]string{},
		},
		{
			name:   "partially existing keys",
			input:  map[int]string{1: "a", 2: "b", 3: "c"},
			keys:   []int{1, 3, 5},
			expect: map[int]string{1: "a", 3: "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterByKeys(tt.input, tt.keys...)
			assert.Equal(t, tt.expect, got)
		})
	}
}

func TestFilterValues(t *testing.T) {
	assert.Equal(t,
		map[int]int{2: 2, 3: 2},
		FilterValues(map[int]int{1: 1, 2: 2, 3: 2, 4: 3}, func(v int) bool { return v%2 == 0 }))
	assert.Equal(t,
		map[int]int{},
		FilterValues(map[int]int{}, func(v int) bool { return v%2 == 0 }))
	assert.Equal(t,
		map[int]int{},
		FilterValues(map[int]int{1: 1, 2: 2, 3: 2, 4: 3}, func(v int) bool { return v > 100 }))
	assert.Equal(t,
		map[int]int{1: 1, 2: 2, 3: 2, 4: 3},
		FilterValues(map[int]int{1: 1, 2: 2, 3: 2, 4: 3}, func(v int) bool { return v > 0 }))
}

func TestFilterByValues(t *testing.T) {
	tests := []struct {
		name   string
		input  map[string]int
		values []int
		expect map[string]int
	}{
		{
			name:   "basic filtering",
			input:  map[string]int{"a": 1, "b": 2, "c": 1, "d": 3},
			values: []int{1, 3},
			expect: map[string]int{"a": 1, "c": 1, "d": 3},
		},
		{
			name:   "empty values",
			input:  map[string]int{"a": 1, "b": 2},
			values: []int{},
			expect: map[string]int{},
		},
		{
			name:   "non-existent values",
			input:  map[string]int{"a": 1, "b": 2},
			values: []int{3, 4},
			expect: map[string]int{},
		},
		{
			name:   "duplicate values",
			input:  map[string]int{"a": 1, "b": 2, "c": 1, "d": 1},
			values: []int{1},
			expect: map[string]int{"a": 1, "c": 1, "d": 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterByValues(tt.input, tt.values...)
			assert.Equal(t, tt.expect, got)
		})
	}
}

func TestReject(t *testing.T) {
	assert.Equal(t,
		map[int]int{1: 1, 2: 2},
		Reject(map[int]int{1: 1, 2: 2, 3: 2, 4: 3}, func(k, v int) bool { return (k+v)%2 != 0 }))
	assert.Equal(t,
		map[int]int{},
		Reject(map[int]int{}, func(k, v int) bool { return (k+v)%2 != 0 }))
	assert.Equal(t,
		map[int]int{},
		Reject(map[int]int{1: 1, 2: 2, 3: 2, 4: 3}, func(k, v int) bool { return k+v < 100 }))
	assert.Equal(t,
		map[int]int{1: 1, 2: 2, 3: 2, 4: 3},
		Reject(map[int]int{1: 1, 2: 2, 3: 2, 4: 3}, func(k, v int) bool { return k+v < 0 }))
}

func TestRejectKeys(t *testing.T) {
	assert.Equal(t,
		map[int]int{2: 2, 4: 3},
		RejectKeys(map[int]int{1: 1, 2: 2, 3: 2, 4: 3}, func(k int) bool { return k%2 != 0 }))
	assert.Equal(t,
		map[int]int{},
		RejectKeys(map[int]int{}, func(k int) bool { return k%2 != 0 }))
	assert.Equal(t,
		map[int]int{},
		RejectKeys(map[int]int{1: 1, 2: 2, 3: 2, 4: 3}, func(k int) bool { return k < 100 }))
	assert.Equal(t,
		map[int]int{1: 1, 2: 2, 3: 2, 4: 3},
		RejectKeys(map[int]int{1: 1, 2: 2, 3: 2, 4: 3}, func(k int) bool { return k < 0 }))
}

func TestRejectByKeys(t *testing.T) {
	tests := []struct {
		name   string
		input  map[int]string
		keys   []int
		expect map[int]string
	}{
		{
			name:   "basic rejection",
			input:  map[int]string{1: "a", 2: "b", 3: "c", 4: "d"},
			keys:   []int{1, 3},
			expect: map[int]string{2: "b", 4: "d"},
		},
		{
			name:   "empty keys",
			input:  map[int]string{1: "a", 2: "b"},
			keys:   []int{},
			expect: map[int]string{1: "a", 2: "b"},
		},
		{
			name:   "non-existent keys",
			input:  map[int]string{1: "a", 2: "b"},
			keys:   []int{3, 4},
			expect: map[int]string{1: "a", 2: "b"},
		},
		{
			name:   "partially existing keys",
			input:  map[int]string{1: "a", 2: "b", 3: "c"},
			keys:   []int{1, 3, 5},
			expect: map[int]string{2: "b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RejectByKeys(tt.input, tt.keys...)
			assert.Equal(t, tt.expect, got)
		})
	}
}

func TestRejectValues(t *testing.T) {
	assert.Equal(t,
		map[int]int{2: 2, 3: 2},
		RejectValues(map[int]int{1: 1, 2: 2, 3: 2, 4: 3}, func(v int) bool { return v%2 != 0 }))
	assert.Equal(t,
		map[int]int{},
		RejectValues(map[int]int{}, func(v int) bool { return v%2 == 0 }))
	assert.Equal(t,
		map[int]int{},
		RejectValues(map[int]int{1: 1, 2: 2, 3: 2, 4: 3}, func(v int) bool { return v < 100 }))
	assert.Equal(t,
		map[int]int{1: 1, 2: 2, 3: 2, 4: 3},
		RejectValues(map[int]int{1: 1, 2: 2, 3: 2, 4: 3}, func(v int) bool { return v < 0 }))
}

func TestRejectByValues(t *testing.T) {
	tests := []struct {
		name   string
		input  map[string]int
		values []int
		expect map[string]int
	}{
		{
			name:   "basic rejection",
			input:  map[string]int{"a": 1, "b": 2, "c": 1, "d": 3},
			values: []int{1, 3},
			expect: map[string]int{"b": 2},
		},
		{
			name:   "empty values",
			input:  map[string]int{"a": 1, "b": 2},
			values: []int{},
			expect: map[string]int{"a": 1, "b": 2},
		},
		{
			name:   "non-existent values",
			input:  map[string]int{"a": 1, "b": 2},
			values: []int{3, 4},
			expect: map[string]int{"a": 1, "b": 2},
		},
		{
			name:   "duplicate values in input",
			input:  map[string]int{"a": 1, "b": 2, "c": 1, "d": 1},
			values: []int{1},
			expect: map[string]int{"b": 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RejectByValues(tt.input, tt.values...)
			assert.Equal(t, tt.expect, got)
		})
	}
}

func TestFold(t *testing.T) {
	assert.Equal(t,
		6,
		fold(map[int]int{1: 1, 2: 2}, func(acc, k, v int) int { return acc + k + v }, 0))
	assert.Equal(t,
		9,
		fold(map[int]int{1: 1, 2: 2}, func(acc, k, v int) int { return acc + k + v }, 3))
	assert.Equal(t,
		0,
		fold(map[int]int{}, func(acc, k, v int) int { return acc + k + v }, 0))
	assert.Equal(t,
		3,
		fold(map[int]int{}, func(acc, k, v int) int { return acc + k + v }, 3))
}

func TestFoldKeys(t *testing.T) {
	assert.Equal(t,
		3,
		foldKeys(map[int]int{1: 2, 2: 4}, gvalue.Add[int], 0))
	assert.Equal(t,
		5,
		foldKeys(map[int]int{1: 2, 2: 4}, gvalue.Add[int], 2))
	assert.Equal(t,
		0,
		foldKeys(map[int]int{}, gvalue.Add[int], 0))
	assert.Equal(t,
		2,
		foldKeys(map[int]int{}, gvalue.Add[int], 2))
}

func TestFoldValues(t *testing.T) {
	assert.Equal(t,
		6,
		foldValues(map[int]int{1: 2, 2: 4}, gvalue.Add[int], 0))
	assert.Equal(t,
		8,
		foldValues(map[int]int{1: 2, 2: 4}, gvalue.Add[int], 2))
	assert.Equal(t,
		0,
		foldValues(map[int]int{}, gvalue.Add[int], 0))
	assert.Equal(t,
		2,
		foldValues(map[int]int{}, gvalue.Add[int], 2))
}

func TestReduceKeys(t *testing.T) {
	assert.Equal(t,
		goption.OK(3),
		reduceKeys(map[int]int{1: 2, 2: 4}, gvalue.Add[int]))
	assert.Equal(t,
		goption.Nil[int](),
		reduceKeys(map[int]int{}, gvalue.Add[int]))
}

func TestReduceValues(t *testing.T) {
	assert.Equal(t,
		goption.OK(6),
		reduceValues(map[int]int{1: 2, 2: 4}, gvalue.Add[int]))
	assert.Equal(t,
		goption.Nil[int](),
		reduceValues(map[int]int{}, gvalue.Add[int]))
}

func TestKeys(t *testing.T) {
	{
		keys := Keys(map[int]string{1: "1", 2: "2", 3: "3", 4: "4"})
		gslice.Sort(keys)
		assert.Equal(t, []int{1, 2, 3, 4}, keys)
	}
	assert.Equal(t, []int{}, Keys(map[int]string{}))
	assert.Equal(t, []int{}, Keys[int, string](nil))
}

func TestValues(t *testing.T) {
	{
		keys := Values(map[int]string{1: "1", 2: "2", 3: "3", 4: "4"})
		gslice.Sort(keys)
		assert.Equal(t, []string{"1", "2", "3", "4"}, keys)
	}
	assert.Equal(t, []string{}, Values(map[int]string{}))
	assert.Equal(t, []string{}, Values[int, string](nil))
}

func TestOrderedKeys(t *testing.T) {
	assert.Equal(t, []int{1, 2, 3, 4}, OrderedKeys(map[int]string{1: "1", 2: "2", 3: "3", 4: "4"}))
	assert.Equal(t, []int{}, OrderedKeys(map[int]string{}))
	assert.Equal(t, []int{}, OrderedKeys[int, string](nil))
}

func TestOrderedValues(t *testing.T) {
	assert.Equal(t, []string{"1", "2", "3", "4"}, OrderedValues(map[int]string{1: "1", 2: "2", 3: "3", 4: "4"}))
	assert.Equal(t, []string{}, OrderedValues(map[int]string{}))
	assert.Equal(t, []string{}, OrderedValues[int, string](nil))
}

func TestItems(t *testing.T) {
	{
		items := Items(map[int]string{1: "1", 2: "2", 3: "3", 4: "4"})
		sort.Slice(items, func(i, j int) bool {
			return items[i].First < items[j].First
		})
		assert.Equal(t, tuple.S2[int, string]{{1, "1"}, {2, "2"}, {3, "3"}, {4, "4"}}, items)
	}
	assert.Equal(t, tuple.S2[int, string]{}, Items(map[int]string{}))
	assert.Equal(t, tuple.S2[int, string]{}, Items[int, string](nil))
}

func TestOrderedItems(t *testing.T) {
	{
		items := OrderedItems(map[int]string{1: "1", 2: "2", 3: "3", 4: "4"})
		assert.Equal(t, tuple.S2[int, string]{{1, "1"}, {2, "2"}, {3, "3"}, {4, "4"}}, items)
	}
	assert.Equal(t, tuple.S2[int, string]{}, OrderedItems(map[int]string{}))
	assert.Equal(t, tuple.S2[int, string]{}, OrderedItems[int, string](nil))
}

func TestMerge(t *testing.T) {
	assert.Equal(t, map[int]int{1: 1, 2: 2, 3: 3, 4: 4},
		Merge(map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, nil))
	assert.Equal(t, map[int]int{1: 1, 2: 2, 3: 3, 4: 4},
		Merge(nil, map[int]int{1: 1, 2: 2, 3: 3, 4: 4}))
	assert.Equal(t, map[int]int{}, Merge[int, int](nil, nil))
	assert.Equal(t, map[int]int{1: 1, 2: 2, 3: 3, 4: 4},
		Merge(map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, map[int]int{1: 1, 2: 2, 3: 3, 4: 4}))
	assert.Equal(t, map[int]int{1: 1, 2: 2, 3: 3, 4: 4},
		Merge(map[int]int{1: 0, 2: 0}, map[int]int{1: 1, 2: 2, 3: 3, 4: 4}))
	assert.Equal(t, map[int]int{1: 1, 2: 2, 3: 3, 4: 4},
		Merge(map[int]int{1: 1, 2: 1}, map[int]int{2: 2, 3: 3, 4: 4}))
}

func TestLoad(t *testing.T) {
	assert.Equal(t, goption.Nil[int](), Load[int, int](nil, 1))
	assert.Equal(t, goption.OK(1),
		Load(map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, 1))
	assert.Equal(t, goption.OK(2),
		Load(map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, 2))
	assert.Equal(t, goption.Nil[int](),
		Load(map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, 5))
}

func TestLoadOrStore(t *testing.T) {
	assert.Panic(t, func() {
		_, _ = LoadOrStore(nil, 1, "1")
	})
	{
		m := map[int]int{1: 1, 2: 2, 3: 3, 4: 4}
		assert.Equal(t, tuple.Make2(1, true),
			tuple.Make2(LoadOrStore(m, 1, 100)))
		assert.Equal(t, map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, m)
	}
	{
		m := map[int]int{1: 1, 2: 2, 3: 3, 4: 4}
		assert.Equal(t, tuple.Make2(2, true),
			tuple.Make2(LoadOrStore(m, 2, 100)))
		assert.Equal(t, map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, m)
	}
	{
		m := map[int]int{1: 1, 2: 2, 3: 3, 4: 4}
		assert.Equal(t, tuple.Make2(100, false),
			tuple.Make2(LoadOrStore(m, 5, 100)))
		assert.Equal(t, map[int]int{1: 1, 2: 2, 3: 3, 4: 4, 5: 100}, m)
	}
}

func TestLoadOrStoreLazy(t *testing.T) {
	assert.Panic(t, func() {
		_, _ = LoadOrStoreLazy(nil, 1, func() string { return "1" })
	})

	lazy100 := func() int { return 100 }

	{
		m := map[int]int{1: 1, 2: 2, 3: 3, 4: 4}
		assert.Equal(t, tuple.Make2(1, true),
			tuple.Make2(LoadOrStoreLazy(m, 1, lazy100)))
		assert.Equal(t, map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, m)
	}
	{
		m := map[int]int{1: 1, 2: 2, 3: 3, 4: 4}
		assert.Equal(t, tuple.Make2(2, true),
			tuple.Make2(LoadOrStoreLazy(m, 2, lazy100)))
		assert.Equal(t, map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, m)
	}
	{
		m := map[int]int{1: 1, 2: 2, 3: 3, 4: 4}
		assert.Equal(t, tuple.Make2(100, false),
			tuple.Make2(LoadOrStoreLazy(m, 5, lazy100)))
		assert.Equal(t, map[int]int{1: 1, 2: 2, 3: 3, 4: 4, 5: 100}, m)
	}
}

func TestLoadAndDelete(t *testing.T) {
	{
		assert.Equal(t, goption.Nil[int](), LoadAndDelete[int, int](nil, 1))
	}
	{
		m := map[int]int{1: 1, 2: 2, 3: 3, 4: 4}
		assert.Equal(t, goption.OK(1), LoadAndDelete(m, 1))
		assert.Equal(t, map[int]int{2: 2, 3: 3, 4: 4}, m)
	}
	{
		m := map[int]int{1: 1, 2: 2, 3: 3, 4: 4}
		assert.Equal(t, goption.OK(2), LoadAndDelete(m, 2))
		assert.Equal(t, map[int]int{1: 1, 3: 3, 4: 4}, m)
	}
	{
		m := map[int]int{1: 1, 2: 2, 3: 3, 4: 4}
		assert.Equal(t, goption.Nil[int](), LoadAndDelete(m, 5))
		assert.Equal(t, map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, m)
	}
}

func TestEqual(t *testing.T) {
	assert.True(t, Equal(
		map[int]int{1: 1, 2: 2, 3: 3, 4: 4},
		map[int]int{1: 1, 2: 2, 3: 3, 4: 4}))
	assert.False(t, Equal(
		map[int]int{1: 1, 2: 2, 3: 3},
		map[int]int{1: 1, 2: 2, 3: 3, 4: 4}))
	assert.False(t, Equal(
		map[int]int{1: 1, 2: 2, 3: 3, 4: 4},
		map[int]int{1: 1, 2: 2, 3: 3}))
	assert.False(t, Equal(
		map[int]int{1: 1, 2: 2, 3: 3, 4: 4},
		map[int]int{1: 1, 2: 2, 3: 3, 4: 5}))
	assert.False(t, Equal(map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, nil))
	assert.False(t, Equal(nil, map[int]int{1: 1, 2: 2, 3: 3, 4: 4}))
	assert.True(t, Equal(map[int]int{}, map[int]int{}))
	assert.True(t, Equal(nil, map[int]int{}))
	assert.True(t, Equal(map[int]int{}, nil))
	assert.True(t, Equal[int, int](nil, nil))
}

func TestEqualStrict(t *testing.T) {
	assert.True(t, EqualStrict(
		map[int]int{1: 1, 2: 2, 3: 3, 4: 4},
		map[int]int{1: 1, 2: 2, 3: 3, 4: 4}))
	assert.False(t, EqualStrict(
		map[int]int{1: 1, 2: 2, 3: 3},
		map[int]int{1: 1, 2: 2, 3: 3, 4: 4}))
	assert.False(t, EqualStrict(
		map[int]int{1: 1, 2: 2, 3: 3, 4: 4},
		map[int]int{1: 1, 2: 2, 3: 3}))
	assert.False(t, EqualStrict(
		map[int]int{1: 1, 2: 2, 3: 3, 4: 4},
		map[int]int{1: 1, 2: 2, 3: 3, 4: 5}))
	assert.False(t, EqualStrict(map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, nil))
	assert.False(t, EqualStrict(nil, map[int]int{1: 1, 2: 2, 3: 3, 4: 4}))
	assert.True(t, EqualStrict(map[int]int{}, map[int]int{}))
	assert.False(t, EqualStrict(nil, map[int]int{}))
	assert.False(t, EqualStrict(map[int]int{}, nil))
	assert.True(t, EqualStrict[int, int](nil, nil))
}

func TestEqualBy(t *testing.T) {
	eq := gvalue.Equal[int]
	assert.True(t, EqualBy(
		map[int]int{1: 1, 2: 2, 3: 3, 4: 4},
		map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, eq))
	assert.False(t, EqualBy(
		map[int]int{1: 1, 2: 2, 3: 3},
		map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, eq))
	assert.False(t, EqualBy(
		map[int]int{1: 1, 2: 2, 3: 3, 4: 4},
		map[int]int{1: 1, 2: 2, 3: 3}, eq))
	assert.False(t, EqualBy(
		map[int]int{1: 1, 2: 2, 3: 3, 4: 4},
		map[int]int{1: 1, 2: 2, 3: 3, 4: 5}, eq))
	assert.False(t, EqualBy(map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, nil, eq))
	assert.False(t, EqualBy(nil, map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, eq))
	assert.True(t, EqualBy(map[int]int{}, map[int]int{}, eq))
	assert.True(t, EqualBy(nil, map[int]int{}, eq))
	assert.True(t, EqualBy(map[int]int{}, nil, eq))
	assert.True(t, EqualBy[int](nil, nil, eq))

	anyEq := func(v1, v2 any) bool { return v1 == v2 }
	assert.True(t, EqualBy(
		map[int]any{1: 1, 2: 2, 3: 3, 4: 4},
		map[int]any{1: 1, 2: 2, 3: 3, 4: 4}, anyEq))
}

func TestEqualStrictBy(t *testing.T) {
	eq := gvalue.Equal[int]
	assert.True(t, EqualStrictBy(
		map[int]int{1: 1, 2: 2, 3: 3, 4: 4},
		map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, eq))
	assert.False(t, EqualStrictBy(
		map[int]int{1: 1, 2: 2, 3: 3},
		map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, eq))
	assert.False(t, EqualStrictBy(
		map[int]int{1: 1, 2: 2, 3: 3, 4: 4},
		map[int]int{1: 1, 2: 2, 3: 3}, eq))
	assert.False(t, EqualStrictBy(
		map[int]int{1: 1, 2: 2, 3: 3, 4: 4},
		map[int]int{1: 1, 2: 2, 3: 3, 4: 5}, eq))
	assert.False(t, EqualStrictBy(map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, nil, eq))
	assert.False(t, EqualStrictBy(nil, map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, eq))
	assert.True(t, EqualStrictBy(map[int]int{}, map[int]int{}, eq))
	assert.False(t, EqualStrictBy(nil, map[int]int{}, eq))
	assert.False(t, EqualStrictBy(map[int]int{}, nil, eq))
	assert.True(t, EqualStrictBy[int](nil, nil, eq))

	anyEq := func(v1, v2 any) bool { return v1 == v2 }
	assert.True(t, EqualStrictBy(
		map[int]any{1: 1, 2: 2, 3: 3, 4: 4},
		map[int]any{1: 1, 2: 2, 3: 3, 4: 4}, anyEq))
}

func TestClone(t *testing.T) {
	assert.Equal(t, map[int]int{1: 1, 2: 2}, Clone(map[int]int{1: 1, 2: 2}))
	var nilMap map[int]int
	assert.Equal(t, map[int]int{}, Clone(map[int]int{}))
	assert.NotEqual(t, nil, Clone(map[int]int{}))
	assert.Equal(t, nil, Clone(nilMap))
	assert.NotEqual(t, map[int]int{}, Clone(nilMap))

	// Test new type.
	type I2I map[int]int
	assert.Equal(t, I2I{1: 1, 2: 2}, Clone(I2I{1: 1, 2: 2}))
	assert.Equal(t, "gmap.I2I", fmt.Sprintf("%T", Clone(I2I{})))

	// Test shallow clone.
	src := map[int]*int{1: gptr.Of(1), 2: gptr.Of(2)}
	dst := Clone(src)
	assert.Equal(t, src, dst)
	assert.True(t, src[1] == dst[1])
	assert.True(t, src[2] == dst[2])
}

func TestCloneBy(t *testing.T) {
	id := func(v int) int { return v }

	assert.Equal(t, map[int]int{1: 1, 2: 2}, CloneBy(map[int]int{1: 1, 2: 2}, id))
	var nilMap map[int]int
	assert.Equal(t, map[int]int{}, CloneBy(map[int]int{}, id))
	assert.NotEqual(t, nil, CloneBy(map[int]int{}, id))
	assert.Equal(t, nil, CloneBy(nilMap, id))
	assert.NotEqual(t, map[int]int{}, CloneBy(nilMap, id))

	// Test deep type.
	type I2I map[int]int
	assert.Equal(t, I2I{1: 1, 2: 2}, CloneBy(I2I{1: 1, 2: 2}, id))
	assert.Equal(t, "gmap.I2I", fmt.Sprintf("%T", CloneBy(I2I{}, id)))

	// Test deep clone.
	src := map[int]*int{1: gptr.Of(1), 2: gptr.Of(2)}
	dst := CloneBy(src, gptr.Clone[int])
	assert.Equal(t, src, dst)
	assert.False(t, src[1] == dst[1])
	assert.False(t, src[2] == dst[2])
}

func TestInvert(t *testing.T) {
	assert.Equal(t, map[int]string{}, Invert(map[string]int{}))
	assert.Equal(t, map[int]string{1: "1", 2: "2"}, Invert(map[string]int{"1": 1, "2": 2}))

	// Test custom type.
	type X struct{ Foo int }
	type Y struct{ Bar int }

	assert.Equal(t, map[Y]X{{Bar: 2}: {Foo: 1}}, Invert(map[X]Y{{Foo: 1}: {Bar: 2}}))
}

func TestInvertBy(t *testing.T) {
	assert.Equal(t, map[int]string{}, InvertBy(map[string]int{}, nil))
	assert.Equal(t, map[int]string{1: "1", 2: "2"}, InvertBy(map[string]int{"1": 1, "2": 2}, nil))
	assert.Equal(t, map[int]string{1: "1"},
		InvertBy(map[string]int{"1": 1, "": 1}, DiscardZero(DiscardOld[int, string]())))
}

func TestInvertGroup(t *testing.T) {
	assert.Equal(t, map[int][]string{}, InvertGroup(map[string]int{}))
	assert.Equal(t, map[int][]string{1: {"1"}, 2: {"2"}}, InvertGroup(map[string]int{"1": 1, "2": 2}))
	inversion := InvertGroup(map[string]int{"1": 1, "2": 1})
	assert.Equal(t, len(inversion), 1)
	assert.True(t, gslice.ContainsAll(inversion[1], "1", "2"))
}

func TestLoadKey(t *testing.T) {
	assert.Equal(t, goption.Nil[int](), LoadKey[int, string](nil, ""))
	assert.Equal(t, goption.OK(1),
		LoadKey(map[int]string{1: "1", 2: "2", 3: "3", 4: "4"}, "1"))
	assert.Equal(t, goption.OK(2),
		LoadKey(map[int]string{1: "1", 2: "2", 3: "3", 4: "4"}, "2"))
	assert.Equal(t, goption.Nil[int](),
		LoadKey(map[int]string{1: "1", 2: "2", 3: "3", 4: "4"}, "5"))
}

func TestLoadBy(t *testing.T) {
	assert.Equal(t, goption.Nil[string](),
		LoadBy[int, string](nil, func(k int, v string) bool {
			return v == ""
		}))
	assert.Equal(t, goption.OK("1"),
		LoadBy(map[int]string{1: "1", 2: "2", 3: "3", 4: "4"}, func(k int, v string) bool {
			return k == 1
		}))
	assert.Equal(t, goption.OK("2"),
		LoadBy(map[int]string{1: "1", 2: "2", 3: "3", 4: "4"}, func(k int, v string) bool {
			return v == "2"
		}))
	assert.Equal(t, goption.Nil[string](),
		LoadBy(map[int]string{1: "1", 2: "2", 3: "3", 4: "4"}, func(k int, v string) bool {
			return k == 0 || v == ""
		}))
}

func TestLoadKeyBy(t *testing.T) {
	assert.Equal(t, goption.Nil[int](),
		LoadKeyBy[int, string](nil, func(k int, v string) bool {
			return v == ""
		}))
	assert.Equal(t, goption.OK(1),
		LoadKeyBy(map[int]string{1: "1", 2: "2", 3: "3", 4: "4"}, func(k int, v string) bool {
			return k == 1
		}))
	assert.Equal(t, goption.OK(2),
		LoadKeyBy(map[int]string{1: "1", 2: "2", 3: "3", 4: "4"}, func(k int, v string) bool {
			return v == "2"
		}))
	assert.Equal(t, goption.Nil[int](),
		LoadKeyBy(map[int]string{1: "1", 2: "2", 3: "3", 4: "4"}, func(k int, v string) bool {
			return k == 0 || v == ""
		}))
}

func TestLoadItemBy(t *testing.T) {
	assert.Equal(t, goption.Nil[tuple.T2[int, string]](),
		LoadItemBy[int, string](nil, func(k int, v string) bool {
			return v == ""
		}))
	assert.Equal(t, goption.OK(tuple.Make2(1, "1")),
		LoadItemBy(map[int]string{1: "1", 2: "2", 3: "3", 4: "4"}, func(k int, v string) bool {
			return k == 1
		}))
	assert.Equal(t, goption.OK(tuple.Make2(2, "2")),
		LoadItemBy(map[int]string{1: "1", 2: "2", 3: "3", 4: "4"}, func(k int, v string) bool {
			return v == "2"
		}))
	assert.Equal(t, goption.Nil[tuple.T2[int, string]](),
		LoadItemBy(map[int]string{1: "1", 2: "2", 3: "3", 4: "4"}, func(k int, v string) bool {
			return k == 0 || v == ""
		}))
}

func TestContains(t *testing.T) {
	assert.False(t, Contains(map[int]string{}, 1))
	assert.False(t, Contains[int, string](nil, 1))
	assert.True(t, Contains(map[int]string{1: "", 2: ""}, 1))
	assert.False(t, Contains(map[int]string{1: "", 2: ""}, 3))
}

func TestContainsAny(t *testing.T) {
	assert.False(t, ContainsAny(map[int]string{}, 1))
	assert.False(t, ContainsAny[int, string](nil, 1))
	assert.False(t, ContainsAny[int, string](nil))
	assert.True(t, ContainsAny(map[int]string{1: "", 2: ""}, 1, 2))
	assert.True(t, ContainsAny(map[int]string{1: "", 2: ""}, 1, 3))
	assert.False(t, ContainsAny(map[int]string{1: "", 2: ""}, 3, 4))
}

func TestContainsAll(t *testing.T) {
	assert.False(t, ContainsAll(map[int]string{}, 1))
	assert.False(t, ContainsAll[int, string](nil, 1))
	assert.True(t, ContainsAll[int, string](nil))
	assert.True(t, ContainsAll(map[int]string{1: "", 2: ""}, 1, 2))
	assert.False(t, ContainsAll(map[int]string{1: "", 2: ""}, 1, 3))
}

func TestLoadAll(t *testing.T) {
	assert.Equal(t, nil, LoadAll(map[int]int{}, 1, 2, 3))
	assert.Equal(t, nil, LoadAll(map[int]string{1: "1", 2: "2"}))
	assert.Equal(t, []string{"1", "2"},
		LoadAll(map[int]string{1: "1", 2: "2", 3: "3"}, 1, 2))
	assert.Equal(t, nil,
		LoadAll(map[int]string{1: "1", 2: "2", 3: "3"}, 1, 4))
}

func TestLoadAny(t *testing.T) {
	assert.Equal(t, goption.Nil[int](), LoadAny(map[int]int{}, 1, 2, 3))
	assert.Equal(t, goption.Nil[string](), LoadAny(map[int]string{1: "1", 2: "2"}))
	assert.Equal(t, goption.OK("1"),
		LoadAny(map[int]string{1: "1", 2: "2", 3: "3"}, 1, 2))
	assert.Equal(t, goption.OK("2"),
		LoadAny(map[int]string{1: "1", 2: "2", 3: "3"}, 2, 1))
	assert.Equal(t, goption.OK("1"),
		LoadAny(map[int]string{1: "1", 2: "2", 3: "3"}, 9, 1))
	assert.Equal(t, goption.Nil[string](),
		LoadAny(map[int]string{1: "1", 2: "2", 3: "3"}, 9, 10))
}

func TestLoadSome(t *testing.T) {
	assert.Equal(t, nil, LoadSome(map[int]int{}, 1, 2, 3))
	assert.Equal(t, nil, LoadSome(map[int]string{1: "1", 2: "2"}))
	assert.Equal(t, []string{"1", "2"},
		LoadSome(map[int]string{1: "1", 2: "2", 3: "3"}, 1, 2))
	assert.Equal(t, []string{"1"},
		LoadSome(map[int]string{1: "1", 2: "2", 3: "3"}, 1, 4))
}

func TestSum(t *testing.T) {
	assert.Equal(t, 0, Sum(map[int]int{}))
	assert.Equal(t, 6, Sum(map[string]int{"1": 1, "2": 2, "3": 3}))
}

func TestSumBy(t *testing.T) {
	type Foo struct {
		Value int
	}
	getValue := func(foo Foo) int {
		return foo.Value
	}
	assert.Equal(t, 0, SumBy(map[int]Foo{}, getValue))
	assert.Equal(t, 6, SumBy(map[string]Foo{"1": {1}, "2": {2}, "3": {3}}, getValue))
}

func TestAvg(t *testing.T) {
	assert.Equal(t, 0.0, Avg(map[int]int{}))
	assert.Equal(t, 2.0, Avg(map[string]int{"1": 1, "2": 2, "3": 3}))
}

func TestAvgBy(t *testing.T) {
	type Foo struct {
		Value int
	}
	getValue := func(foo Foo) int {
		return foo.Value
	}
	assert.Equal(t, 0.0, AvgBy(map[int]Foo{}, getValue))
	assert.Equal(t, 2.0, AvgBy(map[string]Foo{"1": {1}, "2": {2}, "3": {3}}, getValue))
}

func TestMax(t *testing.T) {
	assert.Equal(t, 3, Max(map[string]int{"1": 1, "2": 2, "3": 3}).Value())
}

func TestMaxBy(t *testing.T) {
	type Foo struct {
		Value int
	}
	// FIXME: use greflect
	less := func(foo1, foo2 Foo) bool {
		return foo1.Value < foo2.Value
	}
	assert.Equal(t, Foo{3}, MaxBy(map[string]Foo{"1": {1}, "2": {2}, "3": {3}}, less).Value())
}

func TestMin(t *testing.T) {
	assert.Equal(t, 1, Min(map[string]int{"1": 1, "2": 2, "3": 3}).Value())
}

func TestMinBy(t *testing.T) {
	type Foo struct {
		Value int
	}
	// FIXME: use greflect
	less := func(foo1, foo2 Foo) bool {
		return foo1.Value < foo2.Value
	}
	assert.Equal(t, Foo{1}, MinBy(map[string]Foo{"1": {1}, "2": {2}, "3": {3}}, less).Value())
}

func TestMinMax(t *testing.T) {
	assert.Equal(t, tuple.Make2(1, 3), MinMax(map[string]int{"1": 1, "2": 2, "3": 3}).Value())
}

func TestMinMaxBy(t *testing.T) {
	type Foo struct {
		Value int
	}
	// FIXME: use greflect
	less := func(foo1, foo2 Foo) bool {
		return foo1.Value < foo2.Value
	}
	assert.Equal(t, tuple.Make2(Foo{1}, Foo{3}), MinMaxBy(map[string]Foo{"1": {1}, "2": {2}, "3": {3}}, less).Value())
}

func TestChunk(t *testing.T) {
	{
		m := map[int]string{1: "1", 2: "2", 3: "3", 4: "4", 5: "5"}
		chunks := Chunk(m, 2)
		assert.Equal(t, 3, len(chunks))
		assert.Equal(t, 2, len(chunks[0]))
		assert.Equal(t, 2, len(chunks[1]))
		assert.Equal(t, 1, len(chunks[2]))
		// TODO: Check equal
	}
}

func TestDivide(t *testing.T) {
	{
		m := map[int]string{1: "1", 2: "2", 3: "3", 4: "4", 5: "5"}
		chunks := Divide(m, 2)
		assert.Equal(t, 2, len(chunks))
		assert.Equal(t, 3, len(chunks[0]))
		assert.Equal(t, 2, len(chunks[1]))
		// TODO: Check equal
	}
}

func TestPtrOf(t *testing.T) {
	{
		m := map[int]string{1: "1", 2: "2"}
		ptrs := PtrOf(m)
		assert.Equal(t, map[int]*string{1: gptr.Of("1"), 2: gptr.Of("2")}, ptrs)
	}

	// Test modifying pointer.
	{
		m := map[int]string{1: "1", 2: "2"}
		ptrs := PtrOf(m)
		*ptrs[1] = ""
		assert.Equal(t, "", *ptrs[1])
		assert.Equal(t, "1", m[1])
	}
}

func TestIndirect(t *testing.T) {
	assert.Equal(t,
		map[int]string{1: "1", 2: "2"},
		Indirect(map[int]*string{1: gptr.Of("1"), 2: gptr.Of("2"), 3: nil}))
	assert.Equal(t,
		map[int]string{1: "1", 2: "2"},
		Indirect(map[int]*string{1: gptr.Of("1"), 2: gptr.Of("2")}))
}

func TestIndirectOr(t *testing.T) {
	assert.Equal(t,
		map[int]string{1: "1", 2: "2", 3: ""},
		IndirectOr(map[int]*string{1: gptr.Of("1"), 2: gptr.Of("2"), 3: nil}, ""))
	assert.Equal(t,
		map[int]string{1: "1", 2: "2"},
		IndirectOr(map[int]*string{1: gptr.Of("1"), 2: gptr.Of("2")}, ""))
}

func TestTypeAssert(t *testing.T) {
	assert.Equal(t, map[int]int{1: 1, 2: 2},
		TypeAssert[int](map[int]any{1: 1, 2: 2}))
	assert.Equal(t, map[int]any{1: 1, 2: 2},
		TypeAssert[any](map[int]int{1: 1, 2: 2}))

	assert.Panic(t, func() {
		TypeAssert[float64](map[int]int{1: 1, 2: 2})
	})
}

func TestUnion(t *testing.T) {
	assert.Equal(t, map[int]int{1: 1, 2: 2, 3: 3, 4: 4},
		Union(map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, nil))
	assert.Equal(t, map[int]int{1: 1, 2: 2, 3: 3, 4: 4},
		Union(nil, map[int]int{1: 1, 2: 2, 3: 3, 4: 4}))

	// Empty
	assert.Equal(t, map[int]string{}, Union[int, string]())
	assert.Equal(t, map[int]int{}, Union[int, int](nil))
	assert.Equal(t, map[int]int{}, Union[int, int](nil, nil))
	assert.Equal(t, map[int]int{}, Union[int, int](nil, nil, nil))

	// New value replace old.
	assert.Equal(t, map[int]int{1: 3},
		Union(
			map[int]int{1: 1},
			map[int]int{1: 2},
			map[int]int{},
			map[int]int{1: 3}))
	assert.Equal(t, map[int]int{1: 1, 2: 2, 3: 3, 4: 4},
		Union(
			map[int]int{1: 1, 2: 2, 3: 3, 4: 4},
			map[int]int{1: 1, 2: 2, 3: 3, 4: 4}))
	assert.Equal(t, map[int]int{1: 1, 2: 2, 3: 3, 4: 4},
		Union(
			map[int]int{1: 0, 2: 0},
			map[int]int{1: 1, 2: 2, 3: 3, 4: 4}))
	assert.Equal(t, map[int]int{1: 1, 2: 2, 3: 3, 4: 4},
		Union(
			map[int]int{1: 1, 2: 1},
			map[int]int{2: 2, 3: 3, 4: 4}))
}

func TestUnionBy(t *testing.T) {
	assert.Equal(t, map[int]int{1: 1, 2: 2, 3: 3, 4: 4},
		UnionBy(gslice.Of(map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, nil), nil))
	assert.Equal(t, map[int]int{1: 1, 2: 2, 3: 3, 4: 4},
		Union(nil, map[int]int{1: 1, 2: 2, 3: 3, 4: 4}))

	// Empty
	assert.Equal(t, map[int]string{}, UnionBy[int, string, map[int]string](nil, nil))
	assert.Equal(t, map[int]int{}, UnionBy([]map[int]int{nil}, nil))
	assert.Equal(t, map[int]int{}, UnionBy([]map[int]int{nil, nil}, nil))
	assert.Equal(t, map[int]int{}, UnionBy([]map[int]int{nil, nil, nil}, nil))

	// Nil [ConflictFunc] causes PANIC.
	assert.Panic(t, func() {
		_ = UnionBy(
			gslice.Of(
				map[int]int{1: 1},
				map[int]int{1: 2},
				map[int]int{},
				map[int]int{1: 3}),
			nil)
	})
	assert.NotPanic(t, func() {
		assert.Equal(t, map[int]int{1: 3}, UnionBy(
			gslice.Of(
				map[int]int{1: 1},
				map[int]int{1: 2},
				map[int]int{},
				map[int]int{1: 3}),
			DiscardOld[int, int]()))
	})

	assert.Equal(t, map[int]int{1: 3},
		UnionBy(
			gslice.Of(
				map[int]int{1: 1},
				map[int]int{1: 2},
				map[int]int{},
				map[int]int{1: 3}),
			DiscardOld[int, int]()))
	assert.Equal(t, map[int]int{1: 1},
		UnionBy(
			gslice.Of(
				map[int]int{1: 1},
				map[int]int{1: 2},
				map[int]int{},
				map[int]int{1: 3}),
			DiscardNew[int, int]()))
	assert.Equal(t, map[int]int{1: 3},
		UnionBy(
			gslice.Of(
				map[int]int{1: 0},
				map[int]int{1: 1},
				map[int]int{1: 2},
				map[int]int{},
				map[int]int{1: 3},
				map[int]int{1: 0}),
			DiscardZero(DiscardOld[int, int]())))
	assert.Equal(t, map[int]int{1: 1},
		UnionBy(
			gslice.Of(
				map[int]int{1: 0},
				map[int]int{1: 1},
				map[int]int{1: 2},
				map[int]int{},
				map[int]int{1: 3},
				map[int]int{1: 0}),
			DiscardZero(DiscardNew[int, int]())))

	{ // Test custom map type.
		type M map[int]int
		assert.Equal(t, M{1: 2},
			UnionBy(
				gslice.Of(
					M{1: 0},
					M{1: 2},
					M{1: 0}),
				DiscardZero[int, int](nil)))
	}
}

func TestDiff(t *testing.T) {
	assert.Equal(t, map[int]string{}, Diff(map[int]string{}))
	assert.Equal(t, map[int]string{1: "1"}, Diff(map[int]string{1: "1"}))
	assert.Equal(t, map[int]string{1: "1", 2: "2"}, Diff(map[int]string{1: "1", 2: "2"}))
	assert.Equal(t, map[int]string{1: "1"}, Diff(map[int]string{1: "1"}, nil))
	assert.Equal(t, map[int]string{1: "1"}, Diff(map[int]string{1: "1"}, nil, nil, nil))

	assert.Equal(t, map[int]int{2: 2, 3: 3},
		Diff(
			map[int]int{1: 1, 2: 2, 3: 3},
			map[int]int{1: 2}, map[int]int{1: 3}))
	assert.Equal(t, map[int]int{},
		Diff(
			map[int]int{1: 1, 2: 2, 3: 3, 4: 4},
			map[int]int{1: 1, 2: 2, 3: 3, 4: 4}))
	assert.Equal(t, map[int]int{},
		Diff(
			map[int]int{1: 1, 2: 2, 3: 3, 4: 4},
			map[int]int{1: 1, 2: 2}, map[int]int{3: 3, 4: 4}))
	assert.Equal(t, map[int]int{1: 1, 5: 5},
		Diff(
			map[int]int{1: 1, 2: 1, 5: 5},
			map[int]int{2: 2, 3: 3, 4: 4}))
}

func TestIntersect(t *testing.T) {
	assert.Equal(t, map[int]int{}, Intersect[int, int]())
	assert.Equal(t, map[int]int{}, Intersect[int, int](nil))
	assert.Equal(t, map[int]int{}, Intersect[int, int](nil, nil))
	assert.Equal(t, map[int]int{}, Intersect[int, int](nil, nil, nil))

	assert.Equal(t, map[int]int{}, Intersect(nil, map[int]int{1: 1}, nil))
	assert.Equal(t, map[int]int{}, Intersect(map[int]int{1: 1}, nil, nil))
	assert.Equal(t, map[int]int{}, Intersect(nil, nil, map[int]int{1: 1}, nil))

	assert.Equal(t, map[int]int{1: 1, 2: 2},
		Intersect(map[int]int{1: 1, 2: 2}, map[int]int{1: 1, 2: 2, 3: 3}))
	assert.Equal(t, map[int]int{1: 1, 2: 2},
		Intersect(map[int]int{1: 1, 2: 2, 3: 3}, map[int]int{1: 1, 2: 2}))

	// New value replaces old one.
	assert.Equal(t, map[int]int{1: 1, 2: -1},
		Intersect(map[int]int{1: 1, 2: 2}, map[int]int{1: 1, 2: -1, 3: 3}))
	assert.Equal(t, map[int]int{1: 1, 2: -1},
		Intersect(map[int]int{1: 1, 2: 2, 3: 3}, map[int]int{1: 1, 2: -1}))

	assert.Equal(t, map[int]int{1: 3},
		Intersect(
			map[int]int{1: 1, 2: 2, 3: 3},
			map[int]int{1: 2},
			map[int]int{1: 3}))
	assert.Equal(t, map[int]string{1: "1"}, Intersect(map[int]string{1: "1"}))

	assert.Equal(t, map[int]int{1: 1, 2: 2, 3: 3, 4: 4},
		Intersect(
			map[int]int{1: 1, 2: 2, 3: 3, 4: 4},
			map[int]int{1: 1, 2: 2, 3: 3, 4: 4}))
	assert.Equal(t, map[int]int{},
		Intersect(
			map[int]int{1: 1, 2: 2, 3: 3, 4: 4},
			map[int]int{1: 1, 2: 2},
			map[int]int{3: 3, 4: 4}))
	assert.Equal(t, map[int]int{1: 1, 2: 2, 5: 5},
		Intersect(
			map[int]int{1: 1, 2: 1, 5: 5},
			map[int]int{2: 2, 3: 3, 4: 4, 1: 1, 5: 5}))
}

func TestIntersectBy(t *testing.T) {
	// Empty
	assert.Equal(t, map[int]int{}, IntersectBy[int, int, map[int]int](nil, nil))
	assert.Equal(t, map[int]int{}, IntersectBy([]map[int]int{}, nil))
	assert.Equal(t, map[int]int{}, IntersectBy([]map[int]int{nil}, nil))
	assert.Equal(t, map[int]int{}, IntersectBy([]map[int]int{nil, nil}, nil))
	assert.Equal(t, map[int]int{}, IntersectBy([]map[int]int{nil, nil, nil}, nil))

	assert.Equal(t, map[int]int{}, IntersectBy(gslice.Of(nil, map[int]int{1: 1}, nil), nil))
	assert.Equal(t, map[int]int{}, IntersectBy(gslice.Of(map[int]int{1: 1}, nil, nil), nil))
	assert.Equal(t, map[int]int{}, IntersectBy(gslice.Of(nil, nil, map[int]int{1: 1}, nil), nil))

	// Nil [ConflictFunc] causes PANIC.
	assert.Panic(t, func() {
		_ = IntersectBy(
			gslice.Of(
				map[int]int{1: 1},
				map[int]int{1: 2},
				map[int]int{1: 3}),
			nil)
	})
	assert.NotPanic(t, func() {
		_ = IntersectBy(
			gslice.Of(
				map[int]int{1: 1},
				map[int]int{1: 2},
				map[int]int{1: 3}),
			DiscardOld[int, int]())
	})

	assert.Equal(t, map[int]int{1: 1, 2: -1},
		IntersectBy(
			gslice.Of(
				map[int]int{1: 1, 2: 2},
				map[int]int{1: 1, 2: -1, 3: 3}),
			DiscardOld[int, int]()))
	assert.Equal(t, map[int]int{1: 1, 2: 2},
		IntersectBy(
			gslice.Of(
				map[int]int{1: 1, 2: 2},
				map[int]int{1: 1, 2: -1, 3: 3}),
			DiscardNew[int, int]()))

	{ // Test custom map type.
		type M map[int]int
		assert.Equal(t, M{1: 1, 2: 2},
			IntersectBy(
				gslice.Of(
					M{1: 1, 2: 2},
					M{1: 1, 2: -1, 3: 3}),
				DiscardNew[int, int]()))
	}
}

func TestCompact(t *testing.T) {
	assert.Equal(t, map[int]int{}, Compact[int, int](nil))
	assert.Equal(t, map[int]int{}, Compact(map[int]int{}))
	assert.Equal(t,
		map[int]string{1: "foo", 3: "bar"},
		Compact(map[int]string{0: "", 1: "foo", 2: "", 3: "bar"}))
	assert.Equal(t,
		map[int]string{0: "foo", 1: "foo", 2: "bar", 3: "bar"},
		Compact(map[int]string{0: "foo", 1: "foo", 2: "bar", 3: "bar"}))
	assert.Equal(t,
		map[int]string{},
		Compact(map[int]string{0: "", 1: "", 2: "", 3: ""}))
}

func TestToSlice(t *testing.T) {
	f := func(k, v int) string {
		return fmt.Sprintf("%d: %d", k, v)
	}

	assert.Equal(t,
		[]string{"1: 1", "2: 2", "3: 3"},
		gslice.SortClone(
			ToSlice(map[int]int{1: 1, 2: 2, 3: 3}, f)))
	assert.Equal(t,
		[]string{},
		gslice.SortClone(
			ToSlice(map[int]int{}, f)))

	assert.NotPanic(t, func() {
		assert.Equal(t,
			[]string{},
			gslice.SortClone(
				ToSlice(map[int]any{}, func(k int, v any) string { panic("panic") })))
	})
}

func TestToOrderedSlice(t *testing.T) {
	f := func(k, v int) string {
		return fmt.Sprintf("%d: %d", k, v)
	}

	assert.Equal(t,
		[]string{"1: 1", "2: 2", "3: 3"},
		ToOrderedSlice(map[int]int{1: 1, 2: 2, 3: 3}, f))
	assert.Equal(t,
		[]string{},
		ToOrderedSlice(map[int]int{}, f))

	assert.NotPanic(t, func() {
		assert.Equal(t,
			[]string{},
			ToOrderedSlice(map[int]any{}, func(k int, v any) string { panic("panic") }))
	})
}

func TestFilterMapKeys(t *testing.T) {
	parseInt := func(s string) (int, bool) {
		ki, err := strconv.ParseInt(s, 10, 64)
		return int(ki), err == nil
	}
	assert.Equal(t,
		map[int]string{1: "1", 2: "2", 4: "b"},
		FilterMapKeys(map[string]string{"1": "1", "2": "2", "a": "3", "4": "b", "c": "c"}, parseInt))
	assert.Equal(t,
		map[int]string{4: "b"},
		FilterMapKeys(map[string]string{"a": "3", "4": "b"}, parseInt))
	assert.Equal(t,
		map[int]string{1: "1", 2: "2"},
		FilterMapKeys(map[string]string{"1": "1", "2": "2"}, parseInt))
	assert.Equal(t,
		map[int]string{},
		FilterMapKeys(map[string]string{}, parseInt))
	assert.Equal(t,
		map[int]string{},
		FilterMapKeys((map[string]string)(nil), parseInt))
}

func TestTryFilterMapKeys(t *testing.T) {
	parseInt := strconv.Atoi
	assert.Equal(t,
		map[int]string{1: "1", 2: "2", 4: "b"},
		TryFilterMapKeys(map[string]string{"1": "1", "2": "2", "a": "3", "4": "b", "c": "c"}, parseInt))
	assert.Equal(t,
		map[int]string{4: "b"},
		TryFilterMapKeys(map[string]string{"a": "3", "4": "b"}, parseInt))
	assert.Equal(t,
		map[int]string{1: "1", 2: "2"},
		TryFilterMapKeys(map[string]string{"1": "1", "2": "2"}, parseInt))
	assert.Equal(t,
		map[int]string{},
		TryFilterMapKeys(map[string]string{}, parseInt))
	assert.Equal(t,
		map[int]string{},
		TryFilterMapKeys((map[string]string)(nil), parseInt))
}

func TestFilterMapValues(t *testing.T) {
	parseInt := func(s string) (int, bool) {
		ki, err := strconv.ParseInt(s, 10, 64)
		return int(ki), err == nil
	}
	assert.Equal(t,
		map[string]int{"1": 1, "2": 2, "a": 3},
		FilterMapValues(map[string]string{"1": "1", "2": "2", "a": "3", "4": "b", "c": "c"}, parseInt))
	assert.Equal(t,
		map[string]int{"a": 3},
		FilterMapValues(map[string]string{"a": "3", "4": "b"}, parseInt))
	assert.Equal(t,
		map[string]int{"1": 1, "2": 2},
		FilterMapValues(map[string]string{"1": "1", "2": "2"}, parseInt))
	assert.Equal(t,
		map[string]int{},
		FilterMapValues(map[string]string{}, parseInt))
	assert.Equal(t,
		map[string]int{},
		FilterMapValues((map[string]string)(nil), parseInt))
}

func TestTryFilterMapValues(t *testing.T) {
	parseInt := strconv.Atoi
	assert.Equal(t,
		map[string]int{"1": 1, "2": 2, "a": 3},
		TryFilterMapValues(map[string]string{"1": "1", "2": "2", "a": "3", "4": "b", "c": "c"}, parseInt))
	assert.Equal(t,
		map[string]int{"a": 3},
		TryFilterMapValues(map[string]string{"a": "3", "4": "b"}, parseInt))
	assert.Equal(t,
		map[string]int{"1": 1, "2": 2},
		TryFilterMapValues(map[string]string{"1": "1", "2": "2"}, parseInt))
	assert.Equal(t,
		map[string]int{},
		TryFilterMapValues(map[string]string{}, parseInt))
	assert.Equal(t,
		map[string]int{},
		TryFilterMapValues((map[string]string)(nil), parseInt))
}

func TestDiscardOld(t *testing.T) {
	assert.Equal(t, "new", DiscardOld[int, string]()(10, "old", "new"))
}

func TestDiscardNew(t *testing.T) {
	assert.Equal(t, "old", DiscardNew[int, string]()(10, "old", "new"))
}

func TestDiscardZero(t *testing.T) {
	assert.Equal(t, "new", DiscardZero(DiscardOld[int, string]())(10, "old", "new"))
	assert.Equal(t, "old", DiscardZero(DiscardOld[int, string]())(10, "old", ""))
	assert.Equal(t, "new", DiscardZero(DiscardOld[int, string]())(10, "", "new"))
	assert.Equal(t, "", DiscardZero(DiscardOld[int, string]())(10, "", ""))

	assert.Equal(t, "old", DiscardZero(DiscardNew[int, string]())(10, "old", "new"))
	assert.Equal(t, "old", DiscardZero(DiscardNew[int, string]())(10, "old", ""))
	assert.Equal(t, "new", DiscardZero(DiscardNew[int, string]())(10, "", "new"))
	assert.Equal(t, "", DiscardZero(DiscardNew[int, string]())(10, "", ""))

	assert.Equal(t, "new", DiscardZero[int, string](nil)(10, "old", "new"))
	assert.Equal(t, "", DiscardZero[int, string](nil)(10, "", ""))
	assert.Equal(t, "", DiscardZero[int, string](nil)(10, "", ""))
	assert.Equal(t, "new", DiscardZero[int, string](nil)(10, "", "new"))
	assert.Equal(t, "old", DiscardZero[int, string](nil)(10, "old", ""))
}

func TestDiscardNil(t *testing.T) {
	assert.Equal(t, gptr.Of("new"), DiscardNil(DiscardOld[int, *string]())(10, gptr.Of("old"), gptr.Of("new")))
	assert.Equal(t, gptr.Of("old"), DiscardNil(DiscardOld[int, *string]())(10, gptr.Of("old"), nil))
	assert.Equal(t, gptr.Of("new"), DiscardNil(DiscardOld[int, *string]())(10, nil, gptr.Of("new")))
	assert.Equal(t, nil, DiscardNil(DiscardOld[int, *string]())(10, nil, nil))

	assert.Equal(t, gptr.Of("old"), DiscardNil(DiscardNew[int, *string]())(10, gptr.Of("old"), gptr.Of("new")))
	assert.Equal(t, gptr.Of("old"), DiscardNil(DiscardNew[int, *string]())(10, gptr.Of("old"), nil))
	assert.Equal(t, gptr.Of("new"), DiscardNil(DiscardNew[int, *string]())(10, nil, gptr.Of("new")))
	assert.Equal(t, nil, DiscardNil(DiscardNew[int, *string]())(10, nil, nil))

	assert.Equal(t, gptr.Of("new"), DiscardNil[int, string](nil)(10, gptr.Of("old"), gptr.Of("new")))
	assert.Equal(t, nil, DiscardNil[int, string](nil)(10, nil, nil))
	assert.Equal(t, gptr.Of("new"), DiscardNil[int, string](nil)(10, nil, gptr.Of("new")))
	assert.Equal(t, gptr.Of("old"), DiscardNil[int, string](nil)(10, gptr.Of("old"), nil))
}

func TestCount(t *testing.T) {
	assert.Equal(t, 0, Count(map[int]string{}, "2"))
	assert.Equal(t, 1, Count(map[int]string{1: "1", 2: "2", 3: "3"}, "2"))
	assert.Equal(t, 2, Count(map[int]string{1: "1", 2: "2", 3: "2"}, "2"))
	assert.Equal(t, 3, Count(map[int]string{1: "2", 2: "2", 3: "2"}, "2"))
	assert.Equal(t, 1, Count(map[int]string{1: "2", 2: "2", 3: "3"}, "3"))
	assert.Equal(t, 0, Count(map[int]string{1: "2", 2: "2", 3: "4"}, "3"))
}

func TestCountBy(t *testing.T) {
	f := func(k int, v string) bool {
		i, _ := strconv.Atoi(v)
		return k%2 == 1 && i%2 == 1
	}
	assert.Equal(t, 0, CountBy(map[int]string{}, f))
	assert.Equal(t, 2, CountBy(map[int]string{1: "1", 2: "2", 3: "3"}, f))
	assert.Equal(t, 1, CountBy(map[int]string{1: "1", 2: "2", 3: "2"}, f))
	assert.Equal(t, 1, CountBy(map[int]string{1: "1", 2: "2", 4: "3"}, f))
}

func TestCountValueBy(t *testing.T) {
	f := func(v string) bool {
		i, _ := strconv.Atoi(v)
		return i%2 == 1
	}
	assert.Equal(t, 0, CountValueBy(map[int]string{}, f))
	assert.Equal(t, 2, CountValueBy(map[int]string{1: "1", 2: "2", 3: "3"}, f))
	assert.Equal(t, 1, CountValueBy(map[int]string{1: "1", 2: "2", 3: "2"}, f))
	assert.Equal(t, 2, CountValueBy(map[int]string{1: "1", 2: "2", 4: "3"}, f))
}

func TestPop(t *testing.T) {
	{
		assert.Equal(t, goption.Nil[int](), Pop[int, int](nil))
	}
	{
		m := map[int]int{1: 1, 2: 2, 3: 3, 4: 4}
		assert.True(t, Pop(m).IsOK())
		assert.Equal(t, len(m), 3)
		assert.True(t, Pop(m).IsOK())
		assert.Equal(t, len(m), 2)
		assert.True(t, Pop(m).IsOK())
		assert.Equal(t, len(m), 1)
		assert.True(t, Pop(m).IsOK())
		assert.Equal(t, len(m), 0)
		assert.False(t, Pop(m).IsOK())
		assert.Equal(t, len(m), 0)
		assert.False(t, Pop(m).IsOK())
		assert.Equal(t, len(m), 0)
	}
	{
		m := map[int]int{1: 1}
		assert.Equal(t, goption.OK(1), Pop(m))
		assert.Equal(t, m, map[int]int{})
	}
}

func TestPopItem(t *testing.T) {
	{
		assert.Equal(t, goption.Nil[tuple.T2[int, int]](), PopItem[int, int](nil))
	}
	{
		m := map[int]int{1: 1, 2: 2, 3: 3, 4: 4}
		assert.True(t, PopItem(m).IsOK())
		assert.Equal(t, len(m), 3)
		assert.True(t, PopItem(m).IsOK())
		assert.Equal(t, len(m), 2)
		assert.True(t, PopItem(m).IsOK())
		assert.Equal(t, len(m), 1)
		assert.True(t, PopItem(m).IsOK())
		assert.Equal(t, len(m), 0)
		assert.False(t, PopItem(m).IsOK())
		assert.Equal(t, len(m), 0)
		assert.False(t, PopItem(m).IsOK())
		assert.Equal(t, len(m), 0)
	}
	{
		m := map[string]int{"1": 1}
		assert.Equal(t, goption.OK(tuple.Make2("1", 1)), PopItem(m))
		assert.Equal(t, m, map[string]int{})
	}
}

func TestPeek(t *testing.T) {
	{
		assert.Equal(t, goption.Nil[int](), Peek[int, int](nil))
	}
	{
		m := map[int]int{1: 1, 2: 2, 3: 3, 4: 4}
		assert.True(t, Peek(m).IsOK())
		assert.Equal(t, len(m), 4)
	}
	{
		m := map[int]int{1: 1}
		assert.Equal(t, goption.OK(1), Peek(m))
		assert.Equal(t, m, map[int]int{1: 1})
	}
}

func TestPeekItem(t *testing.T) {
	{
		assert.Equal(t, goption.Nil[tuple.T2[int, int]](), PeekItem[int, int](nil))
	}
	{
		m := map[int]int{1: 1, 2: 2, 3: 3, 4: 4}
		assert.True(t, PeekItem(m).IsOK())
		assert.Equal(t, len(m), 4)
	}
	{
		m := map[string]int{"1": 1}
		assert.Equal(t, goption.OK(tuple.Make2("1", 1)), PeekItem(m))
		assert.Equal(t, m, map[string]int{"1": 1})
	}
}
