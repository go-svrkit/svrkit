// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package gutil

import (
	"fmt"
	"image"
	"math"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
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
