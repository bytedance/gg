# gg：Go 泛型基础库

[![GoDoc](https://godoc.org/github.com/bytedance/gg?status.svg)](https://godoc.org/github.com/bytedance/gg)
[![Go Report Card](https://goreportcard.com/badge/github.com/bytedance/gg)](https://goreportcard.com/report/github.com/bytedance/gg)
[![Go Coverage](https://codecov.io/gh/bytedance/gg/branch/main/graph/badge.svg)](https://codecov.io/gh/bytedance/gg)
[![License](https://img.shields.io/github/license/bytedance/gg)](https://github.com/bytedance/gg/blob/main/LICENSE)

[English](README.md) | 简体中文

🔥`bytedance/gg` 是字节跳动开发的 Go 语言泛型基础库，基于 Go 1.18+ 泛型特性，提供高效、类型安全且丰富的泛型数据结构与工具函数。

❓泛型库为什么叫 gg？

取 **G**o **G**enerics（Go 泛型）的首字母，简短顺口。

❓为什么选择 gg？

- 稳定可靠：字节跳动研发团队必备的工具依赖库，内部有 1w+ 代码仓库引用。
- 使用方便：以简单自洽为设计原则，根据功能分包，模块化，语义直观且统一，学习成本低。
- 高性能：提供了高性能并发数据结构，性能比标准库快 10+ 倍。
- 无三方依赖：泛型库不会引入任何三方依赖。
- 版本控制：遵循语义化版本号SemVer，保证向后兼容。

## 🚀 安装

```sh
go get github.com/bytedance/gg
```

## 🔎 目录

- [泛型函数式编程](#-泛型函数式编程)
  - [goption](#goption)：选项类型，简化 `(T, bool)` 返回值的处理
  - [gresult](#gresult)：结果类型，简化 `(T, error)` 返回值的处理
- [泛型数据处理](#-泛型数据处理)
  - [gcond](#gcond)：条件运算
  - [gvalue](#gvalue)：处理值 `T`
  - [gptr](#gptr)：处理指针 `*T`
  - [gslice](#gslice)：处理切片 `[]T`
  - [gmap](#gmap)：处理散列表 `map[K]V`
  - [gfunc](#gfunc)：处理函数 `func`
  - [gconv](#gconv)：数据类型转换
  - [gson](#gson)：处理 JSON 数据
- [泛型标准库封装](#-泛型标准库封装)
  - [gsync](#gsync)：封装 `sync` 标准库
- [泛型数据结构](#-泛型数据结构)
  - [tuple](#tuple)：元组的实现，提供了 2～10 元组的定义
  - [set](#set)：集合的实现，基于 `map[T]struct{}`
  - [skipset](#skipset)：基于 skiplist 实现的高性能并发集合，在 Go 1.24 以下版本比标准库 `sync.Map` 快 ~15 倍
  - [skipmap](#skipmap)：基于 skiplist 实现的高性能并发散列表，在 Go 1.24 以下版本比标准库 `sync.Map` 快 ~10 倍

## ✨ 泛型函数式编程

### goption

选项类型，简化 `(T, bool)` 返回值的处理

引用：

```go
import (
    "github.com/bytedance/gg/goption"
)
```

示例:

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

提供了结果类型，用来简化 `(T, error)` 返回值的处理

引用：

```go
import (
    "github.com/bytedance/gg/gresult"
)
```

示例:

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

## ✨ 泛型数据处理

### gcond：

条件运算

引用：

```go
import (
    "github.com/bytedance/gg/gcond"
)
```

示例：

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

处理值 `T`

引用：
```go
import (
    "github.com/bytedance/gg/gvalue"
)
```

示例1：零值处理

```go
a := gvalue.Zero[int]()
// 0
gvalue.IsZero(a)
// true
b := gvalue.Zero[*int]()
// nil
gvalue.IsNil(b)
// true
gvalue.Or(0, 1, 2)
// 1
```

示例2：数学运算

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

示例3：数值比较

```go
gvalue.Equal(1, 1)
// true
gvalue.Between(2, 1, 3)
// true
```

示例4：类型断言

```go
gvalue.TypeAssert[int](any(1))
// 1
gvalue.TryAssert[int](any(1))
// 1 true
```

### gptr

处理指针 `*T`

引用：

```go
import (
    "github.com/bytedance/gg/gptr"
)
```

示例：

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

处理切片 `[]T`

引用：

```go
import (
    "github.com/bytedance/gg/gslice"
)
```

示例1：高阶函数

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

示例2：增删改查

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
gslice.Get([]int{1, 2, 3, 4, 5}, -1).Value() // 负索引
// 5
```

示例3：分块操作

```go
gslice.Range(1, 5)
// [1, 2, 3, 4]
gslice.RangeWithStep(5, 1, -2)
// [5, 3]
gslice.Take([]int{1, 2, 3, 4, 5}, 2)
// [1, 2]
gslice.Take([]int{1, 2, 3, 4, 5}, -2)
// [4, 5]
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

示例4：数学运算

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

示例5：转换为map

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

示例6：集合操作

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

示例7：排序操作

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

处理散列表 `map[K]V`

引用：
```go
import (
    "github.com/bytedance/gg/gmap"
)
```

示例1：键值获取

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

示例2：高阶函数

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

示例3：增删改查

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

示例4：分块操作

```go
Chunk(map[int]int{1: 2, 2: 3, 3: 4, 4: 5, 5: 6}, 2)
// possible output: [{1:2, 2:3}, {3:4, 4:5}, {5:6}]
Divide(map[int]int{1: 2, 2: 3, 3: 4, 4: 5, 5: 6}, 2)
// possible output: [{1:2, 2:3, 3:4}, {4:5, 5:6}]
```

示例5：数学运算

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

示例6：集合操作

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

处理函数 `func`

引用：

```go
import (
    "github.com/bytedance/gg/gfunc"
)
```

示例1：偏函数

```go
add := Partial2(gvalue.Add[int]) // 将 Add 转化为偏函数
add1 := add.Partial(1)           // 绑定第一个参数为 1
add1(0)                          // 0 + 1 = 1
// 1
add1(1)                          // add1 可以重复使用, 1 + 1 = 2
// 2
add1n2 := add1.PartialR(2)       // 绑定最后一个参数为 2，所有参数都为固定值
add1n2()                         // 1 + 2 = 3
// 3
```

### gconv

数据类型转换

引用：

```go
import (
    "github.com/bytedance/gg/gconv"
)
```

示例：

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

处理 JSON 数据

引用：

```go
import (
    "github.com/bytedance/gg/gson"
)
```

示例：

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

## ✨ 泛型标准库封装

### gsync

封装 `sync` 标准库

引用：

```go
import (
    "github.com/bytedance/gg/gstd/gsync"
)
```

示例1：`gsync.Map` 封装了 `sync.Map`

```go
sm := gsync.Map[string, int]{}
sm.Store("k", 1)
sm.Load("k")
// 1 true
sm.LoadO("k").Value()
// 1
sm.Store("k", 2)
sm.Load("k")
// 2 true
sm.LoadAndDelete("k")
// 2 true
sm.Load("k")
// 0 false
sm.LoadOrStore("k", 3)
// 3 false
sm.Load("k")
// 3 true
sm.ToMap()
// {"k":3}
```

示例2：`gsync.Pool` 封装了 `sync.Pool`

```go
pool := Pool[*int]{
    New: func() *int {
        i := 1
        return &i
    },
}
a := pool.Get()
*a
// 1
*a = 2
pool.Put(a)
*pool.Get()
// 可能的结果: 1 或 2
```

示例3：`gsync.OnceXXX` 封装了 `sync.Once`

```go
onceFunc := gsync.OnceFunc(func() { fmt.Println("OnceFunc") })
onceFunc()
// "OnceFunc"
onceFunc()
// (no output)
onceFunc()
// (no output)

i := 1
onceValue := gsync.OnceValue(func() int { i++; return i })
onceValue()
// 2
onceValue()
// 2

onceValues := gsync.OnceValues(func() (int, error) { i++; return i, nil })
onceValues()
// 3 nil
onceValues()
// 3 nil
```

## ✨ 泛型数据结构

### tuple

元组的实现，提供了 2～10 元组的定义

引用

```go
import (
    "github.com/bytedance/gg/collection/tuple"
)
```

示例：

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

s.Unzip()
// ["red", "green", "blue"] [14, 15, 16]
```

### set

集合的实现，基于 map[T]struct{}

引用

```go
import (
    "github.com/bytedance/gg/collection/set"
)
```

示例：

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

基于 skiplist 实现的高性能并发集合，在 Go 1.24 以下版本比标准库 `sync.Map` 快 ~15 倍

⚠️ 注意：Go 1.24 及更高版本，建议使用标准库 `sync.Map`，在约 90% 的使用场景中，其性能优于 `skipset`

引用

```go
import (
    "github.com/bytedance/gg/collection/skipset"
)
```

示例：

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

基于 skiplist 实现的高性能并发散列表，在 Go 1.24 以下版本比标准库 `sync.Map` 快 ~15 倍

⚠️ 注意：Go 1.24 及更高版本，建议使用标准库 `sync.Map`，在约 90% 的使用场景中，其性能优于 `skipmap`

引用

```go
import (
    "github.com/bytedance/gg/collection/skipmap"
)
```

示例：

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

## 许可证

`gg` 采用 Apache-2.0 许可证。详情请参阅 [LICENSE](LICENSE)。

2025 © 字节跳动
