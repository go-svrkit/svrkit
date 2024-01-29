// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package pool

import (
	"testing"
)

func TestAllocBuffer(t *testing.T) {
	var b = AllocBytes(12)
	println("len", len(b))
	println("cap", cap(b))
	FreeBytes(b)
	b = AllocBytes(16)
	println("len", len(b))
	println("cap", cap(b))
}
