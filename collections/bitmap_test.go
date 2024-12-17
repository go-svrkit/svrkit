// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package collections

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitWordsCount(t *testing.T) {
	assert.Equal(t, 1, BitWordsCount(1))
	assert.Equal(t, 1, BitWordsCount(64))
	assert.Equal(t, 2, BitWordsCount(65))
	assert.Equal(t, 2, BitWordsCount(128))
	assert.Equal(t, 4, BitWordsCount(256))
	assert.Equal(t, 5, BitWordsCount(257))
}

func TestBitBytesCount(t *testing.T) {
	assert.Equal(t, 1, BitBytesCount(1))
	assert.Equal(t, 1, BitBytesCount(8))
	assert.Equal(t, 2, BitBytesCount(9))
	assert.Equal(t, 8, BitBytesCount(64))
	assert.Equal(t, 9, BitBytesCount(65))
}

func TestBitMap64_TestBit(t *testing.T) {
	tests := []struct {
		input []uint64
		idx   int
		want  bool
	}{
		{[]uint64{0b001001}, 100, false},
		{[]uint64{0b001001}, 0, true},
		{[]uint64{0b001001}, 1, false},
		{[]uint64{0b001001}, 3, true},
		{[]uint64{0x0, 0b1}, 64, true},
	}
	for i, tt := range tests {
		var bm = BitMap64(tt.input)
		var got = bm.TestBit(tt.idx)
		assert.Equalf(t, tt.want, got, "case %d index %d", i+1, tt.idx)
	}
}

func TestBitMap64_SetBit(t *testing.T) {
	var bm = BitMap64([]uint64{0b0})
	for i := 0; i < bitsPerWord; i++ {
		assert.False(t, bm.TestBit(i))
		bm.SetBit(i)
		assert.True(t, bm.TestBit(i))
	}
	assert.Equal(t, uint64(0xFFFFFFFFFFFFFFFF), bm[0])

	bm = BitMap64([]uint64{0b1, 0b1, 0b1})
	for _, i := range []int{0, 64, 128} {
		assert.True(t, bm.TestBit(i))
		bm.SetBit(i)
		assert.True(t, bm.TestBit(i))
	}
	for _, i := range []int{1, 65, 129} {
		assert.False(t, bm.TestBit(i))
		bm.SetBit(i)
		assert.True(t, bm.TestBit(i))
	}
}

func TestBitMap64_ClearBit(t *testing.T) {
	var bm = BitMap64([]uint64{0b1, 0b1, 0b1})
	for _, i := range []int{0, 64, 128} {
		assert.True(t, bm.TestBit(i))
		bm.ClearBit(i)
		assert.False(t, bm.TestBit(i))
	}
	for _, i := range []int{1, 65, 129} {
		assert.False(t, bm.TestBit(i))
		bm.ClearBit(i)
		assert.False(t, bm.TestBit(i))
	}
}

func TestBitMap64_FlipBit(t *testing.T) {
	var bm = BitMap64([]uint64{0b1, 0b1, 0b1})
	for _, i := range []int{0, 64, 128} {
		assert.True(t, bm.TestBit(i))
		bm.FlipBit(i)
		assert.False(t, bm.TestBit(i))
		bm.FlipBit(i)
		assert.True(t, bm.TestBit(i))
	}
	for _, i := range []int{1, 65, 129} {
		assert.False(t, bm.TestBit(i))
		bm.FlipBit(i)
		assert.True(t, bm.TestBit(i))
		bm.FlipBit(i)
		assert.False(t, bm.TestBit(i))
	}
}

func TestBitMap64_TestAndSetBit(t *testing.T) {
	var bm = BitMap64([]uint64{0b1, 0b1, 0b1})
	for _, i := range []int{0, 64, 128} {
		assert.True(t, bm.TestAndSetBit(i))
		assert.True(t, bm.TestBit(i))
	}
	for _, i := range []int{1, 65, 129} {
		assert.False(t, bm.TestAndSetBit(i))
		assert.True(t, bm.TestBit(i))
	}
}

func TestBitMap64_TestAndClearBit(t *testing.T) {
	var bm = BitMap64([]uint64{0b1, 0b1, 0b1})
	for _, i := range []int{0, 64, 128} {
		assert.True(t, bm.TestAndClearBit(i))
		assert.False(t, bm.TestBit(i))
	}
	for _, i := range []int{1, 65, 129} {
		assert.False(t, bm.TestAndClearBit(i))
		assert.False(t, bm.TestBit(i))
	}
}

func TestBitMap64_OnesCount(t *testing.T) {
	tests := []struct {
		input    []uint64
		expected int
	}{
		{nil, 0},
		{[]uint64{}, 0},
		{[]uint64{0b1}, 1},
		{[]uint64{0b100000001, 0b100000001}, 4},
		{[]uint64{0b101010101, 0b101010101}, 10},
	}
	for i, tt := range tests {
		var bm = BitMap64(tt.input)
		assert.Equalf(t, tt.expected, bm.OnesCount(), "case-%d", i+1)
	}
}

func TestBitMap64_IsZero(t *testing.T) {
	tests := []struct {
		input    []uint64
		expected bool
	}{
		{nil, true},
		{[]uint64{}, true},
		{[]uint64{0b1}, false},
		{[]uint64{0b0, 0b0, 0b0}, true},
		{[]uint64{0b0, 0b1, 0b0}, false},
	}
	for i, tt := range tests {
		var bm = BitMap64(tt.input)
		assert.Equalf(t, tt.expected, bm.IsZero(), "case-%d", i+1)
	}
}

func TestBitMap64_String(t *testing.T) {
	tests := []struct {
		input    []uint64
		expected string
	}{
		{nil, ""},
		{[]uint64{}, ""},
		{[]uint64{0b1}, "1000000000000000000000000000000000000000000000000000000000000000"},
		{[]uint64{0b1001001}, "1001001000000000000000000000000000000000000000000000000000000000"},
	}
	for i, tt := range tests {
		var bm = BitMap64(tt.input)
		assert.Equalf(t, tt.expected, bm.String(), "case-%d", i+1)
	}
}

func TestBitMap8_TestBit(t *testing.T) {
	tests := []struct {
		input []uint8
		idx   int
		want  bool
	}{
		{[]uint8{0b001001}, 100, false},
		{[]uint8{0b001001}, 0, true},
		{[]uint8{0b001001}, 1, false},
		{[]uint8{0b001001}, 3, true},
		{[]uint8{0x0, 0b1}, 8, true},
	}
	for i, tt := range tests {
		var bm = BitMap8(tt.input)
		var got = bm.TestBit(tt.idx)
		assert.Equalf(t, tt.want, got, "case %d index %d", i+1, tt.idx)
	}
}

func TestBitMap8_SetBit(t *testing.T) {
	var bm = BitMap8([]uint8{0b0})
	for i := 0; i < bitsPerByte; i++ {
		assert.False(t, bm.TestBit(i))
		bm.SetBit(i)
		assert.True(t, bm.TestBit(i))
	}
	assert.Equal(t, uint8(0xFF), bm[0])

	bm = BitMap8([]uint8{0b1, 0b1, 0b1})
	for _, i := range []int{0, 8, 16} {
		assert.True(t, bm.TestBit(i))
		bm.SetBit(i)
		assert.True(t, bm.TestBit(i))
	}
	for _, i := range []int{1, 9, 17} {
		assert.False(t, bm.TestBit(i))
		bm.SetBit(i)
		assert.True(t, bm.TestBit(i))
	}
}

func TestBitMap8_ClearBit(t *testing.T) {
	var bm = BitMap8([]uint8{0b1, 0b1, 0b1})
	for _, i := range []int{0, 8, 16} {
		assert.True(t, bm.TestBit(i))
		bm.ClearBit(i)
		assert.False(t, bm.TestBit(i))
	}
	for _, i := range []int{1, 9, 17} {
		assert.False(t, bm.TestBit(i))
		bm.ClearBit(i)
		assert.False(t, bm.TestBit(i))
	}
}

func TestBitMa8_FlipBit(t *testing.T) {
	var bm = BitMap8([]uint8{0b1, 0b1, 0b1})
	for _, i := range []int{0, 8, 16} {
		assert.True(t, bm.TestBit(i))
		bm.FlipBit(i)
		assert.False(t, bm.TestBit(i))
		bm.FlipBit(i)
		assert.True(t, bm.TestBit(i))
	}
	for _, i := range []int{1, 9, 17} {
		assert.False(t, bm.TestBit(i))
		bm.FlipBit(i)
		assert.True(t, bm.TestBit(i))
		bm.FlipBit(i)
		assert.False(t, bm.TestBit(i))
	}
}

func TestBitMap8_TestAndSetBit(t *testing.T) {
	var bm = BitMap8([]uint8{0b1, 0b1, 0b1})
	for _, i := range []int{0, 8, 16} {
		assert.True(t, bm.TestAndSetBit(i))
		assert.True(t, bm.TestBit(i))
	}
	for _, i := range []int{1, 9, 17} {
		assert.False(t, bm.TestAndSetBit(i))
		assert.True(t, bm.TestBit(i))
	}
}

func TestBitMap8_TestAndClearBit(t *testing.T) {
	var bm = BitMap8([]uint8{0b1, 0b1, 0b1})
	for _, i := range []int{0, 8, 16} {
		assert.True(t, bm.TestAndClearBit(i))
		assert.False(t, bm.TestBit(i))
	}
	for _, i := range []int{1, 9, 17} {
		assert.False(t, bm.TestAndClearBit(i))
		assert.False(t, bm.TestBit(i))
	}
}

func TestBitMap8_OnesCount(t *testing.T) {
	tests := []struct {
		input    []uint8
		expected int
	}{
		{nil, 0},
		{[]uint8{}, 0},
		{[]uint8{0b1}, 1},
		{[]uint8{0b10000001, 0b10000001}, 4},
		{[]uint8{0b01010101, 0b01010101}, 8},
	}
	for i, tt := range tests {
		var bm = BitMap8(tt.input)
		assert.Equalf(t, tt.expected, bm.OnesCount(), "case-%d", i+1)
	}
}

func TestBitMap8_IsZero(t *testing.T) {
	tests := []struct {
		input    []uint8
		expected bool
	}{
		{nil, true},
		{[]uint8{}, true},
		{[]uint8{0b1}, false},
		{[]uint8{0b0, 0b0, 0b0}, true},
		{[]uint8{0b0, 0b1, 0b0}, false},
	}
	for i, tt := range tests {
		var bm = BitMap8(tt.input)
		assert.Equalf(t, tt.expected, bm.IsZero(), "case-%d", i+1)
	}
}

func TestBitMap8_String(t *testing.T) {
	tests := []struct {
		input    []uint8
		expected string
	}{
		{nil, ""},
		{[]uint8{}, ""},
		{[]uint8{0b1}, "10000000"},
		{[]uint8{0b1101001}, "10010110"},
		{[]uint8{0b1101001, 0b1001111}, "1001011011110010"},
	}
	for i, tt := range tests {
		var bm = BitMap8(tt.input)
		assert.Equalf(t, tt.expected, bm.String(), "case-%d", i+1)
	}
}
