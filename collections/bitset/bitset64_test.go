// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package bitset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitWordsCount(t *testing.T) {
	tests := []struct {
		input int
		want  int
	}{
		{1, 1},
		{64, 1},
		{65, 2},
		{128, 2},
		{129, 3},
		{256, 4},
	}
	for i, tt := range tests {
		var got = BitWordsCount(tt.input)
		assert.Equalf(t, tt.want, got, "case %d", i+1)
	}
}

func Test_TestWordsBit(t *testing.T) {
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
		var got = TestWordsBit(tt.input, tt.idx)
		assert.Equalf(t, tt.want, got, "case %d index %d", i+1, tt.idx)
	}
}

func Test_SetWordsBit(t *testing.T) {
	tests := []struct {
		input []uint64
		idx   int
	}{
		{[]uint64{0b001001}, 100},
		{[]uint64{0b001001}, 0},
		{[]uint64{0b001001}, 1},
		{[]uint64{0b001001}, 33},
		{[]uint64{0x0, 0b1}, 64},
		{[]uint64{0x0, 0b1}, 128},
	}
	for i, tt := range tests {
		var outRange = tt.idx >= len(tt.input)*BitsPerWord
		var oldBit = TestWordsBit(tt.input, tt.idx)
		SetWordsBit(tt.input, tt.idx)
		var newBit = TestWordsBit(tt.input, tt.idx)
		if outRange {
			assert.Equalf(t, oldBit, newBit, "case %d index %d: old bit should equal to new bit", i+1, tt.idx)
		} else {
			assert.Truef(t, newBit, "case %d index %d: new bit should been set", i+1, tt.idx)
		}
	}
}

func Test_ClearWordsBit(t *testing.T) {
	tests := []struct {
		input []uint64
		idx   int
	}{
		{[]uint64{0b1, 0b1, 0b1}, 0},
		{[]uint64{0b1, 0b1, 0b1}, 1},
		{[]uint64{0b1, 0b1, 0b1}, 64},
		{[]uint64{0b1, 0b1, 0b1}, 65},
		{[]uint64{0b1, 0b1, 0b1}, 128},
		{[]uint64{0b1, 0b1, 0b1}, 128},
	}
	for i, tt := range tests {
		var outRange = tt.idx >= len(tt.input)*BitsPerWord
		var oldBit = TestWordsBit(tt.input, tt.idx)
		ClearWordsBit(tt.input, tt.idx)
		var newBit = TestWordsBit(tt.input, tt.idx)
		if outRange {
			assert.Equalf(t, oldBit, newBit, "case %d index %d: old bit should not been set", i+1, tt.idx)
		} else {
			assert.Falsef(t, newBit, "case %d index %d: new bit should been cleared", i+1, tt.idx)
		}
	}
}

func Test_FlipWordsBit(t *testing.T) {
	tests := []struct {
		input []uint64
		idx   int
	}{
		{[]uint64{0b1, 0b1, 0b1}, 0},
		{[]uint64{0b1, 0b1, 0b1}, 1},
		{[]uint64{0b1, 0b1, 0b1}, 64},
		{[]uint64{0b1, 0b1, 0b1}, 65},
		{[]uint64{0b1, 0b1, 0b1}, 128},
		{[]uint64{0b1, 0b1, 0b1}, 128},
	}
	for i, tt := range tests {
		var outRange = tt.idx >= len(tt.input)*BitsPerWord
		var oldBit = TestWordsBit(tt.input, tt.idx)
		FlipWordsBit(tt.input, tt.idx)
		var newBit = TestWordsBit(tt.input, tt.idx)
		if outRange {
			assert.Equalf(t, oldBit, newBit, "case %d index %d: old bit should not been set", i+1, tt.idx)
		} else {
			assert.NotEqualf(t, oldBit, newBit, "case %d index %d: old bit should not equal to new bit", i+1, tt.idx)
		}
	}
}

func Test_TestAndSetWordsBit(t *testing.T) {
	tests := []struct {
		input []uint64
		idx   int
	}{
		{[]uint64{0b1, 0b1, 0b1}, 0},
		{[]uint64{0b1, 0b1, 0b1}, 1},
		{[]uint64{0b1, 0b1, 0b1}, 64},
		{[]uint64{0b1, 0b1, 0b1}, 65},
		{[]uint64{0b1, 0b1, 0b1}, 128},
		{[]uint64{0b1, 0b1, 0b1}, 1024},
	}
	for i, tt := range tests {
		var outRange = tt.idx >= len(tt.input)*BitsPerWord
		var oldBit = TestWordsBit(tt.input, tt.idx)
		var out = TestAndSetWordsBit(tt.input, tt.idx)
		var newBit = TestWordsBit(tt.input, tt.idx)
		assert.Equal(t, oldBit, out, "case %d index %d: old bit should equal to out", i+1, tt.idx)
		if outRange {
			assert.Equalf(t, oldBit, newBit, "case %d index %d: old bit should not been set", i+1, tt.idx)
		} else {
			assert.Truef(t, newBit, "case %d index %d: new bit should be set", i+1, tt.idx)
		}
	}
}

func Test_TestAndClearWordsBit(t *testing.T) {
	tests := []struct {
		input []uint64
		idx   int
	}{
		{[]uint64{0b1, 0b1, 0b1}, 0},
		{[]uint64{0b1, 0b1, 0b1}, 1},
		{[]uint64{0b1, 0b1, 0b1}, 64},
		{[]uint64{0b1, 0b1, 0b1}, 65},
		{[]uint64{0b1, 0b1, 0b1}, 128},
		{[]uint64{0b1, 0b1, 0b1}, 1024},
	}
	for i, tt := range tests {
		var outRange = tt.idx >= len(tt.input)*BitsPerWord
		var oldBit = TestWordsBit(tt.input, tt.idx)
		var out = TestAndClearWordsBit(tt.input, tt.idx)
		var newBit = TestWordsBit(tt.input, tt.idx)
		assert.Equal(t, oldBit, out, "case %d index %d: old bit should equal to out", i+1, tt.idx)
		if outRange {
			assert.Equalf(t, oldBit, newBit, "case %d index %d: old bit should not been set", i+1, tt.idx)
		} else {
			assert.Falsef(t, newBit, "case %d index %d: new bit should not be set", i+1, tt.idx)
		}
	}
}

func Test_OnesCountWords(t *testing.T) {
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
		var out = OnesCountWords(tt.input)
		assert.Equalf(t, tt.expected, out, "case-%d", i+1)
	}
}

func Test_IsAllWordsZero(t *testing.T) {
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
		var got = IsAllWordsZero(tt.input)
		assert.Equalf(t, tt.expected, got, "case-%d", i+1)
	}
}

func Test_WordsToString(t *testing.T) {
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
		var out = WordsToString(tt.input)
		assert.Equalf(t, tt.expected, out, "case-%d", i+1)
	}
}
