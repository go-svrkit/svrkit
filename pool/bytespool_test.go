// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package pool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllocBytesSmallSize(t *testing.T) {
	var buf = AllocBytes(0)
	assert.Nil(t, buf)

	// small size
	buf = AllocBytes(6)
	assert.NotNil(t, buf)
	assert.Equal(t, 6, len(buf))
}

func TestAllocBytesLargeSize(t *testing.T) {
	// large size
	var buf = AllocBytes(32000)
	assert.NotNil(t, buf)
	assert.Equal(t, 32000, len(buf))
}

func TestAllocBytesFree(t *testing.T) {
	for i := 0; i < len(class_to_size); i++ {
		var size = class_to_size[i] + 1
		var buf = AllocBytes(int(size))
		assert.NotNil(t, buf)
		assert.Equal(t, int(size), len(buf))
		FreeBytes(buf)
	}
}
