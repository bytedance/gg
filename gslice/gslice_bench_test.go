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
	"math/rand"
	"strconv"
	"testing"

	"github.com/bytedance/gg/collection/set"
	"github.com/bytedance/gg/internal/iter"
)

func BenchmarkMap10(b *testing.B) {
	benchmarkMapN(b, 10)
}

func BenchmarkMap100(b *testing.B) {
	benchmarkMapN(b, 100)
}

func BenchmarkMap1000(b *testing.B) {
	benchmarkMapN(b, 100)
}

func benchmarkMapN(b *testing.B, n int) {
	s := []int{}
	for i := 0; i < n; i++ {
		s = append(s, i)
	}

	b.Run("baseline", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r := make([]string, 0, n)
			for _, v := range s {
				r = append(r, strconv.Itoa(v))
			}
			_ = r
		}
	})
	b.Run("gslice", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Map(s, strconv.Itoa)
		}
	})
	b.Run("iter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = iter.ToSlice(iter.Map(strconv.Itoa, iter.StealSlice(s)))
		}
	})
	b.Run("reflect", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = reflectMap(s, func(i any) any { return strconv.Itoa(i.(int)) }).([]string)
		}
	})
}

func BenchmarkShuffle(b *testing.B) {
	b.Run("gslice", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s := iter.ToSlice(iter.Range(0, 100))
			Shuffle(s)
			_ = s
		}
	})
	b.Run("math/rand", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s := iter.ToSlice(iter.Range(0, 100))
			rand.Shuffle(len(s), func(i, j int) {
				s[i], s[j] = s[j], s[i]
			})
			_ = s
		}
	})
}

func BenchmarkShuffle_Parallel(b *testing.B) {
	b.Run("gslice", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				s := iter.ToSlice(iter.Range(0, 100))
				Shuffle(s)
				_ = s
			}
		})
	})
	b.Run("math/rand", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				s := iter.ToSlice(iter.Range(0, 100))
				rand.Shuffle(len(s), func(i, j int) {
					s[i], s[j] = s[j], s[i]
				})
				_ = s
			}
		})
	})
}

func oldUnion[S ~[]T, T comparable](ss ...S) S {
	if len(ss) == 0 {
		return S{}
	}
	if len(ss) == 1 {
		return Uniq(ss[0])
	}
	members := set.New[T]()
	ret := S{} // TODO: Guess a cap.
	for _, s := range ss {
		for _, v := range s {
			if members.Add(v) {
				ret = append(ret, v)
			}
		}
	}
	return ret
}

func BenchmarkUnion(b *testing.B) {
	// 1. all different
	ss1 := [][]int{
		Range(0, 10),
		Range(10, 20),
	}
	b.Run("new-union-diff-2-10", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Union(ss1...)
		}
	})
	b.Run("old-union-diff-2-10", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			oldUnion(ss1...)
		}
	})
	ss2 := [][]int{
		Range(0, 100),
		Range(100, 200),
		Range(200, 300),
		Range(300, 400),
		Range(400, 500),
	}
	b.Run("new-union-diff-5-100", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Union(ss2...)
		}
	})
	b.Run("old-union-diff-5-100", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			oldUnion(ss2...)
		}
	})

	// 2. all same
	ss3 := [][]int{
		Repeat(0, 10),
		Repeat(0, 10),
	}
	b.Run("new-union-same-2-10", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Union(ss3...)
		}
	})
	b.Run("old-union-same-2-10", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			oldUnion(ss3...)
		}
	})
	ss4 := [][]int{
		Repeat(0, 100),
		Repeat(0, 100),
		Repeat(0, 100),
		Repeat(0, 100),
		Repeat(0, 100),
	}
	b.Run("new-union-same-5-100", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Union(ss4...)
		}
	})
	b.Run("old-union-same-5-100", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			oldUnion(ss4...)
		}
	})

	// 3. half different
	ss5 := [][]int{
		Range(0, 100),
		Range(0, 100),
	}
	b.Run("new-union-half-2-100", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Union(ss5...)
		}
	})
	b.Run("old-union-half-2-100", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			oldUnion(ss5...)
		}
	})
}
