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
	"math"
	"strconv"
	"testing"
)

func checkSignedString[T ~int8 | ~int16 | ~int32 | ~int64](t *testing.T, bitSize int, values []string) {
	t.Helper()
	for _, s := range values {
		want, wantErr := strconv.ParseInt(formatDecimalString(s), 10, bitSize)
		if wantErr != nil {
			want = 0 // gconv promises the zero value on conversion failure.
		}
		got, gotErr := ToE[T](s)
		if (gotErr != nil) != (wantErr != nil) || int64(got) != want {
			t.Errorf("ToE[%T](%q) = %v, %v; want %v and error=%v", got, s, got, gotErr, want, wantErr != nil)
		} else if gotErr != nil && gotErr.Error() != wantErr.Error() {
			t.Errorf("ToE[%T](%q) error = %q, want %q", got, s, gotErr, wantErr)
		}
		gotBytes, gotBytesErr := ToE[T]([]byte(s))
		if (gotBytesErr != nil) != (wantErr != nil) || int64(gotBytes) != want {
			t.Errorf("ToE[%T]([]byte(%q)) = %v, %v; want %v and error=%v", gotBytes, s, gotBytes, gotBytesErr, want, wantErr != nil)
		} else if gotBytesErr != nil && gotBytesErr.Error() != wantErr.Error() {
			t.Errorf("ToE[%T]([]byte(%q)) error = %q, want %q", gotBytes, s, gotBytesErr, wantErr)
		}
		if wantErr != nil {
			if got := To[T](s); got != 0 {
				t.Errorf("To[%T](%q) = %v, want zero value", got, s, got)
			}
			if got := ToPtr[T](s); got != nil {
				t.Errorf("ToPtr[%T](%q) = %v, want nil", *got, s, got)
			}
		}
	}
}

func checkUnsignedString[T ~uint8 | ~uint16 | ~uint32 | ~uint64](t *testing.T, bitSize int, values []string) {
	t.Helper()
	for _, s := range values {
		want, wantErr := strconv.ParseUint(formatDecimalString(s), 10, bitSize)
		if wantErr != nil {
			want = 0 // gconv promises the zero value on conversion failure.
		}
		got, gotErr := ToE[T](s)
		if (gotErr != nil) != (wantErr != nil) || uint64(got) != want {
			t.Errorf("ToE[%T](%q) = %v, %v; want %v and error=%v", got, s, got, gotErr, want, wantErr != nil)
		} else if gotErr != nil && gotErr.Error() != wantErr.Error() {
			t.Errorf("ToE[%T](%q) error = %q, want %q", got, s, gotErr, wantErr)
		}
		gotBytes, gotBytesErr := ToE[T]([]byte(s))
		if (gotBytesErr != nil) != (wantErr != nil) || uint64(gotBytes) != want {
			t.Errorf("ToE[%T]([]byte(%q)) = %v, %v; want %v and error=%v", gotBytes, s, gotBytes, gotBytesErr, want, wantErr != nil)
		} else if gotBytesErr != nil && gotBytesErr.Error() != wantErr.Error() {
			t.Errorf("ToE[%T]([]byte(%q)) error = %q, want %q", gotBytes, s, gotBytesErr, wantErr)
		}
		if wantErr != nil {
			if got := To[T](s); got != 0 {
				t.Errorf("To[%T](%q) = %v, want zero value", got, s, got)
			}
			if got := ToPtr[T](s); got != nil {
				t.Errorf("ToPtr[%T](%q) = %v, want nil", *got, s, got)
			}
		}
	}
}

func checkFloatString[T ~float32 | ~float64](t *testing.T, bitSize int, values []string) {
	t.Helper()
	for _, s := range values {
		want, wantErr := strconv.ParseFloat(s, bitSize)
		if wantErr != nil {
			want = 0 // gconv promises the zero value on conversion failure.
		} else if bitSize == 32 {
			want = float64(float32(want))
		}
		got, gotErr := ToE[T](s)
		got64 := float64(got)
		if (gotErr != nil) != (wantErr != nil) || math.Float64bits(got64) != math.Float64bits(want) {
			t.Errorf("ToE[%T](%q) = %v, %v; want %v and error=%v", got, s, got, gotErr, want, wantErr != nil)
		} else if gotErr != nil && gotErr.Error() != wantErr.Error() {
			t.Errorf("ToE[%T](%q) error = %q, want %q", got, s, gotErr, wantErr)
		}
		if wantErr != nil {
			if got := To[T](s); got != 0 {
				t.Errorf("To[%T](%q) = %v, want zero value", got, s, got)
			}
			if got := ToPtr[T](s); got != nil {
				t.Errorf("ToPtr[%T](%q) = %v, want nil", *got, s, got)
			}
		}
	}
}

func TestToENumberStringRange(t *testing.T) {
	signed8 := []string{"-129", "-128", "127", "128", "255", "256", "127.0", "128.0"}
	unsigned8 := []string{"-1", "0", "255", "256", "511", "255.0", "256.0"}
	floats := []string{"3.4028235e38", "3.4028236e38", "1e39", "-1e39", "1e308", "1e309"}

	t.Run("int8", func(t *testing.T) { checkSignedString[int8](t, 8, signed8) })
	t.Run("named int8", func(t *testing.T) { checkSignedString[MyInt8](t, 8, signed8) })
	t.Run("int16", func(t *testing.T) {
		checkSignedString[int16](t, 16, []string{"-32769", "-32768", "32767", "32768"})
	})
	t.Run("named int16", func(t *testing.T) {
		checkSignedString[MyInt16](t, 16, []string{"-32769", "-32768", "32767", "32768"})
	})
	t.Run("int32", func(t *testing.T) {
		checkSignedString[int32](t, 32, []string{"-2147483649", "-2147483648", "2147483647", "2147483648"})
	})
	t.Run("named int32", func(t *testing.T) {
		checkSignedString[MyInt32](t, 32, []string{"-2147483649", "-2147483648", "2147483647", "2147483648"})
	})
	t.Run("int64", func(t *testing.T) {
		checkSignedString[int64](t, 64, []string{"-9223372036854775809", "-9223372036854775808", "9223372036854775807", "9223372036854775808"})
	})
	t.Run("named int64", func(t *testing.T) {
		checkSignedString[MyInt64](t, 64, []string{"-9223372036854775809", "-9223372036854775808", "9223372036854775807", "9223372036854775808"})
	})
	t.Run("uint8", func(t *testing.T) { checkUnsignedString[uint8](t, 8, unsigned8) })
	t.Run("named uint8", func(t *testing.T) { checkUnsignedString[MyUint8](t, 8, unsigned8) })
	t.Run("uint16", func(t *testing.T) {
		checkUnsignedString[uint16](t, 16, []string{"-1", "0", "65535", "65536"})
	})
	t.Run("named uint16", func(t *testing.T) {
		checkUnsignedString[MyUint16](t, 16, []string{"-1", "0", "65535", "65536"})
	})
	t.Run("uint32", func(t *testing.T) {
		checkUnsignedString[uint32](t, 32, []string{"-1", "0", "4294967295", "4294967296"})
	})
	t.Run("named uint32", func(t *testing.T) {
		checkUnsignedString[MyUint32](t, 32, []string{"-1", "0", "4294967295", "4294967296"})
	})
	t.Run("uint64", func(t *testing.T) {
		checkUnsignedString[uint64](t, 64, []string{"-1", "0", "18446744073709551615", "18446744073709551616"})
	})
	t.Run("named uint64", func(t *testing.T) {
		checkUnsignedString[MyUint64](t, 64, []string{"-1", "0", "18446744073709551615", "18446744073709551616"})
	})
	t.Run("float32", func(t *testing.T) { checkFloatString[float32](t, 32, floats) })
	t.Run("named float32", func(t *testing.T) { checkFloatString[MyFloat32](t, 32, floats) })
	t.Run("float64", func(t *testing.T) { checkFloatString[float64](t, 64, floats) })
	t.Run("named float64", func(t *testing.T) { checkFloatString[MyFloat64](t, 64, floats) })
}

func TestToEIntegerStringMatchesStrconv(t *testing.T) {
	values := make([]string, 0, 2049)
	for i := -1024; i <= 1024; i++ {
		values = append(values, strconv.Itoa(i))
	}

	checkSignedString[int8](t, 8, values)
	checkSignedString[MyInt8](t, 8, values)
	checkSignedString[int16](t, 16, values)
	checkSignedString[MyInt16](t, 16, values)
	checkUnsignedString[uint8](t, 8, values)
	checkUnsignedString[MyUint8](t, 8, values)
	checkUnsignedString[uint16](t, 16, values)
	checkUnsignedString[MyUint16](t, 16, values)
}
