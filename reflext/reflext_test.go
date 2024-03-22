package reflext

import (
	"bytes"
	"image"
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"gopkg.in/svrkit.v1/reflext/rt"
)

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

func TestEnumerateAllTypes(t *testing.T) {
	var allTypes = EnumerateAllTypes()
	typ := allTypes["runtime.g"] // runtime G
	assert.NotNil(t, typ)
	assert.Truef(t, typeHasField(typ, "goid", reflect.Uint64), "runtime.g should have field goid")
	assert.Truef(t, typeHasField(typ, "sched", reflect.Struct), "runtime.g should have field sched")

}

func TestGetFunc(t *testing.T) {
	var timeNowFunc func() (int64, int32)
	GetFunc(&timeNowFunc, "time.now")
	sec, nsec := timeNowFunc()
	if sec == 0 && nsec == 0 {
		t.Error("Expected nonzero result from time.now().")
	}
}

func TestRTPackEface(t *testing.T) {
	var pt = image.Point{X: 1234, Y: 5678}
	var eface = rt.UnpackEface(pt)
	assert.NotNil(t, eface.Typ)
	assert.Equal(t, eface.Typ.Size_, unsafe.Sizeof(pt))
	assert.Equal(t, eface.Typ.PtrBytes, uintptr(0))
	assert.NotNil(t, eface.Data)

	var val = rt.PackEface(&eface)
	ptt, ok := val.(image.Point)
	assert.True(t, ok)
	assert.Equal(t, pt, ptt)
	assert.Equal(t, &pt, &ptt)
}

func TestRTPackReflectType(t *testing.T) {
	var buf bytes.Buffer
	var rt1 = reflect.TypeOf(buf)
	var gotyp = rt.UnpackReflectType(rt1)
	assert.NotNil(t, gotyp)
	assert.Equal(t, gotyp.Size_, unsafe.Sizeof(buf))
	assert.Greater(t, int(gotyp.PtrBytes), 0)

	var rt2 = rt.PackReflectType(gotyp)
	assert.Equal(t, rt1.String(), rt2.String())
}

func TestRTMap(t *testing.T) {
	var mm = map[int]string{1234: "1234", 5678: "5678"}
	var gomap = *(**rt.GoMap)(unsafe.Pointer(&mm))
	assert.NotNil(t, gomap)
	assert.Equal(t, gomap.Count, len(mm))
}
