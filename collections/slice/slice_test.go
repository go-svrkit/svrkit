// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package slice

import (
	"testing"
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
