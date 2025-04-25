// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package heap

import (
	"math/rand"
	"testing"
)

func heapVerify(t *testing.T, h []int, i int) {
	t.Helper()
	n := len(h)
	j1 := 2*i + 1
	j2 := 2*i + 2
	if j1 < n {
		if h[j1] < h[i] {
			t.Errorf("heap invariant invalidated [%d] = %d > [%d] = %d", i, h[i], j1, h[j1])
			return
		}
		heapVerify(t, h, j1)
	}
	if j2 < n {
		if h[j2] < h[i] {
			t.Errorf("heap invariant invalidated [%d] = %d > [%d] = %d", i, h[i], j1, h[j2])
			return
		}
		heapVerify(t, h, j2)
	}
}

func TestHeapify0(t *testing.T) {
	var h []int
	for i := 20; i > 0; i-- {
		h = append(h, 0) // all elements are the same
	}
	Heapify(h)
	heapVerify(t, h, 0)

	for i := 1; len(h) > 0; i++ {
		var x int
		h, x = HeapPop(h)
		heapVerify(t, h, 0)
		if x != 0 {
			t.Errorf("%d.th pop got %d; want %d", i, x, 0)
		}
	}
}

func TestHeapify1(t *testing.T) {
	var h []int
	for i := 20; i > 0; i-- {
		h = append(h, i) // all elements are different
	}
	Heapify(h)
	heapVerify(t, h, 0)

	for i := 1; len(h) > 0; i++ {
		var x int
		h, x = HeapPop(h)
		heapVerify(t, h, 0)
		if x != i {
			t.Errorf("%d.th pop got %d; want %d", i, x, i)
		}
	}
}

func Test(t *testing.T) {
	var h []int
	heapVerify(t, h, 0)

	for i := 20; i > 10; i-- {
		h = append(h, i)
	}
	Heapify(h)
	heapVerify(t, h, 0)

	for i := 10; i > 0; i-- {
		h = HeapPush(h, i)
		heapVerify(t, h, 0)
	}

	for i := 1; len(h) > 0; i++ {
		var x int
		h, x = HeapPop(h)
		if i < 20 {
			h = HeapPush(h, 20+i)
		}
		heapVerify(t, h, 0)
		if x != i {
			t.Errorf("%d.th pop got %d; want %d", i, x, i)
		}
	}
}

func TestRemove0(t *testing.T) {
	var h []int
	for i := 0; i < 10; i++ {
		h = append(h, i)
	}
	heapVerify(t, h, 0)

	for len(h) > 0 {
		var x int
		i := len(h) - 1
		h, x = HeapRemove(h, i)
		if x != i {
			t.Errorf("Remove(%d) got %d; want %d", i, x, i)
		}
		heapVerify(t, h, 0)
	}
}

func TestRemove1(t *testing.T) {
	var h []int
	for i := 0; i < 10; i++ {
		h = append(h, i)
	}
	heapVerify(t, h, 0)

	for i := 0; len(h) > 0; i++ {
		var x int
		h, x = HeapRemove(h, 0)
		if x != i {
			t.Errorf("Remove(0) got %d; want %d", x, i)
		}
		heapVerify(t, h, 0)
	}
}

func TestRemove2(t *testing.T) {
	N := 10

	var h []int
	for i := 0; i < N; i++ {
		h = append(h, i)
	}
	heapVerify(t, h, 0)

	m := make(map[int]bool)
	for len(h) > 0 {
		var x int
		h, x = HeapRemove(h, (len(h)-1)/2)
		m[x] = true
		heapVerify(t, h, 0)
	}

	if len(m) != N {
		t.Errorf("len(m) = %d; want %d", len(m), N)
	}
	for i := 0; i < len(m); i++ {
		if !m[i] {
			t.Errorf("m[%d] doesn't exist", i)
		}
	}
}

func BenchmarkDup(b *testing.B) {
	const n = 10000
	h := make([]int, 0, n)
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			h = append(h, 0) // all elements are the same
		}
		for len(h) > 0 {
			h, _ = HeapPop(h)
		}
	}
}

func TestFix(t *testing.T) {
	var h []int
	heapVerify(t, h, 0)

	for i := 200; i > 0; i -= 10 {
		h = HeapPush(h, i)
	}
	heapVerify(t, h, 0)

	if h[0] != 10 {
		t.Fatalf("Expected head to be 10, was %d", h[0])
	}
	h[0] = 210
	HeapFix(h, 0)
	heapVerify(t, h, 0)

	for i := 100; i > 0; i-- {
		elem := rand.Intn(len(h))
		if i&1 == 0 {
			h[elem] *= 2
		} else {
			h[elem] /= 2
		}
		HeapFix(h, elem)
		heapVerify(t, h, 0)
	}
}
