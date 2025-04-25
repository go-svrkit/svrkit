package bitset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BitBytesCount(t *testing.T) {
	tests := []struct {
		input int
		want  int
	}{
		{1, 1},
		{8, 1},
		{9, 2},
		{64, 8},
		{65, 9},
	}
	for i, tt := range tests {
		var got = BitBytesCount(tt.input)
		assert.Equalf(t, tt.want, got, "case %d ", i+1)
	}
}

func Test_TestBytesBit(t *testing.T) {
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
		var got = TestBytesBit(tt.input, tt.idx)
		assert.Equalf(t, tt.want, got, "case %d index %d", i+1, tt.idx)
	}
}

func Test_SetBytesBit(t *testing.T) {
	tests := []struct {
		input []uint8
		idx   int
	}{
		{[]uint8{0b001001}, 100},
		{[]uint8{0b001001}, 0},
		{[]uint8{0b001001}, 1},
		{[]uint8{0b001001}, 3},
		{[]uint8{0x0, 0b1}, 9},
		{[]uint8{0x0, 0b1}, 12},
	}
	for i, tt := range tests {
		var outRange = tt.idx >= len(tt.input)*BitsPerByte
		var oldBit = TestBytesBit(tt.input, tt.idx)
		SetBytesBit(tt.input, tt.idx)
		var newBit = TestBytesBit(tt.input, tt.idx)
		if outRange {
			assert.Equalf(t, oldBit, newBit, "case %d index %d: old bit should equal to new bit", i+1, tt.idx)
		} else {
			assert.Truef(t, newBit, "case %d index %d: new bit should been set", i+1, tt.idx)
		}
	}
}

func Test_ClearBytesBit(t *testing.T) {
	tests := []struct {
		input []uint8
		idx   int
	}{
		{[]uint8{0b1, 0b1, 0b1}, 0},
		{[]uint8{0b1, 0b1, 0b1}, 1},
		{[]uint8{0b1, 0b1, 0b1}, 64},
		{[]uint8{0b1, 0b1, 0b1}, 65},
		{[]uint8{0b1, 0b1, 0b1}, 128},
		{[]uint8{0b1, 0b1, 0b1}, 128},
	}
	for i, tt := range tests {
		var outRange = tt.idx >= len(tt.input)*BitsPerByte
		var oldBit = TestBytesBit(tt.input, tt.idx)
		ClearBytesBit(tt.input, tt.idx)
		var newBit = TestBytesBit(tt.input, tt.idx)
		if outRange {
			assert.Equalf(t, oldBit, newBit, "case %d index %d: old bit should not been set", i+1, tt.idx)
		} else {
			assert.Falsef(t, newBit, "case %d index %d: new bit should been cleared", i+1, tt.idx)
		}
	}
}

func Test_FlipBytesBit(t *testing.T) {
	tests := []struct {
		input []uint8
		idx   int
	}{
		{[]uint8{0b1, 0b1, 0b1}, 0},
		{[]uint8{0b1, 0b1, 0b1}, 1},
		{[]uint8{0b1, 0b1, 0b1}, 64},
		{[]uint8{0b1, 0b1, 0b1}, 65},
		{[]uint8{0b1, 0b1, 0b1}, 128},
		{[]uint8{0b1, 0b1, 0b1}, 128},
	}
	for i, tt := range tests {
		var outRange = tt.idx >= len(tt.input)*BitsPerByte
		var oldBit = TestBytesBit(tt.input, tt.idx)
		FlipBytesBit(tt.input, tt.idx)
		var newBit = TestBytesBit(tt.input, tt.idx)
		if outRange {
			assert.Equalf(t, oldBit, newBit, "case %d index %d: old bit should not been set", i+1, tt.idx)
		} else {
			assert.NotEqualf(t, oldBit, newBit, "case %d index %d: old bit should not equal to new bit", i+1, tt.idx)
		}
	}
}

func Test_TestAndSetBytesBit(t *testing.T) {
	tests := []struct {
		input []uint8
		idx   int
	}{
		{[]uint8{0b1, 0b1, 0b1}, 0},
		{[]uint8{0b1, 0b1, 0b1}, 1},
		{[]uint8{0b1, 0b1, 0b1}, 64},
		{[]uint8{0b1, 0b1, 0b1}, 65},
		{[]uint8{0b1, 0b1, 0b1}, 128},
		{[]uint8{0b1, 0b1, 0b1}, 128},
	}
	for i, tt := range tests {
		var outRange = tt.idx >= len(tt.input)*BitsPerByte
		var oldBit = TestBytesBit(tt.input, tt.idx)
		var out = TestAndSetBytesBit(tt.input, tt.idx)
		var newBit = TestBytesBit(tt.input, tt.idx)
		assert.Equal(t, oldBit, out, "case %d index %d: old bit should equal to out", i+1, tt.idx)
		if outRange {
			assert.Equalf(t, oldBit, newBit, "case %d index %d: old bit should not been set", i+1, tt.idx)
		} else {
			assert.Truef(t, newBit, "case %d index %d: new bit should be set", i+1, tt.idx)
		}
	}
}

func Test_TestAndClearBytesBit(t *testing.T) {
	tests := []struct {
		input []uint8
		idx   int
	}{
		{[]uint8{0b1, 0b1, 0b1}, 0},
		{[]uint8{0b1, 0b1, 0b1}, 1},
		{[]uint8{0b1, 0b1, 0b1}, 64},
		{[]uint8{0b1, 0b1, 0b1}, 65},
		{[]uint8{0b1, 0b1, 0b1}, 128},
		{[]uint8{0b1, 0b1, 0b1}, 129},
	}
	for i, tt := range tests {
		var outRange = tt.idx >= len(tt.input)*BitsPerByte
		var oldBit = TestBytesBit(tt.input, tt.idx)
		var out = TestAndClearBytesBit(tt.input, tt.idx)
		var newBit = TestBytesBit(tt.input, tt.idx)
		assert.Equal(t, oldBit, out, "case %d index %d: old bit should equal to out", i+1, tt.idx)
		if outRange {
			assert.Equalf(t, oldBit, newBit, "case %d index %d: old bit should not been set", i+1, tt.idx)
		} else {
			assert.Falsef(t, newBit, "case %d index %d: new bit should not be set", i+1, tt.idx)
		}
	}
}

func Test_OnesCountBytes(t *testing.T) {
	tests := []struct {
		input    []uint8
		expected int
	}{
		{nil, 0},
		{[]uint8{}, 0},
		{[]uint8{0b1}, 1},
		{[]uint8{0b10000001, 0b10000001}, 4},
		{[]uint8{0b10101011, 0b11010101}, 10},
	}
	for i, tt := range tests {
		var out = OnesCountBytes(tt.input)
		assert.Equalf(t, tt.expected, out, "case-%d", i+1)
	}
}

func Test_IsAllBytesZero(t *testing.T) {
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
		var got = IsAllBytesZero(tt.input)
		assert.Equalf(t, tt.expected, got, "case-%d", i+1)
	}
}
