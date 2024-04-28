// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package zutil

import (
	"fmt"
	"image"
	"math"
	"reflect"
	"testing"
)

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
				t.Errorf("JSONStringify() = %v, want %v", got, tt.want)
			}
		})
	}
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
