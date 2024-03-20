package reflext

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
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

	// goroutine ID
	var offset = getFieldOffset(typ, "goid")
	var base = getg()
	var goid = *(*int64)(unsafe.Add(unsafe.Pointer(base), offset))
	t.Logf("goroutine ID: %d", goid)
}

func TestGetFunc(t *testing.T) {
	var ptr = getg()
	t.Logf("%d", ptr)
	var timeNowFunc func() (int64, int32)
	GetFunc(&timeNowFunc, "time.now")
	sec, nsec := timeNowFunc()
	if sec == 0 && nsec == 0 {
		t.Error("Expected nonzero result from time.now().")
	}
}