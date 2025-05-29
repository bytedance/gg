// Copyright 2025 Bytedance Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gsync

import "sync"

// OnceFunc returns a function that invokes f only once when returned
// function is firstly called.
func OnceFunc(f func()) func() {
	var once sync.Once
	return func() {
		once.Do(func() { f() })
	}
}

// OnceValue returns a function as value getter.
// Value is returned by function f, and f is invoked only once when returned
// function is firstly called.
//
// This function can be used to lazily initialize a value, as replacement of
// the packages-level init function. For example:
//
//	var DB *sql.DB
//
//	func init() {
//		// 💡 NOTE: DB is initialized here.
//		DB, _ = sql.Open("mysql", "user:password@/dbname")
//	}
//
//	func main() {
//		DB.Query(...)
//	}
//
// Can be rewritten to:
//
//	var DB = OnceValue(func () *sql.DB {
//		return gresult.Of(sql.Open("mysql", "user:password@/dbname")).Value()
//	})
//
//	func main() {
//		// 💡 NOTE: DB is *LAZILY* initialized here.
//		DB().Query(...)
//	}
//
// 💡 HINT:
//
//   - See also https://github.com/golang/go/issues/56102
func OnceValue[T any](f func() T) func() T {
	var (
		once sync.Once
		v    T
	)
	return func() T {
		once.Do(func() { v = f() })
		return v
	}
}

// OnceValues returns a function as values getter.
// Values are returned by function f, and f is invoked only once when returned
// function is firstly called.
func OnceValues[T1, T2 any](f func() (T1, T2)) func() (T1, T2) {
	var (
		once sync.Once
		v1   T1
		v2   T2
	)
	return func() (T1, T2) {
		once.Do(func() { v1, v2 = f() })
		return v1, v2
	}
}
