// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package reflext

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

func indirect(v reflect.Value) reflect.Value {
	for {
		var kind = v.Kind()
		if kind == reflect.Interface || kind == reflect.Ptr {
			v = v.Elem()
		} else {
			break
		}
	}
	return v
}

// SetFieldByName 设置struct的field
func SetFieldByName(obj any, fileName string, fieldVal any) error {
	var rval = indirect(reflect.ValueOf(obj))
	if rval.Kind() != reflect.Struct {
		return fmt.Errorf("SetFieldByName: obj must be a struct")
	}
	var field = rval.FieldByName(fileName)
	if !field.IsValid() || !field.CanSet() {
		return fmt.Errorf("SetFieldByName: cannot set field %s", fileName)
	}
	field.Set(reflect.ValueOf(fieldVal))
	return nil
}

// GetStructFieldNames 获取struct的所有field名
func GetStructFieldNames(rType reflect.Type) []string {
	if rType.Kind() != reflect.Struct {
		return nil
	}
	var names = make([]string, 0, rType.NumField())
	for i := 0; i < rType.NumField(); i++ {
		names = append(names, rType.Field(i).Name)
	}
	return names
}

// GetStructFieldValues 获取struct的所有field值
func GetStructFieldValues(val reflect.Value) []any {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil
	}
	var numField = val.NumField()
	var values = make([]any, 0, numField)
	for i := 0; i < numField; i++ {
		var field = val.Field(i)
		values = append(values, field.Interface())
	}
	return values
}

// GetStructFieldValueMap 获取struct的所有field值和value
func GetStructFieldValueMap(val reflect.Value) map[string]any {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil
	}
	var rType = val.Type()
	var mm = map[string]any{}
	var numField = val.NumField()
	for i := 0; i < numField; i++ {
		var field = val.Field(i)
		var name = rType.Field(i).Name
		mm[name] = field.Interface()
	}
	return mm
}

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
