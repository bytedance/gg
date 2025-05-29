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

package gsync

import (
	"fmt"
)

func ExamplePool() {
	pool := Pool[*int]{
		New: func() *int {
			i := 1
			return &i
		},
	}
	a := pool.Get()
	fmt.Println(*a) // 1
	*a = 2
	pool.Put(a)
	_ = *pool.Get() // possible output: 1 or 2

	// Output:
	// 1
}
