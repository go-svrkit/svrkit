// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package pool

import (
	"sync"
	"unsafe"
)

// ArenaPageSize is _PageSize at runtime/malloc.go
const ArenaPageSize = 8192

// ArenaPool
// 一次申请一个block（N个元素的数组），然后从block数组里再逐个按需分配，试图减缓GC的压力
type ArenaPool[T any] struct {
	guard sync.Mutex
	off   int
	block []T
}

func NewArenaPool[T any]() *ArenaPool[T] {
	var dummy T
	var blockSize = ArenaPageSize / unsafe.Sizeof(dummy)
	if blockSize <= 1 {
		blockSize = 4 // no more than 32KB
	}
	return &ArenaPool[T]{
		block: make([]T, blockSize),
	}
}

func NewArenaPoolWith[T any](blockSize int) *ArenaPool[T] {
	return &ArenaPool[T]{
		block: make([]T, blockSize),
	}
}

// Size block的大小
func (a *ArenaPool[T]) Size() int {
	return len(a.block)
}

// Len 已分配的元素个数
func (a *ArenaPool[T]) Len() int {
	return a.off
}

// Cap 剩余可分配空间
func (a *ArenaPool[T]) Cap() int {
	return len(a.block) - a.off
}

// Ptr 返回block的指针
func (a *ArenaPool[T]) Ptr() unsafe.Pointer {
	return unsafe.Pointer(&a.block[0])
}

// Alloc 从block里分配一个
func (a *ArenaPool[T]) Alloc() *T {
	a.guard.Lock()
	defer a.guard.Unlock()

	var size = len(a.block)
	if size-a.off == 0 {
		a.block = make([]T, size)
		a.off = 0
	}
	var obj = &a.block[a.off]
	a.off++
	return obj
}

// AllocN 从block里分配n个
func (a *ArenaPool[T]) AllocN(n int) []T {
	a.guard.Lock()
	defer a.guard.Unlock()

	var size = len(a.block)
	if size < n {
		return nil // too big
	}
	if size-a.off < n {
		a.block = make([]T, size)
		a.off = 0
	}
	var obj = a.block[a.off : a.off+n]
	a.off += n
	return obj
}

func (a *ArenaPool[T]) Free(*T) {
	// do nothing
}
