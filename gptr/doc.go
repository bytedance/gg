// Package gptr provides generic operations for pointers.
//
// 💡 HINT: We provide similar functionality for different types in different packages.
// For example, [github.com/bytedance/gg/gslice.Clone] for copying slice while
// [github.com/bytedance/gg/gmap.Clone] for copying map.
//
//   - Use [github.com/bytedance/gg/gslice] for slice operations.
//   - Use [github.com/bytedance/gg/gmap] for map operations.
//   - Use [github.com/bytedance/gg/gvalue] for value operations.
//   - …
//
// # Operations
//
//   - Reference (T → *T): [Of], [OfNotZero], …
//   - Dereference (*T → T): [Indirect], [IndirectOr], …
//   - Predicate: (*T → bool): [Equal], [EqualTo], [IsNil], …
package gptr
