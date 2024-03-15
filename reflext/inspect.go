// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package reflext

import (
	"reflect"
	"unsafe"
)

var allTypes = EnumerateAllTypes()

func EnumerateAllTypes() map[string]reflect.Type {
	var types = make(map[string]reflect.Type)
	var typ = reflect.TypeOf(0)
	var iface = (*EmptyInterface)(unsafe.Pointer(&typ))

	sections, offset := typelinks()
	for i, offs := range offset {
		rodata := sections[i]
		for _, off := range offs {
			iface.Word = resolveTypeOff(rodata, off)
			if typ.Kind() == reflect.Ptr && len(typ.Elem().Name()) > 0 {
				types[typ.String()] = typ
				types[typ.Elem().String()] = typ.Elem()
			}
		}
	}
	return types
}

// TypeForName find the type(exported and unexported) by package path and name.
func TypeForName(pathName string) reflect.Type {
	if typ, ok := allTypes[pathName]; ok {
		return typ
	}
	return nil
}
