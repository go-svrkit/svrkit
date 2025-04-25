// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package gutil

import (
	"fmt"
	"image"
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZeroOf(t *testing.T) {
	assert.Equal(t, false, ZeroOf[bool]())
	assert.Equal(t, "", ZeroOf[string]())
	assert.Equal(t, 0, ZeroOf[int]())
	assert.Equal(t, uint(0), ZeroOf[uint]())
	assert.Equal(t, int8(0), ZeroOf[int8]())
	assert.Equal(t, uint8(0), ZeroOf[uint8]())
	assert.Equal(t, int16(0), ZeroOf[int16]())
	assert.Equal(t, uint16(0), ZeroOf[uint16]())
	assert.Equal(t, int32(0), ZeroOf[int32]())
	assert.Equal(t, uint32(0), ZeroOf[uint32]())
	assert.Equal(t, float32(0), ZeroOf[float32]())
	assert.Equal(t, float64(0), ZeroOf[float64]())
	assert.Equal(t, complex64(0), ZeroOf[complex64]())
	assert.Equal(t, image.Point{}, ZeroOf[image.Point]())
	assert.Equal(t, []int(nil), ZeroOf[[]int]())
	assert.Equal(t, map[int]int(nil), ZeroOf[map[int]int]())
}

func TestMD5Sum(t *testing.T) {
	tests := []struct {
		input []byte
		want  string
	}{
		{[]byte("hello"), "5d41402abc4b2a76b9719d911017c592"},
		{[]byte("world"), "7d793037a0760186574b0282f2f435e7"},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			if got := MD5Sum(tt.input); got != tt.want {
				t.Errorf("MD5Sum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSHA1Sum(t *testing.T) {
	tests := []struct {
		input []byte
		want  string
	}{
		{[]byte("hello"), "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d"},
		{[]byte("world"), "7c211433f02071597741e6ff5a8ea34789abbf43"},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			if got := SHA1Sum(tt.input); got != tt.want {
				t.Errorf("SHA1Sum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSHA256Sum(t *testing.T) {
	tests := []struct {
		input []byte
		want  string
	}{
		{[]byte("hello"), "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			if got := SHA256Sum(tt.input); got != tt.want {
				t.Errorf("SHA1Sum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsFileExist(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"", false},
		{"abcdefgxyz", false},
		{"./helper_test.go", true},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			if got := IsFileExist(tt.input); got != tt.want {
				t.Errorf("IsFileExist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadToLines(t *testing.T) {
	tests := []struct {
		input   string
		want    []string
		wantErr bool
	}{
		{"", []string(nil), false},
		{"a\nb\nc", []string{"a", "b", "c"}, false},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			var rd = strings.NewReader(tt.input)
			got, err := ReadToLines(rd)
			if err != nil {
				if tt.wantErr {
					return
				}
				t.Fatalf("ReadToLines() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadToLines() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONParse(t *testing.T) {
	tests := []struct {
		input  string
		want   interface{}
		hasErr bool
	}{
		{"xx", 0, true},
		{"1", 1, false},
		{"true", true, false},
		{"127", int8(127), false},
		{"65535", uint16(65535), false},
		{"2147483647", int32(math.MaxInt32), false},
		{"9223372036854775807", int64(math.MaxInt64), false},
		{"3.14", float32(3.14), false},
		{`{"X":12, "Y":34}`, image.Point{X: 12, Y: 34}, false},
		{`{"X":12,"Y":34}`, map[string]int{"X": 12, "Y": 34}, false},
		{"[4294967295,18446744073709551615]", []uint64{math.MaxUint32, math.MaxUint64}, false},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			var rval = reflect.New(reflect.TypeOf(tt.want))
			var err = JSONParse(tt.input, rval.Interface())
			if err != nil {
				if tt.hasErr {
					return
				}
				t.Fatalf("JSONParse() error = %v, wantErr %v", err, tt.hasErr)
			}
			var got = rval.Elem().Interface()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PrettyBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONStringify(t *testing.T) {
	tests := []struct {
		input interface{}
		want  string
	}{
		{1, "1"},
		{true, "true"},
		{math.MaxInt64, "9223372036854775807"},
		{image.Point{X: 12, Y: 34}, `{"X":12,"Y":34}`},
		{map[string]int{"X": 12, "Y": 34}, `{"X":12,"Y":34}`},
		{[]uint64{math.MaxUint32, math.MaxUint64}, `[4294967295,18446744073709551615]`},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			if got := JSONStringify(tt.input); got != tt.want {
				t.Errorf("case %d JSONStringify: %v, want %v", i+1, got, tt.want)
			}
		})
	}
}

func TestPrettyBytes(t *testing.T) {
	tests := []struct {
		input int64
		want  string
	}{
		{0, "0B"},
		{KiB, "1KiB"},
		{-KiB, "-1KiB"},
		{KiB + 100, "1.1KiB"},
		{MiB, "1MiB"},
		{MiB + 10*KiB, "1.01MiB"},
		{GiB, "1GiB"},
		{GiB + 100*MiB, "1.098GiB"},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			if got := PrettyBytes(tt.input); got != tt.want {
				t.Errorf("PrettyBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseByteCount(t *testing.T) {
	tests := []struct {
		input  string
		want   int64
		wantOK bool
	}{
		{"", 0, false},
		{"0", 0, true},
		{"0B", 0, true},
		{"64B", 64, true},
		{"1KiB", KiB, true},
		{"1MiB", MiB, true},
		{"1GiB", GiB, true},
		{"1TiB", TiB, true},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i)
		t.Run(name, func(t *testing.T) {
			got, ok := ParseByteCount(tt.input)
			assert.Equal(t, tt.wantOK, ok)
			if ok {
				assert.Equalf(t, tt.want, got, "ParseByteCount(%v)", tt.input)
			}
		})
	}
}
