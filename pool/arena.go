// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package pool

import (
	"sync"
)

// ArenaPool
// 一次申请一个block（N个元素的数组），然后从block数组里再逐个按需分配，
// block分配完了就丢掉（交给GC)，再申请另一个block；
// 这样对runtime来说每次malloc都是以N个元素大小的单位，可以减缓GC的压力
type ArenaPool[T any] struct {
	guard sync.Mutex
	idx   int
	block []T
}

func NewArenaPool[T any](blockSize int) *ArenaPool[T] {
	return &ArenaPool[T]{
		block: make([]T, blockSize),
	}
}

func (a *ArenaPool[T]) Alloc() *T {
	a.guard.Lock()
	var size = len(a.block)
	var ret = &a.block[a.idx]
	a.idx++
	if a.idx >= size {
		a.block = make([]T, size)
		a.idx = 0
	}
	a.guard.Unlock()
	return ret
}

func (a *ArenaPool[T]) Free(*T) {
	// do nothing
}
