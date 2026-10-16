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

package iter

import (
	"math"
	"reflect"
	"testing"
)

func TestRangeWithStepIntegerOverflow(t *testing.T) {
	ascending := ToSlice(RangeWithStep[int8](math.MinInt8, math.MaxInt8, 1))
	if len(ascending) != 255 || ascending[0] != math.MinInt8 || ascending[len(ascending)-1] != math.MaxInt8-1 {
		t.Fatalf("int8 ascending full range: len=%d first=%d last=%d", len(ascending), ascending[0], ascending[len(ascending)-1])
	}

	descending := ToSlice(RangeWithStep[int8](math.MaxInt8, math.MinInt8, -1))
	if len(descending) != 255 || descending[0] != math.MaxInt8 || descending[len(descending)-1] != math.MinInt8+1 {
		t.Fatalf("int8 descending full range: len=%d first=%d last=%d", len(descending), descending[0], descending[len(descending)-1])
	}

	if got := ToSlice(RangeWithStep[int32](math.MinInt32, math.MaxInt32, 1<<30)); !reflect.DeepEqual(got, []int32{math.MinInt32, -1 << 30, 0, 1 << 30}) {
		t.Fatalf("int32 cross-zero range = %v", got)
	}

	if got := ToSlice(RangeWithStep[int8](120, 127, 5)); !reflect.DeepEqual(got, []int8{120, 125}) {
		t.Fatalf("int8 overflow after last value = %v", got)
	}
	if got := ToSlice(RangeWithStep[int8](-120, -128, -5)); !reflect.DeepEqual(got, []int8{-120, -125}) {
		t.Fatalf("int8 underflow after last value = %v", got)
	}
}

func TestRangeWithStepPartialReadsAcrossOverflowingInterval(t *testing.T) {
	i := RangeWithStep[int8](math.MinInt8, math.MaxInt8, 1)
	var got []int8
	for {
		part := i.Next(7)
		if len(part) == 0 {
			break
		}
		got = append(got, part...)
	}
	if len(got) != 255 || got[0] != math.MinInt8 || got[len(got)-1] != math.MaxInt8-1 {
		t.Fatalf("partial int8 full range: len=%d first=%d last=%d", len(got), got[0], got[len(got)-1])
	}
	if got := i.Next(1); got != nil {
		t.Fatalf("exhausted iterator returned %v", got)
	}
}

func TestRangeWithStepRejectsNaNAndStopsOnStagnation(t *testing.T) {
	for _, i := range []Iter[float64]{
		RangeWithStep(math.NaN(), 1.0, 0.1),
		RangeWithStep(0.0, math.NaN(), 0.1),
		RangeWithStep(0.0, 1.0, math.NaN()),
	} {
		if got := i.Next(ALL); len(got) != 0 {
			t.Fatalf("NaN range returned %v", got)
		}
	}

	// At this magnitude, adding 1 does not change the float64 value. The
	// iterator must terminate instead of yielding the same value forever.
	start := float64(1 << 53)
	got := ToSlice(RangeWithStep(start, start+2, 1.0))
	if !reflect.DeepEqual(got, []float64{start}) {
		t.Fatalf("stagnating float range = %v", got)
	}
}
