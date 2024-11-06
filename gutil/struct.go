// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package gutil

import (
	"math/rand"
	"strconv"
	"strings"
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

// In returns true if the value is in the range
func (r *Range) In(val int) bool {
	return r.Min <= val && val <= r.Max
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

func ParseRange(text, sep string) Range {
	var r Range
	if text == "" {
		return r
	}
	var parts = strings.Split(text, sep)
	if len(parts) != 2 {
		return r
	}
	minVal, _ := strconv.Atoi(parts[0])
	maxVal, _ := strconv.Atoi(parts[1])
	if minVal > maxVal {
		maxVal, minVal = minVal, maxVal
	}
	r.Min = minVal
	r.Max = maxVal
	return r
}
