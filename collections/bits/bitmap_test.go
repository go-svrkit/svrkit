// Copyright © 2017 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE

package bits

import (
	"testing"
)

func wordToBytes(words []uint64) []uint8 {
	if len(words) == 0 {
		return nil
	}
	var b = make([]uint8, len(words)*8)
	var x = 0
	for i := 0; i < len(words); i++ {
		var v = words[i]
		// little endian
		b[x] = byte(v)
		b[x+1] = byte(v >> 8)
		b[x+2] = byte(v >> 16)
		b[x+3] = byte(v >> 24)
		b[x+4] = byte(v >> 32)
		b[x+5] = byte(v >> 40)
		b[x+6] = byte(v >> 48)
		b[x+7] = byte(v >> 56)
		x += 8
	}
	return b
}

func TestBitMapTestBit(t *testing.T) {
	tests := []struct {
		name     string
		input    []uint64 // 位组
		expected []int    // 位为1的索引
	}{
		{"one word", []uint64{0b01}, []int{0}},
		{"multi words", []uint64{0b1000000010000000100000001}, []int{0, 8, 16, 24}},
		{"two words", []uint64{0b11, 0b11}, []int{0, 1, 64, 65}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var bm1 = BitMap64(tc.input)
			var bm2 = BitMap8(wordToBytes(tc.input))
			for _, idx := range tc.expected {
				if !bm1.TestBit(idx) {
					t.Fatalf("BitMap64: %b unexpected bit result at %d", tc.input, idx)
				}
				if !bm2.TestBit(idx) {
					t.Fatalf("BitMap8: %b unexpected bit result at %d", tc.input, idx)
				}
			}
		})
	}
}

func TestBitMapSetBit(t *testing.T) {
	tests := []struct {
		name     string
		bitSize  int   // 位长度
		expected []int // 设置为1的位的索引
	}{
		{"one word", 64, []int{0, 1, 2, 3, 4, 5, 6, 60, 61, 62, 63}},
		{"two words", 73, []int{0, 8, 16, 24, 32, 40, 48, 56, 64, 72}},
	}
	var runTest = func(t *testing.T, bm BitMapper, idx int) {
		bm.SetBit(idx)
		if !bm.TestBit(idx) {
			t.Fatalf("%T: unexpected bit result at %d", bm, idx)
		}
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var bm1 = NewBitMap64(tc.bitSize)
			var bm2 = NewBitMap8(tc.bitSize)
			for _, idx := range tc.expected {
				runTest(t, bm1, idx)
				runTest(t, bm2, idx)
			}
		})
	}
}

func TestBitMapMustSetBit(t *testing.T) {
	tests := []struct {
		name     string
		expected []int // 设置为1的位的索引
	}{
		{"simple1", []int{0, 1, 2, 3, 4, 5, 6, 60, 61, 62, 63, 64, 128, 1023}},
		{"simple2", []int{0, 8, 16, 24, 32, 40, 48, 56, 64, 72, 128, 1023}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var bm1 = NewBitMap64(0) // 都初始为0
			var bm2 = NewBitMap8(0)
			for _, idx := range tc.expected {
				bm1 = bm1.MustSetBit(idx)
				if !bm1.TestBit(idx) {
					t.Fatalf("BitMap64: unexpected bit result at %d", idx)
				}
				bm2 = bm2.MustSetBit(idx)
				if !bm2.TestBit(idx) {
					t.Fatalf("BitMap8: unexpected bit result at %d", idx)
				}
			}
			t.Logf("length of BitMap64 %d", len(bm1))
			t.Logf("length of BitMap8 %d", len(bm2))
		})
	}
}

func TestBitMapClearBit(t *testing.T) {
	tests := []struct {
		name    string
		input   []uint64 // 位组
		indexes []int    // 设置为0的位的索引
	}{
		{"one word", []uint64{0b01010101010}, []int{1, 3, 5, 7, 9}},
		{"two words", []uint64{0b1111, 0b1111}, []int{1, 2, 3, 4, 64, 65, 66, 67}},
	}
	var runTest = func(t *testing.T, bm BitMapper, idx int) {
		bm.ClearBit(idx)
		if bm.TestBit(idx) {
			t.Fatalf("%T: unexpected bit result at %d", bm, idx)
		}
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var bm1 = BitMap64(tc.input)
			var bm2 = BitMap8(wordToBytes(tc.input))
			for _, idx := range tc.indexes {
				runTest(t, bm1, idx)
				runTest(t, bm2, idx)
			}
		})
	}
}

func TestBitMapFlipBit(t *testing.T) {
	tests := []struct {
		name    string
		bitSize int   // 位长度
		indexes []int // 测试的索引
	}{
		{"one word", 64, []int{0, 8, 16, 24, 32, 40, 63}},
		{"two words", 100, []int{0, 10, 20, 30, 40, 50, 99}},
	}

	var runTest = func(t *testing.T, bm BitMapper, idx int) {
		// 设置前应该为0
		if bm.TestBit(idx) {
			t.Fatalf("%T: unexpected bit result at %d", bm, idx)
		}
		// 把0翻转为1
		bm.FlipBit(idx)
		if !bm.TestBit(idx) {
			t.Fatalf("%T: unexpected bit result at %d", bm, idx)
		}
		// 把1翻转为0
		bm.FlipBit(idx)
		if bm.TestBit(idx) {
			t.Fatalf("%T: unexpected bit result at %d", bm, idx)
		}
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var bm1 = NewBitMap64(tc.bitSize)
			var bm2 = NewBitMap8(tc.bitSize)
			for _, idx := range tc.indexes {
				runTest(t, bm1, idx)
				runTest(t, bm2, idx)
			}
		})
	}
}

func TestBitMapTestAndSetBit(t *testing.T) {
	tests := []struct {
		bitSize int   // 位长度
		indexes []int // 测试的索引
	}{
		{64, []int{0, 8, 16, 24, 32, 40, 63}},
		{100, []int{0, 10, 20, 30, 40, 50, 99}},
	}
	var runTest = func(t *testing.T, bm BitMapper, idx int) {
		// 设置前应该为0
		if bm.TestBit(idx) {
			t.Fatalf("%T: unexpected bit result at %d", bm, idx)
		}
		// 没有设置过，应该为0
		if bm.TestAndSetBit(idx) {
			t.Fatalf("%T: unexpected bit result at %d", bm, idx)
		}
		// 设置后应该为1
		if !bm.TestBit(idx) {
			t.Fatalf("%T: unexpected bit result at %d", bm, idx)
		}
		// 再次设置应该为1
		if !bm.TestAndSetBit(idx) {
			t.Fatalf("%T: unexpected bit result at %d", bm, idx)
		}
	}
	for _, tc := range tests {
		var bm1 = NewBitMap64(tc.bitSize)
		var bm2 = NewBitMap8(tc.bitSize)
		for _, idx := range tc.indexes {
			runTest(t, bm1, idx)
			runTest(t, bm2, idx)
		}
	}
}

func TestBitMapTestAndClearBit(t *testing.T) {
	tests := []struct {
		input   []uint64 // 位组
		indexes []int    // 测试的索引
	}{
		{[]uint64{0b111100001111}, []int{1, 2, 3, 4, 5, 10}},
		{[]uint64{0b01010101, 0b10101010}, []int{0, 8, 64, 65}},
	}
	var runTest = func(t *testing.T, bm BitMapper, idx int) {
		var old = bm.TestBit(idx)
		if old {
			if !bm.TestAndClearBit(idx) {
				t.Fatalf("%T: unexpected bit result at %d", bm, idx)
			}
		}
		if bm.TestAndClearBit(idx) {
			t.Fatalf("%T: unexpected bit result at %d", bm, idx)
		}
	}
	for _, tc := range tests {
		var bm1 = BitMap64(tc.input)
		var bm2 = BitMap8(wordToBytes(tc.input))
		for _, idx := range tc.indexes {
			runTest(t, bm1, idx)
			runTest(t, bm2, idx)
		}
	}
}

func TestBitMapOnesCount(t *testing.T) {
	tests := []struct {
		input    []uint64 // 位组
		expected int      // 为1的位的数量
	}{
		{nil, 0},
		{[]uint64{0b111100001111}, 8},
		{[]uint64{0b1111, 0b1111}, 8},
	}
	var runTest = func(t *testing.T, bm BitMapper, expected int) {
		var cnt = bm.OnesCount()
		if cnt != expected {
			t.Fatalf("%T: unexpected bit count %d != %d", bm, cnt, expected)
		}
	}
	for _, tc := range tests {
		var bm1 = BitMap64(tc.input)
		var bm2 = BitMap8(wordToBytes(tc.input))
		runTest(t, bm1, tc.expected)
		runTest(t, bm2, tc.expected)
	}
}
