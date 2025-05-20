package iter

import (
	"container/list"
	"fmt"
	"strconv"

	"github.com/bytedance/gg/gvalue"
)

type listIter[T any] struct {
	e *list.Element
}

func FromList[T any](l *list.List) Iter[T] {
	return &listIter[T]{l.Front()}
}

func (i *listIter[T]) Next(n int) []T {
	var next []T
	j := 0
	for i.e != nil {
		next = append(next, i.e.Value.(T))
		i.e = i.e.Next()
		j++
		if n != ALL && j >= n {
			break
		}
	}
	return next
}

// Iter for container list.
func ExampleIter_impl() {
	l := list.New()
	l.PushBack(0)
	l.PushBack(1)
	l.PushBack(2)

	s := ToSlice(
		Filter(gvalue.IsNotZero[int],
			FromList[int](l)))
	fmt.Println(s)

	// Output:
	// [1 2]
}

// Convert an int slice to string slice.
func Example() {
	s := ToSlice(
		Map(strconv.Itoa,
			Filter(gvalue.IsNotZero[int],
				FromSlice([]int{0, 1, 2, 3, 4}))))
	fmt.Printf("%q\n", s)

	// Output:
	// ["1" "2" "3" "4"]
}
