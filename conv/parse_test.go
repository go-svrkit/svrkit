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

func TestParseI8(t *testing.T) {
	tests := []struct {
		input    string
		hasErr   bool
		expected int8
	}{
		{"xyz", true, 0},
		{"0", false, 0},
		{"127", false, 127},
		{"128", true, 0},
		{"-128", false, math.MinInt8},
		{"-129", true, 0},
		{"-2147483649", true, 0},
	}
	for i, tc := range tests {
		n, err := ParseI8(tc.input)
		if tc.hasErr {
			assert.NotNil(t, err, "test case %d", i)
		} else {
			assert.Equal(t, tc.expected, n, "test case %d", i)
		}
	}
}

func TestMustParseI8(t *testing.T) {
	assert.Equal(t, int8(0), MustParseI8("0"))
	assert.Equal(t, int8(127), MustParseI8("127"))
	assert.Equal(t, int8(-128), MustParseI8("-128"))
	assert.Panics(t, func() {
		MustParseI8("128")
	})
	assert.Panics(t, func() {
		MustParseI8("-129")
	})
}

func TestParseU8(t *testing.T) {
	tests := []struct {
		input    string
		hasErr   bool
		expected uint8
	}{
		{"xyz", true, 0},
		{"0", false, 0},
		{"127", false, 127},
		{"255", false, 255},
		{"-1", true, 0},
		{"256", true, 0},
	}
	for i, tc := range tests {
		n, err := ParseU8(tc.input)
		if tc.hasErr {
			assert.NotNil(t, err, "test case %d", i)
		} else {
			assert.Equal(t, tc.expected, n, "test case %d", i)
		}
	}
}

func TestMustParseU8(t *testing.T) {
	assert.Equal(t, uint8(0), MustParseU8("0"))
	assert.Equal(t, uint8(127), MustParseU8("127"))
	assert.Equal(t, uint8(255), MustParseU8("255"))
	assert.Panics(t, func() {
		MustParseU8("-1")
	})
	assert.Panics(t, func() {
		MustParseU8("256")
	})
}

func TestParseI16(t *testing.T) {
	tests := []struct {
		input    string
		hasErr   bool
		expected int16
	}{
		{"xyz", true, 0},
		{"0", false, 0},
		{"32767", false, math.MaxInt16},
		{"-32768", false, math.MinInt16},
		{"32768", true, 0},
		{"-32769", true, 0},
	}
	for i, tc := range tests {
		n, err := ParseI16(tc.input)
		if tc.hasErr {
			assert.NotNil(t, err, "test case %d", i)
		} else {
			assert.Equal(t, tc.expected, n, "test case %d", i)
		}
	}
}

func TestMustParseI16(t *testing.T) {
	assert.Equal(t, int16(0), MustParseI16("0"))
	assert.Equal(t, int16(32767), MustParseI16("32767"))
	assert.Equal(t, int16(-32768), MustParseI16("-32768"))
	assert.Panics(t, func() {
		MustParseI16("-32769")
	})
	assert.Panics(t, func() {
		MustParseI16("32768")
	})
}

func TestParseU16(t *testing.T) {
	tests := []struct {
		input    string
		hasErr   bool
		expected uint16
	}{
		{"xyz", true, 0},
		{"0", false, 0},
		{"65535", false, math.MaxUint16},
		{"65536", true, 0},
		{"-1", true, 0},
	}
	for i, tc := range tests {
		n, err := ParseU16(tc.input)
		if tc.hasErr {
			assert.NotNil(t, err, "test case %d", i)
		} else {
			assert.Equal(t, tc.expected, n, "test case %d", i)
		}
	}
}

func TestMustParseU16(t *testing.T) {
	assert.Equal(t, uint16(0), MustParseU16("0"))
	assert.Equal(t, uint16(65535), MustParseU16("65535"))
	assert.Panics(t, func() {
		MustParseU16("-1")
	})
	assert.Panics(t, func() {
		MustParseU16("65536")
	})
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

func TestMustParseI32(t *testing.T) {
	assert.Equal(t, int32(0), MustParseI32("0"))
	assert.Equal(t, int32(2147483647), MustParseI32("2147483647"))
	assert.Equal(t, int32(-2147483648), MustParseI32("-2147483648"))
	assert.Panics(t, func() {
		MustParseI32("-2147483649")
	})
	assert.Panics(t, func() {
		MustParseI32("2147483648")
	})
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

func TestMustParseU32(t *testing.T) {
	assert.Equal(t, uint32(0), MustParseU32("0"))
	assert.Equal(t, uint32(4294967295), MustParseU32("4294967295"))
	assert.Panics(t, func() {
		MustParseU32("-1")
	})
	assert.Panics(t, func() {
		MustParseU32("4294967296")
	})
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

func TestMustParseI64(t *testing.T) {
	assert.Equal(t, int64(0), MustParseI64("0"))
	assert.Equal(t, int64(9223372036854775807), MustParseI64("9223372036854775807"))
	assert.Equal(t, int64(-9223372036854775808), MustParseI64("-9223372036854775808"))
	assert.Panics(t, func() {
		MustParseI64("-9223372036854775809")
	})
	assert.Panics(t, func() {
		MustParseI64("9223372036854775808")
	})
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

func TestMustParseU64(t *testing.T) {
	assert.Equal(t, uint64(0), MustParseU64("0"))
	assert.Equal(t, uint64(18446744073709551615), MustParseU64("18446744073709551615"))
	assert.Panics(t, func() {
		MustParseU64("-1")
	})
	assert.Panics(t, func() {
		MustParseU64("18446744073709551616")
	})
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

func TestMustParseF32(t *testing.T) {
	assert.Equal(t, float32(0), MustParseF32("0"))
	assert.Equal(t, float32(3.14), MustParseF32("3.14"))
}

func TestParseFloat64(t *testing.T) {
	tests := []struct {
		input    string
		hasErr   bool
		expected float64
	}{
		{"xyz", true, 0},
		{"0", false, 0},
		{"3.40282346638528859811704183484516925440e+38", false, math.MaxFloat32},
		{"1.401298464324817070923729583289916131280e-45", false, math.SmallestNonzeroFloat32},
	}
	for i, tc := range tests {
		n, err := ParseF64(tc.input)
		if tc.hasErr {
			assert.NotNil(t, err, "test case %d", i)
		} else {
			assert.Equal(t, tc.expected, n, "test case %d", i)
		}
	}
}

func TestMustParseF64(t *testing.T) {
	assert.Equal(t, float64(0), MustParseF64("0"))
	assert.Equal(t, float64(3.14), MustParseF64("3.14"))
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
