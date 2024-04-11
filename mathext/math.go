// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package mathext

import (
	"math"
	"math/bits"
)

func Max[T Integer](x, y T) T {
	if x < y {
		return y
	}
	return x
}

func Min[T Integer](x, y T) T {
	if x > y {
		return y
	}
	return x
}

func Abs[T Integer](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

func Dim[T Integer](x, y T) T {
	return Max[T](x-y, 0)
}

func SafeDiv[T Number](x, y T) T {
	var zero T
	if y == zero {
		return zero
	}
	return x / y
}

func SafeMulUint64(a, b uint64) (product uint64, overflow bool) {
	hi, lo := bits.Mul64(a, b)
	if hi != 0 {
		return 0, true
	}
	return lo, false
}

func SafeMulInt64(a, b int64) (product int64, overflow bool) {
	var sign int64 = 1
	if a < 0 {
		a = -a
		sign = -sign
	}
	if b < 0 {
		b = -b
		sign = -sign
	}
	hi, lo := bits.Mul64(uint64(a), uint64(b))
	if hi != 0 || lo > math.MaxInt64 {
		return 0, true
	}
	if sign < 0 {
		return -int64(lo), false
	}
	return int64(lo), false
}
