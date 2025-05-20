# iter - Iterator and Operations


* **Import Path**

    `import "github.com/bytedance/gg/internal/iter"`


Package *iter* provides definition of generic iterator `Iter` and high-order functions
for [operating](#operations) it. Most of the operations are lazy.


* Iterator helps us you iterating various data structure in same way.
* “Operation” is operation for processing elements of iterator.
* “Lazy” means that the evaluation of an operation is delayed until its
value is needed

We can create iterator from Golang’s slice, map or even channel,
and apply various [operations](#operations) on it.

Beside, user can [Implement your custom Iter](#operation-custom-data).

This packages is greatly inspired by Haskell.

## Quick Start


1. Import `"github.com/bytedance/gg/internal/iter"`.
2. Use `FromSlice` source to create a iterator of int slice.
3. Use `Filter` operation to filter the zero values.
4. Use `Map` operation to convert int to string.
5. Use `ToSlice` sink to convert elements of iterator to slice, evaluation is done here.

```go
package main

import (
        "fmt"

        "github.com/bytedance/gg/gvalue"
        "github.com/bytedance/gg/internal/iter"
)

func main() {
        s := iter.ToSlice(
                iter.Map(strconv.Itoa,
                        iter.Filter(gvalue.IsZero[int],
                                iter.FromSlice([]int{0, 1, 2, 3, 4}))))
        fmt.Printf("%q\n", s)

        // Output:
        // ["1" "2" "3" "4"]
}
```

## API References

### Sources

Source is used to craete iterator from data (like slice, map, and etc).
Source functions are usually named “FromXxx” if “Xxx” is a noun.

| Data                                               | Sources                                            |
| -------------------------------------------------- | -------------------------------------------------- |
| Slice                                              | `FromSlice`                                          |
| Map                                                | `FromMap`, `FromMapKeys`, `FromMapValues`                |
| Channel                                            | `FromChan`                                           |
| Others                                             | `Range`, `RangeWithStep`, `Repeat`                       |
### Sinks

Sink is used to extract iterator to data (like slice, map, and etc).
Sink functions are usually named “ToXxx”.

| Data                                               | Sinks                                              |
| -------------------------------------------------- | -------------------------------------------------- |
| Slice                                              | `ToSlice`                                            |
| Map                                                | `ToMap`                                              |
| Channel                                            | `ToChan`, `ToBufferedChan`                             |
### Operations

Operation is used to process elements of iterator.

#### General

These operations can be applied on iterator of any type:

| Kind                                               | Operations                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| -------------------------------------------------- |--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Map                                                | `Map`, `MapInplace`, `FlatMap`, `Cast`, `FilterMap`                                                                                                                                                                                                                                                                                                                                                                             |
| Filter                                             | `Filter`, `Find`, `DistinctWith`, `DistinctOrderedWith`                                                                                                                                                                                                                                                                                                                                                                                                                           |
| Fold                                               | `Fold`, `FoldWith`, `Count`, `All`, `Any`                                                                                                                                                                                                                                                                                                                                                                                                 |
| Iteration                                          | `ForEach`, `ForEachIndexed`                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               |
| Substream                                          | `Take`, `Drop`, `TakeWhile`, `DropWhile`                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| List                                               | `Head`, `Last`, `Reverse`, `Prepend`, `Append`, `Concat`, `Intersperse`, `SortWith`, `At`, `Chunk` |
| Zipping                                            | `Zip`, `ZipWith`                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |
#### Type-specified

These operations can only be applied on iterator of specified type:

| Type                                               | Usage                                              |
| -------------------------------------------------- | -------------------------------------------------- |
| `comparable`                                         | `Contains`, `Distinct`, `DistinctOrdered`, `Remove`, `RemoveN` |
| `constraints.Ordered`                                | `Max`, `Min`, `Sort`                                     |
| `~bool`                                              | `And`, `Or`                                            |
| `~string`                                            | `Join`                                               |
| `constraints.Integer` \| `constraints.Float`           | `Sum`, `Avg`                                           |
## Tips

### Use method chaining

As we see, operations are functions. We have to nest function calls when we have
multiple operations to be applied (as we do in [Quick Start](#quick-start)).

`Method chaining` is a convenient style for simplifying nested calls.
[stream - Stream Processing](../stream/README.md) provides a series of iterator wrappers with method
chaining support.

### Operating Custom Data

If you want to apply operation on data structures which are not included in Sources
just implement an Iter interface for them.
Take [container/list](https://pkg.go.dev/container/list) as an example:

Implement `Iter` interface for `container/list`:

```go
import (
        "container/list"
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
```

Then you can apply various operations:

```go
l := list.New()
l.PushBack(0)
l.PushBack(1)
l.PushBack(2)

l := list.New()
l.PushBack(0)
l.PushBack(1)
l.PushBack(2)

i := FromList(l)
s := ToSlice(Filter(gvalue.IsZero[int], i))
fmt.Println(s)

// Output:
// [1 2]
```

### Partial Application

[partial - Partial Application](../../gfunc/README.md) implements [Partial Application](https://en.wikipedia.org/wiki/Partial_application)
which can simplify our use of higher-order functions:

```go
add := gvalue.Add[int]                 // instantiate a int version of Add function
add1 := gfunc.Partial2(add).Partial(1) // bind the first argument to 1
s := ToSlice(
        Map(add1,
            FromSlice([]int{1, 2, 3, 4})))

fmt.Println(s)
// Output:
// [2 3 4 5]
```

### Generics Utils

`gvalue` provides a lots of helpful functions and predicates
with generics support.
