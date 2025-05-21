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

package skipset

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/bytedance/gg/internal/assert"
	"github.com/bytedance/gg/internal/fastrand"
)

func TestOrdered(t *testing.T) {
	testIntSet(t, func() anyskipset[int] {
		return New[int]()
	})
	testIntSetDesc(t, func() anyskipset[int] {
		return NewDesc[int]()
	})
	testStringSet(t, func() anyskipset[string] {
		return New[string]()
	})
}

func TestFunc(t *testing.T) {
	x := NewFunc(func(a, b float64) bool {
		return a < b || (math.IsNaN(a) && !math.IsNaN(b))
	})
	x.Add(math.NaN())
	x.Add(3)
	x.Add(1)
	x.Add(math.NaN())
	x.Add(2)
	x.Add(math.NaN())
	if x.Len() != 4 {
		t.Fatal(x.Len())
	}
	expected := []float64{math.NaN(), 1, 2, 3}
	var i int
	x.Range(func(value float64) bool {
		if i == 0 && !math.IsNaN(value) {
			t.Fatal("invalid range")
		}
		if i >= 1 && value != expected[i] {
			t.Fatal("invalid range")
		}
		i++
		return true
	})
	actual := x.ToSlice()
	equal := func(a, b float64) bool {
		return a == b || (math.IsNaN(a) && math.IsNaN(b))
	}
	for j := 0; j < 4; j++ {
		if !equal(expected[j], actual[j]) {
			t.Fatal("invalid toslice")
		}
	}

	testIntSet(t, func() anyskipset[int] {
		return NewFunc(func(a, b int) bool {
			return a < b
		})
	})
	testIntSetDesc(t, func() anyskipset[int] {
		return NewFunc(func(a, b int) bool {
			return a > b
		})
	})
	testStringSet(t, func() anyskipset[string] {
		return NewFunc(func(a, b string) bool {
			return a < b
		})
	})
}

type anyskipset[T any] interface {
	Add(v T) bool
	Remove(v T) bool
	Contains(v T) bool
	Range(f func(v T) bool)
	Len() int
	ToSlice() []T
}

// Test suites.

func testIntSet(t *testing.T, newset func() anyskipset[int]) {
	// Correctness.
	l := newset()
	if l.Len() != 0 {
		t.Fatal("invalid length")
	}
	if l.Contains(0) {
		t.Fatal("invalid contains")
	}

	if !l.Add(0) || l.Len() != 1 {
		t.Fatal("invalid add")
	}
	if !l.Contains(0) {
		t.Fatal("invalid contains")
	}
	if !l.Remove(0) || l.Len() != 0 {
		t.Fatal("invalid remove")
	}

	if !l.Add(20) || l.Len() != 1 {
		t.Fatal("invalid add")
	}
	if !l.Add(22) || l.Len() != 2 {
		t.Fatal("invalid add")
	}
	if !l.Add(21) || l.Len() != 3 {
		t.Fatal("invalid add")
	}

	expected := []int{20, 21, 22}
	var i int
	l.Range(func(score int) bool {
		if score != expected[i] {
			t.Fatal("invalid range")
		}
		i++
		return true
	})
	if !assert.Equal(t, expected, l.ToSlice()) {
		t.Fatal()
	}

	if !l.Remove(21) || l.Len() != 2 {
		t.Fatal("invalid remove")
	}

	expected = []int{20, 22}
	i = 0
	l.Range(func(score int) bool {
		if score != expected[i] {
			t.Fatal("invalid range")
		}
		i++
		return true
	})
	if !assert.Equal(t, expected, l.ToSlice()) {
		t.Fatal("invalid toslice")
	}

	const num = math.MaxInt16
	// Make rand shuffle array.
	// The testArray contains [1,num]
	testArray := make([]int, num)
	testArray[0] = num + 1
	for i := 1; i < num; i++ {
		// We left 0, because it is the default score for head and tail.
		// If we check the skipset contains 0, there must be something wrong.
		testArray[i] = int(i)
	}
	for i := len(testArray) - 1; i > 0; i-- { // Fisher–Yates shuffle
		j := fastrand.Uint32n(uint32(i + 1))
		testArray[i], testArray[j] = testArray[j], testArray[i]
	}

	// Concurrent add.
	var wg sync.WaitGroup
	for i := 0; i < num; i++ {
		i := i
		wg.Add(1)
		go func() {
			l.Add(testArray[i])
			wg.Done()
		}()
	}
	wg.Wait()
	if l.Len() != int(num) {
		t.Fatalf("invalid length expected %d, got %d", num, l.Len())
	}

	// Don't contains 0 after concurrent addion.
	if l.Contains(0) {
		t.Fatal("contains 0 after concurrent addion")
	}

	// Concurrent contains.
	for i := 0; i < num; i++ {
		i := i
		wg.Add(1)
		go func() {
			if !l.Contains(testArray[i]) {
				wg.Done()
				panic(fmt.Sprintf("add doesn't contains %d", i))
			}
			wg.Done()
		}()
	}
	wg.Wait()

	// Concurrent remove.
	for i := 0; i < num; i++ {
		i := i
		wg.Add(1)
		go func() {
			if !l.Remove(testArray[i]) {
				wg.Done()
				panic(fmt.Sprintf("can't remove %d", i))
			}
			wg.Done()
		}()
	}
	wg.Wait()
	if l.Len() != 0 {
		t.Fatalf("invalid length expected %d, got %d", 0, l.Len())
	}

	// Test all methods.
	const smallRndN = 1 << 8
	for i := 0; i < 1<<16; i++ {
		wg.Add(1)
		go func() {
			r := fastrand.Uint32n(num)
			if r < 333 {
				l.Add(int(fastrand.Uint32n(smallRndN)) + 1)
			} else if r < 666 {
				l.Contains(int(fastrand.Uint32n(smallRndN)) + 1)
			} else if r != 999 {
				l.Remove(int(fastrand.Uint32n(smallRndN)) + 1)
			} else {
				var pre int
				l.Range(func(score int) bool {
					if score <= pre { // 0 is the default value for header and tail score
						panic("invalid content")
					}
					pre = score
					return true
				})
			}
			wg.Done()
		}()
	}
	wg.Wait()

	// Correctness 2.
	var (
		x     = newset()
		y     = newset()
		count = 10000
	)

	for i := 0; i < count; i++ {
		x.Add(i)
	}

	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			x.Range(func(score int) bool {
				if x.Remove(score) {
					if !y.Add(score) {
						panic("invalid add")
					}
				}
				return true
			})
			wg.Done()
		}()
	}
	wg.Wait()
	if x.Len() != 0 || y.Len() != count {
		t.Fatal("invalid length")
	}

	// Concurrent Add and Remove in small zone.
	x = newset()
	var (
		addcount    uint64 = 0
		removecount uint64 = 0
	)

	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 1000; i++ {
				if fastrand.Uint32n(2) == 0 {
					if x.Remove(int(fastrand.Uint32n(10))) {
						atomic.AddUint64(&removecount, 1)
					}
				} else {
					if x.Add(int(fastrand.Uint32n(10))) {
						atomic.AddUint64(&addcount, 1)
					}
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	if addcount < removecount {
		panic("invalid count")
	}
	if addcount-removecount != uint64(x.Len()) {
		panic("invalid count")
	}

	pre := -1
	x.Range(func(score int) bool {
		if score <= pre {
			panic("invalid content")
		}
		pre = score
		return true
	})

	// Correctness 3.
	s1 := newset()
	var s2 sync.Map
	var counter uint64
	for i := 0; i <= 10000; i++ {
		wg.Add(1)
		go func() {
			if fastrand.Uint32n(2) == 0 {
				r := fastrand.Uint32()
				s1.Add(int(r))
				s2.Store(int(r), nil)
			} else {
				r := atomic.AddUint64(&counter, 1)
				s1.Add(int(r))
				s2.Store(int(r), nil)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	s1.Range(func(value int) bool {
		_, ok := s2.Load(value)
		if !ok {
			t.Fatal(value)
		}
		return true
	})
	s2.Range(func(key, value interface{}) bool {
		k := key.(int)
		if !s1.Contains(k) {
			t.Fatal(value)
		}
		return true
	})
}

func testIntSetDesc(t *testing.T, newsetdesc func() anyskipset[int]) {
	s := newsetdesc()
	nums := []int{-1, 0, 5, 12}
	for _, v := range nums {
		s.Add(v)
	}
	expected := []int{12, 5, 0, -1}
	i := 0
	s.Range(func(value int) bool {
		if expected[i] != value {
			t.Fatal("error")
		}
		i++
		return true
	})
	if !assert.Equal(t, expected, s.ToSlice()) {
		t.Fatal("invalid toslice")
	}
}

func testStringSet(t *testing.T, newset func() anyskipset[string]) {
	x := newset()
	if !x.Add("111") || x.Len() != 1 {
		t.Fatal("invalid")
	}
	if !x.Add("222") || x.Len() != 2 {
		t.Fatal("invalid")
	}
	if x.Add("111") || x.Len() != 2 {
		t.Fatal("invalid")
	}
	if !x.Contains("111") || !x.Contains("222") {
		t.Fatal("invalid")
	}
	if !x.Remove("111") || x.Len() != 1 {
		t.Fatal("invalid")
	}
	if !x.Remove("222") || x.Len() != 0 {
		t.Fatal("invalid")
	}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		i := i
		go func() {
			if !x.Add(strconv.Itoa(i)) {
				panic("invalid")
			}
			wg.Done()
		}()
	}
	wg.Wait()

	tmp := make([]int, 0, 100)
	x.Range(func(val string) bool {
		res, _ := strconv.Atoi(val)
		tmp = append(tmp, res)
		return true
	})
	sort.Ints(tmp)
	expected := make([]string, 0, 100)
	for i := 0; i < 100; i++ {
		if i != tmp[i] {
			t.Fatal("invalid")
		}
		expected = append(expected, strconv.Itoa(i))
	}
	sort.Strings(expected)
	if !assert.Equal(t, expected, x.ToSlice()) {
		t.Fatal("invalid toslice")
	}
}
