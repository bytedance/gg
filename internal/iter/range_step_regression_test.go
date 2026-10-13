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
	"testing"
)

// Regression for review feedback on bytedance/gg#45: a NaN produced by the
// iteration itself (e.g. -Inf + +Inf) must terminate the iterator instead of
// appending NaN forever. Take is bounded, so even if this guard regresses the
// test cannot hang or OOM the runner.
func TestRangeWithStepNaNMidIterationTerminates(t *testing.T) {
	it := RangeWithStep(math.Inf(-1), 0.0, math.Inf(1))
	got := Take(3, it)
	if len(got) != 1 || !math.IsInf(got[0], -1) || math.IsNaN(got[0]) {
		t.Fatalf("expected exactly [-Inf] before exhaustion, got %v", got)
	}
	if rest := ToSlice(it); len(rest) != 0 {
		t.Fatalf("expected exhaustion after -Inf, remaining %v", rest)
	}

	// Reverse direction: +Inf start with a -Inf step.
	it2 := RangeWithStep(math.Inf(1), 0.0, math.Inf(-1))
	got2 := Take(3, it2)
	if len(got2) != 1 || !math.IsInf(got2[0], +1) || math.IsNaN(got2[0]) {
		t.Fatalf("expected exactly [+Inf] before exhaustion, got %v", got2)
	}
}

// Regression for review feedback on bytedance/gg#45: after a batch that
// consumes the range exactly, the following Next must return nil (the
// documented empty-iterator contract), not a non-nil empty slice.
func TestRangeExactBatchConsumptionReturnsNil(t *testing.T) {
	it := Range(0, 2)
	if first := it.Next(2); len(first) != 2 || first[0] != 0 || first[1] != 1 {
		t.Fatalf("expected [0 1], got %v", first)
	}
	if second := it.Next(1); second != nil {
		t.Fatalf("expected nil after exact batch consumption, got %#v", second)
	}
	if third := it.Next(1); third != nil {
		t.Fatalf("expected nil on repeated read, got %#v", third)
	}
}
