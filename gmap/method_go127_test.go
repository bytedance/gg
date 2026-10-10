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

package gmap

import (
	"errors"
	"strconv"
	"testing"

	"github.com/bytedance/gg/internal/assert"
)

func TestWrap(t *testing.T) {
	m := Wrap(map[int]int{1: 2})
	assert.Equal(t, map[int]int{1: 2}, m.Unwrap())
	assert.Equal(t, 1, m.Len())
	assert.Equal(t, map[int]int{1: 2}, m.Clone().Unwrap())
	assert.Equal(t, map[int]int{1: 4}, m.CloneBy(func(v int) int { return v * 2 }).Unwrap())
	assert.Equal(t, 0, New[string, int](4).Len())
}

func TestMMap(t *testing.T) {
	assert.Equal(t, map[int]string{1: "2"}, Wrap(map[int]int{1: 2}).MapValues(strconv.Itoa).Unwrap())
	assert.Equal(t, map[string]int{"1": 2}, Wrap(map[int]int{1: 2}).MapKeys(strconv.Itoa).Unwrap())
	assert.Equal(t, map[string]string{"1": "2"}, Wrap(map[int]int{1: 2}).Map(
		func(k, v int) (string, string) { return strconv.Itoa(k), strconv.Itoa(v) },
	).Unwrap())

	// The key benefit of generic methods: key and value types change in a chain.
	assert.Equal(t,
		map[string]string{"1": "4"},
		Wrap(map[int]int{1: 2}).MapValues(func(v int) int { return v * 2 }).Map(
			func(k, v int) (string, string) { return strconv.Itoa(k), strconv.Itoa(v) },
		).Unwrap(),
	)

	// TryMap keeps the error.
	assert.Equal(t, map[int]string{1: "2"}, Wrap(map[int]int{1: 2}).TryMapValues(func(v int) (string, error) {
		return strconv.Itoa(v), nil
	}).Value().Unwrap())
	assert.True(t, Wrap(map[int]int{1: 2}).TryMapValues(func(v int) (string, error) {
		return "", errors.New("woof!")
	}).IsErr())
	assert.True(t, Wrap(map[int]int{1: 2}).TryMapKeys(func(k int) (string, error) {
		return "", errors.New("woof!")
	}).IsErr())
	assert.True(t, Wrap(map[int]int{1: 2}).TryMap(func(k, v int) (string, string, error) {
		return "", "", errors.New("woof!")
	}).IsErr())
}

func TestMFilter(t *testing.T) {
	m := Wrap(map[int]int{1: 2, 2: 3, 3: 4})
	assert.Equal(t, map[int]int{2: 3, 3: 4}, m.Filter(func(k, v int) bool { return k+v > 3 }).Unwrap())
	assert.Equal(t, map[int]int{1: 2}, m.Reject(func(k, v int) bool { return k+v > 3 }).Unwrap())
	assert.Equal(t, map[int]int{1: 2, 2: 3}, m.FilterKeys(func(k int) bool { return k < 3 }).Unwrap())
	assert.Equal(t, map[int]int{3: 4}, m.RejectKeys(func(k int) bool { return k < 3 }).Unwrap())
	assert.Equal(t, map[int]int{1: 2}, m.FilterByKeys(1).Unwrap())
	assert.Equal(t, map[int]int{2: 3, 3: 4}, m.RejectByKeys(1).Unwrap())
	assert.Equal(t, map[int]int{1: 2}, m.FilterValues(func(v int) bool { return v < 3 }).Unwrap())
	assert.Equal(t, map[int]int{2: 3, 3: 4}, m.RejectValues(func(v int) bool { return v < 3 }).Unwrap())
	assert.Equal(t, map[string]int{"2": 3}, m.FilterMapKeys(func(k int) (string, bool) {
		return strconv.Itoa(k), k == 2
	}).Unwrap())
	assert.Equal(t, map[int]string{1: "2"}, m.FilterMapValues(func(v int) (string, bool) {
		return strconv.Itoa(v), v == 2
	}).Unwrap())
}

func TestMLoad(t *testing.T) {
	m := Wrap(map[int]int{1: 2, 2: 3})
	assert.Equal(t, 2, m.Load(1).Value())
	assert.True(t, m.Load(9).IsNil())
	assert.True(t, m.Contains(1))
	assert.True(t, m.ContainsAny(1, 9))
	assert.False(t, m.ContainsAll(1, 9))
	assert.Equal(t, []int{2, 3}, m.LoadAll(1, 2))
	assert.Equal(t, 2, m.LoadAny(9, 1).Value())
	assert.Equal(t, []int{2}, m.LoadSome(1, 9))
	assert.Equal(t, 3, m.LoadBy(func(k, v int) bool { return k == 2 }).Value())
	assert.Equal(t, 2, m.LoadKeyBy(func(k, v int) bool { return k == 2 }).Value())

	item := m.LoadItemBy(func(k, v int) bool { return k == 2 })
	key, value := item.Value().Values()
	assert.Equal(t, 2, key)
	assert.Equal(t, 3, value)

	v, loaded := m.Clone().LoadOrStore(1, 9)
	assert.Equal(t, 2, v)
	assert.True(t, loaded)

	v, loaded = m.Clone().LoadOrStore(9, 9)
	assert.Equal(t, 9, v)
	assert.False(t, loaded)

	v, loaded = m.Clone().LoadOrStoreLazy(9, func() int { return 9 })
	assert.Equal(t, 9, v)
	assert.False(t, loaded)

	assert.Equal(t, 2, m.Clone().LoadAndDelete(1).Value())
	assert.True(t, m.Clone().LoadAndDelete(9).IsNil())

	assert.True(t, m.Clone().Peek().IsOK())
	assert.True(t, m.Clone().Pop().IsOK())
	assert.True(t, m.Clone().PopItem().IsOK())
}

func TestMKeysValues(t *testing.T) {
	m := Wrap(map[int]int{1: 2})
	assert.Equal(t, 1, len(m.Keys()))
	assert.Equal(t, []int{2}, m.Values())
	assert.Equal(t, 1, len(m.Items()))
	assert.Equal(t, []string{"1:2"}, m.ToSlice(func(k, v int) string {
		return strconv.Itoa(k) + ":" + strconv.Itoa(v)
	}))
	assert.Equal(t, map[int]int{1: 2}, Wrap(map[int]any{1: 2}).TypeAssert[int]().Unwrap())
}

func TestMSetOp(t *testing.T) {
	m := Wrap(map[int]int{1: 2, 2: 3})
	assert.Equal(t, map[int]int{1: 2, 2: 4}, m.Union(Wrap(map[int]int{2: 4})).Unwrap())
	// DiscardOld keeps the newer value on conflict.
	assert.Equal(t, map[int]int{1: 2, 2: 4}, m.UnionBy([]M[int, int]{Wrap(map[int]int{2: 4})}, DiscardOld[int, int]()).Unwrap())
	assert.Equal(t, map[int]int{1: 2}, m.Diff(Wrap(map[int]int{2: 3})).Unwrap())
	assert.Equal(t, map[int]int{2: 4}, m.Intersect(Wrap(map[int]int{2: 4})).Unwrap())
	assert.Equal(t, map[int]int{2: 3}, m.IntersectBy([]M[int, int]{Wrap(map[int]int{2: 4})}, DiscardNew[int, int]()).Unwrap())
	assert.Equal(t, map[int]int{1: 2, 2: 4}, m.Merge(Wrap(map[int]int{2: 4})).Unwrap())
}

func TestMChunk(t *testing.T) {
	chunks := Wrap(map[int]int{1: 2, 2: 3, 3: 4, 4: 5, 5: 6}).Chunk(2)
	assert.Equal(t, 3, len(chunks))
	assert.Equal(t, 5, chunks[0].Len()+chunks[1].Len()+chunks[2].Len())

	divided := Wrap(map[int]int{1: 2, 2: 3, 3: 4, 4: 5, 5: 6}).Divide(2)
	assert.Equal(t, 2, len(divided))
	assert.Equal(t, 1, Wrap(map[int]int{1: 2, 2: 3}).CountBy(func(k, v int) bool { return v > 2 }))
	assert.Equal(t, 1, Wrap(map[int]int{1: 2, 2: 3}).CountValueBy(func(v int) bool { return v > 2 }))
}

func TestMMath(t *testing.T) {
	m := Wrap(map[string]int{"a": 1, "b": 3})
	assert.Equal(t, 3, m.MaxBy(func(a, b int) bool { return a < b }).Value())
	assert.Equal(t, 1, m.MinBy(func(a, b int) bool { return a < b }).Value())
	assert.Equal(t, 8, m.SumBy(func(v int) int { return v * 2 }))
	assert.Equal(t, 4.0, m.AvgBy(func(v int) int { return v * 2 }))

	min, max := m.MinMaxBy(func(a, b int) bool { return a < b }).Value().Values()
	assert.Equal(t, 1, min)
	assert.Equal(t, 3, max)
}

func TestMFilterMapFamily(t *testing.T) {
	m := Wrap(map[int]int{1: 2, 2: 3})

	assert.Equal(t, map[int]int{1: 4}, m.FilterMap(func(k, v int) (int, int, bool) {
		return k, v * 2, k == 1
	}).Unwrap())

	// The item is dropped when f returns an error.
	assert.Equal(t, map[int]int{1: 4}, m.TryFilterMap(func(k, v int) (int, int, error) {
		if k == 2 {
			return 0, 0, errors.New("skip")
		}
		return k, v * 2, nil
	}).Unwrap())

	assert.Equal(t, map[string]int{"1": 2}, m.TryFilterMapKeys(func(k int) (string, error) {
		if k == 2 {
			return "", errors.New("skip")
		}
		return strconv.Itoa(k), nil
	}).Unwrap())

	assert.Equal(t, map[int]string{1: "2"}, m.TryFilterMapValues(func(v int) (string, error) {
		if v == 3 {
			return "", errors.New("skip")
		}
		return strconv.Itoa(v), nil
	}).Unwrap())
}

func TestMTryMapSuccess(t *testing.T) {
	m := Wrap(map[int]int{1: 2})

	assert.Equal(t, map[string]string{"1": "2"}, m.TryMap(func(k, v int) (string, string, error) {
		return strconv.Itoa(k), strconv.Itoa(v), nil
	}).Value().Unwrap())

	assert.Equal(t, map[string]int{"1": 2}, m.TryMapKeys(func(k int) (string, error) {
		return strconv.Itoa(k), nil
	}).Value().Unwrap())

	assert.Equal(t, map[int]string{1: "2"}, m.TryMapValues(func(v int) (string, error) {
		return strconv.Itoa(v), nil
	}).Value().Unwrap())
}

func TestMEqualByAndPeekItem(t *testing.T) {
	m := Wrap(map[int]int{1: 2})
	eq := func(a, b int) bool { return a == b }

	assert.True(t, m.EqualBy(map[int]int{1: 2}, eq))
	assert.False(t, m.EqualBy(map[int]int{1: 3}, eq))

	// PeekItem does not delete the item.
	key, value := m.PeekItem().Value().Values()
	assert.Equal(t, 1, key)
	assert.Equal(t, 2, value)
	assert.Equal(t, 1, m.Len())
}

func TestMUnwrapVariants(t *testing.T) {
	assert.Equal(t, map[int]int{1: 2}, Wrap(map[int]int{1: 2}).Unwrap())
	assert.Equal(t, map[int]int{1: 2}, WrapMK(map[int]int{1: 2}).Unwrap())
	assert.Equal(t, map[string]int{"a": 1}, WrapMV(map[string]int{"a": 1}).Unwrap())
	assert.Equal(t, map[string]int{"a": 1}, WrapMO(map[string]int{"a": 1}).Unwrap())
	assert.Equal(t, map[string]int{"a": 1}, WrapMN(map[string]int{"a": 1}).Unwrap())
}

func TestMK(t *testing.T) {
	mk := WrapMK(map[int]int{2: 3, 1: 2})
	assert.Equal(t, []int{1, 2}, mk.OrderedKeys())
	assert.Equal(t, []int{2, 3}, mk.OrderedValues())
	assert.Equal(t, 2, len(mk.OrderedItems()))
	assert.Equal(t, []string{"1", "2"}, mk.ToOrderedSlice(func(k, v int) string { return strconv.Itoa(k) }))

	// MK is converted back to M to keep chaining.
	assert.Equal(t, map[int]string{1: "2", 2: "3"}, mk.M().MapValues(strconv.Itoa).Unwrap())
}

func TestMV(t *testing.T) {
	mv := WrapMV(map[string]int{"a": 1, "b": 1})
	assert.Equal(t, map[int]string{1: "a"}, WrapMV(map[string]int{"a": 1}).Invert().Unwrap())
	assert.Equal(t, map[int]string{1: "a"}, WrapMV(map[string]int{"a": 1}).InvertBy(DiscardOld[int, string]()).Unwrap())
	// The order of the grouped keys is not guaranteed.
	grouped := mv.InvertGroup()
	assert.Equal(t, 1, len(grouped))
	assert.Equal(t, 2, len(grouped[1]))
	assert.Equal(t, map[string]int{"a": 1, "b": 1}, mv.FilterByValues(1).Unwrap())
	assert.Equal(t, map[string]int{"b": 2}, WrapMV(map[string]int{"a": 1, "b": 2}).RejectByValues(1).Unwrap())
	assert.Equal(t, "a", WrapMV(map[string]int{"a": 1}).LoadKey(1).Value())
	assert.True(t, mv.Equal(map[string]int{"a": 1, "b": 1}))
	assert.Equal(t, 2, mv.Count(1))
	assert.Equal(t, map[string]int{"b": 2}, WrapMV(map[string]int{"a": 0, "b": 2}).Compact().Unwrap())

	// MV is converted back to M to keep chaining.
	assert.Equal(t, map[string]string{"a": "1", "b": "1"}, mv.M().MapValues(strconv.Itoa).Unwrap())
}

func TestMO(t *testing.T) {
	mo := WrapMO(map[string]int{"a": 1, "b": 3})
	assert.Equal(t, 3, mo.Max().Value())
	assert.Equal(t, 1, mo.Min().Value())

	min, max := mo.MinMax().Value().Values()
	assert.Equal(t, 1, min)
	assert.Equal(t, 3, max)

	// Ordered values are comparable, so MO is convertible to MV as well.
	assert.True(t, mo.MV().Equal(map[string]int{"a": 1, "b": 3}))
	assert.Equal(t, map[string]string{"a": "1", "b": "3"}, mo.M().MapValues(strconv.Itoa).Unwrap())
}

func TestMN(t *testing.T) {
	mn := WrapMN(map[string]int{"a": 1, "b": 2})
	assert.Equal(t, 3, mn.Sum())
	assert.Equal(t, 1.5, mn.Avg())

	// Numbers are ordered and comparable, so MN is convertible to both.
	assert.Equal(t, 2, mn.MO().Max().Value())
	assert.True(t, mn.MV().Equal(map[string]int{"a": 1, "b": 2}))
	assert.Equal(t, map[string]string{"a": "1", "b": "2"}, mn.M().MapValues(strconv.Itoa).Unwrap())
}
