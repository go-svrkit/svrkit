// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package util

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
