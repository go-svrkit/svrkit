// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package reflext

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"unsafe"
)

func EnumerateAllTypes() map[string]reflect.Type {
	var types = make(map[string]reflect.Type)
	sections, offsets := typelinks()
	for i, offs := range offsets {
		rodata := sections[i]
		for _, base := range offs {
			var typeAddr = resolveTypeOff(rodata, base)
			typ := reflect.TypeOf(*(*interface{})(unsafe.Pointer(&typeAddr)))
			var kind = typ.Kind()
			for indirect := 0; indirect < 3 && kind == reflect.Ptr; indirect++ {
				typ = typ.Elem()
				kind = typ.Kind()
			}
			// we only care struct types
			if kind != reflect.Struct {
				continue
			}
			var name = typ.String()
			if !strings.HasPrefix(name, "struct ") { // skip unnamed struct
				types[typ.String()] = typ
			}
		}
	}
	return types
}

// GetFunc gets the function defined by the given fully-qualified name. The
// outFuncPtr parameter should be a pointer to a function with the appropriate
// type (e.g. the address of a local variable), and is set to a new function
// value that calls the specified function. If the specified function does not
// exist, outFuncPtr is not set and an error is returned.
func GetFunc(outFuncPtr interface{}, name string) error {
	codePtr, err := FindFuncWithName(name)
	if err != nil {
		return err
	}
	CreateFuncForCodePtr(outFuncPtr, codePtr)
	return nil
}

// FuncP Convenience struct for modifying the underlying code pointer of a function
// value. The actual struct has other values, but always starts with a code
// pointer.
type FuncP struct {
	pc uintptr
}

// CreateFuncForCodePtr is given a code pointer and creates a function value
// that uses that pointer. The outFun argument should be a pointer to a function
// of the proper type (e.g. the address of a local variable), and will be set to
// the result function value.
func CreateFuncForCodePtr(outFuncPtr interface{}, pc uintptr) {
	outFuncVal := reflect.ValueOf(outFuncPtr).Elem()
	newFuncVal := reflect.MakeFunc(outFuncVal.Type(), nil)
	funcValuePtr := reflect.ValueOf(newFuncVal).FieldByName("ptr").Pointer()
	funcPtr := (*FuncP)(unsafe.Pointer(funcValuePtr))
	funcPtr.pc = pc
	outFuncVal.Set(newFuncVal)
}

// FindFuncWithName searches through the moduledata table created by the linker
// and returns the function's code pointer. If the function was not found, it
// returns an error.
func FindFuncWithName(name string) (uintptr, error) {
	for md := &firstmoduledata; md != nil; md = md.next {
		for _, ftab := range md.ftab {
			f := (*runtime.Func)(unsafe.Pointer(&md.pclntable[ftab.funcoff]))
			if f.Name() == name {
				return f.Entry(), nil
			}
		}
	}
	return 0, fmt.Errorf("invalid function name: %s", name)
}
