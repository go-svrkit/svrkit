// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package conv

import (
	"unsafe"
)

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
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
