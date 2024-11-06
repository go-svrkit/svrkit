// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package slice

import (
	"math/rand"
	"sort"
	"testing"

	"golang.org/x/exp/slices"
)

func TestOrderedInsert(t *testing.T) {
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

func TestOrderedDelete(t *testing.T) {
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

func TestOrderedContains(t *testing.T) {
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

func TestOrderedUnion(t *testing.T) {
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

func TestOrderedIntersect(t *testing.T) {
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
