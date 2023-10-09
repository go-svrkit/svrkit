// Copyright © 2021 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package pool

import (
	"log"
	"math/bits"
	"sync"
)

// 固定大小的pool，code taken from github.com/vitessio/vitess bucketpool with modification
type sizedPool struct {
	size int
	pool sync.Pool
}

func newSizedPool(size int) *sizedPool {
	return &sizedPool{
		size: size,
		pool: sync.Pool{
			New: func() interface{} { return makeSizedBuffer(size) },
		},
	}
}

func makeSizedBuffer(size int) []byte {
	return make([]byte, size)
}

func makeSizedPools(minSize, maxSize int) []*sizedPool {
	const multiplier = 2 // double the size
	var pools []*sizedPool
	curSize := minSize
	for curSize < maxSize {
		pools = append(pools, newSizedPool(curSize))
		curSize *= multiplier
	}
	pools = append(pools, newSizedPool(maxSize))
	return pools
}

// BucketPool is actually multiple pools which store buffers of specific size.
type BucketPool struct {
	minSize int
	maxSize int
	pools   []*sizedPool
}

func NewBucketPool(minSize, maxSize int) *BucketPool {
	if maxSize < minSize {
		panic("maxSize less than minSize")
	}
	return &BucketPool{
		minSize: minSize,
		maxSize: maxSize,
		pools:   makeSizedPools(minSize, maxSize),
	}
}

func (p *BucketPool) findPool(size int) *sizedPool {
	if size > p.maxSize {
		return nil
	}
	div, rem := bits.Div64(0, uint64(size), uint64(p.minSize))
	idx := bits.Len64(div)
	if rem == 0 && div != 0 && (div&(div-1)) == 0 {
		idx = idx - 1
	}
	return p.pools[idx]
}

// Get returns a []byte slice which has len size.
// If there is no bucket with buffers >= size, slice will be allocated.
func (p *BucketPool) Get(size int) []byte {
	sp := p.findPool(size)
	if sp == nil {
		return makeSizedBuffer(size)
	}
	buf := sp.pool.Get().([]byte)
	return buf[:size]
}

// Put returns a slice to some bucket. Discards slice for which there is no bucket
func (p *BucketPool) Put(b []byte) {
	sp := p.findPool(cap(b))
	if sp == nil {
		return
	}
	if sp.size == cap(b) {
		sp.pool.Put(b)
	} else {
		log.Panicf("unexpected pool buffer, cap %d", cap(b))
	}
}
