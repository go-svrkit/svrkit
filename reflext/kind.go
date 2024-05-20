// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package reflext

import (
	"reflect"
)

func IsInteger(kind reflect.Kind) bool {
	return IsSignedInteger(kind) || IsUnsignedInteger(kind)
}

func IsSignedInteger(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int, reflect.Int64:
		return true
	default:
		return false
	}
}

func IsUnsignedInteger(kind reflect.Kind) bool {
	switch kind {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint, reflect.Uint64:
		return true
	default:
		return false
	}
}

func IsFloat(kind reflect.Kind) bool {
	switch kind {
	case reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func IsNumber(kind reflect.Kind) bool {
	return IsInteger(kind) || IsFloat(kind)
}

// IsPrimitive 是否基本类型
func IsPrimitive(kind reflect.Kind) bool {
	switch kind {
	case reflect.Bool, reflect.String:
		return true
	default:
		return IsNumber(kind)
	}
}

// IsInterfaceNil interface是否是nil
func IsInterfaceNil(c any) bool {
	if c == nil {
		return true
	}
	var rv = reflect.ValueOf(c)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Array, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return rv.IsNil()
	default:
		return false
	}
}
