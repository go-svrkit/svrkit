package reflext

import (
	"fmt"
	"image"
	"reflect"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetFieldByName(t *testing.T) {
	var pt = &image.Point{X: 1, Y: 2}
	var iface fmt.Stringer = pt
	var err = SetFieldByName(iface, "X", 3)
	assert.Nil(t, err, "SetFieldByName failed")
	assert.Equal(t, 3, pt.X)
	assert.Equal(t, 2, pt.Y)
}

func TestGetStructFieldNames(t *testing.T) {
	var pt image.Point
	var names = GetStructFieldNames(reflect.TypeOf(pt))
	assert.Equal(t, 2, len(names))
	slices.Sort(names)
	assert.Equal(t, "X", names[0])
	assert.Equal(t, "Y", names[1])
}

func TestGetStructFieldValues(t *testing.T) {
	var pt = image.Point{X: 11, Y: 22}
	var values = GetStructFieldValues(reflect.ValueOf(pt))
	assert.Equal(t, 2, len(values))
	assert.Equal(t, 11, values[0].(int))
	assert.Equal(t, 22, values[1].(int))
}

func getFieldOffset(typ reflect.Type, fieldName string) uintptr {
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Name == fieldName {
			return field.Offset
		}
	}
	return 0
}

func typeHasField(typ reflect.Type, fieldName string, kind reflect.Kind) bool {
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Name == fieldName && field.Type.Kind() == kind {
			return true
		}
	}
	return false
}

func TestEnumerateAllStructs(t *testing.T) {
	var allTypes = EnumerateAllStructs()
	typ := allTypes["runtime.g"] // runtime G
	assert.NotNil(t, typ)
	assert.Truef(t, typeHasField(typ, "goid", reflect.Uint64), "runtime.g should have field goid")
	assert.Truef(t, typeHasField(typ, "sched", reflect.Struct), "runtime.g should have field sched")

}
