# stream - Stream Processing

**WARNING**: This package is **experimental** and may change in the future.

Package *stream* provides a `Stream` type and its [variants](#stream-variants).
All of them are wrappers of [iter - Iterator and Operations](../../iter/README.md).
With these wrappers we can call [operations](#sources-operations-and-sinks)
in [method chaining](https://en.wikipedia.org/wiki/Method_chaining) style.

## Quick Start


1. Import `"github.com/bytedance/gg/internal/stream"`.
2. Use `FromSlice` to construct a stream of int slice.
3. Use `Filter` to filter the zero values.
4. Use `ToSlice` to convert filtered stream to slice, evaluation is done here.

```go
package main

import (
        "fmt"

        "github.com/bytedance/gg/gvalue"
        "github.com/bytedance/gg/internal/stream"
)

func main() {
        s := stream.FromSlice([]int{0, 1, 2, 3, 4}).    // Construct a stream from int slice
                Filter(gvalue.IsNotZero[int]).              // Filter zero value lazily
                ToSlice()                               // Evaluate and convert back to slice

        fmt.Println(s)
        // Output:
        // [1 2 3 4]
}
```

## API References

### Sources, Operations and Sinks

All of sources, operations and sinks have their corresponding functions in
[iter - Iterator and Operations](../../iter/README.md). For example:


* `Stream.Map` is corresponding to `iter.Map`
* `String.Join` is corresponding to `iter.Join`

So we won’t repeat them here.

### Stream Variants

`Stream` has different variants depending on the element type.
Variants may have additional sources or operations or sinks.

| Type                                               | Variant                                            |
| -------------------------------------------------- | -------------------------------------------------- |
| `any`                                                | `Stream`                                             |
| `comparable`                                         | `Comparable`                                         |
| `constraints.Ordered`                                | `Orderable`                                          |
| `~bool`                                              | `Bool`                                               |
| `~string`                                            | `String`                                             |
| `map[comparable]any`                                 | `KV`                                                 |
| `map[constraints.Ordered]any`                        | `OrderableKV`                                        |
## Limitation

### Can not transform Stream from one type to another

We known `iter.Map` can transform an `iter.Iter`
from type `F` to `T`.  But in this package, `Stream.Map` can’t do.

According to the [Type Parameters Proposal](https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#No-parameterized-methods), method can not have additional
type parameter. The following code is invalid for now:

```go
func (s *Stream[F]) Map[T any](f func(F) T) *Stream[T]
// ERROR: methods cannot have type parameters
```

This means that we can not transform a stream from one type to another by
method chaining.

[golang/go#49085](https://github.com/golang/go/issues/49085) discussed this matter, but there is no one given a good
plan yet.

### Lacking a better way for defining type-specialized variants

`fold` operation can be used for `any` type, but `join` operation is only for `~string`.
We can not add a `Join` method for `Stream` because its type constraint is
already `[T any]` but not `[T ~string]`.

So we create a variant `String[T ~string]` which has a `Stream` embedded in,
then we implement a `Join` method for it.
But the parameters of the method inherited from `Stream` is still `stream.Stream[T]`,
not `stream.String[T]`. For now, we rewrite these methods by code generation.
