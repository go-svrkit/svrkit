// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package reflext

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"unsafe"

	"gopkg.in/svrkit.v1/reflext/rt"
)

func EnumerateAllStructs() map[string]reflect.Type {
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

func EnumerateAllFuncPCs() map[string]uintptr {
	var all = map[string]uintptr{}
	for md := rt.GetFirstModuleData(); md != nil; md = md.Next {
		for _, ftab := range md.Ftab {
			if int(ftab.Funcoff) >= len(md.Pclntable) {
				continue
			}
			f := (*runtime.Func)(unsafe.Pointer(&md.Pclntable[ftab.Funcoff]))
			if f != nil {
				if name := f.Name(); name != "" {
					all[name] = f.Entry()
				}
			}
		}
	}
	return all
}

func GetFunc(outFuncPtr interface{}, name string) error {
	pc, err := FindFuncPCWithName(name)
	if err != nil {
		return err
	}
	CreateFuncForPC(outFuncPtr, pc)
	return nil
}

type FuncP struct {
	pc uintptr
}

func CreateFuncForPC(outFuncPtr interface{}, pc uintptr) {
	outFuncVal := reflect.ValueOf(outFuncPtr).Elem()
	newFuncVal := reflect.MakeFunc(outFuncVal.Type(), nil)
	funcValuePtr := reflect.ValueOf(newFuncVal).FieldByName("ptr").Pointer()
	funcPtr := (*FuncP)(unsafe.Pointer(funcValuePtr))
	funcPtr.pc = pc
	outFuncVal.Set(newFuncVal)
}

func FindFuncPCWithName(name string) (uintptr, error) {
	for md := rt.GetFirstModuleData(); md != nil; md = md.Next {
		for _, ftab := range md.Ftab {
			f := (*runtime.Func)(unsafe.Pointer(&md.Pclntable[ftab.Funcoff]))
			if f != nil && f.Name() == name {
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
