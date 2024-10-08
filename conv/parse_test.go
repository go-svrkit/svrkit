// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package conv

import (
	"fmt"
	"math"
	"reflect"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseBool(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", false},
		{"0", false},
		{"-", false},
		{"xyz", false},
		{"1", true},
		{"y", true},
		{"Y", true},
		{"t", true},
		{"T", true},
		{"on", true},
		{"ON", true},
		{"yes", true},
		{"YES", true},
		{"true", true},
		{"TRUE", true},
		{"True", true},
	}
	for i, tc := range tests {
		var b = ParseBool(tc.input)
		assert.Equal(t, tc.expected, b, "test case %d", i)
	}
}

func TestParseInt32(t *testing.T) {
	tests := []struct {
		input    string
		hasErr   bool
		expected int32
	}{
		{"xyz", true, 0},
		{"0", false, 0},
		{"1234", false, 1234},
		{"2147483647", false, math.MaxInt32},
		{"-2147483648", false, math.MinInt32},
		{"2147483648", true, 0},
		{"-2147483649", true, 0},
	}
	for i, tc := range tests {
		n, err := ParseI32(tc.input)
		if tc.hasErr {
			assert.NotNil(t, err, "test case %d", i)
		} else {
			assert.Equal(t, tc.expected, n, "test case %d", i)
		}
	}
}

func TestParseUint32(t *testing.T) {
	tests := []struct {
		input    string
		hasErr   bool
		expected uint32
	}{
		{"-1", true, 0},
		{"0", false, 0},
		{"1234", false, 1234},
		{"4294967295", false, math.MaxUint32},
		{"4294967296", true, 0},
	}
	for i, tc := range tests {
		n, err := ParseU32(tc.input)
		if tc.hasErr {
			assert.NotNil(t, err, "test case %d", i)
		} else {
			assert.Equal(t, tc.expected, n, "test case %d", i)
		}
	}
}

func TestParseInt64(t *testing.T) {
	tests := []struct {
		input    string
		hasErr   bool
		expected int64
	}{
		{"xyz", true, 0},
		{"0", false, 0},
		{"1234", false, 1234},
		{"9223372036854775807", false, math.MaxInt64},
		{"-9223372036854775808", false, math.MinInt64},
		{"9223372036854775808", true, 0},
		{"-9223372036854775809", true, 0},
	}
	for i, tc := range tests {
		n, err := ParseI64(tc.input)
		if tc.hasErr {
			assert.NotNil(t, err, "test case %d", i)
		} else {
			assert.Equal(t, tc.expected, n, "test case %d", i)
		}
	}
}

func TestParseUint64(t *testing.T) {
	tests := []struct {
		input    string
		hasErr   bool
		expected uint64
	}{
		{"-1", true, 0},
		{"0", false, 0},
		{"1234", false, 1234},
		{"18446744073709551615", false, math.MaxUint64},
		{"18446744073709551616", true, 0},
	}
	for i, tc := range tests {
		n, err := ParseU64(tc.input)
		if tc.hasErr {
			assert.NotNil(t, err, "test case %d", i)
		} else {
			assert.Equal(t, tc.expected, n, "test case %d", i)
		}
	}
}

func TestParseFloat32(t *testing.T) {
	tests := []struct {
		input    string
		hasErr   bool
		expected float32
	}{
		{"xyz", true, 0},
		{"0", false, 0},
		{"3.40282346638528859811704183484516925440e+38", false, math.MaxFloat32},
		{"1.401298464324817070923729583289916131280e-45", false, math.SmallestNonzeroFloat32},
	}
	for i, tc := range tests {
		n, err := ParseF32(tc.input)
		if tc.hasErr {
			assert.NotNil(t, err, "test case %d", i)
		} else {
			assert.Equal(t, tc.expected, n, "test case %d", i)
		}
	}
}

func TestParseTo(t *testing.T) {
	{
		b, err := ParseTo[bool]("1")
		assert.Nil(t, err)
		assert.Equal(t, true, b)
	}
	{
		s, err := ParseTo[string]("1234")
		assert.Nil(t, err)
		assert.Equal(t, "1234", s)
	}
	{
		n, err := ParseTo[int]("-1234567")
		assert.Nil(t, err)
		assert.Equal(t, -1234567, n)
	}
	{
		n, err := ParseTo[uint]("1234567")
		assert.Nil(t, err)
		assert.Equal(t, uint(1234567), n)
	}
	{
		f, err := ParseTo[float32]("3.14159")
		assert.Nil(t, err)
		assert.Equal(t, float32(3.14159), f)
	}
}

func TestParseSlice(t *testing.T) {
	const sep = "|"
	assert.Equal(t, len(ParseSlice[int]("", sep)), 0)
	assert.Equal(t, len(ParseSlice[int](sep, sep)), 0)
	assert.True(t, slices.Equal([]int{1}, ParseSlice[int]("1", sep)))
	assert.True(t, slices.Equal([]int{1}, ParseSlice[int]("1||", sep)))
	assert.True(t, slices.Equal([]int{1}, ParseSlice[int]("||1||", sep)))
	assert.True(t, slices.Equal([]int{1, 2, 3}, ParseSlice[int]("|1|2|3|", sep)))

	assert.True(t, slices.Equal([]string{"usr", "local", "bin"}, ParseSlice[string]("/usr/local/bin", "/")))
}

func TestParseKeyValues(t *testing.T) {
	tests := []struct {
		input string
		want1 []int
		want2 []int
	}{
		{"", nil, nil},
		{"", []int{}, []int{}},
	}
	for i, tc := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			out1, out2 := ParseKeyValues[int, int](tc.input, SepEqualSign, SepComma)
			assert.True(t, slices.Equal(out1, tc.want1))
			assert.True(t, slices.Equal(out2, tc.want2))
		})
	}
}

func TestParseMap(t *testing.T) {
	tests := []struct {
		input    string
		expected map[int]int
	}{
		{"", map[int]int{}},
		{"1:2", map[int]int{1: 2}},
		{"|1:2|", map[int]int{1: 2}},
		{"||1:2||||", map[int]int{1: 2}},
		{"1:2|3:4", map[int]int{1: 2, 3: 4}},
		{"  1 : 2 | 3 : 4|  ", map[int]int{1: 2, 3: 4}},
	}
	for i, tc := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			var out = ParseMap[int, int](tc.input, SepColon, SepVerticalBar)
			if !reflect.DeepEqual(out, tc.expected) {
				t.Fatalf("unexpected ParseMap(%s)  result: %v != %v", tc.input, out, tc.expected)
			}
		})

	}
}

func BenchmarkParseI32(b *testing.B) {
	var total int64
	var val int32
	for i := 0; i < b.N; i++ {
		val, _ = ParseI32("1234567890")
		total += int64(val)
		val, _ = ParseI32("2147483647")
		total += int64(val)
		val, _ = ParseI32("-2147483648")
		total += int64(val)
	}
}

func BenchmarkParseToInt(b *testing.B) {
	var total int64
	var val int32
	for i := 0; i < b.N; i++ {
		val, _ = ParseTo[int32]("1234567890")
		total += int64(val)
		val, _ = ParseTo[int32]("2147483647")
		total += int64(val)
		val, _ = ParseTo[int32]("-2147483648")
		total += int64(val)
	}
}

// AMD Ryzen 5 6-Core Processor
// BenchmarkParseI32-4     	    22760205                49.79 ns/op
// BenchmarkParseToInt-4        19730870                59.22 ns/op
