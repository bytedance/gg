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

package gconv

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"reflect"
	"strings"
	"testing"
	"time"
	"unsafe"

	"github.com/bytedance/gg/gcond"
	"github.com/bytedance/gg/gptr"
	"github.com/bytedance/gg/gresult"
	"github.com/bytedance/gg/gvalue"
	"github.com/bytedance/gg/internal/constraints"
)

func TestTo(t *testing.T) {
	testToBaseType(t)
	testToApproximatedType(t)
}

func testToBaseType(t *testing.T) {
	testToBool[bool](t)
	testToNumber[int](t)
	testToNumber[int8](t)
	testToNumber[int16](t)
	testToNumber[int32](t)
	testToNumber[int64](t)
	testToNumber[uint](t)
	testToNumber[uint8](t)
	testToNumber[uint16](t)
	testToNumber[uint32](t)
	testToNumber[uint64](t)
	testToNumber[uintptr](t)
	testToNumber[float32](t)
	testToNumber[float64](t)
	testToString[string](t)
}

type (
	MyBool       bool
	MyInt        int
	MyInt8       int8
	MyInt16      int16
	MyInt32      int32
	MyInt64      int64
	MyUint       uint
	MyUint8      uint8
	MyUint16     uint16
	MyUint32     uint32
	MyUint64     uint64
	MyUintptr    uintptr
	MyFloat32    float32
	MyFloat64    float64
	MyComplex64  complex64
	MyComplex128 complex128
	MyString     string
	MyBytes      []byte
)

func testToApproximatedType(t *testing.T) {
	testToBool[MyBool](t)
	testToNumber[MyInt](t)
	testToNumber[MyInt8](t)
	testToNumber[MyInt16](t)
	testToNumber[MyInt32](t)
	testToNumber[MyInt64](t)
	testToNumber[MyUint](t)
	testToNumber[MyUint8](t)
	testToNumber[MyUint16](t)
	testToNumber[MyUint32](t)
	testToNumber[MyUint64](t)
	testToNumber[MyUintptr](t)
	testToNumber[MyFloat32](t)
	testToNumber[MyFloat64](t)
	testToString[MyString](t)
}

func testToBool[T ~bool](t *testing.T) {
	t.Run(fmt.Sprintf("to bool %T", gvalue.Zero[T]()), func(t *testing.T) {
		tests := []struct {
			name    string
			input   any
			want    T
			wantErr bool
		}{
			{name: "bool false -> bool", input: false, want: false, wantErr: false},
			{name: "bool true -> bool", input: true, want: true, wantErr: false},
			{name: "MyBool false -> bool", input: MyBool(false), want: false, wantErr: false},
			{name: "MyBool true -> bool", input: MyBool(true), want: true, wantErr: false},
			{name: "nil -> bool", input: nil, want: false, wantErr: false},
			{name: "int 0 -> bool", input: 0, want: false, wantErr: false},
			{name: "int 1 -> bool", input: 1, want: true, wantErr: false},
			{name: "MyInt 0 -> bool", input: MyInt(0), want: false, wantErr: false},
			{name: "MyInt 1 -> bool", input: MyInt(1), want: true, wantErr: false},
			{name: "int8 0 -> bool", input: int8(0), want: false, wantErr: false},
			{name: "int8 1 -> bool", input: int8(1), want: true, wantErr: false},
			{name: "MyInt8 0 -> bool", input: MyInt8(0), want: false, wantErr: false},
			{name: "MyInt8 1 -> bool", input: MyInt8(1), want: true, wantErr: false},
			{name: "int16 0 -> bool", input: int16(0), want: false, wantErr: false},
			{name: "int16 1 -> bool", input: int16(1), want: true, wantErr: false},
			{name: "MyInt16 0 -> bool", input: MyInt16(0), want: false, wantErr: false},
			{name: "MyInt16 1 -> bool", input: MyInt16(1), want: true, wantErr: false},
			{name: "int32 0 -> bool", input: int32(0), want: false, wantErr: false},
			{name: "int32 1 -> bool", input: int32(1), want: true, wantErr: false},
			{name: "MyInt32 0 -> bool", input: MyInt32(0), want: false, wantErr: false},
			{name: "MyInt32 1 -> bool", input: MyInt32(1), want: true, wantErr: false},
			{name: "int64 0 -> bool", input: int64(0), want: false, wantErr: false},
			{name: "int64 1 -> bool", input: int64(1), want: true, wantErr: false},
			{name: "MyInt64 0 -> bool", input: MyInt64(0), want: false, wantErr: false},
			{name: "MyInt64 1 -> bool", input: MyInt64(1), want: true, wantErr: false},
			{name: "uint 0 -> bool", input: uint(0), want: false, wantErr: false},
			{name: "uint 1 -> bool", input: uint(1), want: true, wantErr: false},
			{name: "MyUint 0 -> bool", input: MyUint(0), want: false, wantErr: false},
			{name: "MyUint 1 -> bool", input: MyUint(1), want: true, wantErr: false},
			{name: "uint8 0 -> bool", input: uint8(0), want: false, wantErr: false},
			{name: "uint8 1 -> bool", input: uint8(1), want: true, wantErr: false},
			{name: "MyUint8 0 -> bool", input: MyUint8(0), want: false, wantErr: false},
			{name: "MyUint8 1 -> bool", input: MyUint8(1), want: true, wantErr: false},
			{name: "uint16 0 -> bool", input: uint16(0), want: false, wantErr: false},
			{name: "uint16 1 -> bool", input: uint16(1), want: true, wantErr: false},
			{name: "MyUint16 0 -> bool", input: MyUint16(0), want: false, wantErr: false},
			{name: "MyUint16 1 -> bool", input: MyUint16(1), want: true, wantErr: false},
			{name: "uint32 0 -> bool", input: uint32(0), want: false, wantErr: false},
			{name: "uint32 1 -> bool", input: uint32(1), want: true, wantErr: false},
			{name: "MyUint32 0 -> bool", input: MyUint32(0), want: false, wantErr: false},
			{name: "MyUint32 1 -> bool", input: MyUint32(1), want: true, wantErr: false},
			{name: "uint64 0 -> bool", input: uint64(0), want: false, wantErr: false},
			{name: "uint64 1 -> bool", input: uint64(1), want: true, wantErr: false},
			{name: "MyUint64 0 -> bool", input: MyUint64(0), want: false, wantErr: false},
			{name: "MyUint64 1 -> bool", input: MyUint64(1), want: true, wantErr: false},
			{name: "uintptr 0 -> bool", input: uintptr(0), want: false, wantErr: false},
			{name: "uintptr 1 -> bool", input: uintptr(1), want: true, wantErr: false},
			{name: "MyUintptr 0 -> bool", input: MyUintptr(0), want: false, wantErr: false},
			{name: "MyUintptr 1 -> bool", input: MyUintptr(1), want: true, wantErr: false},
			{name: "float32 0 -> bool", input: float32(0), want: false, wantErr: false},
			{name: "float32 1 -> bool", input: float32(1), want: true, wantErr: false},
			{name: "MyFloat32 0 -> bool", input: MyFloat32(0), want: false, wantErr: false},
			{name: "MyFloat32 1 -> bool", input: MyFloat32(1), want: true, wantErr: false},
			{name: "float64 0 -> bool", input: float64(0), want: false, wantErr: false},
			{name: "float64 1 -> bool", input: float64(1), want: true, wantErr: false},
			{name: "MyFloat64 0 -> bool", input: MyFloat64(0), want: false, wantErr: false},
			{name: "MyFloat64 1 -> bool", input: MyFloat64(1), want: true, wantErr: false},
			{name: "complex64 0+0i -> bool", input: complex64(0), want: false, wantErr: false},
			{name: "complex64 1+0i -> bool", input: complex64(1), want: true, wantErr: false},
			{name: "complex64 0+1i -> bool", input: complex64(1i), want: true, wantErr: false},
			{name: "MyComplex64 0+0i -> bool", input: MyComplex64(0), want: false, wantErr: false},
			{name: "MyComplex64 1+0i -> bool", input: MyComplex64(1), want: true, wantErr: false},
			{name: "MyComplex64 0+1i -> bool", input: MyComplex64(1i), want: true, wantErr: false},
			{name: "complex128 0+0i -> bool", input: complex128(0), want: false, wantErr: false},
			{name: "complex128 1+0i -> bool", input: complex128(1), want: true, wantErr: false},
			{name: "complex128 0+1i -> bool", input: complex128(1i), want: true, wantErr: false},
			{name: "MyComplex128 0+0i -> bool", input: MyComplex128(0), want: false, wantErr: false},
			{name: "MyComplex128 1+0i -> bool", input: MyComplex128(1), want: true, wantErr: false},
			{name: "MyComplex128 0+1i -> bool", input: MyComplex128(1i), want: true, wantErr: false},
			{name: "string false -> bool", input: "false", want: false, wantErr: false},
			{name: "string true -> bool", input: "true", want: true, wantErr: false},
			{name: "MyString false -> bool", input: MyString("false"), want: false, wantErr: false},
			{name: "MyString true -> bool", input: MyString("true"), want: true, wantErr: false},
			{name: "[]byte false -> bool", input: []byte("false"), want: false, wantErr: false},
			{name: "[]byte true -> bool", input: []byte("true"), want: true, wantErr: false},
			{name: "MyBytes false -> bool", input: MyBytes("false"), want: false, wantErr: false},
			{name: "MyBytes true -> bool", input: MyBytes("true"), want: true, wantErr: false},
			{name: "other -> bool", input: []int{}, want: false, wantErr: true},
			{name: "*string true -> bool", input: gptr.Of("true"), want: true, wantErr: false},
			{name: "**string true -> bool", input: gptr.Of(gptr.Of("true")), want: true, wantErr: false},
			{name: "***string true -> bool", input: gptr.Of(gptr.Of(gptr.Of("true"))), want: true, wantErr: false},
			{name: "array -> bool", input: [2]int{}, want: false, wantErr: true},
			{name: "unbuffered chan -> bool", input: make(chan int), want: false, wantErr: true},
			{name: "buffered chan -> bool", input: make(chan int, 1), want: false, wantErr: true},
			{name: "func -> bool", input: func() {}, want: false, wantErr: true},
			{name: "map -> bool", input: map[int]int{}, want: false, wantErr: true},
			{name: "slice -> bool", input: []int{}, want: false, wantErr: true},
			{name: "struct -> bool", input: strings.Builder{}, want: false, wantErr: true},
			{name: "unsafe pointer -> bool", input: unsafe.Pointer(gptr.Of(1)), want: false, wantErr: true},
		}
		for _, tt := range tests {
			if got := To[T](tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("To() = %v, want %v", got, tt.want)
			}

			wantPtr := gcond.If(tt.wantErr, nil, &tt.want)
			if got := ToPtr[T](tt.input); !reflect.DeepEqual(got, wantPtr) {
				t.Errorf("ToPtr() = %v, want %v", got, wantPtr)
			}

			wantR := gcond.If(tt.wantErr, gresult.Err[T](errUnsupported), gresult.OK(tt.want))
			if got := ToR[T](tt.input); !reflect.DeepEqual(got, wantR) {
				t.Errorf("ToR() = %v, want %v", got, wantR)
			}

			got, err := ToE[T](tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToE() got = %v, want %v", got, tt.want)
			}
		}
	})
}

func testToNumber[T constraints.Number](t *testing.T) {
	t.Run(fmt.Sprintf("to number %T", gvalue.Zero[T]()), func(t *testing.T) {
		tests := []struct {
			name    string
			input   any
			want    T
			wantErr bool
		}{
			{name: "bool false -> %T", input: false, want: T(0), wantErr: false},
			{name: "bool true -> %T", input: true, want: T(1), wantErr: false},
			{name: "MyBool false -> %T", input: MyBool(false), want: T(0), wantErr: false},
			{name: "MyBool true -> %T", input: MyBool(true), want: T(1), wantErr: false},
			{name: "nil -> %T", input: nil, want: T(0), wantErr: false},
			{name: "int 0 -> %T", input: 0, want: T(0), wantErr: false},
			{name: "int 1 -> %T", input: 1, want: T(1), wantErr: false},
			{name: "MyInt 0 -> %T", input: MyInt(0), want: T(0), wantErr: false},
			{name: "MyInt 1 -> %T", input: MyInt(1), want: T(1), wantErr: false},
			{name: "int8 0 -> %T", input: int8(0), want: T(0), wantErr: false},
			{name: "int8 1 -> %T", input: int8(1), want: T(1), wantErr: false},
			{name: "MyInt8 0 -> %T", input: MyInt8(0), want: T(0), wantErr: false},
			{name: "MyInt8 1 -> %T", input: MyInt8(1), want: T(1), wantErr: false},
			{name: "int16 0 -> %T", input: int16(0), want: T(0), wantErr: false},
			{name: "int16 1 -> %T", input: int16(1), want: T(1), wantErr: false},
			{name: "MyInt16 0 -> %T", input: MyInt16(0), want: T(0), wantErr: false},
			{name: "MyInt16 1 -> %T", input: MyInt16(1), want: T(1), wantErr: false},
			{name: "int32 0 -> %T", input: int32(0), want: T(0), wantErr: false},
			{name: "int32 1 -> %T", input: int32(1), want: T(1), wantErr: false},
			{name: "MyInt32 0 -> %T", input: MyInt32(0), want: T(0), wantErr: false},
			{name: "MyInt32 1 -> %T", input: MyInt32(1), want: T(1), wantErr: false},
			{name: "int64 0 -> %T", input: int64(0), want: T(0), wantErr: false},
			{name: "int64 1 -> %T", input: int64(1), want: T(1), wantErr: false},
			{name: "MyInt64 0 -> %T", input: MyInt64(0), want: T(0), wantErr: false},
			{name: "MyInt64 1 -> %T", input: MyInt64(1), want: T(1), wantErr: false},
			{name: "uint 0 -> %T", input: uint(0), want: T(0), wantErr: false},
			{name: "uint 1 -> %T", input: uint(1), want: T(1), wantErr: false},
			{name: "MyUint 0 -> %T", input: MyUint(0), want: T(0), wantErr: false},
			{name: "MyUint 1 -> %T", input: MyUint(1), want: T(1), wantErr: false},
			{name: "uint8 0 -> %T", input: uint8(0), want: T(0), wantErr: false},
			{name: "uint8 1 -> %T", input: uint8(1), want: T(1), wantErr: false},
			{name: "MyUint8 0 -> %T", input: MyUint8(0), want: T(0), wantErr: false},
			{name: "MyUint8 1 -> %T", input: MyUint8(1), want: T(1), wantErr: false},
			{name: "uint16 0 -> %T", input: uint16(0), want: T(0), wantErr: false},
			{name: "uint16 1 -> %T", input: uint16(1), want: T(1), wantErr: false},
			{name: "MyUint16 0 -> %T", input: MyUint16(0), want: T(0), wantErr: false},
			{name: "MyUint16 1 -> %T", input: MyUint16(1), want: T(1), wantErr: false},
			{name: "uint32 0 -> %T", input: uint32(0), want: T(0), wantErr: false},
			{name: "uint32 1 -> %T", input: uint32(1), want: T(1), wantErr: false},
			{name: "MyUint32 0 -> %T", input: MyUint32(0), want: T(0), wantErr: false},
			{name: "MyUint32 1 -> %T", input: MyUint32(1), want: T(1), wantErr: false},
			{name: "uint64 0 -> %T", input: uint64(0), want: T(0), wantErr: false},
			{name: "uint64 1 -> %T", input: uint64(1), want: T(1), wantErr: false},
			{name: "MyUint64 0 -> %T", input: MyUint64(0), want: T(0), wantErr: false},
			{name: "MyUint64 1 -> %T", input: MyUint64(1), want: T(1), wantErr: false},
			{name: "uintptr 0 -> %T", input: uintptr(0), want: T(0), wantErr: false},
			{name: "uintptr 1 -> %T", input: uintptr(1), want: T(1), wantErr: false},
			{name: "MyUintptr 0 -> %T", input: MyUintptr(0), want: T(0), wantErr: false},
			{name: "MyUintptr 1 -> %T", input: MyUintptr(1), want: T(1), wantErr: false},
			{name: "float32 0 -> %T", input: float32(0), want: T(0), wantErr: false},
			{name: "float32 1 -> %T", input: float32(1), want: T(1), wantErr: false},
			{name: "MyFloat32 0 -> %T", input: MyFloat32(0), want: T(0), wantErr: false},
			{name: "MyFloat32 1 -> %T", input: MyFloat32(1), want: T(1), wantErr: false},
			{name: "float64 0 -> %T", input: float64(0), want: T(0), wantErr: false},
			{name: "float64 1 -> %T", input: float64(1), want: T(1), wantErr: false},
			{name: "MyFloat64 0 -> %T", input: MyFloat64(0), want: T(0), wantErr: false},
			{name: "MyFloat64 1 -> %T", input: MyFloat64(1), want: T(1), wantErr: false},
			{name: "string 0 -> %T", input: "0", want: T(0), wantErr: false},
			{name: "string 1 -> %T", input: "1", want: T(1), wantErr: false},
			{name: "string 0.0 -> %T", input: "0.0", want: T(0), wantErr: false},
			{name: "string 1.0 -> %T", input: "1.0", want: T(1), wantErr: false},
			{name: "MyString 0 -> %T", input: MyString("0"), want: T(0), wantErr: false},
			{name: "MyString 1 -> %T", input: MyString("1"), want: T(1), wantErr: false},
			{name: "MyString 0.0 -> %T", input: MyString("0.0"), want: T(0), wantErr: false},
			{name: "MyString 1.0 -> %T", input: MyString("1.0"), want: T(1), wantErr: false},
			{name: "json.Number 0 -> %T", input: json.Number("0"), want: T(0), wantErr: false},
			{name: "json.Number 1 -> %T", input: json.Number("1"), want: T(1), wantErr: false},
			{name: "time.Duration 0 -> %T", input: time.Duration(0), want: T(0), wantErr: false},
			{name: "time.Duration 1 -> %T", input: time.Duration(1), want: T(1), wantErr: false},
			{name: "[]byte 0.0 -> %T", input: []byte("0.0"), want: T(0), wantErr: false},
			{name: "[]byte 1.0 -> %T", input: []byte("1.0"), want: T(1), wantErr: false},
			{name: "MyBytes 0.0 -> %T", input: MyBytes("0.0"), want: T(0), wantErr: false},
			{name: "MyBytes 1.0 -> %T", input: MyBytes("1.0"), want: T(1), wantErr: false},
			{name: "other -> %T", input: []int{}, want: T(0), wantErr: true},
			{name: "*string 1.0 -> %T", input: gptr.Of("1.0"), want: T(1), wantErr: false},
			{name: "**string 1.0 -> %T", input: gptr.Of(gptr.Of("1.0")), want: T(1), wantErr: false},
			{name: "***string 1.0 -> %T", input: gptr.Of(gptr.Of(gptr.Of("1.0"))), want: T(1), wantErr: false},
		}
		for _, tt := range tests {
			if got := To[T](tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("To() = %v, want %v", got, tt.want)
			}

			wantPtr := gcond.If(tt.wantErr, nil, &tt.want)
			if got := ToPtr[T](tt.input); !reflect.DeepEqual(got, wantPtr) {
				t.Errorf("ToPtr() = %v, want %v", got, wantPtr)
			}

			wantR := gcond.If(tt.wantErr, gresult.Err[T](errUnsupported), gresult.OK(tt.want))
			if got := ToR[T](tt.input); !reflect.DeepEqual(got, wantR) {
				t.Errorf("ToR() = %v, want %v", got, wantR)
			}

			got, err := ToE[T](tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToE() got = %v, want %v", got, tt.want)
			}
		}
	})
}

func testToString[T ~string](t *testing.T) {
	t.Run(fmt.Sprintf("to string %T", gvalue.Zero[T]()), func(t *testing.T) {
		tests := []struct {
			name    string
			input   any
			want    T
			wantErr bool
		}{
			{name: "bool false -> string", input: false, want: T("false"), wantErr: false},
			{name: "bool true -> string", input: true, want: T("true"), wantErr: false},
			{name: "MyBool false -> string", input: MyBool(false), want: T("false"), wantErr: false},
			{name: "MyBool true -> string", input: MyBool(true), want: T("true"), wantErr: false},
			{name: "nil -> string", input: nil, want: T(""), wantErr: false},
			{name: "int 0 -> string", input: 0, want: T("0"), wantErr: false},
			{name: "int 1 -> string", input: 1, want: T("1"), wantErr: false},
			{name: "MyInt 0 -> string", input: MyInt(0), want: T("0"), wantErr: false},
			{name: "MyInt 1 -> string", input: MyInt(1), want: T("1"), wantErr: false},
			{name: "int8 0 -> string", input: int8(0), want: T("0"), wantErr: false},
			{name: "int8 1 -> string", input: int8(1), want: T("1"), wantErr: false},
			{name: "MyInt8 0 -> string", input: MyInt8(0), want: T("0"), wantErr: false},
			{name: "MyInt8 1 -> string", input: MyInt8(1), want: T("1"), wantErr: false},
			{name: "int16 0 -> string", input: int16(0), want: T("0"), wantErr: false},
			{name: "int16 1 -> string", input: int16(1), want: T("1"), wantErr: false},
			{name: "MyInt16 0 -> string", input: MyInt16(0), want: T("0"), wantErr: false},
			{name: "MyInt16 1 -> string", input: MyInt16(1), want: T("1"), wantErr: false},
			{name: "int32 0 -> string", input: int32(0), want: T("0"), wantErr: false},
			{name: "int32 1 -> string", input: int32(1), want: T("1"), wantErr: false},
			{name: "MyInt32 0 -> string", input: MyInt32(0), want: T("0"), wantErr: false},
			{name: "MyInt32 1 -> string", input: MyInt32(1), want: T("1"), wantErr: false},
			{name: "int64 0 -> string", input: int64(0), want: T("0"), wantErr: false},
			{name: "int64 1 -> string", input: int64(1), want: T("1"), wantErr: false},
			{name: "MyInt64 0 -> string", input: MyInt64(0), want: T("0"), wantErr: false},
			{name: "MyInt64 1 -> string", input: MyInt64(1), want: T("1"), wantErr: false},
			{name: "uint 0 -> string", input: uint(0), want: T("0"), wantErr: false},
			{name: "uint 1 -> string", input: uint(1), want: T("1"), wantErr: false},
			{name: "MyUint 0 -> string", input: MyUint(0), want: T("0"), wantErr: false},
			{name: "MyUint 1 -> string", input: MyUint(1), want: T("1"), wantErr: false},
			{name: "uint8 0 -> string", input: uint8(0), want: T("0"), wantErr: false},
			{name: "uint8 1 -> string", input: uint8(1), want: T("1"), wantErr: false},
			{name: "MyUint8 0 -> string", input: MyUint8(0), want: T("0"), wantErr: false},
			{name: "MyUint8 1 -> string", input: MyUint8(1), want: T("1"), wantErr: false},
			{name: "uint16 0 -> string", input: uint16(0), want: T("0"), wantErr: false},
			{name: "uint16 1 -> string", input: uint16(1), want: T("1"), wantErr: false},
			{name: "MyUint16 0 -> string", input: MyUint16(0), want: T("0"), wantErr: false},
			{name: "MyUint16 1 -> string", input: MyUint16(1), want: T("1"), wantErr: false},
			{name: "uint32 0 -> string", input: uint32(0), want: T("0"), wantErr: false},
			{name: "uint32 1 -> string", input: uint32(1), want: T("1"), wantErr: false},
			{name: "MyUint32 0 -> string", input: MyUint32(0), want: T("0"), wantErr: false},
			{name: "MyUint32 1 -> string", input: MyUint32(1), want: T("1"), wantErr: false},
			{name: "uint64 0 -> string", input: uint64(0), want: T("0"), wantErr: false},
			{name: "uint64 1 -> string", input: uint64(1), want: T("1"), wantErr: false},
			{name: "MyUint64 0 -> string", input: MyUint64(0), want: T("0"), wantErr: false},
			{name: "MyUint64 1 -> string", input: MyUint64(1), want: T("1"), wantErr: false},
			{name: "uintptr 0 -> string", input: uintptr(0), want: T("0"), wantErr: false},
			{name: "uintptr 1 -> string", input: uintptr(1), want: T("1"), wantErr: false},
			{name: "MyUintptr 0 -> string", input: MyUintptr(0), want: T("0"), wantErr: false},
			{name: "MyUintptr 1 -> string", input: MyUintptr(1), want: T("1"), wantErr: false},
			{name: "float32 0 -> string", input: float32(0), want: T("0"), wantErr: false},
			{name: "float32 1 -> string", input: float32(1), want: T("1"), wantErr: false},
			{name: "MyFloat32 0 -> string", input: MyFloat32(0), want: T("0"), wantErr: false},
			{name: "MyFloat32 1 -> string", input: MyFloat32(1), want: T("1"), wantErr: false},
			{name: "float64 0 -> string", input: float64(0), want: T("0"), wantErr: false},
			{name: "float64 1 -> string", input: float64(1), want: T("1"), wantErr: false},
			{name: "MyFloat64 0 -> string", input: MyFloat64(0), want: T("0"), wantErr: false},
			{name: "MyFloat64 1 -> string", input: MyFloat64(1), want: T("1"), wantErr: false},
			{name: "string -> string", input: "xxx", want: T("xxx"), wantErr: false},
			{name: "MyString -> string", input: MyString("xxx"), want: T("xxx"), wantErr: false},
			{name: "[]byte -> string", input: []byte("xxx"), want: T("xxx"), wantErr: false},
			{name: "MyBytes -> string", input: MyBytes("xxx"), want: T("xxx"), wantErr: false},
			{name: "json.Number 0 -> %T", input: json.Number("0"), want: T("0"), wantErr: false},
			{name: "json.Number 1 -> %T", input: json.Number("1"), want: T("1"), wantErr: false},
			{name: "fmt.Stringer -> string", input: net.IPv4(8, 8, 8, 8), want: T("8.8.8.8"), wantErr: false},
			{name: "error -> string", input: errors.New("zzz"), want: T("zzz"), wantErr: false},
			{name: "other -> string", input: []int{}, want: T(""), wantErr: true},
			{name: "*bool true -> string", input: gptr.Of(true), want: T("true"), wantErr: false},
			{name: "**bool true -> string", input: gptr.Of(gptr.Of(true)), want: T("true"), wantErr: false},
			{name: "***bool true -> string", input: gptr.Of(gptr.Of(gptr.Of(true))), want: T("true"), wantErr: false},
		}
		for _, tt := range tests {
			if got := To[T](tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("To() = %v, want %v", got, tt.want)
			}

			wantPtr := gcond.If(tt.wantErr, nil, &tt.want)
			if got := ToPtr[T](tt.input); !reflect.DeepEqual(got, wantPtr) {
				t.Errorf("ToPtr() = %v, want %v", got, wantPtr)
			}

			wantR := gcond.If(tt.wantErr, gresult.Err[T](errUnsupported), gresult.OK(tt.want))
			if got := ToR[T](tt.input); !reflect.DeepEqual(got, wantR) {
				t.Errorf("ToR() = %v, want %v", got, wantR)
			}

			got, err := ToE[T](tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToE() got = %v, want %v", got, tt.want)
			}
		}
	})
}
