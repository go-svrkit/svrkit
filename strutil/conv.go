// Copyright © 2020 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package strutil

import (
	"cmp"
	"fmt"
	"log"
	"math"
	"strconv"

	"gopkg.in/svrkit.v1/slog"
)

// ParseBool parse string to bool
func ParseBool(s string) bool {
	if len(s) == 0 {
		return false
	}
	switch s {
	case "y", "Y", "on", "ON", "yes", "YES":
		return true
	}
	b, _ := strconv.ParseBool(s) // "1", "t", "T", "true", "TRUE", "True"
	return b
}

// ParseI8 parse string to int8
func ParseI8(s string) (int8, error) {
	n, err := ParseI32(s)
	if err != nil {
		return 0, err
	}
	if n > math.MaxInt8 || n < math.MinInt8 {
		return 0, fmt.Errorf("ParseI8: value %s out of range", s)
	}
	return int8(n), nil
}

// MustParseI8 parse string to int8, panic if error
func MustParseI8(s string) int8 {
	n, err := ParseI32(s)
	if err != nil || n > math.MaxInt8 || n < math.MinInt8 {
		log.Panicf("MustParseI8: value %s out of range", s)
		return 0
	}
	return int8(n)
}

// ParseU8 parse string to uint8
func ParseU8(s string) (uint8, error) {
	n, err := ParseI32(s)
	if err != nil {
		return 0, err
	}
	if n > math.MaxUint8 || n < 0 {
		return 0, fmt.Errorf("ParseU8: value %s out of range", s)
	}
	return uint8(n), nil
}

// MustParseU8 parse string to uint8, panic if error
func MustParseU8(s string) uint8 {
	n, err := ParseI32(s)
	if err != nil || n > math.MaxUint8 || n < 0 {
		log.Panicf("MustParseU8: value %s out of range", s)
	}
	return uint8(n)
}

// ParseI16 parse string to int16
func ParseI16(s string) (int16, error) {
	n, err := ParseI32(s)
	if err != nil {
		return 0, err
	}
	if n > math.MaxInt16 || n < math.MinInt16 {
		return 0, fmt.Errorf("ParseI16: value %s out of range", s)
	}
	return int16(n), nil
}

// MustParseI16 parse string to int16, panic if error
func MustParseI16(s string) int16 {
	n, err := ParseI32(s)
	if err != nil || n > math.MaxInt16 || n < math.MinInt16 {
		log.Panicf("MustParseI16: value %s out of range", s)
	}
	return int16(n)
}

// ParseU16 parse string to uint16
func ParseU16(s string) (uint16, error) {
	n, err := ParseI32(s)
	if err != nil {
		return 0, err
	}
	if n > math.MaxUint16 || n < 0 {
		return 0, fmt.Errorf("ParseU16: value %s out of range", s)
	}
	return uint16(n), nil
}

// MustParseU16 parse string to uint16, panic if error
func MustParseU16(s string) uint16 {
	n, err := ParseI32(s)
	if err != nil || n > math.MaxUint16 || n < 0 {
		log.Panicf("MustParseU16: value %s out of range", s)
		return 0
	}
	return uint16(n)
}

// ParseI32 parse string to int32
func ParseI32(s string) (int32, error) {
	if s == "" {
		return 0, nil
	}
	n, err := strconv.ParseInt(s, 10, 32)
	return int32(n), err
}

// MustParseI32 parse string to int32, panic if error
func MustParseI32(s string) int32 {
	n, err := ParseI32(s)
	if err != nil {
		log.Panicf("MustParseI32: cannot parse [%s] to int32: %v", s, err)
		return 0
	}
	return n
}

// ParseU32 parse string to uint32
func ParseU32(s string) (uint32, error) {
	if s == "" {
		return 0, nil
	}
	n, err := strconv.ParseUint(s, 10, 32)
	return uint32(n), err
}

// MustParseU32 parse string to uint32, panic if error
func MustParseU32(s string) uint32 {
	n, err := ParseU32(s)
	if err != nil {
		log.Panicf("MustParseU32: cannot parse [%s] to uint32: %v", s, err)
		return 0
	}
	return n
}

// ParseI64 parse string to int64
func ParseI64(s string) (int64, error) {
	if s == "" {
		return 0, nil
	}
	return strconv.ParseInt(s, 10, 64)
}

// MustParseI64 parse string to int64, panic if error
func MustParseI64(s string) int64 {
	n, err := ParseI64(s)
	if err != nil {
		log.Panicf("MustParseI64: cannot parse [%s] to uint64: %v", s, err)
		return 0
	}
	return n
}

// ParseU64 parse string to uint64
func ParseU64(s string) (uint64, error) {
	if s == "" {
		return 0, nil
	}
	return strconv.ParseUint(s, 10, 64)
}

// MustParseU64 parse string to uint64, panic if error
func MustParseU64(s string) uint64 {
	n, err := ParseU64(s)
	if err != nil {
		log.Panicf("MustParseU64: cannot parse [%s] to uint64: %v", s, err)
		return 0
	}
	return n
}

// ParseF32 parse string to float32
func ParseF32(s string) (float32, error) {
	if s == "" {
		return 0, nil
	}
	f, err := strconv.ParseFloat(s, 32)
	return float32(f), err
}

// MustParseF32 parse string to float32, panic if error
func MustParseF32(s string) float32 {
	f, err := ParseF32(s)
	if err != nil {
		log.Panicf("MustParseF32: cannot parse [%s] to float", s)
		return 0
	}
	return f
}

// ParseF64 parse string to float64
func ParseF64(s string) (float64, error) {
	if s == "" {
		return 0, nil
	}
	return strconv.ParseFloat(s, 64)
}

// MustParseF64 parse string to float64, panic if error
func MustParseF64(s string) float64 {
	f, err := ParseF64(s)
	if err != nil {
		log.Panicf("MustParseF64: cannot parse [%s] to double: %v", s, err)
		return 0
	}
	return f
}

// ParseTo parse string to any number type
// this routine may not have the best performance, but it has the best usability.
func ParseTo[T cmp.Ordered | bool](s string) (T, error) {
	var zero T
	if s == "" {
		return zero, nil
	}
	switch any(zero).(type) {
	case string:
		return any(s).(T), nil
	case bool:
		var b = ParseBool(s)
		return any(b).(T), nil
	case int8:
		if n, err := ParseI8(s); err != nil {
			return zero, err
		} else {
			return any(n).(T), nil
		}
	case uint8:
		if n, err := ParseU8(s); err != nil {
			return zero, err
		} else {
			return any(n).(T), nil
		}
	case int16:
		if n, err := ParseI16(s); err != nil {
			return zero, err
		} else {
			return any(n).(T), nil
		}
	case uint16:
		if n, err := ParseU16(s); err != nil {
			return zero, err
		} else {
			return any(n).(T), nil
		}
	case int32:
		if n, err := ParseI32(s); err != nil {
			return zero, err
		} else {
			return any(n).(T), nil
		}
	case uint32:
		if n, err := ParseU32(s); err != nil {
			return zero, err
		} else {
			return any(n).(T), nil
		}
	case int:
		if n, err := ParseI64(s); err != nil {
			return zero, err
		} else {
			return any(int(n)).(T), nil
		}
	case uint:
		if n, err := ParseU64(s); err != nil {
			return zero, err
		} else {
			return any(uint(n)).(T), nil
		}
	case int64:
		if n, err := ParseI64(s); err != nil {
			return zero, err
		} else {
			return any(n).(T), nil
		}
	case uint64:
		if n, err := ParseU64(s); err != nil {
			return zero, err
		} else {
			return any(n).(T), nil
		}
	case float32:
		if n, err := ParseF32(s); err != nil {
			return zero, err
		} else {
			return any(n).(T), nil
		}
	case float64:
		if n, err := ParseF64(s); err != nil {
			return zero, err
		} else {
			return any(n).(T), nil
		}
	}
	return zero, fmt.Errorf("ParseTo: unsupported type %T", zero)
}

func MustParseTo[T cmp.Ordered | bool](s string) T {
	val, err := ParseTo[T](s)
	if err != nil {
		slog.Panicf("cannot parse %s to %T: %v", s, val, err)
	}
	return val
}
