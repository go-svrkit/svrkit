// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package collections

import (
	"cmp"
)

// Heapify establishes the heap invariants required by the other routines in this package.
// Init is idempotent with respect to the heap invariants
// and may be called whenever the heap invariants may have been invalidated.
func Heapify[S ~[]E, E cmp.Ordered](h S) {
	n := len(h)
	for i := n/2 - 1; i >= 0; i-- {
		heapDown(h, i, n)
	}
}

// HeapPush pushes the element x onto the heap.
func HeapPush[S ~[]E, E cmp.Ordered](h S, x E) S {
	h = append(h, x)
	heapUp(h, len(h)-1)
	return h
}

// HeapPop removes and returns the minimum element (according to Less) from the heap.
// Pop is equivalent to Remove(h, 0).
func HeapPop[S ~[]E, E cmp.Ordered](h S) (S, E) {
	n := len(h) - 1
	h[0], h[n] = h[n], h[0]
	heapDown(h, 0, n)
	x := h[len(h)-1]
	h = h[:len(h)-1]
	return h, x
}

// HeapRemove removes and returns the element at index i from the heap.
func HeapRemove[S ~[]E, E cmp.Ordered](h S, i int) (S, E) {
	n := len(h) - 1
	if n != i {
		h[i], h[n] = h[n], h[i]
		if !heapDown(h, i, n) {
			heapUp(h, i)
		}
	}
	x := h[len(h)-1]
	h = h[:len(h)-1]
	return h, x
}

// HeapFix re-establishes the heap ordering after the element at index i has changed its Value.
// Changing the Value of the element at index i and then calling Fix is equivalent to,
// but less expensive than, calling Remove(h, i) followed by a Push of the new Value.
// The complexity is O(log n) where n = h.Len().
func HeapFix[S ~[]E, E cmp.Ordered](h S, i int) {
	if !heapDown(h, i, len(h)) {
		heapUp(h, i)
	}
}

func heapUp[S ~[]E, E cmp.Ordered](h S, j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !(h[j] < h[i]) {
			break
		}
		h[j], h[i] = h[i], h[j] // swap
		j = i
	}
}

func heapDown[S ~[]E, E cmp.Ordered](h S, i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && h[j2] < h[j1] {
			j = j2 // = 2*i + 2  // right child
		}
		if !(h[j] < h[i]) {
			break
		}
		h[j], h[i] = h[i], h[j] // swap
		i = j
	}
	return i > i0
}
