// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

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

func TestMathTruncate(t *testing.T) {
	var tests = []struct {
		input     float64
		precision int
		expected  float64
	}{
		{1.2345678, 1, 1.2},
		{1.2345678, 2, 1.23},
		{1.2345678, 3, 1.234},
		{-1.2345678, 1, -1.2},
		{-1.2345678, 2, -1.23},
		{-1.2345678, 3, -1.234},
	}
	for _, tc := range tests {
		var output = Truncate(tc.input, tc.precision)
		if output != tc.expected {
			t.Fatalf("%f %d expected %f, but got %f", tc.input, tc.precision, tc.expected, output)
		}
	}
}

func TestRound(t *testing.T) {
	for _, c := range []struct {
		number   float64
		decimals int
		result   float64
	}{
		{0.1111, 1, 0.1},
		{-0.1111, 2, -0.11},
		{5.3253, 3, 5.325},
		{5.3258, 3, 5.326},
		{5.3253, 0, 5.0},
		{5.55, 1, 5.6},
	} {
		m := RoundFloat(c.number, c.decimals)
		if m != c.result {
			t.Errorf("%.1f != %.1f", m, c.result)
		}
	}
}

func roundHalfSimple(f float64) float64 {
	return math.Floor(f + 0.5)
}

func TestMathRoundHalf(t *testing.T) {
	var roundTests = []struct {
		input    float64
		expected float64
	}{
		{-0.49999999999999994, -0.0}, // -0.5+epsilon
		{-0.5, -1},
		{-0.5000000000000001, -1}, // -0.5-epsilon
		{0, 0},
		{0.49999999999999994, 0}, // 0.5-epsilon
		{0.5, 1},
		{0.0, 0.0},
		{0.5000000000000001, 1},  // 0.5+epsilon
		{1.390671161567e-309, 0}, // denormal
		{2.2517998136852485e+15, 2.251799813685249e+15}, // 1 bit fraction
		{4.503599627370497e+15, 4.503599627370497e+15},  // large integer
		{math.Inf(-1), math.Inf(-1)},
		{math.Inf(1), math.Inf(1)},
		{math.NaN(), math.NaN()},
	}
	for i, tc := range roundTests {
		//var output = RoundFloat(tc.input, 0)
		var output = roundHalfSimple(tc.input)
		if math.IsNaN(output) {
			if !math.IsNaN(tc.expected) {
				t.Fatalf("%d: %f => %f != %f", i+1, tc.input, output, tc.expected)
			}
		} else {
			if output != tc.expected {
				t.Fatalf("%d: %f => %f != %f", i+1, tc.input, output, tc.expected)
			}
		}
	}
}

func BenchmarkRound(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RoundFloat(0.1111, 1)
	}
}
