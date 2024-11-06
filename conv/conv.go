// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package conv

import (
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
