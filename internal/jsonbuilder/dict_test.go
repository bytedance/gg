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

package jsonbuilder

import (
	"encoding/json"
	"testing"

	"github.com/bytedance/gg/internal/assert"
)

func TestDict(t *testing.T) {
	{
		s := map[int]string{1: "1", 2: "2", 3: "3", 4: "4", 5: "5"}
		d := NewDict()
		for k, v := range s {
			assert.Nil(t, d.Store(k, v))
		}
		d.Sort()
		bs, err := d.Build()
		assert.Nil(t, err)
		assert.Equal(t, []byte(`{"1":"1","2":"2","3":"3","4":"4","5":"5"}`), bs)
	}

	{
		s := map[int]string{}
		d := NewDict()
		for k, v := range s {
			assert.Nil(t, d.Store(k, v))
		}
		bs, err := d.Build()
		assert.Nil(t, err)
		assert.Equal(t, []byte(`{}`), bs)
	}

	{
		s := map[int]string{1: "1"}
		d := NewDict()
		for k, v := range s {
			assert.Nil(t, d.Store(k, v))
		}
		bs, err := d.Build()
		assert.Nil(t, err)
		assert.Equal(t, []byte(`{"1":"1"}`), bs)
	}

	{
		d := NewDict()
		assert.Nil(t, d.Store("a", "b"))
		assert.NotNil(t, d.Store(1.4, "b")) // float key is not supported
		assert.Nil(t, d.Store(1, 1.2))
		assert.Nil(t, d.Store("w", []int{1}))

		bs, err := d.Build()
		assert.Nil(t, err)
		t.Log(string(bs))
		assert.Equal(t, []byte(`{"a":"b","1":1.2,"w":[1]}`), bs)
	}

	{
		var d *Dict
		bs, err := d.Build()
		assert.Nil(t, err)
		assert.Equal(t, []byte(`null`), bs)
	}
}

func TestDict2(t *testing.T) {

	testCases := []map[string]interface{}{
		{
			"Name":       "John Doe",
			"Age":        float64(30),
			"Profession": "Software Engineer",
		},
		{
			"Name":       "John Doe",
			"Age":        float64(1 / 3),
			"Profession": "Software Engineer",
		},
		{
			"customerID": "1234",
			"amount":     99.99,
			"items": []interface{}{
				map[string]interface{}{"id": float64(1), "name": "Apple", "price": 1.99},
				map[string]interface{}{"id": float64(2), "name": "Banana", "price": 2.99},
			},
		},
		{
			"author": "Jane Austen",
			"books": []interface{}{
				"Pride and Prejudice",
				"Sense and Sensibility",
				"Emma",
			},
		},
		{
			"pi":       3.14159265359,
			"e":        2.71828182846,
			"sqrt_two": 1.41421356237,
			"phi":      1.61803398875,
		},
	}

	for _, c := range testCases {
		dict := NewDict()
		for k, v := range c {
			dict.Store(k, v)
		}
		out, err := dict.Build()
		assert.Nil(t, err)
		var m map[string]interface{}
		assert.Nil(t, json.Unmarshal(out, &m))
		t.Log(c, len(c))
		t.Log(m, len(m))
		assert.Equal(t, c, m)
	}

}
