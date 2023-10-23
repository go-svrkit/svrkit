// Copyright © 2022 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package reflext

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"gopkg.in/svrkit.v1/strutil"
)

// ParseCallExprArgs 解析call表达式的参数列表
func ParseCallExprArgs(call *ast.CallExpr) ([]string, error) {
	if len(call.Args) == 0 {
		return nil, nil
	}
	var args = make([]string, 0, len(call.Args))
	// 把常量放入参数列表
	var parseLiteral = func(literal *ast.BasicLit) error {
		var param = literal.Value
		switch literal.Kind {
		case token.STRING, token.CHAR:
			var err error
			if param, err = strconv.Unquote(param); err != nil {
				return err
			}
			param = strings.TrimSpace(param)
		}
		args = append(args, param)
		return nil
	}
	// 参数只能是一元运算表达式或者常量，不支持嵌套和名字解析
	for _, arg := range call.Args {
		switch v := arg.(type) {
		case *ast.UnaryExpr:
			// 一元运算表达式只能是常量
			literal, ok := v.X.(*ast.BasicLit)
			if !ok {
				return nil, fmt.Errorf("ParseCallExprArgs: invalid unary expr")
			}
			// 一元运算符号只接受【+-】
			switch v.Op {
			case token.SUB:
				literal.Value = "-" + literal.Value // 这里修改了ast里值
			case token.ADD:
			default:
				return nil, fmt.Errorf("ParseCallExprArgs: invalid unary expr")
			}
			if er := parseLiteral(literal); er != nil {
				return nil, er
			}

		case *ast.BasicLit:
			if er := parseLiteral(v); er != nil {
				return nil, er
			}
		}
	}
	return args, nil
}

// ConvParamToType 转换函数所有参数类型
func ConvParamToType(rType reflect.Type, input string) (val reflect.Value, err error) {
	var kind = rType.Kind()
	if IsPrimitive(kind) {
		return ParseBaseKindToType(rType, input)
	}
	switch kind {
	case reflect.Ptr:
		return ConvParamToType(rType.Elem(), input)
	case reflect.Slice, reflect.Map, reflect.Struct:
		var b = unsafe.Slice(unsafe.StringData(input), len(input))
		var rd = bytes.NewReader(b)
		var dec = json.NewDecoder(rd)
		dec.UseNumber()
		val = reflect.New(rType)
		err = dec.Decode(val.Interface())
		return
	default:
		err = fmt.Errorf("ConvParamToType: unrecognized type %s", rType.String())
		return
	}
}

// ParseInputArgs 解析传参
func ParseInputArgs(fnType reflect.Type, args []string) (result []reflect.Value, err error) {
	if len(args) == 0 {
		return
	}
	var numIn = fnType.NumIn()
	var isVariadic = fnType.IsVariadic() // 可变参数
	var isArgsMatch = len(args) == numIn
	if isVariadic {
		isArgsMatch = len(args) >= numIn-1 // 最后一个参数是variadic
	}
	if !isArgsMatch {
		err = fmt.Errorf("method %v parameters and arguments not match", fnType.String())
		return
	}
	for i := 0; i < numIn; i++ {
		var paramType = fnType.In(i)
		// last variadic parameter
		if isVariadic && i+1 == numIn {
			paramType = paramType.Elem() // paramType is []T
			for j := i; j < len(args); j++ {
				var value reflect.Value
				if value, err = ConvParamToType(paramType, args[i]); err != nil {
					return
				}
				result = append(result, value)
			}
			return
		} else {
			var value reflect.Value
			if value, err = ConvParamToType(paramType, args[i]); err != nil {
				return
			}
			result = append(result, value)
		}
	}
	return
}

// ParseBaseKindToType 转换函数所有参数类型
func ParseBaseKindToType(rType reflect.Type, input string) (val reflect.Value, err error) {
	val = reflect.New(rType).Elem()
	switch rType.Kind() {
	case reflect.String:
		val.SetString(input)

	case reflect.Int:
		var n int64
		if n, err = strutil.ParseI64(input); err == nil {
			val.SetInt(n)
		}

	case reflect.Int8:
		var n int8
		if n, err = strutil.ParseI8(input); err == nil {
			val.SetInt(int64(n))
		}

	case reflect.Int16:
		var n int16
		if n, err = strutil.ParseI16(input); err == nil {
			val.SetInt(int64(n))
		}

	case reflect.Int32:
		var n int32
		if n, err = strutil.ParseI32(input); err == nil {
			val.SetInt(int64(n))
		}

	case reflect.Int64:
		var n int64
		if n, err = strutil.ParseI64(input); err == nil {
			val.SetInt(int64(n))
		}

	case reflect.Uint:
		var n uint64
		if n, err = strconv.ParseUint(input, 10, 64); err == nil {
			val.SetInt(int64(n))
		}

	case reflect.Uint8:
		var n uint8
		if n, err = strutil.ParseU8(input); err == nil {
			val.SetInt(int64(n))
		}

	case reflect.Uint16:
		var n uint16
		if n, err = strutil.ParseU16(input); err == nil {
			val.SetInt(int64(n))
		}

	case reflect.Uint32:
		var n uint32
		if n, err = strutil.ParseU32(input); err == nil {
			val.SetInt(int64(n))
		}

	case reflect.Uint64:
		var n uint64
		if n, err = strutil.ParseU64(input); err == nil {
			val.SetInt(int64(n))
		}

	case reflect.Float32:
		var f float32
		if f, err = strutil.ParseF32(input); err == nil {
			val.SetFloat(float64(f))
		}

	case reflect.Float64:
		var f float64
		if f, err = strutil.ParseF64(input); err == nil {
			val.SetFloat(f)
		}

	default:
		err = fmt.Errorf("ParseBaseKindToType: unrecognized type %v", rType.Kind())
	}
	return
}
