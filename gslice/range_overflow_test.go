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
	"testing"
)

func TestRangeWithStepAcrossIntegerLimits(t *testing.T) {
	got := RangeWithStep[int8](math.MinInt8, math.MaxInt8, 1)
	if len(got) != 255 || got[0] != math.MinInt8 || got[len(got)-1] != math.MaxInt8-1 {
		t.Fatalf("RangeWithStep[int8](MinInt8, MaxInt8, 1): len=%d first=%d last=%d", len(got), got[0], got[len(got)-1])
	}

	got16 := RangeWithStep[int16](math.MinInt16, math.MaxInt16, 1)
	if len(got16) != 65535 || got16[0] != math.MinInt16 || got16[len(got16)-1] != math.MaxInt16-1 {
		t.Fatalf("RangeWithStep[int16](MinInt16, MaxInt16, 1): len=%d first=%d last=%d", len(got16), got16[0], got16[len(got16)-1])
	}
}
