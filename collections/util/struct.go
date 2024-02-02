// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package util

import (
	"math/rand"
)

// Pair is a type that provides a way to store two heterogeneous objects as a single unit.
type Pair[T1, T2 any] struct {
	First  T1
	Second T2
}

// MakePair creates a Pair object, deducing the target type from the types of arguments
func MakePair[T1, T2 any](a T1, b T2) Pair[T1, T2] {
	return Pair[T1, T2]{
		First:  a,
		Second: b,
	}
}

// Range contains a min value and a max value
type Range struct {
	Min int
	Max int
}

// Mid returns the middle value of the range
func (r *Range) Mid() int {
	return (r.Min + r.Max) / 2
}

// Rand returns a random value in the range
func (r *Range) Rand() int {
	if r.Min == r.Max {
		return r.Min
	}
	var val = rand.Intn(r.Max - r.Min + 1)
	return r.Min + val
}
