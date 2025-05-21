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

package fastrand

import (
	"testing"
)

func TestAll(t *testing.T) {
	_ = Uint32()

	p := make([]byte, 1000)
	n, err := Read(p)
	if n != len(p) || err != nil || (p[0] == 0 && p[1] == 0 && p[2] == 0) {
		t.Fatal()
	}

	a := Perm(100)
	for i := range a {
		var find bool
		for _, v := range a {
			if v == i {
				find = true
			}
		}
		if !find {
			t.Fatal()
		}
	}

	Shuffle(len(a), func(i, j int) {
		a[i], a[j] = a[j], a[i]
	})
	for i := range a {
		var find bool
		for _, v := range a {
			if v == i {
				find = true
			}
		}
		if !find {
			t.Fatal()
		}
	}

	Shuffle2(a)
	for i := range a {
		var find bool
		for _, v := range a {
			if v == i {
				find = true
			}
		}
		if !find {
			t.Fatal()
		}
	}
}
