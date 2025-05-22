# gg: Go Generics

[![GoDoc](https://godoc.org/github.com/bytedance/gg?status.svg)](https://godoc.org/github.com/bytedance/gg)
[![Go Report Card](https://goreportcard.com/badge/github.com/bytedance/gg)](https://goreportcard.com/report/github.com/bytedance/gg)
[![Go Coverage](https://codecov.io/gh/bytedance/gg/branch/main/graph/badge.svg)](https://codecov.io/gh/bytedance/gg)
[![License](https://img.shields.io/github/license/bytedance/gg)](https://github.com/bytedance/gg/blob/main/LICENSE)

English | [简体中文](README.zh-CN.md)

🔥`bytedance-gg` is a basic library of generics for Go language developed by ByteDance. It is based on the Go 1.18+ generic features and provides efficient, type-safe and rich generic data structures and tool functions.

❓**Why this name?**

Take the first letter of **G**o **G**enerics, short and simple.

❓Why choose gg?

- Stable and reliable: It is a necessary tool library for ByteDance R&D team, and it has 1w+ repository references inside.
- Easy to use: With the design principle of simplicity and self-consistent, subcontracted according to functions, modular, semantic intuitive and unified, and low learning cost.
- High Performance: Provides high-performance concurrent data structures, with performance 10+ times faster than standard library.
- No three-party dependencies: Generic libraries will not introduce any three-party dependencies.
- Version control: Follow the SemVer, guaranteeing backward compatibility.

## 🚀 Install

```sh
go get github.com/bytedance/gg
```

## 🔎 Table of contents

- [Generic Functional Programming](#-generic-functional-programming)
    - [goption](#goption)：Option type, simplifying the processing of `(T, bool)`
    - [gresult](#gresult)：Result type, simplifying the processing of `(T, error)`
- [Generic Data Processing](#-generic-data-processing)
    - [gcond](#gcond)：Conditional operation
    - [gvalue](#gvalue)：Processing value `T`
    - [gptr](#gptr)：Processing pointer `*T`
    - [gslice](#gslice)：Process slice `[]T`
    - [gmap](#gmap)：Processing map `map[K]V`
    - [gfunc](#gfunc)：Processing function `func`
    - [gconv](#gconv)：Data type conversion
    - [gson](#gson)：Processing `JSON`
- [Generic Data Structures](#-generic-data-structures)
    - [tuple](#tuple)：The implementation of tuples provides the definition of 2 to 10 tuples
    - [set](#set)：The implementation of the collection is based on `map[T]struct{}`
    - [skipset](#skipset)：High-performance concurrent set based on skiplist are ~15 times faster than `sync.Map`
    - [skipmap](#skipmap)：High-performance concurrent map implemented based on skiplist, ~10 times faster than `sync.Map`

## ✨ Generic Functional Programming

### goption

Option type, simplifying the processing of `(T, bool)`

Usage：

```go
import (
    "github.com/bytedance/gg/goption"
)
```

Example:

```go
goption.Of(1, true).Value()
// 1
goption.Nil[int]().IsNil()
// true
goption.Nil[int]().ValueOr(10)
// 10
goption.OK(1).IsOK()
// true
goption.OK(1).ValueOrZero()
// 1
goption.OfPtr((*int)(nil)).Ptr()
// nil
goption.Map(goption.OK(1), strconv.Itoa).Get()
// "1" true
```

### gresult

Result type, simplifying the processing of `(T, error)`

Usage：

```go
import (
    "github.com/bytedance/gg/gresult"
)
```

Example:

```go
gresult.Of(strconv.Atoi("1")).Value()
// 1
gresult.Err[int](io.EOF).IsErr()
// true
gresult.Err[int](io.EOF).ValueOr(10)
// 10
gresult.OK(1).IsOK()
// true
gresult.OK(1).ValueOrZero()
// 1
gresult.Of(strconv.Atoi("x")).Option().Get()
// 0 false
gresult.Map(gresult.OK(1), strconv.Itoa).Get()
// "1" nil
```

## ✨ Generic Data Processing

### gcond：

Conditional operation

Usage：

```go
import (
    "github.com/bytedance/gg/gcond"
)
```

Example：

```go
gcond.If(true, 1, 2)
// 1
var a *struct{ A int }
getA := func() int { return a.A }
get1 := func() int { return 1 }
gcond.IfLazy(a != nil, getA, get1)
// 1
gcond.IfLazyL(a != nil, getA, 1)
// 1
gcond.IfLazyR(a == nil, 1, getA)
// 1

gcond.Switch[string](3).
    Case(1, "1").
    CaseLazy(2, func() string { return "3" }).
    When(3, 4).Then("3/4").
    When(5, 6).ThenLazy(func() string { return "5/6" }).
    Default("other"))
// 3/4
```

### gvalue

Processing value `T`

Usage：
```go
import (
    "github.com/bytedance/gg/gvalue"
)
```

Example1：Zero Value

```go
a := gvalue.Zero[int]()
// 0
gvalue.IsZero(a)
// true
b := gvalue.Zero[*int]()
// nil
gvalue.IsNil(b)
// true
```

Example2：Math Operation

```go
gvalue.Max(1, 2, 3)
// 3
gvalue.Min(1, 2, 3)
// 1
gvalue.MinMax(1, 2, 3)
// 1 3
gvalue.Clamp(5, 1, 10)
// 5
gvalue.Add(1, 2)
// 3
```

Example3：Comparison

```go
gvalue.Equal(1, 1)
// true
gvalue.Between(2, 1, 3)
// true
```

Example4：Type Assertion

```go
gvalue.TypeAssert[int](any(1))
// 1
gvalue.TryAssert[int](any(1))
// 1 true
```

Example5：Once

```go
once := gvalue.Once(func() {
    fmt.Println("once")
})
once()
// "once"
once()
// (no output)
```

### gptr

Processing pointer `*T`

Usage：

```go
import (
    "github.com/bytedance/gg/gptr"
)
```

Example：

```go
a := Of(1)
gptr.Indirect(a)
// 1

b := OfNotZero(1)
gptr.IsNotNil(b)
// true
gptr.IndirectOr(b, 2)
// 1
gptr.Indirect(gptr.Map(b, strconv.Itoa))
// "1"

c := OfNotZero(0)
// nil
gptr.IsNil(c)
// true
gptr.IndirectOr(c, 2)
// 2
```

### gslice

Process slice `[]T`

Usage：

```go
import (
    "github.com/bytedance/gg/gslice"
)
```

Example1：High-order Function

```go
gslice.Map([]int{1, 2, 3, 4, 5}, strconv.Itoa)
// ["1", "2", "3", "4", "5"]
isEven := func(i int) bool { return i%2 == 0 }
gslice.Filter([]int{1, 2, 3, 4, 5}, isEven)
// [2, 4]
gslice.Reduce([]int{1, 2, 3, 4, 5}, gvalue.Add[int].Value())
// 15
gslice.Any([]int{1, 2, 3, 4, 5}, isEven)
// true
gslice.All([]int{1, 2, 3, 4, 5}, isEven)
// false
```

Example2：CURD Operation

```go
gslice.Contains([]int{1, 2, 3, 4, 5}, 2)
// true
gslice.ContainsAny([]int{1, 2, 3, 4, 5}, 2, 6)
// true
gslice.ContainsAll([]int{1, 2, 3, 4, 5}, 2, 6)
// false
gslice.Index([]int{1, 2, 3, 4, 5}, 3.Value())
// 2
gslice.Find([]int{1, 2, 3, 4, 5}, isEven).Value()
// 2
gslice.First([]int{1, 2, 3, 4, 5}).Value()
// 1
gslice.Get([]int{1, 2, 3, 4, 5}, 1).Value()
// 2
gslice.Get([]int{1, 2, 3, 4, 5}, -1).Value() // Access element with negative index
// 5
```

Example3：Partion Operation

```go
gslice.Take([]int{1, 2, 3, 4, 5}, 2)
// [1, 2]
gslice.Slice([]int{1, 2, 3, 4, 5}, 1, 3)
// [2, 3]
gslice.Chunk([]int{1, 2, 3, 4, 5}, 2)
// [[1, 2], [3, 4], [5]]
gslice.Divide([]int{1, 2, 3, 4, 5}, 2)
// [[1, 2, 3], [4, 5]]
gslice.Concat([]int{1, 2}, []int{3, 4, 5})
// [1, 2, 3, 4, 5]
gslice.Flatten([][]int{{1, 2}, {3, 4, 5}})
// [1, 2, 3, 4, 5]
gslice.Partition([]int{1, 2, 3, 4, 5}, isEven)
// [2, 4], [1, 3, 5]
```

Example4：Math Operation

```go
gslice.Max([]int{1, 2, 3, 4, 5}).Value()
// 5
gslice.Min([]int{1, 2, 3, 4, 5}).Value()
// 1
gslice.MinMax([]int{1, 2, 3, 4, 5}).Value().Values()
// 1 5
gslice.Sum([]int{1, 2, 3, 4, 5})
// 15
```

Example5：Convert to map

```go
ToMap([]int{1, 2, 3, 4, 5}, func(i int) (string, int) { return strconv.Itoa(i), i })
// {"1":1, "2":2, "3":3, "4":4, "5":5}
ToMapValues([]int{1, 2, 3, 4, 5}, strconv.Itoa)
// {"1":1, "2":2, "3":3, "4":4, "5":5}
GroupBy([]int{1, 2, 3, 4, 5}, func(i int) string {
  if i%2 == 0 {
    return "even"
  } else {
    return "odd"
  }
})
// {"even":[2,4], "odd":[1,3,5]}
```

Example6：Set Operation

```go
gslice.Union([]int{1, 2, 3}, []int{3, 4, 5})
// [1, 2, 3, 4, 5]
gslice.Intersect([]int{1, 2, 3}, []int{3, 4, 5})
// [3]
gslice.Diff([]int{1, 2, 3}, []int{3, 4, 5})
// [1, 2]
gslice.Uniq([]int{1, 1, 2, 2, 3})
// [1, 2, 3]
gslice.Dup([]int{1, 1, 2, 2, 3})
// [1, 2]
```

Example7：Re-order Operation

```go
s1 := []int{5, 1, 2, 3, 4}
s2, s3, s4 := Clone(s1), Clone(s1), Clone(s1)
Sort(s1)
// [1, 2, 3, 4, 5]
SortBy(s2, func(i, j int) bool { return i > j })
// [5, 4, 3, 2, 1]
StableSortBy(s3, func(i, j int) bool { return i > j })
// [5, 4, 3, 2, 1]
Reverse(s4)
// [4, 3, 2, 1, 5]
```

### gmap

Processing map `map[K]V`

Usage：
```go
import (
    "github.com/bytedance/gg/gmap"
)
```

Example1：Keys / Values Getter

```go
gmap.Keys(map[int]int{1: 2})
// [1]
gmap.Values(map[int]int{1: 2})
// [2]
gmap.Items(map[int]int{1: 2}).Unzip()
// [1] [2]
gmap.OrderedKeys(map[int]int{1: 2, 2: 3, 3: 4})
// [1, 2, 3]
gmap.OrderedValues(map[int]int{1: 2, 2: 3, 3: 4})
// [2, 3, 4]
gmap.OrderedItems(map[int]int{1: 2, 2: 3, 3: 4}).Unzip()
// [1, 2, 3] [2, 3, 4]
f := func(k, v int) string { return strconv.Itoa(k) + ":" + strconv.Itoa(v) }
gmap.ToSlice(map[int]int{1: 2}, f)
// ["1:2"]
gmap.ToOrderedSlice(map[int]int{1: 2, 2: 3, 3: 4}, f)
// ["1:2", "2:3", "3:4"]
```

Example2：High-order Function

```go
Map(map[int]int{1: 2, 2: 3, 3: 4}, func(k int, v int) (string, string) {
    return strconv.Itoa(k), strconv.Itoa(k + 1)
})
// {"1":"2", "2":"3", "3":"4"}
Filter(map[int]int{1: 2, 2: 3, 3: 4}, func(k int, v int) bool {
    return k+v > 3
})
// {"2":2, "3":3}
```

Example3：CURD Operation

```go
gmap.Contains(map[int]int{1: 2, 2: 3, 3: 4}, 1)
// true
gmap.ContainsAny(map[int]int{1: 2, 2: 3, 3: 4}, 1, 4)
// true
gmap.ContainsAll(map[int]int{1: 2, 2: 3, 3: 4}, 1, 4)
// false
gmap.Load(map[int]int{1: 2, 2: 3, 3: 4}, 1).Value()
// 2
gmap.LoadAny(map[int]int{1: 2, 2: 3, 3: 4}, 1, 4).Value()
// 2
gmap.LoadAll(map[int]int{1: 2, 2: 3, 3: 4}, 1, 4)
// []
gmap.LoadSome(map[int]int{1: 2, 2: 3, 3: 4}, 1, 4)
// [2]
```

Example4：Partion Operation

```go
Chunk(map[int]int{1: 2, 2: 3, 3: 4, 4: 5, 5: 6}, 2)
// possible output: [{1:2, 2:3}, {3:4, 4:5}, {5:6}]
Divide(map[int]int{1: 2, 2: 3, 3: 4, 4: 5, 5: 6}, 2)
// possible output: [{1:2, 2:3, 3:4}, {4:5, 5:6}]
```

Example5：Math Operation

```go
gmap.Max(map[int]int{1: 2, 2: 3, 3: 4}).Value()
// 4
gmap.Min(map[int]int{1: 2, 2: 3, 3: 4}).Value()
// 2
gmap.MinMax(map[int]int{1: 2, 2: 3, 3: 4}).Value().Values()
// 2 4
gmap.Sum(map[int]int{1: 2, 2: 3, 3: 4})
// 9
```

Example6：Set Operation

```go
gmap.Union(map[int]int{1: 2, 2: 3, 3: 4}, map[int]int{3: 14, 4: 15, 5: 16})
// {1:2, 2:3, 3:14, 4:15, 5:16}
gmap.Intersect(map[int]int{1: 2, 2: 3, 3: 4}, map[int]int{3: 14, 4: 15, 5: 16})
// {3:14}
gmap.Diff(map[int]int{1: 2, 2: 3, 3: 4}, map[int]int{3: 14, 4: 15, 5: 16})
// {1:2, 2:3}
gmap.UnionBy(gslice.Of(map[int]int{1: 2, 2: 3, 3: 4}, map[int]int{3: 14, 4: 15, 5: 16}), DiscardNew[int, int]())
// {1:2, 2:3, 3:4, 4:15, 5:16}
gmap.IntersectBy(gslice.Of(map[int]int{1: 2, 2: 3, 3: 4}, map[int]int{3: 14, 4: 15, 5: 16}), DiscardNew[int, int]())
// {3:4}
```

### gfunc

Processing function `func`

Usage：

```go
import (
    "github.com/bytedance/gg/gfunc"
)
```

Example1：Partial Application

```go
add := Partial2(gvalue.Add[int]) // convert the Add function into a partial function
add1 := add.Partial(1)           // Bind (i.e., "freeze") the first argument to 1
add1(0)                          // 0 + 1 = 1
// 1
add1(1)                          // Reuse the partially applied function: 1 + 1 = 2
// 2
add1n2 := add1.PartialR(2)       // Bind the remaining (rightmost) argument to 2; all arguments are now fixed
add1n2()                         // 1 + 2 = 3
```

### gconv

Data type conversion

Usage：

```go
import (
    "github.com/bytedance/gg/gconv"
)
```

Example：

```go
gconv.To[string](1)
// "1"
gconv.To[int]("1")
// 1
gconv.To[int]("x")
// 0
gconv.To[bool]("true")
// true
gconv.To[bool]("x")
// false
gconv.To[int](gptr.Of(gptr.Of(gptr.Of("1"))))
// 1
type myInt int
type myString string
gconv.To[myInt](myString("1"))
// 1
gconv.To[myString](myInt(1))
// "1"

gconv.ToE[int]("x")
// 0 strconv.ParseInt: parsing "x": invalid syntax
```

### gson

Processing `JSON`

Usage：

```go
import (
    "github.com/bytedance/gg/gson"
)
```

Example：

```go
type testStruct struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}
testcase := testStruct{Name: "test", Age: 10}

gson.Marshal(testcase)
// []byte(`{"name":"test","age":10}`) nil
gson.MarshalString(testcase)
// `{"name":"test","age":10}` nil
gson.ToString(testcase)
// `{"name":"test","age":10}`
gson.MarshalIndent(testcase, "", "  ")
// "{\n  \"name\": \"test\",\n  \"age\": 10\n}" nil
gson.ToStringIndent(testcase, "", "  ")
// "{\n  \"name\": \"test\",\n  \"age\": 10\n}"
gson.Valid(`{"name":"test","age":10}`)
// true
gson.Unmarshal[testStruct](`{"name":"test","age":10}`)
// {test 10} nil
```

## ✨ Generic Data Structures

### tuple

The implementation of tuples provides the definition of 2 to 10 tuples

Usage

```go
import (
    "github.com/bytedance/gg/collection/tuple"
)
```

Example：

```go
addr := Make2("localhost", 8080)
fmt.Printf("%s:%d\n", addr.First, addr.Second)
// localhost:8080

s := Zip2([]string{"red", "green", "blue"}, []int{14, 15, 16})
for _, v := range s {
    fmt.Printf("%s:%d\n", v.First, v.Second)
}
// red:14
// green:15
// blue:16

fmt.Println(s.Unzip())
// ["red", "green", "blue"] [14, 15, 16]
```

### set

Set implementation based on `map[T]struct{}`

Usage

```go
import (
    "github.com/bytedance/gg/collection/set"
)
```

Example：

```go
s := New(10, 10, 12, 15)
s.Len()
// 3
s.Add(10)
// false
s.Add(11)
// true
s.Remove(11) && s.Remove(12)
// true

s.ContainsAny(10, 15)
// true
s.ContainsAny(11, 12)
// false
s.ContainsAny()
// false
s.ContainsAll(10, 15)
// true
s.ContainsAll(10, 11)
// false
s.ContainsAll()
// true

len(s.ToSlice())
// 2
```

### skipset

High-performance concurrent sets based on skiplist are ~15 times faster than sync.Map

Usage

```go
import (
    "github.com/bytedance/gg/collection/skipset"
)
```

Example：

```go
s := skipset.New[int]()
s.Add(10)
// true
s.Add(10)
// false
s.Add(11)
// true
s.Add(12)
// true
s.Len()
// 3

s.Contains(10)
// true
s.Remove(10)
// true
s.Contains(10)
// false

s.ToSlice()
// [11, 12]

var wg sync.WaitGroup
wg.Add(1000)
for i := 0; i < 1000; i++ {
    i := i
    go func() {
        defer wg.Done()
        s.Add(i)
    }()
}
wg.Wait()
s.Len()
// 1000
```

### skipmap

High-performance concurrent hash list implemented based on skiplist, ~10 times faster than `sync.Map` below Go 1.23.

After Go 1.24, please consider using the std `sync.Map`, which has better performance compared to skipmap in about 90% of use cases.

Usage

```go
import (
    "github.com/bytedance/gg/collection/skipmap"
)
```

Example：

```go
s := New[string, int]()
s.Store("a", 0)
s.Store("a", 1)
s.Store("b", 2)
s.Store("c", 3)
s.Len()
// 3

s.Load("a")
// 1 true
s.LoadAndDelete("a")
// 1 true
s.LoadOrStore("a", 11)
// 11 false

gson.ToString(s.ToMap())
// {"a":11, "b":2, "c": 3}

s.Delete("a")
s.Delete("b")
s.Delete("c")
var wg sync.WaitGroup
wg.Add(1000)
for i := 0; i < 1000; i++ {
    i := i
    go func() {
        defer wg.Done()
        s.Store(strconv.Itoa(i), i)
    }()
}
wg.Wait()
s.Len()
// 1000
```

## License

`gg` is licensed under the Apache-2.0 license. See [LICENSE](LICENSE) for details.

2025 © Bytedance
