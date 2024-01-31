// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package util

import (
	"math/rand"
)

type Pair[T1, T2 any] struct {
	First  T1
	Second T2
}

func MakePair[T1, T2 any](a T1, b T2) Pair[T1, T2] {
	return Pair[T1, T2]{
		First:  a,
		Second: b,
	}
}

type Range struct {
	Min int
	Max int
}

func (r *Range) Mid() int {
	return (r.Min + r.Max) / 2
}

func (r *Range) Rand() int {
	if r.Min == r.Max {
		return r.Min
	}
	var val = rand.Intn(r.Max - r.Min + 1)
	return r.Min + val
}
