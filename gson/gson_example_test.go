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

package gson

import (
	"fmt"

	"github.com/bytedance/gg/gresult"
)

func Example() {
	type testStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	testcase := testStruct{Name: "test", Age: 10}

	fmt.Println(string(gresult.Of(Marshal(testcase)).Value()))                 // `{"name":"test","age":10}`
	fmt.Println(gresult.Of(MarshalString(testcase)).Value())                   // `{"name":"test","age":10}`
	fmt.Println(ToString(testcase))                                            // `{"name":"test","age":10}`
	fmt.Println(string(gresult.Of(MarshalIndent(testcase, "", "  ")).Value())) // "{\n  \"name\": \"test\",\n  \"age\": 10\n}"
	fmt.Println(ToStringIndent(testcase, "", "  "))                            // "{\n  \"name\": \"test\",\n  \"age\": 10\n}"
	fmt.Println(Valid(`{"name":"test","age":10}`))                             // true
	fmt.Println(Unmarshal[testStruct](`{"name":"test","age":10}`))             // {test 10} nil

	// Output:
	// {"name":"test","age":10}
	// {"name":"test","age":10}
	// {"name":"test","age":10}
	// {
	//   "name": "test",
	//   "age": 10
	// }
	// {
	//   "name": "test",
	//   "age": 10
	// }
	// true
	// {test 10} <nil>
}
