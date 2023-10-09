// Copyright © 2022 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package reflext

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"gopkg.in/svrkit.v1/strutil"
)

// ParseCallExpr 把一个函数调用解析为函数名和参数
func ParseCallExpr(expr string) (fnName string, params []string, outErr error) {
	node, err := parser.ParseExpr(expr)
	if err != nil {
		outErr = err
		return
	}
	// 只能是调用表达式
	call, ok := node.(*ast.CallExpr)
	if !ok {
		outErr = fmt.Errorf("only CallExpr allowed")
		return
	}
	// 函数名称只能是identifier
	fn, ok := call.Fun.(*ast.Ident)
	if !ok {
		outErr = fmt.Errorf("command is not identifier")
		return
	}
	fnName = fn.Name
	params = make([]string, 0, len(call.Args))

	// 把常量放入参数列表
	var putLiteralArg = func(literal *ast.BasicLit) error {
		var param = literal.Value
		switch literal.Kind {
		case token.STRING, token.CHAR:
			if param, err = strconv.Unquote(param); err != nil {
				return fmt.Errorf("cannot unquote: %w", err)
			}
			param = strings.TrimSpace(param)
		}
		params = append(params, param)
		return nil
	}
	// 参数只能是一元运算表达式或者常量，不支持嵌套和eval
	for i, arg := range call.Args {
		switch v := arg.(type) {
		case *ast.UnaryExpr:
			// 一元运算表达式只能是常量
			literal, ok := v.X.(*ast.BasicLit)
			if !ok {
				outErr = fmt.Errorf("argument %d is not literal", i)
				return
			}
			// 一元运算符号只接受【+-】
			switch v.Op {
			case token.SUB:
				literal.Value = "-" + literal.Value // 这里修改了ast里值
			case token.ADD:
			default:
				outErr = fmt.Errorf("unrecognized argument expression %T", v)
				return
			}
			if outErr = putLiteralArg(literal); outErr != nil {
				return
			}

		case *ast.BasicLit:
			if outErr = putLiteralArg(v); outErr != nil {
				return
			}
		}
	}
	return
}

// CastParamToType 转换函数所有参数类型
func CastParamToType(rType reflect.Type, s string) (result reflect.Value, err error) {
	var kind = rType.Kind()
	if IsPrimitive(kind) {
		return ParseBaseKindToType(rType, s)
	}
	switch kind {
	case reflect.Ptr:
		return CastParamToType(rType.Elem(), s)
	case reflect.Slice, reflect.Map, reflect.Struct:
		var b = unsafe.Slice(unsafe.StringData(s), len(s))
		var rd = bytes.NewReader(b)
		var dec = json.NewDecoder(rd)
		dec.UseNumber()
		result = reflect.New(rType)
		err = dec.Decode(result.Interface())
		return
	default:
		err = fmt.Errorf("unrecognized type %s", rType.String())
		return
	}
}

// ParseInputArgs 解析传参
func ParseInputArgs(fnType reflect.Type, args []string) (result []reflect.Value, err error) {
	var numIn = fnType.NumIn()
	var isVariadic = fnType.IsVariadic() // 可变参数
	var isArgsMatch = len(args) == numIn
	if isVariadic {
		isArgsMatch = len(args) >= numIn-1 // 最后一个参数是variadic
	}
	if isArgsMatch {
		err = fmt.Errorf("method %v expect parameters not match", fnType.String())
		return
	}
	for i := 0; i < numIn; i++ {
		var paramType = fnType.In(i)
		// last variadic parameter
		if isVariadic && i+1 == numIn {
			paramType = paramType.Elem() // paramType is []T
			for j := i; j < len(args); j++ {
				var value reflect.Value
				if value, err = CastParamToType(paramType, args[i]); err != nil {
					return
				}
				result = append(result, value)
			}
			return
		} else {
			var value reflect.Value
			if value, err = CastParamToType(paramType, args[i]); err != nil {
				return
			}
			result = append(result, value)
		}
	}
	return
}

// ParseBaseKindToType 转换函数所有参数类型
func ParseBaseKindToType(rType reflect.Type, s string) (ret reflect.Value, err error) {
	switch rType.Kind() {
	case reflect.String:
		ret = reflect.ValueOf(s)

	case reflect.Int:
		var n int64
		if n, err = strutil.ParseI64(s); err == nil {
			ret = reflect.ValueOf(n)
		}

	case reflect.Int8:
		var n int8
		if n, err = strutil.ParseI8(s); err == nil {
			ret = reflect.ValueOf(n)
		}

	case reflect.Int16:
		var n int16
		if n, err = strutil.ParseI16(s); err == nil {
			ret = reflect.ValueOf(n)
		}

	case reflect.Int32:
		var n int32
		if n, err = strutil.ParseI32(s); err == nil {
			ret = reflect.ValueOf(n)
		}

	case reflect.Int64:
		var n int64
		if n, err = strutil.ParseI64(s); err == nil {
			ret = reflect.ValueOf(n)
		}

	case reflect.Uint:
		var n uint64
		if n, err = strconv.ParseUint(s, 10, 64); err == nil {
			ret = reflect.ValueOf(uint(n))
		}

	case reflect.Uint8:
		var n uint8
		if n, err = strutil.ParseU8(s); err == nil {
			ret = reflect.ValueOf(n)
		}

	case reflect.Uint16:
		var n uint16
		if n, err = strutil.ParseU16(s); err == nil {
			ret = reflect.ValueOf(n)
		}

	case reflect.Uint32:
		var n uint32
		if n, err = strutil.ParseU32(s); err == nil {
			ret = reflect.ValueOf(n)
		}

	case reflect.Uint64:
		var n uint64
		if n, err = strutil.ParseU64(s); err == nil {
			ret = reflect.ValueOf(n)
		}

	case reflect.Float32:
		var f float32
		if f, err = strutil.ParseF32(s); err == nil {
			ret = reflect.ValueOf(f)
		}

	case reflect.Float64:
		var f float64
		if f, err = strutil.ParseF64(s); err == nil {
			ret = reflect.ValueOf(f)
		}

	default:
		err = fmt.Errorf("unrecognized type %v", rType.Kind())
	}
	return
}
