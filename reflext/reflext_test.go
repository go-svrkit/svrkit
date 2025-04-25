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
