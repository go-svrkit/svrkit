// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package pool

import (
	"testing"
)

func TestAllocBytes(t *testing.T) {
	var buf = AllocBytes(0)
	FreeBytes(buf)

	for i := 0; i < len(class_to_size); i++ {
		var size = class_to_size[i] + 1
		buf = AllocBytes(int(size))
		FreeBytes(buf)
	}
}
