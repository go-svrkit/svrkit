// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package reflext

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"unsafe"

	"gopkg.in/svrkit.v1/reflext/rt"
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

func GetFunc(outFuncPtr interface{}, name string) error {
	codePtr, err := FindFuncWithName(name)
	if err != nil {
		return err
	}
	CreateFuncForCodePtr(outFuncPtr, codePtr)
	return nil
}

type FuncP struct {
	pc uintptr
}

func CreateFuncForCodePtr(outFuncPtr interface{}, pc uintptr) {
	outFuncVal := reflect.ValueOf(outFuncPtr).Elem()
	newFuncVal := reflect.MakeFunc(outFuncVal.Type(), nil) // reflect.Value
	//(*rt.RValue)(unsafe.Pointer(&newFuncVal)).Ptr = unsafe.Pointer(pc)
	funcValuePtr := reflect.ValueOf(newFuncVal).FieldByName("ptr").Pointer() // reflect.Value.ptr
	funcPtr := (*FuncP)(unsafe.Pointer(funcValuePtr))                        // reflect.Value.ptr = pc
	funcPtr.pc = pc
	outFuncVal.Set(newFuncVal)
}

func FindFuncWithName(name string) (uintptr, error) {
	for md := rt.GetFirstModuleData(); md != nil; md = md.Next {
		for _, ftab := range md.Ftab {
			f := (*runtime.Func)(unsafe.Pointer(&md.Pclntable[ftab.Funcoff]))
			if f.Name() == name {
				return f.Entry(), nil
			}
		}
	}
	return 0, fmt.Errorf("invalid function name: %s", name)
}

//go:linkname typelinks reflect.typelinks
func typelinks() (sections []unsafe.Pointer, offset [][]int32)

//go:linkname resolveTypeOff reflect.resolveTypeOff
func resolveTypeOff(rtype unsafe.Pointer, off int32) unsafe.Pointer
