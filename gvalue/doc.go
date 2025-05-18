// Package gvalue provides generic operations for go values.
//
// 💡 HINT: We provide similar functionality for different types in different packages.
// For example, [github.com/bytedance/gg/gslice.Clone] for copying slice while
// [github.com/bytedance/gg/gmap.Clone] for copying map.
//
//   - Use [github.com/bytedance/gg/gslice] for slice operations.
//   - Use [github.com/bytedance/gg/gmap] for map operations.
//   - Use [github.com/bytedance/gg/gptr] for pointer operations.
//   - …
//
// # Operations
//
//   - Math operations: [Max], [Min], [MinMax], [Clamp], …
//   - Type assertion (T1 → T2): [TypeAssert], [TryAssert], …
//   - Predicate: (T → bool): [Equal], [Greater], [Less], [Between], [IsNil], [IsZero], …
package gvalue
