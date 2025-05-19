# gfunc - Operations of Functions

Package *gfunc* provides operations of functions.

* **Import Path**

    `import "github.com/bytedance/gg/gfunc"`

Package *partial* implements [Partial Application](https://en.wikipedia.org/wiki/Partial_application) of functions.

Partial Application refers to the process of fixing a number of arguments to a function,
producing another function of smaller arity.

## Quick Start


1. Import `"github.com/bytedance/gg/gfunc"`.
2. Create `FuncN` (function with N parameters, such as `Func2`)
3. Use method `Func2.Partial` or `Func2.PartialR` to bind parameters, producing a `FuncN-1` (such as `Func1`)

```go
package main

import (
        "fmt"

        "github.com/bytedance/gg/gfunc"
)

func main() {
        f := func(a, b int) int {
                return a + b
        }
        add := gfunc.Partial2(f)   // Cast f to "partial application"-able function
        add1 := add.Partial(1)    // Bind argument a to 1
        fmt.Println(add1(0))      // 1 + 0 = 1
        fmt.Println(add1(1))      // add1 can be reused, 1 + 1 = 2
        add1n2 := add1.Partial(2) // Bind argument b to 2, all arguments are fixed
        fmt.Println(add1n2())     // 1 + 2 = 3
        // Output:
        // 1
        // 2
        // 3
}
```

This is just an example.
In real world, we no need to write such `add` function for specify type,
`gvalue` provides a generics `gvalue.Add` functions and
various utils for writing generic code.

## Limitation

### No type inference for composite literals

According to the [Type Parameters Proposal](https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#type-inference-for-composite-literals), Type inference for composite
literals is not support at least in Go1.18. So we can not easily cast a
function to “partial application”-able type.

Fortunately, [Type inference for functions is supported](https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#type-inference), we provides
`MakeN` functions for casting without specify all type parameters explicitly.

### Arity Limitation

If you have a need for n-ary (where n > 10) functions, please file an issue.
