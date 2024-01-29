// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package pool

const (
	MaxBufSize = 32 << 10 // 32K
	MinBufSize = 8
)

var (
	bucketPool  = NewBucketPool(MinBufSize, MaxBufSize)
	emptyBuffer = make([]byte, 0)
)

func AllocBytes(size int) []byte {
	switch {
	case size <= 0:
		return emptyBuffer
	case size < MaxBufSize:
		return bucketPool.Get(size)
	default:
		return make([]byte, size)
	}
}

func FreeBytes(buf []byte) {
	if cap(buf) > 0 && cap(buf) < MaxBufSize {
		bucketPool.Put(buf)
	}
}
