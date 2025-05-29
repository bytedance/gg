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

package skipmap

import (
	"encoding/json"
	"reflect"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/bytedance/gg/goption"
	"github.com/bytedance/gg/internal/assert"
	"github.com/bytedance/gg/internal/constraints"
	"github.com/bytedance/gg/internal/fastrand"
)

func TestOrdered(t *testing.T) {
	testSkipMapInt(t, func() anyskipmap[int] { return New[int, any]() })
	testSkipMapInt(t, func() anyskipmap[int] { return NewDesc[int, any]() })
	testSkipMapInt(t, func() anyskipmap[int] { return NewFunc[int, any](func(a, b int) bool { return a < b }) })
	testSkipMapString(t, func() anyskipmap[string] { return New[string, any]() })
	testSkipMapString(t, func() anyskipmap[string] { return NewDesc[string, any]() })
	testSkipMapString(t, func() anyskipmap[string] { return NewFunc[string, any](func(a, b string) bool { return a < b }) })
	testSkipMapToMap(t, func() orderedskipmap[int] { return New[int, any]() })
	testSkipMapToMap(t, func() orderedskipmap[int] { return NewDesc[int, any]() })
}

type anyskipmap[T any] interface {
	Store(key T, value any)
	Load(key T) (any, bool)
	Delete(key T) bool
	LoadAndDelete(key T) (any, bool)
	LoadOrStore(key T, value any) (any, bool)
	LoadOrStoreLazy(key T, f func() any) (any, bool)
	Range(f func(key T, value any) bool)
	Len() int
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}

type orderedskipmap[T constraints.Ordered] interface {
	anyskipmap[T]
	ToMap() map[T]any
}

func testSkipMapInt(t *testing.T, newset func() anyskipmap[int]) {
	m := newset()

	// Correctness.
	m.Store(123, "123")
	v, ok := m.Load(123)
	if !ok || v != "123" || m.Len() != 1 {
		t.Fatal("invalid")
	}

	m.Store(123, "456")
	v, ok = m.Load(123)
	if !ok || v != "456" || m.Len() != 1 {
		t.Fatal("invalid")
	}

	m.Store(123, 456)
	v, ok = m.Load(123)
	if !ok || v != 456 || m.Len() != 1 {
		t.Fatal("invalid")
	}

	m.Delete(123)
	v, ok = m.Load(123)
	if ok || m.Len() != 0 || v != nil {
		t.Fatal("invalid")
	}

	v, loaded := m.LoadOrStore(123, 456)
	if loaded || v != 456 || m.Len() != 1 {
		t.Fatal("invalid")
	}

	v, loaded = m.LoadOrStore(123, 789)
	if !loaded || v != 456 || m.Len() != 1 {
		t.Fatal("invalid")
	}

	v, ok = m.Load(123)
	if !ok || v != 456 || m.Len() != 1 {
		t.Fatal("invalid")
	}

	v, ok = m.LoadAndDelete(123)
	if !ok || v != 456 || m.Len() != 0 {
		t.Fatal("invalid")
	}

	_, ok = m.LoadOrStore(123, 456)
	if ok || m.Len() != 1 {
		t.Fatal("invalid")
	}

	m.LoadOrStore(456, 123)
	if ok || m.Len() != 2 {
		t.Fatal("invalid")
	}

	m.Range(func(key int, _ interface{}) bool {
		if key == 123 {
			m.Store(123, 123)
		} else if key == 456 {
			m.LoadAndDelete(456)
		}
		return true
	})

	v, ok = m.Load(123)
	if !ok || v != 123 || m.Len() != 1 {
		t.Fatal("invalid")
	}

	// Concurrent.
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		i := i
		wg.Add(1)
		go func() {
			m.Store(i, int(i+1000))
			wg.Done()
		}()
	}
	wg.Wait()
	wg.Add(1)
	go func() {
		m.Delete(600)
		wg.Done()
	}()
	wg.Wait()
	wg.Add(1)
	var count int64
	go func() {
		m.Range(func(_ int, _ interface{}) bool {
			atomic.AddInt64(&count, 1)
			return true
		})
		wg.Done()
	}()
	wg.Wait()

	val, ok := m.Load(500)
	if !ok || reflect.TypeOf(val).Kind().String() != "int" || val.(int) != 1500 {
		t.Fatal("fail")
	}

	_, ok = m.Load(600)
	if ok {
		t.Fatal("fail")
	}

	if m.Len() != 999 || int(count) != m.Len() {
		t.Fatal("fail")
	}
	// Correctness 2.
	var m1 sync.Map
	m2 := newset()
	var v1, v2 interface{}
	var ok1, ok2 bool
	for i := 0; i < 100000; i++ {
		rd := int(fastrand.Uint32n(10))
		r1, r2 := int(fastrand.Uint32n(100)), int(fastrand.Uint32n(100))
		if rd == 0 {
			m1.Store(r1, r2)
			m2.Store(r1, r2)
		} else if rd == 1 {
			v1, ok1 = m1.LoadAndDelete(r1)
			v2, ok2 = m2.LoadAndDelete(r1)
			if ok1 != ok2 || v1 != v2 {
				t.Fatal(rd, v1, ok1, v2, ok2)
			}
		} else if rd == 2 {
			v1, ok1 = m1.LoadOrStore(r1, r2)
			v2, ok2 = m2.LoadOrStore(r1, r2)
			if ok1 != ok2 || v1 != v2 {
				t.Fatal(rd, v1, ok1, v2, ok2, "input -> ", r1, r2)
			}
		} else if rd == 3 {
			m1.Delete(r1)
			m2.Delete(r1)
		} else if rd == 4 {
			m2.Range(func(key int, value interface{}) bool {
				v, ok := m1.Load(key)
				if !ok || v != value {
					t.Fatal(v, ok, key, value)
				}
				return true
			})
		} else {
			v1, ok1 = m1.Load(r1)
			v2, ok2 = m2.Load(r1)
			if ok1 != ok2 || v1 != v2 {
				t.Fatal(rd, v1, ok1, v2, ok2)
			}
		}
	}
	// Correctness 3. (LoadOrStore)
	// Only one LoadOrStore can successfully insert its key and value.
	// And the returned value is unique.
	mp := newset()
	tmpmap := newset()
	samekey := 123
	var added int64
	for i := 1; i < 1000; i++ {
		wg.Add(1)
		go func() {
			v := fastrand.Int63()
			actual, loaded := mp.LoadOrStore(samekey, v)
			if !loaded {
				atomic.AddInt64(&added, 1)
			}
			tmpmap.Store(int(actual.(int64)), nil)
			wg.Done()
		}()
	}
	wg.Wait()
	if added != 1 {
		t.Fatal("only one LoadOrStore can successfully insert a key and value")
	}
	if tmpmap.Len() != 1 {
		t.Fatal("only one value can be returned from LoadOrStore")
	}
	// Correctness 4. (LoadAndDelete)
	// Only one LoadAndDelete can successfully get a value.
	mp = newset()
	tmpmap = newset()
	samekey = 123
	added = 0 // int64
	mp.Store(samekey, 555)
	for i := 1; i < 1000; i++ {
		wg.Add(1)
		go func() {
			value, loaded := mp.LoadAndDelete(samekey)
			if loaded {
				atomic.AddInt64(&added, 1)
				if value != 555 {
					panic("invalid")
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	if added != 1 {
		t.Fatal("Only one LoadAndDelete can successfully get a value")
	}

	// Correctness 5. (LoadOrStoreLazy)
	mp = newset()
	tmpmap = newset()
	samekey = 123
	added = 0
	var fcalled int64
	valuef := func() interface{} {
		atomic.AddInt64(&fcalled, 1)
		return fastrand.Int63()
	}
	for i := 1; i < 1000; i++ {
		wg.Add(1)
		go func() {
			actual, loaded := mp.LoadOrStoreLazy(samekey, valuef)
			if !loaded {
				atomic.AddInt64(&added, 1)
			}
			tmpmap.Store(int(actual.(int64)), nil)
			wg.Done()
		}()
	}
	wg.Wait()
	if added != 1 || fcalled != 1 {
		t.Fatal("only one LoadOrStoreLazy can successfully insert a key and value")
	}
	if tmpmap.Len() != 1 {
		t.Fatal("only one value can be returned from LoadOrStoreLazy")
	}
}

func testSkipMapString(t *testing.T, newset func() anyskipmap[string]) {
	m := newset()

	// Correctness.
	m.Store("123", "123")
	v, ok := m.Load("123")
	if !ok || v != "123" || m.Len() != 1 {
		t.Fatal("invalid")
	}

	m.Store("123", "456")
	v, ok = m.Load("123")
	if !ok || v != "456" || m.Len() != 1 {
		t.Fatal("invalid")
	}

	m.Store("123", 456)
	v, ok = m.Load("123")
	if !ok || v != 456 || m.Len() != 1 {
		t.Fatal("invalid")
	}

	m.Delete("123")
	_, ok = m.Load("123")
	if ok || m.Len() != 0 {
		t.Fatal("invalid")
	}

	_, ok = m.LoadOrStore("123", 456)
	if ok || m.Len() != 1 {
		t.Fatal("invalid")
	}

	v, ok = m.Load("123")
	if !ok || v != 456 || m.Len() != 1 {
		t.Fatal("invalid")
	}

	v, ok = m.LoadAndDelete("123")
	if !ok || v != 456 || m.Len() != 0 {
		t.Fatal("invalid")
	}

	_, ok = m.LoadOrStore("123", 456)
	if ok || m.Len() != 1 {
		t.Fatal("invalid")
	}

	m.LoadOrStore("456", 123)
	if ok || m.Len() != 2 {
		t.Fatal("invalid")
	}

	m.Range(func(key string, value interface{}) bool {
		if key == "123" {
			m.Store("123", 123)
		} else if key == "456" {
			m.LoadAndDelete("456")
		}
		return true
	})

	v, ok = m.Load("123")
	if !ok || v != 123 || m.Len() != 1 {
		t.Fatal("invalid")
	}

	// Concurrent.
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		i := i
		wg.Add(1)
		go func() {
			n := strconv.Itoa(i)
			m.Store(n, int(i+1000))
			wg.Done()
		}()
	}
	wg.Wait()
	var count2 int64
	m.Range(func(key string, value interface{}) bool {
		atomic.AddInt64(&count2, 1)
		return true
	})
	m.Delete("600")
	var count int64
	m.Range(func(key string, value interface{}) bool {
		atomic.AddInt64(&count, 1)
		return true
	})

	val, ok := m.Load("500")
	if !ok || reflect.TypeOf(val).Kind().String() != "int" || val.(int) != 1500 {
		t.Fatal("fail")
	}

	_, ok = m.Load("600")
	if ok {
		t.Fatal("fail")
	}

	if m.Len() != 999 || int(count) != m.Len() {
		t.Fatal("fail", m.Len(), count, count2)
	}
}

func testSkipMapToMap(t *testing.T, newset func() orderedskipmap[int]) {
	m := newset()

	// Correctness.
	m.Store(1, "123")
	mm := m.ToMap()
	if !assert.Equal(t, map[int]any{1: "123"}, mm) {
		t.Fatal("invalid")
	}

	m.Store(2, "456")
	mm = m.ToMap()
	if !assert.Equal(t, map[int]any{1: "123", 2: "456"}, mm) {
		t.Fatal("invalid")
	}

	m.Delete(2)
	mm = m.ToMap()
	if !assert.Equal(t, map[int]any{1: "123"}, mm) {
		t.Fatal("invalid")
	}

	m.LoadOrStore(3, "789")
	mm = m.ToMap()
	if !assert.Equal(t, map[int]any{1: "123", 3: "789"}, mm) {
		t.Fatal("invalid")
	}

	m.LoadAndDelete(3)
	mm = m.ToMap()
	if !assert.Equal(t, map[int]any{1: "123"}, mm) {
		t.Fatal("invalid")
	}
}

type BoxedString struct {
	v string
}

func (s BoxedString) MarshalText() ([]byte, error) {
	return []byte(s.v), nil
}

func (s *BoxedString) UnmarshalText(bs []byte) error {
	s.v = string(bs)
	return nil
}

func TestFuncMarshalJSON(t *testing.T) {
	{ // Test TextMarshaler for builtin map.
		expect := `{"foo":"bar"}`
		s := map[BoxedString]string{
			{"foo"}: "bar",
		}
		bs, err := json.Marshal(s)
		if err != nil {
			t.Error(err)
		} else if string(bs) != expect {
			t.Errorf("except JSON %s, get: %s", expect, string(bs))
		}
	}

	{ // Test TextMarshaler.
		expect := `{"foo":"bar"}`
		s := NewFunc[BoxedString, string](func(a, b BoxedString) bool { return a.v < b.v })
		s.Store(BoxedString{"foo"}, "bar")
		bs, err := json.Marshal(s)
		if err != nil {
			t.Error(err)
		} else if string(bs) != expect {
			t.Errorf("except JSON %s, get: %s", expect, string(bs))
		}
	}

	{ // Test string variant.
		expect := `{"foo":"bar"}`
		type MyString string
		s := NewFunc[MyString, string](func(a, b MyString) bool { return a < b })
		s.Store("foo", "bar")
		bs, err := json.Marshal(s)
		if err != nil {
			t.Error(err)
		} else if string(bs) != expect {
			t.Logf("except JSON %s, get: %s", expect, string(bs))
		}
	}

	var nilMap *FuncMap[int, any]
	testSkipMapIntMarshalJSONNil(t, func() anyskipmap[int] { return nilMap })
	testSkipMapIntMarshalJSON(t, func() anyskipmap[int] { return NewFunc[int, any](func(a, b int) bool { return a < b }) }, false)
	testSkipMapIntMarshalJSON(t, func() anyskipmap[uint] { return NewFunc[uint, any](func(a, b uint) bool { return a < b }) }, false)
	testSkipMapStringMarshalJSON(t, func() anyskipmap[string] { return NewFunc[string, any](func(a, b string) bool { return a < b }) }, false)

	testSkipMapIntMarshalJSON(t, func() anyskipmap[int] { return NewFunc[int, any](func(a, b int) bool { return a > b }) }, true)
	testSkipMapIntMarshalJSON(t, func() anyskipmap[uint] { return NewFunc[uint, any](func(a, b uint) bool { return a > b }) }, true)
	testSkipMapStringMarshalJSON(t, func() anyskipmap[string] { return NewFunc[string, any](func(a, b string) bool { return a > b }) }, true)
}

func TestFuncUnmarshalJSON(t *testing.T) {
	{ // Test TextUnmarshaler for builtin map.
		data := []byte(`{"foo":"bar"}`)
		s := map[BoxedString]string{}
		if err := json.Unmarshal(data, &s); err != nil {
			t.Error(err)
		} else if len(s) != 1 || s[BoxedString{"foo"}] != "bar" {
			bs, _ := json.Marshal(s)
			t.Errorf("expect %s, get : %s", string(data), string(bs))
		}
	}

	{ // Test TextUnmarshaler.
		data := []byte(`{"foo":"bar"}`)
		s := NewFunc[BoxedString, string](func(a, b BoxedString) bool { return a.v < b.v })
		if err := json.Unmarshal(data, s); err != nil {
			t.Error(err)
		} else if s.Len() != 1 || goption.Of(s.Load(BoxedString{"foo"})).Value() != "bar" {
			bs, _ := json.Marshal(s)
			t.Errorf("expect %s, get : %s", string(data), string(bs))
		}
	}

	{ // Test string variant.
		data := []byte(`{"foo":"bar"}`)
		type MyString string
		s := NewFunc[MyString, string](func(a, b MyString) bool { return a < b })
		if err := json.Unmarshal(data, s); err != nil {
			t.Error(err)
		} else if s.Len() != 1 || goption.Of(s.Load("foo")).Value() != "bar" {
			bs, _ := json.Marshal(s)
			t.Errorf("expect %s, get : %s", string(data), string(bs))
		}
	}

	testSkipMapIntUnmarshalJSON(t, func() anyskipmap[int] { return NewFunc[int, any](func(a, b int) bool { return a < b }) })
	testSkipMapIntUnmarshalJSON(t, func() anyskipmap[uint] { return NewFunc[uint, any](func(a, b uint) bool { return a < b }) })
	testSkipMapStringUnmarshalJSON(t, func() anyskipmap[string] { return NewFunc[string, any](func(a, b string) bool { return a < b }) })
}

func TestOrderedMarshalJSON(t *testing.T) {
	var nilMap *OrderedMap[int, any]
	testSkipMapIntMarshalJSONNil(t, func() anyskipmap[int] { return nilMap })
	testSkipMapIntMarshalJSON(t, func() anyskipmap[int] { return New[int, any]() }, false)
	testSkipMapIntMarshalJSON(t, func() anyskipmap[uint] { return New[uint, any]() }, false)
	testSkipMapStringMarshalJSON(t, func() anyskipmap[string] { return New[string, any]() }, false)
}

func TestOrderedUnmarshalJSON(t *testing.T) {
	testSkipMapIntUnmarshalJSON(t, func() anyskipmap[int] { return New[int, any]() })
	testSkipMapIntUnmarshalJSON(t, func() anyskipmap[uint] { return New[uint, any]() })
	testSkipMapStringUnmarshalJSON(t, func() anyskipmap[string] { return New[string, any]() })
}

func TestOrderedDescMarshalJSON(t *testing.T) {
	var nilMap *OrderedMapDesc[int, any]
	testSkipMapIntMarshalJSONNil(t, func() anyskipmap[int] { return nilMap })
	testSkipMapIntMarshalJSON(t, func() anyskipmap[int] { return NewDesc[int, any]() }, true)
	testSkipMapIntMarshalJSON(t, func() anyskipmap[uint] { return NewDesc[uint, any]() }, true)
	testSkipMapStringMarshalJSON(t, func() anyskipmap[string] { return NewDesc[string, any]() }, true)
}

func TestOrderedDescUnmarshalJSON(t *testing.T) {
	testSkipMapIntUnmarshalJSON(t, func() anyskipmap[int] { return NewDesc[int, any]() })
	testSkipMapIntUnmarshalJSON(t, func() anyskipmap[uint] { return NewDesc[uint, any]() })
	testSkipMapStringUnmarshalJSON(t, func() anyskipmap[string] { return NewDesc[string, any]() })
}

func testSkipMapIntMarshalJSONNil[T int](t *testing.T, newset func() anyskipmap[T]) {
	expect := `null`
	m := newset()
	if bs, err := json.Marshal(m); err != nil {
		t.Error(err)
	} else if string(bs) != expect {
		t.Fatalf("except %s, get: %s", expect, string(bs))
	}
}

func testSkipMapIntMarshalJSON[T int | uint](t *testing.T, newset func() anyskipmap[T], desc bool) {
	{ // Test empty map.
		expect := `{}`
		m := newset()
		if bs, err := json.Marshal(m); err != nil {
			t.Error(err)
		} else if string(bs) != expect {
			t.Fatalf("except %s, get: %s", expect, string(bs))
		}
	}

	{ // Test empty map.
		var expect string
		if desc {
			expect = `{"3":"condy","2":"bob","1":"alice"}`
		} else {
			expect = `{"1":"alice","2":"bob","3":"condy"}`
		}
		m := newset()
		m.Store(1, "alice")
		m.Store(2, "bob")
		m.Store(3, "condy")
		if bs, err := json.Marshal(m); err != nil {
			t.Error(err)
		} else if string(bs) != expect {
			t.Fatalf("except %s, get: %s", expect, string(bs))
		}
	}
}

func testSkipMapStringMarshalJSON(t *testing.T, newset func() anyskipmap[string], desc bool) {
	{ // Test empty map.
		expect := `{}`
		m := newset()
		if bs, err := json.Marshal(m); err != nil {
			t.Error(err)
		} else if string(bs) != expect {
			t.Fatalf("except %s, get: %s", expect, string(bs))
		}
	}

	{
		var expect string
		if desc {
			expect = `{"3":"condy","2":"bob","1":"alice"}`
		} else {
			expect = `{"1":"alice","2":"bob","3":"condy"}`
		}
		m := newset()
		m.Store("1", "alice")
		m.Store("2", "bob")
		m.Store("3", "condy")
		if bs, err := json.Marshal(m); err != nil {
			t.Error(err)
		} else if string(bs) != expect {
			t.Fatalf("except %s, get: %s", expect, string(bs))
		}
	}
}

func testSkipMapStringUnmarshalJSON(t *testing.T, newset func() anyskipmap[string]) {
	{ // Test empty map.
		data := []byte(`{}`)
		m := newset()
		if err := json.Unmarshal(data, m); err != nil {
			t.Error(err)
		} else if m.Len() != 0 {
			t.Fatal("except a empty map")
		}
	}

	{
		data := []byte(`{"1":"alice","2":"bob","3":"condy"}`)
		m := newset()
		if err := json.Unmarshal(data, m); err != nil {
			t.Error(err)
		} else if m.Len() != 3 ||
			goption.Of(m.Load("1")).Value() != "alice" ||
			goption.Of(m.Load("2")).Value() != "bob" ||
			goption.Of(m.Load("3")).Value() != "condy" {
			bs, _ := json.Marshal(m)
			t.Fatalf("expect: %s", string(bs))
		}
	}
}

func testSkipMapIntUnmarshalJSON[T int | uint](t *testing.T, newset func() anyskipmap[T]) {
	{ // Test empty map.
		data := []byte(`{}`)
		m := newset()
		if err := json.Unmarshal(data, m); err != nil {
			t.Error(err)
		} else if m.Len() != 0 {
			t.Fatal("except a empty map")
		}
	}

	{
		data := []byte(`{"1":"alice","2":"bob","3":"condy"}`)
		m := newset()
		if err := json.Unmarshal(data, m); err != nil {
			t.Error(err)
		} else if m.Len() != 3 ||
			goption.Of(m.Load(1)).Value() != "alice" ||
			goption.Of(m.Load(2)).Value() != "bob" ||
			goption.Of(m.Load(3)).Value() != "condy" {
			bs, _ := json.Marshal(m)
			t.Fatalf("expect: %s", string(bs))
		}
	}
}
