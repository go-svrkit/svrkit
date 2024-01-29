// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package util

import (
	"cmp"
)

type Complex interface {
	~complex64 | ~complex128
}

// Comparable compares its two arguments for order. Returns a negative integer, zero, or a positive integer as the first argument is less than, equal to, or greater than the second.
// The implementor must ensure that signum(compare(x, y)) == -signum(compare(y, x)) for all x and y. (This implies that compare(x, y) must throw an exception if and only if compare(y, x) throws an exception.)
// The implementor must also ensure that the relation is transitive: ((compare(x, y)>0) && (compare(y, z)>0)) implies compare(x, z)>0.
// Finally, the implementor must ensure that compare(x, y)==0 implies that signum(compare(x, z))==signum(compare(y, z)) for all z.
type Comparable interface {
	// CompareTo returns an integer comparing two Comparables.
	// a.CompareTo(b) < 0 implies a < b
	// a.CompareTo(b) > 0 implies a > b
	// a.CompareTo(b) == 0 implies a == b
	CompareTo(b Comparable) int
}

// Comparator compares its two arguments for order. Returns a negative integer, zero, or a positive integer
// as the first argument is less than, equal to, or greater than the second. In the foregoing description,
// the notation sgn(expression) designates the mathematical signum function, which is defined to return
// one of -1, 0, or 1 according to whether the value of expression is negative, zero or positive.
//
// The implementor must ensure that sgn(compare(x, y)) == -sgn(compare(y, x)) for all x and y.
// (This implies that compare(x, y) must throw an exception if and only if compare(y, x) throws an exception.)
// The implementor must also ensure that the relation is transitive: ((compare(x, y)>0) && (compare(y, z)>0)) implies compare(x, z)>0.
// Finally, the implementor must ensure that compare(x, y)==0 implies that sgn(compare(x, z))==sgn(compare(y, z)) for all z.
// It is generally the case, but not strictly required that (compare(x, y)==0) == (x.equals(y)). Generally speaking,
// any comparator that violates this condition should clearly indicate this fact.
type Comparator[T any] func(a, b T) int

func OrderedCmp[T cmp.Ordered](a, b T) int {
	return cmp.Compare(a, b)
}

func Reversed[T any](cmp Comparator[T]) Comparator[T] {
	return func(a, b T) int {
		return -cmp(a, b)
	}
}

func BoolCmp(a, b bool) int {
	if a == b {
		return 0
	}
	if !a && b {
		return -1
	}
	return 1
}

func Complex64Cmp(a, b complex64) int {
	if a == b {
		return 0
	}
	if real(a) < real(a) {
		return -1
	}
	if real(a) == real(b) && imag(a) < imag(b) {
		return -1
	}
	return 1
}

func Complex128Cmp(a, b complex128) int {
	if a == b {
		return 0
	}
	if real(a) < real(a) {
		return -1
	}
	if real(a) == real(b) && imag(a) < imag(b) {
		return -1
	}
	return 1
}

func ZeroOf[T any]() T {
	var zero T
	return zero
}
