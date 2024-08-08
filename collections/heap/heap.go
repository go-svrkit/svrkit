// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package heap

import (
	"cmp"
)

// Heapify establishes the heap invariants required by the other routines in this package.
// Init is idempotent with respect to the heap invariants
// and may be called whenever the heap invariants may have been invalidated.
func Heapify[S ~[]E, E cmp.Ordered](h S) {
	n := len(h)
	for i := n/2 - 1; i >= 0; i-- {
		down(h, i, n)
	}
}

// Push pushes the element x onto the heap.
func Push[S ~[]E, E cmp.Ordered](h S, x E) S {
	h = append(h, x)
	up(h, len(h)-1)
	return h
}

// Pop removes and returns the minimum element (according to Less) from the heap.
// Pop is equivalent to Remove(h, 0).
func Pop[S ~[]E, E cmp.Ordered](h S) (S, E) {
	n := len(h) - 1
	h[0], h[n] = h[n], h[0]
	down(h, 0, n)
	x := h[len(h)-1]
	h = h[:len(h)-1]
	return h, x
}

// Remove removes and returns the element at index i from the heap.
func Remove[S ~[]E, E cmp.Ordered](h S, i int) (S, E) {
	n := len(h) - 1
	if n != i {
		h[i], h[n] = h[n], h[i]
		if !down(h, i, n) {
			up(h, i)
		}
	}
	x := h[len(h)-1]
	h = h[:len(h)-1]
	return h, x
}

// Fix re-establishes the heap ordering after the element at index i has changed its value.
// Changing the value of the element at index i and then calling Fix is equivalent to,
// but less expensive than, calling Remove(h, i) followed by a Push of the new value.
// The complexity is O(log n) where n = h.Len().
func Fix[S ~[]E, E cmp.Ordered](h S, i int) {
	if !down(h, i, len(h)) {
		up(h, i)
	}
}

func up[S ~[]E, E cmp.Ordered](h S, j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !(h[j] < h[i]) {
			break
		}
		h[j], h[i] = h[i], h[j] // swap
		j = i
	}
}

func down[S ~[]E, E cmp.Ordered](h S, i0, n int) bool {
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
