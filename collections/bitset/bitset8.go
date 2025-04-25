// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package bitset

import (
	"math/bits"
	"strings"
)

const BitsPerByte = 8 // 用uint8表示

func BitBytesCount(bitSize int) int {
	return (bitSize + BitsPerByte - 1) / BitsPerByte
}

type BitMap8 []uint8

func TestBytesBit(bm []uint8, i int) bool {
	if i >= 0 && i < len(bm)*BitsPerByte {
		return bm[i/BitsPerByte]&(1<<(i%BitsPerByte)) != 0
	}
	return false
}

func SetBytesBit(bm []uint8, i int) {
	if i >= 0 && i < len(bm)*BitsPerByte {
		var v = uint8(1) << (i % BitsPerByte)
		bm[i/BitsPerByte] |= v
	}
}

// MustSetBytesBit 指定位是否为1，并且自动增长数组
func MustSetBytesBit(bm []uint8, i int) BitMap8 {
	var n = BitBytesCount(i + 1)
	if n > len(bm) {
		var newb = make(BitMap8, n)
		copy(newb, bm)
		bm = newb
	}
	SetBytesBit(bm, i)
	return bm
}

// ClearBytesBit clears the bit at the given index.
func ClearBytesBit(bm []uint8, i int) {
	if i >= 0 && i < len(bm)*BitsPerByte {
		var v = uint8(1) << (i % BitsPerByte)
		bm[i/BitsPerByte] &= ^v
	}
}

// FlipBytesBit flips the bit at the given index.
func FlipBytesBit(bm []uint8, i int) {
	if i >= 0 && i < len(bm)*BitsPerByte {
		bm[i/BitsPerByte] ^= 1 << (i % BitsPerByte)
	}
}

// TestAndSetBytesBit returns the old Value of the bit at the given index and sets it to 1.
func TestAndSetBytesBit(bm []uint8, i int) bool {
	var v = uint8(1) << (i % BitsPerByte)
	var index = i / BitsPerByte
	if index >= 0 && index < len(bm) {
		var old = bm[index]
		bm[index] |= v
		return old&v != 0
	}
	return false
}

// TestAndClearBytesBit returns the old Value of the bit at the given index and clears it.
func TestAndClearBytesBit(bm []uint8, i int) bool {
	var v = uint8(1) << (i % BitsPerByte)
	var index = i / BitsPerByte
	if index >= 0 && index < len(bm) {
		var old = bm[index]
		bm[index] &= ^v
		return old&v != 0
	}
	return false
}

// OnesCountBytes returns the number of bits set to 1.
func OnesCountBytes(bm []uint8) int {
	var count int
	for i := 0; i < len(bm); i++ {
		if bm[i] != 0 {
			count += bits.OnesCount8(bm[i])
		}
	}
	return count
}

// IsZero 是否所有位都是0
func IsAllBytesZero(bm []uint8) bool {
	for i := 0; i < len(bm); i++ {
		if bm[i] != 0 {
			return false
		}
	}
	return true
}

func BytesToString(bm []uint8) string {
	var size = len(bm) * BitsPerByte
	var sb strings.Builder
	sb.Grow(size)
	for i := 0; i < size; i++ {
		if TestBytesBit(bm, i) {
			sb.WriteByte('1')
		} else {
			sb.WriteByte('0')
		}
	}
	return sb.String()
}

// FormatBytes 按指定宽度对齐
func FormatBytes(bm []uint8, width int) string {
	var size = len(bm) * BitsPerByte
	var sb strings.Builder
	sb.Grow(size + size/width + 1)
	var n = 0
	for i := 0; i < size; i++ {
		if n%width == 0 {
			sb.WriteByte('\n')
		}
		n++
		if TestBytesBit(bm, i) {
			sb.WriteByte('1')
		} else {
			sb.WriteByte('0')
		}
	}
	return sb.String()
}
