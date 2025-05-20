package skipmap

import (
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/bytedance/gg/internal/jsonbuilder"
)

// MarshalJSON returns s as the JSON encoding of s.
func (s *OrderedMap[keyT, valueT]) MarshalJSON() ([]byte, error) {
	if s == nil {
		return []byte("null"), nil
	}

	enc := jsonbuilder.NewDict()
	var err error
	s.Range(func(key keyT, value valueT) bool {
		err = enc.Store(key, value)
		return err == nil
	})
	if err != nil {
		return nil, err
	}
	return enc.Build()
}

// UnmarshalJSON sets *s to a copy of data.
func (s *OrderedMap[keyT, valueT]) UnmarshalJSON(data []byte) error {
	m := make(map[keyT]valueT)
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	for k, v := range m {
		s.Store(k, v)
	}
	return nil
}

// MarshalJSON returns s as the JSON encoding of s.
func (s *OrderedMapDesc[keyT, valueT]) MarshalJSON() ([]byte, error) {
	if s == nil {
		return []byte("null"), nil
	}

	enc := jsonbuilder.NewDict()
	var err error
	s.Range(func(key keyT, value valueT) bool {
		err = enc.Store(key, value)
		return err == nil
	})
	if err != nil {
		return nil, err
	}
	return enc.Build()
}

// UnmarshalJSON sets *s to a copy of data.
func (s *OrderedMapDesc[keyT, valueT]) UnmarshalJSON(data []byte) error {
	m := make(map[keyT]valueT)
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	for k, v := range m {
		s.Store(k, v)
	}
	return nil
}

// MarshalJSON returns s as the JSON encoding of s.
func (s *FuncMap[keyT, valueT]) MarshalJSON() ([]byte, error) {
	if s == nil {
		return []byte("null"), nil
	}

	enc := jsonbuilder.NewDict()
	var err error
	s.Range(func(key keyT, value valueT) bool {
		err = enc.Store(key, value)
		return err == nil
	})
	if err != nil {
		return nil, err
	}
	return enc.Build()
}

// UnmarshalJSON sets *s to a copy of data.
func (s *FuncMap[keyT, valueT]) UnmarshalJSON(data []byte) error {
	var (
		m            = make(map[string]valueT)
		zk           keyT
		unmarshalKey func(string) (keyT, error)
	)

	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	// See also: [encoding/json.(*decodeState).object]
	if _, ok := any(&zk).(encoding.TextUnmarshaler); ok {
		unmarshalKey = func(s string) (keyT, error) {
			var key keyT
			// TODO: Unsafe conv
			err := any(&key).(encoding.TextUnmarshaler).UnmarshalText([]byte(s))
			return key, err
		}
	} else {
		rk := reflect.ValueOf(zk)
		kt := rk.Type()
		switch rk.Kind() {
		case reflect.String:
			unmarshalKey = func(s string) (keyT, error) {
				return reflect.ValueOf(s).Convert(kt).Interface().(keyT), nil
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			unmarshalKey = func(s string) (keyT, error) {
				n, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					return zk, err
				}
				if rk.OverflowInt(n) {
					return zk, fmt.Errorf("%s overflow type %T", s, zk)
				}
				return reflect.ValueOf(n).Convert(kt).Interface().(keyT), nil
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			unmarshalKey = func(s string) (keyT, error) {
				n, err := strconv.ParseUint(s, 10, 64)
				if err != nil {
					return zk, err
				}
				if rk.OverflowUint(n) {
					return zk, fmt.Errorf("%s overflows type %T", s, zk)
				}
				return reflect.ValueOf(n).Convert(kt).Interface().(keyT), nil
			}
		default:
			return fmt.Errorf("unexpected key type: %T", zk)
		}
	}

	for ks, v := range m {
		k, err := unmarshalKey(ks)
		if err != nil {
			return err
		}
		s.Store(k, v)
	}

	return nil
}
