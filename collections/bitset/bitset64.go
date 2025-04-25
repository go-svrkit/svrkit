// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package bitset

import (
	"math/bits"
	"strings"
)

const BitsPerWord = 64 // 用uint64表示

func BitWordsCount(bitSize int) int {
	return (bitSize + BitsPerWord - 1) / BitsPerWord
}

// TestWordsBit checks if the bit at the given index which starting from LSB is set.
func TestWordsBit(bm []uint64, i int) bool {
	if i >= 0 && i < len(bm)*BitsPerWord {
		return bm[i/BitsPerWord]&(1<<(i%BitsPerWord)) != 0
	}
	return false
}

func SetWordsBit(bm []uint64, i int) {
	if i >= 0 && i < len(bm)*BitsPerWord {
		var v = uint64(1) << (i % BitsPerWord)
		bm[i/BitsPerWord] |= v
	}
}

// MustSetWordsBit 指定位是否为1，并且自动增长数组
func MustSetWordsBit(bm []uint64, i int) []uint64 {
	var n = BitWordsCount(i + 1)
	if n > len(bm) {
		var newbm = make([]uint64, n)
		copy(newbm, bm)
		bm = newbm
	}
	SetWordsBit(bm, i)
	return bm
}

// ClearWordsBit clears the bit at the given index.
func ClearWordsBit(bm []uint64, i int) {
	if i >= 0 && i < len(bm)*BitsPerWord {
		var v = uint64(1) << (i % BitsPerWord)
		bm[i/BitsPerWord] &= ^v
	}
}

// FlipWordsBit flips the bit at the given index.
func FlipWordsBit(bm []uint64, i int) {
	if i >= 0 && i < len(bm)*BitsPerWord {
		bm[i/BitsPerWord] ^= 1 << (i % BitsPerWord)
	}
}

// TestAndSetWordsBit Set a bit and return its old value
func TestAndSetWordsBit(bm []uint64, i int) bool {
	var v = uint64(1) << (i % BitsPerWord)
	var index = i / BitsPerWord
	if index >= 0 && index < len(bm) {
		var old = bm[index]
		bm[index] |= v
		return old&v != 0
	}
	return false
}

// TestAndClearWordsBit Clear a bit and return its old value
func TestAndClearWordsBit(bm []uint64, i int) bool {
	var v = uint64(1) << (i % BitsPerWord)
	var index = i / BitsPerWord
	if index >= 0 && index < len(bm) {
		var old = bm[index]
		bm[index] &= ^v
		return old&v != 0
	}
	return false
}

// OnesCountWords returns the number of bits set to 1.
func OnesCountWords(bm []uint64) int {
	var count int
	for i := 0; i < len(bm); i++ {
		if bm[i] != 0 {
			count += bits.OnesCount64(bm[i])
		}
	}
	return count
}

// IsAllWordsZero 是否所有位都是0
func IsAllWordsZero(bm []uint64) bool {
	for i := 0; i < len(bm); i++ {
		if bm[i] != 0 {
			return false
		}
	}
	return true
}

// WordsToString returns a string representation of the bitmap from LSB to MSB.
func WordsToString(bm []uint64) string {
	var size = len(bm) * BitsPerWord
	var sb strings.Builder
	sb.Grow(size)
	for i := 0; i < size; i++ {
		if TestWordsBit(bm, i) {
			sb.WriteByte('1')
		} else {
			sb.WriteByte('0')
		}
	}
	return sb.String()
}

// FormatWords 按指定宽度对齐
func FormatWords(bm []uint64, width int) string {
	var size = len(bm) * BitsPerWord
	var sb strings.Builder
	sb.Grow(size + size/width + 1)
	var n = 0
	for i := 0; i < size; i++ {
		if n%width == 0 {
			sb.WriteByte('\n')
		}
		n++
		if TestWordsBit(bm, i) {
			sb.WriteByte('1')
		} else {
			sb.WriteByte('0')
		}
	}
	return sb.String()
}
