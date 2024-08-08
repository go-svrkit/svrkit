package algext

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZeroOf(t *testing.T) {
	assert.Equal(t, false, ZeroOf[bool]())
	assert.Equal(t, "", ZeroOf[string]())
	assert.Equal(t, 0, ZeroOf[int]())
	assert.Equal(t, uint(0), ZeroOf[uint]())
	assert.Equal(t, int8(0), ZeroOf[int8]())
	assert.Equal(t, uint8(0), ZeroOf[uint8]())
	assert.Equal(t, int16(0), ZeroOf[int16]())
	assert.Equal(t, uint16(0), ZeroOf[uint16]())
	assert.Equal(t, int32(0), ZeroOf[int32]())
	assert.Equal(t, uint32(0), ZeroOf[uint32]())
	assert.Equal(t, float32(0), ZeroOf[float32]())
	assert.Equal(t, float64(0), ZeroOf[float64]())
	assert.Equal(t, complex64(0), ZeroOf[complex64]())
	assert.Equal(t, image.Point{}, ZeroOf[image.Point]())
	assert.Equal(t, []int(nil), ZeroOf[[]int]())
	assert.Equal(t, map[int]int(nil), ZeroOf[map[int]int]())
}
