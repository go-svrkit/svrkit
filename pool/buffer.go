// Copyright © 2021 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package pool

const (
	MaxBufSize = 32 << 10 // 32K
	MinBufSize = 8
)

var (
	defaultPool = NewBucketPool(MinBufSize, MaxBufSize)
	emptyBuffer = make([]byte, 0)
)

func AllocBuffer(size int) []byte {
	switch {
	case size <= 0:
		return emptyBuffer
	case size < MaxBufSize:
		return defaultPool.Get(size)
	default:
		return make([]byte, size)
	}
}

func FreeBuffer(buf []byte) {
	if cap(buf) > 0 && cap(buf) < MaxBufSize {
		defaultPool.Put(buf)
	}
}
