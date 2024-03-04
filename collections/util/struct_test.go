// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package util

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakePair(t *testing.T) {
	tests := []struct {
		input1 int
		input  string
		want   Pair[int, string]
	}{
		{0, "", Pair[int, string]{0, ""}},
		{123, "456", Pair[int, string]{123, "456"}},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			var pair = MakePair(tt.input1, tt.input)
			if !reflect.DeepEqual(pair, tt.want) {
				t.Fatalf("MakePair() = %v, want %v", pair, tt.want)
			}
		})
	}
}

func TestRange_Mid(t *testing.T) {
	tests := []struct {
		min  int
		max  int
		want int
	}{
		{0, 0, 0},
		{1, 1, 1},
		{1, 3, 2},
		{1, 8, 4},
		{1, 7, 4},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			var ra = Range{Min: tt.min, Max: tt.max}
			var got = ra.Mid()
			if got != tt.want {
				t.Fatalf("Range.Mid() = %v, want %v", got, tt.want)
			}
		})
	}
}

// 结果是否符合随机区间
func expectRandomResult(t *testing.T, name string, min, max int, action func() int) {
	var total int
	var record = make(map[int]int) // 每个结果出现的次数
	for i := 0; i < 1000; i++ {
		var v = action()
		assert.LessOrEqual(t, v, max)
		assert.GreaterOrEqual(t, v, min)
		record[v]++
		total++
	}
	var count = max - min + 1
	for i := 0; i < count; i++ {
		var v = min + i
		assert.Greaterf(t, record[v], 0, name)
		assert.LessOrEqualf(t, record[v], 1000, name)
		delete(record, v)
	}
	assert.Equal(t, len(record), 0)
}

func TestRange_Rand(t *testing.T) {
	tests := []struct {
		min int
		max int
	}{
		{0, 0},
		{1, 1},
		{1, 3},
		{1, 8},
		{1, 7},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			var ra = Range{Min: tt.min, Max: tt.max}
			expectRandomResult(t, name, tt.min, tt.max, ra.Rand)
		})
	}
}
