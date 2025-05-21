#!/bin/sh

# Generate Comparable[comparable] from Stream[any]
 go run gen.go \
    -parent "Stream" \
    -parent-types "T any" \
    -child "Comparable" \
    -child-types "T comparable" \
    -ignore-paths "github.com/bytedance/gg/goption github.com/bytedance/gg/collection/tuple"

# Generate Bool[~bool] from from Comparable[comparable]
go run gen.go \
    -parent "Comparable" \
    -parent-types "T comparable" \
    -child "Bool" \
    -child-types "T ~bool" \
    -ignore-funcs "Uniq UniqBy Remove RemoveN" # Distinct funcs are meaningless for bool type.

# Generate Orderable[constraints.Ordered] from Comparable[comparable]
go run gen.go \
    -parent "Comparable" \
    -parent-types "T comparable" \
    -child "Orderable" \
    -child-types "T constraints.Ordered" \
    -import-paths "github.com/bytedance/gg/internal/constraints"

# Generate String[~string] from Orderable[constraints.Ordered]
go run gen.go \
    -parent "Orderable" \
    -parent-types "T constraints.Ordered" \
    -child "String" \
    -child-types "T ~string" \
    -ignore-paths "github.com/bytedance/gg/internal/constraints github.com/bytedance/gg/goption github.com/bytedance/gg/collection/tuple"

# Generate KV[comparable, any] from Stream[any]
go run gen.go \
    -parent "Stream" \
    -parent-types "T any" \
    -child "KV" \
    -child-types "K comparable V any" \
    -child-elem "tuple.T2[K, V]" \
    -import-paths "github.com/bytedance/gg/collection/tuple" \
    -ignore-paths "github.com/bytedance/gg/goption" \
    -ignore-funcs "FromMapValues Repeat Map Fold Reduce Filter ForEach All Any Zip Intersperse Append Prepend Find TakeWhile DropWhile SortBy UniqBy DistinctOrderedWith"

# Generate OrderableKV[constraints.Ordered, any] from KV[comparable, any]
go run gen.go \
    -parent "KV" \
    -parent-types "K comparable V any" \
    -parent-elem "tuple.T2[K, V]" \
    -child "OrderableKV" \
    -child-types "K constraints.Ordered V any" \
    -child-elem "tuple.T2[K, V]" \
    -import-paths "github.com/bytedance/gg/internal/constraints" \
    -ignore-paths "github.com/bytedance/gg/goption" \
    -ignore-funcs "FromMap"

# Generate Number[constraints.Number] from Orderable[constraints.Ordered]
go run gen.go \
    -parent "Orderable" \
    -parent-types "T constraints.Ordered" \
    -child "Number" \
    -child-types "T constraints.Number" \
    -ignore-paths "github.com/bytedance/gg/goption github.com/bytedance/gg/collection/tuple"
