// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package conv

import (
	"encoding/json"
	"fmt"
	"reflect"
	"unsafe"
)

type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Integer interface {
	Signed | Unsigned
}

type Float interface {
	~float32 | ~float64
}

type Number interface {
	Integer | Float
}

// StrAsBytes returns the bytes backing a string, it is the caller's responsibility not to mutate the bytes returned.
// see https://pkg.go.dev/unsafe#Pointer rule(6)
func StrAsBytes(s string) []byte {
	if len(s) == 0 {
		return nil
	}
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// BytesAsStr returns the string view of byte slice
func BytesAsStr(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// BoolTo converts a bool to integer
func BoolTo[T Integer](b bool) T {
	if b {
		return 1
	}
	return 0
}

// IntToBool converts an integer to bool
func IntToBool[T Integer](v T) bool {
	return v != 0
}

// ConvTo 转换为整数
func ConvTo[T Integer](val any) T {
	switch v := val.(type) {
	case int8:
		return T(v)
	case uint8:
		return T(v)
	case int16:
		return T(v)
	case uint16:
		return T(v)
	case int:
		return T(v)
	case uint:
		return T(v)
	case int32:
		return T(v)
	case uint32:
		return T(v)
	case int64:
		return T(v)
	case uint64:
		return T(v)
	case float32:
		return T(v)
	case float64:
		return T(v)
	case bool:
		if v {
			return T(1)
		}
		return T(0)
	case string:
		n, _ := ParseTo[T](v)
		return n
	}
	return 0
}

func ReflectConvAny(rtype reflect.Type, body string) (rval reflect.Value, err error) {
	var kind = rtype.Kind()
	switch kind {
	case reflect.Ptr:
		return ReflectConvAny(rtype.Elem(), body)

	case reflect.Bool:
		rval = reflect.ValueOf(ParseBool(body))

	case reflect.String:
		rval = reflect.ValueOf(body)

	case reflect.Int:
		var n int64
		if n, err = ParseI64(body); err != nil {
			return
		}
		rval = reflect.ValueOf(int(n))

	case reflect.Uint:
		var n uint64
		if n, err = ParseU64(body); err != nil {
			return
		}
		rval = reflect.ValueOf(uint(n))

	case reflect.Int8:
		var n int8
		if n, err = ParseI8(body); err != nil {
			return
		}
		rval = reflect.ValueOf(n)

	case reflect.Uint8:
		var n uint8
		if n, err = ParseU8(body); err != nil {
			return
		}
		rval = reflect.ValueOf(n)

	case reflect.Int16:
		var n int16
		if n, err = ParseI16(body); err != nil {
			return
		}
		rval = reflect.ValueOf(n)

	case reflect.Uint16:
		var n uint16
		if n, err = ParseU16(body); err != nil {
			return
		}
		rval = reflect.ValueOf(n)

	case reflect.Int32:
		var n int32
		if n, err = ParseI32(body); err != nil {
			return
		}
		rval = reflect.ValueOf(n)

	case reflect.Uint32:
		var n uint32
		if n, err = ParseU32(body); err != nil {
			return
		}
		rval = reflect.ValueOf(n)

	case reflect.Int64:
		var n int64
		if n, err = ParseI64(body); err != nil {
			return
		}
		rval = reflect.ValueOf(n)

	case reflect.Uint64:
		var n uint64
		if n, err = ParseU64(body); err != nil {
			return
		}
		rval = reflect.ValueOf(n)

	case reflect.Float32:
		var f float32
		if f, err = ParseF32(body); err != nil {
			return
		}
		rval = reflect.ValueOf(f)

	case reflect.Float64:
		var f float64
		if f, err = ParseF64(body); err != nil {
			return
		}
		rval = reflect.ValueOf(f)
		return

	case reflect.Slice, reflect.Map, reflect.Struct:
		rval = reflect.New(rtype)
		if body != "" {
			if err = json.Unmarshal(StrAsBytes(body), rval.Interface()); err != nil {
				return
			}
		}
		rval = rval.Elem()

	default:
		err = fmt.Errorf("unsupported type %v", kind)
	}
	return
}
