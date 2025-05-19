# tuple - Tuple Types

- **Import Path**

    `import "github.com/bytedance/gg/collection/tuple"`


## Quick Start

```go
package main

import (
        "fmt"

        "github.com/bytedance/gg/collection/tuple"
)

func main() {
        addr := tuple.Make2("localhost", 8080)
        fmt.Printf("%s:%d\n", addr.First, addr.Second)
        // Output:
        // localhost:8080
}
```

## Limitation

### No type inference for composite literals

According to the [Type Parameters Proposal](https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#type-inference-for-composite-literals),
type inference for composite literals is not support at least in Go1.18.
So we can not easily cast a function to “partial application”-able type.

Fortunately, [Type inference for functions is supported](https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#type-inference),
use `MakeN` functions for casting without specify all type parameters explicitly.

### Arity Limitation

If you have a need for n-ary (where n > 10) tuple, please file an issue.
