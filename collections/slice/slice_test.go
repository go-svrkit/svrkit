// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package slice

import (
	"fmt"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlice_IndexOf(t *testing.T) {
	var a = []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for i := 0; i < len(a); i++ {
		var idx = OrderedIndexOf(a, a[i])
		if idx != i {
			t.Fatalf("a.IndexOf(%d) = %d, want %d", a[i], idx, i)
		}
		if idx >= 0 {
			if !OrderedContains(a, a[i]) {
				t.Fatalf("a.Contains(%d) = false, want true", a[i])
			}
		}
	}
}

func TestSlice_Shrink(t *testing.T) {
	var a = Int32Slice{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var b = Shrink(a)
	if len(a) != len(b) {
		t.Fatalf("len(a) = %d, len(b) = %d", len(a), len(b))
	}
	if len(b) != cap(b) {
		t.Fatalf("len(b) = %d, cap(b) = %d", len(b), cap(b))
	}
}

func TestSlice_Shuffle(t *testing.T) {
	var a = Int32Slice{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	Shuffle(a)
	t.Logf("%v", a)
}

func TestShrinkTypedSlice(t *testing.T) {
	var a, b []int64
	for i := 0; i < 10; i++ {
		a = append(a, int64(i))
	}
	b = Shrink(a)
	if len(b) != len(a) {
		t.Fatalf("len(b) = %d, want 10", len(b))
	}
	if len(b) != cap(b) {
		t.Fatalf("len(b) = %d, cap(b) = %d, want len(b) == cap(b)", len(b), cap(b))
	}
	for i := 0; i < len(b); i++ {
		if a[i] != b[i] {
			t.Fatalf("a[%d] = %d, b[%d] = %d, want a[%d] == b[%d]", i, a[i], i, b[i], i, i)
		}
	}
}

func TestSlice_InsertAt(t *testing.T) {
	tests := []struct {
		A        Int32Slice
		I        int
		N        int32
		expected Int32Slice
	}{
		{[]int32{}, 0, 1, []int32{1}},
		{[]int32{1}, -1, 2, []int32{1, 2}},
		{[]int32{1}, 2, 2, []int32{1, 2}},
		{[]int32{1, 3, 4}, 1, 2, []int32{1, 2, 3, 4}},
		{[]int32{2, 3, 4}, 0, 1, []int32{1, 2, 3, 4}},
		{[]int32{1, 2, 3}, 3, 4, []int32{1, 2, 3, 4}},
	}
	for i, tc := range tests {
		var output = InsertAt(tc.A, tc.I, tc.N)
		if !slices.Equal(output, tc.expected) {
			t.Fatalf("test %d not equal", i+1)
		}
	}
}

func TestSlice_RemoveAt(t *testing.T) {
	tests := []struct {
		A        Int32Slice
		I        int
		expected Int32Slice
	}{
		{[]int32{}, 0, []int32{}},
		{[]int32{1}, -1, []int32{1}},
		{[]int32{1}, 1, []int32{1}},
		{[]int32{1, 2, 3, 4, 5}, 4, []int32{1, 2, 3, 4}},
		{[]int32{1, 2, 3, 4, 5}, 0, []int32{5, 2, 3, 4}},
		{[]int32{1, 2, 3, 4, 5}, 2, []int32{1, 2, 5, 4}},
	}
	for i, tc := range tests {
		var output = RemoveAt(tc.A, tc.I)
		if !slices.Equal(output, tc.expected) {
			t.Fatalf("test %d not equal, expect %v, got %v", i+1, tc.expected, output)
		}
	}
}

func TestRemoveFirst(t *testing.T) {
	tests := []struct {
		A        Int32Slice
		I        int32
		expected Int32Slice
	}{
		{[]int32{}, 0, []int32{}},
		{[]int32{1}, 0, []int32{1}},
		{[]int32{1}, 1, []int32{}},
		{[]int32{1}, 1, []int32{}},
		{[]int32{1, 2, 3, 4, 5}, 2, []int32{1, 3, 4, 5}},
		{[]int32{1, 2, 2, 4, 5}, 2, []int32{1, 2, 4, 5}},
		{[]int32{1, 2, 3, 4, 5}, 5, []int32{1, 2, 3, 4}},
		{[]int32{1, 2, 3, 2, 1}, 1, []int32{2, 3, 2, 1}},
	}
	for i, tc := range tests {
		var output = RemoveFirst(tc.A, tc.I)
		if !slices.Equal(output, tc.expected) {
			t.Fatalf("test %d not equal, expect [%v], got %v", i+1, tc.expected, output)
		}
	}
}

func TestSortAndRemoveDup(t *testing.T) {
	tests := []struct {
		input []int
		want  []int
	}{
		{nil, nil},
		{[]int{}, []int{}},
		{[]int{1}, []int{1}},
		{[]int{1, 3, 2}, []int{1, 2, 3}},
		{[]int{1, 2, 2, 2, 3, 3}, []int{1, 2, 3}},
		{[]int{1, 2, 2, 2, 2, 3}, []int{1, 2, 3}},
		{[]int{1, 1, 3, 3, 2}, []int{1, 2, 3}},
		{[]int{1, 1, 3, 3, 2, 2}, []int{1, 2, 3}},
		{[]int{1, 3, 2, 1, 2, 3}, []int{1, 2, 3}},
	}

	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			out := SortAndRemoveDup(tt.input)
			assert.True(t, slices.Equal(out, tt.want))
		})
	}
}

func TestIsAllElemZero(t *testing.T) {
	{
		tests := []struct {
			input []int
			want  bool
		}{
			{nil, true},
			{[]int{}, true},
			{[]int{0, 0}, true},
			{[]int{0, 1, -2}, false},
			{[]int{0, 1, 0}, false},
		}
		for i, tt := range tests {
			var output = IsAllElemZero(tt.input)
			assert.Equalf(t, output, tt.want, "int case-%d", i+1)
		}
	}
	{
		tests := []struct {
			input []int64
			want  bool
		}{
			{nil, true},
			{[]int64{}, true},
			{[]int64{0, 0}, true},
			{[]int64{0, 1, -2}, false},
			{[]int64{0, 1, 0}, false},
		}
		for i, tt := range tests {
			var output = IsAllElemZero(tt.input)
			assert.Equalf(t, output, tt.want, "int64 case-%d", i+1)
		}
	}
	{
		tests := []struct {
			input []bool
			want  bool
		}{
			{nil, true},
			{[]bool{}, true},
			{[]bool{false, false}, true},
			{[]bool{false, true, false}, false},
		}
		for i, tt := range tests {
			var output = IsAllElemZero(tt.input)
			assert.Equalf(t, output, tt.want, "bool case-%d", i+1)
		}
	}
	{
		tests := []struct {
			input []string
			want  bool
		}{
			{nil, true},
			{[]string{}, true},
			{[]string{"", ""}, true},
			{[]string{"", "123", ""}, false},
		}
		for i, tt := range tests {
			var output = IsAllElemZero(tt.input)
			assert.Equalf(t, output, tt.want, "string case-%d", i+1)
		}
	}
}
