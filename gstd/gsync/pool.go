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

import "sync"

// Pool wraps [sync.Pool].
type Pool[T any] struct {
	New     func() T
	p       sync.Pool
	newOnce sync.Once
}

func (p *Pool[T]) init() {
	p.newOnce.Do(func() {
		p.p.New = func() any {
			return p.New()
		}
	})
}

// Get wraps [sync.Pool.Get].
func (p *Pool[T]) Get() T {
	p.init()
	return p.p.Get().(T)
}

// Put wraps [sync.Pool.Put].
func (p *Pool[T]) Put(x T) {
	p.init()
	p.p.Put(x)
}
