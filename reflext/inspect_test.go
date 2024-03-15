package reflext

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeForName(t *testing.T) {
	typ := TypeForName("runtime.g")
	assert.NotNil(t, typ)

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		t.Logf("%s: %s", field.Name, field.Type.String())
	}
}
