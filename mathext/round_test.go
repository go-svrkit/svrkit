// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package mathext

import (
	"math"
	"testing"
)

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

func BenchmarkRound(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RoundFloat(0.1111, 1)
	}
}
