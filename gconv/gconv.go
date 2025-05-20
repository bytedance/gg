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

// Package gconv helps users convert any type to the convertible type.
package gconv

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/bytedance/gg/gcond"
	"github.com/bytedance/gg/gresult"
	"github.com/bytedance/gg/gvalue"
	"github.com/bytedance/gg/internal/constraints"
)

var errUnsupported = errors.New("unsupported type conversion")

type convertible interface {
	~bool | constraints.Number | ~string
}

// To converts any type to the convertible type.
// If the conversion is not supported, a zero value is returned.
//
// 🚀 EXAMPLE:
//
//	To[bool]("true") ⏩ true
//	To[int64]("1.0") ⏩ 1
//	To[string](true) ⏩ "true"
//	p := gptr.Of(gptr.Of(gptr.Of('a')))
//	To[int64](p)     ⏩ 97
//	To[int64]("a")   ⏩ 0
//
// ⚠️ WARNING: byte is an alias for uint8, rune is an alias for int32.
//
//	To[string]('a')  ⏩ "97"
//	To[string]('中') ⏩ "20013"
func To[T convertible, V any](v V) T {
	t, _ := ToE[T](v)
	return t
}

// ToPtr converts any type to a pointer of the convertible type.
// If the conversion is not supported, nil is returned.
//
// 🚀 EXAMPLE:
//
//	ToPtr[bool]("true") ⏩ (*bool)(true)
//	ToPtr[int64]("1.0") ⏩ (*int64)(1)
//	ToPtr[string](true) ⏩ (*string)("true")
//	p := gptr.Of(gptr.Of(gptr.Of('a')))
//	ToPtr[int64](p)     ⏩ (*int64)(97)
//	ToPtr[int64]("a")   ⏩ (*int64)(nil)
//
// ⚠️ WARNING: byte is an alias for uint8, rune is an alias for int32.
//
//	ToPtr[string]('a')  ⏩ (*string)("97")
//	ToPtr[string]('中') ⏩ (*string)("20013")
func ToPtr[T convertible, V any](v V) *T {
	t, err := ToE[T](v)
	if err != nil {
		return nil
	}
	return &t
}

// ToR converts any type to gresult.R.
//
// 🚀 EXAMPLE:
//
//	ToR[bool]("true") ⏩ gresult.OK(true)
//	ToR[int64]("1.0") ⏩ gresult.OK[int64](1)
//	ToR[string](true) ⏩ gresult.OK("true")
//	p := gptr.Of(gptr.Of(gptr.Of('a')))
//	ToR[int64](p)     ⏩ gresult.OK[int64](97)
//	ToR[int64]("a")   ⏩ gresult.Err[int64]("strconv.ParseInt: parsing \"a\": invalid syntax")
//
// ⚠️ WARNING: byte is an alias for uint8, rune is an alias for int32.
//
//	ToR[string]('a')  ⏩ gresult.OK("97")
//	ToR[string]('中') ⏩ gresult.OK("20013")
func ToR[T convertible, V any](v V) gresult.R[T] {
	return gresult.Of(ToE[T](v))
}

// ToE converts any type to the convertible type.
// If the conversion is not supported, a zero value and an error is returned.
//
// 🚀 EXAMPLE:
//
//	ToE[bool]("true") ⏩ true, nil
//	ToE[int64]("1.0") ⏩ 1, nil
//	ToE[string](true) ⏩ "true", nil
//	p := gptr.Of(gptr.Of(gptr.Of('a')))
//	ToE[int64](p)     ⏩ 97, nil
//	ToE[int64]("a")   ⏩ 0, "strconv.ParseInt: parsing \"a\": invalid syntax"
//
// ⚠️ WARNING: byte is an alias for uint8, rune is an alias for int32.
//
//	ToE[string]('a')  ⏩ "97", nil
//	ToE[string]('中') ⏩ "20013", nil
func ToE[T convertible, V any](v V) (T, error) {
	t := gvalue.Zero[T]()
	switch any(t).(type) {
	case bool:
		return assertT[T](toBool(v))
	case int:
		return assertT[T](toNumber[int](v))
	case int8:
		return assertT[T](toNumber[int8](v))
	case int16:
		return assertT[T](toNumber[int16](v))
	case int32:
		return assertT[T](toNumber[int32](v))
	case int64:
		return assertT[T](toNumber[int64](v))
	case uint:
		return assertT[T](toNumber[uint](v))
	case uint8:
		return assertT[T](toNumber[uint8](v))
	case uint16:
		return assertT[T](toNumber[uint16](v))
	case uint32:
		return assertT[T](toNumber[uint32](v))
	case uint64:
		return assertT[T](toNumber[uint64](v))
	case uintptr:
		return assertT[T](toNumber[uintptr](v))
	case float32:
		return assertT[T](toNumber[float32](v))
	case float64:
		return assertT[T](toNumber[float64](v))
	case string:
		return assertT[T](toString(v))
	default:
		switch reflect.TypeOf(t).Kind() {
		case reflect.Bool:
			return convertT[T](toBool(v))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return convertT[T](toNumber[int64](v))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return convertT[T](toNumber[uint64](v))
		case reflect.Float32, reflect.Float64:
			return convertT[T](toNumber[float64](v))
		case reflect.String:
			return convertT[T](toString(v))
		default:
			return t, errUnsupported
		}
	}
}

func assertT[T any](a any, err error) (T, error) {
	return a.(T), err
}

func convertT[T any](a any, err error) (T, error) {
	t := gvalue.Zero[T]()
	return reflect.ValueOf(a).Convert(reflect.TypeOf(t)).Interface().(T), err
}

func toBool(a any) (bool, error) {
	a = indirect(a)
	switch v := a.(type) {
	case bool:
		return v, nil
	case nil:
		return false, nil
	case int:
		return v != 0, nil
	case int8:
		return v != 0, nil
	case int16:
		return v != 0, nil
	case int32:
		return v != 0, nil
	case int64:
		return v != 0, nil
	case uint:
		return v != 0, nil
	case uint8:
		return v != 0, nil
	case uint16:
		return v != 0, nil
	case uint32:
		return v != 0, nil
	case uint64:
		return v != 0, nil
	case uintptr:
		return v != 0, nil
	case float32:
		return v != 0, nil
	case float64:
		return v != 0, nil
	case complex64:
		return v != 0, nil
	case complex128:
		return v != 0, nil
	case string:
		return strconv.ParseBool(v)
	case []byte:
		return strconv.ParseBool(string(v))
	default:
		rt := reflect.TypeOf(a)
		switch rt.Kind() {
		case reflect.Bool:
			return reflect.ValueOf(a).Bool(), nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return reflect.ValueOf(a).Int() != 0, nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return reflect.ValueOf(a).Uint() != 0, nil
		case reflect.Float32, reflect.Float64:
			return reflect.ValueOf(a).Float() != 0, nil
		case reflect.Complex64, reflect.Complex128:
			return reflect.ValueOf(a).Complex() != 0, nil
		case reflect.String:
			return strconv.ParseBool(reflect.ValueOf(a).String())
		case reflect.Slice:
			if rt.Elem().Kind() == reflect.Uint8 {
				return strconv.ParseBool(string(reflect.ValueOf(a).Bytes()))
			}
			return false, errUnsupported
		default:
			return false, errUnsupported
		}
	}
}

type number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr | float32 | float64
}

func toNumber[T number](a any) (T, error) {
	a = indirect(a)
	switch v := a.(type) {
	case bool:
		return gcond.If[T](v, 1, 0), nil
	case nil:
		return 0, nil
	case int:
		return T(v), nil
	case int8:
		return T(v), nil
	case int16:
		return T(v), nil
	case int32:
		return T(v), nil
	case int64:
		return T(v), nil
	case uint:
		return T(v), nil
	case uint8:
		return T(v), nil
	case uint16:
		return T(v), nil
	case uint32:
		return T(v), nil
	case uint64:
		return T(v), nil
	case uintptr:
		return T(v), nil
	case float32:
		return T(v), nil
	case float64:
		return T(v), nil
	case string:
		return parseNumber[T](v)
	case []byte:
		return parseNumber[T](string(v))
	default:
		rt := reflect.TypeOf(a)
		switch rt.Kind() {
		case reflect.Bool:
			return gcond.If[T](reflect.ValueOf(a).Bool(), 1, 0), nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return T(reflect.ValueOf(a).Int()), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return T(reflect.ValueOf(a).Uint()), nil
		case reflect.Float32, reflect.Float64:
			return T(reflect.ValueOf(a).Float()), nil
		case reflect.String:
			return parseNumber[T](reflect.ValueOf(a).String())
		case reflect.Slice:
			if rt.Elem().Kind() == reflect.Uint8 {
				return parseNumber[T](string(reflect.ValueOf(a).Bytes()))
			}
			return 0, errUnsupported
		default:
			return 0, errUnsupported
		}
	}
}

func parseNumber[T number](s string) (T, error) {
	t := gvalue.Zero[T]()
	switch any(t).(type) {
	case int, int8, int16, int32, int64:
		tt, err := strconv.ParseInt(formatDecimalString(s), 10, 64)
		return T(tt), err
	case uint, uint8, uint16, uint32, uint64, uintptr:
		tt, err := strconv.ParseUint(formatDecimalString(s), 10, 64)
		return T(tt), err
	case float32, float64:
		tt, err := strconv.ParseFloat(s, 64)
		return T(tt), err
	default:
		return 0, errUnsupported
	}
}

func toString(a any) (string, error) {
	a = indirectToStringerOrError(a)
	switch v := a.(type) {
	case bool:
		return strconv.FormatBool(v), nil
	case nil:
		return "", nil
	case int:
		return strconv.Itoa(v), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case int32:
		return strconv.FormatInt(int64(v), 10), nil
	case int16:
		return strconv.FormatInt(int64(v), 10), nil
	case int8:
		return strconv.FormatInt(int64(v), 10), nil
	case uint:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint64:
		return strconv.FormatUint(v, 10), nil
	case uint32:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(v), 10), nil
	case uintptr:
		return strconv.FormatUint(uint64(v), 10), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32), nil
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	case fmt.Stringer:
		return v.String(), nil
	case error:
		return v.Error(), nil
	default:
		rt := reflect.TypeOf(a)
		switch rt.Kind() {
		case reflect.Bool:
			return strconv.FormatBool(reflect.ValueOf(a).Bool()), nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return strconv.FormatInt(reflect.ValueOf(a).Int(), 10), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return strconv.FormatUint(reflect.ValueOf(a).Uint(), 10), nil
		case reflect.Float32, reflect.Float64:
			return strconv.FormatFloat(reflect.ValueOf(a).Float(), 'f', -1, 64), nil
		case reflect.String:
			return reflect.ValueOf(a).String(), nil
		case reflect.Slice:
			if rt.Elem().Kind() == reflect.Uint8 {
				return string(reflect.ValueOf(a).Bytes()), nil
			}
			return "", errUnsupported
		default:
			return "", errUnsupported
		}
	}
}

// From html/template/content.go
// Copyright 2011 The Go Authors. All rights reserved.
//
// indirect returns the value, after dereferencing as many times
// as necessary to reach the base type (or nil).
func indirect(a any) any {
	if a == nil {
		return nil
	}
	if t := reflect.TypeOf(a); t.Kind() != reflect.Pointer {
		// Avoid creating a reflect.Value if it's not a pointer.
		return a
	}
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Pointer && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}

var (
	errorType       = reflect.TypeOf((*error)(nil)).Elem()
	fmtStringerType = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
)

// From html/template/content.go
// Copyright 2011 The Go Authors. All rights reserved.
//
// indirectToStringerOrError returns the value, after dereferencing as many times
// as necessary to reach the base type (or nil) or an implementation of fmt.Stringer
// or error.
func indirectToStringerOrError(a any) any {
	if a == nil {
		return nil
	}
	v := reflect.ValueOf(a)
	for !v.Type().Implements(fmtStringerType) && !v.Type().Implements(errorType) && v.Kind() == reflect.Pointer && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}

func formatDecimalString(numStr string) string {
	parts := strings.SplitN(numStr, ".", 2)
	if len(parts) == 1 {
		return numStr
	}
	integerPart, decimalPart := parts[0], parts[1]

	if integerPart == "" {
		integerPart = "0"
	}
	decimalPart = strings.TrimRight(decimalPart, "0")
	if decimalPart == "" {
		return integerPart
	}
	return integerPart + "." + decimalPart
}
