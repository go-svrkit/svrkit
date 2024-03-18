package reflext

import (
	"testing"
)

/*
	func TestRangeAllTypes(t *testing.T) {
		var count = 0
		RangeAllTypes(func(typ reflect.Type) bool {
			//t.Logf("%s", typ.String())
			count++
			return true
		})
		assert.True(t, count > 0)
		t.Logf("all types count: %d", count)
	}

	func TestRangePackageTypes(t *testing.T) {
		var count = 0
		RangePackageTypes("sync", func(typ reflect.Type) bool {
			//t.Logf("%s", typ.String())
			count++
			return true
		})
		assert.True(t, count > 0)
		t.Logf("all sync pacakge types count: %d", count)
	}

	func TestTypeForName(t *testing.T) {
		// runtime package
		{
			typ := TypeForName("runtime.g")
			assert.NotNil(t, typ)
			assert.Truef(t, typeHasField(typ, "goid", reflect.Uint64), "runtime.g should have field goid")
			assert.Truef(t, typeHasField(typ, "sched", reflect.Struct), "runtime.g should have field goid")
		}
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
*/
func TestGetFunc(t *testing.T) {
	var timeNowFunc func() (int64, int32)
	GetFunc(&timeNowFunc, "time.now")
	sec, nsec := timeNowFunc()
	if sec == 0 && nsec == 0 {
		t.Error("Expected nonzero result from time.now().")
	}
}
