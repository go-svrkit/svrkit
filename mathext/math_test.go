// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package mathext

import (
	"math"
	"testing"
)

func TestMathAbs(t *testing.T) {
	tests := []struct {
		a        int
		expected int
	}{
		{10, 10},
		{-10, 10},
		{0, 0},
	}
	for i, tc := range tests {
		if output := Abs(tc.a); output != tc.expected {
			t.Fatalf("unexpected case %d output: %v != %v", i, output, tc.expected)
		}
	}
}

func TestMathMax(t *testing.T) {
	tests := []struct {
		a        int
		b        int
		expected int
	}{
		{10, 11, 11},
		{11, 11, 11},
		{-10, -11, -10},
		{-10, -10, -10},
		{0, 0, 0},
	}
	for i, tc := range tests {
		if output := Max(tc.a, tc.b); output != tc.expected {
			t.Fatalf("unexpected case %d output: %v != %v", i, output, tc.expected)
		}
	}
}

func TestMathMin(t *testing.T) {
	tests := []struct {
		a        int
		b        int
		expected int
	}{
		{10, 11, 10},
		{11, 11, 11},
		{-10, -11, -11},
		{-10, -10, -10},
		{0, 0, 0},
	}
	for i, tc := range tests {
		if output := Min(tc.a, tc.b); output != tc.expected {
			t.Fatalf("unexpected case %d output: %v != %v", i, output, tc.expected)
		}
	}
}

func TestMathDim(t *testing.T) {
	tests := []struct {
		a        int
		b        int
		expected int
	}{
		{10, 11, 0},
		{11, 11, 0},
		{-10, -11, 1},
		{-10, -10, 0},
		{0, 0, 0},
		{11, 10, 1},
	}
	for i, tc := range tests {
		if output := Dim(tc.a, tc.b); output != tc.expected {
			t.Fatalf("unexpected case %d output: %v != %v", i, output, tc.expected)
		}
	}
}

func TestSafeMulInt64(t *testing.T) {
	testCases := []struct {
		input1    int64
		input2    int64
		expected  int64
		expectErr bool
	}{
		{0, 0, 0, false},
		{-1, 1, -1, false},
		{math.MinInt32, 1, math.MinInt32, false},
		{math.MinInt32, math.MinInt32, math.MinInt32 * math.MinInt32, false},
		{math.MaxInt32, math.MaxInt32, math.MaxInt32 * math.MaxInt32, false},
		{math.MaxInt64, 1, math.MaxInt64, false},
		{math.MaxInt64, 2, 0, true},
		{math.MinInt64, 2, 0, true},
	}
	for i, tc := range testCases {
		r, overflow := SafeMulInt64(tc.input1, tc.input2)
		if !overflow && !tc.expectErr {
			t.Fatalf("case %d: unexpected error: %v", i, overflow)
			continue
		}
		if r != tc.expected {
			t.Fatalf("case %d: expected %d, got %d", i, tc.expected, r)
		}
	}
}

func TestSafeMulUint64(t *testing.T) {
	testCases := []struct {
		input1    uint64
		input2    uint64
		expected  uint64
		expectErr bool
	}{
		{math.MaxUint64, 1, math.MaxUint64, false},
		{math.MaxUint64, 2, 0, true},
	}
	for i, tc := range testCases {
		r, overflow := SafeMulUint64(tc.input1, tc.input2)
		if !overflow && !tc.expectErr {
			t.Fatalf("case %d: unexpected error: %v", i, overflow)
			continue
		}
		if r != tc.expected {
			t.Fatalf("case %d: expected %d, got %d", i, tc.expected, r)
		}
	}
}

func BenchmarkSafeMulInt64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SafeMulInt64(int64(i), int64(i))
	}
}

func BenchmarkSafeMulUint64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SafeMulUint64(uint64(i), uint64(i))
	}
}
