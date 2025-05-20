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

//go:build go1.20
// +build go1.20

package gmap

import (
	"errors"
	"strconv"
	"testing"

	"github.com/bytedance/gg/gresult"
	"github.com/bytedance/gg/internal/assert"
)

// errors.Join is introduced in 1.20+
func TestTryMap(t *testing.T) {
	f := func(k, v string) (int, int, error) {
		ki, kerr := strconv.Atoi(k)
		vi, verr := strconv.Atoi(v)
		return ki, vi, errors.Join(kerr, verr)
	}
	assert.Equal(t,
		gresult.OK(map[int]int{}),
		TryMap(map[string]string{}, f))
	assert.Equal(t,
		gresult.OK(map[int]int{1: 1, 2: 2}),
		TryMap(map[string]string{"1": "1", "2": "2"}, f))
	assert.Equal(t,
		"strconv.Atoi: parsing \"a\": invalid syntax",
		TryMap(map[string]string{"1": "1", "2": "a"}, f).Err().Error())
}

func TestFilterMap(t *testing.T) {
	parseInt := func(k, v string) (int, int, bool) {
		ki, kerr := strconv.ParseInt(k, 10, 64)
		vi, verr := strconv.ParseInt(v, 10, 64)
		return int(ki), int(vi), errors.Join(kerr, verr) == nil
	}
	assert.Equal(t,
		map[int]int{1: 1, 2: 2},
		FilterMap(map[string]string{"1": "1", "2": "2", "a": "3", "4": "b", "c": "c"}, parseInt))
	assert.Equal(t,
		map[int]int{},
		FilterMap(map[string]string{"a": "3", "4": "b"}, parseInt))
	assert.Equal(t,
		map[int]int{1: 1, 2: 2},
		FilterMap(map[string]string{"1": "1", "2": "2"}, parseInt))
	assert.Equal(t,
		map[int]int{},
		FilterMap(map[string]string{}, parseInt))
	assert.Equal(t,
		map[int]int{},
		FilterMap(nil, parseInt))
}

func TestTryFilterMap(t *testing.T) {
	parseInt := func(k, v string) (int, int, error) {
		ki, kerr := strconv.ParseInt(k, 10, 64)
		vi, verr := strconv.ParseInt(v, 10, 64)
		return int(ki), int(vi), errors.Join(kerr, verr)
	}
	assert.Equal(t,
		map[int]int{1: 1, 2: 2},
		TryFilterMap(map[string]string{"1": "1", "2": "2", "a": "3", "4": "b", "c": "c"}, parseInt))
	assert.Equal(t,
		map[int]int{},
		TryFilterMap(map[string]string{"a": "3", "4": "b"}, parseInt))
	assert.Equal(t,
		map[int]int{1: 1, 2: 2},
		TryFilterMap(map[string]string{"1": "1", "2": "2"}, parseInt))
	assert.Equal(t,
		map[int]int{},
		TryFilterMap(map[string]string{}, parseInt))
	assert.Equal(t,
		map[int]int{},
		TryFilterMap(nil, parseInt))
}
