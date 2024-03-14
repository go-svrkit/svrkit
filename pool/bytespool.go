// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package pool

import (
	"sync"
)

var bytesPool = createLeveledBytesPool()

func AllocBytes(size int) []byte {
	if size == 0 {
		return nil
	}
	if size <= MaxSmallSize {
		var sizeClass = getSizeClass(size)
		if sizeClass > 0 {
			return bytesPool[sizeClass-1].Get().([]byte)
		}
	}
	return make([]byte, size)
}

func FreeBytes(buf []byte) {
	var size = len(buf)
	if cap(buf) >= MaxSmallSize && len(buf) <= MaxSmallSize {
		var sizeClass = getSizeClass(size)
		if sizeClass > 0 {
			bytesPool[sizeClass-1].Put(buf)
		}
	}
}

func createLeveledBytesPool() *[NumSizeClasses - 1]sync.Pool {
	var pools [NumSizeClasses - 1]sync.Pool
	for i := 0; i < NumSizeClasses-1; i++ {
		size := int(class_to_size[i+1])
		pools[i].New = func() interface{} {
			return make([]byte, size)
		}
	}
	return &pools
}

func getSizeClass(size int) uint8 {
	var sizeclass uint8
	if size <= smallSizeMax-8 {
		sizeclass = size_to_class8[divRoundUp(size, smallSizeDiv)]
	} else {
		sizeclass = size_to_class128[divRoundUp(size-smallSizeMax, largeSizeDiv)]
	}
	return sizeclass
}

// divRoundUp returns ceil(n / a).
func divRoundUp(n, a int) int {
	// a is generally a power of two. This will get inlined and
	// the compiler will optimize the division.
	return (n + a - 1) / a
}

// see $GOROOT/src/runtime/sizeclasses.go
const (
	MaxSmallSize   = 32768
	smallSizeDiv   = 8
	smallSizeMax   = 1024
	largeSizeDiv   = 128
	NumSizeClasses = 68
)

var class_to_size = [NumSizeClasses]uint16{0, 8, 16, 24, 32, 48, 64, 80, 96, 112, 128, 144, 160, 176, 192, 208, 224, 240, 256, 288, 320, 352, 384, 416, 448, 480, 512, 576, 640, 704, 768, 896, 1024, 1152, 1280, 1408, 1536, 1792, 2048, 2304, 2688, 3072, 3200, 3456, 4096, 4864, 5376, 6144, 6528, 6784, 6912, 8192, 9472, 9728, 10240, 10880, 12288, 13568, 14336, 16384, 18432, 19072, 20480, 21760, 24576, 27264, 28672, 32768}
var size_to_class8 = [smallSizeMax/smallSizeDiv + 1]uint8{0, 1, 2, 3, 4, 5, 5, 6, 6, 7, 7, 8, 8, 9, 9, 10, 10, 11, 11, 12, 12, 13, 13, 14, 14, 15, 15, 16, 16, 17, 17, 18, 18, 19, 19, 19, 19, 20, 20, 20, 20, 21, 21, 21, 21, 22, 22, 22, 22, 23, 23, 23, 23, 24, 24, 24, 24, 25, 25, 25, 25, 26, 26, 26, 26, 27, 27, 27, 27, 27, 27, 27, 27, 28, 28, 28, 28, 28, 28, 28, 28, 29, 29, 29, 29, 29, 29, 29, 29, 30, 30, 30, 30, 30, 30, 30, 30, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32}
var size_to_class128 = [(MaxSmallSize-smallSizeMax)/largeSizeDiv + 1]uint8{32, 33, 34, 35, 36, 37, 37, 38, 38, 39, 39, 40, 40, 40, 41, 41, 41, 42, 43, 43, 44, 44, 44, 44, 44, 45, 45, 45, 45, 45, 45, 46, 46, 46, 46, 47, 47, 47, 47, 47, 47, 48, 48, 48, 49, 49, 50, 51, 51, 51, 51, 51, 51, 51, 51, 51, 51, 52, 52, 52, 52, 52, 52, 52, 52, 52, 52, 53, 53, 54, 54, 54, 54, 55, 55, 55, 55, 55, 56, 56, 56, 56, 56, 56, 56, 56, 56, 56, 56, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 58, 58, 58, 58, 58, 58, 59, 59, 59, 59, 59, 59, 59, 59, 59, 59, 59, 59, 59, 59, 59, 59, 60, 60, 60, 60, 60, 60, 60, 60, 60, 60, 60, 60, 60, 60, 60, 60, 61, 61, 61, 61, 61, 62, 62, 62, 62, 62, 62, 62, 62, 62, 62, 62, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67, 67}
