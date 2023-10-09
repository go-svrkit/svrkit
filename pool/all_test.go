// Copyright © 2020 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package pool

import (
	"testing"
)

func TestAllocBuffer(t *testing.T) {
	var b = AllocBuffer(12)
	println("len", len(b))
	println("cap", cap(b))
	FreeBuffer(b)
	b = AllocBuffer(16)
	println("len", len(b))
	println("cap", cap(b))
}
