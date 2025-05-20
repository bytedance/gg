package skipset

import (
	"math"
	"sync"
	"testing"

	"github.com/bytedance/gg/internal/fastrand"
)

const (
	initsize = 1 << 10 // for `contains` `1Remove9Add90Contains` `1Range9Remove90Add900Contains`
	randN    = math.MaxUint32
)

func BenchmarkInt64(b *testing.B) {
	all := []benchTask[int64]{{
		name: "skipset", New: func() anyskipset[int64] {
			return New[int64]()
		}}}
	all = append(all, benchTask[int64]{
		name: "skipset(func)", New: func() anyskipset[int64] {
			return NewFunc(func(a, b int64) bool {
				return a < b
			})
		}})
	all = append(all, benchTask[int64]{
		name: "sync.Map", New: func() anyskipset[int64] {
			return new(anySyncMap[int64])
		}})
	rng := fastrand.Int63
	benchAdd(b, rng, all)
	bench30Add70Contains(b, rng, all)
	bench1Remove9Add90Contains(b, rng, all)
	bench1Range9Remove90Add900Contains(b, rng, all)
}

func benchAdd[T any](b *testing.B, rng func() T, benchTasks []benchTask[T]) {
	for _, v := range benchTasks {
		b.Run("Add/"+v.name, func(b *testing.B) {
			s := v.New()
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					s.Add(rng())
				}
			})
		})
	}
}

func bench30Add70Contains[T any](b *testing.B, rng func() T, benchTasks []benchTask[T]) {
	for _, v := range benchTasks {
		b.Run("30Add70Contains/"+v.name, func(b *testing.B) {
			s := v.New()
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					u := fastrand.Uint32n(10)
					if u < 3 {
						s.Add(rng())
					} else {
						s.Contains(rng())
					}
				}
			})
		})
	}
}

func bench1Remove9Add90Contains[T any](b *testing.B, rng func() T, benchTasks []benchTask[T]) {
	for _, v := range benchTasks {
		b.Run("1Remove9Add90Contains/"+v.name, func(b *testing.B) {
			s := v.New()
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					u := fastrand.Uint32n(100)
					if u < 9 {
						s.Add(rng())
					} else if u == 10 {
						s.Remove(rng())
					} else {
						s.Contains(rng())
					}
				}
			})
		})
	}
}

func bench1Range9Remove90Add900Contains[T any](b *testing.B, rng func() T, benchTasks []benchTask[T]) {
	for _, v := range benchTasks {
		b.Run("1Range9Remove90Add900Contains/"+v.name, func(b *testing.B) {
			s := v.New()
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					u := fastrand.Uint32n(1000)
					if u == 0 {
						s.Range(func(score T) bool {
							return true
						})
					} else if u > 10 && u < 20 {
						s.Remove(rng())
					} else if u >= 100 && u < 190 {
						s.Add(rng())
					} else {
						s.Contains(rng())
					}
				}
			})
		})
	}
}

type benchTask[T any] struct {
	name string
	New  func() anyskipset[T]
}

type anySyncMap[T any] struct {
	data sync.Map
}

func (m *anySyncMap[T]) Add(x T) bool {
	m.data.Store(x, struct{}{})
	return true
}

func (m *anySyncMap[T]) Contains(x T) bool {
	_, ok := m.data.Load(x)
	return ok
}

func (m *anySyncMap[T]) Remove(x T) bool {
	m.data.Delete(x)
	return true
}

func (m *anySyncMap[T]) Range(f func(value T) bool) {
	m.data.Range(func(key, _ any) bool {
		return !f(key.(T))
	})
}

func (m *anySyncMap[T]) Len() int {
	var i int
	m.data.Range(func(_, _ any) bool {
		i++
		return true
	})
	return i
}

func (m *anySyncMap[T]) ToSlice() []T {
	s := make([]T, 0, m.Len())
	m.data.Range(func(key, _ any) bool {
		s = append(s, key.(T))
		return true
	})
	return s
}
