// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package bits

import (
	"math/bits"
	"strings"
)

const (
	bitsPerWord = 64 // 用uint64表示
	bitsPerByte = 8  // 用uint8表示
)

func BitWordsCount(bitSize int) int {
	return (bitSize + bitsPerWord - 1) / bitsPerWord
}

func BitBytesCount(bitSize int) int {
	return (bitSize + bitsPerByte - 1) / bitsPerByte
}

type BitMap64 []uint64

func NewBitMap64(bitSize int) BitMap64 {
	return make(BitMap64, BitWordsCount(bitSize))
}

// TestBit checks if the bit at the given index which starting from LSB is set.
func (bm BitMap64) TestBit(i int) bool {
	if i >= 0 && i < len(bm)*bitsPerWord {
		return bm[i/bitsPerWord]&(1<<(i%bitsPerWord)) != 0
	}
	return false
}

func (bm BitMap64) SetBit(i int) {
	var v = uint64(1) << (i % bitsPerWord)
	bm[i/bitsPerWord] |= v // 这里不进行边界检查
}

// MustSetBit 指定位是否为1，并且自动增长数组
func (bm BitMap64) MustSetBit(i int) BitMap64 {
	var n = BitWordsCount(i + 1)
	if n > len(bm) {
		var newb = make(BitMap64, n)
		copy(newb, bm)
		bm = newb
	}
	bm.SetBit(i)
	return bm
}

// ClearBit clears the bit at the given index.
func (bm BitMap64) ClearBit(i int) {
	var v = uint64(1) << (i % bitsPerWord)
	bm[i/bitsPerWord] &= ^v // 这里不进行边界检查
}

// FlipBit flips the bit at the given index.
func (bm BitMap64) FlipBit(i int) {
	bm[i/bitsPerWord] ^= 1 << (i % bitsPerWord)
}

// TestAndSetBit returns the old value of the bit at the given index and sets it to 1.
func (bm BitMap64) TestAndSetBit(i int) bool {
	var v = uint64(1) << (i % bitsPerWord)
	var index = i / bitsPerWord
	var old = bm[index]
	bm[index] |= v
	return old&v != 0
}

// TestAndClearBit returns the old value of the bit at the given index and clears it.
func (bm BitMap64) TestAndClearBit(i int) bool {
	var v = uint64(1) << (i % bitsPerWord)
	var index = i / bitsPerWord
	var old = bm[index]
	bm[index] &= ^v
	return old&v != 0
}

// OnesCount returns the number of bits set to 1.
func (bm BitMap64) OnesCount() int {
	var count int
	for i := 0; i < len(bm); i++ {
		if bm[i] != 0 {
			count += bits.OnesCount64(bm[i])
		}
	}
	return count
}

// IsZero 是否所有位都是0
func (bm BitMap64) IsZero() bool {
	for i := 0; i < len(bm); i++ {
		if bm[i] != 0 {
			return false
		}
	}
	return true
}

// String returns a string representation of the bitmap from LSB to MSB.
func (bm BitMap64) String() string {
	var size = len(bm) * bitsPerWord
	var sb strings.Builder
	sb.Grow(size)
	for i := 0; i < size; i++ {
		if bm.TestBit(i) {
			sb.WriteByte('1')
		} else {
			sb.WriteByte('0')
		}
	}
	return sb.String()
}

// FormattedString 按指定宽度对齐
func (bm BitMap64) FormattedString(width int) string {
	var size = len(bm) * bitsPerWord
	var sb strings.Builder
	var n = 0
	for i := 0; i < size; i++ {
		if n%width == 0 {
			sb.WriteByte('\n')
		}
		n++
		if bm.TestBit(i) {
			sb.WriteByte('1')
		} else {
			sb.WriteByte('0')
		}
	}
	return sb.String()
}

///////////////////////////////////////////////////////////////////////////////////////////

type BitMap8 []uint8

func NewBitMap8(bitSize int) BitMap8 {
	return make(BitMap8, BitBytesCount(bitSize))
}

func (bm BitMap8) TestBit(i int) bool {
	if i >= 0 && i < len(bm)*bitsPerByte {
		return bm[i/bitsPerByte]&(1<<(i%bitsPerByte)) != 0
	}
	return false
}

func (bm BitMap8) SetBit(i int) {
	var v = uint8(1) << (i % bitsPerByte)
	bm[i/bitsPerByte] |= v // 这里不进行边界检查
}

// MustSetBit 指定位是否为1，并且自动增长数组
func (bm BitMap8) MustSetBit(i int) BitMap8 {
	var n = BitBytesCount(i + 1)
	if n > len(bm) {
		var newb = make(BitMap8, n)
		copy(newb, bm)
		bm = newb
	}
	bm.SetBit(i)
	return bm
}

// ClearBit clears the bit at the given index.
func (bm BitMap8) ClearBit(i int) {
	var v = uint8(1) << (i % bitsPerByte)
	bm[i/bitsPerByte] &= ^v // 这里不进行边界检查
}

// FlipBit flips the bit at the given index.
func (bm BitMap8) FlipBit(i int) {
	bm[i/bitsPerByte] ^= 1 << (i % bitsPerByte)
}

// TestAndSetBit returns the old value of the bit at the given index and sets it to 1.
func (bm BitMap8) TestAndSetBit(i int) bool {
	var v = uint8(1) << (i % bitsPerByte)
	var index = i / bitsPerByte
	var old = bm[index]
	bm[index] |= v
	return old&v != 0
}

// TestAndClearBit returns the old value of the bit at the given index and clears it.
func (bm BitMap8) TestAndClearBit(i int) bool {
	var v = uint8(1) << (i % bitsPerByte)
	var index = i / bitsPerByte
	var old = bm[index]
	bm[index] &= ^v
	return old&v != 0
}

// OnesCount returns the number of bits set to 1.
func (bm BitMap8) OnesCount() int {
	var count int
	for i := 0; i < len(bm); i++ {
		if bm[i] != 0 {
			count += bits.OnesCount8(bm[i])
		}
	}
	return count
}

// IsZero 是否所有位都是0
func (bm BitMap8) IsZero() bool {
	for i := 0; i < len(bm); i++ {
		if bm[i] != 0 {
			return false
		}
	}
	return true
}

func (bm BitMap8) String() string {
	var size = len(bm) * bitsPerByte
	var sb strings.Builder
	sb.Grow(size)
	for i := 0; i < size; i++ {
		if bm.TestBit(i) {
			sb.WriteByte('1')
		} else {
			sb.WriteByte('0')
		}
	}
	return sb.String()
}

// FormattedString 按指定宽度对齐
func (bm BitMap8) FormattedString(width int) string {
	var size = len(bm) * bitsPerByte
	var sb strings.Builder
	var n = 0
	for i := 0; i < size; i++ {
		if n%width == 0 {
			sb.WriteByte('\n')
		}
		n++
		if bm.TestBit(i) {
			sb.WriteByte('1')
		} else {
			sb.WriteByte('0')
		}
	}
	return sb.String()
}
