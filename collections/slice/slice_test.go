// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package slice

import (
	"fmt"
	"math/rand"
	"slices"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_OrderedIndexOf(t *testing.T) {
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

func Test_Shrink(t *testing.T) {
	var a = Int32Slice{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var b = Shrink(a)
	if len(a) != len(b) {
		t.Fatalf("len(a) = %d, len(b) = %d", len(a), len(b))
	}
	if len(b) != cap(b) {
		t.Fatalf("len(b) = %d, cap(b) = %d", len(b), cap(b))
	}

	//
	{
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

}

func Test_Shuffle(t *testing.T) {
	var origin = []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var clone = slices.Clone(origin)
	Shuffle(clone)
	assert.Equal(t, len(clone), len(origin))
	slices.Sort(clone)
	assert.Equalf(t, origin, clone, "origin: %v", origin)
}

func Test_InsertAt(t *testing.T) {
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

func Test_RemoveAt(t *testing.T) {
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

func Test_RemoveFirst(t *testing.T) {
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

func Test_SortAndRemoveDup(t *testing.T) {
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

func Test_IsAllZeroElem(t *testing.T) {
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
			var output = IsAllZeroElem(tt.input)
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
			var output = IsAllZeroElem(tt.input)
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
			var output = IsAllZeroElem(tt.input)
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
			var output = IsAllZeroElem(tt.input)
			assert.Equalf(t, output, tt.want, "string case-%d", i+1)
		}
	}
}

func Test_OrderedInsert(t *testing.T) {
	var idSet []int

	for i := 50; i > 0; i-- {
		idSet = OrderedPutIfAbsent(idSet, i)
	}
	if !sort.IsSorted(IntSlice(idSet)) {
		t.Error("OrderedIDSet is not sorted")
	}
	for i := 1; i < 50; i++ {
		if !OrderedContains(idSet, i) {
			t.Error("OrderedIDSet does not contain", i)
		}
	}
	var n = len(idSet)
	for i := 1; i < 50; i++ {
		idSet = OrderedPutIfAbsent(idSet, i)
	}
	if len(idSet) != n {
		t.Error("OrderedIDSet does not contain duplicates")
	}
}

func Test_OrderedDelete(t *testing.T) {
	var idSet []int
	for i := 1; i <= 100; i++ {
		idSet = OrderedPutIfAbsent(idSet, i)
	}
	for i := 1; i <= 90; i++ {
		idSet = OrderedDelete(idSet, i)
	}
	if len(idSet) > 10 {
		t.Error("OrderedIDSet size error")
	}
	idSet = nil
	for i := 1; i <= 100; i++ {
		idSet = OrderedPutIfAbsent(idSet, i)
	}

	var deleted []int
	for i := 0; i < 80; i++ {
		var idx = rand.Int() % len(idSet)
		deleted = append(deleted, idSet[idx])
		idSet = slices.Delete(idSet, idx, idx+1)
	}
	for _, n := range deleted {
		if OrderedContains(idSet, n) {
			t.Error("OrderedIDSet does not contain", n)
		}
	}
}

func Test_OrderedContains(t *testing.T) {
	var idSet Int32Slice
	for i := 1; i <= 100; i++ {
		idSet = append(idSet, int32(i))
	}

	var deleted []int32
	for i := 0; i < 80; i++ {
		var idx = rand.Int() % len(idSet)
		deleted = append(deleted, idSet[idx])
		idSet = slices.Delete(idSet, idx, idx+1)
	}
	for _, n := range deleted {
		if OrderedContains(idSet, n) {
			t.Error("OrderedIDSet does not contain", n)
		}
	}
}

func Test_OrderedUnion(t *testing.T) {
	tests := []struct {
		input1   []int32
		input2   []int32
		expected []int32
	}{
		{[]int32{1, 2, 3}, nil, []int32{1, 2, 3}},
		{[]int32{1, 2, 3}, []int32{1, 2, 3}, []int32{1, 2, 3}},
		{[]int32{1, 2, 3, 4, 5}, []int32{4, 5, 6, 7, 8}, []int32{1, 2, 3, 4, 5, 6, 7, 8}},
		{[]int32{4, 5, 6, 7, 8}, []int32{1, 2, 3, 4, 5}, []int32{1, 2, 3, 4, 5, 6, 7, 8}},
		{[]int32{1, 3, 5, 7}, []int32{2, 4, 6, 8}, []int32{1, 2, 3, 4, 5, 6, 7, 8}},
		{[]int32{2, 4, 6, 8}, []int32{1, 3, 5, 7}, []int32{1, 2, 3, 4, 5, 6, 7, 8}},
	}
	for _, tc := range tests {
		var result = OrderedUnion(tc.input1, tc.input2)
		if len(result) != len(tc.expected) {
			t.Errorf("Expected %v, got %v", tc.expected, result)
		}
		for i := 0; i < len(result); i++ {
			if result[i] != tc.expected[i] {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		}
	}
}

func Test_OrderedIntersect(t *testing.T) {
	tests := []struct {
		input1   []int32
		input2   []int32
		expected []int32
	}{
		{[]int32{1, 2, 3}, nil, nil},
		{[]int32{1, 2, 3}, []int32{4, 5, 6}, nil},
		{[]int32{1, 2, 3, 4, 5}, []int32{4, 5, 6, 7, 8}, []int32{4, 5}},
		{[]int32{4, 5, 6, 7, 8}, []int32{1, 2, 3, 4, 5}, []int32{4, 5}},
	}
	for _, tc := range tests {
		var result = OrderedIntersect(tc.input1, tc.input2)
		if len(result) != len(tc.expected) {
			t.Errorf("Expected %v, got %v", tc.expected, result)
		}
		for i := 0; i < len(result); i++ {
			if result[i] != tc.expected[i] {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		}
	}
}
